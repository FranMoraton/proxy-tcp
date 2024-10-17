package workers

import (
	"myapp/tcp"
)

// Estructura para almacenar el mensaje recibido de la conexión TCP
type Message struct {
	Content  string
	Response chan string // Canal para enviar la respuesta de vuelta al handler HTTP
}

// Iniciar los workers
func StartWorkers(workerCount int, msgChan chan Message, apiURL string) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for msg := range msgChan {
				response, err := tcp.SendTCPMessage(msg.Content)
				if err != nil {
					msg.Response <- response
					continue
				}
				msg.Response <- response

				// Llamar a otra API HTTP con la respuesta recibida
				tcp.CallAnotherAPI(apiURL, response)
			}
		}()
	}
}
