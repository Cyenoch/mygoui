package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type SeparatorOptions struct {
	Vertical bool
	Length   float32
	Label    string
}

// Separator fills the parent's cross axis unless Length is supplied. An unnamed
// decorative rule stays out of AX; a named rule is exposed as a semantic group
// because MyGo's splitter role implies interactive resizing.
// Source: https://ui.shadcn.com/r/styles/base-rhea/separator.json
func Separator(c *myui.Context, o SeparatorOptions) myui.Element {
	e := myui.Box(c).Shrink(0).Background(theme.Colors(c).Border)
	if o.Vertical {
		e.Width(theme.HairlineWidth).AlignSelf(myui.Stretch)
		if o.Length > 0 {
			e.Height(o.Length)
		}
	} else {
		e.Height(theme.HairlineWidth).FillWidth()
		if o.Length > 0 {
			e.Width(o.Length)
		}
	}
	if o.Label != "" {
		e.Label(o.Label).Role(myui.RoleGroup)
	} else {
		e.Role(myui.RoleNone)
	}
	return e
}
