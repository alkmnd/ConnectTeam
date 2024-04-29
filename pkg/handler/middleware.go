package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	accessCtx           = "access"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, access, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set("userId", userId)
	c.Set("access", access)
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return uuid.Nil, errors.New("user id not found")
	}

	castId, ok := id.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user id is of invalid type")
	}

	return castId, nil
}

func getUserAccess(c *gin.Context) (string, error) {
	access, ok := c.Get(accessCtx)
	if !ok {
		return "", errors.New("user access not found")
	}

	accessString, ok := access.(string)
	if !ok {
		return "", errors.New("user access is of invalid type")
	}

	return accessString, nil
}
