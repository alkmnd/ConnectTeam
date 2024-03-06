package handler

import (
	connectteam "ConnectTeam"
	"log"
	"net/http"
	"strconv"

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

	user, err := h.services.User.GetUserById(assert_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":          id,
		"email":       user.Email,
		"description": user.Description,
		// "phone_number": user.PhoneNumber,
		"first_name":    user.FirstName,
		"second_name":   user.SecondName,
		"access":        user.Access,
		"company_name":  user.CompanyName,
		"company_info":  user.CompanyInfo,
		"company_url":   user.CompanyURL,
		"company_logo":  user.CompanyLogo,
		"profile_image": user.ProfileImage,
	})
}

func (h *Handler) getUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var user connectteam.UserPublic
	user, err = h.services.User.GetUserById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":          id,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"second_name": user.SecondName,
	})
}

func (h *Handler) restorePasswordAuthorized(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.RestorePasswordAuthorized(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

type changeAccessInput struct {
	Id     int                     `json:"id,string" binding:"required"`
	Access connectteam.AccessLevel `json:"access" binding:"required"`
}

func (h *Handler) changeAccessWithId(c *gin.Context) {
	var input changeAccessInput
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "User id is not found")
		return
	}
	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	if input.Id == id {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.UpdateAccessWithId(input.Id, input.Access); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.Plan.DeletePlan(input.Id); err != nil {
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

	list, err := h.services.User.GetUsersList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getUsersListResponse{
		Data: list,
	})
}

type changePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
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

	err = h.services.User.UpdatePassword(input.OldPassword, input.NewPassword, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

type changeEmailInput struct {
	NewEmail string `json:"new_email" binding:"required"`
	Code     string `json:"code" binding:"required"`
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

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

type sendCodeInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	err = h.services.UpdatePersonalData(id, input)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
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

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

const MAX_UPLOAD_SIZE = 5 << 20

var IMAGE_TYPES = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

type uploadResponse struct {
	Status string `json:"status"`
	Msg    string `json:"message,omitempty"`
	URL    string `json:"url,omitempty"`
}

func (h *Handler) uploadProfileImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MAX_UPLOAD_SIZE)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("suka")
		c.JSON(http.StatusBadRequest, &uploadResponse{
			Status: "error",
			Msg:    err.Error(),
		})
		return
	}

	defer file.Close()

	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)

	if _, ex := IMAGE_TYPES[fileType]; !ex {
		c.JSON(http.StatusBadRequest, &uploadResponse{
			Status: "error",
			Msg:    "file type is not supported",
		})
		return
	}
	url, err := h.services.Upload(c.Request.Context(), file, fileHeader.Size, fileType)
	if err != nil {
		c.JSON(http.StatusBadRequest, &uploadResponse{
			Status: "error",
			Msg:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &uploadResponse{
		Status: "ok",
		URL:    url,
	})
}
