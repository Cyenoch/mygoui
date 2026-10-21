package form

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type SelectOption struct {
	Value, Label, Group string
	Disabled            bool
}

type NativeSelectOptions struct {
	ControlOptions
	Small       bool
	Placeholder string
	OnChange    func(string)
}

// NativeSelect uses the operating system's menu popup, including native
// keyboard traversal, selection, disabled items and grouped option submenus.
// The frame itself remains a themeable MyGo popup button.
func NativeSelect(c *ui.Context, value *string, o NativeSelectOptions, options []SelectOption) ui.Element {
	colors := theme.Colors(c)
	label := o.Placeholder
	for _, option := range options {
		if option.Value == *value {
			label = option.Label
			break
		}
	}
	height := float32(32)
	if o.Small {
		height = 28
	}
	e := ui.ButtonBase(c).Role(ui.RolePopUpButton).Gap(8).Padding(4, 10).
		MinHeight(max(height, theme.Metrics(c).CompactHeight)).MinWidth(0).Radius(theme.Radii().XXL).
		Background(colors.Input.Alpha(0.5)).TextColor(colors.Foreground).
		FontSize(theme.Scale(c, 14)).Justify(ui.Start).Value(label)
	configureControl(c, e, o.ControlOptions)
	inputBorder(c, e, o.ControlOptions)
	e.Children(func() {
		color := colors.Foreground
		if label == o.Placeholder && *value == "" {
			color = colors.MutedForeground
		}
		ui.Text(c, label).TextColor(color).SingleLine().Grow(1).MinWidth(0)
		ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.MoveTo(r.X+4, r.Y+6).LineTo(r.X+8, r.Y+10).LineTo(r.X+12, r.Y+6)
			p.StrokePath(&path, 1.5, colors.MutedForeground)
		})
	})
	e.Menu(func(menu *ui.Menu) {
		choose := func(menu *ui.Menu, option SelectOption) {
			if menu.Item(option.Label).Disabled(option.Disabled || o.ReadOnly).
				Checked(option.Value == *value).Chosen() && option.Value != *value {
				*value = option.Value
				if o.OnChange != nil {
					o.OnChange(*value)
				}
			}
		}
		for i, option := range options {
			if option.Group == "" {
				choose(menu, option)
				continue
			}
			seen := false
			for _, prior := range options[:i] {
				if prior.Group == option.Group {
					seen = true
					break
				}
			}
			if seen {
				continue
			}
			menu.Submenu(option.Group, func(sub *ui.Menu) {
				for _, child := range options {
					if child.Group == option.Group {
						choose(sub, child)
					}
				}
			})
		}
	})
	return e
}
