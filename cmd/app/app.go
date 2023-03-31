package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/handler"
	"github.com/rarecircles/backend-challenge-go/cmd/helper"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/logging"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DAO    dao.DaoInterface
}

var app *App

func NewApp() *App {

	db := createDbConnection()
	runMigration(db)

	return &App{
		Router: mux.NewRouter().StrictSlash(true),
		DAO: &dao.Dao{
			DB: db,
		},
	}
}

func runMigration(db *gorm.DB) {
	err := db.AutoMigrate(&model.Token{})
	if err != nil {
		errorMsg := fmt.Sprintf("could not run the migration %v", err)
		panic(errorMsg)
	}
}

func createDbConnection() *gorm.DB {
	dsn := createDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		errorMsg := fmt.Sprintf("could not connect to the database %v", err)
		panic(errorMsg)
	}
	return db
}

func createDSN() string {

	host := helper.GetEnv("DB_HOST", "localhost")
	user := helper.GetEnv("DB_USER", "shamsazad")
	password := helper.GetEnv("DB_PASSWORD", "123")
	dbname := helper.GetEnv("DB_NAME", "rare_circle")
	port := helper.GetEnv("DB_PORT", "5432")

	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
}

func (a *App) HandleRequests() {
	if app == nil {
		app = NewApp()
	}
	setup()
	log.Println("inside app")
	app.Router.HandleFunc("/tokens", handler.GetTokens(app.DAO)).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":10000", app.Router))
}

// This method read the address, extract tokens from address and load it to DB
func setup() {
	zLog := logging.MustCreateLoggerWithServiceName("challenge")

	rpcURL := helper.GetEnv("RPC_CURL", "https://eth-mainnet.alchemyapi.io/v2/")
	rpcTOKEN := helper.GetEnv("RPC_TOKEN", "sscKIY7xv-5uHTJtHVRdmpS-y-Q3rBVA")
	addChannel := make(chan string, 10)
	tokenChannel := make(chan model.TokenDTO, 10)

	addLoader := helper.NewAddLoader(addChannel)
	rpcClient := rpc.NewClient(rpcURL + rpcTOKEN)

	// scan and load address
	go func() {
		scanAddress(addLoader, zLog)
		defer close(addChannel)
	}()
	// Extract tokens from address
	go func() {
		helper.ParseTokenData(rpcClient, addChannel, tokenChannel, zLog)
		defer close(tokenChannel)
	}()

	go func() {
		helper.SeedDB(app.DAO, tokenChannel, zLog)
	}()
}

func scanAddress(loader helper.IAddLoader, zLog *zap.Logger) {
	err := loader.ScanAndLoadAddressFile()
	if err != nil {
		zLog.Fatal(err.Error())
	}
}
