package repository

import (
	"ConnectTeam/pkg/repository/models"
	"ConnectTeam/pkg/repository/notification_service"
	cache "ConnectTeam/pkg/repository/redis"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type NotificationManager struct {
	notificationCache   cache.Cache
	notificationService *notification_service.NotificationService
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

func (n *NotificationManager) SendNotification(userId uuid.UUID) {
	data := map[string]interface{}{
		"user_id": userId.String(),
	}

	// Преобразование данных в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("JSON marshal:", err)
	}

	// Отправка данных через WebSocket в формате JSON
	if n.notificationService.Conn != nil {
		err = n.notificationService.Conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			println(err.Error())
		}
	}
}

func NewNotification(client *redis.Client, service *notification_service.NotificationService) *NotificationManager {
	return &NotificationManager{
		notificationCache:   cache.NewNotificationCache(client),
		notificationService: service,
	}
}
