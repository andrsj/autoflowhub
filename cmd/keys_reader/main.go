package main

import (
	"flag"
	"log"

	"github.com/mrlutik/autoflowhub/pkg/keygen/usecase"
)

func main() {
	pathPtr := flag.String("path", "", "Path to dir")
	flag.Parse()

	if *pathPtr == "" {
		log.Println("Please provide a --path")
		flag.PrintDefaults()
		return
	}

	keysReader := usecase.NewKeysReader(*pathPtr)
	addresses, err := keysReader.GetAllAddresses()
	if err != nil {
		log.Fatal(err)
	}

	for i, address := range addresses {
		log.Printf("%d: %s", i, address)
	}
}
