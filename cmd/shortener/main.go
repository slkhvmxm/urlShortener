package main

import (
	"log"
	"net/http"

	"github.com/slkhvmxm/urlShortener/internal/handler"
)

func main() {
	shortener := handler.NewUrlShortener()
	mux := http.NewServeMux()
	mux.HandleFunc("/", shortener.MainHandler)

	log.Printf("Server statrting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
