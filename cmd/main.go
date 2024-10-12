package main

import (
	"net/http"

	"github.com/mohdjishin/SplitWise/config"
	_ "github.com/mohdjishin/SplitWise/docs"
	_ "github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/routes"
	"github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap/zapcore"
)

// @title SplitWise API
// @version 1.0
// @description This is an API for managing splits.
// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		logger.LoggerInstance.Error("Error starting server", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
	}

}

func run() error {
	defer logger.LoggerInstance.Sync()
	port := config.GetConfig().Port
	logger.LoggerInstance.Info("Starting server on port " + port)
	return http.ListenAndServe(":"+port, routes.NewRouter())
}
