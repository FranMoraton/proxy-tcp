package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	port := ":9000"

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Servidor TCP escuchando en", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error al leer del cliente:", err)
		return
	}

	message := string(buffer[:n])
	fmt.Println("Mensaje recibido:", message)

	response := "Mensaje recibido: " + message
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error al enviar respuesta:", err)
	}
}
