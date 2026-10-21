package mygoui

import myui "github.com/egoist/mygo/ui"

// AspectRatio derives height from the native layout width, with children clipped
// to the viewport. Callers can set Width/MaxWidth on the returned Element.
// Source: https://ui.shadcn.com/r/styles/base-rhea/aspect-ratio.json
func AspectRatio(c *myui.Context, ratio float32, body func()) myui.Element {
	e := myui.Box(c).FillWidth().MinWidth(0).AspectRatio(ratio).Clip()
	if body != nil {
		e.Children(body)
	}
	return e
}
