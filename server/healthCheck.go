package server

import (
	"github.com/korvised/ilog-shop/pkg/response"
	"github.com/labstack/echo/v4"
)

type healthCheck struct {
	App    string `json:"app"`
	Status string `json:"status"`
}

func (s *server) healthCheckService(c echo.Context) error {
	return response.SuccessResponse(c, 200, &healthCheck{
		App:    s.cfg.App.Name,
		Status: "OK",
	})
}
