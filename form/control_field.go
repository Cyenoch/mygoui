package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ControlFieldKind is the control a ControlField puts in its pressable row.
type ControlFieldKind uint8

const (
	ControlFieldSwitch ControlFieldKind = iota
	ControlFieldCheckbox
)

// ControlFieldOptions configures a label, description and switch or check box
// in one pressable row.
//
// Checked is caller-owned; a click anywhere in the row, or Space while it has
// the focus, flips it. The indicator sits at the row's trailing edge.
type ControlFieldOptions struct {
	Checked     *bool
	Label       string
	Description string
	Error       string
	Kind        ControlFieldKind
	Required    bool
	Invalid     bool
	Disabled    bool
	OnChange    func(bool)
	Handle      *ui.Handle
}

// ControlField is a single native toggle whose whole row—label, description
// and indicator—is the target.
func ControlField(c *ui.Context, o ControlFieldOptions) ui.Element {
	checked := o.Checked
	if checked == nil {
		var off bool
		checked = &off
		o.Disabled = true
	}
	t := theme.Colors(c)
	disabled := o.Disabled
	errText := choiceErrorText(o.Invalid, o.Error)

	root := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	root.Children(func() {
		var row ui.Element
		if o.Kind == ControlFieldCheckbox {
			row = ui.CheckboxBase(c, checked)
		} else {
			row = ui.SwitchBase(c, checked)
		}
		row = row.Gap(theme.Space12).FillWidth().MinWidth(0).AlignItems(ui.Center).
			FocusRing(false).TextColor(t.Foreground).Label(o.Label)
		if o.Handle != nil {
			row.Bind(o.Handle)
		}
		if o.Description != "" {
			row.Description(o.Description)
		}
		if errText != "" {
			row.Error(errText)
		}
		if o.OnChange != nil && o.Checked != nil && !disabled {
			row.OnChange(func() { o.OnChange(*checked) })
		}

		row.Children(func() {
			ui.Column(c).MinWidth(0).FillWidth().Gap(theme.Space2).Children(func() {
				if o.Label != "" {
					Label(c, o.Label, LabelOptions{
						Required: o.Required, Invalid: o.Invalid,
					})
				}
				if o.Description != "" {
					Description(c, o.Description)
				}
			})
			if o.Kind == ControlFieldCheckbox {
				checkboxControl(c, checked, nil, o.Invalid, ControlPrimary, row)
			} else {
				switchIndicator(c, checked, SwitchMedium, disabled, row)
			}
		})
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return style.Disabled(root, disabled)
}
