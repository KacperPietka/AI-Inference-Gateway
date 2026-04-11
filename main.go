package main

import (
	"encoding/json" // for turning GO structs/maps into JSON
	"log"           //Package for printing to terminal with timestamps
	"net/http"      // GO's built-in HTTP server (very powerful)
)

// This defines the shape of the incoming request
type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

// This defines the shape of the response
type GenerateResponse struct {
	Prompt string `json:"prompt"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code})
}

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

func generateHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body into our struct
	var req GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that the prompt is not empty
	if req.Prompt == "" {
		writeError(w, "prompt is requeired", http.StatusBadRequest)
		return
	}

	// echo the request back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GenerateResponse{Prompt: req.Prompt})
}

func main() {
	// Registers the handler ("when someone hits /health, call healthHandler!")
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/generate", generateHandler)

	log.Println("Server starting on port 8080...")

	//Starts the server on port 8080 (the nil means "use default router" which in this case is /health as set above)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
