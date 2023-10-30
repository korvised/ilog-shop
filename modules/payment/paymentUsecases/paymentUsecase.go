package paymentUsecases

import "github.com/korvised/ilog-shop/modules/payment/paymentRepositories"

type (
	PaymentUsecaseService interface {
	}

	paymentUsecase struct {
		paymentRepository paymentRepositories.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepository paymentRepositories.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{paymentRepository}
}
