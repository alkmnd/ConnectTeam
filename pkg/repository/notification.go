package repository

import (
	"ConnectTeam/pkg/repository/models"
	cache "ConnectTeam/pkg/repository/redis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type NotificationManager struct {
	notificationCache cache.Cache
}

func (n *NotificationManager) GetNotifications(userId uuid.UUID) ([]models.Notification, error) {
	return n.notificationCache.HGet(userId.String())

}

func (n *NotificationManager) ReadNotifications(userId uuid.UUID) error {
	notifications, err := n.notificationCache.HGet(userId.String())
	if err != nil {
		return err
	}
	for i := range notifications {
		notifications[i].IsRead = true
		_ = n.notificationCache.HSet(userId.String(), notifications[i])
	}

	return nil
}

func (n *NotificationManager) CreateGameCancelNotification(gameId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Type:    models.CancelGameNotification,
		Payload: gameId.String(),
		Date:    time.Now(),
	}

	return n.notificationCache.HSet(userId.String(), notification)
}

func (n *NotificationManager) CreateGameStartNotification(gameId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Type:    models.StartGameNotification,
		Payload: gameId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification)
}

func (n *NotificationManager) CreateGameInviteNotification(gameId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Type:    models.InviteGameNotification,
		Payload: gameId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification)
}

func (n *NotificationManager) CreateSubInviteNotification(planId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Type:    models.InviteSubNotification,
		Payload: planId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification)
}

func (n *NotificationManager) CreateDeleteFromSubNotification(holderId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Type:    models.DeleteFromSubNotification,
		Payload: holderId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification)
}

func NewNotification(client *redis.Client) *NotificationManager {
	return &NotificationManager{
		notificationCache: cache.NewNotificationCache(client),
	}
}
