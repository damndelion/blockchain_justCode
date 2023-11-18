package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(address string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
