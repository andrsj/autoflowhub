package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/mrlutik/autoflowhub/pkg/keygen/formatter"
)

type Executer interface {
	ExecuteCommand(context.Context, string, ...string) ([]byte, error)
}

type KeysClient struct {
	client         Executer
	containerName  string
	homePath       string
	keyringBackend string
	dirOfKeys      string
}

func NewKeysClient(exe Executer, containerName, homePath, keyringBackend, dirKeys string) *KeysClient {
	return &KeysClient{
		client:         exe,
		containerName:  containerName,
		homePath:       fmt.Sprintf("--home=%s", homePath),
		keyringBackend: fmt.Sprintf("--keyring-backend=%s", keyringBackend),
		dirOfKeys:      dirKeys,
	}
}

func (k *KeysClient) ListOfKeys() ([]string, error) {
	out, err := k.client.ExecuteCommand(context.Background(), k.containerName,
		"sekaid",
		"keys",
		"list",
		k.homePath,
		k.keyringBackend,
		"--output=json",
	)
	if err != nil {
		return nil, err
	}

	keys, err := formatter.ParseListCMDOutput(out)
	if err != nil {
		return nil, err
	}

	addresses := make([]string, 0, len(keys))
	for _, key := range keys {
		addresses = append(addresses, key.Address)
	}

	return addresses, err
}

func (k *KeysClient) GenerateKeys(count int) ([]string, error) {
	
	if count <= 0 {
		return nil, fmt.Errorf("the count '%d' needs to be positive", count)
	}

	log.Printf("Adding '%d' users\n", count)

	addresses := make([]string, 0, count)
	for i := 0; i < count; i++ {
		newUUID := uuid.New()
		log.Printf("Creating account #%d: '%s'\n", i, newUUID)

		address, err := k.addKey(newUUID.String())
		if err != nil {
			return nil, fmt.Errorf("can't add new key (exit loop): %w", err)
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (k *KeysClient) addKey(name string) (string, error) {
	out, err := k.client.ExecuteCommand(context.Background(), k.containerName,
		"sekaid",
		"keys",
		"add",
		name,
		k.homePath,
		k.keyringBackend,
		"--output=json",
	)
	if err != nil {
		return "", err
	}

	key, err := formatter.ParseAddCMDOutput(out)
	if err != nil {
		return "", err
	}

	k.createKeyFiles(key)

	return "", err
}

func (k KeysClient) createKeyFiles(key *formatter.KeyWithMnemonic) error {
	dir, err := filepath.Abs(k.dirOfKeys)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		log.Println("Creating a folder:", dir)
		os.Mkdir(dir, 0o755) // or 'os.MkdirAll(dir, 0755)' to create parent directories as needed
	}

	addressFileName := fmt.Sprintf("%s/%s.address", dir, key.Name)
	addressFile, err := os.Create(addressFileName)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %w", addressFileName, err)
	}
	defer addressFile.Close()

	mnemonicFileName := fmt.Sprintf("%s/%s.mnemonic", dir, key.Name)
	mnemonicFile, err := os.Create(mnemonicFileName)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %w", mnemonicFileName, err)
	}
	defer mnemonicFile.Close()

	_, err = io.WriteString(addressFile, key.Address)
	if err != nil {
		return fmt.Errorf("error writing to file '%s': %w", addressFileName, err)
	}

	_, err = io.WriteString(mnemonicFile, key.Mnemonic)
	if err != nil {
		return fmt.Errorf("error writing to file '%s': %w", mnemonicFileName, err)
	}

	log.Printf("Files ['%s', '%s'] written successfully\n", addressFileName, mnemonicFileName)

	return nil
}
