package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createSchema(db); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original_path TEXT,
			library_path TEXT NOT NULL,
			filename TEXT NOT NULL UNIQUE,
			hash TEXT NOT NULL,
			date_taken DATETIME,
			camera_model TEXT,
			latitude REAL,
			longitude REAL,
			import_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_photos_hash ON photos(hash);`,
		`CREATE INDEX IF NOT EXISTS idx_photos_date_taken ON photos(date_taken);`,
		`CREATE TABLE IF NOT EXISTS metadata_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			photo_id INTEGER,
			field_name TEXT NOT NULL,
			old_value TEXT,
			new_value TEXT,
			changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (photo_id) REFERENCES photos(id)
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	log.Println("Database schema initialized.")
	return nil
}
