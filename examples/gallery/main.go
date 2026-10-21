// Command gallery opens the native mygoui component gallery.
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoui/internal/gallery"
	"github.com/Cyenoch/mygoui/internal/gallery/locales"
)

func main() {
	mygo.App.SetName("mygoui Gallery")
	mygo.App.OnWindowAllClosed(func() { mygo.App.Quit() })
	mygo.App.WhenReady(func() {
		state := gallery.New()
		window := mygo.NewWindow(mygo.WindowOptions{
			Title: locales.ComponentGalleryWindowTitle(), Width: 1120, Height: 820,
			Content: ui.View(func(c *ui.Context) { gallery.View(c, state) }),
		})
		mygo.App.OnLocaleChanged(func(string) {
			window.SetTitle(locales.ComponentGalleryWindowTitle())
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
