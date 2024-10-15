package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	my_redis "github.com/mikandro/url_shortener/internal/redis"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() (*my_redis.Client, func()) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := my_redis.NewClient(s.Addr(), "", 1)

	return client, func() {
		client.Close()
		s.Close()
	}
}

func TestShortenUrl(t *testing.T) {
	redisClient, teardown := setupTestRedis()
	defer teardown()

	handler := &UrlHandler{RedisClient: redisClient}

	reqBody := `{"url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ShortenUrl(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var responseBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&responseBody)
	assert.NoError(t, err)

	shortURL, ok := responseBody["short_url"]
	assert.True(t, ok)
	assert.NotEmpty(t, shortURL)
}

func TestRedirect(t *testing.T) {
	redisClient, teardown := setupTestRedis()
	defer teardown()

	handler := &UrlHandler{RedisClient: redisClient}

	shortCode := "abcd1234"
	longURL := "https://example.com"
	ctx := context.Background()
	err := redisClient.RedisClient.Set(ctx, shortCode, longURL, 0).Err()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	w := httptest.NewRecorder()

	handler.RedirectUrl(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Check that the response is a redirect
	assert.Equal(t, http.StatusFound, res.StatusCode)
}
