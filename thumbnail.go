package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ThumbnailHandler struct {
	libraryPath string
}

func NewThumbnailHandler(libraryPath string) *ThumbnailHandler {
	return &ThumbnailHandler{libraryPath: libraryPath}
}

func (h *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Path should be like /thumbnail/2023-06-10_16-13-14.HEIC
	filename := strings.TrimPrefix(r.URL.Path, "/thumbnail/")
	if filename == "" {
		http.Error(w, "missing filename", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.libraryPath, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	// For now, HEIC thumbnails will fail here since imaging doesn't support them.
	// But let's at least handle JPGs and PNGs.
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".heic" {
		// Return a placeholder for HEIC for now
		http.Error(w, "HEIC thumbnails not yet supported on-the-fly", http.StatusNotImplemented)
		return
	}

	src, err := imaging.Open(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to open image: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a 300x300 thumbnail
	thumbnail := imaging.Fill(src, 300, 300, imaging.Center, imaging.Lanczos)

	// Set content type
	w.Header().Set("Content-Type", "image/jpeg")

	// Encode to writer
	err = imaging.Encode(w, thumbnail, imaging.JPEG)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode thumbnail: %v", err), http.StatusInternalServerError)
	}
}
