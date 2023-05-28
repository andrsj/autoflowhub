package main

import (
	"fmt"
	"sync"

	"github.com/docker/docker/client"
	"github.com/mrlutik/autoflowhub/pkg/keygen/usecase"

	"github.com/mrlutik/autoflowhub/internal/docker"
)

func main() {
	disruptSum := (7000000 / 4) * 100
	KeysPath := "./keyTest2"
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	reader := usecase.NewKeysReader(KeysPath)
	list, err := reader.GetAllAddresses()
	if err != nil {
		panic(err)
	}
	var arr []*docker.User = make([]*docker.User, len(list))
	for i := range list {
		arr[i] = &docker.User{Key: list[i], Balance: 0}
	}
	// var arr []*docker.User = make([]*docker.User, 2500)
	// for i := 0; i < 2500; i++ {
	// 	arr[i] = &docker.User{Key: "pepelaugh", Balance: 0}
	// }
	waitGroup := &sync.WaitGroup{}
	fmt.Println(arr[0])
	docker.DisruptTokensBetweenAllAccounts(client, waitGroup, disruptSum, arr)
	waitGroup.Wait()
	for u := range arr {
		fmt.Println(arr[u])
	}

}
