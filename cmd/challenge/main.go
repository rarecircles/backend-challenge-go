package main

import (
	"flag"
	"go.uber.org/zap"
)

const Name = "challenge"

func main() {
	config := getConfig(Name)
	zlog := getLogger(Name)

	rpc.SetLogger(zlog)
	eth.SetLogger(zlog)


