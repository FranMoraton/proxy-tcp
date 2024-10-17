package tcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

// Llamar a otra API HTTP con la respuesta recibida
func CallAnotherAPI(apiURL, response string) {
	body, err := json.Marshal(map[string]string{"response": response})
	if err != nil {
		fmt.Println("Error al serializar el cuerpo:", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error al llamar a la otra API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("La otra API devolvió un estado %s\n", resp.Status)
	} else {
		fmt.Println("Llamada a la otra API completada con éxito.")
	}
}

// Establecer conexión TCP y enviar un mensaje
func SendTCPMessage(content string) (string, error) {
	var conn net.Conn
	var err error

	// Leer la dirección TCP desde las variables de entorno
	tcpAddress := os.Getenv("TCP_ADDRESS")
	if tcpAddress == "" {
		return "Error: TCP_ADDRESS no está configurada", fmt.Errorf("TCP_ADDRESS no está configurada")
	}

	// Intentar conectarse al servidor TCP usando la variable de entorno
	for {
		conn, err = net.Dial("tcp", tcpAddress)
		if err != nil {
			fmt.Println("Error al conectar al servidor TCP, reintentando en 2 segundos...")
			time.Sleep(2 * time.Second) // Esperar antes de reintentar
			continue
		}
		break // Conexión establecida, salir del bucle
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
