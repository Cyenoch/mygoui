// Command basic demonstrates a small native application using mygoui.
package main

import (
	"fmt"
	"log"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	mygo.App.SetName("HeroUI native")
	mygo.App.OnWindowAllClosed(func() { mygo.App.Quit() })
	mygo.App.WhenReady(func() {
		count := 0
		mygo.NewWindow(mygo.WindowOptions{
			Title: "HeroUI native", Width: 520, Height: 380,
			Content: ui.View(func(c *ui.Context) {
				theme.Apply(c)
				colors := theme.Colors(c)
				ui.Column(c).Fill().Padding(theme.Space24).Gap(theme.Space16).
					Background(colors.Background).TextColor(colors.Foreground).Children(func() {
					mygoui.Card(c, mygoui.CardOptions{}, func() {
						mygoui.CardHeader(c, func() {
							mygoui.CardTitle(c, "Native HeroUI components")
							mygoui.CardDescription(c, "Go state, MyGo controls, HeroUI appearance.")
						})
						mygoui.CardContent(c, func() {
							mygoui.Paragraph(c, fmt.Sprintf("Count: %d", count), false)
							mygoui.Chip(c, mygoui.ChipOptions{
								Label: "Live", Color: mygoui.ColorSuccess, Variant: mygoui.ChipSoft,
							})
						})
						mygoui.CardFooter(c, func() {
							mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Increment", OnClick: func() { count++ },
							})
							mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Reset", Variant: mygoui.ButtonSecondary,
								OnClick: func() { count = 0 },
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
