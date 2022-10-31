package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	HttpPort     string `mapstructure:"HTTP_PORT"`
	RpcUrl       string `mapstructure:"RPC_URL"`
	RpcApiToken  string `mapstructure:"RPC_API_TOKEN"`
	EsHost       string `mapstructure:"ES_HOST"`
	DataFilePath string `mapstructure:"DATA_FILE_PATH"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
