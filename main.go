package main

import (
	"bufio"
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
		fmt.Println("Create Schema:", err)
	}
	glossary, err := os.Open("./stash/glossary/glossary.jsonl")
	if err != nil {
		fmt.Println("Opening glossary:", err)
		os.Exit(1)
	}
	create := "create"
	batchSize := 5
	params := &api.ImportDocumentsParams{
		Action:    &create,
		BatchSize: &batchSize,
	}
	_, err = client.Collection(schema.Name).Documents().ImportJsonl(glossary, params)
	if err != nil {
		fmt.Println("Import Glossary Json:", err)
	}
	fmt.Println("Done.")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text to search:")
	text, _ := reader.ReadString('\n')
	fmt.Printf("Searching %s\n", text)
	searchParams := &api.SearchCollectionParams{
		Q:       text,
		QueryBy: "title, description",
	}
	searchResult, err := client.Collection(schema.Name).Documents().Search(searchParams)
	if err != nil {
		fmt.Printf("Search %s returned: %v", text, err)
		os.Exit(1)
	}
	if searchResult.Found != nil {
		fmt.Printf("Found %d documents\n", *searchResult.Found)
	}
}
