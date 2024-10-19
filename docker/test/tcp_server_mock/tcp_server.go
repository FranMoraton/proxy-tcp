package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Definir el puerto TCP
	port := ":9000"

	// Escuchar en el puerto especificado
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Servidor TCP escuchando en", port)

	for {
		// Aceptar conexiones entrantes
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:", err)
			continue
		}

		// Manejar la conexión en una goroutine
		go handleConnection(conn)
	}
}

// Función para manejar conexiones
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Leer el mensaje del cliente
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error al leer del cliente:", err)
		return
	}

	// Mostrar el mensaje recibido
	message := string(buffer[:n])
	fmt.Println("Mensaje recibido:", message)

	// Responder al cliente
	response := "Mensaje recibido: " + message
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error al enviar respuesta:", err)
	}
}
