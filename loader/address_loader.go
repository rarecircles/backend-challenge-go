package loader

import (
	"bufio"
	"go.uber.org/zap"
	"os"
)

type AddressLoader interface {
	LoadAddressFile(filePath string) error
}

type AddressLoaderImpl struct {
	outputChan chan string
	zLog       *zap.Logger
}

func NewAddressLoader(output chan string, logger *zap.Logger) AddressLoader {
	return &AddressLoaderImpl{
		outputChan: output,
		zLog:       logger,
	}
}

func (a AddressLoaderImpl) LoadAddressFile(filePath string) error {

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		if err = f.Close(); err != nil {
			a.zLog.Fatal(err.Error())
		}
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		if len(s.Text()) > 55 {
			a.outputChan <- s.Text()[12:54]
		}
	}
	return s.Err()
}
