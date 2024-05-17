package repository

import (
	"ConnectTeam/pkg/repository/models"
	"errors"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type YooClient struct {
	yooClient *yookassa.Client
}

func NewYookassaClient(yooClient *yookassa.Client) *YooClient {
	return &YooClient{yooClient: yooClient}
}

func (c *YooClient) CreatePayment(paymentRequest models.PaymentRequest) (models.PaymentResponse, error) {
	paymentHandler := yookassa.NewPaymentHandler(c.yooClient)
	payment, err := paymentHandler.CreatePayment(&yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    paymentRequest.Amount.Value,
			Currency: "RUB",
		},
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: paymentRequest.Confirmation.ReturnUrl,
		},
		Description: "Test payment",
		Capture:     true,
		Metadata:    paymentRequest.MetaData,
	})

	if err != nil {
		return models.PaymentResponse{}, err
	}

	return models.PaymentResponse{
		Id:     payment.ID,
		Status: string(payment.Status),
		Paid:   payment.Paid,
		Amount: models.Amount{
			Value:    payment.Amount.Value,
			Currency: payment.Amount.Currency,
		},
		Confirmation: payment.Confirmation,
		CreatedAt:    payment.CreatedAt,
		Description:  payment.Description,
		Recipient: models.Recipient{
			AccountId: payment.Recipient.AccountId,
			GatewayId: payment.Recipient.GatewayId,
		},
		Refundable: payment.Refundable,
		Test:       payment.Test,
	}, err
}

func (c *YooClient) GetPayment(orderId string) (models.PaymentResponse, error) {
	paymentHandler := yookassa.NewPaymentHandler(c.yooClient)
	payment, _ := paymentHandler.FindPayment(orderId)
	if payment == nil {
		return models.PaymentResponse{}, errors.New("payment not found")
	}
	m := payment.Metadata.(map[string]interface{})
	var metadata models.MetaData
	if userId, ok := m["user_id"].(string); ok {
		metadata.UserId = userId
	}
	if planType, ok := m["plan_type"].(string); ok {
		metadata.PlanType = planType
	}

	return models.PaymentResponse{
		Id:     payment.ID,
		Status: string(payment.Status),
		Paid:   payment.Paid,
		Amount: models.Amount{
			Value:    payment.Amount.Value,
			Currency: payment.Amount.Currency,
		},
		Confirmation: payment.Confirmation,
		CreatedAt:    payment.CreatedAt,
		Description:  payment.Description,
		Metadata:     metadata,
		Recipient: models.Recipient{
			AccountId: payment.Recipient.AccountId,
			GatewayId: payment.Recipient.GatewayId,
		},
		Refundable: payment.Refundable,
		Test:       payment.Test,
	}, nil
}
