package helper

import "os"

func GetEnv(key string, defaultValue string) string {
	configValue := os.Getenv(key)
	if configValue == "" {
		return defaultValue
	} else {
		return configValue
	}
}
