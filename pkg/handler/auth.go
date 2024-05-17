package handler

import (
	connectteam "ConnectTeam/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input connectteam.UserSignUpRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect format")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) refreshToken(c *gin.Context) {
	refreshToken := c.Param("refresh_token")
	id, err := h.services.Authorization.ParseRefreshToken(refreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := h.services.User.GetUserCredentials(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Create access token.
	_, token, err, _ := h.services.Authorization.GenerateAccessToken(user.Email, user.PasswordHash)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send response.
	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token,
	})

}

type restorePasswordInput struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) restorePassword(c *gin.Context) {
	var input restorePasswordInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	println(input.Email)

	err := h.services.RestorePassword(input.Email)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) verifyEmailOnRegistration(c *gin.Context) {
	var input connectteam.VerifyEmail

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ifExists, err := h.services.CheckIfExist(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if ifExists {
		newErrorResponse(c, http.StatusForbidden, "email is already used")
		return
	}
	err = h.services.Authorization.VerifyEmail(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

type signInWithPhoneNumInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type signInWithEmailInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// signIn process sign in request.
func (h *Handler) signIn(c *gin.Context) {
	var input signInWithEmailInput

	// Validate input.
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create access token.
	access, token, err, id := h.services.Authorization.GenerateAccessToken(input.Email, h.services.GeneratePasswordHash(input.Password))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Creates refresh token.
	refreshToken, err := h.services.Authorization.GenerateRefreshToken(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send response.
	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  token,
		"refresh_token": refreshToken,
		"access":        access,
		"user_id":       id,
	})
}

func (h *Handler) signInWithPhoneNumber(c *gin.Context) {
	var input signInWithPhoneNumInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	access, token, err, _ := h.services.Authorization.GenerateAccessToken(input.PhoneNumber, h.services.Authorization.GeneratePasswordHash(input.Password))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":  token,
		"access": access,
	})
}
