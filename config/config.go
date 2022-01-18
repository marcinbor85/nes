package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Alternate(priorityVal string, key string, defaultVal string) string {
	if priorityVal != "" {
		return priorityVal
	}
	configVal := Get(key)
	if configVal != "" {
		return configVal
	}
	return defaultVal
}

func Init(fileName string) {
	_ = godotenv.Load(fileName)
}

func Get(key string) string {
	return os.Getenv(key)
}
