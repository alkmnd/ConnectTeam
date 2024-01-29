package handler

import (
	connectteam "ConnectTeam"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getCurrentUser(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	assert_id, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "Incorrect auth header")
		return
	}

	user, err := h.services.UserInterface.GetUserById(assert_id) 
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}



	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
		"email": user.Email, 
		// "phone_number": user.PhoneNumber, 
		"first_name": user.FirstName, 
		"second_name": user.SecondName, 
		"access": user.Access,
		"comppany_name": user.CompanyName, 
		"image": user.ProfileImage,
	})
}

type changeAccessInput struct {
	Id int `json:"id,string" binding "required"` 
	NewAccess string `json:"access" binding "required"`
}

func (h *Handler) changeAccessWithId(c *gin.Context) {
	var input changeAccessInput
	_, ok_id := c.Get(userCtx)
	if !ok_id {
		newErrorResponse(c, http.StatusInternalServerError, "User id is not found")
		return
	}
	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	if access != "admin" {
		newErrorResponse(c, http.StatusInternalServerError, "Access error")
		return
	}
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	if err := h.services.UserInterface.UpdateAccessWithId(input.Id, input.NewAccess); err != nil{
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
	} 
}

type getUsersListResponse struct {
	Data []connectteam.UserPublic `json:"data"`
}

func (h *Handler) getUsersList(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	list, err := h.services.UserInterface.GetUsersList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
	c.JSON(http.StatusOK, getUsersListResponse {
		Data: list,
	})
} 

type changePasswordInput struct {
	OldPassword string `json:"old_password" binding "required"`
	NewPassword string `json:"new_password" binding "required"`
}
func (h *Handler) changePassword(c *gin.Context) {
	var input changePasswordInput
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	err = h.services.UserInterface.UpdatePassword(input.OldPassword, input.NewPassword, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
} 

type changeEmailInput struct {
	NewEmail string `json:"new_email" binding "required"`
	Code string `json:"code" binding "required"`
}
func (h *Handler) changeEmail(c *gin.Context) {
	var input changeEmailInput
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	err = h.services.UpdateEmail(id, input.NewEmail, input.Code)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
}

type sendCodeInput struct {
	Email string `json:"email" binding "required"`
	Password string `json:"password" binding "required"`
}
func (h *Handler) verifyEmailOnChange(c *gin.Context) {
	var input sendCodeInput 
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	err = h.services.CheckEmailOnChange(id, input.Email, input.Password)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
}

func (h *Handler) changePersonalData(c *gin.Context) {
	var input connectteam.UserPersonalInfo

	id, err := getUserId(c)
	if err != nil {
		println("1")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	err = h.services.UpdatePersonalData(id,input)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
}

func (h *Handler) changeCompanyData(c *gin.Context) {
	var input connectteam.UserCompanyData

	id, err := getUserId(c)
	if err != nil {
		println("1")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	err = h.services.UpdateCompanyData(id, input)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}


}