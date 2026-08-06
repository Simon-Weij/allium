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
		pingHandler := handlers.NewPingHandler(cfg)
		r.Get("/ping", pingHandler.HandlePing).With(
			option.Summary("Ping"),
			option.Request(new(handlers.SubsonicResponse)),
		)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
