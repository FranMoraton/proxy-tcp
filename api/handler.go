package api

import (
	"encoding/json"
	"myapp/workers"
	"net/http"
)

// RequestBody estructura de la petición entrante
type RequestBody struct {
	Message string `json:"message"`
}

// HandleSend maneja la solicitud POST /send
func HandleSend(w http.ResponseWriter, r *http.Request, msgChan chan workers.Message) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el Content-Type sea application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type no soportado", http.StatusUnsupportedMediaType)
		return
	}

	// Decodificar el cuerpo JSON
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Crear mensaje y enviar al canal
	msg := workers.Message{
		Content:  reqBody.Message,
		Response: make(chan string),
	}

	select {
	case msgChan <- msg:
		// Mensaje enviado, devolver OK
		w.Write([]byte("OK"))

		// Esperar la respuesta del worker
		go func() {
			response := <-msg.Response
			println("Respuesta del TCP:", response)
		}()
	default:
		// Si el canal está cerrado
		http.Error(w, "No se puede enviar el mensaje, el servidor está cerrando", http.StatusInternalServerError)
	}
}
