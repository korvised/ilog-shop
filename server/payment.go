package server

import (
	"github.com/korvised/ilog-shop/modules/payment/paymentHandlers"
	"github.com/korvised/ilog-shop/modules/payment/paymentRepositories"
	"github.com/korvised/ilog-shop/modules/payment/paymentUsecases"
)

func (s *server) paymentService() {
	repo := paymentRepositories.NewPaymentRepository(s.db)
	usecase := paymentUsecases.NewPaymentUsecase(repo)
	httpHandler := paymentHandlers.NewPaymentHttpHandler(s.cfg, usecase)
	queueHandler := paymentHandlers.NewPaymentQueueHandler(s.cfg, usecase)

	_ = httpHandler
	_ = queueHandler

	router := s.app.Group("/api/v1/payment")

	// Health check
	router.GET("", s.healthCheckService)
}
