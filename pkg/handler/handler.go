package handler

import (
	"ConnectTeam/pkg/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	auth := router.Group("/auth")
	{
		auth.POST("/verify-user", h.verifyUser)
		auth.POST("/verify-email", h.verifyEmailOnRegistration)
		// verify := auth.Group("/verify")
		// {
		// 	verify.POST("/user", h.verifyUser)
		// 	verify.POST("/phone", h.verifyPhone)
		// 	verify.POST("/email", h.verifyEmail)
		// }
		auth.POST("/sign-up", h.signUp)
		signIn := auth.Group("sign-in")
		{
			signIn.POST("/email", h.signInWithEmail)
			signIn.POST("/phone", h.signInWithPhoneNumber)
		}

		auth.PATCH("/password", h.restorePassword)

	}

	// do verification

	userApi := router.Group("/users", h.userIdentity)
	{
		userApi.GET("/:id")
		userApi.GET("/me", h.getCurrentUser)
		userApi.PATCH("/access", h.changeAccessWithId)
		userApi.GET("/list", h.getUsersList)
		userApi.PATCH("/change-password", h.changePassword)
		userApi.GET("/password", h.restorePasswordAuthorized)
		userApi.POST("/verify-email", h.verifyEmailOnChange)
		userApi.PATCH("/change-email", h.changeEmail)
		userApi.PATCH("/info", h.changePersonalData)
		userApi.PATCH("/company", h.changeCompanyData)
		userApi.PATCH("/upload-image", h.uploadProfileImage)
		// userApi.GET("/plan", h.getUserPlan)
		// userApi.POST("/plan", h.sendPlanRequest)
	}

	plan := router.Group("/plans", h.userIdentity)
	{
		plan.GET("/current", h.getUserActivePlan)
		plan.POST("/purchase", h.selectPlan)
		plan.GET("/active", h.getUsersPlans)
		plan.PATCH("/:id", h.confirmPlan)
		plan.POST("/:user_id", h.setPlan)
		plan.DELETE("/:id", h.deleteUserPlan)
		plan.POST("/trial", h.getTrial)
		plan.GET("/", h.getUserSubscriptions)
		plan.GET("/validate/:code", h.validateInvitationCode)
		plan.GET("members/:code", h.getMembers)
		// delete plan
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
	questions := router.Group("/questions", h.userIdentity)
	{
		questions.DELETE("/:id", h.deleteQuestion)
		questions.PATCH("/:id", h.updateQuestion)
	}

	return router
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		print(err.Error())
		return
	}
	defer conn.Close()

	for {
		_, message, _ := conn.ReadMessage()

		conn.WriteMessage(websocket.TextMessage, message)
		go messageHandler(message)
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
