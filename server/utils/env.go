package utils

import (
	"os"
	"strconv"
)

func GetEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func GetEnvOrDefaultNumber(key string, defaultValue int64) int64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	n, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return defaultValue
	}
	return n
}
