package handler

import (
	connectteam "ConnectTeam"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) createTopic(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var input connectteam.Topic
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Topic.CreateTopic(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getTopicsResponse struct {
	Data []connectteam.Topic `json:"data"`
}

func (h *Handler) getAllTopics(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	topics, err := h.services.Topic.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getTopicsResponse{
		Data: topics,
	})
}

func (h *Handler) deleteTopic(c *gin.Context) {

	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}

	if err := h.services.DeleteTopic(id); err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

type updateTopicInput struct {
	Title string `json:"title" binding:"required"`
}

func (h *Handler) updateTopic(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input updateTopicInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len([]rune(input.Title)) < 3 {
		newErrorResponse(c, http.StatusInternalServerError, "Incorrect title")
		return
	}

	err = h.services.Topic.UpdateTopic(id, input.Title)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
