package service_handler

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/service"
	"ConnectTeam/pkg/service/models"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	services *service.Service
	ApiKey   string
}

func NewHandler(services *service.Service, apiKey string) *Handler {
	return &Handler{services: services, ApiKey: apiKey}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	httpService := router.Group("/api", h.clientIdentity)
	{
		game := httpService.Group("/games")
		{
			game.GET("/:id", h.getGame)
			game.PATCH("/start/:id", h.startGame)
			game.POST("/results", h.saveResults)
			game.PATCH("/end/:id", h.endGame)
		}
		topic := httpService.Group("/topics")
		{
			topic.GET("/:id", h.getTopic)
			topic.GET("list/:limit", h.getTopicsWithLimit)
		}
		question := httpService.Group("/questions")
		{
			question.GET("/{:id}/{:limit}", h.getRandWithLimit)
			//question.GET("/{:id}/{:limit}", h.getQuestionTags)
		}
		user := httpService.Group("/users")
		{
			user.GET("/:id/plan", h.getUserActivePlan)
		}
	}

	return router
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
	Results map[uuid.UUID]Rates `json:"results"`
}

type Rates struct {
	Value int         `json:"value"`
	Tags  []uuid.UUID `json:"tags"`
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

	for i, v := range input.Results {
		_ = h.services.SaveResults(gameId, i, v.Value)
		for k := range input.Results[i].Tags {
			_ = h.services.CreateTagsUsers(i, gameId, input.Results[i].Tags[k])
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
	jsonTopic, err := json.Marshal(topic)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, jsonTopic)
}

func (h *Handler) getRandWithLimit(c *gin.Context) {
	topicId, err := uuid.Parse(c.Param("topic_id"))
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

	c.JSON(http.StatusOK, getTopicsResponse{
		Data: topics,
	})
}

type getTopicsResponse struct {
	Data []connectteam.Topic `json:"data"`
}
