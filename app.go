package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"photoo/internal/db"
	"photoo/internal/library"
	"photoo/internal/models"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	db       *sql.DB
	manager  *library.Manager
	thumbH   *library.ThumbnailHandler
	uiLogs   []string
	uiErrors []string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		uiLogs:   make([]string, 0),
		uiErrors: make([]string, 0),
	}
}

func (a *App) SetThumbnailHandler(h *library.ThumbnailHandler) {
	a.thumbH = h
}

// GetThumbnail returns a base64 encoded thumbnail for a photo
func (a *App) GetThumbnail(filename string) (string, error) {
	libPath := a.manager.LibraryPath
	cachePath := filepath.Join(libPath, ".thumbnails")
	// replace all slashes (both / and \) with underscores for flat cache
	safeName := strings.ReplaceAll(filename, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	cacheFilename := safeName + ".thumb.jpg"
	cacheFullPath := filepath.Join(cachePath, cacheFilename)

	if _, err := os.Stat(cacheFullPath); os.IsNotExist(err) {
		// Try to trigger generation if handler is available
		if a.thumbH != nil {
			// Fake a request to trigger generation
			req := httptest.NewRequest("GET", "/thumbnail/"+filename, nil)
			rr := httptest.NewRecorder()
			a.thumbH.ServeHTTP(rr, req)
		}
	}

	data, err := os.ReadFile(cacheFullPath)
	if err != nil {
		return "", err
	}

	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(data), nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dbName := "photoo.db"
	libDir := "library"

	if os.Getenv("PHOTOO_SELF_TEST") == "true" {
		dbName = "photoo_self_test.db"
		libDir = "library_self_test"
		os.Remove(dbName)
		os.RemoveAll(libDir)
		os.MkdirAll(libDir, 0755)
		os.MkdirAll(filepath.Join(libDir, ".thumbnails"), 0755)
	}

	if a.db == nil {
		// Initialize DB
		dbConn, err := db.InitDB(dbName)
		if err != nil {
			log.Fatal(err)
		}
		a.db = dbConn
	}

	if a.manager == nil {
		libPath, err := filepath.Abs(libDir)
		if err != nil {
			libPath = libDir
		}
		// Initialize Manager
		manager, err := library.NewManager(libPath, a.db)
		if err != nil {
			log.Fatal(err)
		}
		a.manager = manager
	}

	if os.Getenv("PHOTOO_SELF_TEST") == "true" {
		go a.runSelfTest()
	}
}

func (a *App) runSelfTest() {
	time.Sleep(30 * time.Second) // Give it plenty of time to load
	fmt.Println("[AUTO] Starting Self Test...")

	wd, _ := os.Getwd()
	testPath := filepath.Join(wd, "test_data/source_digital_camera")
	fmt.Printf("[AUTO] Importing from %s\n", testPath)

	// Send command to frontend to trigger import
	a.SendCommand("trigger_import", testPath)

	time.Sleep(10 * time.Second) // Wait for import and thumbnails

	// Verify folder structure manually in logs
	fmt.Println("[AUTO] Verifying folder structure...")
	filepath.Walk(a.manager.LibraryPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fmt.Printf("[AUTO] Found file: %s\n", path)
		}
		return nil
	})

	fmt.Println("[AUTO] Requesting thumbnail inspection...")
	a.SendCommand("inspect_thumbnails", nil)
}
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// GetPhotos returns all photos from the database
func (a *App) GetPhotos() ([]models.Photo, error) {
	return a.GetPhotosPaged(0, 1000000)
}

// GetPhotosPaged returns a page of photos from the database
func (a *App) GetPhotosPaged(offset, limit int) ([]models.Photo, error) {
	rows, err := a.db.Query("SELECT id, original_path, library_path, filename, hash, date_taken, camera_model, latitude, longitude, import_date FROM photos ORDER BY date_taken DESC LIMIT ? OFFSET ?", limit, offset)
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

	// 1. Initial pass to count total candidate files
	var totalCandidates int
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".png" || ext == ".heic" {
			totalCandidates++
		}
		return nil
	})

	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "import:start", map[string]interface{}{
			"total": totalCandidates,
		})
	}

	count := 0
	duplicates := 0
	errors := 0

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Only import JPG/PNG/HEIC
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" && ext != ".png" && ext != ".heic" {
			return nil
		}

		_, err = a.manager.ImportPhoto(path)
		if err == nil {
			count++
		} else if strings.Contains(err.Error(), "duplicate photo detected") {
			duplicates++
		} else {
			errors++
			fmt.Printf("[BACKEND] Import error for %s: %v\n", path, err)
		}

		// Emit progress
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "import:progress", map[string]interface{}{
				"current":    count + duplicates + errors,
				"total":      totalCandidates,
				"imported":   count,
				"duplicates": duplicates,
				"errors":     errors,
				"lastPath":   filepath.Base(path),
			})
		}
		return nil
	})

	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "import:end", map[string]interface{}{
			"imported":   count,
			"duplicates": duplicates,
			"errors":     errors,
			"total":      totalCandidates,
		})
	}

	return count, err
}

// UpdatePhotoDate updates the capture date of a photo
func (a *App) UpdatePhotoDate(photoID int64, newDate string) error {
	parsedDate, err := time.Parse(time.RFC3339, newDate)
	if err != nil {
		// Try other formats if RFC3339 fails (e.g. from datetime-local input)
		parsedDate, err = time.Parse("2006-01-02T15:04", newDate)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
	}

	return a.manager.UpdateMetadata(photoID, "date_taken", parsedDate)
}

// LogFrontendError allows the frontend to log errors to the Go terminal
func (a *App) LogFrontendError(message string) {
	fmt.Printf("[FRONTEND ERROR] %s\n", message)
	a.uiErrors = append(a.uiErrors, fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), message))
	if len(a.uiErrors) > 100 {
		a.uiErrors = a.uiErrors[1:]
	}
}

// LogUIState allows the frontend to report its full state for debugging/automation
func (a *App) LogUIState(state string) {
	fmt.Printf("[UI STATE] %s\n", state)
	a.uiLogs = append(a.uiLogs, fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), state))
	if len(a.uiLogs) > 100 {
		a.uiLogs = a.uiLogs[1:]
	}
}

// GetDiagnostics returns system diagnostic information
func (a *App) GetDiagnostics() map[string]interface{} {
	libPath, _ := filepath.Abs("library")
	dbPath, _ := filepath.Abs("photoo.db")

	stats, _ := os.Stat(libPath)
	libExists := stats != nil

	var photoCount int
	a.db.QueryRow("SELECT COUNT(*) FROM photos").Scan(&photoCount)

	return map[string]interface{}{
		"library_path":   libPath,
		"library_exists": libExists,
		"db_path":        dbPath,
		"photo_count":    photoCount,
		"wails_context":  a.ctx != nil,
	}
}

// GetAutomationLogs returns the captured logs and errors for analysis
func (a *App) GetAutomationLogs() map[string]interface{} {
	var thumbHistory []string
	if a.thumbH != nil {
		thumbHistory = a.thumbH.History
	}
	return map[string]interface{}{
		"ui_logs":           a.uiLogs,
		"ui_errors":         a.uiErrors,
		"thumbnail_history": thumbHistory,
	}
}

// SendCommand allows the backend to send an automation command to the frontend
// Useful for driving the UI from tests or scripts
func (a *App) SendCommand(action string, payload interface{}) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "automation:command", map[string]interface{}{
		"action":  action,
		"payload": payload,
	})
}
