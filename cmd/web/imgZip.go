package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Zip Implementaion:
// Issues: slow loads
func handleBatchImageLoad(w http.ResponseWriter, r *http.Request) {
	// Set Headers
	w.Header().Set("Content-Type", "application/zip")
	w.(http.Flusher).Flush()

	// Create a buffer to write archive to
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	reqCount++
	for i := 1; i <= TOTAL_IMAGES; i++ {
		// Contruct filePath
		fileName := fmt.Sprintf("img_%d.png", i)
		filePath := filepath.Join("cmd", "assets", "low-res", fileName)

		// Read img file
		imgData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("ERROR reading img_%d: %v\n", i, err)
			continue
		}

		// Create file in archive
		writer, err := zipWriter.Create(fileName)
		if err != nil {
			log.Printf("ERROR creating zip entry for %s: %v\n", fileName, err)
		}

		// Write img data to archive
		if _, err := writer.Write(imgData); err != nil {
			log.Printf("ERROR writing %s to zip: %v\n", fileName, err)
			continue
		}
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		http.Error(w, "Error creating zip", http.StatusInternalServerError)
		return
	}

	// Write the zip file to response
	if _, err := io.Copy(w, buf); err != nil {
		log.Printf("Error sending zip: %v", err)
	}

	//Logs
	log.Printf("req from: %v", r.RemoteAddr)
	log.Printf("[ORIGIN SERVER]> Addr: %s, ReqCount: %d", r.RemoteAddr, reqCount)
}
