package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
	id, erro := strconv.Atoi(vars["bookId"])

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

func getBook(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Path)

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["bookId"])

	for _, book := range Books {
		if book.Id == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, erro := strconv.Atoi(vars["bookId"])

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

func setupRoutes(router *mux.Router) {

	router.HandleFunc("/", routeIndex)
	router.HandleFunc("/books", routeListBooks).Methods("GET")
	router.HandleFunc("/books/{bookId}", getBook).Methods("GET")
	router.HandleFunc("/books", createBooks).Methods("POST")
	router.HandleFunc("/books/{bookId}", alterBook).Methods("PUT")
	router.HandleFunc("/books/{bookId}", removeBook).Methods("DELETE")

}
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func setupServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(jsonMiddleware)
	setupRoutes(router)

	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
func main() {
	setupServer()

}
