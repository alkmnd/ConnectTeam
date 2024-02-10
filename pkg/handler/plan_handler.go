package handler

import (
	connectteam "ConnectTeam"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserPlan(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	userPlan, err := h.services.Plan.GetUserPlan(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "User has no plan")
		return 
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"plan_type":userPlan.PlanType, 
		"user_id":userPlan.UserId,
		"holder_id":userPlan.HolderId, 
		"expiry_date":userPlan.ExpiryDate,
		"plan_access":userPlan.PlanAccess,
		"confirmed":userPlan.Confirmed,
	})
}

func (h *Handler) selectPlan(c *gin.Context) {
	var input connectteam.UserPlan
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
	println(input.Confirmed)

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	input.UserId = id
	input.HolderId = id
	input.PlanAccess = "holder"

	plan, err := h.services.CreatePlan(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": plan.UserId,
		"holder_id": plan.HolderId, 
		"plan_type": plan.PlanType, 
		"plan_access": plan.PlanAccess, 
		"confirmed": plan.Confirmed, 
		"duration": plan.Duration,
		"expiry_date": plan.ExpiryDate,
	})
}

type getUsersPlansResponse struct {
	Data []connectteam.UserPlan `json:"data"`
}

func (h *Handler) getUsersPlans(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	if access != string(connectteam.Admin) && access != string(connectteam.Superadmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return 
	}
	
	list, err := h.services.Plan.GetUsersPlans()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}
	c.JSON(http.StatusOK, getUsersPlansResponse {
		Data: list,
	})
}

func (h *Handler) confirmPlan(c *gin.Context) {
	_, err := getUserId(c)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	access, err := getUserAccess(c)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.Admin) && access != string(connectteam.Superadmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return 
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}

	err = h.services.Plan.ConfirmPlan(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

type newPlanInput struct {
	Duration int `json:"duration"`
	PlanType string `json:"plan_type"`
}
func (h *Handler) setPlan(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} 

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.Admin) && access != string(connectteam.Superadmin)  {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var input newPlanInput
	if err := c.BindJSON(&input); err != nil {
		println("1")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		println("2")
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}
	err = h.services.SetPlanByAdmin(user_id, input.Duration, input.PlanType)
	println(input.PlanType)
	if err != nil {
		println("3")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteUserPlan(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} 

	access, err := getUserAccess(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if access != string(connectteam.Admin) && access != string(connectteam.Superadmin)  {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}

	err = h.services.DeletePlan(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}