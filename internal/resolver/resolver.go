package resolver

import (
	"github.com/Simon-Weij/allium/internal/config"
	"resty.dev/v3"
)

type Resolver struct {
	client     *resty.Client
	cfg        config.Config
	downloader Downloader
}

func NewResolver(cfg config.Config) *Resolver {
	resolver := &Resolver{
		client: resty.New(),
		cfg:    cfg,

		downloader: nil,
	}
	resolver.downloader = resolver

	return resolver
}
