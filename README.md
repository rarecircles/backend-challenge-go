# RareCircles Coding Challenge

## Requirements

Design an API endpoint that provides token information based on token title.

- the endpoint is exposed at `/tokens`
- the partial (or complete) token title is passed as a query string parameter `q`
- the endpoint returns a JSON response with an array suggested matches
  - each suggestion has a name
  - each suggestion has a symbol
  - each suggestion has a address
  - each suggestion has a decimals
  - each suggestion has a total supply
- at-least 2 go-test needs to be implemented
- concurrency should be applied where appropriate  
- feel free to add more features if you like!

#### Notes
- you have a list of tokens accounts defined inthe `addresses.jsonl` file that should be used as a seed to your application
- we have included a `eth` lirbary that should help with any decoding and rpc calls

#### Sample responses

These responses are meant to provide guidance. The exact values can vary based on the data source and scoring algorithm.

**Near match**

    GET /tokens?q=rare

```json
{
  "tokens": [
    {
      "name": "RareCircles",
      "symbol": "RCI",
      "address": "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
      "decimals": 18,
      "totalSupply": 1000000000,
    },
    {
      "name": "Rareible",
      "symbol": "RRI",
      "address": "e9c8934ebd00bf73b0e961d1ad0794fb22837206",
      "decimals": 9,
      "totalSupply": 100,
    },
  ]
}
```

**No match**

    GET /tokens?q=ThereIsNoTokenHere

```json
{
  "tokens": []
}
```


### Non-functional

- All code should be written in Goland, Typescript or PHP.
- Mitigations to handle high levels of traffic should be implemented.
- Challenge is submitted as pull request against this repo ([fork it](https://help.github.com/articles/fork-a-repo/) and [create a pull request](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)).
- Documentation and maintainability is a plus.

## Dataset

You can find the necessary dataset along with its description and documentation in the [`data`](data/) directory.

## Evaluation

We will use the following criteria to evaluate your solution:

- Capacity to follow instructions
- Developer Experience (how easy it is to run your solution locally, how clear your documentation is, etc)
- Solution correctness
- Performance
- Tests (quality and coverage)
- Attention to detail
- Ability to make sensible assumptions

It is ok to ask us questions!

We know that the time for this project is limited and it is hard to create a "perfect" solution, so we will consider that along with your experience when evaluating the submission.

## Getting Started

### Prerequisites

You are going to need:

- `Git`
- `go`

### Starting the application

To start a local server run:

```
go run ./cmd/challenger
```


# About This

This is a golang backend code challenge by Ian Mah.

### Requirements & dependency
* Golang runtime, with proper environment variables configured.
* Git
* *nix (suggested, but Windows Platform should also be OK, not tested)
### How to start
* clone code
```shell
git clone https://github.com/stgmsa/backend-challenge-go
```
* modify backend-challenge-go/env/env.go and point the sqlite data file to where it really locates.
e.g. if you checkout this repository at /home/someone/go/backend-challenge-go
then you should modify the last line of env.go like this:
```go
var SqliteFile = flag.String("sqlite-file", "/home/someone/go/backend-challenge-go/cache.sqlite", "sqlite file")
```

* launch the program by typing the following command in your terminal
```shell
go run ./cmd/challenge
```

* to run a go test
```text
cd gotest && go test
```

### Some hints might be useful
* this server runs at port 8080, make sure it is not occupied by other programs.
* if you are using a Mac running OSX 10.15 or later, and cloned the code into a sub folder of your home folder, you might to give a file read/write permission when launches this.


### Folder structure
```text
.
├── README.md
├── cache.sqlite                  // sqlite file for caching token entries
├── challenge
├── cmd
│   └── challenge
│       ├── logging.go
│       └── main.go               // entry point of go
├── data
│   └── addresses.jsonl           // address data entries, not modified
├── env
│   └── env.go                    // parameters and settings.
├── eth
│   ├── crypto.go
│   ├── decoder.go
│   ├── decoder_log.go
│   ├── encoder.go
│   ├── errors.go
│   ├── interface.go
│   ├── log_event.go
│   ├── logging.go
│   ├── method.go
│   ├── nft.go
│   ├── numbers.go
│   ├── rpc
│   │   ├── errors.go
│   │   ├── json_encode.go
│   │   ├── json_indent.go
│   │   ├── json_scanner.go
│   │   ├── json_stream.go
│   │   ├── logging.go
│   │   ├── rpc.go
│   │   ├── token.go
│   │   └── token_helper.go
│   ├── token.go
│   ├── types.go
│   └── utils.go
├── go.mod
├── go.sum
├── gotest                              // go test directory
│   └── handler_test.go                 // go test
├── handlers
│   └── handlers.go                     // registered handlers of endpoint
├── logging
│   ├── context.go
│   ├── core.go
│   └── logging.go
└── types
    └── types.go                        // some structs defined and used by json encoding/decoding
```

### What will this program do
1. initial environment and data file
2. read entries from data file and gets 10 new addresses (for further API call)
3. call token API asynchronous and send results to channel for further consumption
4. launch cache goroutine, fetch from channel and save each item to sqlite.
5. serve general API query using data stored in sqlite

### Some issues already known
1. Q: why sqlite with absolute path
A: sqlite is lightweight and portable language providing sql-like interface, easier for small application development.
But for this situation, a database system / search enging supporting reverse-indexing would be a better solution 
(e.g. postgre/ ElasticSearch) for productive environment, but this may introduce more components. 
For productive environment, I would say yes to use a replacement instead of sqlite, for now, No. because it's more portable
and does not require a docker / image shipping

2. Q: why load 10 new items at the beginning of this program 
A: instead of using fixed items stored in DB, this will increase the data entries on each start. More likely to be 
a normal program activity.

3. Q: Any issue with sqlite?
A: Yes. 1st time using sqlite with the help of google, that's why I finally used absolute path for sqlite data file.
Because in go test and go run (launches the server), the data file path using relative path are different and I found this
until the last minute of debugging.

4. Why native builtin http server
A: simple, easier, less code, although we have better choices with enhanced routing performance
and rich libraries and encapsulation.

### Performance
1. Performance is not good. We didn't use a server based data storage system (e.g. mysql, ElasticSearch, with cache enabled
and better tunned algorithm for searching). but it will do with normal loads.

### Endpoints and APIs
```text
endpoint: /token
queryparameter: q
```
* if access this endpoint without a "q" parameter, it will generate a HTTP 400 error.
* if access this endpoint with q= (q equals to empty) will return all entries in sqlite.
* q is case insensitive

### Sensitive data
* configurations and API keys are hard coded in go files. that is not a good idea for developing.
* some data paths can also be found in source file.

### DDL
```sql
CREATE TABLE IF NOT EXISTS "tokens"
(
    id          INTEGER      primary key autoincrement,
    name        VARCHAR(128) default "" not null,
    symbol      VARCHAR(128) default "" not null,
    address     VARCHAR(128) default "" not null,
    decimals    BIGINT       default 0 not null,
    totalSupply VARCHAR(128) default "" not null
);
```
