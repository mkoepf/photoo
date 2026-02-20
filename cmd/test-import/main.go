package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"photoo/internal/db"
	"photoo/internal/library"
)

func main() {
	// Initialize DB
	dbConn, err := db.InitDB("photoo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Initialize Manager
	manager, err := library.NewManager("library", dbConn)
	if err != nil {
		log.Fatal(err)
	}

	// Scan test_data
	testDir := "test_data"
	err = filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Only import JPG/PNG/HEIC
		ext := filepath.Ext(path)
		if ext != ".JPG" && ext != ".jpg" && ext != ".png" && ext != ".HEIC" {
			return nil
		}

		fmt.Printf("Importing: %s...\n", path)
		photo, err := manager.ImportPhoto(path)
		if err != nil {
			fmt.Printf("  Error importing %s: %v\n", path, err)
			return nil
		}
		fmt.Printf("  Success! Filename in library: %s (Date: %s)\n", photo.Filename, photo.DateTaken)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
