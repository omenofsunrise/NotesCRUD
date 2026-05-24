package handlers

import (
	"encoding/json"
	"net/http"
	"notes-crud/models"

	"github.com/gorilla/mux"
)

var notes = []models.Note{
	{ID: "19msdlk432v;4432", CreatedAt: "10.10.2005", Content: "HELLO!"},
	{ID: "219msdlk432vdasds;4432", CreatedAt: "10.11.2005", Content: "Goodbye!"},
}

// GET notes
func GetNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// POST notes
func AddNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ContentType", "application/json")
	var newNote models.Note
	json.NewDecoder(r.Body).Decode(&newNote)
	notes = append(notes, newNote)
	json.NewEncoder(w).Encode(newNote)
}

// GET notes/{id}
func GetNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var id = params["id"]
	var note *models.Note
	for _, v := range notes {
		if v.ID == id {
			note = &v
			break
		}
	}

	if note == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "not found",
		})
		return
	}
	json.NewEncoder(w).Encode(note)
}

// DELETE notes/{id}
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id required"})
	}

	var note *models.Note
	for i, v := range notes {
		if v.ID == id {
			note = &v
			notes = append(notes[:i], notes[i+1:]...)
		}
	}
	if note == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "note to delete not found"})
	} else {
		json.NewEncoder(w).Encode(notes)
	}
}
