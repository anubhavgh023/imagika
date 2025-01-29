package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Image Packet
type ImagePacket struct {
	head int
	buf  []byte
}

const PORT = ":8080"
const TOTAL_IMAGES = 15

var reqCount int

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
	log.Printf("Addr: %s, ReqCount: %d, Requested [%s-res img_id: %s]", r.RemoteAddr, reqCount, res, imgID)
}

func handleBatchImageLoad(w http.ResponseWriter, r *http.Request) {
	// Set Headers
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "application/octet-stream")

	reqCount++

	for i := 1; i <= TOTAL_IMAGES; i++ {
		// Contruct filePath
		filePath := filepath.Join("cmd", "assets", "low-res", fmt.Sprintf("img_%s.png", strconv.Itoa(i)))
		fmt.Println(filePath)

		// Open the img file
		file, err := os.Open(filePath)
		buf, _ := os.ReadFile(filePath)
		fmt.Println(len(buf))
		fmt.Println(buf[:20])
		fmt.Println(buf[(len(buf) - 20):])
		if err != nil {
			log.Printf("ERROR opening img_%d: %v\n", i, err)
			continue
		}

		// Get file info to determine size
		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("ERROR getting file info for img_%d.png: %v\n", i, err)
			continue
		}

		// Write the img size(4 byte header)
		size := uint32(fileInfo.Size())
		if err := binary.Write(w, binary.BigEndian, size); err != nil {
			log.Printf("ERROR writing size for img_%d.png: %v\n", i, err)
			continue
		}

		// Stream the img data to client
		_, err = io.Copy(w, file)
		file.Close()
		if err != nil {
			log.Printf("ERROR streaming img_%d.png: %v\n", i, err)
			return
		}
	}

	//Logs
	log.Printf("Addr: %s, ReqCount: %d", r.RemoteAddr, reqCount)
}

func main() {
	fs := http.FileServer(http.Dir("ui/dist"))
	http.Handle("/", fs)

	http.HandleFunc("/api/images/all", handleBatchImageLoad)
	http.HandleFunc("/api/images/{res}/{id}", handleImageLoad)

	log.Println("Stating server at PORT", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
