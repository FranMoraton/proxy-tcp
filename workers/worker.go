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

// Iniciar los workers con el contexto
func StartWorkers(ctx context.Context, workerCount int, msgChan chan Message, apiURL string, tcpAddress string) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done(): // Si el contexto se cancela, salir del worker
					fmt.Println("Worker finalizado.")
					return
				case msg, ok := <-msgChan:
					if !ok {
						return // Canal cerrado, finalizar worker
					}
					response, err := tcp.SendTCPMessage(ctx, msg.Content, tcpAddress)
					if err != nil {
						msg.Response <- response
						continue
					}
					msg.Response <- response

					// Llamar a otra API HTTP con la respuesta recibida
					tcp.CallAnotherAPI(apiURL, response)
				}
			}
		}()
	}
}
