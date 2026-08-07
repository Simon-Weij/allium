package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/handlers"
)

func main() {
	cfg := config.NewConfig()

	router := chi.NewRouter()
	r := chiopenapi.NewRouter(router,
		option.WithOpenAPIVersion(openapi.Version320),
		option.WithTitle(cfg.Name),
		option.WithVersion(cfg.Version),
		option.WithStoplightElements(),
	)

	r.Route("/rest", func(r chiopenapi.Router) {
		systemHandler := handlers.NewSystemHandler(cfg)
		r.Get("/ping", systemHandler.HandlePing).With(
			option.Summary("Ping"),
			option.Request(new(handlers.SubsonicResponse)),
		)
		r.Get("/getLicense", systemHandler.HandleGetLicense).With(
			option.Summary("Get the license"),
			option.Request(new(handlers.LicenseResponse)),
		)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
