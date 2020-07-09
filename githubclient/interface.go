package githubclient

import "context"

var _ Interface = (*Client)(nil)

type Interface interface {
	FetchConfig(ctx context.Context, owner, repo, path string) (string, error)
}
