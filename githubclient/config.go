package githubclient

import (
	"context"
)

// FetchConfig fetch yaml config
func (c *Client) FetchConfig(ctx context.Context, owner, repo, path string) (string, error) {
	fileContent, _, _, err := c.client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return "", err
	}
	var config string
	if fileContent.Content != nil {
		config, err = fileContent.GetContent()
		if err != nil {
			return "", err
		}
	}
	return config, nil
}
