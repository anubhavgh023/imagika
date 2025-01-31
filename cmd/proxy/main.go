package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	ORIGIN_SERVER = "http://185.18.221.19:8080"
	PORT          = ":8089"
)

var reqCount int

func proxyPass(w http.ResponseWriter, r *http.Request) {
	reqCount++
	body, _ := io.ReadAll(r.Body)
	data := string(body)

	newReq, _ := http.NewRequest(r.Method, r.URL.String(), strings.NewReader(data))

	url, _ := url.Parse(ORIGIN_SERVER)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, newReq)

	//Logs
	log.Printf("[PROXY SERVER]> Addr: %s, ReqCount: %d", r.RemoteAddr, reqCount)
}

func main() {
	http.HandleFunc("/", proxyPass)

	log.Printf("Starting proxy server at %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
