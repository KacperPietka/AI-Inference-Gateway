# AI Inference Gateway

A self-hosted API gateway written in Go that sits in front of AI models. It centralizes control over model access, handles intelligent routing, rate limiting, caching, and observability in one place.

## The Problem

Without a gateway, every application talks directly to model providers:

```
App A -> OpenAI API
App B -> Gemini API
App C -> Ollama
```

Problems:
- Every app reimplements the same logic
- No central cost visibility
- No rate limiting or quota management
- No intelligent routing between providers

With the gateway:

```
App -> Gateway -> Ollama (local, fast, free)
               -> Gemini API (cloud, capable)
```

---

## Features

- Intelligent model routing — short prompts go to Ollama, long or code prompts go to Gemini
- Multi-provider support — Ollama and Gemini API with a unified interface
- Per-user rate limiting backed by Redis
- Response caching — same prompt returns instantly without calling the model
- Structured JSON logging with request ID tracing
- Health checks for all dependencies
- Graceful shutdown
- Deployed on GKE with rolling updates and automated deploy script

---

## Architecture

```
Request
    |
    v
RequestID Middleware
    |
    v
Logger Middleware
    |
    v
Timeout Middleware
    |
    v
RateLimit Middleware (Redis)
    |
    v
Cache Check (Redis)
    |
    +--> HIT  -> return instantly (X-Cache: HIT)
    |
    +--> MISS -> Router
                    |
                    +--> short prompt  -> Ollama (tinyllama)
                    +--> long prompt   -> Gemini (gemini-2.0-flash)
                    +--> code prompt   -> Gemini (gemini-2.0-flash)
                    +--> fallback      -> Ollama
```

### Project Structure

```
inference-gateway/
├── cache/          <- cache interface, key generation, Redis implementation
├── config/         <- environment-based configuration
├── errors/         <- custom error types and sentinels
├── handlers/       <- HTTP handlers (generate, health, models, version)
├── interfaces/     <- shared interfaces (ModelProvider)
├── middleware/     <- logging, rate limiting, timeout, request ID
├── models/         <- Ollama and Gemini clients
├── ratelimit/      <- Redis rate limiter
├── router/         <- model routing rules and logic
├── types/          <- shared types and context keys
├── k8s/            <- Kubernetes manifests
├── Dockerfile
├── docker-compose.yml
└── deploy.sh
```

---

## Getting Started

### Prerequisites

- Go 1.26+
- Docker and Docker Compose
- Ollama with at least one model pulled
- Gemini API key (optional, from aistudio.google.com)

### Run locally

```bash
docker run -d -p 6379:6379 --name redis redis:7-alpine
ollama serve
go run main.go
```

With Gemini routing enabled:
```bash
GEMINI_API_KEY=your-key go run main.go
```

### Run with Docker Compose

```bash
docker compose up --build
```

---

## Routing Logic

The gateway routes requests automatically based on prompt characteristics:

| Condition | Provider | Model |
|---|---|---|
| Contains code keywords | Gemini | gemini-2.0-flash |
| Prompt >= 50 characters | Gemini | gemini-2.0-flash |
| Prompt < 50 characters | Ollama | tinyllama |
| Fallback | Ollama | tinyllama |

Users can always override routing by specifying a model explicitly in the request.

If `GEMINI_API_KEY` is not set, all requests route to Ollama.

---

## API

### POST /generate

```json
{
  "prompt": "Explain Kubernetes in one sentence",
  "model": "tinyllama"
}
```

Response:
```json
{
  "response": "Kubernetes is...",
  "model": "tinyllama",
  "cached": false
}
```

Response headers:
```
X-Cache: MISS
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 9
X-Request-ID: req-a3f9b2
```

### GET /health

```json
{
  "status": "healthy",
  "uptime": "5m32s",
  "ollama": { "status": "reachable", "model": "tinyllama" },
  "redis": { "status": "reachable" }
}
```

### GET /models

```json
{
  "models": [
    { "name": "tinyllama", "size": 637697920 }
  ]
}
```

### GET /version

```json
{
  "version": "v0.1.0",
  "go_version": "go1.26.1",
  "built_at": "2026-04-27T09:00:00Z",
  "cache_ttl": "3600s"
}
```

---

## Deployment

```bash
./deploy.sh
```

Builds the image tagged with the git commit hash, pushes to Google Artifact Registry, and performs a rolling update on GKE with a health check verification at the end.

Scale down when not in use:
```bash
kubectl scale deployment gateway --replicas=0 -n inference-gateway
kubectl scale deployment redis --replicas=0 -n inference-gateway
```

---

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.26 |
| Local model provider | Ollama |
| Cloud model provider | Gemini API |
| Cache and rate limiting | Redis |
| Orchestration | Kubernetes / GKE |
| Image registry | Google Artifact Registry |
| Logging | Go slog (structured JSON) |