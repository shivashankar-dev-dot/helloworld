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
	todos  = make(map[int]Todo)
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
		var newItem Todo
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		} else {
			fmt.Printf("POST request received %v", err)
		}
		newItem.ID = nextID
		todos[nextID] = newItem
		nextID++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newItem)

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
