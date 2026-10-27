package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// CheckboxGroupState is the caller-owned aggregate of a group's check boxes.
// Items holds the same bools the group's CheckboxOptions.Checked point at;
// the group derives its optional select-all check box from them.
type CheckboxGroupState struct {
	Items []*bool
	all   bool
}

// Checked reports whether the group has items and every one is checked.
func (s *CheckboxGroupState) Checked() bool {
	any := false
	for _, item := range s.Items {
		if item == nil {
			continue
		}
		any = true
		if !*item {
			return false
		}
	}
	return any
}

// Indeterminate reports whether some, but not all, of the items are checked.
func (s *CheckboxGroupState) Indeterminate() bool {
	on, off := false, false
	for _, item := range s.Items {
		if item == nil {
			continue
		}
		if *item {
			on = true
		} else {
			off = true
		}
	}
	return on && off
}

// Set checks or clears every non-nil item.
func (s *CheckboxGroupState) Set(checked bool) {
	for _, item := range s.Items {
		if item != nil {
			*item = checked
		}
	}
}

// CheckboxGroupOptions configures a group of related check boxes. State and
// SelectAllLabel together add a mixed-state select-all row; leave
// either empty for a plain group whose items keep their own state.
//
// Disabled disables the whole group, including body-built items, and paints
// it with the native disabled appearance. Invalid marks the group's error;
// each item must also be given Invalid for its control outline.
type CheckboxGroupOptions struct {
	Label          string
	Description    string
	Error          string
	Orientation    GroupOrientation
	Variant        ControlVariant
	Disabled       bool
	Invalid        bool
	Required       bool
	State          *CheckboxGroupState
	SelectAllLabel string
}

// CheckboxGroup renders a labelled group of check boxes built by body, in a
// vertical or horizontal layout.
func CheckboxGroup(c *ui.Context, o CheckboxGroupOptions, body func()) ui.Element {
	errText := choiceErrorText(o.Invalid, o.Error)
	disabled := o.Disabled

	root := ui.Column(c).MinWidth(0).FillWidth().Gap(theme.Space4).
		Role(ui.RoleGroup).Label(o.Label)
	if o.Description != "" {
		root.Description(o.Description)
	}
	if errText != "" {
		root.Error(errText)
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
		items := choiceItems(c, o.Orientation)
		items.Children(func() {
			if o.State != nil && o.SelectAllLabel != "" {
				st := o.State
				st.all = st.Checked() && !st.Indeterminate()
				mixed := st.Indeterminate()
				Checkbox(c, CheckboxOptions{
					Checked:       &st.all,
					Indeterminate: &mixed,
					Label:         o.SelectAllLabel,
					Invalid:       o.Invalid,
					Disabled:      disabled,
					Variant:       o.Variant,
					OnChange:      st.Set,
				})
			}
			body()
		})
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return style.Disabled(root, disabled)
}

// choiceItems builds the container a choice group's options flow in.
func choiceItems(c *ui.Context, orientation GroupOrientation) ui.Element {
	if orientation == GroupHorizontal {
		return ui.Row(c).MinWidth(0).FillWidth().Wrap().AlignItems(ui.Start).Gap(theme.Space12)
	}
	return ui.Column(c).MinWidth(0).FillWidth().Gap(theme.Space12)
}
