package githubclient

import (
	"context"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// Client Github Client
type Client struct {
	client *github.Client
}

func NewClient(ctx context.Context, githubToken string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
	}
}
