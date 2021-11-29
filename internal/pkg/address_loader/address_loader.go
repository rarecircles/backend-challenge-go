package address_loader

import (
	"bufio"
	"go.uber.org/zap"
	"os"
)

type AddressLoader interface {
	LoadAddressFile() error
}

type AddressLoaderImpl struct {
	outputChan chan string
	zLog *zap.Logger
}

func NewAddressLoader(output chan string, logger *zap.Logger) AddressLoader {
	return &AddressLoaderImpl{
		outputChan: output,
		zLog: logger,
	}
}

func (ad AddressLoaderImpl) LoadAddressFile() error {
	filePath := "data/addresses.jsonl"

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		if err = f.Close(); err != nil {
			ad.zLog.Fatal(err.Error())
		}
	}()

	s := bufio.NewScanner(f)
	//i := 0 // <- uncomment this to see lines read.
	for s.Scan() {
		if len(s.Text()) > 55 {
			//ad.zLog.Info(fmt.Sprintf("address file: reading line - %d", i)) // <- uncomment this to see lines read.
			ad.outputChan <- s.Text()[12:54]
			//i++ // <- uncomment this to see lines read.
		}
	}
	return s.Err()
}
