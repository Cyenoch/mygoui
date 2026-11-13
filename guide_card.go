package mygoui

import (
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// GuideCardOptions configures a dismissible tip. Cover is an already-decoded
// picture shown across the top at its own ratio; CoverContent builds a custom
// cover instead. OnDismiss adds the close button: the caller stops rendering
// the card.
type GuideCardOptions struct {
	Title, Description string
	Cover              ui.ImageSource
	CoverContent       func()
	Actions            func()
	OnDismiss          func()
}

// GuideCard is a card introducing a feature: cover, title, description,
// actions and a close button.
func GuideCard(c *ui.Context, o GuideCardOptions) ui.Element {
	t := theme.Colors(c)
	covered := o.CoverContent != nil || hasImage(o.Cover)
	return Card(c, CardOptions{Label: o.Title}, func() {
		switch {
		case o.CoverContent != nil:
			ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Children(o.CoverContent)
		case covered:
			w, h := ui.ImageSize(o.Cover)
			ui.Image(c, o.Cover).FillWidth().AspectRatio(w / h).Shrink(0).
				Fit(ui.Cover).Role(ui.RoleNone)
		}
		if o.Title != "" || o.Description != "" {
			header := CardHeader(c, func() {
				if o.Title != "" {
					CardTitle(c, o.Title)
				}
				if o.Description != "" {
					CardDescription(c, o.Description)
				}
			})
			if o.OnDismiss != nil && !covered {
				// Keep the title clear of the close button.
				header.PaddingEnd(theme.Space24 + theme.Scale(c, theme.Space24))
			}
		}
		if o.Actions != nil {
			footer := CardFooter(c, o.Actions)
			if o.Title != "" || o.Description != "" {
				footer.Padding(0, theme.Space24, theme.Space24)
			}
		}
		if o.OnDismiss != nil {
			closer := ui.Box(c).Absolute().Top(theme.Space8).End(theme.Space8).
				Radius(theme.RadiusScale().XL)
			if covered {
				// A cover can be any color; the button keeps its own ground.
				closer.Background(t.Surface)
			}
			closer.Children(func() {
				CloseButton(c, CloseButtonOptions{Label: locale.M().Close, OnClick: o.OnDismiss})
			})
		}
	})
}
