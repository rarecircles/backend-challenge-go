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

type tokenResponse struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Address     string `json:"address"`
	Decimals    uint64 `json:"decimals"`
	TotalSupply int64  `json:"totalSupply"`
}

func (app *application) seedDataAsync(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		app.logger.Fatal("Unable to read file", zap.Error(err))
	}
	f := strings.NewReader(string(b))
	dec := json.NewDecoder(f)
	app.db = createInMemoryDb()
	for {
		tA := tokenAddress{}
		if err := dec.Decode(&tA); err != nil {
			if err == io.EOF {
				break
			}
			app.logger.Fatal("failed to get address in file", zap.Error(err))
			panic(err)
		}
		tR := make(chan tokenResponse)
		go func(address, url string) {
			add, err := eth.NewAddress(address)
			if err != nil {
				app.logger.Fatal("Unable to create eth address", zap.Error(err))
				return
			}
			client := rpc.NewClient(url)
			detail, err := client.GetERC20(add)
			if err != nil {
				app.logger.Fatal("Unable to get token detail", zap.Error(err))
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
	if err := txn.Insert("address", tr); err != nil {
		app.logger.Error("Unable to insert record", zap.Error(err))
	}

	// Commit the transaction
	txn.Commit()
}

func getData(app *application, name string) []tokenResponse {
	resp := make([]tokenResponse, 0)
	// Create read-only transaction
	txn := app.db.Txn(false)
	defer txn.Abort()

	itr, err := txn.Get("address", "Name_prefix", name)
	if err != nil {
		app.logger.Error("Unable to get record in database", zap.Error(err))
	}

	for obj := itr.Next(); obj != nil; obj = itr.Next() {
		p := obj.(tokenResponse)
		resp = append(resp, p)
	}

	return resp
}

func createInMemoryDb() *memdb.MemDB {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"address": &memdb.TableSchema{
				Name: "address",
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
