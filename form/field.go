package form

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type FieldOrientation uint8

const (
	FieldVertical FieldOrientation = iota
	FieldHorizontal
	FieldResponsive
)

type FieldOptions struct {
	Label, Description, Error string
	Disabled                  bool
	Orientation               FieldOrientation
}

type FieldSetOptions struct {
	Legend, Description string
	Disabled, Compact   bool
}

// Field uses MyGo's maintained label association: clicking its label focuses
// or activates the real input, and assistive technology reads help and errors
// on that input. The body should build its primary input before other actions.
func Field(c *ui.Context, o FieldOptions, body func()) ui.Element {
	colors := theme.Colors(c)
	build := func() ui.Element {
		return ui.Field(c, o.Label, body).Gap(12).FillWidth().MinWidth(0).
			FontSize(theme.Scale(c, 14)).TextColor(colors.Foreground).
			Description(o.Description).Error(o.Error).Disabled(o.Disabled)
	}
	if o.Orientation == FieldHorizontal || o.Orientation == FieldResponsive {
		var field ui.Element
		ui.Form(c, func() {
			field = build()
			if o.Orientation == FieldResponsive {
				field.Wrap()
			}
		}).FillWidth().MinWidth(0)
		return field
	}
	return build()
}

func FieldSet(c *ui.Context, o FieldSetOptions, body func()) ui.Element {
	gap := float32(24)
	if o.Compact {
		gap = 12
	}
	return ui.Fieldset(c, o.Legend, func() {
		if o.Description != "" {
			Description(c, o.Description)
		}
		body()
	}).Gap(gap).FillWidth().MinWidth(0).Disabled(o.Disabled).
		FontSize(theme.Scale(c, 14)).TextColor(theme.Colors(c).Foreground)
}

func Group(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).Gap(24).FillWidth().MinWidth(0).Children(body)
}

func FieldContent(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).Gap(4).Grow(1).MinWidth(0).Children(body)
}

func FieldLegend(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).FontSize(theme.Scale(c, 16)).FontWeight(500).
		TextColor(theme.Colors(c).Foreground)
}

func FieldTitle(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).FontSize(theme.Scale(c, 14)).FontWeight(500).
		TextColor(theme.Colors(c).Foreground)
}

// Label names a target already built in this pass and requests its focus on
// click. Field also supplies native toggle activation and help/error linkage.
func Label(c *ui.Context, text string, target *ui.Handle) ui.Element {
	e := FieldTitle(c, text)
	if target != nil {
		control := c.Resolve(*target)
		if control.Valid() {
			control.Label(text)
			if control.IsDisabled() {
				e.Opacity(0.5)
				return e
			}
		}
		e.Cursor(ui.CursorText).OnClick(func() { target.Focus(c) })
	}
	return e
}

func Description(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).FontSize(theme.Scale(c, 14)).LineHeight(1.5).
		TextColor(theme.Colors(c).MutedForeground)
}

// Error omits empty messages and displays each distinct validation once.
// Use FieldOptions.Error or ControlOptions.Error to associate it with input AX.
func Error(c *ui.Context, messages ...string) ui.Element {
	text := errorText(messages)
	if text == "" {
		return ui.Element{}
	}
	return ui.Text(c, text).FontSize(theme.Scale(c, 14)).
		TextColor(theme.Colors(c).Destructive).Role(ui.RoleStatus)
}

func FieldSeparator(c *ui.Context, label string) ui.Element {
	colors := theme.Colors(c)
	return ui.Row(c).Gap(8).MinHeight(20).FillWidth().Children(func() {
		ui.Box(c).Height(1).Grow(1).Background(colors.Border)
		if label != "" {
			Description(c, label)
			ui.Box(c).Height(1).Grow(1).Background(colors.Border)
		}
	})
}
