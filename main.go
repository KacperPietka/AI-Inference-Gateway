package main

import (
	"bytes"
	"encoding/json" // for turning GO structs/maps into JSON
	"log"           //Package for printing to terminal with timestamps
	"net/http"      // GO's built-in HTTP server (very powerful)
	"time"
)

// This defines the shape of the incoming request
type GenerateRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

// This defines the shape of the response
type GenerateResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Cached   bool   `json:"cached"`
}

// This defines the shape of the error
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// This defines Ollama Request strcture
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// This defines the ollama response structure
type OllamaResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
}

// Config
const (
	ollamaURL    = "http://localhost:11434/api/generate"
	defaultModel = "mistral"
	serverPort   = ":8080"
)

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code})
}

func callOllama(prompt string, model string) (*OllamaResponse, error) {
	ollamaReq := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	// encode the request to JSON bytes
	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}

	// Create a HTTP client with a timeout to not let servers hang forever
	client := &http.Client{Timeout: 60 * time.Second}

	// Make the POST request to Ollama
	resp, err := client.Post(ollamaURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // close the body when done

	// decode Ollama's response into Go's struct
	var ollamaResp OllamaResponse
	err = json.NewDecoder(resp.Body).Decode(&ollamaResp)
	if err != nil {
		return nil, err
	}

	return &ollamaResp, nil
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

	// Use default model if none specified
	if req.Model == "" {
		req.Model = defaultModel
	}

	// Call Ollama
	log.Printf("calling ollama with model=%s prompt=%q", req.Model, req.Prompt)
	ollamaResp, err := callOllama(req.Prompt, req.Model)
	if err != nil {
		log.Printf("ollama error: %v", err)
		writeError(w, "failed to call model", http.StatusInternalServerError)
		return
	}

	// Return clean response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GenerateResponse{
		Response: ollamaResp.Response,
		Model:    ollamaResp.Model,
		Cached:   false, // this will change when we apply caching
	})
}

func main() {
	// Registers the handler ("when someone hits /health, call healthHandler!")
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/generate", generateHandler)

	log.Println("Server starting on port 8080...")

	//Starts the server on port 8080 (the nil means "use default router" which in this case is /health as set above)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
