// Command basic demonstrates a small native application using mygoui.
package main

import (
	"fmt"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
)

func main() {
	mygo.App.SetName("Component example")
	mygo.App.OnWindowAllClosed(func() { mygo.App.Quit() })
	mygo.App.WhenReady(func() {
		count := 0
		mygo.NewWindow(mygo.WindowOptions{
			Title: "Component example", Width: 480, Height: 320,
			Content: ui.View(func(c *ui.Context) {
				theme.Apply(c)
				colors := theme.Colors(c)
				ui.Column(c).Fill().Padding(24).Gap(16).
					Background(colors.Background).TextColor(colors.Foreground).Children(func() {
					mygoui.Card(c, mygoui.CardOptions{}, func() {
						mygoui.CardHeader(c, mygoui.CardHeaderOptions{
							Title: "Native components", Description: "Go state, MyGo controls, shadcn style.",
						})
						mygoui.CardContent(c, func() {
							ui.Text(c, fmt.Sprintf("Count: %d", count))
						})
						mygoui.CardFooter(c, func() {
							mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Increment", OnClick: func() { count++ },
							})
						})
					})
				})
			}),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
