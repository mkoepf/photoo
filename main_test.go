package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"photoo/internal/db"
	"photoo/internal/library"
	"testing"
)

func TestThumbnailIntegration(t *testing.T) {
	// 1. Setup temporary library
	tempDir, err := os.MkdirTemp("", "photoo-integration-*")
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
	os.MkdirAll(libPath, 0755)

	// Create a nested test image
	nestedDir := filepath.Join(libPath, "2024/01/01")
	os.MkdirAll(nestedDir, 0755)

	wd, _ := os.Getwd()
	realPhoto := filepath.Join(wd, "test_data/source_digital_camera/RIMG0018.JPG")
	data, err := os.ReadFile(realPhoto)
	if err != nil {
		t.Skip("Skipping test: real test photo not found")
		return
	}

	targetPhoto := filepath.Join(nestedDir, "test.jpg")
	os.WriteFile(targetPhoto, data, 0644)

	// 2. Setup Handler (simulating main.go AssetServer logic)
	thumbHandler := library.NewThumbnailHandler(libPath)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thumbHandler.ServeHTTP(w, r)
	}))
	defer server.Close()

	// 3. Test requests
	testCases := []struct {
		url            string
		expectedStatus int
	}{
		{"/thumbnail/2024/01/01/test.jpg", http.StatusOK},
		{"/thumbnail//2024/01/01/test.jpg", http.StatusOK},
		{"/thumbnail/missing.jpg", http.StatusNotFound},
	}

	for _, tc := range testCases {
		resp, err := http.Get(server.URL + tc.url)
		if err != nil {
			t.Errorf("Request %s failed: %v", tc.url, err)
			continue
		}
		if resp.StatusCode != tc.expectedStatus {
			t.Errorf("URL %s: expected status %d, got %d", tc.url, tc.expectedStatus, resp.StatusCode)
		}
		resp.Body.Close()
	}
}
