package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mrlutik/autoflowhub/internal/adapters"
)

func main() {
	var accessToken = os.Getenv("GITHUB_TOKEN")

	adapter := adapters.NewGitHubAdapter(accessToken)

	repositories := []struct {
		Owner string
		Repo  string
	}{
		{"KiraCore", "sekai"},
		{"KiraCore", "interx"},
		{"KiraCore", "kira"},
		{"KiraCore", "miro"},
		{"KiraCore", "tools"},
	}

	var wg sync.WaitGroup
	results := make(chan string)

	for _, repo := range repositories {
		wg.Add(1)
		go func(owner, repo string) {
			defer wg.Done()

			latestRelease, err := adapter.GetLatestRelease(owner, repo)
			if err != nil {
				log.Printf("Error fetching latest release for %s/%s: %v\n", owner, repo, err)
				return
			}

			results <- fmt.Sprintf("Latest release for %s/%s:\t%s\n", owner, repo, *latestRelease.TagName)
		}(repo.Owner, repo.Repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Print(result)
	}

}
