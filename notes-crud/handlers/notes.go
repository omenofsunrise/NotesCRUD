package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notes-crud/models"
	"notes-crud/repository"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type NoteHandler struct {
	r *repository.NoteRepo
}

func NewHanlderNote(repo *repository.NoteRepo) *NoteHandler {
	return &NoteHandler{r: repo}
}

// GET notes
func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	notes, e := h.r.GetAll()
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("an error by getting notes: %v", e)
	}
	json.NewEncoder(w).Encode(notes)
}

// POST notes
func (h *NoteHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var note models.Note
	json.NewDecoder(r.Body).Decode(&note)
	note.Content = strings.Trim(note.Content, " \n\t\r")
	if len(note.Content) > 300 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "note content limit 300 chars is up"})
		return
	}
	if note.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "note content is empty"})
		return
	}
	newNote, e := h.r.Create(note.Content)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error: %v", e)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	json.NewEncoder(w).Encode(newNote)
}

// GET notes/{id}
func (h *NoteHandler) GetNoteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	idUud, e := uuid.Parse(id)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid uuid"})
		return
	}
	note, e := h.r.GetById(idUud)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error getting note by id: %v", e)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}
	if note == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
		return
	}

	json.NewEncoder(w).Encode(note)
}

// DELETE notes/{id}
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id required"})
		return
	}
	noteUUID, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid uuid"})
		return
	}
	success, err := h.r.Delete(noteUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "unhandled db error"})
		log.Printf("error: %v", err)
		return
	}

	if !success {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "record not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"success": "true"})
}

// PUT /notes/{id}
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.Note

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	idStr := req.Id

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id required"})
		return
	}

	noteID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid uuid"})
		return
	}

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "content cannot be empty"})
		return
	}

	if len(req.Content) > 300 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "content too long (max 300)"})
		return
	}

	updated, err := h.r.UpdateNote(noteID, req.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		log.Printf("Error updating note: %v", err)
		return
	}

	if !updated {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "note not found"})
		return
	}

	json.NewEncoder(w).Encode(updated)
}
