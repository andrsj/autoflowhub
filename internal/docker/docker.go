package docker

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunTransaction(dockerClient *client.Client, containerName, source, destination, amount, demon string, wg *sync.WaitGroup, txAmount, sleepTimeBetweenTx int) {
	//    sekaid tx bank send $SOURCE $DESTINATION "${AMOUNT}${DENOM}" --keyring-backend=test --chain-id=$NETWORK_NAME --fees "${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME | txAwait 180
	for i := 0; i < txAmount; i++ {
		command := "sekaid tx bank send " + source + " " + destination + " " + amount + "ukex --keyring-backend=test --chain-id=$NETWORK_NAME --fees 100ukex --output=json --yes --home=$SEKAID_HOME "
		fmt.Println(command)
		out, err := ExecCommandInContainer(containerName, []string{`bash`, `-c`, command}, dockerClient)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		time.Sleep(time.Second * time.Duration(sleepTimeBetweenTx))
	}
	wg.Done()
	//sekaid tx bank send kira1rh9wufxa9sxre53zapsvajsvyqvfvrhmckndta kira1cv89u97thwm837uzlcxh8yc0hele0jxwsklk93 100ukex --keyring-backend=test --chain-id=$NETWORK_NAME --fees 100ukex --output=json --yes --home=$SEKAID_HOME | txAwait 180

}

func ExecCommandInContainer(containerID string, command []string, Cli *client.Client) ([]byte, error) {
	execCreateResponse, err := Cli.ContainerExecCreate(context.Background(), containerID, types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, err
	}
	execAttachConfig := types.ExecStartCheck{}
	resp, err := Cli.ContainerExecAttach(context.Background(), execCreateResponse.ID, execAttachConfig)
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	// Read the output
	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		panic(err)
	}
	return output, nil
}
