package main

import (
	"context"
	"fmt"
	"myapp/internal/api"
	"myapp/internal/config"
	"myapp/internal/workers"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error al cargar la configuración:", err)
		return
	}

	// Crear contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Crear canal de mensajes
	msgChan := make(chan workers.Message)

	// Crear un wait group para los workers
	var wg sync.WaitGroup

	// Iniciar los workers
	workers.StartWorkers(ctx, cfg.WorkerCount, msgChan, cfg.APIURL, cfg.TCPAddress, &wg)

	// Crear y configurar el servidor HTTP
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		api.HandleSend(w, r, msgChan)
	})

	// Manejo de señales del sistema para una finalización ordenada
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Señal de terminación recibida. Cerrando el servidor...")

		cancel()       // Cancelar contexto para los workers
		close(msgChan) // Cerrar el canal de mensajes
		wg.Wait()      // Esperar a que los workers finalicen

		// Detener el servidor HTTP de manera ordenada
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			fmt.Println("Error al cerrar el servidor:", err)
		}
	}()

	fmt.Println("Servidor HTTP escuchando en :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
