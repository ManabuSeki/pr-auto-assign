package githubclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFetchConfig(t *testing.T) {
	t.Run("get content", func(t *testing.T) {
		t.Helper()

		client, mux, _, teardown := setup()
		defer teardown()
		mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./testdata/fetch_config.json")
		})

		got, err := client.FetchConfig(context.Background(), "o", "r", "p", "dummy")
		if err != nil {
			t.Fatal(err)
		}

		want := &ReviewConfig{
			MustReviewers:     []string{"reviewerA"},
			Reviewers:         []string{"reviewerA", "reviewerB", "reviewerC"},
			NumberOfReviewers: 2,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("review config mismatch(-want +got):\n%s", diff)
		}
	})
}
