package main

import (
	"context"
	"m/githubclient"
	"math/rand"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var errCount int
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Error("GITHUB_TOKEN is missing")
		errCount++
	}

	tmpID := os.Getenv("CI_PULL_REQUEST_ID")
	if tmpID == "" {
		log.Error("CI_PULL_REQUEST_ID is missing")
		errCount++
	}
	prID, err := strconv.ParseInt(tmpID, 10, 64)
	if err != nil {
		panic(err)
	}
	own := os.Getenv("CI_REPO_OWNER")
	if own == "" {
		log.Error("CI_REPO_OWNER is missing")
		errCount++
	}
	repo := os.Getenv("CI_REPO_NAME")
	if repo == "" {
		log.Error("CI_REPO_NAME is missing")
		errCount++
	}
	configPath := os.Getenv("CI_CONFIG_PATH")
	if configPath == "" {
		log.Error("CI_CONFIG_PATH is missing")
		errCount++
	}

	if errCount != 0 {
		log.Errorf("failed to initialize action, error count: %d", errCount)
	}

	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	client := githubclient.NewClient(ctx, token)
	config, err := client.FetchConfig(ctx, own, repo, configPath, token)
	if err != nil {
		panic(err)
	}
	if err := client.HandlePullRequest(ctx, own, repo, int(prID), config); err != nil {
		panic(err)
	}
}
