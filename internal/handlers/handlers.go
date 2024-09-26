package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mikandro/url_shortener/internal/redis"
	"github.com/mikandro/url_shortener/internal/shortener"
)

type UrlHandler struct {
	RedisClient *redis.Client
}

type Url struct {
	Url string `json:"url"`
}

func (h *UrlHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	var longUrl Url
	if err := json.NewDecoder(r.Body).Decode(&longUrl); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	shortCode := shortener.GenerateShortCode(longUrl.Url)

	// Store the url in Redis (for example, as a JSON string)
	err := h.RedisClient.RedisClient.Set(ctx, longUrl.Url, shortCode, 0).Err()
	if err != nil {
		log.Printf("Error saving url in db %e", err)
		http.Error(w, "Could not save url", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(longUrl)
}
