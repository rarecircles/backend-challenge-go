package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	slice := []string{
		"string1",
		"string2",
		"string3",
		"string4",
		"string5",
		"string6",
	}
	chunk := Chunk(slice, 3)
	assert.Equal(t, len(chunk), 2)
	assert.Equal(t, len(chunk[0]), 3)
	assert.Equal(t, len(chunk[1]), 3)
}

func TestDecodeAddressesFile(t *testing.T) {
	jsonFile, err := os.Open("../data/addresses.jsonl")
	assert.Nil(t, err)

	addresses := DecodeAddressJsonL(jsonFile)
	assert.Equal(t, len(addresses), 14139)
}
