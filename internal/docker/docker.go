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

type Userp struct {
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
			// if wallet+1 > total {
			// 	return
			// }
			for i := wallet + 1; i < total; i++ {

				RunTransaction(dockerClient, "validator", users[wallet].Key, users[i].Key, strconv.Itoa(amountToOneAcc), "ukex", 1, 10000)
				fmt.Println(i, wallet, total, "TOOOOTAL")

			}
			wg.Done()
		}(w)
		fmt.Println("FIRST WALLET INERATION NUMBER ", n)
		fmt.Println("TIME STARTED ", timestarted)
		fmt.Println("TIME SINCE   ", time.Since(timestarted))
	}
}

type User struct {
	Key     string
	Balance int
	Sent    bool
}

func AlgoritmTesting(users []*User, amountToOneAcc int) {
	// totalUsers := len(users)
	// // totalAmountOftokens := totalUsers * amountToOneAcc
	// var firstIterationWallets [4]int
	// firstgenerationDivider := totalUsers / 4
	// firstIterationWallets[0] = 0
	// firstIterationWallets[1] = firstgenerationDivider
	// firstIterationWallets[2] = firstgenerationDivider + firstgenerationDivider
	// firstIterationWallets[3] = firstgenerationDivider + firstgenerationDivider + firstgenerationDivider
	// tokensSpreding(users, amountToOneAcc, totalAmountOftokens, 4)
	// VALIDATOR.Balance = 99999999
	// infiniteUser := &User{Balance: 1<<63 - 1}
	// transferTokens(infiniteUser, users, 200, 2)
	// fmt.Println(infiniteUser.Balance, COUNT)
	// fmt.Println(users)
	// for _, user := range users {
	// 	fmt.Printf("Пользователь %v: баланс - %d\n", user.Key, user.Balance)
	// }

}
func TestFunc(users []*User) {
	// users := make([]*User, 10000)
	for i := range users {
		users[i] = &User{Balance: 0}
	}
	for _, b := range users {
		fmt.Println(b)
	}

	// бесконечный пользователь
	infiniteUser := &User{Balance: 1<<63 - 1}

	queue := []*User{infiniteUser}
	queue = append(queue, users...)

	i := 0
	for len(queue) > 0 {
		user := queue[0]
		queue = queue[1:]

		if user.Sent {
			continue
		}

		recipients := []*User{}
		if i < len(users) {
			recipients = append(recipients, users[i])
			queue = append(queue, users[i])
			i++
		}

		transferTokens(user, recipients, 200)
	}

	for i, user := range users {
		if i > 10 { // Выводим только первые 10 пользователей для упрощения
			break
		}
		fmt.Printf("User[%d] balance: %d\n", i, user.Balance)
	}
	fmt.Printf("InfiniteUser balance: %d\n", infiniteUser.Balance)
}

var COUNT int

func transferTokens(from *User, to []*User, amount int) {
	if from.Sent {
		return
	}

	from.Balance -= amount * len(to)
	from.Sent = true
	for i := range to {
		to[i].Balance += amount
	}
}

// func distributeCoins(users []*User, senderIdx, coins int) {
// 	if coins == 0 || len(users) == 0 {
// 		return
// 	}

// 	if senderIdx >= len(users) || users[senderIdx].Balance < coins {
// 		return
// 	}

// 	users[senderIdx].Balance -= coins

// 	receiver1 := (senderIdx + 1) % len(users)
// 	receiver2 := (senderIdx + 2) % len(users)

// 	users[receiver1].Balance += coins / 2
// 	users[receiver2].Balance += coins - (coins / 2)

// 	remainingCoins := coins / 2
// 	distributeCoins(users, receiver1, remainingCoins)
// 	distributeCoins(users, receiver2, coins-remainingCoins)
// }

// var VALIDATOR TestUsers

// func tokensSpreding(users []*TestUsers, amountToOneAcc, totalAmount, divider int) {
// 	var curentLayerOfAccs []*TestUsers
// 	position := 0
// 	for i := 0; i < divider; i++ {
// 		fmt.Println(len(users), divider)
// 		position = (len(users) / divider) * i
// 		fmt.Println(position, "position")
// 		curentLayerOfAccs = append(curentLayerOfAccs, users[position])
// 		sendTokens(&VALIDATOR, curentLayerOfAccs[i], totalAmount/divider)
// 	}

// 	fmt.Println(curentLayerOfAccs, "curentLayerOfAccs")
// }
// func sendTokens(userFrom, userTo *TestUsers, amountToSend int) {
// 	userFrom.Balance -= amountToSend
// 	userTo.Balance += amountToSend
// }

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
