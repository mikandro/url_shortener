package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/mikandro/url_shortener/internal/handlers"
	my_redis "github.com/mikandro/url_shortener/internal/redis"
)

func NewRouter(redisClient *my_redis.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"}, // Allow requests from frontend
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value for preflight requests cache
	}))
	urlHandler := &handlers.UrlHandler{
		RedisClient: redisClient,
	}

	r.Post("/url", urlHandler.ShortenUrl)
	r.Get("/{short_url}", urlHandler.RedirectUrl)
	r.Delete("/{short_url}", urlHandler.DeleteUrl)
	return r
}
