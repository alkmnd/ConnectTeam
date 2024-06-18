package models

import (
	"time"
)

type PaymentResponse struct {
	Id           string     `json:"id"`
	Status       string     `json:"status"`
	Paid         bool       `json:"paid"`
	Amount       Amount     `json:"amount"`
	Confirmation Confirmer  `json:"confirmation"`
	CreatedAt    *time.Time `json:"created_at"`
	Description  string     `json:"description"`
	Metadata     MetaData
	Recipient    Recipient `json:"recipient"`
	Refundable   bool      `json:"refundable"`
	Test         bool      `json:"test"`
}

type PaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnUrl string `json:"return_url"`
	} `json:"confirmation"`
	Capture     bool     `json:"capture"`
	Description string   `json:"description"`
	MetaData    MetaData `json:"metadata"`
}
type MetaData struct {
	UserId   string `json:"user_id"`
	PlanType string `json:"plan_type"`
}
type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}
type Confirmer interface {
}

type Recipient struct {
	AccountId string `json:"account_id"`
	GatewayId string `json:"gateway_id"`
}
