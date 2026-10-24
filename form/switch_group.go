package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SwitchGroupOptions configures a set of related switches. HeroUI spaces the
// group's rows by 16 DIP and the group from its label by 24 DIP.
type SwitchGroupOptions struct {
	Label       string
	Description string
	Error       string
	Orientation GroupOrientation
	Disabled    bool
	Invalid     bool
	Required    bool
	FullWidth   bool
}

// SwitchGroup lays out caller-owned switches built by body, vertically or in
// a wrapped row.
func SwitchGroup(c *ui.Context, o SwitchGroupOptions, body func()) ui.Element {
	errText := choiceErrorText(o.Invalid, o.Error)
	disabled := o.Disabled

	root := ui.Column(c).MinWidth(0).Gap(theme.Space24).
		Role(ui.RoleGroup).Label(o.Label)
	if o.FullWidth {
		root.FillWidth()
	}
	if o.Description != "" {
		root.Description(o.Description)
	}
	if errText != "" {
		root.Error(errText)
	}
	root.Children(func() {
		if o.Label != "" || o.Description != "" || errText != "" {
			ui.Column(c).MinWidth(0).Gap(theme.Space4).Children(func() {
				if o.Label != "" {
					Label(c, o.Label, LabelOptions{
						Required: o.Required, Invalid: o.Invalid,
					})
				}
				if o.Description != "" {
					Description(c, o.Description)
				}
				if errText != "" {
					FieldError(c, errText)
				}
			})
		}
		items := choiceItems(c, o.Orientation)
		items.Children(body)
	})
	return style.Disabled(root, disabled)
}
