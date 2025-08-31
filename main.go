package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Checking if it is deployed!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	http.ListenAndServe(":"+port, nil)
}
