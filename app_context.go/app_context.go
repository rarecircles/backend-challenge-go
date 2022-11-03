package appContext

import "github.com/rarecircles/backend-challenge-go/eth/rpc"

type appContext struct {
	rpcClient rpc.Client
}

var appCtx *appContext

func Init(rpcUrl string) {
	appCtx = &appContext{}
	appCtx.rpcClient = initRPCClient(rpcUrl)
}

func initRPCClient(url string) rpc.Client {
	return *rpc.NewClient(url)
}

func GetRPCClient() rpc.Client {
	return appCtx.rpcClient
}
