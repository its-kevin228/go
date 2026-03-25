package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB


func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "books.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			pages INTEGER NOT NULL CHECK(pages > 0)
		)
	`)
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = db.Exec(`
			INSERT INTO books (title, author, pages) VALUES
			('1984', 'George Orwell', 328),
			('Le Petit Prince', 'Antoine de Saint-Exupery', 96),
			('Clean Code', 'Robert C. Martin', 464)
		`)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAllBooks() ([]Book, error) {
	rows, err := db.Query("SELECT id, title, author, pages FROM books ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Pages)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func getBookByID(id int) (Book, bool, error) {
	var book Book
	err := db.QueryRow("SELECT id, title, author, pages FROM books WHERE id = ?", id).Scan(&book.ID, &book.Title, &book.Author, &book.Pages)
	if err == sql.ErrNoRows {
		return Book{}, false, nil
	}
	if err != nil {
		return Book{}, false, err
	}
	return book, true, nil
}

func createBook(payload createBookRequest) (Book, error) {
	result, err := db.Exec("INSERT INTO books (title, author, pages) VALUES (?, ?, ?)", payload.Title, payload.Author, payload.Pages)
	if err != nil {
		return Book{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Book{}, err
	}

	return Book{ID: int(id), Title: payload.Title, Author: payload.Author, Pages: payload.Pages}, nil
}

func updateBook(id int, payload createBookRequest) (Book, bool, error) {
	result, err := db.Exec("UPDATE books SET title = ?, author = ?, pages = ? WHERE id = ?", payload.Title, payload.Author, payload.Pages, id)
	if err != nil {
		return Book{}, false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Book{}, false, err
	}
	if affected == 0 {
		return Book{}, false, nil
	}

	return Book{ID: id, Title: payload.Title, Author: payload.Author, Pages: payload.Pages}, true, nil
}

func deleteBook(id int) (bool, error) {
	result, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}
