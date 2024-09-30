package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikandro/url_shortener/internal/handlers"
	my_redis "github.com/mikandro/url_shortener/internal/redis"
)

func NewRouter(redisClient *my_redis.Client) *chi.Mux {
	r := chi.NewRouter()

	urlHandler := &handlers.UrlHandler{
		RedisClient: redisClient,
	}

	r.Post("/url", urlHandler.ShortenUrl)
	r.Get("/{short_url}", urlHandler.RedirectUrl)
	return r
}
