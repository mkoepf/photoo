package library

import (
	"os"
	"path/filepath"
	"photoo/internal/db"
	"testing"
)

func TestImportPhoto(t *testing.T) {
	// 1. Setup in-memory DB
	testDB, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init in-memory DB: %v", err)
	}
	defer testDB.Close()

	// 2. Create temp library dir
	tempLib, err := os.MkdirTemp("", "photoo-lib-*")
	if err != nil {
		t.Fatalf("Failed to create temp lib: %v", err)
	}
	defer os.RemoveAll(tempLib)

	manager, err := NewManager(tempLib, testDB)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// 3. Create a dummy photo file
	tempDir, err := os.MkdirTemp("", "photoo-source-*")
	if err != nil {
		t.Fatalf("Failed to create temp source dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "test.jpg")
	if err := os.WriteFile(srcPath, []byte("fake-photo-content"), 0644); err != nil {
		t.Fatalf("Failed to write dummy photo: %v", err)
	}

	// 4. Test Import
	photo, err := manager.ImportPhoto(srcPath)
	if err != nil {
		t.Fatalf("ImportPhoto failed: %v", err)
	}

	if photo.ID == 0 {
		t.Errorf("Expected non-zero photo ID")
	}

	// Verify file exists in library
	if _, err := os.Stat(photo.LibraryPath); os.IsNotExist(err) {
		t.Errorf("Imported file does not exist at %s", photo.LibraryPath)
	}

	// 5. Test Duplicate Detection
	_, err = manager.ImportPhoto(srcPath)
	if err == nil {
		t.Errorf("Expected error when importing duplicate, got nil")
	}

	// 6. Test UpdateMetadata
	err = manager.UpdateMetadata(photo.ID, "camera_model", "New Camera")
	if err != nil {
		t.Fatalf("UpdateMetadata failed: %v", err)
	}

	var updatedModel string
	err = testDB.QueryRow("SELECT camera_model FROM photos WHERE id = ?", photo.ID).Scan(&updatedModel)
	if err != nil {
		t.Fatalf("Failed to query updated model: %v", err)
	}
	if updatedModel != "New Camera" {
		t.Errorf("Expected camera_model 'New Camera', got '%s'", updatedModel)
	}
}
