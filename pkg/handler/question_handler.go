package handler

import (
	connectteam "ConnectTeam/models"
	"ConnectTeam/pkg/service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) createQuestion(c *gin.Context) {
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
		newErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}

	var input connectteam.Question
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect format")
		return
	}

	topicId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	id, err := h.services.Question.CreateQuestion(input.Content, topicId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) deleteQuestion(c *gin.Context) {
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

	err = h.services.Question.DeleteQuestion(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

type getQuestionsResponse struct {
	Data []models.Question `json:"data"`
}

func (h *Handler) getAllQuestions(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	topicId, err := uuid.Parse(c.Param("id"))
	questions, err := h.services.Question.GetAll(topicId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getQuestionsResponse{
		Data: questions,
	})
}

type updateQuestionInput struct {
	NewContent string `json:"new_content" binding:"required,min=1,max=300"`
}

func (h *Handler) updateQuestion(c *gin.Context) {
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
		newErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input updateQuestionInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect format")
		return
	}

	q, err := h.services.Question.UpdateQuestion(input.NewContent, id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, connectteam.Question{
		Id:      q.Id,
		TopicId: q.TopicId,
		Content: q.Content,
	})
}

type getTagsResponse struct {
	Data []models.Tag `json:"data"`
}

func (h *Handler) getAllTags(c *gin.Context) {
	tags, err := h.services.GetAllTags()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getTagsResponse{
		Data: tags,
	})
}

type updateQuestionTagsInput struct {
	Tags []models.Tag `json:"tags"`
}

func (h *Handler) updateQuestionTags(c *gin.Context) {
	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var input updateQuestionTagsInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	questionId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}
	tags, err := h.services.UpdateQuestionTags(questionId, input.Tags)
	c.JSON(http.StatusOK, getTagsResponse{
		Data: tags,
	})
}
