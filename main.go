package main

import (
	"context"
	"math/rand"
	"os"
	"pr_auto_assign/githubclient"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	version string
)

func main() {
	log.Infof("pr_auto_assign version %s", version)
	var errCount int
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Error("GITHUB_TOKEN is missing")
		errCount++
	}

	own := os.Getenv("CI_REPO_OWNER")
	if own == "" {
		log.Errorf("CI_REPO_OWNER is missing")
		errCount++
	}
	repo := os.Getenv("CI_REPO_NAME")
	if repo == "" {
		log.Errorf("CI_REPO_NAME is missing")
		errCount++
	}
	configPath := os.Getenv("CI_CONFIG_PATH")
	if configPath == "" {
		log.Errorf("CI_CONFIG_PATH is missing")
		errCount++
	}
	ref := os.Getenv("CI_REPO_REF")
	if ref == "" {
		log.Errorf("CI_REPO_REF is missing")
		errCount++
	}
	tmpID := os.Getenv("CI_PULL_REQUEST_ID")
	if tmpID == "" {
		log.Errorf("CI_PULL_REQUEST_ID is missing")
		errCount++
	}
	prID, err := strconv.ParseInt(tmpID, 10, 64)
	if err != nil {
		panic(err)
	}

	if errCount != 0 {
		log.Errorf("failed to initialize action, error count: %d", errCount)
	}

	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	client := githubclient.NewClient(ctx, token)
	config, err := client.FetchConfig(ctx, own, repo, configPath, ref)
	if err != nil {
		panic(err)
	}
	if err := client.HandlePullRequest(ctx, own, repo, int(prID), config); err != nil {
		panic(err)
	}
}
