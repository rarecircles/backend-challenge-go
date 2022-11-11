package main

import (
	"encoding/json"
	"github.com/hashicorp/go-memdb"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
)

const addressTable = "address"

type tokenResponse struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Address     string `json:"address"`
	Decimals    uint64 `json:"decimals"`
	TotalSupply int64  `json:"totalSupply"`
}

// seedDataAsync takes in a path that contains a file with json data of addresses
// stores in an in memory database
//
//line by line, retrieves information about the addresses
func (app *application) seedDataAsync(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		app.logger.Fatal("Unable to read file", zap.Error(err))
		return
	}
	f := strings.NewReader(string(b))
	dec := json.NewDecoder(f)
	for {
		tA := tokenAddress{}
		if err := dec.Decode(&tA); err != nil {
			if err == io.EOF {
				break
			}
			app.logger.Error("failed to get address in file", zap.Error(err))
			panic(err)
		}
		tR := make(chan tokenResponse)
		go func(address, url string) {
			add, err := eth.NewAddress(address)
			if err != nil {
				app.logger.Error("Unable to create eth address", zap.Error(err))
				return
			}
			client := rpc.NewClient(url)
			detail, err := client.GetERC20(add)
			if err != nil {
				app.logger.Error("Unable to get token detail", zap.Error(err))
				return
			}

			tR <- tokenResponse{
				Name:        detail.Name,
				Symbol:      detail.Symbol,
				Address:     detail.Address.String(),
				Decimals:    detail.Decimals,
				TotalSupply: detail.TotalSupply.Int64(),
			}
		}(tA.Address, app.rpcUrl)

		data := <-tR
		if data.Name != "" {
			go writeOneData(app, data)
		}

	}
}

func writeOneData(app *application, tr tokenResponse) {
	txn := app.db.Txn(true)
	if err := txn.Insert(addressTable, tr); err != nil {
		app.logger.Error("Unable to insert record", zap.Error(err))
	}

	// Commit the transaction
	txn.Commit()
}

// getData queries a database to retrieve token information for query parameter that is prefixed by
// name input
func getData(app *application, name string) []tokenResponse {
	return getDataFromDatabase(app, name)
}

func createInMemoryDb() *memdb.MemDB {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			addressTable: &memdb.TableSchema{
				Name: addressTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Address"},
					},
					"Name": &memdb.IndexSchema{
						Name:    "Name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name", Lowercase: true},
					},
				},
			},
		},
	}

	// Create a new database
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return db
}

func getDataFromDatabase(app *application, name string) []tokenResponse {
	// Create read-only transaction
	txn := app.db.Txn(false)
	defer txn.Abort()

	itr, err := txn.Get(addressTable, "Name_prefix", strings.ToLower(name))
	if err != nil {
		app.logger.Error("Unable to get record in database", zap.Error(err))
	}

	resp := make([]tokenResponse, 0)
	for obj := itr.Next(); obj != nil; obj = itr.Next() {
		p := obj.(tokenResponse)
		resp = append(resp, p)
	}
	return resp

}
