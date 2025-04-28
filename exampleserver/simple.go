package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: Method=%s, URL=%s\n", r.Method, r.URL.String())

	w.Write([]byte(fmt.Sprintf("Request received: Method=%s, URL=%s\n", r.Method, r.URL.String())))
}

func main() {
	port := flag.String("port", "9000", "Port to listen on")
	flag.Parse()

	http.HandleFunc("/", handler)

	log.Printf("Starting server on port %s...\n", *port)

	if err := http.ListenAndServe(":"+(*port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
