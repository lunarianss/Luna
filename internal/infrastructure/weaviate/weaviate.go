package weaviate

import (
	"fmt"
	"sync"

	"github.com/lunarianss/Luna/infrastructure/log"
	weaviateBase "github.com/lunarianss/Luna/infrastructure/weaviate"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

var (
	once           sync.Once
	WeaviateClient *weaviate.Client
)

func GetWeaviateClient(opt *options.WeaviateOptions) (*weaviate.Client, error) {

	var (
		err            error
		weaviateClient *weaviate.Client
	)
	once.Do(func() {
		weaviateClient, err = weaviateBase.NewWeaviateClient(opt)

		if err != nil {
			log.Error(err)
		}

		WeaviateClient = weaviateClient
	})

	if WeaviateClient == nil || err != nil {
		return nil, fmt.Errorf("failed to get weaviate client factory, vdbFactory: %+v, error: %w", WeaviateClient, err)
	}
	return WeaviateClient, nil
}
