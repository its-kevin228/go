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

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/api/status", apiStatusHandler)
	http.HandleFunc("/api/books", apiBooksHandler)
	fmt.Println("Serveur lancé sur http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur serveur:", err)
	}
}
