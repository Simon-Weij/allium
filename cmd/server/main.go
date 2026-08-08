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
	cfg, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.Authenticate(*cfg))

	r.Route("/rest", func(r chi.Router) {
		systemHandler := handlers.NewSystemHandler(*cfg)
		r.Get("/ping.view", systemHandler.HandlePing)
		r.Get("/getLicense.view", systemHandler.HandleGetLicense)
		r.Get("/getOpenSubsonicExtensions.view", systemHandler.HandleGetOpenSubsonicExtensions)

		userManagementHandler := handlers.NewUserManagementHandler(*cfg)
		r.Get("/getUser.view", userManagementHandler.HandleGetUser)

		jukeboxHandler := handlers.NewJukeboxHandler(*cfg)
		r.HandleFunc("/jukeboxControl.view", jukeboxHandler.HandleJukeboxStatus)
		r.Get("/jukeboxControl.view", jukeboxHandler.HandleJukeboxPlaylist)
	})

	slog.Info("starting app...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
