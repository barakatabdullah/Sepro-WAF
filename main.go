package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Create a reverse proxy
	targetURL, err := url.Parse("http://localhost:5173") // Replace with your target URL
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Start the server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Println("Reverse proxy server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
