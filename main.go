package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	ctx := context.Background()
	cfg := weaviate.Config{
		Host:   "2dpdqizcszm2ro5npi3ilq.c0.us-west3.gcp.weaviate.cloud", // Replace with your Weaviate URL
		Scheme: "https",
		AuthConfig: auth.ApiKey{
			Value: "4lSsrDa2yQcojyZpDZqgDFb4jPVyKDKvaPZv",
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Check the connection
	ready, err := client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", ready)

	// Define the collection
	classObj := &models.Class{
		Class: "Question",
	}

	isExist, err := client.Schema().ClassExistenceChecker().WithClassName(classObj.Class).Do(ctx)

	if err != nil {
		panic(err)
	}

	if !isExist {
		// add the collection
		err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	w, err := client.Data().Creator().
		WithClassName(classObj.Class).
		WithProperties(map[string]interface{}{
			"question":    "12 with vector This vector DB is OSS and supports automatic property type inference on import",
			"answer":      "Weaviate", // schema properties can be omitted
			"newProperty": 123,        // will be automatically added as a number property
		}).WithVector([]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.01, 0.02, 0.1}).
		Do(ctx)

	if err != nil {
		panic(err)
	}
	// the returned value is a wrapped object
	b, err := json.MarshalIndent(w.Object, "", "")

	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
