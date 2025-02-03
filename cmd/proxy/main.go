package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"unsafe"

	"golang.org/x/time/rate"
)

// var cache = make(map[string][]byte) // {img_id : img_data}
var cache = NewLRUCache(10) // LRU Cache
var mu sync.Mutex
var reqCount int

// var ORIGIN_SERVER_URL = os.Getenv("REMOTE_HOST_ADDR")
const (
	PROXY_SERVER_PORT = ":9090"
	ORIGIN_SERVER_URL = "http://localhost:8080"
)

// Rate limiter
func rateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client't IP addr
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("ERROR rateLimiter:", err)
			return
		}

		mu.Lock()
		if _, ok := clients[ip]; !ok {
			clients[ip] = rate.NewLimiter(30, 60)
			fmt.Println(clients[ip])
		}

		if !clients[ip].Allow() {
			mu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		mu.Unlock()
		next(w, r)
	})
}

func getCacheMemoryUsage() float64 {
	var totalSize int64

	for key, value := range cache.cacheMap {
		totalSize += int64(len(key)) + int64(unsafe.Sizeof(key))           // Key string size
		totalSize += int64(len(value.value)) + int64(unsafe.Sizeof(value)) // Byte slice size
	}
	return float64(totalSize) / (1024 * 1024) // Convert to MB
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	reqCount++
	fmt.Println("-------------------------")
	fmt.Println("CACHE len:", len(cache.cacheMap))
	fmt.Printf("CACHE Memory Usage: %.2f MB\n", getCacheMemoryUsage())
	fmt.Println("-------------------------")
	mu.Unlock()

	// Step 1: Parse the oring server URL
	target, err := url.Parse(ORIGIN_SERVER_URL)
	if err != nil {
		log.Fatal("ERROR in parsing ORIGIN_SERVER_URL")
	}
	log.Printf("Req Count: %d, Incomming req to proxy:%v %v", reqCount, r.Method, r.URL.Path)

	//[CACHE] Add r.URL.Path to cache if empty
	if v, ok := cache.Get(r.URL.Path); ok {
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
	// log.Println("Forwarding request to origin server:", newReq.Method, newReq.URL.String())

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
	cache.Put(r.URL.Path, body)
	mu.Unlock()

	// Step 6: Cody headers from the ORIGIN_SERVER
	// log.Println("Received resp from ORIGIN_SERVER:", resp.Status)
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
	http.Handle("/", rateLimiter(proxyHandler))
	log.Println("Starting [PROXY SERVER] on", PROXY_SERVER_PORT)
	if err := http.ListenAndServe(PROXY_SERVER_PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
