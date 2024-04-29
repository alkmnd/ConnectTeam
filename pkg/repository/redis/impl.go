package redis

import (
	"ConnectTeam/pkg/repository/models"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type NotificationCache struct {
	Cache      *redis.Client
	expiration time.Duration
}

func (n NotificationCache) HSet(key string, notification models.Notification) error {
	ctx := context.Background()
	err := n.Cache.HSet(ctx, key, notification, n.expiration).Err()
	return err
}

func (n NotificationCache) Get(key string) (models.Notification, error) {
	var notification models.Notification
	ctx := context.Background()
	err := n.Cache.Get(ctx, key).Scan(&notification)
	return notification, err
}

func NewNotificationCache(client *redis.Client) Cache {
	return &NotificationCache{
		Cache:      client,
		expiration: time.Hour * 24 * 7,
	}
}
