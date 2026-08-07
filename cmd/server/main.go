package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/handlers"
	"github.com/Simon-Weij/allium/internal/middleware"
)

func main() {
	cfg := config.NewConfig()

	r := chi.NewRouter()
	r.Use(middleware.RequireApiKey(cfg))

	r.Route("/rest", func(r chi.Router) {
		systemHandler := handlers.NewSystemHandler(cfg)
		r.Get("/ping", systemHandler.HandlePing)
		r.Get("/getLicense", systemHandler.HandleGetLicense)
	})

	slog.Info("starting app...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
