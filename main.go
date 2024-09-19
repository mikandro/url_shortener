package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("url shortener."))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Route("/url", func(r chi.Router) {
			r.Post("/shorten", shortenUrl)
		})
	})

	http.ListenAndServe(":3000", r)
}

type Url struct {
	Url string `json:"url"`
}

type UrlShortenRequest struct {
	*Url
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	data := &UrlShortenRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
}

func (u *UrlShortenRequest) Bind(r *http.Request) error {
	// error to avoid a nil pointer dereference.
	if u.Url == nil {
		return errors.New("missing required Url fields")
	}
	return nil
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
