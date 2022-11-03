package service

import (
	appContext "github.com/rarecircles/backend-challenge-go/app_context.go"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
)

type ServerDependencies struct {
	TokenService TokenService
	RpcClient    rpc.Client
}

func InstantiateServerDependencies() *ServerDependencies {
	rpcClient := appContext.GetRPCClient()
	tokenService := NewTokenService(rpcClient)

	return &ServerDependencies{
		RpcClient:    rpcClient,
		TokenService: tokenService,
	}
}
