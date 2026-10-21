package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type CardOptions struct {
	Small   bool
	Spacing float32
	Label   string
}

type CardHeaderOptions struct {
	Title, Description string
	Action             func()
}

// Card owns the inset and section spacing so independently composable native
// rows and columns need no ambient component context or duplicated padding.
// Source: https://ui.shadcn.com/r/styles/base-rhea/card.json
func Card(c *myui.Context, o CardOptions, body func()) myui.Element {
	spacing := float32(theme.Space20)
	if o.Small {
		spacing = theme.Space16
	}
	if o.Spacing > 0 {
		spacing = o.Spacing
	}
	t := theme.Colors(c)
	return myui.Column(c).MinWidth(0).FillWidth().Padding(spacing).Gap(spacing).
		Radius(min(theme.Radii().XXXXL, theme.Space24)).Clip().
		Background(t.Card).TextColor(t.CardForeground).Border(theme.HairlineWidth, t.Border).
		Label(o.Label).Children(body)
}

// CardHeader puts an optional native action column beside the title and description.
func CardHeader(c *myui.Context, o CardHeaderOptions) myui.Element {
	return myui.Row(c).FillWidth().MinWidth(0).AlignItems(myui.Start).Gap(theme.Space12).Children(func() {
		myui.Column(c).Grow(1).MinWidth(0).Gap(theme.Space6).Children(func() {
			if o.Title != "" {
				CardTitle(c, o.Title)
			}
			if o.Description != "" {
				CardDescription(c, o.Description)
			}
		})
		if o.Action != nil {
			o.Action()
		}
	})
}

func CardTitle(c *myui.Context, title string) myui.Element {
	return Typography(c, title, Heading4)
}

func CardDescription(c *myui.Context, description string) myui.Element {
	return Typography(c, description, Muted)
}

func CardContent(c *myui.Context, body func()) myui.Element {
	return myui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(body)
}

func CardFooter(c *myui.Context, body func()) myui.Element {
	return myui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Space8).Children(body)
}
