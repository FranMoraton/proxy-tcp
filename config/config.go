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
		workerCount = 5
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8000/another-api"
	}

	tcpAddress := os.Getenv("TCP_ADDRESS")
	if tcpAddress == "" {
		tcpAddress = "localhost:9000"
	}

	return &Config{
		WorkerCount: workerCount,
		APIURL:      apiURL,
		TCPAddress:  tcpAddress,
	}, nil
}
