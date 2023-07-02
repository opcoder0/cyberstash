package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/cyberstash")
	viper.AddConfigPath("$HOME/.cyberstash")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	typesenseServer := viper.GetString("typesense-server")
	typesenseAPIKey := viper.GetString("typesense-api-key")
	fmt.Println("Loaded config", typesenseServer, typesenseAPIKey)
	fmt.Println("Initialising client...")
	client := typesense.NewClient(
		typesense.WithServer(typesenseServer),
		typesense.WithAPIKey(typesenseAPIKey),
	)
	schema := &api.CollectionSchema{
		Name: "references",
		Fields: []api.Field{
			{
				Name: "title",
				Type: "string",
			},
			{
				Name: "description",
				Type: "string",
			},
			{
				Name: "url",
				Type: "string",
			},
		},
	}
	if _, err = client.Collections().Create(schema); err != nil {
		fmt.Println("Schema not created", err)
		os.Exit(1)
	}
	// fmt.Println("Created schema: ", collectionResponse.Name)
}
