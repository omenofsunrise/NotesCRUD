package repository

import (
	"database/sql"
	"fmt"
	"notes-crud/models"
	"time"

	"github.com/google/uuid"
)

type PostgresNoteRepo struct {
	DB *sql.DB
}

type NoteRepository interface {
	GetAll() ([]models.Note, error)
	GetById(id uuid.UUID) (*models.Note, error)
	Create(content string) (*models.Note, error)
	Delete(id uuid.UUID) (bool, error)
	UpdateNote(id uuid.UUID, content string) (bool, error)
}

func NewPostgresNoteRepo(db *sql.DB) *PostgresNoteRepo {
	return &PostgresNoteRepo{DB: db}
}

func (r *PostgresNoteRepo) GetAll() ([]models.Note, error) {
	query := `
		SELECT id, content, created_at
		FROM notes`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения записей: %v", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err = rows.Scan(&note.Id, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *PostgresNoteRepo) GetById(id uuid.UUID) (*models.Note, error) {
	query := `
		SELECT id, content, created_at
		FROM notes
		WHERE id = $1`

	var note models.Note
	e := r.DB.QueryRow(query, id).Scan(&note.Id, &note.Content, &note.CreatedAt)
	if e == sql.ErrNoRows {
		return nil, nil
	}
	if e != nil {
		return nil, fmt.Errorf("ошибка получения заметки: %v", e)
	}

	return &note, nil
}

func (r *PostgresNoteRepo) Create(content string) (*models.Note, error) {
	now := time.Now()
	query := `
		INSERT INTO notes (content, created_at)
		VALUES ($1, $2)
		RETURNING id, content, created_at`

	var note models.Note
	err := r.DB.QueryRow(query, content, now).Scan(&note.Id, &note.Content, &note.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting row: %v", err)
	}

	return &note, nil
}

func (r *PostgresNoteRepo) Delete(id uuid.UUID) (bool, error) {
	query := `
		DELETE 
		FROM notes
		WHERE id = $1`

	result, e := r.DB.Exec(query, id)
	if e != nil {
		return false, e
	}

	n, e := result.RowsAffected()
	return n > 0, e
}

func (r *PostgresNoteRepo) UpdateNote(id uuid.UUID, content string) (bool, error) {
	query := `
        UPDATE notes 
        SET content = $1
        WHERE id = $2`

	result, err := r.DB.Exec(query, content, id)
	if err != nil {
		return false, fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
