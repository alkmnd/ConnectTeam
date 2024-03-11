package handler

import (
	connectteam "ConnectTeam"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserActivePlan(c *gin.Context) {
	id, err := getUserId(c)
	println(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userPlan, err := h.services.Plan.GetUserActivePlan(id)

	if err != nil {
		c.Status(204)
		return
	}

	var invitationCode string

	if userPlan.PlanType == "premium" &&
		userPlan.Status == connectteam.Active {
		invitationCode = userPlan.InvitationCode
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":              userPlan.Id,
		"plan_type":       userPlan.PlanType,
		"user_id":         userPlan.UserId,
		"holder_id":       userPlan.HolderId,
		"expiry_date":     userPlan.ExpiryDate,
		"plan_access":     userPlan.PlanAccess,
		"status":          userPlan.Status,
		"invitation_code": invitationCode,
		"is_trial":        userPlan.IsTrial,
	})
}

func (h *Handler) selectPlan(c *gin.Context) {
	var input connectteam.UserPlan
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	input.UserId = id
	input.HolderId = id
	input.PlanAccess = "holder"
	input.Status = connectteam.OnConfirm
	input.ExpiryDate = time.Time{}

	plan, err := h.services.CreatePlan(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":          plan.Id,
		"user_id":     plan.UserId,
		"holder_id":   plan.HolderId,
		"plan_type":   plan.PlanType,
		"plan_access": plan.PlanAccess,
		"status":      plan.Status,
		"duration":    plan.Duration,
		"expiry_date": plan.ExpiryDate,
		"is_trial":    false,
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

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	list, err := h.services.Plan.GetUsersPlans()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getUsersPlansResponse{
		Data: list,
	})
}

type getUserSubscriptionsResponse struct {
	Data []connectteam.UserPlan `json:"data"`
}

func (h *Handler) getUserSubscriptions(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	plans, err := h.services.Plan.GetUserSubscriptions(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getUserSubscriptionsResponse{
		Data: plans,
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

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
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
	PlanType   string `json:"plan_type" binding:"required"`
	ExpiryDate string `json:"expiry_date" binding:"required"`
}

func (h *Handler) getTrial(c *gin.Context) {
	println("hui")
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	println("pizda")

	subscriptionExists, err := h.services.Plan.CheckIfSubscriptionExists(id)
	if err != nil {
		println("1")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if subscriptionExists {
		println("meow")
		newErrorResponse(c, http.StatusForbidden, "user could not get trial")
		return
	}

	plan, err := h.services.Plan.CreateTrialPlan(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":          plan.Id,
		"user_id":     plan.UserId,
		"holder_id":   plan.HolderId,
		"plan_type":   plan.PlanType,
		"plan_access": plan.PlanAccess,
		"status":      plan.Status,
		"duration":    plan.Duration,
		"expiry_date": plan.ExpiryDate,
		"is_trial":    plan.IsTrial,
	})

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

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var input newPlanInput
	if err := c.BindJSON(&input); err != nil {
		println("1")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}
	err = h.services.SetPlanByAdmin(userId, input.PlanType, input.ExpiryDate)
	println(input.PlanType)
	if err != nil {
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

	if access != string(connectteam.Admin) && access != string(connectteam.SuperAdmin) {
		newErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}

	err = h.services.DeletePlan(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) addUserToPlan(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	code := c.Param("code")
	holderId, err := h.services.GetHolderWithInvitationCode(code)
	if holderId == 0 {
		newErrorResponse(c, http.StatusNotFound, "incorrect invitation code")
		return
	}

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if holderId == id {
		newErrorResponse(c, http.StatusForbidden, "holder and user is equal")
		return
	}
	holderPlan, err := h.services.GetUserActivePlan(holderId)
	if holderPlan.Status != connectteam.Active {
		newErrorResponse(c, http.StatusForbidden, "invitor subscription is not active")
		return
	}

	members, err := h.services.GetMembers(code)
	if len(members) == 4 {
		newErrorResponse(c, http.StatusForbidden, "max number of members")
		return
	}

	plan, err := h.services.AddUserToAdvanced(holderPlan, id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":          plan.Id,
		"user_id":     plan.UserId,
		"holder_id":   plan.HolderId,
		"plan_type":   plan.PlanType,
		"plan_access": plan.PlanAccess,
		"status":      plan.Status,
		"duration":    plan.Duration,
		"expiry_date": plan.ExpiryDate,
	})
}

func (h *Handler) validateInvitationCode(c *gin.Context) {
	code := c.Param("code")
	id, err := h.services.GetHolderWithInvitationCode(code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if id == 0 || len(code) == 0 {
		newErrorResponse(c, http.StatusNotFound, "incorrect invitation code")
		return
	}

	var userPlan connectteam.UserPlan
	userPlan, err = h.services.Plan.GetUserActivePlan(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userPlan.Status != connectteam.Active {
		newErrorResponse(c, http.StatusForbidden, "subscription is not active")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"holder_id": id,
	})
}

type getMembersResponse struct {
	Data []connectteam.UserPublic `json:"data"`
}

func (h *Handler) getMembers(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	code := c.Param("code")

	users, err := h.services.Plan.GetMembers(code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getMembersResponse{
		Data: users,
	})
}

func (h *Handler) deleteUserFromSub(c *gin.Context) {
	holderId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	plan, err := h.services.GetUserPlan(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if holderId != plan.HolderId {
		newErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	err = h.services.DeleteUserFromSub(plan.Id)
	if holderId != plan.HolderId {
		newErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
