package form

import (
	"strings"

	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ControlOptions configures the actual input, including its accessible name.
// Use Field to show and connect a visible label, description and error.
type ControlOptions struct {
	Label, Description, Error            string
	Disabled, ReadOnly, Invalid, Loading bool
	Handle                               *ui.Handle
}

type InputOptions struct {
	ControlOptions
	Placeholder string
	Password    bool
	OnChange    func(string)
	OnSubmit    func(string)
}

type TextareaOptions struct {
	InputOptions
	MinLines, MaxLines int
}

func Input(c *ui.Context, value *string, o InputOptions) ui.Element {
	e := ui.TextInput(c, value)
	configureInput(c, e, o, false)
	if o.Password {
		e.Password()
	}
	inputActions(e, value, o)
	return e
}

func Textarea(c *ui.Context, value *string, o TextareaOptions) ui.Element {
	e := ui.TextArea(c, value)
	configureInput(c, e, o.InputOptions, true)
	if o.MinLines > 0 {
		e.Lines(o.MinLines, max(o.MinLines, o.MaxLines))
	}
	inputActions(e, value, o.InputOptions)
	return e
}

func inputActions(e ui.Element, value *string, o InputOptions) {
	if o.ReadOnly {
		return
	}
	if o.OnChange != nil {
		e.OnChange(func() { o.OnChange(*value) })
	}
	if o.OnSubmit != nil {
		e.OnSubmit(func() { o.OnSubmit(*value) })
	}
}

func configureInput(c *ui.Context, e ui.Element, o InputOptions, multiline bool) {
	colors := theme.Colors(c)
	e.FontSize(theme.Scale(c, 14)).TextColor(colors.Foreground).
		Placeholder(o.Placeholder).ReadOnly(o.ReadOnly).
		Padding(4, theme.Space10).MinHeight(theme.Metrics(c).InputHeight).
		MinWidth(0).FillWidth().Radius(theme.Radii().XXL).Background(colors.Input.Alpha(0.5))
	if multiline {
		e.PaddingY(8).MinHeight(max(float32(64), theme.Scale(c, 40)+16))
	}
	configureControl(c, e, o.ControlOptions)
	inputBorder(c, e, o.ControlOptions)
}

func configureControl(c *ui.Context, e ui.Element, o ControlOptions) {
	e.Disabled(o.Disabled || o.Loading)
	if o.Label != "" {
		e.Label(o.Label)
	}
	if o.Description != "" {
		e.Description(o.Description)
	}
	if o.Error != "" {
		e.Error(o.Error)
	} else if o.Invalid {
		e.Error(locales.ComponentInvalidValue())
	}
	if o.Loading {
		e.Description(locales.ComponentLoading())
	}
	if o.ReadOnly {
		e.Description(locales.ComponentReadOnly())
	}
	if o.Handle != nil {
		e.Bind(o.Handle)
	}
}

func inputBorder(c *ui.Context, e ui.Element, o ControlOptions) {
	colors := theme.Colors(c)
	border := colors.Input.Alpha(0)
	invalid := o.Invalid || o.Error != ""
	if e.FocusVisible() || e.FocusWithin() {
		border = colors.Ring
	}
	if invalid {
		border = colors.Destructive
	}
	e.Border(1, border)
	if invalid {
		e.Shadow(0, 0, 0, 3, colors.Destructive.Alpha(0.2))
	}
}

func errorText(messages []string) string {
	var text strings.Builder
	for i, message := range messages {
		if message == "" {
			continue
		}
		duplicate := false
		for _, prior := range messages[:i] {
			if message == prior {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}
		if text.Len() > 0 {
			text.WriteByte('\n')
		}
		text.WriteString(message)
	}
	return text.String()
}
