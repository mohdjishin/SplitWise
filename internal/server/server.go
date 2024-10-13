package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		<-stop
		s.Shutdown()
	}()

	log.Info("Starting server on port " + s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("Server stopped gracefully")
	return nil
}

func (s *Server) Shutdown() {
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = log.Sync()

	log.Info("Logger synced...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error("Error during server shutdown", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
	}
}
