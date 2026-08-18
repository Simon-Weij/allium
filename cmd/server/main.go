package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/database"
	"github.com/Simon-Weij/allium/internal/handlers"
	"github.com/Simon-Weij/allium/internal/metadata"
	"github.com/Simon-Weij/allium/internal/middleware"
	"github.com/Simon-Weij/allium/internal/users"
)

const (
	readTimeout                    = 45 * time.Second
	writeTimeout                   = 30 * time.Second
	idleTimeout                    = 120 * time.Second
	contextTimeout                 = 5 * time.Second
	dataDirPermissions os.FileMode = 0o750
)

//go:generate sqlc generate -f ../../sqlc.yml
func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	cfg, err := config.ParseConfig()
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}

	db, err := database.InitialiseDatabase(ctx, cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to initialise database: %w", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := os.MkdirAll(cfg.Data, dataDirPermissions); err != nil {
		return fmt.Errorf("could not create data directory %s: %w", cfg.Data, err)
	}

	queries := sqlc.New(db)

	server := handlers.NewServer(*cfg, queries, metadata.NewMetadata(*cfg), users.NewUserClient(queries))

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
		router.Get("/getUsers.view", server.HandleGetUsers)

		// Clarification
		router.Get("/search3.view", server.HandleSearch3)
		router.Get("/getCoverArt.view", server.HandleGetCoverArt)
		router.Get("/stream.view", server.HandleStream)
		router.Get("/stream.view/", server.HandleStream)

		// Browsing
		router.Get("/getAlbum.view", server.HandleGetAlbum)
		router.Get("/getArtist.view", server.HandleGetArtist)

		// Media annotation
		router.Get("/scrobble.view", server.HandleScrobble)
	})

	slog.Info("starting app...")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}

	return nil
}
