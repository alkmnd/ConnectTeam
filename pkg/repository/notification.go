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
		err = n.notificationCache.HSet(userId.String(), notifications[i].Id, notifications[i])
		if err != nil {
			println(err.Error())
		}
	}

	return nil
}

func (n *NotificationManager) CreateGameCancelNotification(gameId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Id:      uuid.New(),
		Type:    models.CancelGameNotification,
		Payload: gameId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification.Id, notification)
}

func (n *NotificationManager) CreateGameStartNotification(gameId uuid.UUID, userId uuid.UUID) error {

	notification := models.Notification{
		Id:      uuid.New(),
		Type:    models.StartGameNotification,
		Payload: gameId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification.Id, notification)
}

func (n *NotificationManager) CreateGameInviteNotification(gameId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Id:      uuid.New(),
		Type:    models.InviteGameNotification,
		Payload: gameId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification.Id, notification)
}

func (n *NotificationManager) CreateSubInviteNotification(planId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Id:      uuid.New(),
		Type:    models.InviteSubNotification,
		Payload: planId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification.Id, notification)
}

func (n *NotificationManager) CreateDeleteFromSubNotification(holderId uuid.UUID, userId uuid.UUID) error {
	notification := models.Notification{
		Id:      uuid.New(),
		Type:    models.DeleteFromSubNotification,
		Payload: holderId.String(),
		Date:    time.Now(),
		IsRead:  false,
	}

	return n.notificationCache.HSet(userId.String(), notification.Id, notification)
}

func NewNotification(client *redis.Client) *NotificationManager {
	return &NotificationManager{
		notificationCache: cache.NewNotificationCache(client),
	}
}
