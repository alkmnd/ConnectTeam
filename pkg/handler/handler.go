package handler

import (
	"ConnectTeam/pkg/service"

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
			verify.POST("/verify-user", h.verifyUser)			
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

	api := router.Group("/api", h.userIdentity)
	{
		api.GET("me", h.getCurrentUser)
	}

	return router
}