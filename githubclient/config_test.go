package githubclient

import (
	"context"
	"net/http"
	"testing"
)

func TestFetchConfig(t *testing.T) {
	t.Run("get content", func(t *testing.T) {
		t.Helper()

		client, mux, _, teardown := setup()
		defer teardown()
		mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./testdata/fetch_config.json")
		})

		got, err := client.FetchConfig(context.Background(), "o", "r", "p")
		if err != nil {
			t.Fatal(err)
		}

		want := "Hello, world!"
		if want != got {
			t.Errorf("unexpected config want %s got %s", want, got)
		}
	})
}
