package api_test

import (
	"bytes"
	"encoding/json"
	"myapp/api"
	"myapp/workers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleSend_Success(t *testing.T) {
	// Crear un canal ficticio para pasar mensajes
	msgChan := make(chan workers.Message, 1)

	// Crear un cuerpo de solicitud válido
	reqBody := map[string]string{
		"message": "Hello, World!",
	}
	body, _ := json.Marshal(reqBody)

	// Crear una solicitud HTTP POST con el JSON
	req, err := http.NewRequest("POST", "/send", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("No se pudo crear la solicitud: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Crear un ResponseRecorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Ejecutar la función HandleSend
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.HandleSend(w, r, msgChan)
	})

	handler.ServeHTTP(rr, req)

	// Verificar que la respuesta tenga el código de estado correcto
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Código de estado incorrecto: got %v want %v", status, http.StatusOK)
	}

	// Verificar que el mensaje haya sido enviado al canal
	select {
	case msg := <-msgChan:
		if msg.Content != "Hello, World!" {
			t.Errorf("El mensaje en el canal es incorrecto: got %v want %v", msg.Content, "Hello, World!")
		}
	default:
		t.Error("No se recibió ningún mensaje en el canal")
	}
}

func TestHandleSend_InvalidMethod(t *testing.T) {
	// Crear una solicitud HTTP GET en lugar de POST
	req, err := http.NewRequest("GET", "/send", nil)
	if err != nil {
		t.Fatalf("No se pudo crear la solicitud: %v", err)
	}

	// Crear un ResponseRecorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Ejecutar la función HandleSend
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.HandleSend(w, r, nil)
	})

	handler.ServeHTTP(rr, req)

	// Verificar que la respuesta tenga el código de estado correcto
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Código de estado incorrecto: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandleSend_InvalidContentType(t *testing.T) {
	// Crear una solicitud HTTP POST pero con un Content-Type incorrecto
	req, err := http.NewRequest("POST", "/send", nil)
	if err != nil {
		t.Fatalf("No se pudo crear la solicitud: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	// Crear un ResponseRecorder para capturar la respuesta
	rr := httptest.NewRecorder()

	// Ejecutar la función HandleSend
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.HandleSend(w, r, nil)
	})

	handler.ServeHTTP(rr, req)

	// Verificar que la respuesta tenga el código de estado correcto
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("Código de estado incorrecto: got %v want %v", status, http.StatusUnsupportedMediaType)
	}
}
