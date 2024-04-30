package payment_gateway

import "github.com/rvinnie/yookassa-sdk-go/yookassa"

type Config struct {
	ShopId string
	ApiKey string
}

func NewYookassaClient(cfg Config) *yookassa.Client {
	return yookassa.NewClient(cfg.ShopId, cfg.ApiKey)
}
