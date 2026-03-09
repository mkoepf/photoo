package main

import (
	"os"
	"path/filepath"
	"photoo/internal/db"
	"testing"
	"time"
)

func TestGetPhotosPaged(t *testing.T) {
	// 1. Setup temporary DB
	tempDir, err := os.MkdirTemp("", "photoo-app-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	dbConn, err := db.InitDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer dbConn.Close()

	// 2. Insert dummy data
	for i := 1; i <= 10; i++ {
		_, err := dbConn.Exec(
			"INSERT INTO photos (original_path, library_path, filename, hash, date_taken, camera_model) VALUES (?, ?, ?, ?, ?, ?)",
			"orig/path",
			"path/to/photo",
			string(rune('a'+i)),
			"hash",
			time.Now().Add(time.Duration(-i)*time.Hour),
			"Mock Camera",
		)
		if err != nil {
			t.Fatalf("Failed to insert dummy photo %d: %v", i, err)
		}
	}

	app := &App{db: dbConn}

	// 3. Test first page
	photos, err := app.GetPhotosPaged(0, 3)
	if err != nil {
		t.Fatalf("GetPhotosPaged failed: %v", err)
	}
	if len(photos) != 3 {
		t.Errorf("Expected 3 photos, got %d", len(photos))
	}

	// 4. Test second page
	photos2, err := app.GetPhotosPaged(3, 3)
	if err != nil {
		t.Fatalf("GetPhotosPaged failed: %v", err)
	}
	if len(photos2) != 3 {
		t.Errorf("Expected 3 photos on second page, got %d", len(photos2))
	}

	// Ensure different items on different pages
	if photos[0].ID == photos2[0].ID {
		t.Errorf("Duplicate item on different pages: ID %d", photos[0].ID)
	}

	// 5. Test last partial page
	photos3, err := app.GetPhotosPaged(9, 3)
	if err != nil {
		t.Fatalf("GetPhotosPaged failed: %v", err)
	}
	if len(photos3) != 1 {
		t.Errorf("Expected 1 photo on last page, got %d", len(photos3))
	}
}
