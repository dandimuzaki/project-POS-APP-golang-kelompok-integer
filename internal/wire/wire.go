package wire

import (
	"travel-api/internal/adaptor"
	"travel-api/internal/data/repository"
	"travel-api/internal/usecase"
	"travel-api/pkg/utils"

	mCustom "travel-api/pkg/middleware"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	Route *chi.Mux
}

func Wiring(repo *repository.Repository, log *zap.Logger, config utils.Configuration) *App {
	r := chi.NewRouter()

	usecase := usecase.NewUsecase(repo, log)
	handler := adaptor.NewHandler(usecase, log, config)
	mw := mCustom.NewMiddlewareCustom(usecase, log)
	r.Mount("/api/v1", ApiV1(&handler, mw))

	return &App{
		Route: r,
	}
}

func ApiV1(handler *adaptor.Handler, mw mCustom.MiddlewareCustom) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(mw.Logging)

	r.Route("/tours", func(r chi.Router) {
		r.Get("/", handler.TourHandler.GetAllTours)
	})

	r.Route("/schedules", func(r chi.Router) {
		r.Get("/{id}", handler.TourHandler.GetTourDetails)
	})
	
	return r
}