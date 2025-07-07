package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	sharedConfig "korun.io/shared/config"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(cfg *sharedConfig.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     100,
		PoolTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		DialTimeout:  time.Second * 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	slog.Info("Connected to Redis successfully")
	return &Client{rdb: rdb}, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to get key from Redis: %w", err)
	}
	return val, nil
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := c.rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %w", err)
	}
	return nil
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	err := c.rdb.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete keys from Redis: %w", err)
	}
	return nil
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key in Redis: %w", err)
	}
	return val, nil
}

func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	err := c.rdb.SAdd(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to add members to set in Redis: %w", err)
	}
	return nil
}

func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := c.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get members from set in Redis: %w", err)
	}
	return members, nil
}

func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	err := c.rdb.SRem(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to remove members from set in Redis: %w", err)
	}
	return nil
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.rdb.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key in Redis: %w", err)
	}
	return nil
}

func (c *Client) Close() error {
	if err := c.rdb.Close(); err != nil {
		slog.Error("Failed to close Redis connection", "error", err)
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	slog.Info("Redis connection closed successfully")
	return nil
}
