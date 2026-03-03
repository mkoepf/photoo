package library

import (
	"image"
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

	// 4. Perform HTTP request
	// Note: The handler expects the full URL path starting with /thumbnail/
	req := httptest.NewRequest("GET", "/thumbnail/"+photo.Filename, nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// 5. Verify Response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "image/jpeg" {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, "image/jpeg")
	}

	if len(rr.Body.Bytes()) == 0 {
		t.Error("Handler returned empty body")
	}

	// Verify it's actually an image
	_, format, err := image.Decode(rr.Body)
	if err != nil {
		t.Fatalf("Failed to decode response body as image: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("Expected jpeg response, got %s", format)
	}
}
