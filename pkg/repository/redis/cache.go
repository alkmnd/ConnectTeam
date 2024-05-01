package redis

import (
	"ConnectTeam/pkg/repository/models"
)

type Cache interface {
	HSet(key string, notification models.Notification) error
	Get(key string) (models.Notification, error)
	HGet(key string) ([]models.Notification, error)
}
