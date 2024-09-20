package main

import (
	"log"
	"net/http"

	"github.com/go-chi/render"

	"github.com/mikandro/url_shortener/internal/config"
	"github.com/mikandro/url_shortener/internal/redis"
	"github.com/mikandro/url_shortener/internal/router"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis client
	redisClient := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer redisClient.Close() // Ensure the client is closed when the app exits

	// Set up the router
	r := router.NewRouter(redisClient)

	// Start the server
	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("url shortener."))
	})
}

type Url struct {
	Url string `json:"url"`
}

type UrlShortenRequest struct {
	*Url
}

type ShortenUrlHandler struct {
	RedisClient *redis.Client
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
