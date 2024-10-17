package api

import (
	"net/http"

	"myapp/workers"
)

// HandleSend maneja las solicitudes entrantes en la ruta /send
func HandleSend(w http.ResponseWriter, r *http.Request, msgChan chan workers.Message) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el mensaje del cuerpo de la petición
	msg := workers.Message{Content: r.FormValue("message"), Response: make(chan string)}

	// Enviar el mensaje al canal
	msgChan <- msg

	// Enviar la respuesta "OK" inmediatamente
	w.Write([]byte("OK"))

	// Esperar la respuesta asincrónica
	go func() {
		response := <-msg.Response
		// Aquí podrías manejar la respuesta, como imprimirla o registrarla
		println("Respuesta del TCP:", response)
	}()
}
