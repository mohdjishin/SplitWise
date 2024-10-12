package routes

import (
	"github.com/go-chi/chi"
	mChi "github.com/go-chi/chi/middleware"
	"github.com/mohdjishin/SplitWise/internal/handlers"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// TODO: add route for specific group report
func NewRouter() (r *chi.Mux) {
	r = chi.NewRouter()
	r.Use(
		mChi.Logger,
		mChi.Recoverer,
		mChi.Logger,
		mChi.RequestID,
		mChi.RealIP,
		mChi.Heartbeat("/ping"),
		mChi.Timeout(1),
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
	return
}
