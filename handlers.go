package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			books, err := getAllBooks()
			if err != nil {
				http.Error(w, "Erreur base de données", http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(books)
			if err != nil {
				http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
			}
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Le paramètre 'id' doit être un nombre", http.StatusBadRequest)
			return
		}

		book, found, err := getBookByID(id)
		if err != nil {
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, "Livre introuvable", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(book)
		if err != nil {
			http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		}

	case http.MethodPost:
		var payload createBookRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}

		if payload.Title == "" || payload.Author == "" || payload.Pages <= 0 {
			http.Error(w, "title, author et pages sont obligatoires", http.StatusBadRequest)
			return
		}

		newBook, err := createBook(payload)
		if err != nil {
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(newBook)
		if err != nil {
			http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Le paramètre 'id' est obligatoire", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Le paramètre 'id' doit être un nombre", http.StatusBadRequest)
			return
		}

		deleted, err := deleteBook(id)
		if err != nil {
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}
		if !deleted {
			http.Error(w, "Livre introuvable", http.StatusNotFound)
			return
		}

		response := map[string]string{
			"status":  "ok",
			"message": "Livre supprimé",
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		}

	case http.MethodPut:
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Le paramètre 'id' est obligatoire", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Le paramètre 'id' doit être un nombre", http.StatusBadRequest)
			return
		}

		var payload createBookRequest
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}

		if payload.Title == "" || payload.Author == "" || payload.Pages <= 0 {
			http.Error(w, "title, author et pages sont obligatoires", http.StatusBadRequest)
			return
		}

		updatedBook, found, err := updateBook(id, payload)
		if err != nil {
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, "Livre introuvable", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(updatedBook)
		if err != nil {
			http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		}

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
