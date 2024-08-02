package service_handler

import (
	connectteam "ConnectTeam/models"
	"ConnectTeam/pkg/service"
	"ConnectTeam/pkg/service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type Handler struct {
	services *service.Service
	apiKey   string
}

func NewHandler(services *service.Service, apiKey string) *Handler {
	return &Handler{services: services, apiKey: apiKey}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	httpService := router.Group("/api", h.clientIdentity)
	{
		game := httpService.Group("/games")
		{
			game.GET("/:id", h.getGame)
			game.PATCH("/start/:id", h.startGame)
			game.POST(":id/results", h.saveResults)
			game.PATCH("/end/:id", h.endGame)
			game.GET("/results/:id", h.getResults)
		}
		topic := httpService.Group("/topics")
		{
			topic.GET("/:id", h.getTopic)
			topic.GET("list/:limit", h.getTopicsWithLimit)
		}
		question := httpService.Group("/questions")
		{
			question.GET("/", h.getRandWithLimit)
		}
		user := httpService.Group("/users")
		{
			user.GET("/:id", h.getUserById)
			user.GET("/:id/plan", h.getUserActivePlan)
		}
	}

	return router
}

func (h *Handler) getUserById(c *gin.Context) {
	println("getUserById")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var user connectteam.UserPublic

	user, err = h.services.User.GetUserById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":          id,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"second_name": user.SecondName,
	})
}

func (h *Handler) getUserActivePlan(c *gin.Context) {

	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userPlan, err := h.services.Plan.GetUserActivePlan(userId)

	if err != nil {
		c.Status(204)
		return
	}

	var invitationCode string

	if userPlan.PlanType == "premium" &&
		userPlan.Status == connectteam.Active {
		invitationCode = userPlan.InvitationCode
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":              userPlan.Id,
		"plan_type":       userPlan.PlanType,
		"holder_id":       userPlan.HolderId,
		"expiry_date":     userPlan.ExpiryDate,
		"plan_access":     userPlan.PlanAccess,
		"status":          userPlan.Status,
		"invitation_code": invitationCode,
		"is_trial":        userPlan.IsTrial,
	})
}

func (h *Handler) getGame(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	game, err := h.services.Game.GetGame(gameId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":              game.Id,
		"name":            game.Name,
		"creator_id":      game.CreatorId,
		"start_date":      game.StartDate,
		"status":          game.Status,
		"invitation_code": game.InvitationCode,
	})
}

func (h *Handler) startGame(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))
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

type saveResultsInput struct {
	Results []Rates
}

type Rates struct {
	Value  int         `json:"value"`
	Tags   []uuid.UUID `json:"tags"`
	UserId uuid.UUID   `json:"user_id"`
	Name   string      `json:"name"`
}

func (h *Handler) saveResults(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input saveResultsInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for i := range input.Results {
		_ = h.services.SaveResults(gameId, input.Results[i].UserId, input.Results[i].Value, input.Results[i].Name)
		for k := range input.Results[i].Tags {
			_ = h.services.SaveTagsResults(input.Results[i].UserId, gameId, input.Results[i].Tags[k], input.Results[i].Name)
		}

	}
	c.Status(http.StatusOK)
}

func (h *Handler) endGame(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))
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

func (h *Handler) getTopic(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	topic, err := h.services.GetTopic(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, topic)
}

type getResultsResponse struct {
	Data []models.UserResult `json:"data"`
}

func (h *Handler) getResults(c *gin.Context) {
	gameId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	results, err := h.services.GetResults(gameId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getResultsResponse{
		Data: results,
	})
}

func (h *Handler) getRandWithLimit(c *gin.Context) {
	topicId, err := uuid.Parse(c.Query("topic_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	questions, err := h.services.Question.GetRandWithLimit(topicId, limit)

	c.JSON(http.StatusOK, questions)
}

type getQuestionsResponse struct {
	Data []models.Question `json:"data"`
}

func (h *Handler) getTopicsWithLimit(c *gin.Context) {

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	topics, err := h.services.Topic.GetRandWithLimit(limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, topics)
}

type getTopicsResponse struct {
	Data []connectteam.Topic `json:"data"`
}
