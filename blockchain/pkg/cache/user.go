package cache

import (
	"context"
	"encoding/json"
	"time"

	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/redis/go-redis/v9"
)

type User interface {
	Get(ctx context.Context, key string) (*userEntity.User, error)
	Set(ctx context.Context, key string, value *userEntity.User) error
}

type UserCache struct {
	Expiration time.Duration
	redisCli   *redis.Client
}

func NewUserCache(redisCli *redis.Client, expiration time.Duration) User {
	return &UserCache{
		redisCli:   redisCli,
		Expiration: expiration,
	}
}

func (c *UserCache) Get(ctx context.Context, key string) (*userEntity.User, error) {
	value := c.redisCli.Get(ctx, key).Val()

	if value == "" {
		return nil, nil
	}

	var user *userEntity.User
	err := json.Unmarshal([]byte(value), &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *UserCache) Set(ctx context.Context, key string, value *userEntity.User) error {
	userJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redisCli.Set(ctx, key, string(userJSON), c.Expiration).Err()
}
