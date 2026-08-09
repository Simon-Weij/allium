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

	server := handlers.NewServer(*cfg)

	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.Authenticate(*cfg))

	r.Route("/rest", func(r chi.Router) {
		// System
		r.Get("/ping.view", server.HandlePing)
		r.Get("/getLicense.view", server.HandleGetLicense)
		r.Get("/getOpenSubsonicExtensions.view", server.HandleGetOpenSubsonicExtensions)

		// User management
		r.Get("/getUser.view", server.HandleGetUser)

		// Browsing
		r.Get("/getMusicFolders.view", server.HandleGetMusicFolders)
		r.Get("/getGenres.view", server.HandleGetGenres)
		r.Get("/getCoverArt.view", server.HandleGetCoverArt)

		// Playlists
		r.Get("/getPlaylists.view", server.HandleGetPlaylists)

		// Jukebox
		r.HandleFunc("/jukeboxControl.view", server.HandleJukeboxStatus)
		r.Get("/jukeboxControl.view", server.HandleJukeboxPlaylist)

		// Lists
		r.Get("/getAlbumList2.view", server.HandleGetAlbumList2)
	})

	slog.Info("starting app...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
