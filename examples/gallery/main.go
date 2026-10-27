// Command gallery opens the native MyGo UI component workbench.
package main

import (
	"log"
	"runtime"

	"github.com/Cyenoch/mygoui/fonts"
	"github.com/Cyenoch/mygoui/internal/gallery"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	if err := fonts.Register(); err != nil {
		log.Fatal(err)
	}
	mygo.App.SetName("Coss UI Gallery")
	mygo.App.OnWindowAllClosed(func() { mygo.App.Quit() })
	mygo.App.WhenReady(func() {
		state := gallery.New()
		window := mygo.NewWindow(windowOptions(state))
		state.Window = window
		state.OnAppearanceChange = func(dark bool) {
			source := mygo.ThemeLight
			if dark {
				source = mygo.ThemeDark
			}
			window.SetWindowButtonsTheme(source)
		}
		state.OnAppearanceChange(state.Dark)
		window.Invalidate()
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func windowOptions(state *gallery.State) mygo.WindowOptions {
	o := mygo.WindowOptions{
		Title: "Coss UI Gallery", Width: 1180, Height: 860,
		TitleBarStyle:  mygo.TitleBarHiddenInset,
		TitleBarHeight: gallery.TitleBarHeight,
		Content:        ui.View(func(c *ui.Context) { gallery.View(c, state) }),
	}
	if runtime.GOOS == "darwin" {
		// A standard 14 DIP button is centered in the 52 DIP title bar.
		o.TrafficLightPosition = &mygo.Point{X: 16, Y: 19}
	}
	return o
}
