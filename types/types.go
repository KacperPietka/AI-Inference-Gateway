package types

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

type OllamaModelsReponse struct {
	Models []OllamaModel `json:"models"`
}

type ModelsReponse struct {
	Models []OllamaModel `json:"models"`
}
