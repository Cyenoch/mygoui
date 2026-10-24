// Command gallery opens the native MyGo UI component workbench.
package main

import (
	"log"

	"github.com/Cyenoch/mygoui/internal/gallery"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	mygo.App.SetName("MyGo UI Gallery")
	mygo.App.OnWindowAllClosed(func() { mygo.App.Quit() })
	mygo.App.WhenReady(func() {
		state := gallery.New()
		mygo.NewWindow(mygo.WindowOptions{
			Title: "MyGo UI Gallery", Width: 1180, Height: 860,
			Content: ui.View(func(c *ui.Context) { gallery.View(c, state) }),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
