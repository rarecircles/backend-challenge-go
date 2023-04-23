package storagerepo

import (
	"errors"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/models"
	"go.uber.org/zap"
	"math/big"
	"strconv"
)

type SR struct {
	redisClient iRedisClient
	logger      *zap.Logger
}

func New(l *zap.Logger, c iRedisClient) SR {
	sr := SR{redisClient: c, logger: l}
	return sr
}

func (sr SR) CreateIndex() error {
	_, err := sr.redisClient.Info()
	if err == nil {
		return err // index already exists
	}

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextField("symbol")).
		AddField(redisearch.NewTextField("address")).
		AddField(redisearch.NewNumericField("decimals")).
		AddField(redisearch.NewNumericField("totalSupply"))

	// Create the index with the given schema
	if err := sr.redisClient.CreateIndex(sc); err != nil {
		return err
	}

	return nil
}

func (sr SR) Store(token models.Token) error {
	// Create a document with an id and given score
	doc := redisearch.NewDocument(token.Address.String(), 1.0)
	doc.Set("name", token.Name).
		Set("symbol", token.Symbol).
		Set("address", token.Address.String()).
		Set("decimals", token.Decimals).
		Set("totalSupply", token.TotalSupply.String()).
		Set("baseTokenURI", token.BaseTokenURI)

	// Index the document. The API accepts multiple documents at a time
	if err := sr.redisClient.Index(doc); err != nil {
		return err
	}

	return nil
}

func (sr SR) Search(key string) ([]models.Token, error) {
	if key == "" {
		key = "*"
	} else {
		key = "*" + key + "*" // add wildcards
	}

	// Searching with limit and sorting
	docs, _, err := sr.redisClient.Search(
		redisearch.NewQuery(key).
			Limit(0, 10))

	if err != nil {
		return nil, err
	}

	var tokens []models.Token
	for _, doc := range docs {
		decimals, err := strconv.ParseUint(doc.Properties["decimals"].(string), 10, 64)
		if err != nil {
			panic(err)
		}

		totalSupply := new(big.Int)
		totalSupply, ok := totalSupply.SetString(doc.Properties["totalSupply"].(string), 10)
		if !ok {
			return nil, errors.New("error parsing totalSupply")
		}

		address, _ := eth.NewAddress(doc.Properties["address"].(string))
		token := models.Token{
			Name:         doc.Properties["name"].(string),
			Symbol:       doc.Properties["symbol"].(string),
			BaseTokenURI: doc.Properties["baseTokenURI"].(string),
			Address:      address,
			Decimals:     decimals,
			TotalSupply:  totalSupply,
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (sr SR) GetAllAddresses() (map[string]bool, error) {
	addresses := make(map[string]bool)
	offset := 0
	limit := 10000
	for {
		// Searching with limit and sorting
		docs, count, err := sr.redisClient.Search(
			redisearch.NewQuery("*").
				Limit(offset, limit))

		if err != nil {
			return nil, err
		}

		for _, doc := range docs {
			address, _ := eth.NewAddress(doc.Properties["address"].(string))
			addresses[address.String()] = true
		}

		if count < limit {
			break
		} else {
			offset += limit
		}
	}

	return addresses, nil
}
