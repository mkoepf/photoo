package main

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

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

	// Get absolute path for library to ensure it's found during dev and prod
	libPath, err := filepath.Abs("library")
	if err != nil {
		libPath = "library" // Fallback
	}

	// Create thumbnail handler
	thumbHandler := library.NewThumbnailHandler(libPath)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "photoo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("[BACKEND] Request: method=%s, path=%s, rawQuery=%s\n", r.Method, r.URL.Path, r.URL.RawQuery)
				if strings.Contains(r.URL.Path, "/thumbnail/") {
					fmt.Printf("[BACKEND] HIT /thumbnail/ -> calling handler\n")
					thumbHandler.ServeHTTP(w, r)
					return
				}
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
