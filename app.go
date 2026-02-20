package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"photoo/internal/db"
	"photoo/internal/library"
	"photoo/internal/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	db      *sql.DB
	manager *library.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize DB
	dbConn, err := db.InitDB("photoo.db")
	if err != nil {
		log.Fatal(err)
	}
	a.db = dbConn

	// Initialize Manager
	manager, err := library.NewManager("library", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	a.manager = manager
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// GetPhotos returns all photos from the database
func (a *App) GetPhotos() ([]models.Photo, error) {
	rows, err := a.db.Query("SELECT id, original_path, library_path, filename, hash, date_taken, camera_model, latitude, longitude, import_date FROM photos ORDER BY date_taken DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []models.Photo
	for rows.Next() {
		var p models.Photo
		err := rows.Scan(&p.ID, &p.OriginalPath, &p.LibraryPath, &p.Filename, &p.Hash, &p.DateTaken, &p.CameraModel, &p.Latitude, &p.Longitude, &p.ImportDate)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

// SelectFolder opens a dialog to select a folder
func (a *App) SelectFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Folder to Import Photos",
	})
}

// ImportFromFolder triggers an import process for a folder
func (a *App) ImportFromFolder(folderPath string) (int, error) {
	if folderPath == "" {
		return 0, fmt.Errorf("no folder selected")
	}

	count := 0
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
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

		_, err = a.manager.ImportPhoto(path)
		if err == nil {
			count++
			// Emit event for progress (optional)
			runtime.EventsEmit(a.ctx, "photo-imported", path)
		}
		return nil
	})

	return count, err
}
