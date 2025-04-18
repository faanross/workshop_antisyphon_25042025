package router

import (
	"github.com/faanross/orlokC2/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	r.Use(middleware.UUIDMiddleware)

	r.Get("/command", CommandHandler)

	r.Post("/result", ResultHandler)
}
