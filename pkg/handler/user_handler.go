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

	if access == "admin" {
		if err := c.BindJSON(&input); err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return 
		}

		if err := h.services.UserInterface.ChangeAccessById(input.Id, input.NewAccess); err != nil{
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		} 
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