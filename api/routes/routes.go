package routes

import (
	"github.com/alexshelto/tigres-tracker/api/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	r.Get("/", handlers.Index)
	r.Post("/song/play", handlers.AddSongPlay)
}
