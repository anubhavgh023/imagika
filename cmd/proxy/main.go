// PROXY SERVEr

package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/time/rate"
)

// var cache = make(map[string][]byte) // {img_id : img_data}
var cache = NewLRUCache(30) // LRU Cache
var mu sync.RWMutex
var reqCount int

// var ORIGIN_SERVER_URL = os.Getenv("REMOTE_HOST_ADDR")
const (
	PROXY_SERVER_PORT = ":9090"
	ORIGIN_SERVER_URL = "http://localhost:8080"
	// ORIGIN_SERVER_URL = "http://185.18.221.19:8080"
)

// Rate limiter Middleware
var clients = make(map[string]*rate.Limiter)

func rateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client't IP addr
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("ERROR rateLimiter:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		if _, ok := clients[ip]; !ok {
			clients[ip] = rate.NewLimiter(10, 30) // (r = rate ,b = tokens)
		}
		mu.Unlock()

		if !clients[ip].Allow() {
			log.Println("[RATE LIMITED]: Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	})
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Parse the oring server URL
	target, err := url.Parse(ORIGIN_SERVER_URL)
	if err != nil {
		log.Fatal("ERROR in parsing ORIGIN_SERVER_URL")
	}
	log.Printf("ReqCount: %d, CacheLen: %d, Incomming req to proxy: %v", reqCount, len(cache.cacheMap), r.URL.Path)

	//[CACHE] Add r.URL.Path to cache if empty
	// Using readers-writer lock:
	// The readers–writer lock allows multiple concurrent goroutines
	// to execute the read-only critical section part.
	mu.RLock()
	if v, ok := cache.Get(r.URL.Path); ok {
		// write the v stored in cache to w
		mu.RUnlock()
		log.Printf("[CACHED] Response: %s\n", r.URL.Path)
		if _, err := w.Write(v); err != nil {
			log.Println("ERROR copying resp.Body from:", err)
		}
		return
	}
	mu.RUnlock()

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
	// http.HandleFunc("/", proxyHandler)
	log.Println("Starting [PROXY SERVER] on", PROXY_SERVER_PORT)
	if err := http.ListenAndServe(PROXY_SERVER_PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
