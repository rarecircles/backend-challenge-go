package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rarecircles/backend-challenge-go/env"
	"github.com/rarecircles/backend-challenge-go/types"
	"math/big"
	"net/http"
	"strconv"
)

func TokensHandler(w http.ResponseWriter, r *http.Request) {
	path := *env.SqliteFile
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		msg := "Internal Server Error - Can not connect to SQLite"
		w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(msg))
		return
	}

	type RetT struct {
		Tokens []types.TokenT `json:"tokens"`
	}
	tokens := make([]types.TokenT, 0)

	q, e := r.URL.Query()["q"]
	if !e {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte{})
		_ = db.Close()
		return
	}

	id := 0
	name := ""
	symbol := ""
	address := ""
	decimal := 0
	totalSupply := ""
	rows, err := db.Query("SELECT * FROM `tokens` WHERE `name` LIKE ?", "%"+q[0]+"%")

	if rows != nil {
		for rows.Next() {
			err = rows.Scan(&id, &name, &symbol, &address, &decimal, &totalSupply)
			if err != nil {
			}
			ts := new(big.Int)
			ts.SetString(totalSupply, 10)
			token := types.TokenT{
				Name:        name,
				Symbol:      symbol,
				Address:     address,
				Decimals:    uint64(decimal),
				TotalSupply: ts,
			}
			tokens = append(tokens, token)
		}
	}
	ret := RetT{
		Tokens: tokens,
	}

	_ = db.Close()

	jb := new(bytes.Buffer)
	_ = json.NewEncoder(jb).Encode(ret)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(jb.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jb.Bytes())

	return
}
