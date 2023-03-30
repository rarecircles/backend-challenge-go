package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/helper"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/logging"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
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

	host := getEnv("DB_HOST", "database")
	user := getEnv("DB_USER", "test_rare_circle_user")
	password := getEnv("DB_PASSWORD", "123")
	dbname := getEnv("DB_NAME", "rare_circle")
	port := getEnv("DB_PORT", "5432")

	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
}

func getEnv(key string, defaultValue string) string {
	configValue := os.Getenv(key)
	if configValue == "" {
		return defaultValue
	} else {
		return configValue
	}
}

func (a *App) HandleRequests() {
	if app == nil {
		app = NewApp()
	}
	setup()
	log.Println("inside app")
}

// This method read the address, extract tokens from address and load it to DB
func setup() {
	zLog := logging.MustCreateLoggerWithServiceName("challenge")

	rpcURL := getEnv("RPC_CURL", "")
	rpcTOKEN := getEnv("RPC_TOKEN", "")
	addChannel := make(chan string, 10)
	tokenChannel := make(chan model.TokenDTO, 10)

	addLoader := helper.NewAddLoader(addChannel)
	rpcClient := rpc.NewClient(rpcURL + rpcTOKEN)

	var wg sync.WaitGroup
	wg.Add(2)
	// scan and load address
	go func() {
		defer wg.Done()
		scanAddress(addLoader, zLog)
	}()
	// Extract tokens from address
	go func() {
		defer wg.Done()
		getTokenData(rpcClient, addChannel, tokenChannel, zLog)
	}()

	//searchEngine, err := search_engine.NewElasticsearchIngest(tokenDataChannel, zLog)
	//if err != nil {
	//	zLog.Fatal(err.Error())
	//}

	wg.Wait()
}

func getTokenData(client *rpc.Client, addChannel chan string, tokenChannel chan model.TokenDTO, zLog *zap.Logger) {
	for add := range addChannel {
		addr, err := eth.NewAddress(add)
		if err != nil {
			zLog.Error("eth address doesn't get created " + err.Error())
		}

		ethToken, err := client.GetERC20(addr)
		if err != nil {
			zLog.Error("unable to fetch ERC20: " + err.Error())
		}

		tokenChannel <- model.TokenDTO{
			Name:        ethToken.Name,
			Symbol:      ethToken.Symbol,
			Address:     ethToken.Address.String(),
			Decimals:    ethToken.Decimals,
			TotalSupply: ethToken.TotalSupply,
		}
	}
}

func scanAddress(loader helper.IAddLoader, zLog *zap.Logger) {
	err := loader.ScanAndLoadAddressFile()
	if err != nil {
		zLog.Fatal(err.Error())
	}
}
