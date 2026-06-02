// main.go is the entry point and wires up routes and starts the HTTP server
package main

import (
	"log"
	"net/http"
)

func main() {
	initDB() // connect to Postgres before starting the server

	// mux is the router that maps incoming URL paths to handler functions
	mux := http.NewServeMux()

	// "METHOD /path" pattern 
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /products", getProducts)
	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("POST /products", createProduct)

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
