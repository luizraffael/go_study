package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type ResponseErr struct {
	Error string `json:"erro"`
}
type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var db *sql.DB

func routeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome")
}

func routeListBooks(w http.ResponseWriter, r *http.Request) {
	rows, errorQuery := db.Query("SELECT * FROM books")

	if errorQuery != nil {
		log.Println(errorQuery.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var books []Book = make([]Book, 0)
	for rows.Next() {
		var book Book
		erroScan := rows.Scan(&book.Id, &book.Author, &book.Title)
		if erroScan != nil {
			log.Println("RouteListBooks: erroScan: " + erroScan.Error())
			continue
		}
		books = append(books, book)
	}
	erroCloseRows := rows.Close()

	if erroCloseRows != nil {
		log.Println("Close rows: erroCloseRows: " + erroCloseRows.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(books)
}

func createBooks(w http.ResponseWriter, r *http.Request) {
	//
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

	}
	var newBook Book
	json.Unmarshal(body, &newBook)

	erroValidBook := validateBook(newBook)

	if len(erroValidBook) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseErr{erroValidBook})
		return
	}

	result, erroInsert := db.Exec("INSERT INTO books (author, title) VALUES (?, ?)",
		newBook.Author, newBook.Title)

	id, errorId := result.LastInsertId()

	if errorId != nil {
		log.Println("Errr trying to get ID" + errorId.Error())
	}

	if erroInsert != nil {
		log.Println("Error insert book" + erroInsert.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newBook.Id = int(id)
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(newBook)
}

func editBook(w http.ResponseWriter, r *http.Request) {
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

	row := db.QueryRow(
		"SELECT id, author, title FROM books WHERE id = ?", id)

	var book Book
	erroScan := row.Scan(&book.Id, &book.Author, &book.Title)

	if erroScan != nil {
		log.Println("Erro scan: Scan erro getBook: " + erroScan.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, erroExec := db.Exec("UPDATE books SET author = ?, title = ? WHERE id = ?",
		bookChanged.Author, bookChanged.Title, id)

	if erroExec != nil {
		log.Println("Erro changing book: " + erroExec.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bookChanged)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Path)

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["bookId"])

	row := db.QueryRow(
		"SELECT id, author, title FROM books WHERE id = ?", id)

	var book Book
	erroScan := row.Scan(&book.Id, &book.Author, &book.Title)

	if erroScan != nil {
		log.Println("Erro scan: Scan erro getBook: " + erroScan.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)

}

func removeBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, erro := strconv.Atoi(vars["bookId"])

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id FROM books WHERE id = ?", id)

	var bookId int

	erroScan := row.Scan(&bookId)

	if erroScan != nil {
		log.Println("Erro scan: Scan erro deletBook: " + erroScan.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, erroExec := db.Exec("DELETE FROM books WHERE id = ?", id)

	if erroExec != nil {
		log.Println("Erro deleting book: " + erroExec.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func validateBook(book Book) string {
	if len(book.Author) == 0 || len(book.Author) > 50 {
		return "Author can't be null or have more than 50 characters"
	}
	if len(book.Title) == 0 || len(book.Title) > 100 {
		return "Title can't be null or have more than 100 characters"
	}
	return ""
}
func setupRoutes(router *mux.Router) {

	router.HandleFunc("/", routeIndex)
	router.HandleFunc("/books", routeListBooks).Methods("GET")
	router.HandleFunc("/books/{bookId}", getBook).Methods("GET")
	router.HandleFunc("/books", createBooks).Methods("POST")
	router.HandleFunc("/books/{bookId}", editBook).Methods("PUT")
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
	setupDataBase()
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupDataBase() {

	var erroOpen error

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var stringConnection string = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	db, erroOpen = sql.Open("mysql", stringConnection)
	if erroOpen != nil {
		log.Fatal(erroOpen.Error())
	}
	erroPing := db.Ping()
	if erroPing != nil {
		log.Fatal(erroPing.Error())
	}
}

func main() {
	setupServer()

}
