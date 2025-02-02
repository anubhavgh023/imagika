package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"unsafe"
)

var cache = make(map[string][]byte) // {img_id : img_data}
var mu sync.Mutex

const (
	PROXY_SERVER_PORT = ":9090"
	// ORIGIN_SERVER_URL = "http://185.18.221.19:8080"
	ORIGIN_SERVER_URL = "http://localhost:8080"
)

func getCacheMemoryUsage() float64 {
	mu.Lock()
	defer mu.Unlock()

	var totalSize int64

	for key, value := range cache {
		totalSize += int64(len(key)) + int64(unsafe.Sizeof(key))     // Key string size
		totalSize += int64(len(value)) + int64(unsafe.Sizeof(value)) // Byte slice size
	}
	return float64(totalSize) / (1024 * 1024) // Convert to MB
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------------------")
	fmt.Println("CACHE len:", len(cache))
	fmt.Printf("CACHE Memory Usage: %.2f MB\n", getCacheMemoryUsage())
	fmt.Println("-------------------------")

	// Step 1: Parse the oring server URL
	target, err := url.Parse(ORIGIN_SERVER_URL)
	if err != nil {
		log.Fatal("ERROR in parsing ORIGIN_SERVER_URL")
	}
	log.Println("Incomming req to proxy:", r.Method, r.URL.Path)

	//[CACHE] Add r.URL.Path to cache if empty
	if v, ok := cache[r.URL.Path]; ok {
		// write the v stored in cache to w
		log.Printf("This response was CACHED: %s\n", r.URL.Path)
		if _, err := w.Write(v); err != nil {
			log.Println("ERROR copying resp.Body from:", err)
		}
		return
	}

	// Step 2: Reconstructing the URL for the origin server
	proxyURL := *r.URL
	proxyURL.Scheme = target.Scheme
	proxyURL.Host = target.Host

	// Step 3: Create a new HTTP req for origin server
	newReq, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	if err != nil {
		log.Println("ERROR creating req:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Copy Headers from original req to => newReq
	for key, values := range r.Header {
		for _, v := range values {
			newReq.Header.Add(key, v)
		}
	}

	// Modidy headers to indicate req in being proxied
	newReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	newReq.Host = target.Host
	log.Println("Forwarding request to origin server:", newReq.Method, newReq.URL.String())

	// Step 5: Send the req to origin server
	client := http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		log.Println("ERROR forwarding req:", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	//[CACHE] reading the body and caching
	body, err := io.ReadAll(resp.Body)
	mu.Lock()
	cache[r.URL.Path] = body
	mu.Unlock()

	// Step 6: Cody headers from the ORIGIN_SERVER
	log.Println("Received resp from ORIGIN_SERVER:", resp.Status)
	for key, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	// Step 7: Write body data to w
	// _, err = io.Copy(w, resp.Body)
	if _, err := w.Write(body); err != nil {
		log.Println("ERROR writing to w:", err)
	}
}

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Println("Starting [PROXY SERVER] on", PROXY_SERVER_PORT)
	if err := http.ListenAndServe(PROXY_SERVER_PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
