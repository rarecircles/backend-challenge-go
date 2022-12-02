package sql

import (
	"github.com/go-pg/pg/v10"
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/types"
)

var insertTokenQuery = `
INSERT INTO tokens(
	name,
	symbol,
	address,
	decimals,
	total_supply
) VALUES
	(?,?,?,?,?)
`

func InsertToken(t types.Token) (string, error) {
	// Could convert this to an upsert if we use address as primary key
	db := Connect()
	_, err := db.QueryOne(&t, insertTokenQuery, t.Name, t.Symbol, t.Address, t.Decimals, t.TotalSupply)
	defer db.Close()
	return "", err
}

var tokensFetchQuery = `
SELECT
	name,
	symbol,
	address,
	decimals,
	total_supply
FROM tokens
WHERE name LIKE '%' || ? || '%'
`

func TokensFetch(name string) ([]types.Token, error) {
	db := Connect()
	tokens := []types.Token{}
	_, err := db.Query(&tokens, tokensFetchQuery, name)
	defer db.Close()
	return tokens, err
}

var tokensCountQuery = `
SELECT COUNT(address)
FROM tokens
`

func TokensCount() (uint64, error) {
	db := Connect()
	var count uint64
	_, err := db.QueryOne(pg.Scan(&count), tokensCountQuery)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	return count, err
}
