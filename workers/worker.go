package workers

import (
	"context"
	"fmt"
	"myapp/tcp"
	"sync"
)

// Estructura para almacenar el mensaje recibido de la conexión TCP
type Message struct {
	Content  string
	Response chan string // Canal para enviar la respuesta de vuelta al handler HTTP
}

// Iniciar los workers
func StartWorkers(ctx context.Context, workerCount int, msgChan chan Message, apiURL string, tcpAddress string, wg *sync.WaitGroup) {
	for i := 0; i < workerCount; i++ {
		wg.Add(1) // Incrementar el contador del wait group por cada worker

		go func(id int) {
			defer wg.Done() // Decrementar el contador del wait group cuando el worker finalice

			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d finalizado.\n", id)
					return // Finalizar el worker cuando el contexto se cancele
				case msg, ok := <-msgChan:
					if !ok {
						// Si el canal está cerrado, salir del worker
						fmt.Printf("Worker %d: canal cerrado, finalizando.\n", id)
						return
					}
					// Procesar el mensaje
					response, err := tcp.SendTCPMessage(msg.Content, tcpAddress)
					if err != nil {
						msg.Response <- "Error al enviar mensaje TCP"
						continue
					}
					msg.Response <- response

					// Llamar a otra API HTTP con la respuesta recibida
					tcp.CallAnotherAPI(apiURL, response)
				}
			}
		}(i)
	}
}
