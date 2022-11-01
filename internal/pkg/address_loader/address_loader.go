// Package address_loader loads address files
package address_loader

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type AddressLoader interface {
	Load(filepath string) error
}

type addressLoader struct {
	log *zap.Logger
}

// NewAddressLoader creates AddressLoader
func NewAddressLoader(log *zap.Logger) AddressLoader {
	return &addressLoader{
		log: log,
	}
}

func (al addressLoader) Load(filepath string) error {
	f, err := os.Open(filepath)
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
		al.log.Info(s.Text())
		if len(s.Text()) == 56 {
			al.log.Info(s.Text()[12:54])
		}
	}

	if err := s.Err(); err != nil {
		al.log.Fatal(err.Error())
	}

	return nil
}
