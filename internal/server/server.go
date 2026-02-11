package server

import (
	h "chopper/internal/delivery/http"
	"chopper/internal/domain"
	"chopper/internal/middleware"
	"chopper/internal/usecase"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server            *http.Server
	timeoutToShutdown time.Duration
}

func NewServer(address string, readTimeout, writeTimeout, idleTimeout, timeoutToShutdown time.Duration, serverMode domain.ServerMode, userService *usecase.UserService, dailyNotesService *usecase.DailyNotesService, alertService *usecase.AlertService, authMiddleware *middleware.AuthMiddleware, rateLimiter *middleware.RateLimiter) *Server {
	// создание gin core
	gin.SetMode(string(serverMode))
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// users public
	usersPublic := r.Group("/users")
	usersPublic.Use(rateLimiter.RateLimit())

	// users protected
	usersProtected := r.Group("/users")
	usersProtected.Use(authMiddleware.Auth())
	usersProtected.Use(rateLimiter.RateLimit())

	// notes protected
	notesProtected := r.Group("/notes")
	notesProtected.Use(authMiddleware.Auth())
	notesProtected.Use(rateLimiter.RateLimit())

	// alert protected
	alertProtected := r.Group("/alert")
	alertProtected.Use(authMiddleware.Auth())
	alertProtected.Use(rateLimiter.RateLimit())

	userHandler := h.NewUserHandler(userService)
	userHandler.RegisterRoutes(usersPublic, usersProtected)
	noteHandler := h.NewNoteHandler(dailyNotesService)
	noteHandler.RegisterRoutes(notesProtected)
	alertHandler := h.NewAlertHandler(alertService)
	alertHandler.RegisterRoutes(alertProtected)

	server := &http.Server{
		Addr:         address,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      r,
	}
	return &Server{
		server:            server,
		timeoutToShutdown: timeoutToShutdown,
	}
}

func (s *Server) StartServer() error {
	address := s.server.Addr
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		fmt.Printf("server has been started on %v\n", address)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err) // заменить в будущем
		}
	}()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), s.timeoutToShutdown)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
