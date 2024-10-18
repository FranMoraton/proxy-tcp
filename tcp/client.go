package tcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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

// Establecer conexión TCP y enviar un mensaje, con soporte para cancelación de contexto
func SendTCPMessage(ctx context.Context, content string, tcpAddress string) (string, error) {
	var conn net.Conn
	var err error

	// Intentar conectarse al servidor TCP usando el contexto
	dialer := net.Dialer{
		Timeout: 5 * time.Second, // Timeout para la conexión
	}

	for {
		conn, err = dialer.DialContext(ctx, "tcp", tcpAddress)
		if err != nil {
			fmt.Println("Error al conectar al servidor TCP, reintentando en 2 segundos...")

			select {
			case <-ctx.Done(): // Revisar si el contexto fue cancelado
				return "Conexión cancelada por contexto", ctx.Err()
			case <-time.After(2 * time.Second): // Reintentar después de 2 segundos
				continue
			}
		}
		break // Conexión establecida, salir del bucle
	}
	defer conn.Close()

	// Enviar el mensaje al servidor TCP
	_, err = conn.Write([]byte(content))
	if err != nil {
		return "Error al enviar el mensaje TCP", err
	}

	// Leer la respuesta del servidor TCP con timeout
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Timeout de lectura
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "Error al recibir la respuesta TCP", err
	}

	return string(buf[:n]), nil
}
