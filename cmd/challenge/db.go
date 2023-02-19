package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
)

func DBConnect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    	"password=%s dbname=%s sslmode=disable",
    	host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)

	return db
}

func seedDB(rpcURL string, db *sql.DB) {
	isSeeded := checkIfSeeded(db)

	if !isSeeded {
		// Read json addresses in seed file
		jsonContent, err := ioutil.ReadFile("./data/addresses.jsonl")
		checkErr(err)

		var addrs []string 
		addrs = readJsonString(string(jsonContent))

		// TODO: move api key to a flag
		apiKey := "RtBNZI7jboJBSVutqQidtcUE8Nbw2M6p"
		client := rpc.NewClient(rpcURL + apiKey)

		addTokensToDB(addrs, client, db)
	}
}

func addTokensToDB(addrs []string, client *rpc.Client, db *sql.DB) {
	for _, a := range addrs {
		decodedHexAddr := eth.MustDecodeString(a)

		// Get token for address
		t, err := client.GetERC20(decodedHexAddr)
		if err != nil {
			// Don't want to exit if only 1 address isn't good
			continue
		}

		insStatement := `INSERT INTO tokens (name, symbol, address, decimals, total_supply)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;`
		id := 0
		row := db.QueryRow(insStatement, 
			t.Name, 
			t.Symbol, 
			t.Address.Pretty(),
			strconv.FormatUint(t.Decimals, 10),
			t.TotalSupply.String())
		switch err := row.Scan(&id); err {
			case sql.ErrNoRows:
				fmt.Println("No row inserted!")
			case nil:
				fmt.Println("Inserted row")
			default:
				checkErr(err)
		}
	}
}

func checkIfSeeded(db *sql.DB) bool {
	sqlStatement := `SELECT (id) FROM tokens;`
	var id int
	row := db.QueryRow(sqlStatement)
	
	switch err := row.Scan(&id); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return false
		case nil:
			return true
		default:
			checkErr(err)
	}
	return false
}

func queryToken(queryName string, db *sql.DB) []*TokenModel {
	queryName = "%" + queryName + "%"

	results, err := db.Query(`SELECT * FROM tokens WHERE name LIKE $1`, queryName)
	checkErr(err)

	var tokens []*TokenModel
	for results.Next() {
		var t TokenModel

		err = results.Scan(
			&t.ID, 
			&t.Name,
			&t.Symbol,
			&t.Address,
			&t.Decimals,
			&t.TotalSupply)
		checkErr(err)

		tokens = append(tokens, &t)
	}

	return tokens
}