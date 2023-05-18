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

	type user struct {
		key string
	}
	c := make(chan int)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	go docker.BlockListener(client, "validator", "1578", waitGroup, c)

	//блокуєм виконання за допомогою читання з канала, запис відбудеться лише тоді коли блок досягне певної висоти
	<-c

	var arr []user = []user{{key: "kira1ap96a0dpx0mdjv9mxr7pnryx67hrry5melcs8n"}, {key: "kira1rh9wufxa9sxre53zapsvajsvyqvfvrhmckndta"}}
	for _, b := range arr {
		waitGroup.Add(1)
		go docker.RunTransaction(client, "validator", b.key, b.key, "1", "ukex", waitGroup, 5, 1)

	}

	waitGroup.Wait()

}
