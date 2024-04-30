package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type createPaymentRequest struct {
	Plan      string `json:"plan" binding:"required"`
	ReturnURL string `json:"return_url" binding:"required"`
}

type createPaymentResponse struct {
	ConfirmationURL string `json:"confirmation_url"`
	OrderID         string `json:"order_id"`
}

func (h *Handler) createPayment(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input createPaymentRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	paymentResponse, err := h.services.CreatePayment(id, input.Plan, input.ReturnURL)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, createPaymentResponse{
		ConfirmationURL: paymentResponse.Confirmation.ConfirmationURL,
		OrderID:         paymentResponse.Id,
	})
}
