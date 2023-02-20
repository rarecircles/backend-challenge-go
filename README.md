# RareCircles Coding Challenge - Jordan Coil

## Assumptions and Technical Decisions

- I wrote tests for one of the util functions to show understanding of writing tests in Go, but given the time limitation, I did not write any tests to test methods that touch the database.
- I assumed that passing in an empty query value to the tokens endpoint was not valid, as that would have returned all tokens.
- I'm only running a simple http sever, which might not be able to handle a ton of requests, but given the project is an API with a single GET endpoint, I think it is sufficient. If I were to update the project to be able to handle more requests and more endpoints, I would utilize go routines to dispatch workers to handle larger processes (ie. a POST request that had to process alot of data.)
  - (probably following what is outlined in this [blog article](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/))

## How to run tests

```
go test ./cmd/challenge
```

## How to run project

### setting up the project Database using docker

*These commands were tested on a window machine using cmd*

- ensure the machine has docker installed
- run the following command to start the postgres docker container: 

```
docker run --name postgres-server -e POSTGRES_PASSWORD=secret -d -p 6000:5432 postgres
```

- run the following command to open a terminal window on the postgres container: 

```
docker exec -it postgres-server /bin/bash
```

- run the following command to access the postgres cli, and enter password ``secret``

```
psql -d postgres -U postgres -W
```


- copy and paste the sql statment in the ```migrations/init.sql``` to create the table to store tokens
- you can now clsoe this terminal/cmd window and the postgres server docker container should be running in the background

### running the project

```
go run ./cmd/challenge --api-key <alchemy_api_key>
```

- please provide your own alchemy api key as a flag to the program
- if running for the first time the database will seed in the background. You should be able to access endpoints while this is happening.
  - *note: if you stop the application and run it again, the application will not seed if there are any records in the tokens database*
- test out the endpoint by going to ```localhost:8080/tokens?q=<token_name>```
- replace ```<token_name>``` in the URL params with some full or partial token name
- some json containing info about the matching tokens should be displayed
