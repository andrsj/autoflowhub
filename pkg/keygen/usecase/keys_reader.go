package usecase

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type KeysReader struct {
	pathToKeys string
}

func NewKeysReader(pathToKeys string) *KeysReader {
	return &KeysReader{
		pathToKeys: pathToKeys,
	}
}

func (k *KeysReader) GetAllAddresses() ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.address", k.pathToKeys))
	if err != nil {
		return nil, fmt.Errorf("can't find key's addresses: %w", err)
	}

	var keysAddresses []string
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("can't read file: %w", err)
		}

		keysAddresses = append(keysAddresses, string(content))
	}

	return keysAddresses, nil
}
