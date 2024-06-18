package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisClient(cfg Config) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	client := redis.NewClient(opts)
	return client, nil
}
