package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const PORT = ":8080"
const TOTAL_IMAGES = 15
const SAFE_MODE = true

var reqCount int

// Health check
func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health: OK")
}

func handleImageLoad(w http.ResponseWriter, r *http.Request) {
	reqCount++

	// URL path values
	imgID := r.PathValue("id")
	res := r.PathValue("res")
	filePath := filepath.Join("cmd", "assets", fmt.Sprintf("%s-res", res), fmt.Sprintf("img_%s.png", imgID))

	// Read img
	buf, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("ERROR loading img_%s.png: %v\n", imgID, err)
	}

	//Headers
	w.Header().Set("Content-Type", "image/png")

	//Response
	w.Write(buf)

	//Logs
	if SAFE_MODE {
		log.Printf("Addr: %s, ReqCount: %d, Requested [%s-res img_id: %s]", r.RemoteAddr, reqCount, res, imgID)
	} else {
		log.Printf("ReqCount: %d, Requested [%s-res img_id: %s]", reqCount, res, imgID)
	}
}

func main() {
	fs := http.FileServer(http.Dir("ui/dist"))
	http.Handle("/", fs)

	// http.HandleFunc("/api/images/all", handleImageLoad)
	http.HandleFunc("/api/images/{res}/{id}", handleImageLoad)
	http.HandleFunc("/api/health", healthCheck)

	log.Println("Stating server at PORT", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
