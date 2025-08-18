package main

import (
	"context"
	"fmt"
	"github-forgejo-backup/internal/forgejo"
	"log"
	"os"
	"time"
)

type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HtmlUrl     string `json:"html_url"`
	Language    string `json:"language"`
}

func main() {

	err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fgClient, err := forgejo.NewClient(config.ForgejoBaseUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	apiSettings, err := fgClient.GetAPISettings(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(apiSettings)

}
