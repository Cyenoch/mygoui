package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// RadioGroupOptions configures a labelled radio group. The group owns roving
// focus and arrow-key selection through the native radio group; each radio
// carries its own RadioGroupState.
//
// Disabled disables the whole group and paints it at half opacity. Invalid
// marks the group's error; each radio must also be given Invalid for its
// danger outline, because radio state is explicit.
type RadioGroupOptions struct {
	Label       string
	Description string
	Error       string
	Orientation GroupOrientation
	Disabled    bool
	Invalid     bool
	Required    bool
	FullWidth   bool
}

// RadioGroup is HeroUI's radio group: a labelled container of radio buttons
// that form one Tab stop, with the arrows, Home and End moving the choice.
func RadioGroup(c *ui.Context, o RadioGroupOptions, body func()) ui.Element {
	errText := choiceErrorText(o.Invalid, o.Error)
	disabled := o.Disabled

	root := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	if o.FullWidth {
		root.FillWidth()
	}
	root.Children(func() {
		if o.Label != "" {
			Label(c, o.Label, LabelOptions{
				Required: o.Required, Invalid: o.Invalid,
			})
		}
		if o.Description != "" {
			Description(c, o.Description)
		}
		group := ui.RadioGroup(c, body).Label(o.Label)
		if o.Description != "" {
			group.Description(o.Description)
		}
		if errText != "" {
			group.Error(errText)
		}
		if o.Orientation == GroupHorizontal {
			group.Row().Wrap().Gap(theme.Space16)
		} else {
			// HeroUI's first radio sits 16 DIP under the label block: the
			// group's 4 DIP line gap plus 12 DIP before the item column.
			group.Column().Gap(theme.Space16).Padding(theme.Space12, 0, 0, 0)
		}
		group.Disabled(disabled)
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return style.Disabled(root, disabled)
}
