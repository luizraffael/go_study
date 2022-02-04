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

type Contact struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

var Contacts []Contact = []Contact{
	Contact{
		Id:          1,
		Name:        "Joao",
		PhoneNumber: "7777777777",
	},
	Contact{
		Id:          2,
		Name:        "Jose",
		PhoneNumber: "9999999999",
	},
	Contact{
		Id:          3,
		Name:        "Alan",
		PhoneNumber: "88888888888",
	},
}

func routeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome")
}

func routeListContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(Contacts)
}

func createContacts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		log.Fatal(error)
	}
	var newContact Contact
	json.Unmarshal(body, &newContact)
	newContact.Id = len(Contacts) + 1
	Contacts = append(Contacts, newContact)
	encoder := json.NewEncoder(w)
	encoder.Encode(newContact)
}

func alterContact(w http.ResponseWriter, r *http.Request) {
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

	indexContact := -1
	for index, contact := range Contacts {
		if contact.Id == id {
			indexContact = index
			break
		}
	}
	if indexContact < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	Contacts[indexContact] = contactChanged

	json.NewEncoder(w).Encode(contactChanged)
}

func getContact(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["contactId"])

	for _, contact := range Contacts {
		if contact.Id == id {
			json.NewEncoder(w).Encode(contact)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func removeContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, erro := strconv.Atoi(vars["contactId"])

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indexContact := -1
	for index, contact := range Contacts {
		if contact.Id == id {
			indexContact = index
			break
		}
	}
	if indexContact < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	leftSide := Contacts[0:indexContact]
	rightSide := Contacts[indexContact+1 : len(Contacts)]

	Contacts = append(leftSide, rightSide...)
	w.WriteHeader(http.StatusNoContent)

}

func setupRoutes(router *mux.Router) {

	router.HandleFunc("/", routeIndex)
	router.HandleFunc("/contacts", routeListContacts).Methods("GET")
	router.HandleFunc("/contacts/{contactId}", getContact).Methods("GET")
	router.HandleFunc("/contacts", createContacts).Methods("POST")
	router.HandleFunc("/contacts/{contactId}", alterContact).Methods("PUT")
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
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
func main() {
	setupServer()

}
