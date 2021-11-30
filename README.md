# RareCircles Coding Challenge

This is the version of the Rarecircles coding challenge developed by Jose Camilo.

### Requirements

- Install Docker
- Install Git
- Clone the repository. `git clone git@github.com:jose-camilo/backend-challenge-go.git`
- Checkout the `backend-challenge` branch. `git checkout backend-challenge`

### To start the application
- Open your terminal.
- Go to the repository folder.
- Run these commands. 
```
  make netsetup
  make all
```
The search engine will print on console the tokens it's working with and will take some time to load and index.

#### Tokens Endpoint:
`GET http://0.0.0.0:1234/v1/tokens?q=<token name/partial name>`
 
#### Symbols Endpoint:
`GET http://0.0.0.0:1234/v1/symbols?q=<symbol>`

#### Addresses Endpoint:
`GET http://0.0.0.0:1234/v1/addresses?q=<address>`
Remove `0x` from the address to do a search. `eg. 22f4a547ca569ae4dfee96c7aeff37884e25b1cf`


# Coding Challenge Specifications
## Requirements

Design an API endpoint that provides token information based on token title.

- The endpoint is exposed at `/tokens`
- The partial (or complete) token title is passed as a query string parameter `q`
- The endpoint returns a JSON response with an array suggested matches
  - Each suggestion has a name
  - Each suggestion has a symbol
  - Each suggestion has an address
  - Each suggestion has a decimals
  - Each suggestion has a total supply
- At-least 2 go-test needs to be implemented
- Concurrency should be applied where appropriate  
- Feel free to add more features if you like!

#### Notes
- You have a list of token accounts defined in the `addresses.jsonl` file that should be used as a seed to your application
- We have included an `eth` library that should help with any decoding and rpc calls

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
      "address": "e9c8934ebd00bf73b0e961d1ad0794fb22837206",
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

We know that the time for this project is limited, and it's hard to create a "perfect" solution, so we will consider that along with your experience when evaluating the submission.

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
