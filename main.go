package main

import (
	"log"
	"net/http"
	"vib-library/library"
)

func main() {
	// init lib
	lib, err := library.NewLibrary()
	if err != nil {
		log.Fatalln(err)
	}

	// endpoints
	http.HandleFunc("/AddMember", lib.AddMember)
	http.HandleFunc("/ListMembers", lib.ListMembers)
	http.HandleFunc("/RentBook", lib.RentBook)
	http.HandleFunc("/ReturnBook", lib.ReturnBook)
	http.HandleFunc("/ListAvailableBooks", lib.ListAvailableBooks)

	// listen and serve
	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
