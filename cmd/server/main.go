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
	r.Use(middleware.WithLogging)
	r.Use(middleware.Authenticate(cfg))

	r.Route("/rest", func(r chi.Router) {
		systemHandler := handlers.NewSystemHandler(cfg)
		r.Get("/ping.view", systemHandler.HandlePing)
		r.Get("/getLicense.view", systemHandler.HandleGetLicense)

		userManagementHandler := handlers.NewUserManagementHandler(cfg)
		r.Get("/getUser.view", userManagementHandler.HandleGetUser)
	})

	slog.Info("starting app...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
