package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/docker/docker/client"
	"github.com/mrlutik/autoflowhub/pkg/keygen/usecase"

	"github.com/mrlutik/autoflowhub/internal/docker"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("ERROR: \nUSSAGE: main <arg1> <arg2> <arg3>\narg1=total tx amount, if arg1=0 then default value=7000000/4 \narg2=folder with keys\n")
		os.Exit(1)
	}
	// disruptSum := (7000000 / 4) * 100
	disruptSum, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	if disruptSum == 0 {
		disruptSum = 7000000
	}
	disruptSum = disruptSum * 100
	KeysPath := os.Args[2]

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
