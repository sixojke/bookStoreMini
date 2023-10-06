package main

var bookStore = BookStore{
	books: make([]Book, 0),
}

type Resp struct {
	Message interface{}
	Error   string
}

type Book struct {
	Id     string `json:"id"`
	Author string `json:"author"`
	Name   string `json:"name"`
}

type BookStore struct {
	books []Book
}
