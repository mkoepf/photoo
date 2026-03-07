package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"photoo/internal/library"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Get a sub-filesystem for the frontend assets
	frontendDist, _ := fs.Sub(assets, "frontend/dist")

	// Get absolute path for library to ensure it's found during dev and prod
	libDir := "library"
	if os.Getenv("PHOTOO_SELF_TEST") == "true" {
		libDir = "library_self_test"
	}

	libPath, err := filepath.Abs(libDir)
	if err != nil {
		libPath = libDir // Fallback
	}
	fmt.Printf("[BACKEND] Using library path: %s\n", libPath)

	// Create thumbnail handler
	thumbHandler := library.NewThumbnailHandler(libPath)
	app.SetThumbnailHandler(thumbHandler)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "photoo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: frontendDist,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("[ASSET] Request: %s\n", r.URL.Path)
				thumbHandler.ServeHTTP(w, r)
			}),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
