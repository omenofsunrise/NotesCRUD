package main

import (
	"fmt"
	"log"
	"net/http"
	"notes-crud/database"
	"notes-crud/handlers"
	"notes-crud/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file in project")
	}

	database.InitDB()
	defer database.DB.Close()

	noteRepo := repository.NewPostgresNoteRepo(database.DB)
	noteHandler := handlers.NewHanlderNote(noteRepo)
	r := mux.NewRouter()

	r.HandleFunc("/notes", noteHandler.GetAll).Methods("GET")
	r.HandleFunc("/notes/{id}", noteHandler.GetNoteById).Methods("GET")
	r.HandleFunc("/notes/{id}", noteHandler.DeleteNote).Methods("DELETE")
	r.HandleFunc("/notes", noteHandler.UpdateNote).Methods("PUT")
	r.HandleFunc("/notes", noteHandler.AddNote).Methods("POST")

	fmt.Println("I STARTED SERVER")
	log.Fatal(http.ListenAndServe(":8080", r))
}
