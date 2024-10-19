package workers_test

import (
	"context"
	"myapp/workers"
	"sync"
	"testing"
	"time"
)

func TestStartWorkers_ProcessMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgChan := make(chan workers.Message, 1)
	var wg sync.WaitGroup

	// Iniciar 2 workers
	workers.StartWorkers(ctx, 2, msgChan, "http://example.com", "localhost:9001", &wg)

	// Enviar un mensaje al canal
	msg := workers.Message{
		Content:  "Mensaje de prueba",
		Response: make(chan string),
	}
	msgChan <- msg

	// Esperar la respuesta del worker
	select {
	case response := <-msg.Response:
		// Simulamos una respuesta en este caso
		if response != "Error al conectar al servidor TCP" {
			t.Errorf("Respuesta incorrecta: got %v want %v", response, "Error al conectar al servidor TCP")
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout esperando la respuesta del worker")
	}
}

func TestStartWorkers_CancelContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	msgChan := make(chan workers.Message)
	var wg sync.WaitGroup

	// Iniciar 1 worker
	workers.StartWorkers(ctx, 1, msgChan, "http://example.com", "localhost:9001", &wg)

	// Cancelar el contexto
	cancel()

	// Verificar que el worker termine correctamente
	wg.Wait()
}
