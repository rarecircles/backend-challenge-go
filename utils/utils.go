package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/rarecircles/backend-challenge-go/cmd/challenge/types"
)

func DecodeAddressJsonL(jsonFile *os.File) []types.Address {
	var addresses []types.Address
	dec := json.NewDecoder(jsonFile)
	for {
		var address types.Address

		err := dec.Decode(&address)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		addresses = append(addresses, address)
	}
	return addresses
}

// Chunk breaks down a slice of strings into a slice of multiple slices for batch processing
func Chunk(xs []string, chunkSize int) [][]string {
	if len(xs) == 0 {
		return nil
	}

	divided := make([][]string, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize

	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}

	divided[i] = xs[prev:]
	return divided
}

func Work(limiter chan bool, group *sync.WaitGroup, worker int) {
	limiter <- true
	defer func() {
		<-limiter
	}()
	time.Sleep(time.Second * 1)
}

func Retry(attempts int, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(i+1*2) * time.Second)
		log.Println("retrying after error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
