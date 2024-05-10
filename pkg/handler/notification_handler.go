package handler

import (
	"ConnectTeam/pkg/handler/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type getAllNotificationsResponse struct {
	Data []models.Notification `json:"data"`
}

func (h *Handler) getAllNotifications(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	serviceNotifications, err := h.services.Notification.GetUserNotifications(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var response []models.Notification
	for i := range serviceNotifications {
		response = append(response, models.Notification{
			Type:    serviceNotifications[i].Type,
			Payload: serviceNotifications[i].Payload,
			Date:    serviceNotifications[i].Date,
		})
	}
	c.JSON(http.StatusOK, getAllNotificationsResponse{
		Data: response,
	})

}

func (h *Handler) createGameStartNotification(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.CreateGameStartNotification(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)

}

func (h *Handler) readNotifications(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.ReadNotifications(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
