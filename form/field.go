// Package form provides Coss-inspired form controls with native MyGo editing,
// focus, selection and accessibility. Callers own values, validation state
// and error text; the package performs no I/O.
package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// FieldOptions describes a labelled form field. Error takes precedence over
// Invalid, which uses the library's localized invalid message.
type FieldOptions struct {
	Label, Description, Error string
	Required, Disabled        bool
	Invalid                   bool
	FullWidth                 bool
}

// Field composes a native field so the label names and focuses the control and
// assistive technology reads the description and error from it. The body must
// build exactly one primary control.
func Field(c *ui.Context, o FieldOptions, body func()) ui.Element {
	label := o.Label
	if o.Required && label != "" {
		label += " *"
	}
	f := ui.Field(c, label, body).Gap(theme.Space8).MinWidth(0).
		FontSize(theme.Scale(c, theme.Label().Size)).
		FontWeight(theme.WeightNormal).
		TextColor(theme.Colors(c).Foreground)
	if o.FullWidth {
		f.FillWidth()
	}
	errorText := o.Error
	if errorText == "" && o.Invalid {
		errorText = locale.M().Invalid
	}
	if o.Description == "" && errorText == "" {
		return f
	}
	// Native Field help text uses the current theme's font size minus one.
	// Limit the temporary theme to creating those texts, preserving the
	// native field's description, error and control associations.
	original := c.Theme()
	helpTheme := *original
	helpTheme.FontSize = theme.Scale(c, theme.HelpText().Size) + 1
	c.SetTheme(&helpTheme)
	if o.Description != "" {
		f.Description(o.Description)
	}
	if errorText != "" {
		f.Error(errorText)
	}
	c.SetTheme(original)
	return f
}

// LabelOptions configures a standalone label. Target receives focus when
// clicked; the required marker is visual only.
type LabelOptions struct {
	Required bool
	Invalid  bool
	Disabled bool
	Target   *ui.Handle
}

// Label renders a 14 DIP medium label.
func Label(c *ui.Context, text string, o LabelOptions) ui.Element {
	t := theme.Colors(c)
	ink := t.Foreground
	if o.Invalid {
		ink = t.Danger
	}
	e := style.Text(ui.Text(c, text), c, theme.Label()).TextColor(ink)
	if o.Required {
		e = ui.Row(c).Gap(theme.Space2).Children(func() {
			style.Text(ui.Text(c, text), c, theme.Label()).TextColor(ink)
			style.Text(ui.Text(c, "*"), c, theme.Label()).TextColor(t.Danger)
		})
	}
	if o.Target != nil && !o.Disabled {
		// Focus through the handle, not a resolved element: a label built
		// before its control binds the target later in the same pass, and
		// Resolve only sees bindings made in the pass being built.
		e.Cursor(ui.CursorText).OnClick(func() { o.Target.Focus(c) })
	}
	return style.Disabled(e, o.Disabled)
}

// Description renders muted help text.
func Description(c *ui.Context, text string) ui.Element {
	return style.Text(ui.Text(c, text), c, theme.HelpText()).TextColor(theme.Colors(c).Muted)
}

// FieldError renders a visible validation message. Empty text omits it.
func FieldError(c *ui.Context, text string) ui.Element {
	if text == "" {
		return ui.Element{}
	}
	return style.Text(ui.Text(c, text), c, theme.HelpText()).
		TextColor(theme.Colors(c).DestructiveForeground).Role(ui.RoleStatus)
}

// ErrorMessage is the low-level error message without field layout.
func ErrorMessage(c *ui.Context, text string) ui.Element { return FieldError(c, text) }

// FieldsetOptions groups related controls.
type FieldsetOptions struct {
	Legend, Description string
	Disabled            bool
	FullWidth           bool
}

// Fieldset renders a group with a semibold legend.
func Fieldset(c *ui.Context, o FieldsetOptions, body func()) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).Gap(theme.Space24).MinWidth(0).Shrink(0).Disabled(o.Disabled)
	if o.FullWidth {
		e.FillWidth()
	}
	e.Children(func() {
		if o.Legend != "" {
			style.Text(ui.Text(c, o.Legend), c, theme.Base()).
				FontWeight(theme.WeightSemibold).TextColor(t.Foreground)
		}
		if o.Description != "" {
			Description(c, o.Description)
		}
		ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space16).Children(body)
	})
	return e
}

// FieldGroup spaces a column of fields.
func FieldGroup(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space16).Children(body)
}

// FieldActions is a footer row for form actions.
func FieldActions(c *ui.Context, body func()) ui.Element {
	return ui.Row(c).Gap(theme.Space8).AlignItems(ui.Center).Children(body)
}

// FormOptions configures a form container.
type FormOptions struct {
	Label   string
	Invalid bool
	Error   string
}

// Form is a vertical form container. Submission stays with the caller's
// submit action; the container only lays out a top-level error.
func Form(c *ui.Context, o FormOptions, body func()) ui.Element {
	message := o.Error
	if message == "" && o.Invalid {
		message = locale.M().Invalid
	}
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space24).
		Label(o.Label).Children(func() {
		if message != "" {
			FieldError(c, message)
		}
		body()
	})
}
