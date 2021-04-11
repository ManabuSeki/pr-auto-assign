package githubclient

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"gopkg.in/yaml.v2"
)

// ReviewConfig  config struct
type ReviewConfig struct {
	MustReviewers     []string `yaml:"must_reviewers"`
	Reviewers         []string `yaml:"reviewers"`
	NumberOfReviewers int      `yaml:"number_of_reviewers"`
	AddAssignees      bool     `yaml:"add_assignees"`
}

// FetchConfig fetch yaml config
func (c *Client) FetchConfig(ctx context.Context, owner, repo, path, ref string) (*ReviewConfig, error) {
	ops := github.RepositoryContentGetOptions{
		Ref: ref,
	}
	fileContent, _, _, err := c.client.Repositories.GetContents(ctx, owner, repo, path, &ops)
	if err != nil {
		return nil, err
	}

	if fileContent.Content == nil {
		return nil, fmt.Errorf("config file  is not found")
	}

	confStr, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}
	config, err := c.parseConfig(ctx, confStr)
	return config, nil
}

func (c *Client) parseConfig(ctx context.Context, config string) (*ReviewConfig, error) {
	var conf ReviewConfig
	err := yaml.Unmarshal([]byte(config), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
