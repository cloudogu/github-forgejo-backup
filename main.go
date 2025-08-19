package main

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"context"
	"fmt"
	"github.com/cloudogu/github-forgejo-backup/internal/disk"
	"github.com/google/go-github/v74/github"
	"os"
	"time"
)

import "github.com/gofri/go-github-pagination/githubpagination"

const ua = "github-forgejo-backup/0.1.0"

func main() {

	err := config.Load()
	if err != nil {
		logs.Error("failed loading config", "error", err)
		os.Exit(1)
	}

	githubClient := github.NewClient(
		githubpagination.NewClient(nil,
			githubpagination.WithPerPage(100)),
	).WithAuthToken(config.GithubToken)
	githubClient.UserAgent = ua

	githubRepos := ListAllGithubRepos(githubClient)

	forgejoClient, err := forgejo.NewClient(config.ForgejoBaseUrl,
		forgejo.SetToken(config.ForgejoToken),
		forgejo.SetUserAgent(ua),
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

	// create forgejo orga when missing
	orga, _, err := forgejoClient.CreateOrg(forgejo.CreateOrgOption{
		Name:       config.ForgejoOrga,
		Visibility: forgejo.VisibleTypeLimited,
	})
	if err != nil {
		if err.Error() != fmt.Sprintf("user already exists [name: %s]", config.ForgejoOrga) {
			logs.Error("failed creating forgejo orga", "error", err)
			os.Exit(1)
		}
	} else {
		logs.Info("created orga", "name", orga.UserName)
	}

	// fetch all forgejo repos
	forgejoRepos := ListAllForgejoRepos(forgejoClient, apiSettings)

	// check for and create missing mirrors

	for _, githubRepo := range githubRepos {
		found := false
		for _, forgejoRepo := range forgejoRepos {
			if forgejoRepo.Name == githubRepo.GetName() {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("\nmissing mirror: %s\n", githubRepo.GetName())

			diskUsage, err := disk.UsageOf("/")
			if err != nil {
				logs.Error("failed reading disk stats", "error", err)
				os.Exit(1)
			}

			fmt.Printf("disk usage: %.2f%%\n", diskUsage.UtilizationPct)
			if diskUsage.UtilizationPct >= 90 {
				logs.Error("free space is < 10%", "used %", diskUsage.UtilizationPct)
				os.Exit(1)
			}

			CreateMirror(forgejoClient, githubRepo)
		}
	}

}

func ListAllGithubRepos(client *github.Client) (repos []*github.Repository) {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*1)

	pageRepos, _, err := client.Repositories.ListByOrg(
		ctx,
		"cloudogu",
		&github.RepositoryListByOrgOptions{Sort: "full_name"})

	cancel()

	if err != nil {
		logs.Error("failed fetching github repos", "error", err)
		os.Exit(1)
	}

	repos = append(repos, pageRepos...)

	return repos
}

func ListAllForgejoRepos(client *forgejo.Client, apiSettings *forgejo.GlobalAPISettings) (repos []*forgejo.Repository) {

	page := 1
	for {

		pageRepos, responseInfos, err := client.ListOrgRepos(
			config.ForgejoOrga,
			forgejo.ListOrgReposOptions{
				ListOptions: forgejo.ListOptions{
					Page:     page,
					PageSize: apiSettings.MaxResponseItems,
				},
			})

		if err != nil {
			logs.Error("failed fetching forgejo repos", "error", err)
			os.Exit(1)
		}

		repos = append(repos, pageRepos...)

		if page >= responseInfos.NextPage || len(pageRepos) == 0 {
			break
		}

		page = page + 1
	}

	return repos
}

func CreateMirror(client *forgejo.Client, githubRepo *github.Repository) {

	forgejoRepo, _, err := client.MigrateRepo(forgejo.MigrateRepoOption{
		RepoName:       githubRepo.GetName(),
		RepoOwner:      config.ForgejoOrga,
		CloneAddr:      githubRepo.GetCloneURL(),
		Service:        forgejo.GitServiceGithub,
		AuthToken:      config.GithubToken,
		Mirror:         true,
		Private:        true,
		Wiki:           true,
		Milestones:     true,
		Labels:         true,
		Issues:         true,
		PullRequests:   true,
		Releases:       true,
		MirrorInterval: "1h",
	})
	if err != nil {
		logs.Error("failed creating repo mirror", "error", err)
		os.Exit(1)
	}

	fmt.Printf("created mirror: %s\n", forgejoRepo.Name)

}
