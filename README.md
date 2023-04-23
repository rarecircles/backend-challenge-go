## Running the project

You need to put your Alchemy API key in the config.toml file for the key `ETH_RPC_KEY`.
You can get your free API key from [here](https://www.alchemy.com/).

### With Docker

You can simply run the project with:

```bash
docker compose up --build
```

Or in older versions of docker compose:

```bash
docker-compose up --build
```

The project will immediately start running on port 8080 and you can make requests to it after seeing this message in the logs on stdout.

```
Listening and serving HTTP on 0.0.0.0:8080
```

However, it does take the project some time to fetch all the data from the API and store it, So in the beginning you might get fewer items in response than expected.
<br/>
A sample request can be made with:

```bash
curl http://localhost:8080/tokens\?q\=net | jq
```

### Without Docker

You'll need a running instance of RedisSearch on your machine.
You can either install it directly on your machine or run an instance with Docker.

```bash
docker run -p 6379:6379 redis/redis-stack-server:latest
```

Then you'll need to update the config.toml file with the correct host and port for RedisSearch.
You'll need to also add your Alchemy RPC key to the config.toml file.
<br />
Then you can run the project with:

```bash
go run ./cmd/challenge
```

### Running the tests

After installing Ginkgo and Gomega, you can run the tests with:

```bash
ginkgo -r
```

This project has 5 tests that cover 2 of the packages.
This test coverage is not enough for a production application,
but for demonstration proposes and considering the scope and time limit of this project it hopefully is sufficient.

## Thought process, decisions, and assumptions

### Handling configurations

The first decision was to have all the configurations in a single config.toml file.
I chose to use Viper because it's a very popular library for handling configurations in Go.
With the help of Viper we read and parse the configurations stored in the config.toml file.

In order to have more flexibility configurations can also be overwritten with environment variables.
For instance if we want to overwrite `REDIS_SEARCH_URL` we can do so by setting the environment variable `CHALLENGE_REDIS_SEARCH_URL` to the desired value.
Which we have done in the `Dockerfile` and `docker-compose.yml` files. The prefix is added to distinguish the environment variables related to this project.

### Why layered architecture and the use of interfaces?

Layered architecture is a very common pattern in software development.
The upsides of using this architecture are simplicity, testability, maintainability, and flexibility.

For instance if we later in the future want to use `ElasticSearch` instead of `RedisSearch` or read
the token address from a database instead of a file, we can easily do so by implementing the interfaces and without changing the business logic.

Layered architecture also makes it easier to test the code. We can easily mock the dependencies and test the business logic in isolation.

### Why loaders?

For this program we need to load all the token data from the API and store it in RedisSearch.
We want to do this concurrently in order to make the project able to handle user requests from the beginning instead of waiting for all the data to be loaded.

We might also want to update the list of token addresses every once in a while, so this loader has the ability to refetch
the data after a certain amount of time (like a cronjob). However, we don't overwrite already existing data since smart contracts are immutable and we don't expect them to change often.

Refetching can be turned off by setting `REFETCH_DELAY_HOURS` to zero in configs, so we will just fetch data once the program starts.

### Concurrency

Since accessing the RPC is a blocking operation, we want to do it concurrently in order to make the program faster.
However, we want to have limited concurrency not to overload our server or reach the rate limit of the API.

For this propose we could either have a pool of workers or use a semaphore. I chose to use a semaphore since it's simpler and easier to implement.

### Why RedisSearch?

- Searching and Suggestions are the main features of this program.
- Our program needs to be fast, scalable and efficient.
- All our data can easily fit in memory.
- We need full-text search and prefix, suffix, and infix search.
- We don't need to do any complex queries.

Other candidates and why they were not chosen:

- BigQuery: Can't persist the data, can't sync instances if we have replication, not suitable for our queries (prefix, suffix, infix, full-text search).
- ElasticSearch: Difficult to maintain, not as fast as RedisSearch, needs more resources.
- Postgres: Needs more resource, not as fast as RedisSearch (data is on disk).

RedisSearch is the best choice since our data can easily fit on memory, and we need to do simple queries.
<br/>
In the future if we need to do more complex queries or store much more data, we can easily replace RedisSearch with ElasticSearch because our architecture allows that.

### Why Ginkgo?

Ginkgo is a testing framework for Go designed to help you write expressive tests. It helps us write BDD tests very easily and descriptive.

### Why Gin?

Very fast, easy to use, a lot of features, and a huge community.

### Handling rate limits

Reaching the rate limit of the API is a very common problem in this kind of projects. I handle this problem with two approaches.

- First, I use a semaphore to limit the number of concurrent requests to the API.
- Second, I use a retry mechanism to retry the request if it fails due to rate limit.

The retry mechanism is simply a fixed time delay between each retry.
This is not the best approach since we don't know how long we should wait between each retry.
We can implement a better retry mechanism by using exponential back-off in the future.

### Support all Tokens

The project supports tokens that are either ERC20, ERC721, or ERC1155.

### Persisting data and picking up where we left off

The loader in this project won't refetch the data tha it has already fetched and stored in redis.
We can easily take advantage of this by sharing a redis between instances of the program and persisting its data.

### Separating loader and server

The loader and server can become separated application in the future in order to make the program more scalable and efficient, especially if we want to have multiple instances of the server.
