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

// GET /notes
func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, e := h.r.GetAll()
	if e != nil {
		setErrorResponse(w, http.StatusInternalServerError, "error getting notes")
		log.Printf("an error by getting notes: %v", e)
		return
	}
	setResponse(w, http.StatusOK, notes)
}

// POST /notes
func (h *NoteHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	json.NewDecoder(r.Body).Decode(&note)
	note.Content = strings.Trim(note.Content, " \n\t\r")
	if len(note.Content) > 300 {
		setErrorResponse(w, http.StatusBadRequest, "note content limit 300 chars is up")
		return
	}
	if note.Content == "" {
		setErrorResponse(w, http.StatusBadRequest, "note content is empty")
		return
	}
	newNote, e := h.r.Create(note.Content)
	if e != nil {
		fmt.Printf("error: %v", e)
		setErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	setResponse(w, http.StatusCreated, newNote)
}

// GET /notes/{id}
func (h *NoteHandler) GetNoteById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idUud, e := uuid.Parse(id)
	if e != nil {
		setErrorResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	note, e := h.r.GetById(idUud)
	if e != nil {
		log.Printf("error getting note by id: %v", e)
		setErrorResponse(w, http.StatusInternalServerError, "database error")
		return
	}
	if note == nil {
		setErrorResponse(w, http.StatusNotFound, "not found")
		return
	}

	setResponse(w, http.StatusOK, note)
}

// DELETE /notes/{id}
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		setErrorResponse(w, http.StatusBadRequest, "id required")
		return
	}
	noteUUID, err := uuid.Parse(id)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, "invald uuid")
		return
	}
	success, err := h.r.Delete(noteUUID)
	if err != nil {
		setErrorResponse(w, http.StatusInternalServerError, "unhandled database error")
		log.Printf("error: %v", err)
		return
	}

	if !success {
		setErrorResponse(w, http.StatusNotFound, "record not found")
		return
	}
	setResponse(w, http.StatusNoContent, nil)
}

// PUT /notes/{id}
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var req models.Note

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		setErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	idStr := req.Id

	if idStr == "" {
		setErrorResponse(w, http.StatusBadRequest, "id required")
		return
	}

	noteID, err := uuid.Parse(idStr)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	if req.Content == "" {
		setErrorResponse(w, http.StatusBadRequest, "content cannot be empty")
		return
	}

	if len(req.Content) > 300 {
		setErrorResponse(w, http.StatusBadRequest, "content too long (max 300)")
		return
	}

	updated, err := h.r.UpdateNote(noteID, req.Content)
	if err != nil {
		setErrorResponse(w, http.StatusInternalServerError, "unhanled database error")
		log.Printf("Error updating note: %v", err)
		return
	}

	if !updated {
		setErrorResponse(w, http.StatusNotFound, "note not found")
		return
	}

	setResponse(w, http.StatusOK, req)
}
