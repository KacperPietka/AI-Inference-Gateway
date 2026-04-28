package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	gwerrors "inference-gateway/errors"
	"inference-gateway/types"
)

var buildTime = time.Now().UTC().Format(time.RFC3339)

type VersionHandler struct {
	Version string
}

func NewVersionHandler(version string) *VersionHandler {
	return &VersionHandler{
		Version: version,
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
	})
}
