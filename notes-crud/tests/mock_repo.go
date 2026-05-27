package tests

import (
	"notes-crud/models"

	"github.com/google/uuid"
)

type MockNoteRepo struct {
	notes      map[uuid.UUID]models.Note
	GetAllErr  error
	GetByIdErr error
	CreateErr  error
	DeleteErr  error
	UpdateErr  error
}

func NewMockNoteRepo() *MockNoteRepo {
	return &MockNoteRepo{
		notes: make(map[uuid.UUID]models.Note),
	}
}

func (m *MockNoteRepo) AddTestNote(id uuid.UUID, content string) {
	m.notes[id] = models.Note{
		Id:      id.String(),
		Content: content,
	}
}

func (m *MockNoteRepo) GetAll() ([]models.Note, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	var notes []models.Note
	for _, note := range m.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

func (m *MockNoteRepo) GetById(id uuid.UUID) (*models.Note, error) {
	if m.GetByIdErr != nil {
		return nil, m.GetByIdErr
	}
	if note, ok := m.notes[id]; ok {
		return &note, nil
	}
	return nil, nil
}

func (m *MockNoteRepo) Create(content string) (*models.Note, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	id := uuid.New()
	note := models.Note{
		Id:      id.String(),
		Content: content,
	}
	m.notes[id] = note
	return &note, nil
}

func (m *MockNoteRepo) Delete(id uuid.UUID) (bool, error) {
	if m.DeleteErr != nil {
		return false, m.DeleteErr
	}
	if _, ok := m.notes[id]; !ok {
		return false, nil
	}
	delete(m.notes, id)
	return true, nil
}

func (m *MockNoteRepo) UpdateNote(id uuid.UUID, content string) (bool, error) {
	if m.UpdateErr != nil {
		return false, m.UpdateErr
	}
	if _, ok := m.notes[id]; !ok {
		return false, nil
	}
	note := m.notes[id]
	note.Content = content
	m.notes[id] = note
	return true, nil
}
