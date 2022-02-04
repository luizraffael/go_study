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

type Contact struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

var Contacts []Contact = []Contact{
	Contact{
		Id:     1,
		Name:  "Joao",
		PhoneNumber: "7777777777",
	},
	Contact{
		Id:     2,
		Name:  "Jose",
		PhoneNumber: "9999999999",
	},
	Contact{
		Id:     3,
		Name:  "Alan",
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

func routerContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) == 2 || len(parts) == 3 && parts[2] == "" {
		if r.Method == "GET" {
			routeListContacts(w, r)
		} else if r.Method == "POST" {
			createContacts(w, r)
		}
	} else if len(parts) == 3 || len(parts) == 4 && parts[3] == "" {
		if r.Method == "GET" {
			getContact(w, r)
		} else if r.Method == "DELETE" {
			removeContact(w, r)
		} else if r.Method == "PUT" {
			alterContact(w, r)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func getContact(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	parts := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(parts[2])

	for _, contact := range Contacts {
		if contact.Id == id {
			json.NewEncoder(w).Encode(contact)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func removeContact(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id, erro := strconv.Atoi(parts[2])

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

func setupRoutes() {
	http.HandleFunc("/", routeIndex)
	//http.HandleFunc("/contacts", routeListContacts)
	http.HandleFunc("/contacts/", routerContacts)
	//http.HandleFunc("/contacts/", getContact)
	//http.HandleFunc("/contacts/", removeContact)

}

func setupServer() {
	setupRoutes()
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func main() {
	setupServer()

}
