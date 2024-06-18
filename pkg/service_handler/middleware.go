package service_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// clientIdentity Проверяет api ключ, переданный в хедере.
func (h *Handler) clientIdentity(c *gin.Context) {
	// Получение ключа из хедера запроса.
	apiKey := c.Request.Header.Get("X-API-Key")

	// Проверка ключа.
	if apiKey != h.apiKey {
		newErrorResponse(c, http.StatusUnauthorized, "incorrect api key")
		return
	}
}
