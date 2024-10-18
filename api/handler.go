package api

import (
	"encoding/json"
	"myapp/workers"
	"net/http"
)

// RequestBody defines the expected structure of the incoming JSON payload
type RequestBody struct {
	Message string `json:"message"`
}

// HandleSend maneja las solicitudes entrantes en la ruta /send
func HandleSend(w http.ResponseWriter, r *http.Request, msgChan chan workers.Message) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el Content-Type sea application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type no soportado, debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Decodificar el cuerpo de la petición JSON
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Leer el mensaje del cuerpo de la petición
	msg := workers.Message{
		Content:  reqBody.Message,
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
