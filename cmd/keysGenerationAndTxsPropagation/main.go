package main

import (
	"log"
	"strconv"

	"fmt"
	"sync"

	"github.com/mrlutik/autoflowhub/pkg/keygen/docker"
	"github.com/mrlutik/autoflowhub/pkg/keygen/usecase"
	"github.com/spf13/cobra"

	"github.com/docker/docker/client"

	idocker "github.com/mrlutik/autoflowhub/internal/docker"
)

const (
	use              = "keysgen"
	shortDescription = "CLI application for generate accounts for Kira Network"
)

const longDescription = `This application accepts three parameters: 
home of sekaid, 
keyring-backend value
count of keys which will be generated
directory for saved address and mnemonic for keys
There is no default values!`

var KeyringBackend string
var Home string
var SekaiContainer string
var KeysPath string
var BlockToListen int
var TxAmount int
var Count int

var rootCmd = &cobra.Command{
	Use:   use,
	Short: shortDescription,
	Long:  longDescription,
	Run: func(cmd *cobra.Command, _ []string) {
		Home, _ = cmd.Flags().GetString("home")
		KeyringBackend, _ = cmd.Flags().GetString("keyring-backend")
		KeysPath, _ = cmd.Flags().GetString("keys-dir")
		SekaiContainer, _ = cmd.Flags().GetString("sekai")
		Count, _ = cmd.Flags().GetInt("count")
		BlockToListen, _ = cmd.Flags().GetInt("blockToListen")
		TxAmount, _ = cmd.Flags().GetInt("txAmount")

		if Home == "" || SekaiContainer == "" || KeyringBackend == "" || KeysPath == "" || Count < 0 {
			cmd.Help()
			log.Fatal("Please provide all required parameters: home, backend and positive count")
		}

		log.Println("Sekai Container:", SekaiContainer)
		log.Println("Home:", Home)
		log.Println("Backend:", KeyringBackend)
		log.Println("Directory of keys:", KeysPath)
		log.Println("Count:", Count)
		log.Println("Block to listen:", BlockToListen)
		log.Println("Amount of transactions:", TxAmount)

	},
}

func main() {
	// Usage: keygen --home "/root/.sekaid-testnetwork-1" -k "test" -c 4 -d ./data -s sekai
	rootCmd.PersistentFlags().String("home", "", "Path to sekaid home")
	rootCmd.PersistentFlags().StringP("keys-dir", "d", "", "Keys directory (relative or absolute path)")
	rootCmd.PersistentFlags().StringP("keyring-backend", "k", "", "Keyring backend")
	rootCmd.PersistentFlags().StringP("sekai", "s", "", "Sekaid container name")
	rootCmd.PersistentFlags().IntP("count", "c", 0, "Count of keys which will be added")
	rootCmd.PersistentFlags().IntP("blockToListen", "b", 0, "which block to listen")
	rootCmd.PersistentFlags().IntP("txAmount", "t", 0, "how much transactions from generated users you want")
	fmt.Println("BLOCKBLOCK", BlockToListen)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	reader := usecase.NewKeysReader(KeysPath)
	list, err := reader.GetAllAddresses()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(list), "list before")
	n := len(list)
	if n >= Count {
		list = list[0:Count]
		Count = 0
		fmt.Println(len(list), "len of list ")
	}
	if Count > 0 {
		generating(SekaiContainer, Home, KeyringBackend, KeysPath, Count)
	}

	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)
	c := make(chan int)
	go idocker.BlockListener(client, "validator", strconv.Itoa(BlockToListen), waitGroup, c)
	var arr []*idocker.User = make([]*idocker.User, len(list))
	for i := range list {
		arr[i] = &idocker.User{Key: list[i], Balance: 0}
	}
	fmt.Println(len(arr), "LEEEEEEEEEEEEEEEEEN")
	fmt.Println(2)

	disruptSum := TxAmount * 100
	idocker.DisruptTokensBetweenAllAccounts(client, waitGroup, disruptSum, arr[:])
	// блокуєм виконання за допомогою читання з канала, запис відбудеться лише тоді коли блок досягне певної висоти
	fmt.Println(3)

	<-c
	fmt.Println(4)

	waitGroup.Wait()
	fmt.Println(5)

	for _, u := range arr {
		fmt.Println(u)
	}
	fmt.Println(1)
	txcount := idocker.TransactionSpam(client, waitGroup, TxAmount, arr)
	waitGroup.Wait()
	fmt.Println(txcount)
}

func generating(containerName, homePath, keyringBackend, dirOfKeys string, count int) {
	dockerClient := docker.NewDockerCommandRunner()
	keysUsecase := usecase.NewKeysClient(dockerClient, containerName, homePath, keyringBackend, dirOfKeys)

	var err error
	addresses, err := keysUsecase.GenerateKeys(count)
	if err != nil {
		log.Fatal(err)
	}

	allAddresses, err := keysUsecase.ListOfKeys()
	if err != nil {
		log.Fatal(err)
	}

	// NEXT BUSINESS LOGIC HERE
	// allAddresses includes all users in Kira network
	// addresses includes only generated keys

	log.Println("Checking generated addresses in the list of all...")
	if containsAllStrings(allAddresses, addresses) {
		log.Fatal("Error: not all generated addresses are saved")
	}

	log.Println("All is O'key!")
}

// Temporary helpers which are used for checking program
func sliceContains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func containsAllStrings(slice1 []string, slice2 []string) bool {
	for _, b := range slice2 {
		if !sliceContains(slice1, b) {
			return false
		}
	}
	return true
}
