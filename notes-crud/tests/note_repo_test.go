package tests

import (
	"database/sql"
	"fmt"
	"log"
	"notes-crud/repository"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	if e := godotenv.Load("../.env"); e != nil {
		log.Fatalf("no .env file found")
	}
	connStr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Test DB not responding: %v", err)
	}

	createTableSQL := `
        CREATE TABLE IF NOT EXISTS notes (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `
	if _, err := db.Exec(createTableSQL); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	db.Exec("TRUNCATE notes")

	cleanup := func() {
		db.Exec("TRUNCATE notes")
		db.Close()
	}

	return db, cleanup
}

func TestPostgresNoteRepo_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	content := "Integration test note"
	note, err := repo.Create(content)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if note.Content != content {
		t.Errorf("Expected content '%s', got '%s'", content, note.Content)
	}

	if _, err := uuid.Parse(note.Id); err != nil {
		t.Errorf("Invalid UUID returned: %s", note.Id)
	}

	if note.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestPostgresNoteRepo_GetById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	created, err := repo.Create("Test note for get")
	if err != nil {
		t.Fatalf("Failed to create test note: %v", err)
	}

	note, err := repo.GetById(uuid.MustParse(created.Id))
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}

	if note == nil {
		t.Fatal("Expected note, got nil")
	}

	if note.Content != "Test note for get" {
		t.Errorf("Expected 'Test note for get', got '%s'", note.Content)
	}
}

func TestPostgresNoteRepo_GetById_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	nonExistentID := uuid.New()
	note, err := repo.GetById(nonExistentID)

	if err != nil {
		t.Fatalf("GetById returned error: %v", err)
	}

	if note != nil {
		t.Error("Expected nil, got note")
	}
}

func TestPostgresNoteRepo_GetAll(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	contents := []string{"Note 1", "Note 2", "Note 3"}
	for _, content := range contents {
		_, err := repo.Create(content)
		if err != nil {
			t.Fatalf("Failed to create test note: %v", err)
		}
	}

	notes, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(notes) != 3 {
		t.Errorf("Expected 3 notes, got %d", len(notes))
	}
}

func TestPostgresNoteRepo_UpdateNote(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	created, err := repo.Create("Old content")
	if err != nil {
		t.Fatalf("Failed to create: %v", err)
	}

	id := uuid.MustParse(created.Id)
	updated, err := repo.UpdateNote(id, "New content")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if !updated {
		t.Error("Expected updated=true, got false")
	}

	note, _ := repo.GetById(id)
	if note.Content != "New content" {
		t.Errorf("Expected 'New content', got '%s'", note.Content)
	}
}

func TestPostgresNoteRepo_UpdateNote_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	nonExistentID := uuid.New()
	updated, err := repo.UpdateNote(nonExistentID, "Any content")

	if err != nil {
		t.Fatalf("UpdateNote returned error: %v", err)
	}

	if updated {
		t.Error("Expected updated=false, got true")
	}
}

func TestPostgresNoteRepo_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	created, err := repo.Create("To delete")
	if err != nil {
		t.Fatalf("Failed to create: %v", err)
	}

	id := uuid.MustParse(created.Id)
	deleted, err := repo.Delete(id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if !deleted {
		t.Error("Expected deleted=true, got false")
	}

	note, _ := repo.GetById(id)
	if note != nil {
		t.Error("Note still exists after delete")
	}
}

func TestPostgresNoteRepo_Delete_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository.PostgresNoteRepo{DB: db}

	nonExistentID := uuid.New()
	deleted, err := repo.Delete(nonExistentID)

	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if deleted {
		t.Error("Expected deleted=false, got true")
	}
}
