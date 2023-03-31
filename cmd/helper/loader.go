package helper

import (
	"bufio"
	"os"
)

type IAddLoader interface {
	ScanAndLoadAddressFile() error
}

type AddLoader struct {
	addressOutput chan string
}

func NewAddLoader(addressOutput chan string) IAddLoader {
	return &AddLoader{
		addressOutput: addressOutput,
	}
}

func (a AddLoader) ScanAndLoadAddressFile() error {
	filePath := GetEnv("ADDRESS_FILE_PATH", "data/addresses.jsonl")

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) > 55 {
			a.addressOutput <- scanner.Text()[12:54]
		}
	}
	return scanner.Err()
}
