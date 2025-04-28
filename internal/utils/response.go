package utils

import (
	"encoding/json"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("failed to write JSON error: %v", err)
	}
}
