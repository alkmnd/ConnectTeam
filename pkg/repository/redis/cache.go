package redis

import (
	"ConnectTeam/pkg/repository/models"
	"github.com/google/uuid"
)

type Cache interface {
	HSet(key string, value uuid.UUID, notification models.Notification) error
	Get(key string) (models.Notification, error)
	HGet(key string) ([]models.Notification, error)
}
