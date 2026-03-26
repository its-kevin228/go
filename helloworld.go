package main

import (
	"fmt"
	"net/http"
)

func main() {
	err := initDB()
	if err != nil {
		fmt.Println("Erreur initialisation SQLite:", err)
		return
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/status", statusHandler)
	mux.HandleFunc("/sum", sumHandler)
	mux.HandleFunc("/api/status", apiStatusHandler)
	mux.Handle("/api/books", apiKeyAuthMiddleware(http.HandlerFunc(apiBooksHandler)))

	handler := loggingMiddleware(mux)
	fmt.Println("Serveur lancé sur http://localhost:8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println("Erreur serveur:", err)
	}
}
