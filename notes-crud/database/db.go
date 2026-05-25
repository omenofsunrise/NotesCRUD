package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf(
		"host=%v port=%v password=%v user=%v dbname=%v",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
	)
	var err error
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Ошибка подключения к бд", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB не отвечает!", err)
	}

	fmt.Println("Connected to postgres")

	if e := CreateTable(); e != nil {
		log.Fatalf("ERROR CREATING TABLE: %v", e)
	}
}

func CreateTable() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS notes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	if _, e := DB.Exec(createTableQuery); e != nil {
		return fmt.Errorf("%v", e)
	}

	createIndexQuery := "CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);"
	if _, e := DB.Exec(createIndexQuery); e != nil {
		return fmt.Errorf("%v", e)
	}
	return nil
}
