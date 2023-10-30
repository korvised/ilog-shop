package paymentHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/payment/paymentUsecases"
)

type (
	PaymentHttpHandlerService interface {
	}

	paymentHttpHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecases.PaymentUsecaseService
	}
)

func NewPaymentHttpHandler(cfg *config.Config, paymentUsecase paymentUsecases.PaymentUsecaseService) PaymentHttpHandlerService {
	return &paymentHttpHandler{cfg, paymentUsecase}
}
