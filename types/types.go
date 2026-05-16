package types

type ContextKey string

const RequestIDKey ContextKey = "request_id"
const ModelKey ContextKey = "model"

type VersionResponse struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	BuiltAt   string `json:"built_at"`
	CacheTTL  string `json:"cache_ttl"`
}

// API types
type GenerateRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Cached   bool   `json:"cached"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// Ollama types

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
}

// Models endpoint types

type OllamaModel struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type OllamaModelsResponse struct {
	Models []OllamaModel `json:"models"`
}

type ModelsReponse struct {
	Models []OllamaModel `json:"models"`
}

// Health check types

type OllamaHealth struct {
	Status string `json:"status"`
	Model  string `json:"model"`
}

type RedisHealth struct {
	Status string `json:"status"`
}

type HealthResponse struct {
	Status string       `json:"status"`
	Uptime string       `json:"uptime"`
	Ollama OllamaHealth `json:"ollama"`
	Redis  RedisHealth  `json:"redis"`
}

// Gemini request types
type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// Gemini response types
type GeminiResponsePart struct {
	Text string `json:"text"`
}

type GeminiResponseContent struct {
	Parts []GeminiResponsePart `json:"parts"`
}

type GeminiCandidate struct {
	Content GeminiResponseContent `json:"content"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}
