package main

import (
	"log"

	"github.com/mrlutik/autoflowhub/pkg/keygen/docker"
	"github.com/mrlutik/autoflowhub/pkg/keygen/usecase"
)

func main() {
	log.Println("Hello!")

	dockerClient := docker.NewDockerCommandRunner()
	keysUsecase := usecase.NewKeysClient(dockerClient)

	var err error
	err = keysUsecase.AddKey("test7")
	if err != nil {
		log.Fatal(err)
	}

	err = keysUsecase.ListOfKeys()
	if err != nil {
		log.Fatal(err)
	}
}
