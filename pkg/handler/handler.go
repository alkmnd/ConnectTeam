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
		AllowOrigins:     []string{"http://localhost:5432"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5432"
		},
		MaxAge: 12 * time.Hour,
	}))
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