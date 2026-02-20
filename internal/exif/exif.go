package exif

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Metadata struct {
	DateTaken   time.Time
	CameraModel string
	Latitude    *float64
	Longitude   *float64
}

// GooglePhotosMetadata represents the structure of the .json sidecar files
type GooglePhotosMetadata struct {
	PhotoTakenTime struct {
		Timestamp string `json:"timestamp"`
	} `json:"photoTakenTime"`
	GeoData struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geoData"`
}

func ExtractMetadata(path string) (*Metadata, error) {
	metadata := &Metadata{}

	// 1. Try to read from sidecar JSON (Google Photos style)
	sidecarPath := path + ".supplemental-metadata.json"
	if _, err := os.Stat(sidecarPath); err == nil {
		if sm, err := readGooglePhotosJSON(sidecarPath); err == nil {
			metadata.DateTaken = sm.DateTaken
			metadata.Latitude = sm.Latitude
			metadata.Longitude = sm.Longitude
		}
	} else {
		// Try alternative sidecar name: .jpg.suppl.json or similar
		// For simplicity, we just check the most common ones
	}

	// 2. Extract from EXIF (if possible)
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		x, err := exif.Decode(f)
		if err == nil {
			if dt, err := x.DateTime(); err == nil && metadata.DateTaken.IsZero() {
				metadata.DateTaken = dt
			}
			if model, err := x.Get(exif.Model); err == nil && model != nil {
				metadata.CameraModel = model.String()
			}
			if lat, lon, err := x.LatLong(); err == nil && metadata.Latitude == nil {
				metadata.Latitude = &lat
				metadata.Longitude = &lon
			}
		}
	}

	// 3. Fallback to file modification time
	if metadata.DateTaken.IsZero() {
		info, _ := os.Stat(path)
		metadata.DateTaken = info.ModTime()
	}

	return metadata, nil
}

func readGooglePhotosJSON(path string) (*Metadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var gp GooglePhotosMetadata
	if err := json.Unmarshal(data, &gp); err != nil {
		return nil, err
	}

	m := &Metadata{}
	// Google Photos uses Unix timestamps in seconds
	ts := gp.PhotoTakenTime.Timestamp
	var seconds int64
	fmt.Sscanf(ts, "%d", &seconds)
	if seconds > 0 {
		m.DateTaken = time.Unix(seconds, 0)
	}

	if gp.GeoData.Latitude != 0 || gp.GeoData.Longitude != 0 {
		m.Latitude = &gp.GeoData.Latitude
		m.Longitude = &gp.GeoData.Longitude
	}

	return m, nil
}
