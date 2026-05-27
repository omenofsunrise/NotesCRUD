package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notes-crud/handlers"
	"notes-crud/models"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestGetAll_Success(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	id1 := uuid.New()
	id2 := uuid.New()
	mockRepo.AddTestNote(id1, "Note 1")
	mockRepo.AddTestNote(id2, "Note 2")

	req := httptest.NewRequest("GET", "/notes", nil)
	w := httptest.NewRecorder()
	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var notes []models.Note
	json.NewDecoder(w.Body).Decode(&notes)

	if len(notes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(notes))
	}
}

func TestGetAll_DatabaseError(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	mockRepo.GetAllErr = http.ErrHandlerTimeout
	handler := handlers.NewHanlderNote(mockRepo)

	req := httptest.NewRequest("GET", "/notes", nil)
	w := httptest.NewRecorder()
	handler.GetAll(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

func TestAddNote_Success(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	body := `{"content": "Hello World"}`
	req := httptest.NewRequest("POST", "/notes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddNote(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	var note models.Note
	json.NewDecoder(w.Body).Decode(&note)

	if note.Content != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", note.Content)
	}
}

func TestAddNote_EmptyContent(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	body := `{"content": ""}`
	req := httptest.NewRequest("POST", "/notes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestAddNote_TooLong(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	longContent := strings.Repeat("a", 301)
	body := `{"content": "` + longContent + `"}`
	req := httptest.NewRequest("POST", "/notes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestGetNoteById_Success(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	testID := uuid.New()
	mockRepo.AddTestNote(testID, "Test Note")

	req := httptest.NewRequest("GET", "/notes/"+testID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.GetNoteById(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var note models.Note
	json.NewDecoder(w.Body).Decode(&note)

	if note.Content != "Test Note" {
		t.Errorf("Expected 'Test Note', got '%s'", note.Content)
	}
}

func TestGetNoteById_NotFound(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	nonExistentID := uuid.New()
	req := httptest.NewRequest("GET", "/notes/"+nonExistentID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": nonExistentID.String()})
	w := httptest.NewRecorder()

	handler.GetNoteById(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestGetNoteById_InvalidUUID(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	req := httptest.NewRequest("GET", "/notes/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	w := httptest.NewRecorder()

	handler.GetNoteById(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestDeleteNote_Success(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	testID := uuid.New()
	mockRepo.AddTestNote(testID, "To Delete")

	req := httptest.NewRequest("DELETE", "/notes/"+testID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.DeleteNote(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected 204, got %d", w.Code)
	}

	if _, exists := mockRepo.notes[testID]; exists {
		t.Error("Note should be deleted")
	}
}

func TestDeleteNote_NotFound(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	nonExistentID := uuid.New()
	req := httptest.NewRequest("DELETE", "/notes/"+nonExistentID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": nonExistentID.String()})
	w := httptest.NewRecorder()

	handler.DeleteNote(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestUpdateNote_Success(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	testID := uuid.New()
	mockRepo.AddTestNote(testID, "Old content")

	body := `{"content": "Updated content"}`
	req := httptest.NewRequest("PUT", "/notes/"+testID.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.UpdateNote(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var note models.Note
	json.NewDecoder(w.Body).Decode(&note)

	if note.Content != "Updated content" {
		t.Errorf("Expected 'Updated content', got '%s'", note.Content)
	}

	updatedNote, _ := mockRepo.GetById(testID)
	if updatedNote.Content != "Updated content" {
		t.Error("Note not updated in repository")
	}
}

func TestUpdateNote_InvalidUUID(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	body := `{"content": "New content"}`
	req := httptest.NewRequest("PUT", "/notes/not-a-uuid", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	w := httptest.NewRecorder()

	handler.UpdateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["error"] != "invalid uuid" {
		t.Errorf("Expected 'invalid uuid', got '%s'", resp["error"])
	}
}

func TestUpdateNote_NotFound(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	nonExistentID := uuid.New()
	body := `{"content": "New content"}`
	req := httptest.NewRequest("PUT", "/notes/"+nonExistentID.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": nonExistentID.String()})
	w := httptest.NewRecorder()

	handler.UpdateNote(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestUpdateNote_EmptyContent(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	testID := uuid.New()
	mockRepo.AddTestNote(testID, "Some content")

	body := `{"content": ""}`
	req := httptest.NewRequest("PUT", "/notes/"+testID.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.UpdateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["error"] != "content cannot be empty" {
		t.Errorf("Expected 'content cannot be empty', got '%s'", resp["error"])
	}
}

func TestUpdateNote_TooLong(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	testID := uuid.New()
	mockRepo.AddTestNote(testID, "Some content")

	longContent := strings.Repeat("a", 301)
	body := `{"content": "` + longContent + `"}`
	req := httptest.NewRequest("PUT", "/notes/"+testID.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.UpdateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestUpdateNote_DatabaseError(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	mockRepo.UpdateErr = http.ErrHandlerTimeout
	handler := handlers.NewHanlderNote(mockRepo)

	testID := uuid.New()
	mockRepo.AddTestNote(testID, "Some content")

	body := `{"content": "New content"}`
	req := httptest.NewRequest("PUT", "/notes/"+testID.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.UpdateNote(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

func TestUpdateNote_InvalidJSON(t *testing.T) {
	mockRepo := NewMockNoteRepo()
	handler := handlers.NewHanlderNote(mockRepo)

	testID := uuid.New()
	mockRepo.AddTestNote(testID, "Some content")

	body := `{invalid json}`
	req := httptest.NewRequest("PUT", "/notes/"+testID.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.UpdateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
