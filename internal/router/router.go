package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikandro/url_shortener/internal/handlers"
	"github.com/mikandro/url_shortener/internal/redis"
)

func NewRouter(redisClient *redis.Client) *chi.Mux {
	r := chi.NewRouter()

	urlHandler := &handlers.UrlHandler{
		RedisClient: redisClient,
	}

	r.Post("/url", urlHandler.ShortenUrl)
	return r
}
