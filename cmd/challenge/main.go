package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rarecircles/backend-challenge-go/env"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/handlers"
	"github.com/rarecircles/backend-challenge-go/types"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func LoadAddress(path string, c chan types.AddressT, offset int) {
	f, err := os.Open(path)

	defer func() {
		e := f.Close()
		if e != nil {
			zlog.Error("Load Address Book", zap.String("path", path), zap.String("stage", "close"))
		}
	}()
	if err != nil {
		zlog.Error("Load Address Book", zap.String("path", path), zap.String("stage", "open"))
	}
	var address types.AddressT

	s := bufio.NewScanner(f)
	lineNo := 0
	for s.Scan() {
		jsonL := s.Text()
		if len(jsonL) > 20 {
			lineNo++
			if lineNo > offset && (lineNo-offset <= 10) {
				err := json.Unmarshal([]byte(jsonL), &address)
				if err != nil {
					zlog.Error("Json Unmarshal", zap.String("json", jsonL))
				}
				c <- address
			}
		}
	}
}

func QueryERC20(ca chan types.AddressT, ct chan types.TokenT) {
	rpcURL := *env.FlagRPCURL
	alchemyToken := *env.FlagRPCToken
	for {
		address := <-ca

		a, _ := eth.NewAddress(address.Address)
		c := rpc.NewClient(rpcURL + alchemyToken)

		token, _ := c.GetERC20(a)

		t := types.TokenT{
			Name:        token.Name,
			Symbol:      token.Symbol,
			Address:     token.Address.String(),
			Decimals:    token.Decimals,
			TotalSupply: token.TotalSupply,
		}
		ct <- t

		zlog.Info("token", zap.String("token", token.String()))
	}
}

func GetCount() (int, error) {
	path := *env.SqliteFile
	db, err := sql.Open("sqlite3", path)
	defer db.Close()
	if err != nil {
		return 0, err
	}

	rows, err := db.Query("SELECT COUNT(*) AS `c` FROM `tokens`")
	c := 0
	for rows.Next() {
		err = rows.Scan(&c)
		if err != nil {
			return 0, err
		}
	}
	return c, nil
}

func ToSqlite(c chan types.TokenT) int64 {
	path := *env.SqliteFile
	db, err := sql.Open("sqlite3", path)
	defer db.Close()
	if err != nil {
	}
	stmt, err := db.Prepare("INSERT INTO tokens(`name`, `symbol`, `address`, `decimals`, `totalSupply`) values(?,?,?,?,?)")
	if err != nil {
		fmt.Println("Error prepare")
	}
	for {
		token := <-c
		_, _ = stmt.Exec(token.Name, token.Symbol, token.Address, token.Decimals, token.TotalSupply.String())
	}
}

func main() {
	flag.Parse()
	httpListenAddr := *env.FlagHTTPListenAddr
	addressFilePath := *env.AddressFile

	AddressChan := make(chan types.AddressT, 10)
	TokenChan := make(chan types.TokenT, 10)

	var err error

	count, err := GetCount()
	if err == nil {
		go LoadAddress(addressFilePath, AddressChan, count)
		go QueryERC20(AddressChan, TokenChan)
		go ToSqlite(TokenChan)
	}

	zlog.Info("Running Challenge", zap.String("httpL_listen_addr", httpListenAddr))

	// register /tokens URI, tokens endpoint.
	// We've chosen golang's builtin http server, it's a simple one, easier for development & demonstration.
	// But for online environment, we have better choices. There are many golang web frameworks. e.g. gin
	http.HandleFunc("/tokens", handlers.TokensHandler)
	err = http.ListenAndServe(httpListenAddr, nil)
	if err != nil {
		zlog.Error("Launching HTTP Server",
			zap.String("Port", httpListenAddr),
			zap.String("Err", err.Error()),
		)
	}
}
