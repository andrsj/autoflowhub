package main

import (
	"sync"

	"github.com/docker/docker/client"

	"github.com/mrlutik/autoflowhub/internal/docker"
)

func main() {

	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	type users struct {
		key string
	}
	waitGroup := &sync.WaitGroup{}
	var arr []users = []users{{key: "kira1ap96a0dpx0mdjv9mxr7pnryx67hrry5melcs8n"}, {key: "kira1rh9wufxa9sxre53zapsvajsvyqvfvrhmckndta"}}
	for _, b := range arr {
		waitGroup.Add(1)
		go docker.RunTransaction(client, "validator", b.key, b.key, "1", "ukex", waitGroup, 5, 10)

	}
	waitGroup.Wait()
}
