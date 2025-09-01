package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	todos  []Todo
	nextID = 1
	mu     sync.Mutex
)

func handleItems(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		var todoList []Todo

		for _, todo := range todos {
			todoList = append(todoList, todo)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todoList)
		fmt.Printf("GET request received %v", todoList)

	case http.MethodPost:
		var newTodo Todo

		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			http.Error(w, "Invalid fdsfdf payload", http.StatusBadRequest)
			return
		}

		newTodo.ID = nextID
		todos[nextID] = newTodo
		nextID++

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

func main() {
	http.HandleFunc("/", handleItems)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}
