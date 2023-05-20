package docker

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type User struct {
	Key string
}

func BlockListener(dockerClient *client.Client, containerName, blockToListen string, wg *sync.WaitGroup, c chan int) {
	var currentBlockHeight string
	for currentBlockHeight != blockToListen {
		time.Sleep(time.Second * 1)
		currentBlockHeight = GetBlockHeight(dockerClient, containerName)
		fmt.Println("block to listen ", blockToListen)

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
	// fmt.Println(string(out))
	return value
}
func RunTransaction(dockerClient *client.Client, containerName, source, destination, amount, demon string, txAmount, sleepTimeBetweenTxInMeeleseconds int) {
	//    sekaid tx bank send $SOURCE $DESTINATION "${AMOUNT}${DENOM}" --keyring-backend=test --chain-id=$NETWORK_NAME --fees "${FEE_AMOUNT}${FEE_DENOM}" --output=json --yes --home=$SEKAID_HOME | txAwait 180
	for i := 0; i < txAmount; i++ {
		command := "sekaid tx bank send " + source + " " + destination + " " + amount + "ukex --keyring-backend=test --chain-id=$NETWORK_NAME --fees 100ukex --output=json --yes --home=$SEKAID_HOME "
		fmt.Println(command)
		_, err := ExecCommandInContainer(containerName, []string{`bash`, `-c`, command}, dockerClient)
		if err != nil {
			panic(err)
		}
		// fmt.Println(string(out))
		if sleepTimeBetweenTxInMeeleseconds > 0 {
			time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenTxInMeeleseconds))
		}
	}
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

// func DisruptTokensBetweenAllAccounts(dockerClient *client.Client, users []User, summToAddIntoOneAcc int) {
// 	totalAmount := len(users) * summToAddIntoOneAcc

//		// stringAmountToOneAcc := strconv.Itoa(summToAddIntoOneAcc)
//		// for i, u := range users {
//		// 	fmt.Println(i, u, totalAmount)
//		// 	wg.Add(1)
//		// 	RunTransaction(dockerClient, "validator", "validator", u.Key, stringAmountToOneAcc, "ukex", wg, 1, 10)
//		// }
//		// v > 0 > 1,2 >
//		//1 >>3,4
//		// 2 >>5,6
//		// calculateSendings(dockerClient, summToAddIntoOneAcc, totalAmount, 8000, users)
//	}
func DisruptTokensBetweenAllAccounts(dockerClient *client.Client, wg *sync.WaitGroup, amountToOneAcc int, users []User) {
	totalUsers := len(users)
	totalAmountOftokens := totalUsers * amountToOneAcc
	var firstIterationWallets []int
	fmt.Println(totalUsers)
	divider := totalUsers / 100
	var firstIterationWalletsSum int
	if divider > 0 {
		firstIterationWalletsSum = totalAmountOftokens / divider
	} else {
		divider = totalUsers
		firstIterationWalletsSum = totalAmountOftokens
	}
	timestarted := time.Now()
	fmt.Println("DIVIDERS", divider)
	for i := 0; i < totalUsers; i = i + divider {
		fmt.Println(i)
		firstIterationWallets = append(firstIterationWallets, i)
	}
	fmt.Println("FIRST WALLET", firstIterationWallets)
	for n, w := range firstIterationWallets {

		RunTransaction(dockerClient, "validator", "validator", users[w].Key, strconv.Itoa(firstIterationWalletsSum), "ukex", 1, 10050)
		//sending from 1st generation wallet to second others
		wg.Add(1)
		go func(wallet int) {
			total := wallet + divider

			for i := wallet + 1; i < total-1; i++ {
				RunTransaction(dockerClient, "validator", users[wallet].Key, users[i].Key, strconv.Itoa(amountToOneAcc), "ukex", 1, 10050)
				fmt.Println(wallet, total, "TOOOOTAL")

			}
			wg.Done()
		}(w)
		fmt.Println("FIRST WALLET INERATION NUMBER ", n)
		fmt.Println("TIME STARTED ", timestarted)
		fmt.Println("TIME SINCE   ", time.Since(timestarted))
	}
}

// func sendToSecondGenerationWallets(dockerClient *client.Client, wg *sync.WaitGroup, amountToOneAcc int, users []User) {

// }

// func calculateOwnStake(totalAmount, amountToOneAcc int) int {
// 	ret := totalAmount - amountToOneAcc
// 	if ret <= amountToOneAcc {
// 		return 0
// 	}
// 	return ret
// }

// func findSmallestDivisor(num int) int {
// 	for i := 2; i <= num; i++ {
// 		if num%i == 0 {
// 			return i
// 		}
// 	}
// 	return num
// }
