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

// CallAnotherAPI realiza una llamada POST a otra API
func CallAnotherAPI(apiURL, response string) {
	body, err := json.Marshal(map[string]string{"response": response})
	if err != nil {
		fmt.Println("Error al serializar el cuerpo:", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error al llamar a la API:", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("La API devolvió un estado %s. Respuesta: %s\n", resp.Status, string(respBody))
	} else {
		fmt.Printf("Llamada a la API completada con éxito. Respuesta: %s\n", string(respBody))
	}
}

// SendTCPMessage envía un mensaje TCP y espera una respuesta
func SendTCPMessage(content, tcpAddress string) (string, error) {
	dialer := net.Dialer{Timeout: 40 * time.Second}
	conn, err := dialer.Dial("tcp", tcpAddress)
	if err != nil {
		return "Error al conectar al servidor TCP", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(content))
	if err != nil {
		return "Error al enviar el mensaje TCP", err
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "Error al recibir la respuesta TCP", err
	}

	return string(buf[:n]), nil
}
