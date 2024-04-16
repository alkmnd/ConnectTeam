package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type saveResultsInput struct {
	Results map[int]int `json:"results"`
}

func (h *Handler) saveResults(c *gin.Context) {
	gameId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input saveResultsInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for i, v := range input.Results {
		_ = h.services.SaveResults(gameId, i, v)
		//if err != nil {
		//
		//}
	}
	c.Status(http.StatusOK)
}

func (h *Handler) endGame(c *gin.Context) {
	gameId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.EndGame(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) startGame(c *gin.Context) {
	gameId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.StartGame(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getTopic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	topic, err := h.services.GetTopic(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	jsonTopic, err := json.Marshal(topic)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, jsonTopic)
}

func (h *Handler) getRandWithLimit(c *gin.Context) {
	topicId, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	questions, err := h.services.Question.GetRandWithLimit(topicId, limit)

	c.JSON(http.StatusOK, getQuestionsResponse{
		Data: questions,
	})

}
