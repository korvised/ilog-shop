package paymentHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/payment/paymentUsecases"
)

type (
	PaymentQueueHandlerService interface {
	}

	paymentQueueHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecases.PaymentUsecaseService
	}
)

func NewPaymentQueueHandler(cfg *config.Config, paymentUsecase paymentUsecases.PaymentUsecaseService) PaymentQueueHandlerService {
	return &paymentQueueHandler{cfg, paymentUsecase}
}
