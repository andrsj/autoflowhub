package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mrlutik/autoflowhub/pkg/keygen/formatter"
)

type Executer interface {
	ExecuteCommand(context.Context, string, ...string) ([]byte, error)
}

type KeysClient struct {
	client         Executer
	homePath       string
	keyringBackend string
}

func NewKeysClient(exe Executer) *KeysClient {
	return &KeysClient{
		client:         exe,
		homePath:       "--home=/root/.sekaid-testnetwork-1",
		keyringBackend: "--keyring-backend=test",
	}
}

func (k *KeysClient) ListOfKeys() error {
	out, err := k.client.ExecuteCommand(context.Background(), "sekai",
		"sekaid",
		"keys",
		"list",
		k.homePath,
		k.keyringBackend,
		"--output=json",
	)
	if err != nil {
		return err
	}

	_, err = formatter.ParseListCMDOutput(out)
	if err != nil {
		return err
	}

	return err
}

func (k *KeysClient) AddKey(name string) error {
	out, err := k.client.ExecuteCommand(context.Background(), "sekai",
		"sekaid",
		"keys",
		"add",
		name,
		k.homePath,
		k.keyringBackend,
		"--output=json",
	)
	if err != nil {
		return err
	}

	key, err := formatter.ParseAddCMDOutput(out)
	if err != nil {
		return err
	}

	createKeyFiles(key)

	return err
}

func createKeyFiles(k *formatter.KeyWithMnemonic) error {
	dir := "./data"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println("Creating a folder")
		os.Mkdir(dir, 0o755) // or 'os.MkdirAll(dir, 0755)' to create parent directories as needed
	}

	addressFileName := fmt.Sprintf("%s/%s.address", dir, k.Name)
	addressFile, err := os.Create(addressFileName)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %w", addressFileName, err)
	}
	defer addressFile.Close()

	mnemonicFileName := fmt.Sprintf("%s/%s.mnemonic", dir, k.Name)
	mnemonicFile, err := os.Create(mnemonicFileName)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %w", mnemonicFileName, err)
	}
	defer mnemonicFile.Close()

	_, err = io.WriteString(addressFile, k.Address)
	if err != nil {
		return fmt.Errorf("error writing to file '%s': %w", addressFileName, err)
	}

	_, err = io.WriteString(mnemonicFile, k.Mnemonic)
	if err != nil {
		return fmt.Errorf("error writing to file '%s': %w", mnemonicFileName, err)
	}

	log.Printf("Files ['%s', '%s'] written successfully\n", addressFileName, mnemonicFileName)

	return nil
}
