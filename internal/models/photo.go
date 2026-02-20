package models

import (
	"time"
)

type Photo struct {
	ID           int64     `json:"id"`
	OriginalPath string    `json:"original_path"`
	LibraryPath  string    `json:"library_path"`
	Filename     string    `json:"filename"`
	Hash         string    `json:"hash"` // SHA-256
	DateTaken    time.Time `json:"date_taken"`
	CameraModel  string    `json:"camera_model"`
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	ImportDate   time.Time `json:"import_date"`
}

type MetadataHistory struct {
	ID        int64     `json:"id"`
	PhotoID   int64     `json:"photo_id"`
	FieldName string    `json:"field_name"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}
