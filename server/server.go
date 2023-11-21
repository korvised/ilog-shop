package server

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/middleware/middlewareHandlers"
	"github.com/korvised/ilog-shop/modules/middleware/middlewareRepositories"
	"github.com/korvised/ilog-shop/modules/middleware/middlewareUsecases"
	"github.com/korvised/ilog-shop/pkg/jwtauth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	server struct {
		app *echo.Echo
		db  *mongo.Client
		cfg *config.Config
		m   middlewareHandlers.MiddlewareHandlerService
	}
)

func newMiddleware(cfg *config.Config) middlewareHandlers.MiddlewareHandlerService {
	repo := middlewareRepositories.NewMiddlewareRepository(cfg)
	usecase := middlewareUsecases.NewMiddlewareUsecase(cfg, repo)
	return middlewareHandlers.NewMiddlewareHandler(cfg, usecase)
}

func (s *server) graceFullShutdown(c context.Context, quit <-chan os.Signal) {
	log.Printf("Start sevice: %s \n", s.cfg.App.Name)

	<-quit
	log.Printf("Shutting down service: %s \n", s.cfg.App.Name)

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}

func (s *server) httpListening() {
	if err := s.app.Start(s.cfg.App.Url); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error: %+v", err)
	}
}

func Start(c context.Context, cfg *config.Config, db *mongo.Client) {
	s := &server{
		app: echo.New(),
		db:  db,
		cfg: cfg,
		m:   newMiddleware(cfg),
	}

	jwtauth.SetApiKey(cfg.Jwt.ApiSecretKey)

	// Basic middleware
	// Request TimeOut
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	// CORS
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	// Body Limit
	s.app.Use(middleware.BodyLimit("10M"))

	switch s.cfg.App.Name {
	case "auth":
		s.authService()
	case "player":
		s.playerService()
	case "item":
		s.itemService()
	case "inventory":
		s.inventoryService()
	case "payment":
		s.paymentService()
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Apply Logger
	s.app.Use(middleware.Logger())

	go s.graceFullShutdown(c, quit)

	// Listen and Serve
	s.httpListening()
}
