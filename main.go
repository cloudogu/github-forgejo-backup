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

	// fetch a list of all github repos in the github "cloudogu" orga

	githubClient := github.NewClient(nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancel()

	var githubRepos []*github.Repository
	{
		githubReposPage, _, err := githubClient.Repositories.ListByOrg(ctx, "cloudogu",
			&github.RepositoryListByOrgOptions{Sort: "full_name"})
		if err != nil {
			logs.Error("failed listing repos", "error", err)
			os.Exit(1)
		}
		githubRepos = append(githubRepos, githubReposPage...)
		// TODO: fetch more pages
	}

	fmt.Println("### github repos")
	for _, repo := range githubRepos {
		fmt.Printf("%s\t%s\n", repo.GetName(), repo.GetCloneURL())
	}

	// fetch a list of all forgejo repos in the forgejo "cloudogu" orga

	forgejoClient, err := forgejo.NewClient(config.ForgejoBaseUrl,
		forgejo.SetToken(config.ForgejoToken),
		forgejo.SetUserAgent("github-forgejo-backup/0.1.0"),
	)
	if err != nil {
		logs.Error("failed creating forgejo client", "error", err)
		os.Exit(1)
	}

	apiSettings, _, err := forgejoClient.GetGlobalAPISettings()
	if err != nil {
		logs.Error("failed fetching api settings", "error", err)
		os.Exit(1)
	}

	forgejoOrganisation, _, err := forgejoClient.GetOrg("cloudogu")
	if err != nil {
		logs.Error("failed fetching organisation", "error", err)
		os.Exit(1)
	}

	var forgejoRepos []*forgejo.Repository
	{
		forgejoReposPage, _, err := forgejoClient.ListOrgRepos(forgejoOrganisation.UserName,
			forgejo.ListOrgReposOptions{
				ListOptions: forgejo.ListOptions{
					Page:     0,
					PageSize: apiSettings.MaxResponseItems,
				},
			})
		if err != nil {
			logs.Error("failed fetching forgejo repos", "error", err)
			os.Exit(1)
		}
		forgejoRepos = append(forgejoRepos, forgejoReposPage...)
		// TODO: fetch more pages
	}

	fmt.Println("### forgejo repos")
	for _, repo := range forgejoRepos {
		fmt.Printf("%s\t%s\n", repo.Name, repo.CloneURL)
	}

	// TODO: create a new mirror in forgejo in the forgejo "cloudogu" orga for all github repos in the github "cloudogu" orga missing in the forgejo "cloudogu" orga

}
