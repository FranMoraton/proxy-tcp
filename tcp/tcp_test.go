package tcp_test

import (
	"fmt"
	"myapp/tcp"
	"net"
	"testing"
)

func TestSendTCPMessage_Success(t *testing.T) {
	// Iniciar un servidor TCP simulado
	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		t.Fatalf("No se pudo iniciar el servidor TCP simulado: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatalf("No se pudo aceptar la conexión: %v", err)
		}
		defer conn.Close()

		// Leer el mensaje enviado por el cliente
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatalf("Error al leer del cliente: %v", err)
		}

		// Enviar una respuesta al cliente
		conn.Write([]byte(fmt.Sprintf("Recibido: %s", buf[:n])))
	}()

	// Enviar un mensaje al servidor TCP simulado
	response, err := tcp.SendTCPMessage("Hola TCP", "localhost:9001")
	if err != nil {
		t.Fatalf("Error al enviar el mensaje TCP: %v", err)
	}

	// Verificar que la respuesta sea la esperada
	expected := "Recibido: Hola TCP"
	if response != expected {
		t.Errorf("Respuesta incorrecta: got %v want %v", response, expected)
	}
}

func TestSendTCPMessage_Timeout(t *testing.T) {
	// Intentar conectar a un servidor inexistente
	_, err := tcp.SendTCPMessage("Hola TCP", "localhost:9999")

	// Verificar que se produce un error de timeout
	if err == nil {
		t.Error("Se esperaba un error de conexión, pero no se produjo")
	}
}
