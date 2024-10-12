package main

import (
	"net/http"

	"github.com/mohdjishin/SplitWise/config"
	_ "github.com/mohdjishin/SplitWise/docs"
	_ "github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/routes"
	log "github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap/zapcore"
)

// @title SplitWise API
// @version 1.0
// @description This is an API for managing splits.
// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		log.Error("Error starting server", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
		return
	}
}

func run() error {
	defer log.Sync()

	port := config.GetConfig().Port
	log.Info("Starting server on port " + port)

	serverAddr := ":" + port
	if err := http.ListenAndServe(serverAddr, routes.NewRouter()); err != nil {
		return err
	}
	return nil
}
