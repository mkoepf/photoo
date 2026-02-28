package main

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jdeng/goheif"
)

type ThumbnailHandler struct {
	libraryPath string
	cachePath   string
}

func NewThumbnailHandler(libraryPath string) *ThumbnailHandler {
	cachePath := filepath.Join(libraryPath, ".thumbnails")
	os.MkdirAll(cachePath, 0755)
	return &ThumbnailHandler{
		libraryPath: libraryPath,
		cachePath:   cachePath,
	}
}

func (h *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Path should be like /thumbnail/2023-06-10_16-13-14.HEIC
	filename := strings.TrimPrefix(r.URL.Path, "/thumbnail/")
	if filename == "" {
		http.Error(w, "missing filename", http.StatusBadRequest)
		return
	}

	// 1. Check Cache First
	// Thumbnails are always saved as .jpg in the cache
	cacheFilename := filename + ".thumb.jpg"
	cacheFullPath := filepath.Join(h.cachePath, cacheFilename)

	if _, err := os.Stat(cacheFullPath); err == nil {
		// Serve cached thumbnail
		http.ServeFile(w, r, cacheFullPath)
		return
	}

	// 2. Generate if not cached
	fullPath := filepath.Join(h.libraryPath, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	ext := strings.ToLower(filepath.Ext(filename))
	var src image.Image
	var err error

	if ext == ".heic" {
		file, err := os.Open(fullPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to open file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		src, err = goheif.Decode(file)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode HEIC: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		src, err = imaging.Open(fullPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to open image: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Create a 300x300 thumbnail
	thumbnail := imaging.Fill(src, 300, 300, imaging.Center, imaging.Lanczos)

	// 3. Save to Cache
	err = imaging.Save(thumbnail, cacheFullPath)
	if err != nil {
		// Log error but continue serving the generated thumbnail
		fmt.Printf("Failed to save thumbnail to cache: %v\n", err)
	}

	// Set content type
	w.Header().Set("Content-Type", "image/jpeg")

	// Encode to writer
	err = imaging.Encode(w, thumbnail, imaging.JPEG)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode thumbnail: %v", err), http.StatusInternalServerError)
	}
}
