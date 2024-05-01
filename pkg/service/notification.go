package service

import (
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/service/models"
	"github.com/google/uuid"
)

type NotificationService struct {
	notificationRepo repository.Notification
}

func (n NotificationService) GetUserNotifications(userId uuid.UUID) (notifications []models.Notification, err error) {
	repoNotifications, err := n.notificationRepo.GetNotifications(userId)
	if err != nil {
		return notifications, err
	}

	for i := range repoNotifications {
		notifications = append(notifications, models.Notification{
			Type:    repoNotifications[i].Type,
			Payload: repoNotifications[i].Payload,
			Date:    repoNotifications[i].Date,
		})
	}

	return notifications, err
}

func NewNotificationService(notificationRepo repository.Notification) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo}
}
