package library

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"vib-library/library/database"
	"vib-library/library/structs"
)

// errors
var (
	ErrMissingParams error = errors.New("missing params")
)

type library struct {
	db *database.Database
}

// NewLibrary creates new library instance, initializes it and returns it
func NewLibrary() (*library, error) {
	db, err := database.NewDatabase()
	if err != nil {
		return nil, err
	}

	return &library{
		db: db,
	}, nil
}

// AddMember adds a member to the library
func (lib *library) AddMember(w http.ResponseWriter, req *http.Request) {
	// list required params
	requiredParams := []string{
		"firstName",
		"lastName",
	}

	// parse form
	if err := req.ParseForm(); err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error parsing form", http.StatusInternalServerError)
		return
	}

	// check required parameters
	if err := checkRequiredParams(req.Form, requiredParams); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create new member instance
	newMember := structs.Member{
		FirstName: req.Form.Get("firstName"),
		LastName:  req.Form.Get("lastName"),
	}

	// add member to database
	if err := lib.db.AddMember(newMember); err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error adding user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (lib *library) ListMembers(w http.ResponseWriter, req *http.Request) {
	// get members from database
	members, err := lib.db.GetMembers()
	if err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error getting members", http.StatusInternalServerError)
		return
	}

	// write members json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(members)
}

func (lib *library) RentBook(w http.ResponseWriter, req *http.Request) {
	// list required params
	requiredParams := []string{
		"memberId",
		"bookId",
	}

	// parse form
	if err := req.ParseForm(); err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error parsing form", http.StatusInternalServerError)
		return
	}

	// check required parameters
	if err := checkRequiredParams(req.Form, requiredParams); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// note book rental in db
	if err := lib.db.RentBook(req.Form.Get("memberId"), req.Form.Get("bookId")); err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error renting book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (lib *library) ReturnBook(w http.ResponseWriter, req *http.Request) {
	// list required params
	requiredParams := []string{
		"rentId",
	}

	// parse form
	if err := req.ParseForm(); err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error parsing form", http.StatusInternalServerError)
		return
	}

	// check required parameters
	if err := checkRequiredParams(req.Form, requiredParams); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// note book return in db
	if err := lib.db.ReturnBook(req.Form.Get("rentId")); err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error returning book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (lib *library) ListAvailableBooks(w http.ResponseWriter, req *http.Request) {
	// get available books from database
	books, err := lib.db.GetAvailableBooks()
	if err != nil {
		log.Println(err) // log internal server error
		http.Error(w, "error getting available books", http.StatusInternalServerError)
		return
	}

	// write available books json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func checkRequiredParams(requestParams url.Values, requiredParams []string) error {
	for _, rp := range requiredParams {
		if strings.TrimSpace(requestParams.Get(rp)) == "" {
			return ErrMissingParams
		}
	}

	return nil
}
