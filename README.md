# RareCircles Coding Challenge

## [Original instructions file](docs/CHALLENGE_INSTRUCTIONS.md)

# Requirements
* Docker

# Running challenge service

To run challenge service:
1. Move .env.example to .env: `mv .env.example .env`
2. Configure `.env` file with necessary information
3. Run `make deploy-challenge`
4. Wait for seed data to populate. Check logs via `make challenge-logs` to determine seed status
5. When seeding is done, make a search query ie. `curl http://localhost:8080/tokens?q=coin`


# Cleaning up
This will remove all related containers, images, and volumes
`make cleanup`

# Unit tests
This will run all unit tests (currently only have unit tests for utils)
`make unit-test`


# Things to do/consider (would need more time):
* Use middleware for rpcClient & Postgres
* Create integration tests
* Convert DAL insert requests into upsert requests (insert or update)
* Implement redis middleware for caching user queries