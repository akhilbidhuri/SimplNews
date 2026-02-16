package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// parseLimit parses limit parameter with default and max values
func parseLimit(limitStr string, defaultLimit, maxLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return defaultLimit
	}

	if limit > maxLimit {
		return maxLimit
	}

	return limit
}

// parseFloat parses a float parameter
func parseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}

	return val
}
