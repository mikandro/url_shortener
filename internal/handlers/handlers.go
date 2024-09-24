package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mikandro/url_shortener/internal/redis"
)

type UrlHandler struct {
	RedisClient *redis.Client
}

type Url struct {
	Url string `json:"url"`
}

func (h *UrlHandler) AddShortUrl(w http.ResponseWriter, r *http.Request) {
	var longUrl Url
	if err := json.NewDecoder(r.Body).Decode(&longUrl); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	urlJSON, err := json.Marshal(longUrl)
	if err != nil {
		http.Error(w, "Error encoding article to JSON", http.StatusInternalServerError)
		return
	}

	// Store the url in Redis (for example, as a JSON string)
	err = h.RedisClient.RedisClient.Set(ctx, longUrl.Url, urlJSON, 0).Err()
	if err != nil {
		log.Printf("Error saving url in db %e", err)
		http.Error(w, "Could not save url", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(longUrl)
}
