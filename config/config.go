package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config estructura para almacenar todas las configuraciones
type Config struct {
	WorkerCount int
	APIURL      string
	TCPAddress  string
}

// LoadConfig carga las variables de entorno y devuelve una configuración
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	// Leer y validar WORKERS
	workerCountStr := os.Getenv("WORKERS")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil || workerCount <= 0 {
		workerCount = 5 // Valor por defecto
	}

	// Leer API_URL
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8000/another-api" // Valor por defecto
	}

	// Leer TCP_ADDRESS
	tcpAddress := os.Getenv("TCP_ADDRESS")
	if tcpAddress == "" {
		tcpAddress = "localhost:9000" // Valor por defecto
	}

	return &Config{
		WorkerCount: workerCount,
		APIURL:      apiURL,
		TCPAddress:  tcpAddress,
	}, nil
}
