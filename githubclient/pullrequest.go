package githubclient

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
)

func (c *Client) HandlePullRequest(ctx context.Context, owner, repo string, prID int, config *ReviewConfig) error {
	// get pullrequest
	pr, _, err := c.client.PullRequests.Get(ctx, owner, repo, prID)
	if err != nil {
		return err
	}

	config = c.removeCreateUser(ctx, config, *pr.GetUser().Login)
	reviewer := c.selectReviewer(ctx, config)
	if err := c.addReviewers(ctx, owner, repo, prID, reviewer); err != nil {
		return err
	}
	// set assignees
	if config.AddAssignees {
		err := c.addAssignees(ctx, owner, repo, prID, *pr.GetUser().Login)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) removeCreateUser(ctx context.Context, config *ReviewConfig, user string) *ReviewConfig {
	var mustReviewer, Reviewer []string
	for _, u := range config.MustReviewers {
		if u != user {
			mustReviewer = append(mustReviewer, u)
		}
	}
	for _, u := range config.Reviewers {
		if u != user {
			Reviewer = append(Reviewer, u)
		}
	}
	return &ReviewConfig{
		MustReviewers:     mustReviewer,
		Reviewers:         Reviewer,
		NumberOfReviewers: config.NumberOfReviewers,
		AddAssignees:      config.AddAssignees,
	}
}

func (c *Client) selectReviewer(ctx context.Context, config *ReviewConfig) []string {
	var reviewers []string

	// All Reviewers
	if config.NumberOfReviewers == 0 {
		reviewers = append(reviewers, config.MustReviewers...)
		reviewers = append(reviewers, config.Reviewers...)
		return c.sliceUnique(reviewers)
	}

	if (config.NumberOfReviewers - len(config.MustReviewers)) <= 0 {
		if config.NumberOfReviewers > 1 {
			reviewers = append(reviewers, c.shuffle(config.MustReviewers)[:config.NumberOfReviewers]...)
		} else {
			reviewers = append(reviewers, c.shuffle(config.MustReviewers)[0])
		}
	}

	if (config.NumberOfReviewers - len(config.MustReviewers)) > 0 {
		if len(config.MustReviewers) > 0 {
			reviewers = append(reviewers, config.MustReviewers...)
		}

		if (config.NumberOfReviewers - len(config.MustReviewers)) > 1 {
			reviewers = append(reviewers, c.shuffle(c.sliceUnique(config.Reviewers))[:config.NumberOfReviewers-len(config.MustReviewers)]...)
		} else {
			reviewers = append(reviewers, c.shuffle(c.sliceUnique(config.Reviewers))[0])
		}
	}

	return reviewers
}

func (c *Client) addReviewers(ctx context.Context, owner, repo string, prID int, reviewers []string) error {
	reviewer := github.ReviewersRequest{
		Reviewers: reviewers,
	}
	_, _, err := c.client.PullRequests.RequestReviewers(ctx, owner, repo, prID, reviewer)
	if err != nil {
		return err
	}
	log.Infof("Added reviewers to PR #%d: %s", prID, strings.Join(reviewers, ","))
	return nil
}

func (c *Client) addAssignees(ctx context.Context, owner, repo string, prID int, user string) error {
	_, resp, err := c.client.Issues.AddAssignees(ctx, owner, repo, prID, []string{user})
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			log.Warn("assignees user not found.")
			return nil
		}
		return err
	}
	log.Infof("Added assignees to PR #%d: %s", prID, user)
	return nil
}
