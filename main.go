package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

	repos, err := fetchOrgRepos(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching repositories: %v\n", err)
		os.Exit(1)
	}

	for _, r := range repos {
		fmt.Printf("%s\t%s\n", r.FullName, r.HtmlUrl)
	}
}

func fetchOrgRepos(ctx context.Context) ([]Repository, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	perPage := 100
	page := 1
	var all []Repository

	for {
		url := fmt.Sprintf("https://api.github.com/orgs/cloudogu/repos?per_page=%d&page=%d", perPage, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "token "+config.GithubToken)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "github-forgejo-backup/0.1.0")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		func() {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
				err = fmt.Errorf("GitHub API error: %s: %s", resp.Status, strings.TrimSpace(string(body)))
				return
			}

			var batch []Repository
			if decErr := json.NewDecoder(resp.Body).Decode(&batch); decErr != nil {
				err = decErr
				return
			}

			all = append(all, batch...)
			if len(batch) < perPage {
				// Last page
				err = nil
				return
			}
			// Continue to next page
			err = nil
		}()
		if err != nil {
			return nil, err
		}

		// If we appended fewer than perPage, we're done (handled above)
		if len(all) < page*perPage {
			break
		}
		page++
	}

	return all, nil
}
