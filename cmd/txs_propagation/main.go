package main

import (
	"fmt"
	"sync"

	"github.com/docker/docker/client"
	"github.com/mrlutik/autoflowhub/pkg/keygen/usecase"

	"github.com/mrlutik/autoflowhub/internal/docker"
)

func main() {
	blockTolisten := 800
	KeysPath := "keydir"
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	reader := usecase.NewKeysReader(KeysPath)
	list, err := reader.GetAllAddresses()
	if err != nil {
		panic(err)
	}
	var arr []*docker.User = make([]*docker.User, 2500)
	for i := range list {
		arr[i] = &docker.User{Key: list[i], Balance: 0}
	}
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	c := make(chan int)
	go docker.BlockListener(client, "validator", blockTolisten, waitGroup, c)
	<-c
	// блокуєм виконання за допомогою читання з канала, запис відбудеться лише тоді коли блок досягне певної висоти

	// var arr []*docker.User = make([]*docker.User, 2500)
	// for i := 0; i < len(arr); i++ {
	// 	arr[i] = &docker.User{Key: "pepega", Balance: 280000}
	// }
	txcount := docker.TransactionSpam(client, waitGroup, 7000000/4, arr)
	waitGroup.Wait()
	fmt.Println(*txcount, "total tx was completed")

}
