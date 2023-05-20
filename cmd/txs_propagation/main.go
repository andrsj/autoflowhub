package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/docker/docker/client"

	"github.com/mrlutik/autoflowhub/internal/docker"
)

func main() {

	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	waitGroup := &sync.WaitGroup{}

	// waitGroup.Add(1)
	// c := make(chan int)
	// go docker.BlockListener(client, "validator", "925", waitGroup, c)
	var arr []docker.User = []docker.User{
		{Key: "kira1u890hktk35k22y256lenaweuc86f23tnsawjl8"},
		{Key: "kira1gf573vs5du5yfj3ck6mrkfplnqzqghek0rqlqa"},
		{Key: "kira1ymcmfmwqnh2vz86w5jf24h02sv7ey3reuaw7w2"},
		{Key: "kira1dwqu7essep0efesv5x738gvg6x509ledk5c6gf"},
		{Key: "kira1hclghrqumgv08zzq2w268agfpsrrjxnec0678y"},
		{Key: "kira1qrtx2cvj4sel43jp5ymzxsg65ktf7ydeypf83g"},
		{Key: "kira1d5acaqhrjf7avvtxt8ljecwtnyadnd2drkra8r"},
		{Key: "kira1ny6qpxczmrggzu343u0mwps54p9zfvytvxy4mz"},
	}
	fmt.Println(arr)
	var usirs [1000]docker.User
	for a, b := range usirs {
		b.Key = string("pepega" + strconv.Itoa(a))
	}
	docker.DisruptTokensBetweenAllAccounts(client, waitGroup, 10000, usirs[:])
	//блокуєм виконання за допомогою читання з канала, запис відбудеться лише тоді коли блок досягне певної висоти
	// <-c

	// for _, b := range arr {
	// 	fmt.Println("add goutine")
	// 	waitGroup.Add(1)
	// 	go func(b *docker.User, wg *sync.WaitGroup) {
	// 		fmt.Println("done")
	// 		docker.RunTransaction(client, "validator", b.Key, b.Key, "1", "ukex", 1, 0)
	// 		wg.Done()
	// 	}(&b, waitGroup)
	// 	time.Sleep(time.Millisecond * 100)

	// }

	waitGroup.Wait()

}
