package library

import (
	"fmt"
	_ "image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"photoo/internal/db"
	"testing"
)

func TestThumbnailHTTPHandler(t *testing.T) {
	// ... (Setup code from existing TestThumbnailHTTPHandler)
}

func TestThumbnailConcurrency(t *testing.T) {
	// 1. Setup temporary environment
	tempDir, err := os.MkdirTemp("", "photoo-concurrency-*")
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
	projectRoot := filepath.Dir(filepath.Dir(wd))
	testPhoto := filepath.Join(projectRoot, "test_data", "source_digital_camera", "RIMG0018.JPG")

	photo, err := manager.ImportPhoto(testPhoto)
	if err != nil {
		t.Fatalf("Failed to import photo: %v", err)
	}

	// 3. Setup the Handler
	handler := NewThumbnailHandler(libPath)

	// 4. Fire multiple concurrent requests for the SAME photo
	// This tests the locking/deduplication mechanism
	const concurrentRequests = 10
	errChan := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/thumbnail/"+photo.Filename, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				errChan <- fmt.Errorf("concurrent request failed with status %d", rr.Code)
			} else {
				errChan <- nil
			}
		}()
	}

	for i := 0; i < concurrentRequests; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("Concurrent request error: %v", err)
		}
	}
}
