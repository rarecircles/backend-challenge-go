package main

import (
	"strings"
	"encoding/json"
	"io"
)

func readJsonString(jsonString string) []string {
	var addrs []string

	f := strings.NewReader(jsonString)
	dec := json.NewDecoder(f)

	for {
		var addr struct {
			Address string `json:"address"`
		}
		err := dec.Decode(&addr)
		if err == io.EOF {
			break
		}
		checkErr(err)

		addrs = append(addrs, string(addr.Address))
	}

	return addrs
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}