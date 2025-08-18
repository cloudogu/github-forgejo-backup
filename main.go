package main

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"context"
	"fmt"
	"github.com/google/go-github/v74/github"
	"os"
	"time"
)

func main() {

	err := config.Load()
	if err != nil {
		logs.Error("failed loading config", "error", err)
		os.Exit(1)
	}

	githubClient := github.NewClient(nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancel()

	repos, _, err := githubClient.Repositories.ListByOrg(ctx, "cloudogu", &github.RepositoryListByOrgOptions{Sort: "full_name"})
	if err != nil {
		logs.Error("failed listing repos", "error", err)
		os.Exit(1)
	}

	for _, repo := range repos {
		fmt.Printf("%s\t%s\n", repo.GetName(), repo.GetCloneURL())
	}

	forgejoClient, err := forgejo.NewClient(config.ForgejoBaseUrl,
		forgejo.SetToken(config.ForgejoToken),
		forgejo.SetUserAgent("github-forgejo-backup/0.1.0"),
	)
	if err != nil {
		logs.Error("failed creating forgejo client", "error", err)
		os.Exit(1)
	}

	organisation, _, err := forgejoClient.GetOrg("cloudogu")
	if err != nil {
		logs.Error("failed fetching organisation", "error", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", organisation.UserName)

}
