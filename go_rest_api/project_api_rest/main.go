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
type Contact struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

var db *sql.DB

func routeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome")
}

func routeListContacts(w http.ResponseWriter, r *http.Request) {
	rows, errorQuery := db.Query("SELECT * FROM contacts")

	if errorQuery != nil {
		log.Println(errorQuery.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contacts []Contact = make([]Contact, 0)
	for rows.Next() {
		var contact Contact
		erroScan := rows.Scan(&contact.Id, &contact.Name, &contact.PhoneNumber)
		if erroScan != nil {
			log.Println("RouteListContacts: erroScan: " + erroScan.Error())
			continue
		}
		contacts = append(contacts, contact)
	}
	erroCloseRows := rows.Close()

	if erroCloseRows != nil {
		log.Println("Close rows: erroCloseRows: " + erroCloseRows.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(contacts)
}

func createContacts(w http.ResponseWriter, r *http.Request) {
	//
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

	}
	var newContact Contact
	json.Unmarshal(body, &newContact)

	erroValidContact := validateContact(newContact)

	if len(erroValidContact) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseErr{erroValidContact})
		return
	}

	result, erroInsert := db.Exec("INSERT INTO contacts (phone_number, name) VALUES (?, ?)",
		newContact.PhoneNumber, newContact.Name)

	id, errorId := result.LastInsertId()

	if errorId != nil {
		log.Println("Errr trying to get ID" + errorId.Error())
	}

	if erroInsert != nil {
		log.Println("Error insert contact" + erroInsert.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newContact.Id = int(id)
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(newContact)
}

func editContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, erro := strconv.Atoi(vars["contactId"])

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, bodyError := ioutil.ReadAll(r.Body)

	if bodyError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contactChanged Contact
	erroJson := json.Unmarshal(body, &contactChanged)
	if erroJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := db.QueryRow(
		"SELECT id, phone_number, name FROM contacts WHERE id = ?", id)

	var contact Contact
	erroScan := row.Scan(&contact.Id, &contact.PhoneNumber, &contact.Name)

	if erroScan != nil {
		log.Println("Erro scan: Scan erro getContact: " + erroScan.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, erroExec := db.Exec("UPDATE contacts SET phone_number = ?, name = ? WHERE id = ?",
		contactChanged.PhoneNumber, contactChanged.Name, id)

	if erroExec != nil {
		log.Println("Erro changing contact: " + erroExec.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contactChanged)
}

func getContact(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Path)

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["contactId"])

	row := db.QueryRow(
		"SELECT id, phone_number, name FROM contacts WHERE id = ?", id)

	var contact Contact
	erroScan := row.Scan(&contact.Id, &contact.PhoneNumber, &contact.Name)

	if erroScan != nil {
		log.Println("Erro scan: Scan erro getContact: " + erroScan.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(contact)

}

func removeContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, erro := strconv.Atoi(vars["contactId"])

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id FROM contacts WHERE id = ?", id)

	var contactId int

	erroScan := row.Scan(&contactId)

	if erroScan != nil {
		log.Println("Erro scan: Scan erro deletContact: " + erroScan.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, erroExec := db.Exec("DELETE FROM contacts WHERE id = ?", id)

	if erroExec != nil {
		log.Println("Erro deleting contact: " + erroExec.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func validateContact(contact Contact) string {
	if len(contact.PhoneNumber) == 0 || len(contact.PhoneNumber) > 50 {
		return "PhoneNumber can't be null or have more than 50 characters"
	}
	if len(contact.Name) == 0 || len(contact.Name) > 100 {
		return "Name can't be null or have more than 100 characters"
	}
	return ""
}
func setupRoutes(router *mux.Router) {

	router.HandleFunc("/", routeIndex)
	router.HandleFunc("/contacts", routeListContacts).Methods("GET")
	router.HandleFunc("/contacts/{contactId}", getContact).Methods("GET")
	router.HandleFunc("/contacts", createContacts).Methods("POST")
	router.HandleFunc("/contacts/{contactId}", editContact).Methods("PUT")
	router.HandleFunc("/contacts/{contactId}", removeContact).Methods("DELETE")

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
