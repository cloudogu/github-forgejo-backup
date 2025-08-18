package forgejo

import "time"

type ApiError struct {
	Message string   `json:"message"`
	Url     string   `json:"url"`
	Errors  []string `json:"errors"`
}

type ApiSettings struct {
	DefaultGitTreesPerPage int `json:"default_git_trees_per_page"`
	DefaultMaxBlobSize     int `json:"default_max_blob_size"`
	DefaultPagingNum       int `json:"default_paging_num"`
	MaxResponseItems       int `json:"max_response_items"`
}

type MigrateRepoOptions struct {
	AuthPassword   int `json:"auth_password"`
	AuthToken      int `json:"auth_token"`
	AuthUsername   int `json:"auth_username"`
	CloneAddr      int `json:"clone_addr"`
	Mirror         int `json:"mirror"`
	MirrorInterval int `json:"mirror_interval"`
	RepoName       int `json:"repo_name"`
	RepoOwner      int `json:"repo_owner"`
	Wiki           int `json:"wiki"`
}

type Repo struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	Description    string    `json:"description"`
	Empty          bool      `json:"empty"`
	Mirror         bool      `json:"mirror"`
	Size           int       `json:"size"`
	Url            string    `json:"url"`
	Link           string    `json:"link"`
	SshUrl         string    `json:"ssh_url"`
	CloneUrl       string    `json:"clone_url"`
	OriginalUrl    string    `json:"original_url"`
	DefaultBranch  string    `json:"default_branch"`
	Archived       bool      `json:"archived"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ArchivedAt     time.Time `json:"archived_at"`
	MirrorInterval string    `json:"mirror_interval"`
	MirrorUpdated  time.Time `json:"mirror_updated"`
}
