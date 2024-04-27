package handler

import (
	connectteam "ConnectTeam"
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/swaggo/files"
	//"github.com/swaggo/gin-swagger"
)

// @Summary      Sign up
// @Description  create user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body connectteam.User true "User created"
// @Success      200  {object}  connectteam.UserPublic
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input connectteam.UserSignUpRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()+"Incorrect format")
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

// func (h *Handler) verifyPhone(c *gin.Context) {
// 	var input connectteam.VerifyPhone

// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	confirmationCode, err := h.services.Authorization.VerifyPhone(input)

// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"confirmationCode": confirmationCode,
// 	})
// }

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

	err := h.services.Authorization.VerifyEmail(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

//func (h *Handler) verifyUser(c *gin.Context) {
//	var input connectteam.VerifyUser
//
//	if err := c.BindJSON(&input); err != nil {
//		newErrorResponse(c, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	println(input.Id)
//
//	err := h.services.Authorization.VerifyUser(input)
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	c.JSON(http.StatusOK, map[string]interface{}{
//		"id": input.Id,
//	})
//}

// func (h *Handler) signUpWithPhone(c *gin.Context) {
// 	var input connectteam.User

// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 	}
// }

type signInWithEmailInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInWithPhoneNumInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func (h *Handler) signInWithEmail(c *gin.Context) {
	var input signInWithEmailInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	access, token, err := h.services.Authorization.GenerateToken(input.Email, input.Password, true)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":  token,
		"access": access,
	})
}

func (h *Handler) signInWithPhoneNumber(c *gin.Context) {
	var input signInWithPhoneNumInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	access, token, err := h.services.Authorization.GenerateToken(input.PhoneNumber, input.Password, false)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":  token,
		"access": access,
	})
}
