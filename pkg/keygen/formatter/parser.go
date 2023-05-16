package formatter

import (
	"encoding/json"
	"fmt"
	"log"
)

type Key struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func ParseListCMDOutput(output []byte) ([]Key, error) {
	var (
		err  error
		keys []Key
	)

	err = json.Unmarshal(output, &keys)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal JSON output: %w", err)
	}

	// Output to console
	formattedOutput, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("can't marshal the JSON: %w", err)
	}
	log.Println(string(formattedOutput))

	return keys, nil
}

type KeyWithMnemonic struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Mnemonic string `json:"mnemonic"`
}

func ParseAddCMDOutput(output []byte) (*KeyWithMnemonic, error) {
	var (
		err error
		key KeyWithMnemonic
	)

	err = json.Unmarshal(output, &key)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal JSON output: %w", err)
	}

	// Output to console
	formattedOutput, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("can't marshal the JSON: %w", err)
	}
	log.Println(string(formattedOutput))

	return &key, nil
}
