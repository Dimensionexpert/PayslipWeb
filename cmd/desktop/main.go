package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, err := NewApp()
	if err != nil {
		fmt.Println("creating app:", err)
		return
	}

	err = wails.Run(&options.App{
		Title:  "PayslipDesktop",
		Width:  1920,
		Height: 1200,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		OnShutdown: app.shutdown,
	})

	if err != nil {
		fmt.Println("running application:", err)
	}
}
