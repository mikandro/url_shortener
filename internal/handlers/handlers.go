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

type ShortenUrlRequest struct {
	Url string `json:"url"`
}

type ShortenUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *UrlHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	var req ShortenUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	shortCode := shortener.GenerateShortCode(req.Url)

	// Store the url in Redis (for example, as a JSON string)
	err := h.RedisClient.RedisClient.Set(ctx, req.Url, shortCode, 0).Err()
	if err != nil {
		log.Printf("Error saving url in db %e", err)
		http.Error(w, "Could not save url", http.StatusInternalServerError)
		return
	}

	response := ShortenUrlResponse{
		ShortUrl: "http://localhost/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to encode the response"})
		return
	}
}
