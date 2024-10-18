package main

import (
	"context"
	"fmt"
	"myapp/api"
	"myapp/workers"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func main() {
	// Crear contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Leer el número de workers desde la variable de entorno
	workerCountStr := os.Getenv("WORKERS")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil || workerCount <= 0 {
		workerCount = 5 // Valor por defecto si no se configura correctamente
	}

	// Leer la URL de la otra API desde la variable de entorno
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8000/another-api" // Valor por defecto
	}

	// Leer la dirección TCP desde las variables de entorno
	tcpAddress := os.Getenv("TCP_ADDRESS")
	if tcpAddress == "" {
		tcpAddress = "localhost:9000" // Valor por defecto
	}

	// Crear un canal para comunicarse entre goroutines
	msgChan := make(chan workers.Message)

	// Crear un wait group para esperar a que los workers terminen
	var wg sync.WaitGroup

	// Iniciar los workers con el contexto
	workers.StartWorkers(ctx, workerCount, msgChan, apiURL, tcpAddress, &wg)

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		api.HandleSend(w, r, msgChan)
	})

	// Manejar la terminación del programa
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs // Esperar la señal de terminación
		fmt.Println("Señal de terminación recibida. Cerrando el servidor...")

		// Cancelar todas las operaciones con el contexto
		cancel()

		// Cerrar el canal de mensajes, permitiendo que los workers terminen
		close(msgChan)

		// Esperar a que todos los workers finalicen
		wg.Wait()
		fmt.Println("Todos los workers han terminado. Saliendo...")
	}()

	fmt.Println("Servidor HTTP escuchando en :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
