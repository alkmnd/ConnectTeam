package service

import (
	"ConnectTeam/pkg/repository"
	repoModels "ConnectTeam/pkg/repository/models"
	"ConnectTeam/pkg/service/models"
	"errors"
	"github.com/google/uuid"
)

type PaymentService struct {
	client repository.Payment
}

func NewPaymentService(client repository.Payment) *PaymentService {
	return &PaymentService{client: client}
}

func (s *PaymentService) CreatePayment(userId uuid.UUID, plan string, returnURL string) (models.PaymentResponse, error) {
	var paymentRequest repoModels.PaymentRequest
	paymentRequest.Confirmation.ReturnUrl = returnURL
	switch plan {
	case "basic-to-advanced":
		paymentRequest.Amount = models.Amount{
			Value:    "50",
			Currency: "RUB",
		}
		paymentRequest.MetaData = repoModels.MetaData{
			UserId:   userId.String(),
			PlanType: "advanced",
		}
	case "basic-to-premium":
		paymentRequest.Amount = models.Amount{
			Value:    "100",
			Currency: "RUB",
		}
		paymentRequest.MetaData = repoModels.MetaData{
			UserId:   userId.String(),
			PlanType: "premium",
		}
	case "advanced-to-premium":
		paymentRequest.Amount = models.Amount{
			Value:    "50",
			Currency: "RUB",
		}
		paymentRequest.MetaData = repoModels.MetaData{
			UserId:   userId.String(),
			PlanType: "premium",
		}
	case "basic":
		paymentRequest.Amount = models.Amount{
			Value:    "100",
			Currency: "RUB",
		}
		paymentRequest.MetaData = repoModels.MetaData{
			UserId:   userId.String(),
			PlanType: plan,
		}
	case "advanced":
		paymentRequest.Amount = models.Amount{
			Value:    "150",
			Currency: "RUB",
		}
		paymentRequest.MetaData = repoModels.MetaData{
			UserId:   userId.String(),
			PlanType: plan,
		}
	case "premium":
		paymentRequest.Amount = models.Amount{
			Value:    "200",
			Currency: "RUB",
		}
		paymentRequest.MetaData = repoModels.MetaData{
			UserId:   userId.String(),
			PlanType: plan,
		}

	default:
		return models.PaymentResponse{}, errors.New("incorrect subscription plan")
	}

	payment, err := s.client.CreatePayment(paymentRequest)
	if err != nil {
		return models.PaymentResponse{}, errors.New("incorrect subscription plan")
	}

	m := payment.Confirmation.(map[string]interface{})
	var confirmation models.Confirmation
	if confirmationURL, ok := m["confirmation_url"].(string); ok {
		confirmation.ConfirmationURL = confirmationURL
	}

	return models.PaymentResponse{
		Id:     payment.Id,
		Status: payment.Status,
		Paid:   payment.Paid,
		Amount: models.Amount{
			Value:    payment.Amount.Value,
			Currency: payment.Amount.Currency,
		},
		Confirmation: confirmation,
		CreatedAt:    payment.CreatedAt,
		Description:  payment.Description,
		Metadata: models.MetaData{
			UserId:   payment.Metadata.UserId,
			PlanType: payment.Metadata.PlanType,
		},
		Recipient: models.Recipient{
			AccountId: payment.Recipient.AccountId,
			GatewayId: payment.Recipient.GatewayId,
		},
		Refundable: payment.Refundable,
		Test:       payment.Test,
	}, nil
}
