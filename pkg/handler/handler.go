package handler

import (
	"ConnectTeam/pkg/service"
    "fmt"
    "net/http"

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
	auth := router.Group("/auth") 
	{
		verify := auth.Group("/verify") 
		{
			verify.POST("/user", h.verifyUser)			
			verify.POST("/phone", h.verifyPhone)
			verify.POST("/email", h.verifyEmail)
		}
		auth.POST("/sign-up", h.signUp)
		signIn := auth.Group("sign-in") 
		{
			signIn.POST("/email", h.signInWithEmail)
			signIn.POST("/phone", h.signInWithPhoneNumber)
		}

	}

	userApi := router.Group("/user", h.userIdentity)
	{
		userApi.GET("/me", h.getCurrentUser)
		userApi.PATCH("/change-access", h.changeAccessById)
		userApi.GET("/list", h.getUsersList)
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

func messageHandler(message []byte)  {
	fmt.Println(string(message))
  }