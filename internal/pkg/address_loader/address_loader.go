// Package address_loader loads address files
package address_loader

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type AddressLoader interface {
	Load(filePath string) error
}

type addressLoader struct {
	log *zap.Logger
	ch  chan<- string
}

// NewAddressLoader creates AddressLoader
func NewAddressLoader(log *zap.Logger, ch chan<- string) AddressLoader {
	return &addressLoader{
		log: log,
		ch:  ch,
	}
}

func (al addressLoader) Load(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open a file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			al.log.Error("failed to close a file " + err.Error())
			return
		}
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		if len(s.Text()) == 56 {
			address := s.Text()[12:54]
			al.ch <- address
		}
	}
	close(al.ch)

	if err := s.Err(); err != nil {
		return fmt.Errorf("address loader scan error: %w", err)
	}

	return nil
}
