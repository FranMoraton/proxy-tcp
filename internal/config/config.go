package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	WorkerCount int
	APIURL      string
	TCPAddress  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	workerCountStr := os.Getenv("WORKERS")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil || workerCount <= 0 {
		return nil, fmt.Errorf("error WORKERS missing in .env file")
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("error API_URL missing in .env file")
	}

	tcpAddress := os.Getenv("TCP_ADDRESS")
	if tcpAddress == "" {
		return nil, fmt.Errorf("error TCP_ADDRESS missing in .env file")
	}

	return &Config{
		WorkerCount: workerCount,
		APIURL:      apiURL,
		TCPAddress:  tcpAddress,
	}, nil
}
