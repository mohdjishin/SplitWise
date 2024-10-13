package app

import (
	"github.com/mohdjishin/SplitWise/config"
	"github.com/mohdjishin/SplitWise/internal/routes"
	"github.com/mohdjishin/SplitWise/internal/server"
)

type App struct {
	server *server.Server
}

func New() *App {
	port := config.GetConfig().Port
	handler := routes.NewRouter()
	server := server.NewServer(port, handler)
	return &App{server: server}
}

func (a *App) Run() error {
	return a.server.Start()
}
