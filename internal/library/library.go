package library

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"photoo/internal/exif"
	"photoo/internal/models"
)

type Manager struct {
	LibraryPath string
	DB          *sql.DB
}

func NewManager(libraryPath string, db *sql.DB) (*Manager, error) {
	if err := os.MkdirAll(libraryPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create library directory: %w", err)
	}
	return &Manager{LibraryPath: libraryPath, DB: db}, nil
}

func (m *Manager) ImportPhoto(sourcePath string) (*models.Photo, error) {
	// 1. Calculate Hash
	hash, err := calculateHash(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	// 2. Check for Duplicates in DB
	var existingID int64
	err = m.DB.QueryRow("SELECT id FROM photos WHERE hash = ?", hash).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("duplicate photo detected (hash: %s)", hash)
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	// 3. Extract Metadata (checks sidecars)
	metadata, err := exif.ExtractMetadata(sourcePath)
	if err != nil {
		metadata = &exif.Metadata{}
		info, _ := os.Stat(sourcePath)
		metadata.DateTaken = info.ModTime()
	}

	// 4. Determine Filename (YYYY-MM-DD_HH-mm-ss)
	ext := filepath.Ext(sourcePath)
	baseFilename := metadata.DateTaken.Format("2006-01-02_15-04-05")
	finalFilename, err := m.findUniqueFilename(baseFilename, ext)
	if err != nil {
		return nil, fmt.Errorf("failed to determine unique filename: %w", err)
	}

	// 5. Copy File to Library
	libraryPath := filepath.Join(m.LibraryPath, finalFilename)
	if err := copyFile(sourcePath, libraryPath); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 6. Save to Database
	photo := &models.Photo{
		OriginalPath: sourcePath,
		LibraryPath:  libraryPath,
		Filename:     finalFilename,
		Hash:         hash,
		DateTaken:    metadata.DateTaken,
		CameraModel:  metadata.CameraModel,
		ImportDate:   time.Now(),
	}

	if metadata.Latitude != nil {
		photo.Latitude = metadata.Latitude
	}
	if metadata.Longitude != nil {
		photo.Longitude = metadata.Longitude
	}

	res, err := m.DB.Exec(
		"INSERT INTO photos (original_path, library_path, filename, hash, date_taken, camera_model, latitude, longitude, import_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		photo.OriginalPath, photo.LibraryPath, photo.Filename, photo.Hash, photo.DateTaken, photo.CameraModel, photo.Latitude, photo.Longitude, photo.ImportDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save photo to database: %w", err)
	}

	id, _ := res.LastInsertId()
	photo.ID = id

	return photo, nil
}

func (m *Manager) UpdateMetadata(photoID int64, field string, newValue interface{}) error {
	// 1. Get current value
	var oldValue string
	query := fmt.Sprintf("SELECT %s FROM photos WHERE id = ?", field)
	err := m.DB.QueryRow(query, photoID).Scan(&oldValue)
	if err != nil {
		return fmt.Errorf("failed to get old value: %w", err)
	}

	// 2. Log in metadata_history
	_, err = m.DB.Exec(
		"INSERT INTO metadata_history (photo_id, field_name, old_value, new_value) VALUES (?, ?, ?, ?)",
		photoID, field, oldValue, fmt.Sprintf("%v", newValue),
	)
	if err != nil {
		return fmt.Errorf("failed to log metadata history: %w", err)
	}

	// 3. Update DB
	updateQuery := fmt.Sprintf("UPDATE photos SET %s = ? WHERE id = ?", field)
	_, err = m.DB.Exec(updateQuery, newValue, photoID)
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	// TODO: Phase 3 - Write back to file EXIF

	return nil
}

func (m *Manager) findUniqueFilename(base, ext string) (string, error) {
	filename := base + ext
	counter := 1
	for {
		path := filepath.Join(m.LibraryPath, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return filename, nil
		}
		filename = fmt.Sprintf("%s_%d%s", base, counter, ext)
		counter++
		if counter > 1000 {
			return "", fmt.Errorf("too many filename collisions for %s", base)
		}
	}
}

func calculateHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
