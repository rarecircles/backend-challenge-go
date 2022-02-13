package env

import "flag"

// This could be written in a configuration file or dynamically loaded from a system env variable.
var FlagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var FlagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")
var FlagRPCToken = flag.String("rpc-token", "6Jo0D-Kh1mLj4jsnEV_xFCegoJCfWwCJ", "Alchemy API Token")
var AddressFile = flag.String("address-file", "./data/addresses.jsonl", "address book")
var SqliteFile = flag.String("sqlite-file", "/Users/aeon/go/src/github.com/stgmsa/backend-challenge-go/cache.sqlite", "sqlite file")
