package repository

import (
	"ConnectTeam/pkg/repository/models"
	redis2 "ConnectTeam/pkg/repository/redis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type NotificationManager struct {
	notificationCache redis2.Cache
}

func (n NotificationManager) GetNotifications(userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (n NotificationManager) CreateGameCancelNotification(gameId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Type:    models.CancelGameNotification,
		Payload: gameId.String(),
	}

	return n.notificationCache.HSet(userId.String(), notification)
}

func (n NotificationManager) CreateGameStartNotification(gameId uuid.UUID, userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (n NotificationManager) CreateGameInviteNotification(gameId uuid.UUID, userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (n NotificationManager) CreateSubInviteNotification(holderId uuid.UUID, userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (n NotificationManager) CreateDeleteFromSubNotification(holderId uuid.UUID, userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func NewNotification(client *redis.Client) *NotificationManager {
	return &NotificationManager{
		notificationCache: redis2.NewNotificationCache(client),
	}
}
