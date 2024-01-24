package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getCurrentUser(c *gin.Context) {
	id, _ := c.Get(userCtx)
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
		"id":id,
		"email": user.Email, 
		"phone_number": user.PhoneNumber, 
		"first_name": user.FirstName, 
		"second_name": user.SecondName, 
		"access": user.Access,
	})
}

type changeAccessInput struct {
	Id int `json:"id,string" binding "required"` 
	NewAccess string `json:"access" binding "required"`
}

func (h *Handler) changeAccessById(c *gin.Context) {
	var input changeAccessInput
	access, _ := c.Get(accessCtx)
	assert_access, ok := access.(string)

	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "Incorrect auth header")
		return
	}

	if assert_access == "admin" {
		if err := c.BindJSON(&input); err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		}

		if err := h.services.UserInterface.ChangeAccessById(input.Id, input.NewAccess); err != nil{
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		} 
	}
}