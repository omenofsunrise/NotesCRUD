package main

import (
	"fmt"
	"log"
	"net/http"
	"notes-crud/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/notes", handlers.GetNotes).Methods("GET")
	r.HandleFunc("/notes/{id}", handlers.GetNote).Methods("GET")
	r.HandleFunc("/notes/{id}", handlers.DeleteNote).Methods("DELETE")
	r.HandleFunc("/notes", handlers.AddNote).Methods("POST")

	fmt.Println("I STARTED SERVER")
	log.Fatal(http.ListenAndServe(":8080", r))
}
