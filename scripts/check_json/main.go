package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"photoo/internal/db"
	"photoo/internal/models"
)

// Simplified App for testing
type TestApp struct {
	db *sql.DB
}

func main() {
	dbConn, err := db.InitDB("photoo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT id, original_path, library_path, filename, hash, date_taken, camera_model, latitude, longitude, import_date FROM photos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var photos []models.Photo
	for rows.Next() {
		var p models.Photo
		err := rows.Scan(&p.ID, &p.OriginalPath, &p.LibraryPath, &p.Filename, &p.Hash, &p.DateTaken, &p.CameraModel, &p.Latitude, &p.Longitude, &p.ImportDate)
		if err != nil {
			log.Fatal(err)
		}
		photos = append(photos, p)
	}

	data, _ := json.MarshalIndent(photos, "", "  ")
	fmt.Println("--- JSON Output of GetPhotos ---")
	fmt.Println(string(data))
}
