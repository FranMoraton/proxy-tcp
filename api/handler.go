package api

import (
	"myapp/workers"
	"net/http"
)

// HandleSend maneja las solicitudes entrantes en la ruta /send
func HandleSend(w http.ResponseWriter, r *http.Request, msgChan chan workers.Message) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el mensaje del cuerpo de la petición
	msg := workers.Message{
		Content:  r.FormValue("message"),
		Response: make(chan string),
	}

	// Verificar si el canal está cerrado antes de enviar
	select {
	case msgChan <- msg:
		// Si el mensaje fue enviado exitosamente
		w.Write([]byte("OK"))

		// Esperar la respuesta asincrónica en una goroutine separada
		go func() {
			response := <-msg.Response
			// Aquí podrías manejar la respuesta, como imprimirla o registrarla
			println("Respuesta del TCP:", response)
		}()
	default:
		// Si el canal está cerrado, devolver un error al cliente
		http.Error(w, "No se puede enviar el mensaje, el servidor está cerrando", http.StatusInternalServerError)
	}
}
