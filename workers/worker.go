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
		wg.Add(1) // Agregar al wait group
		go func(workerID int) {
			defer wg.Done() // Marcar el worker como terminado al final

			for {
				select {
				case msg, ok := <-msgChan:
					if !ok { // Canal cerrado, salir
						fmt.Printf("cerrado por !ok en canal %d \n", workerID)
						return
					}
					response, err := tcp.SendTCPMessage(msg.Content, tcpAddress)
					if err != nil {
						msg.Response <- response
						continue
					}
					msg.Response <- response

					// Llamar a otra API HTTP con la respuesta recibida
					tcp.CallAnotherAPI(apiURL, response)
					fmt.Println("gestionado por ", workerID)
				case <-ctx.Done(): // Manejar la cancelación
					fmt.Printf("Worker %d finalizado.\n", workerID)
					return // Salir del worker
				}
			}
		}(i) // Pasar el ID del worker a la goroutine
	}
}
