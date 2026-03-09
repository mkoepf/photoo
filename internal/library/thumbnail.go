package library

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type ThumbnailHandler struct {
	libraryPath string
	cachePath   string
	History     []string
	mu          sync.Mutex
	semaphore   chan struct{}
	locks       sync.Map // Map of filename -> *sync.Mutex
}

func NewThumbnailHandler(libraryPath string) *ThumbnailHandler {
	cachePath := filepath.Join(libraryPath, ".thumbnails")
	os.MkdirAll(cachePath, 0755)
	return &ThumbnailHandler{
		libraryPath: libraryPath,
		cachePath:   cachePath,
		History:     make([]string, 0),
		semaphore:   make(chan struct{}, 4), // Limit to 4 concurrent decodes
	}
}

func (h *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	h.mu.Lock()
	h.History = append(h.History, fmt.Sprintf("[%s] %s %s", time.Now().Format("15:04:05"), r.Method, path))
	if len(h.History) > 100 {
		h.History = h.History[1:]
	}
	h.mu.Unlock()

	// Support both /thumbnail/ and thumbnail/
	if !strings.HasPrefix(path, "/thumbnail/") && !strings.HasPrefix(path, "thumbnail/") {
		return
	}

	// Path parsing: /thumbnail/filename.ext or thumbnail/filename.ext -> filename.ext
	trimmed := strings.TrimLeft(path, "/")
	if !strings.HasPrefix(trimmed, "thumbnail/") {
		return
	}

	filename := strings.TrimPrefix(trimmed, "thumbnail/")
	filename = strings.TrimLeft(filename, "/")

	if filename == "" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// replace slashes with underscores for flat cache
	safeName := strings.ReplaceAll(filename, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	cacheFilename := safeName + ".thumb.jpg"
	cacheFullPath := filepath.Join(h.cachePath, cacheFilename)

	// 1. Check Cache First
	if _, err := os.Stat(cacheFullPath); err == nil {
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, cacheFullPath)
		return
	}

	// 2. Lock for this specific filename to avoid redundant generation
	actualLock, _ := h.locks.LoadOrStore(filename, &sync.Mutex{})
	lock := actualLock.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	// Re-check cache after acquiring lock (another goroutine might have finished it)
	if _, err := os.Stat(cacheFullPath); err == nil {
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, cacheFullPath)
		return
	}

	// 3. Wait for semaphore to limit total concurrent decodes
	h.semaphore <- struct{}{}
	defer func() { <-h.semaphore }()

	fullPath := filepath.Join(h.libraryPath, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	/* ext := strings.ToLower(filepath.Ext(filename)) */
	var src image.Image
	var err error

	/* HEIC support disabled on ARM64 if it fails to build
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
	*/
	src, err = imaging.Open(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to open image: %v", err), http.StatusInternalServerError)
		return
	}
	// }

	// Create a 300x300 thumbnail
	thumbnail := imaging.Fill(src, 300, 300, imaging.Center, imaging.Lanczos)

	// Save to Cache
	os.MkdirAll(h.cachePath, 0755)
	err = imaging.Save(thumbnail, cacheFullPath)
	if err != nil {
		fmt.Printf("[BACKEND] Failed to save thumbnail to cache: %v\n", err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	err = imaging.Encode(w, thumbnail, imaging.JPEG)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode thumbnail: %v", err), http.StatusInternalServerError)
	}
}
