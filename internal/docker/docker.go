package docker

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func BlockListener(dockerClient *client.Client, containerName, blockToListen string, wg *sync.WaitGroup, c chan int) {
	var currentBlockHeight string
	for currentBlockHeight != blockToListen {
		time.Sleep(time.Second * 1)
		currentBlockHeight = GetBlockHeight(dockerClient, containerName)
		fmt.Println("current", currentBlockHeight)

	}
	c <- 1
	fmt.Println("curent:", currentBlockHeight, "goal:", blockToListen)
	defer wg.Done()

}
func GetBlockHeight(dockerClient *client.Client, containerName string) string {
	command := "sekaid status"
	out, err := ExecCommandInContainer(containerName, []string{`bash`, `-c`, command}, dockerClient)
	if err != nil {
		panic(err)
	}
	str := string(out)
	re := regexp.MustCompile(`"latest_block_height":"(\d+)"`) // Regular expression pattern

	match := re.FindStringSubmatch(str) // Find the match in the string
	var value string
	if len(match) == 2 {
		value = match[1] // Extract the value from the match
		fmt.Println("latest_block_height:", value)
	} else {
		fmt.Println("latest_block_height not found")
	}
	fmt.Println(string(out))
	return value
}
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
	defer wg.Done()
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
