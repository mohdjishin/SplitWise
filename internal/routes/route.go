package routes

import (
	"github.com/go-chi/chi"
	mChi "github.com/go-chi/chi/middleware"
	"github.com/mohdjishin/SplitWise/internal/handlers"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter() (r *chi.Mux) {
	r = chi.NewRouter()
	r.Use(
		mChi.Logger,
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
		})

		r.Route("/payments", func(r chi.Router) {
			r.Post("/", handlers.MarkPayment)
			r.Get("/pending", handlers.GetPendingPayments)

		})

		r.Route("/report", func(r chi.Router) {
			r.Post("/", handlers.GetGroupReport)
			r.Get("/{id}", handlers.GenerateSingleGroupReport)
		})
	})
	return
}
