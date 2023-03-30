package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
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
	err := db.AutoMigrate(&model.TokenDTO{})
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
	log.Println("inside app")
}
