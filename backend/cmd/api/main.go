package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("RSSY API Server")
	log.Println("Starting server on :8080")

	// Basic health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
