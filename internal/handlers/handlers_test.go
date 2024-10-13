package handlers

import (
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
}
