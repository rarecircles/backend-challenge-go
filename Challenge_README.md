# Rarecircles Technical Assessment

## Overview

The aim is to create an API endpoint that provides ethereum ERC20 token information based on token name.

## Highlights

1. Object-oriented design
2. Follows go testing procedures
3. Uses pipeline to parse addresses and read tokens

## Directory Tree

For better understanding the project, here is a directory tree:

```
|-- LICENSE
|-- README.md
|-- Challenge_README.md
|-- go.mod
|-- go.sum
|-- .gitignore
|-- cmd
    |-- challenge
        |-- logging.go
        |-- main.go
|-- internal
    |-- processor
        |-- logging.go
        |-- processor.go
        |-- processor_test.go
|-- test
    |-- input 
        |-- seed_data.jsonl
    |-- output
        |-- ThereIsNoToken
|-- scripts
    |-- build.sh
    |-- test.sh
    |-- concurrent_request.sh
|-- logging (default)  
|-- eth (default)  
|-- data (default)   
```

## Assumptions

- The project runs on Ubuntu 20.04 with Golang version 1.18.
- JSONL files are considered as a valid input. But the code can be extended to other input formats.
- Eth rpc server capable of providing token information

## Go installation

Download and install [Go](https://go.dev/doc/install). Here is an example in Linux:

```
$ wget https://dl.google.com/go/go1.18.4.linux-amd64.tar.gz
$ rm -rf /usr/local/go 
$ tar -C /usr/local -xzf go1.18.4.linux-amd64.tar.gz
```

## Test

Provide a valid rpc-URL before running the test

### With building the source

```
.scripts/build.sh
./build/challenge  ...
```

### Without building the source

```
go run ./cmd/challenge/  ...
```

## For Developers

You are welcome to open issues or pull requests for this project. You can run unit tests with 

```
go test -v ./...
```

You should see something similar as follows:

```
?       github.com/rarecircles/backend-challenge-go/cmd/challenge   [no test files]
?       github.com/rarecircles/backend-challenge-go/eth [no test files]
?       github.com/rarecircles/backend-challenge-go/eth/rpc [no test files]
=== RUN   TestSeedParsing
--- PASS: TestSeedParsing (13.46s)
=== RUN   TestHttp
--- PASS: TestHttp (13.41s)
PASS
ok      github.com/rarecircles/backend-challenge-go/internal/processor  26.871s
?       github.com/rarecircles/backend-challenge-go/logging [no test files]
```

## Code Style

This project is formatted using go fmt tool to keep code style consistent:

```
go fmt ./...
```

## Future Works

In the future, we should build a docker image that enables the application to run on any platform.
Using trie data structure to record tokens for faster search. 
Real-time updates tokens.

## Developer Notes

The default rpc URL (https://eth-mainnet.alchemyapi.io/v2/) does not work (http 404 error, page not found). The developer tested with other Ethereum rpc URL searched online. Those connections does not support high speed token reading. Concurrent requests result in http 429 error (too many requests). Thus, the bootstrapping phase depends on the performance of the Ethereum rpc server.
