package server

import (
	h "chopper/internal/delivery/http"
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

func NewServer(address string, readTimeout, writeTimeout, idleTimeout, timeoutToShutdown time.Duration, userService *usecase.UserService) *Server {
	// создание gin core
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	usersGroup := r.Group("/users")

	userHandler := h.NewUserHandler(userService)
	userHandler.RegisterRoutes(usersGroup)
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
