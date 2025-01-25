package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

const PORT = ":8080"

var db = make(map[string]string) // {remoteAddr : uuid}
var mu = sync.Mutex{}
var reqCount int

func home(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	id, _ := exec.Command("uuidgen").Output()
	remoteAddr := strings.Split(r.RemoteAddr, ":")
	if _, ok := db[r.RemoteAddr]; !ok {
		db[remoteAddr[0]] = string(id)
	} else {
		fmt.Fprintf(w, "User already present with id: %s", string(id))
		return
	}
	reqCount++
	mu.Unlock()

	//Logs
	log.Printf("Request No: %d, In Memory DB Entries: %d\n", reqCount, len(db))

	//Headers
	w.Header().Add("Server", "Go")
	w.Header().Add("Set-cookie", string(id))

	//Responses
	fmt.Fprintf(w, "Health: OK\n")
	fmt.Fprintf(w, "In Memory DB -> Entries: %d\n", len(db))
}

func dbStats(w http.ResponseWriter, r *http.Request) {
	for k, v := range db {
		fmt.Fprintf(w, "id: %v\naddr: %v\n\n", v, k)
	}
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Println("User Info Req for id:", id)
	fmt.Fprintf(w, "addr: %s\n", db[id])
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /dbStats", dbStats)
	mux.HandleFunc("GET /userInfo/{id}", userInfo)

	log.Printf("Starting server on %v", PORT)
	err := http.ListenAndServe(PORT, mux)
	log.Fatal(err)
}
