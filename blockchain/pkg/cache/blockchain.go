package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Blockchain interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type BlockchainCache struct {
	Expiration time.Duration
	redisCli   *redis.Client
}

func NewBlockchainCache(redisCli *redis.Client, expiration time.Duration) Blockchain {
	return &BlockchainCache{
		redisCli:   redisCli,
		Expiration: expiration,
	}
}

func (b *BlockchainCache) Get(ctx context.Context, key string) (string, error) {
	value := b.redisCli.Get(ctx, key).Val()

	if value == "" {
		return "", nil
	}

	return value, nil
}

func (b *BlockchainCache) Set(ctx context.Context, key string, value string) error {
	return b.redisCli.Set(ctx, key, value, b.Expiration).Err()
}
