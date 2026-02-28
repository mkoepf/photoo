package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Get a sub-filesystem for the frontend assets
	frontendDist, _ := fs.Sub(assets, "frontend/dist")

	// Create thumbnail handler
	thumbHandler := NewThumbnailHandler("library")

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "photoo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if len(r.URL.Path) >= 11 && r.URL.Path[:11] == "/thumbnail/" {
					thumbHandler.ServeHTTP(w, r)
					return
				}
				// Default handler for assets
				http.FileServer(http.FS(frontendDist)).ServeHTTP(w, r)
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
