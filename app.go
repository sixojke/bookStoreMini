package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	handler := http.NewServeMux()

	handler.HandleFunc("/hello/", BasicAuth(Logger(helloHandler)))

	handler.HandleFunc("/book/", Logger(bookHandler))

	handler.HandleFunc("/books/", Logger(booksHandler))

	s := http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 * 2 ^ 20 - 128kByte
	}

	log.Fatal(s.ListenAndServe())
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	name := strings.Replace(r.URL.Path, "/hello/", "", 1)
	resp := Resp{
		Message: fmt.Sprintf("hello, %s. Glad to see you again!", name),
	}

	respJson, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGetBook(w, r)
	}

	if r.Method == http.MethodPost {
		handleAddBook(w, r)
	}

	if r.Method == http.MethodPut {
		handleUpdateBook(w, r)
	}

	if r.Method == http.MethodDelete {
		handleDeleteBook(w, r)
	}
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.URL.Path, "/book/", "", 1)

	decoder := json.NewDecoder(r.Body)

	var book Book

	var resp Resp

	err := decoder.Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)

		w.Write(respJson)

		return
	}

	book.Id = id

	err = bookStore.UpdateBook(book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)

		w.Write(respJson)

		return
	}

	resp.Message = book

	respJson, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)

	w.Write(respJson)
}

func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.URL.Path, "/book/", "", 1)

	var resp Resp

	err := bookStore.DeleteBook(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)

		w.Write(respJson)

		return
	}

	booksHandler(w, r)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGetBook(w, r)
	}

	w.WriteHeader(http.StatusOK)

	resp := Resp{
		Message: bookStore.GetBooks(),
	}

	booksJson, _ := json.Marshal(resp)

	w.Write(booksJson)
}

func handleAddBook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var book Book

	var resp Resp

	err := decoder.Decode(&book)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)

		w.Write(respJson)

		return
	}
	err = bookStore.AddBook(book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)

		w.Write(respJson)

		return
	}

	booksHandler(w, r)

}

func handleGetBook(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.URL.Path, "/book/", "", 1)

	book := bookStore.FindBookById(id)

	var resp Resp

	if book == nil {
		w.WriteHeader(http.StatusNotFound)

		resp.Error = fmt.Sprintf("Book with id %s not found", id)

		respJson, _ := json.Marshal(resp)

		w.Write(respJson)

		return
	}

	resp.Message = book

	respJson, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)

	w.Write(respJson)
}

func (s BookStore) FindBookById(id string) *Book {
	for _, book := range s.books {
		if book.Id == id {
			return &book
		}
	}
	return nil
}

func (s BookStore) GetBooks() []Book {
	return s.books
}

func (s *BookStore) AddBook(book Book) error {
	for _, bk := range s.books {
		if bk.Id == book.Id {
			return errors.New(fmt.Sprintf("Book with id %s not found", book.Id))
		}
	}

	s.books = append(s.books, book)
	return nil
}

func (s *BookStore) UpdateBook(book Book) error {
	for i, bk := range s.books {
		if bk.Id == book.Id {
			s.books[i] = book
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Book with id %s not found", book.Id))
}

func (s *BookStore) DeleteBook(id string) error {
	for i, bk := range s.books {
		if bk.Id == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Book with id %s not found", id))
}
