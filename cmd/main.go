package main

import (
	_ "github.com/mohdjishin/SplitWise/docs"
	"github.com/mohdjishin/SplitWise/internal/app"
	_ "github.com/mohdjishin/SplitWise/internal/db"
	log "github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap/zapcore"
)

// @title SplitWise API
// @version 1.0
// @description This is an API for managing splits.
// @host localhost:8080
// @BasePath /
func main() {
	app := app.New()
	if err := app.Run(); err != nil {
		log.Error("Error starting server", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
	}

}
