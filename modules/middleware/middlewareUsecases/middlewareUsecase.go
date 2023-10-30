package middlewareUsecases

import "github.com/korvised/ilog-shop/modules/middleware/middlewareRepositories"

type (
	MiddlewareUsecaseService interface {
	}

	middlewareUsecase struct {
		middlewareRepository middlewareRepositories.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepositories.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewareRepository}
}
