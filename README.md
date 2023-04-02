# RareCircles Coding Challenge

## Assumptions
The eth rpc client is very slow, and it takes more than a second sometimes to 
fetch details of a single token using rpc client, so I paid more attention to 
have cleaner and more performant code that seeding db with 14k. 
I just seeded 50 token address to the database and those addresses are present in
`data/addresses_i.json` file. If we want to add more tokens in the DB,
we will copy more address to the file.

### Prerequisite
* **Docker** - install and run docker. `brew install --cask docker`
* **Golang** - use `brew` to install golang

### How to run assignment
This is a dockerized app, to run this application make sure you have docker installed and running on your local system. Use the `docker-compose.yml` file in the project to run it.
* `cd` to project directory
* Run this command in terminal to up and start the service
* `docker-compose up --build`
* Once the service is up and running, open `postman` and call the `GET` endpoint. Read [API Specs](#api-specs) for more information

### API Specs
This project has one end-point.

#### `GET /tokens `
**Path** `http://localhost:10000/tokens?q="Year"` </br>
Endpoint to fetch a token details, by providing full/partial token title in query parameter.
The search is **case-insensitive**
```json
{
  "tokens": [
    {
      "name": "Yearn Finance Network Token",
      "symbol": "YFIXT",
      "address": "22f4a547ca569ae4dfee96c7aeff37884e25b1cf",
      "decimals": 18,
      "total_supply": 99988412498244662807623413580
    },
    {
      "name": "Yearn Utrade Finance",
      "symbol": "YUF",
      "address": "dbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
      "decimals": 18,
      "total_supply": 5000000000000000000000000
    }
  ]
}
```
### What I could have improved, if I had more time
1) I could have added few more end-points such as searching by address, symbol etc.
2) I should have leveraged elastic-search for searching instead of db.

### Note
Please feel free to reach out to me if you want to discuss any part of the code. Developed by @Shams Azad


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
  "tokesn": [
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
