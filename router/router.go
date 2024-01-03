package router

import (
	"awesomeProject/internal/handler"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRouter(userHandler *handler.UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Post("/api/users", userHandler.Create)
	router.Get("/api/users/{id}", userHandler.GetById)
	router.Put("/api/users/{id}", userHandler.Update)
	router.Delete("/api/users/{id}", userHandler.Delete)

	router.Get("/api/users", userHandler.List)

	return router
}
