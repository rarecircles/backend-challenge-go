package main

import (
	"github.com/rarecircles/backend-challenge-go/cmd/app"
)

func main() {
	app := app.App{}
	app.HandleRequests()
}
