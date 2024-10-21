package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Respuesta desde la API dummy!",
	}

	w.Header().Set("Content-Type", "application/json")

	respBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return
	}

	fmt.Printf("mensaje: %s\n", string(respBody))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/dummy-api", handler)

	port := ":8081"
	log.Printf("Servidor Dummy escuchando en %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %s", err)
	}
}
