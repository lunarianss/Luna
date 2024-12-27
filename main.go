package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

type WeaviateConsole struct {
	client *weaviate.Client
}

func NewWeaviateClient() (*weaviate.Client, error) {
	log.NewStdWithOptions()
	cfg := weaviate.Config{
		Host:       "127.0.0.1:8080",
		Scheme:     "http",
		AuthConfig: auth.ApiKey{Value: "WVF5YThaHlkYwhGUSmCRgsX3tD5ngdN8pkih"},
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
		fmt.Println(color.GreenString("weaviate is live!"))
	}

	return client, nil
}

func (wc *WeaviateConsole) GetAllCollectionDefined() {
	context := context.Background()
	schema, err := wc.client.Schema().Getter().
		Do(context)

	if err != nil {
		panic(err)
	}
	fmt.Sprintln("========= All Collection Start ========")
	util.LogCompleteInfo(schema)
	fmt.Sprintln("========= All Collection End   ========")
}

func (wc *WeaviateConsole) GetCollectionObjects(className string) {
	ctx := context.Background()
	client, err := NewWeaviateClient()

	if err != nil {
		panic(err)
	}

	response, err := client.GraphQL().Get().
		WithClassName(className).
		WithFields(graphql.Field{Name: "text"}, graphql.Field{Name: "doc_id"}, graphql.Field{Name: "app_id"}, graphql.Field{Name: "annotation_id"}, graphql.Field{Name: "_additional", Fields: []graphql.Field{
			{Name: "vector"},
			{Name: "id"},
		}}).
		Do(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Printf("========= Collection %s Objects Start ========\n", className)
	util.LogCompleteInfo(response)
	fmt.Printf("========= Collection %s  Objects End   ========\n", className)
}

func NewConsole() *WeaviateConsole {
	client, err := NewWeaviateClient()

	if err != nil {
		panic(err)
	}

	return &WeaviateConsole{
		client,
	}
}
func main() {
	console := NewConsole()

	// console.GetAllCollectionDefined()

	console.GetCollectionObjects("Vector_index_9a532b60_a004_425e_a952_fa5f33c6cdd1_Node")

	var a []string = nil

	fmt.Println(len(a))

	var c = []string{"你", "好"}

	fmt.Println(c[1:])
}
