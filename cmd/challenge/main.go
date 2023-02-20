package main

import (
	"flag"
	"go.uber.org/zap"
	"database/sql"
	"net/http"
	"encoding/json"
	"log"
)

type TokenModel struct {
	Name               string  `json:"name"`
	Symbol             string  `json:"symbol"`
	Address            string  `json:"address"`
	Decimals           uint64  `json:"decimals"`
	TotalSupply        string  `json:"total_supply"`
}

var db *sql.DB  // global variable to make db available to api endpoints

var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")
var flagAPIKEY = flag.String("api-key", "", "alchemy api key")

const (
	host     = "localhost"
	port     = 6000  // Local Port to use with Dockerized Postgres server
	user     = "postgres"
	password = "secret"
	dbname   = "postgres"
  )

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL
	apiKey := *flagAPIKEY

	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)

	db = DBConnect()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	zlog.Info("Successfully connected to DB")

	// apiKey := "RtBNZI7jboJBSVutqQidtcUE8Nbw2M6p"
	go seedDB(rpcURL, apiKey, db)

	http.HandleFunc("/tokens", tokenEndpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func tokenEndpoint(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	var tokens []*TokenModel
	if q != "" {
		tokens = queryToken(q, db)
	}

	json.NewEncoder(w).Encode(tokens)
}
