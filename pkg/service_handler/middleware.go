package service_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) clientIdentity(c *gin.Context) {
	apiKey := c.Request.Header.Get("X-API-Key")

	if apiKey != h.ApiKey {
		newErrorResponse(c, http.StatusUnauthorized, "incorrect api key")
		return
	}

}
