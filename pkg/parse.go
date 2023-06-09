package pkg

import (
	"os"
	"strconv"
	"time"
)

func getBoolEnv(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	durationValue, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return durationValue
}

func getIntEnv(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
