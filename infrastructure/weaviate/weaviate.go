package weaviate

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

func NewWeaviateClient(opt *options.WeaviateOptions) (*weaviate.Client, error) {
	cfg := weaviate.Config{
		Host:       opt.Endpoint,
		Scheme:     opt.Schema,
		AuthConfig: auth.ApiKey{Value: opt.ApiKey},
		Timeout:    60 * time.Second,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	live, err := client.Misc().LiveChecker().Do(context.Background())

	if err != nil {
		return nil, err
	}

	if live {
		log.Info(color.GreenString("weaviate is live!"))
	}

	return client, nil
}
