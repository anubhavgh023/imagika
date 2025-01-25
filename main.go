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
const MODE = "dev" // dev || prod

// In-Memery DB
var db = make(map[string]string) // {remoteAddr : uuid}

var mu = sync.Mutex{}
var reqCount int

func home(w http.ResponseWriter, r *http.Request) {
	var remoteAddr string
	if MODE == "dev" {
		remoteAddr = r.RemoteAddr
	} else {
		remoteAddr = strings.Split(r.RemoteAddr, ":")[0]
	}

	if _, ok := db[r.RemoteAddr]; !ok {
		mu.Lock()
		id, _ := exec.Command("uuidgen").Output()
		db[remoteAddr] = string(id)
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
	} else {
		fmt.Fprintf(w, "User already present with id: %s", db[remoteAddr])
		return
	}
}

func dbStats(w http.ResponseWriter, r *http.Request) {
	for k, v := range db {
		fmt.Fprintf(w, "addr: %v \nid: %v \n\n", k, v)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /dbStats", dbStats)

	log.Printf("Starting server on %v", PORT)
	err := http.ListenAndServe(PORT, mux)
	log.Fatal(err)
}
