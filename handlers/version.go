package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	gwerrors "inference-gateway/errors"
	"inference-gateway/types"
)

var buildTime = time.Now().UTC().Format(time.RFC3339)

type VersionHandler struct {
	Version     string
	CacheTTLSec int
}

func NewVersionHandler(version string, cacheTTLSec int) *VersionHandler {
	return &VersionHandler{
		Version:     version,
		CacheTTLSec: cacheTTLSec,
	}
}

func (h *VersionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, gwerrors.ErrMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.VersionResponse{
		Version:   h.Version,
		GoVersion: runtime.Version(),
		BuiltAt:   buildTime,
		CacheTTL:  fmt.Sprintf("%ds", h.CacheTTLSec),
	})
}
