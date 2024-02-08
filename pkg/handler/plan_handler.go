package handler

import (
	connectteam "ConnectTeam"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) sendPlanRequest(c *gin.Context) {
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
