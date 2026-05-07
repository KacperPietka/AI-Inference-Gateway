# AI Inference Gateway

A self-hosted API gateway written in Go that sits in front of AI models. It centralizes control over model access, handles rate limiting, caching, and observability in one place.

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
- No reliability guarantees

With the gateway:

```
App -> Gateway -> OpenAI / Gemini / Ollama
```

---

## Features

- Request forwarding to local and remote model providers
- Per-user rate limiting backed by Redis
- Response caching — same prompt returns instantly without calling the model
- Structured JSON logging with request ID tracing
- Health checks for all dependencies
- Graceful shutdown
- Deployed on GKE with rolling updates

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
RateLimit Middleware
    |
    v
Handler -> Model Provider
```

---

## Getting Started

### Prerequisites

- Go 1.26+
- Docker and Docker Compose
- Ollama with at least one model pulled

### Run locally

```bash
docker run -d -p 6379:6379 --name redis redis:7-alpine
ollama serve
go run main.go
```

### Run with Docker Compose

```bash
docker compose up --build
```

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
  "built_at": "2026-04-27T09:00:00Z"
}
```

---

## Deployment

```bash
./deploy.sh
```

Builds the image, pushes to Artifact Registry, and performs a rolling update on GKE.

---

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.26 |
| Model provider | Ollama |
| Cache and rate limiting | Redis |
| Orchestration | Kubernetes / GKE |
| Logging | Go slog |