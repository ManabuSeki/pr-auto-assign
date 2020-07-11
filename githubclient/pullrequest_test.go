package githubclient

import (
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v32/github"
)

func TestClient_selectReviewer(t *testing.T) {
	type args struct {
		config *ReviewConfig
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"all reviewers",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A"},
					Reviewers:         []string{"B", "C", "D"},
					NumberOfReviewers: 0,
				},
			},
			[]string{"A", "B", "C", "D"},
		},
		{
			"must reviewers 1",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A", "B"},
					Reviewers:         []string{"C", "D"},
					NumberOfReviewers: 1,
				},
			},
			[]string{"A"},
		},
		{
			"must reviewers 2",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A", "B"},
					Reviewers:         []string{"C", "D"},
					NumberOfReviewers: 2,
				},
			},
			[]string{"A", "B"},
		},
		{
			"must reviewers 2 reviewer 1",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A", "B"},
					Reviewers:         []string{"C", "D"},
					NumberOfReviewers: 3,
				},
			},
			[]string{"A", "B", "C"},
		},
		{
			"must reviewers 2 reviewer 2",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A", "B"},
					Reviewers:         []string{"C", "D"},
					NumberOfReviewers: 4,
				},
			},
			[]string{"A", "B", "C", "D"},
		},
		{
			"must reviewers 1 reviewer 2",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A"},
					Reviewers:         []string{"C", "D"},
					NumberOfReviewers: 3,
				},
			},
			[]string{"A", "C", "D"},
		},
		{
			"must reviewers 0 reviewer 2",
			args{
				&ReviewConfig{
					MustReviewers:     nil,
					Reviewers:         []string{"A", "B", "C", "D", "E", "F", "G"},
					NumberOfReviewers: 2,
				},
			},
			[]string{"B", "C"},
		},
		{
			"must reviewers 1 reviewer 3",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A"},
					Reviewers:         []string{"B", "C", "D", "E", "F", "G"},
					NumberOfReviewers: 4,
				},
			},
			[]string{"A", "B", "D", "E"},
		},
	}
	c := github.NewClient(nil)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			c := &Client{
				client: c,
			}
			got := c.selectReviewer(ctx, tt.args.config)
			sort.Strings(got)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Client.selectReviewer() mismatch (want- got+):\n%s", diff)
			}
		})
	}
}

func TestClient_removeCreateUser(t *testing.T) {
	type args struct {
		config *ReviewConfig
		user   string
	}
	tests := []struct {
		name string
		args args
		want *ReviewConfig
	}{
		{
			"remove create user",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A"},
					Reviewers:         []string{"B", "C", "D"},
					NumberOfReviewers: 1,
				},
				"A",
			},
			&ReviewConfig{
				MustReviewers:     nil,
				Reviewers:         []string{"B", "C", "D"},
				NumberOfReviewers: 1,
			},
		},
		{
			"remove create user2",
			args{
				&ReviewConfig{
					MustReviewers:     []string{"A"},
					Reviewers:         []string{"B", "C", "D"},
					NumberOfReviewers: 0,
				},
				"B",
			},
			&ReviewConfig{
				MustReviewers:     []string{"A"},
				Reviewers:         []string{"C", "D"},
				NumberOfReviewers: 0,
			},
		},
	}
	c := github.NewClient(nil)
	ctx := context.Background()
	for _, tt := range tests {
		t.Helper()

		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				client: c,
			}
			got := c.removeCreateUser(ctx, tt.args.config, tt.args.user)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Client.removeCreateUser() mismatch (want- got+):\n%s", diff)
			}
		})
	}
}
