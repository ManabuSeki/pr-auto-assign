package githubclient

import "context"

var _ Interface = (*Client)(nil)

type Interface interface {
	FetchConfig(ctx context.Context, owner, repo, path, ref string) (*ReviewConfig, error)
	HandlePullRequest(ctx context.Context, owner, repo string, prID int, config *ReviewConfig) error
}
