package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"

	my_redis "github.com/mikandro/url_shortener/internal/redis"
	"github.com/mikandro/url_shortener/internal/shortener"
)

type UrlHandler struct {
	RedisClient *my_redis.Client
}

type ShortenUrlRequest struct {
	Url string `json:"url"`
}

type ShortenUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type RedirectUrlRequest struct {
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
	err := h.RedisClient.RedisClient.Set(ctx, shortCode, req.Url, 0).Err()
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

func (h *UrlHandler) RedirectUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "short_url")

	ctx := context.Background()
	longUrl, err := h.RedisClient.RedisClient.Get(ctx, shortUrl).Result()
	if err == redis.Nil { // TODO it should be redis.Nil
		// Short code does not exist in Redis
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		// An error occurred while accessing Redis
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longUrl, http.StatusFound)
}

func (h *UrlHandler) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "short_url")

	ctx := context.Background()

	res, err := h.RedisClient.RedisClient.Del(ctx, shortUrl).Result()

	if err == redis.Nil {
		http.Error(w, "Short url not found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if res == 0 {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"message": "Short url was deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
