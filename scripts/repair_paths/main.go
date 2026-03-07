package main

import (
	"fmt"
	"log"
	"photoo/internal/db"
	"strings"
)

func main() {
	fmt.Println("--- Photoo Library Path Repair ---")

	dbConn, err := db.InitDB("photoo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT id, filename, library_path FROM photos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	updatedCount := 0
	for rows.Next() {
		var id int64
		var filename, libPath string
		rows.Scan(&id, &filename, &libPath)

		// If filename is flat but libPath contains subfolders
		if !strings.Contains(filename, "/") && !strings.Contains(filename, "\\") {
			// Find where 'library/' is in the full path
			parts := strings.Split(libPath, "/library/")
			if len(parts) > 1 {
				newFilename := parts[1]
				fmt.Printf("[FIX] Updating ID %d: %s -> %s\n", id, filename, newFilename)
				_, err = dbConn.Exec("UPDATE photos SET filename = ? WHERE id = ?", newFilename, id)
				if err != nil {
					fmt.Printf("[ERR] Failed to update ID %d: %v\n", id, err)
				} else {
					updatedCount++
				}
			}
		}
	}

	fmt.Printf("--- Repair Complete. Updated %d records ---\n", updatedCount)
}
