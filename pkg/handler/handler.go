package handler

import (
	"github.com/gin-contrib/cors"
	"ConnectTeam/pkg/service"
    "fmt"
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
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "http://localhost:5173"
		// },
		MaxAge: 12 * time.Hour,
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

	}

	userApi := router.Group("/users", h.userIdentity)
	{
		userApi.GET("/me", h.getCurrentUser)
		userApi.PATCH("/change-access", h.changeAccessWithId)
		userApi.GET("/list", h.getUsersList)
		userApi.PATCH("/change-password", h.changePassword)
		userApi.POST("/verify-email", h.verifyEmailOnChange)
		userApi.PATCH("/change-email", h.changeEmail)
		userApi.PATCH("/info", h.changePersonalData)
		userApi.PATCH("/company", h.changeCompanyData)
		userApi.GET("/plan", h.getUserPlan)
		userApi.POST("/plan", h.sendPlanRequest)
	}

	plan := router.Group("/plans", h.userIdentity) 
	{
		plan.GET("/current", h.getUserPlan)
		plan.POST("/purchase", h.sendPlanRequest)
		plan.GET("/users-plans", h.getUsersPlans)
	}

	return router
}

func test(c *gin.Context) {

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

func messageHandler(message []byte)  {
	fmt.Println(string(message))
  }