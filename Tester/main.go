package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	db, err := InitDB("data/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := NewApp(db)

	err = InitializeGlobalApp(db)
	if err != nil {
		log.Fatal(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Tester",
		Width:  1024,
		Height: 768,
		// MinWidth:  800,
    	// MinHeight: 600,
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			app.App,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
