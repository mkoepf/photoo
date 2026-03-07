package main

import (
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"photoo/internal/db"
	"photoo/internal/library"
)

func main() {
	fmt.Println("--- Photoo System Diagnosis ---")

	// 1. Check Database
	dbConn, err := db.InitDB("photoo.db")
	if err != nil {
		fmt.Printf("[FAIL] Database connection: %v\n", err)
		os.Exit(1)
	}
	defer dbConn.Close()
	fmt.Println("[OK] Database connection established.")

	// 2. Check Library Directory
	libPath, _ := filepath.Abs("library")
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		fmt.Printf("[FAIL] Library directory not found: %s\n", libPath)
	} else {
		fmt.Printf("[OK] Library directory found: %s\n", libPath)
	}

	// 3. Check Photos in DB vs FS
	rows, err := dbConn.Query("SELECT id, filename, library_path FROM photos")
	if err != nil {
		fmt.Printf("[FAIL] Querying photos: %v\n", err)
	} else {
		count := 0
		missing := 0
		for rows.Next() {
			var id int
			var filename, libPathAttr string
			rows.Scan(&id, &filename, &libPathAttr)
			count++

			if _, err := os.Stat(libPathAttr); os.IsNotExist(err) {
				missing++
			}
		}
		fmt.Printf("[INFO] Photos in DB: %d\n", count)
		if missing > 0 {
			fmt.Printf("[WARN] Photos missing on disk: %d\n", missing)
		} else {
			fmt.Println("[OK] All DB photos found on disk.")
		}
	}

	// 4. Check Thumbnails Cache
	cachePath := filepath.Join(libPath, ".thumbnails")
	files, _ := os.ReadDir(cachePath)
	fmt.Printf("[INFO] Cached thumbnails: %d\n", len(files))

	// 5. Test Thumbnail Generation for first photo
	rows, err = dbConn.Query("SELECT filename FROM photos LIMIT 1")
	if err == nil && rows.Next() {
		var filename string
		rows.Scan(&filename)
		fmt.Printf("[INFO] Testing thumbnail generation for: %s\n", filename)

		handler := library.NewThumbnailHandler(libPath)
		req := httptest.NewRequest("GET", "/thumbnail/"+filename, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code == 200 {
			fmt.Println("[OK] Thumbnail generation successful.")
			// Check if it's in cache now
			files, _ = os.ReadDir(cachePath)
			fmt.Printf("[INFO] Cached thumbnails after test: %d\n", len(files))
		} else {
			fmt.Printf("[FAIL] Thumbnail generation failed with status %d: %s\n", rr.Code, rr.Body.String())
		}

		// Test double slash in path (common in Wails sometimes)
		req = httptest.NewRequest("GET", "/thumbnail//"+filename, nil)
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code == 200 {
			fmt.Println("[OK] Double-slash path handled correctly.")
		} else {
			fmt.Printf("[WARN] Double-slash path failed with status %d\n", rr.Code)
		}
	}

	fmt.Println("--- Diagnosis Complete ---")
}
