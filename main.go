package main

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"context"
	"fmt"
	"github.com/cloudogu/github-forgejo-backup/internal/disk"
	"github.com/cloudogu/github-forgejo-backup/internal/logs"
	"github.com/google/go-github/v74/github"
	"github.com/robfig/cron/v3"
	"log/slog"
	"time"
)

import "github.com/gofri/go-github-pagination/githubpagination"

const ua = "github-forgejo-backup/0.1.0"

func main() {

	err := config.Load()
	if err != nil {
		logs.Fatal("failed loading config", "error", err)
	}

	location, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		logs.Fatal("unable to load timezone", "error", err)
	}

	sched := cron.New(
		cron.WithLocation(location),
		cron.WithLogger(logs.CronLogger{}),
	)

	_, err = sched.AddFunc(config.CronSpec, doRun)
	if err != nil {
		logs.Fatal(err.Error())
	}

	sched.Run()

}

func doRun() {

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
		logs.Fatal("failed creating forgejo client", "error", err)
	}

	apiSettings, _, err := forgejoClient.GetGlobalAPISettings()
	if err != nil {
		logs.Fatal("failed fetching api settings", "error", err)
	}

	// create forgejo orga when missing
	orga, _, err := forgejoClient.CreateOrg(forgejo.CreateOrgOption{
		Name:       config.ForgejoOrga,
		Visibility: forgejo.VisibleTypeLimited,
	})
	if err != nil {
		if err.Error() != fmt.Sprintf("user already exists [name: %s]", config.ForgejoOrga) {
			logs.Fatal("failed creating forgejo orga", "error", err)
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
			slog.Info("missing mirror", "name", githubRepo.GetName())

			diskUsage, err := disk.UsageOf("/")
			if err != nil {
				logs.Fatal("failed reading disk stats", "error", err)
			}

			slog.Info("disk usage", "used", fmt.Sprintf("%.2f%%", diskUsage.UtilizationPct))

			if diskUsage.UtilizationPct >= 90 {
				logs.Fatal("free space is < 10%", "used %", diskUsage.UtilizationPct)
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
		logs.Fatal("failed fetching github repos", "error", err)
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
			logs.Fatal("failed fetching forgejo repos", "error", err)
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
		logs.Fatal("failed creating repo mirror", "error", err)
	}

	slog.Info("created mirror", "name", forgejoRepo.Name)

}
