package workers

import (
	"context"
	"fmt"
	"myapp/internal/tcp"
	"sync"
)

type Message struct {
	Content  string
	Response chan string
}

func StartWorkers(ctx context.Context, workerCount int, msgChan chan Message, apiURL, tcpAddress string, wg *sync.WaitGroup) {
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case msg, ok := <-msgChan:
					if !ok { // Canal cerrado
						fmt.Printf("Worker %d finalizado por cierre de canal.\n", workerID)
						return
					}

					response, err := tcp.SendTCPMessage(msg.Content, tcpAddress)
					if err != nil {
						msg.Response <- response
						continue
					}
					msg.Response <- response

					tcp.CallAnotherAPI(apiURL, response)
					fmt.Printf("Mensaje procesado por worker %d\n", workerID)

				case <-ctx.Done(): // Cancelación recibida
					fmt.Printf("Worker %d finalizado por contexto cancelado.\n", workerID)
					return
				}
			}
		}(i)
	}
}
