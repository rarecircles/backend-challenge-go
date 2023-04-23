package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	BaseUrl           string
	Port              string
	EthRpcUrl         string
	EthRpcKey         string
	RedisSearchUrl    string
	RedisKeyPrefix    string
	AddressesFile     string
	NumWorkers        int
	RefetchDelayHours int
}

func getConfig(name string) Config {
	setupViper(name)
	return getFilledConfig()
}

func setupViper(name string) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	viper.SetEnvPrefix(name)
	viper.AutomaticEnv()
	// NOTE: Configs are either stored in config file without prefix or in env variables with prefix "SERVER_*"
	// NOTE: Environment variables have higher priority than the config file
}

func getFilledConfig() Config {
	return Config{
		BaseUrl:           viper.GetString("BASE_URL"),
		Port:              viper.GetString("PORT"),
		EthRpcUrl:         viper.GetString("ETH_RPC_URL"),
		EthRpcKey:         viper.GetString("ETH_RPC_KEY"),
		RedisSearchUrl:    viper.GetString("REDIS_SEARCH_URL"),
		RedisKeyPrefix:    viper.GetString("REDIS_KEY_PREFIX"),
		AddressesFile:     viper.GetString("ADDRESSES_FILE"),
		NumWorkers:        viper.GetInt("NUM_WORKERS"),
		RefetchDelayHours: viper.GetInt("REFETCH_DELAY_HOURS"),
	}
}
