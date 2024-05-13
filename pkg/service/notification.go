package service

import (
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/service/models"
	"github.com/google/uuid"
)

type NotificationService struct {
	notificationRepo repository.Notification
	gameRepo         repository.Game
}

func (n *NotificationService) GetUserNotifications(userId uuid.UUID) (notifications []models.Notification, err error) {
	repoNotifications, err := n.notificationRepo.GetNotifications(userId)
	if err != nil {
		return notifications, err
	}

	for i := range repoNotifications {
		notifications = append(notifications, models.Notification{
			Type:    repoNotifications[i].Type,
			Payload: repoNotifications[i].Payload,
			Date:    repoNotifications[i].Date,
			IsRead:  repoNotifications[i].IsRead,
		})
	}

	return notifications, err
}

func (n *NotificationService) CreateGameStartNotification(gameId uuid.UUID) error {
	gameMembers, err := n.gameRepo.GetGameParticipants(gameId)
	if err != nil {
		return err
	}
	for i := range gameMembers {
		_ = n.notificationRepo.CreateGameStartNotification(gameId, gameMembers[i].Id)
	}
	return nil
}
func NewNotificationService(notificationRepo repository.Notification,
	gameRepo repository.Game) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo, gameRepo: gameRepo}
}
func (n *NotificationService) ReadNotifications(userId uuid.UUID) error {
	return n.notificationRepo.ReadNotifications(userId)
}
