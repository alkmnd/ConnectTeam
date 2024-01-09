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
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in-email", h.signInWithEmail)
		auth.POST("/sign-in-phone", h.signInWithPhoneNumber)

	}

	// api := router.Group("/api", h.userIdentity)

	return router
}