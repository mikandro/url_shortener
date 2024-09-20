package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mikandro/url_shortener/internal/redis"
)

type UrlHandler struct {
	RedisClient *redis.Client
}

func (h *UrlHandler) AddShortUrl(w http.ResponseWriter, r *http.Request) {
	var shortUrl map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&shortUrl); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	// Store the url in Redis (for example, as a JSON string)
	err := h.RedisClient.RedisClient.Set(ctx, "url:1", shortUrl, 0).Err()
	if err != nil {
		http.Error(w, "Could not save url", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortUrl)
}
