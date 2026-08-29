package resolver

import (
	"github.com/Simon-Weij/allium/internal/config"
	"resty.dev/v3"
)

type Metadata struct {
	client     *resty.Client
	cfg        config.Config
	downloader Downloader
}

func NewMetadata(cfg config.Config) *Metadata {
	metadata := &Metadata{
		client: resty.New(),
		cfg:    cfg,

		downloader: nil,
	}
	metadata.downloader = metadata

	return metadata
}
