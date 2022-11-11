package main

import (
	"os"
	"testing"
)

func setupTest() *application {
	*flagRPCURL = *flagRPCURL + os.Getenv("ethKey")
	app := &application{
		logger: zlog,
		rpcUrl: *flagRPCURL,
	}

	path := os.Getenv("datapath")
	app.db = createInMemoryDb()
	app.seedDataAsync(path)
	return app
}

func Test_getData_ReturnsData(t *testing.T) {
	app := setupTest()

	data := getData(app, "yearn")
	if len(data) <= 0 {
		t.Error("Test failed with empty response")
	}
}

func Test_getData_ReturnsEmptyArray(t *testing.T) {
	app := setupTest()

	data := getData(app, "ThereIsNoTokenHere")

	if len(data) > 0 {
		t.Error("Test failed with token response not empty")
	}
}
