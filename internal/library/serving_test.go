package library

import (
	_ "image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"photoo/internal/db"
	"testing"
)

func TestThumbnailHTTPHandler(t *testing.T) {
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

	// 3. Setup the Handler
	handler := NewThumbnailHandler(libPath)

	// 4. Test various path formats
	paths := []struct {
		path           string
		expectedStatus int
	}{
		{"/thumbnail/" + photo.Filename, http.StatusOK},
		{"/thumbnail//" + photo.Filename, http.StatusOK}, // Double slash
		{"/thumbnail/missing.jpg", http.StatusNotFound},
	}

	for _, tc := range paths {
		req := httptest.NewRequest("GET", tc.path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != tc.expectedStatus {
			t.Errorf("Path %s failed: expected status %v, got %v", tc.path, tc.expectedStatus, rr.Code)
		}
	}
}
