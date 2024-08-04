package handler

import (
	"ConnectTeam/pkg/service"
	"github.com/gin-contrib/cors"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

// NewHandler creates Handler instance.
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://connecteam.ru"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	auth := router.Group("/auth")
	{
		auth.GET("/refresh/:refresh_token", h.refreshToken)
		auth.POST("/verify-email", h.verifyEmailOnRegistration)
		auth.POST("/sign-up", h.signUp)
		signIn := auth.Group("sign-in")
		{
			signIn.POST("/email", h.signIn)
			signIn.POST("/phone", h.signInWithPhoneNumber)
		}

		auth.PATCH("/password", h.restorePassword)

	}

	userApi := router.Group("/users", h.userIdentity)
	{
		userApi.GET("/:id", h.getUserById)
		userApi.GET("/me", h.getCurrentUser)
		userApi.PATCH("/access", h.changeAccessWithId)
		userApi.GET("/list", h.getUsersList)
		userApi.PATCH("/password", h.changePassword)
		userApi.PATCH("/password/restore", h.restorePasswordAuthorized)
		userApi.POST("/verify-email", h.verifyEmailOnChange)
		userApi.PATCH("/change-email", h.changeEmail)
		userApi.PATCH("/info", h.changePersonalData)
		userApi.PATCH("/company", h.changeCompanyData)
		userApi.PATCH("/image", h.uploadProfileImage)
	}

	plan := router.Group("/plans", h.userIdentity)
	{
		plan.GET("/current", h.getUserActivePlan)
		plan.POST("/", h.createPlan)
		plan.GET("/active", h.getUsersPlans)
		plan.POST("/:user_id", h.setPlan)
		plan.DELETE("/cancel/:id", h.deleteUserPlan)
		plan.POST("/trial", h.getTrial)
		plan.GET("/", h.getUserSubscriptions)
		plan.GET(":id/members", h.getMembers)
		plan.POST("/join/:code", h.addUserToPlan)
		plan.DELETE(":id/members/:user_id", h.deleteUserFromSub)
		plan.PATCH("/upgrade/:id", h.upgradePlan)
		plan.POST("/invite/:id", h.inviteMemberToSub)
		plan.GET("/:id", h.getPlan)
	}
	payment := router.Group("/payment", h.userIdentity)
	{
		payment.POST("/", h.createPayment)
	}
	validator := router.Group("/validate")
	{
		validator.GET("/plan/:code", h.validateInvitationCode)
		validator.GET("/game/:code", h.validateGameInvitationCode)
	}

	topic := router.Group("/topics", h.userIdentity)
	{
		topic.POST("/", h.createTopic)
		topic.GET("/", h.getAllTopics)
		topic.DELETE("/:id", h.deleteTopic)
		topic.PATCH("/:id", h.updateTopic)
		questions := topic.Group("/:id/questions")
		{
			questions.GET("/", h.getAllQuestions)
			questions.POST("/", h.createQuestion)
		}
	}
	question := router.Group("/questions", h.userIdentity)
	{
		question.DELETE("/:id", h.deleteQuestion)
		question.PATCH("/:id", h.updateQuestion)
		question.PUT("/:id/tags", h.updateQuestionTags)
	}
	game := router.Group("/games", h.userIdentity)
	{
		game.POST("/", h.createGame)
		game.GET("/all/:page", h.getGames)
		game.GET("/created/:page", h.getCreatedGames)
		game.GET("/:id", h.getGame)
		game.DELETE("/:id", h.deleteGameFromGameList)
		game.POST("/:code", h.addUserAsParticipant)
		game.GET("/:id/results", h.getResults)
		game.GET("/:id/results/:user_id/tags", h.getTagsResults)
		game.GET("/:id/members", h.getGameMembers)
		game.PATCH("/:id/cancel", h.cancelGame)
		game.POST("/invite/:id", h.inviteMemberToGame)
		game.PATCH("/:id/date", h.changeGameStartDate)
		game.PATCH("/:id/name", h.changeGameName)
	}

	tags := router.Group("/tags", h.userIdentity)
	{
		tags.GET("/", h.getAllTags)
	}

	notifications := router.Group("/notifications", h.userIdentity)
	{
		notifications.POST("/games/:game_id/start", h.createGameStartNotification)
		notifications.PATCH("/", h.readNotifications)
		notifications.GET("/", h.getAllNotifications)
	}

	return router
}
