package forgejo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client
}

func NewClient(baseUrl string) (*Client, error) {

	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		baseUrl:    u,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil

}

func (c *Client) applyDefaultHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "github-forgejo-backup/0.1.0")
	req.Header.Set("Accept", "application/json")
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
}

func (c *Client) GetAPISettings(ctx context.Context) (data *ApiSettings, err error) {

	u, err := url.JoinPath(c.baseUrl.String(), "settings", "api")
	if err != nil {
		return nil, fmt.Errorf("failed building api url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	c.applyDefaultHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:

		err = json.NewDecoder(resp.Body).Decode(data)
		if err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		} else {
			return data, nil
		}

	default:
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

}

func (c *Client) GetRepo(ctx context.Context, owner string, repo string) (data *Repo, err error) {

	u, err := url.JoinPath(c.baseUrl.String(), "repos", owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed building api url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	c.applyDefaultHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:

		err = json.NewDecoder(resp.Body).Decode(data)
		if err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		} else {
			return data, nil
		}

	case http.StatusNotFound:

		var apiError *ApiError
		err = json.NewDecoder(resp.Body).Decode(apiError)
		if err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		} else {
			return data, fmt.Errorf("repo not found: %s", apiError)
		}

	default:
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

}

func (c *Client) MigrateRepo(ctx context.Context, opts MigrateRepoOptions) (data *Repo, err error) {

	u, err := url.JoinPath(c.baseUrl.String(), "migrate")
	if err != nil {
		return nil, fmt.Errorf("failed building api url: %w", err)
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(opts)
	if err != nil {
		return nil, fmt.Errorf("failed parsing options: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	c.applyDefaultHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:

		err = json.NewDecoder(resp.Body).Decode(data)
		if err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		} else {
			return data, nil
		}

	case http.StatusForbidden:

		var apiError *ApiError
		err = json.NewDecoder(resp.Body).Decode(apiError)
		if err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		} else {
			return data, fmt.Errorf("action forbidden: %s", apiError)
		}

	case http.StatusConflict:

		return data, fmt.Errorf("repository with same name exists")

	case http.StatusRequestEntityTooLarge:

		return data, fmt.Errorf("quota exceeded: %s (%s %s)", resp.Header.Get("message"), resp.Header.Get("username"), resp.Header.Get("user_id"))

	case http.StatusUnprocessableEntity:

		var apiError *ApiError
		err = json.NewDecoder(resp.Body).Decode(apiError)
		if err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		} else {
			return data, fmt.Errorf("repo not found: %s", apiError)
		}

	default:
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

}
