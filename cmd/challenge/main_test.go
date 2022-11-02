package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rarecircles/backend-challenge-go/eth"
)

func TestParseInput(t *testing.T) {
	jsonFile, err := os.Open("../../data/testaddresses.jsonl")
	if err != nil {
		t.Fatal("Could not open file")
	}

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatal("Could not parse file")
	}

	addresses := ParseInput(data)
	if len(addresses) != 3 || addresses[0].Address != "0x22f4a547ca569ae4dfee96c7aeff37884e25b1cf" {
		t.Fatal("String content do not match expected")
	}
}

func TestFilterResults(t *testing.T) {

	var placeholder []eth.Token

	placeholder = append(placeholder, eth.Token{Name: "Xena"})
	placeholder = append(placeholder, eth.Token{Name: ""})

	var tests = []struct {
		input    string
		expected string
	}{
		{"xe", "Xena"},
		{"", ""},
	}

	var output []eth.Token

	for _, test := range tests {
		output = FilterResults(placeholder, test.input)
		for _, v := range output {
			if v.Name != test.expected {
				t.Error("Test Failed: {} inputted: {} received: {} expected: {}", test.input, output, test.expected)
			}
		}
	}

}
