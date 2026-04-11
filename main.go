package main

import (
	"encoding/json" // for turning GO structs/maps into JSON
	"log"           //Package for printing to terminal with timestamps
	"net/http"      // GO's built-in HTTP server (very powerful)
)

// This is a HTTP handler (always with 2 params)
// w -> the response writer (where you write the responses)
// r -> the incoming request (you read from this (headers, body, URL, etc.))
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Tells the client "what comes back is JSON" ALWAYS SET THIS!!
	w.Header().Set("Content-Type", "application/json")

	// This Sends HTTP status code 200
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}

	//Encodes the map as JSON and writes it directly into the response writer!
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Registers the handler ("when someone hits /health, call healthHandler!")
	http.HandleFunc("/health", healthHandler)

	log.Println("Server starting on port 8080...")

	//Starts the server on port 8080 (the nil means "use default router" which in this case is /health as set above)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
