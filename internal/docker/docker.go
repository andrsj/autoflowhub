package docker

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"

	// "strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type User struct {
	Key     string
	Balance int
	Sent    bool
}

func BlockListener(dockerClient *client.Client, containerName string, blockToListen int, wg *sync.WaitGroup, c chan int) {
	var currentBlockHeight string
	converted := strconv.Itoa(blockToListen)
	for currentBlockHeight != converted {
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
		command := "sekaid tx bank send " + source + " " + destination + " " + amount + "ukex --keyring-backend=test --chain-id=$NETWORK_NAME --fees 100ukex --output=json --yes --home=$SEKAID_HOME " //можна додати --broadcast-mode=async але воно просто ігнорує помилку але не віксить її
		fmt.Println(command)
		out, err := ExecCommandInContainer(containerName, []string{`bash`, `-c`, command}, dockerClient)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
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
func TransactionSpam(dockerClient *client.Client, wg *sync.WaitGroup, txAmount int, users []*User) *int {
	fmt.Println("tx spamming")

	amountOfIterationForOneAcc := txAmount / len(users)
	if amountOfIterationForOneAcc < 1 {
		amountOfIterationForOneAcc = 1
	}
	fmt.Println("TRANSAKTION AMOUNNT", txAmount, len(users), amountOfIterationForOneAcc)
	fmt.Println("TOTAL TRANSAKTION AMOUNNT", len(users)*amountOfIterationForOneAcc)
	txCount := 0
	for i := 0; i < amountOfIterationForOneAcc; i++ {
		for u := range users {
			RunTransaction(dockerClient, "validator", users[u].Key, users[u].Key, "1", "ukex", 1, 0)
			//need to calibrate time in miliseconds
			//for 1750000 transactions *100miliseconds = 48h
			time.Sleep(time.Millisecond * 100)
			fmt.Println("+1tx", u)
			txCount++
		}
		fmt.Println(txCount)
	}
	return &txCount

}
func DisruptTokensBetweenAllAccounts(dockerClient *client.Client, wg *sync.WaitGroup, amountToOneAcc int, users []*User) {
	totalUsers := len(users)
	totalAmountOftokens := totalUsers * amountToOneAcc // + totalUsers*100
	var firstIterationWallets []int
	// fmt.Println(totalUsers)
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
		// fmt.Println(i)
		firstIterationWallets = append(firstIterationWallets, i)
	}
	// fmt.Println("FIRST WALLET", firstIterationWallets)
	for n, w := range firstIterationWallets {
		fmt.Println("FIRST WALLET INERATION NUMBER ", n)
		fmt.Println("TIME STARTED ", timestarted)
		fmt.Println("TIME SINCE   ", time.Since(timestarted))
		RunTransaction(dockerClient, "validator", "validator", users[w].Key, strconv.Itoa(firstIterationWalletsSum), "ukex", 1, 0)
		users[w].Balance += firstIterationWalletsSum
		//sending from 1st generation wallet to 2nd
		wg.Add(1)
		go func(wallet int) {
			total := wallet + divider
			if wallet+1 >= len(users) && total > len(users) {
				fmt.Println(totalAmountOftokens, "total tokens ")
				wg.Done()
				return
			}
			for i := wallet + 1; i < total; i++ {
				if wallet+1 > len(users) && i > len(users) {
					fmt.Println(totalAmountOftokens, "total tokens ")
					return
				}
				// fmt.Println(i, wallet, total, "TOOOOTAL")
				// fmt.Println(users[i], "heeeereeeee", total, firstIterationWallets, len(users))
				// fmt.Println(users[wallet], wallet, "sent from")
				// fmt.Println(users[i], i, "sent to")
				RunTransaction(dockerClient, "validator", users[wallet].Key, users[i].Key, strconv.Itoa(amountToOneAcc), "ukex", 1, 10010)
				users[i].Balance += amountToOneAcc
				users[wallet].Balance -= amountToOneAcc
			}
			defer wg.Done()
		}(w)
		fmt.Println(totalAmountOftokens, "total tokens ")
	}

}
