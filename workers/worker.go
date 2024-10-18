package workers

import (
	"context"
	"fmt"
	"myapp/tcp"
)

// Estructura para almacenar el mensaje recibido de la conexión TCP
type Message struct {
	Content  string
	Response chan string // Canal para enviar la respuesta de vuelta al handler HTTP
}

// Iniciar los workers
func StartWorkers(ctx context.Context, workerCount int, msgChan chan Message, apiURL string, tcpAddress string) {
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d finalizado.\n", id)
					return // Finalizar el worker cuando el contexto se cancele
				case msg := <-msgChan:
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
