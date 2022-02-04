package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var Books []Book = []Book{
	Book{
		Id:     1,
		Title:  "Moreninha",
		Author: "Joao",
	},
	Book{
		Id:     2,
		Title:  "Mao Petra",
		Author: "Macabeus",
	},
	Book{
		Id:     3,
		Title:  "Guarani",
		Author: "Alencar",
	},
}

func routeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome")
}

func routeListBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(Books)
}

func createBooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		log.Fatal(error)
	}
	var newBook Book
	json.Unmarshal(body, &newBook)
	newBook.Id = len(Books) + 1
	Books = append(Books, newBook)
	encoder := json.NewEncoder(w)
	encoder.Encode(newBook)
}

func alterBook(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id, erro := strconv.Atoi(parts[2])

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, bodyError := ioutil.ReadAll(r.Body)

	if bodyError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var bookChanged Book
	erroJson := json.Unmarshal(body, &bookChanged)
	if erroJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indexBook := -1
	for index, book := range Books {
		if book.Id == id {
			indexBook = index
			break
		}
	}
	if indexBook < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	Books[indexBook] = bookChanged

	json.NewEncoder(w).Encode(bookChanged)
}

func routerBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) == 2 || len(parts) == 3 && parts[2] == "" {
		if r.Method == "GET" {
			routeListBooks(w, r)
		} else if r.Method == "POST" {
			createBooks(w, r)
		}
	} else if len(parts) == 3 || len(parts) == 4 && parts[3] == "" {
		if r.Method == "GET" {
			getBook(w, r)
		} else if r.Method == "DELETE" {
			removeBook(w, r)
		} else if r.Method == "PUT" {
			alterBook(w, r)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func getBook(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	parts := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(parts[2])

	for _, book := range Books {
		if book.Id == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id, erro := strconv.Atoi(parts[2])

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indexBook := -1
	for index, book := range Books {
		if book.Id == id {
			indexBook = index
			break
		}
	}
	if indexBook < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	leftSide := Books[0:indexBook]
	rightSide := Books[indexBook+1 : len(Books)]

	Books = append(leftSide, rightSide...)
	w.WriteHeader(http.StatusNoContent)

}

func setupRoutes() {
	http.HandleFunc("/", routeIndex)
	//http.HandleFunc("/books", routeListBooks)
	http.HandleFunc("/books/", routerBooks)
	//http.HandleFunc("/books/", getBook)
	//http.HandleFunc("/books/", removeBook)

}

func setupServer() {
	setupRoutes()

	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func main() {
	setupServer()

}
