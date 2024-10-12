package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	mChi "github.com/go-chi/chi/middleware"
	"github.com/mohdjishin/SplitWise/config"
	_ "github.com/mohdjishin/SplitWise/docs"
	_ "github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/handlers"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	"github.com/mohdjishin/SplitWise/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SplitWise API
// @version 1.0
// @description This is an API for managing splits.
// @host localhost:8080
// @BasePath /
func main() {
	r := chi.NewRouter()
	defer logger.LoggerInstance.Sync()
	// TODO:move this route to a separate file
	r.Use(
		mChi.Recoverer,
		mChi.Logger,
		mChi.RequestID,
		mChi.RealIP,
		mChi.Heartbeat("/ping"),
	)

	r.Post("/auth/register", handlers.Register)
	r.Post("/auth/login", handlers.Login)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/groups", func(r chi.Router) {
			r.Post("/", handlers.CreateGroupWithBill)
			r.Delete("/{id}", handlers.DeleteGroup)
			r.Get("/owned", handlers.ListOwnedGroups)
			r.Post("/{id}/addMembers", handlers.AddUsersToGroup)
			r.Get("/member-groups", handlers.ListMemberGroups)

			r.Post("/report", handlers.GetGroupReport)
			// r.Get("/groups/{groupID}/report", GenerateSingleGroupReport) // Specific group report based on groupID

		})

		r.Route("/payments", func(r chi.Router) {
			r.Post("/", handlers.MarkPayment)
		})
	})

	port := config.GetConfig().Port
	logger.LoggerInstance.Info("Starting server on port " + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
