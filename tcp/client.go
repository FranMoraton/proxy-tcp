package tcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// Llamar a otra API HTTP con la respuesta recibida
func CallAnotherAPI(apiURL, response string) {
	// Serializar el cuerpo de la petición
	body, err := json.Marshal(map[string]string{"response": response})
	if err != nil {
		fmt.Println("Error al serializar el cuerpo:", err)
		return
	}

	// Hacer la petición POST a la otra API
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error al llamar a la otra API:", err)
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta para manejarlo o mostrarlo
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return
	}

	// Comprobar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("La otra API devolvió un estado %s. Respuesta: %s\n", resp.Status, string(respBody))
	} else {
		fmt.Printf("Llamada a la otra API completada con éxito. Respuesta: %s\n", string(respBody))
	}
}

// Establecer conexión TCP y enviar un mensaje con contexto y timeout
func SendTCPMessage(content string, tcpAddress string) (string, error) {
	var conn net.Conn
	var err error

	// Establecer un timeout para la conexión
	dialer := net.Dialer{Timeout: 40 * time.Second}

	// Intentar conectarse al servidor TCP
	conn, err = dialer.Dial("tcp", tcpAddress)
	if err != nil {
		return "Error al conectar al servidor TCP", err
	}
	defer conn.Close()

	// Enviar el mensaje al servidor TCP
	_, err = conn.Write([]byte(content))
	if err != nil {
		return "Error al enviar el mensaje TCP", err
	}

	// Leer la respuesta del servidor TCP
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "Error al recibir la respuesta TCP", err
	}

	return string(buf[:n]), nil
}
