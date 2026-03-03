package library

import (
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"photoo/internal/db"
	"strings"
	"testing"
)

func TestThumbnailServingLogic(t *testing.T) {
	// 1. Setup temporary environment
	tempDir, err := os.MkdirTemp("", "photoo-test-*")
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

	libPath := filepath.Join(tempDir, "library")
	manager, err := NewManager(libPath, dbConn)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Import a real test photo
	wd, _ := os.Getwd()
	// Go up one level to project root if running from internal/library
	projectRoot := filepath.Dir(filepath.Dir(wd))
	testPhoto := filepath.Join(projectRoot, "test_data", "source_digital_camera", "RIMG0018.JPG")

	photo, err := manager.ImportPhoto(testPhoto)
	if err != nil {
		t.Fatalf("Failed to import photo: %v", err)
	}

	// 3. Verify the file exists in the library
	fullPath := filepath.Join(libPath, photo.Filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Fatalf("Imported file does not exist at %s", fullPath)
	}

	// 4. Test logic that matches the ServeHTTP implementation
	requestPath := "/thumbnail/" + photo.Filename
	filename := strings.TrimPrefix(requestPath, "/thumbnail/")
	if filename != photo.Filename {
		t.Errorf("Path parsing failed: expected %s, got %s", photo.Filename, filename)
	}

	testFullPath := filepath.Join(libPath, filename)
	if _, err := os.Stat(testFullPath); os.IsNotExist(err) {
		t.Fatalf("Handler would fail to find file at %s", testFullPath)
	}

	// Try to open it as an image (simulating imaging.Open)
	f, err := os.Open(testFullPath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	_, format, err := image.Decode(f)
	if err != nil {
		t.Fatalf("Failed to decode image: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("Expected jpeg, got %s", format)
	}
}
