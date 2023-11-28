package server

import (
	"github.com/korvised/ilog-shop/modules/payment/paymentHandlers"
	"github.com/korvised/ilog-shop/modules/payment/paymentRepositories"
	"github.com/korvised/ilog-shop/modules/payment/paymentUsecases"
)

func (s *server) paymentService() {
	repo := paymentRepositories.NewPaymentRepository(s.db, s.cfg)
	usecase := paymentUsecases.NewPaymentUsecase(s.cfg, repo)
	httpHandler := paymentHandlers.NewPaymentHttpHandler(s.cfg, usecase)
	queueHandler := paymentHandlers.NewPaymentQueueHandler(s.cfg, usecase)

	_ = queueHandler

	router := s.app.Group("/api/v1")

	// Health check
	router.GET("", s.healthCheckService)

	router.POST("/payment/buy", httpHandler.BuyItem, s.m.Authorization)
	router.POST("/payment/sell", httpHandler.SellItem, s.m.Authorization)
}
