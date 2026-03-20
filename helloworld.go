package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

var books = []Book{
	{ID: 1, Title: "1984", Author: "George Orwell", Pages: 328},
	{ID: 2, Title: "Le Petit Prince", Author: "Antoine de Saint-Exupery", Pages: 96},
	{ID: 3, Title: "Clean Code", Author: "Robert C. Martin", Pages: 464},
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur mon premier serveur Go HTTP !")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Kevin"
	}

	fmt.Fprintln(w, "Bonjour", name, "👋")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Je suis un développeur Go passionné par la création de serveurs web simples et efficaces.")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Le serveur est en ligne et fonctionne correctement.")
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Query().Get("a")
	b := r.URL.Query().Get("b")

	if a == "" || b == "" {
		fmt.Fprintln(w, "Veuillez fournir les paramètres 'a' et 'b' pour effectuer la somme.")
		return
	}

	aNum, err := strconv.Atoi(a)
	if err != nil {
		fmt.Fprintln(w, "Le paramètre 'a' doit être un nombre.")
		return
	}

	bNum, err := strconv.Atoi(b)
	if err != nil {
		fmt.Fprintln(w, "Le paramètre 'b' doit être un nombre.")
		return
	}

	result := aNum + bNum
	fmt.Fprintln(w, "Le résultat de la somme de", aNum, "et", bNum, "est:", result)

}

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	reponse := map[string]string{
		"status":  "OK",
		"message": "Le serveur est en ligne et fonctionne correctement.",
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(reponse)
	if err != nil {
		http.Error(w, "Erreur lors de l'encodage de la réponse JSON", http.StatusInternalServerError)
		return
	}
}

func apiBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/api/status", apiStatusHandler)
	http.HandleFunc("/api/books", apiBooksHandler)
	fmt.Println("Serveur lancé sur http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur serveur:", err)
	}
}
