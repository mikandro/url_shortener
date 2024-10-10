package handlers

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupTestRedis() (*redis.Client, func()) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	return client, func() {
		client.Close()
		s.Close()
	}
}
