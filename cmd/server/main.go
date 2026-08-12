package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/handlers"
	"github.com/Simon-Weij/allium/internal/middleware"
)

const (
	readTimeout                    = 5 * time.Second
	writeTimeout                   = 10 * time.Second
	idleTimeout                    = 120 * time.Second
	dataDirPermissions os.FileMode = 0o750
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(cfg.Data, dataDirPermissions); err != nil {
		panic("could not create data directory " + cfg.Data + ": " + err.Error())
	}

	server := handlers.NewServer(*cfg)

	router := chi.NewRouter()
	router.Use(middleware.WithLogging)
	router.Use(middleware.Authenticate(*cfg))

	router.Route("/rest", func(router chi.Router) {
		// System
		router.Get("/ping.view", server.HandlePing)
		router.Get("/getLicense.view", server.HandleGetLicense)
		router.Get("/getOpenSubsonicExtensions.view", server.HandleGetOpenSubsonicExtensions)

		// User management
		router.Get("/getUser.view", server.HandleGetUser)

		// Clarification
		router.Get("/search3.view", server.HandleSearch3)
		router.Get("/getCoverArt.view", server.HandleGetCoverArt)
	})

	slog.Info("starting app...")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
