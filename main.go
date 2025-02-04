package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	users      []string
	usersMutex sync.Mutex
)

func main() {
	// Crear un ServeMux independiente para el servidor 1
	server1 := http.NewServeMux()
	server1.HandleFunc("/add", addUser)
	server1.HandleFunc("/users", getUsers)

	// Crear un ServeMux independiente para el servidor 2
	server2 := http.NewServeMux()
	server2.HandleFunc("/users", getUsers)

	// Iniciar el servidor 1 en una goroutine
	go func() {
		fmt.Println("Servidor 1 escuchando en :8080")
		http.ListenAndServe(":8080", server1)
	}()

	// Iniciar el servidor 2 en otra goroutine
	go func() {
		fmt.Println("Servidor 2 escuchando en :8081")
		http.ListenAndServe(":8081", server2)
	}()

	// Mantener el programa en ejecución
	select {}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var newUser string
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	users = append(users, newUser)
	usersMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Usuario agregado: %s\n", newUser)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}