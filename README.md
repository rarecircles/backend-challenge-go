# RareCircles Coding Challenge

## Requirements

Design an API endpoint that provides token information based on token title.

- The endpoint is exposed at `/tokens`
- The partial (or complete) token title is passed as a query string parameter `q`
- The endpoint returns a JSON response with an array suggested matches
  - Each suggestion has a name
  - Each suggestion has a symbol
  - Each suggestion has a address
  - Each suggestion has a decimals
  - Each suggestion has a total supply
- At-least 2 go-test needs to be implemented
- Concurrency should be applied where appropriate
- Feel free to add more features if you like!

#### Notes

- You have a list of tokens accounts defined inthe `addresses.jsonl` file that should be used as a seed to your application
- We have included a `eth` lirbary that should help with any decoding and rpc calls

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
      "totalSupply": 1000000000
    },
    {
      "name": "Rareible",
      "symbol": "RRI",
      "address": "0xe9c8934ebd00bf73b0e961d1ad0794fb22837206",
      "decimals": 9,
      "totalSupply": 100
    }
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

- All code should be written in Golang, Typescript or PHP.
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

We know that the time for this project is limited, and it is hard to create a "perfect" solution, so we will consider that along with your experience when evaluating the submission.

## Getting Started

### Prerequisites

You are going to need:

- `Git`
- `go`

### Starting the application with docker-compose

```
make all
```

### To skip seeding data

- set ENV `SEED_DATA=false` in docker-compose

### How to search Elastic Search on local

```
curl -X GET "localhost:9200/tokens/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match": {
      "name": {
        "query": "token_name",
        "fuzziness": "AUTO"
      }
    }
  }
}
'
```

### Endpoints

- Healthcheck
  `GET http://localhost:8080/healthcheck`

- Query Tokens
  `GET http://localhost:8080/tokens?q=<token_name>`

### Folder Structure

```
.
├── Makefile
├── README.md
├── cmd
│   └── challenge
│       └── main.go
├── data
│   └── addresses.jsonl
├── docker
│   └── dockerfile.token-api
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── handler
│   │   │   ├── health_check.go
│   │   │   └── token_grp
│   │   │       ├── dto.go
│   │   │       ├── token_grp.go
│   │   │       └── token_grp_test.go
│   │   ├── middleware
│   │   │   ├── rate_limiter.go
│   │   │   └── timeout.go
│   │   └── server.go
│   ├── pkg
│   │   ├── address_loader
│   │   │   ├── address_loader.go
│   │   │   ├── address_loader_test.go
│   │   │   └── testdata
│   │   │       └── addresses_test.jsonl
│   │   ├── eth
│   │   └── rpc
│   └── service
│       └── search
│           ├── mock
│           │   └── search.go
│           └── search.go
└── pkg
    └── logger
```

### TODO

- Add getToken endpoint`GET http://localhost:8080/tokens/:symbol`
- Add more test code
- Configuration using viper
- Fix TotalSupply long type error issue in ElasticSearch when seeding data
- Set request id and logging middelware for tracing
- Swagger
