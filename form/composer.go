package form

import (
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ComposerOptions configures a multiline draft and its submit action. The
// caller owns the draft and decides whether to clear it in OnSubmit.
type ComposerOptions struct {
	Label, Placeholder, SubmitLabel string
	Disabled, ReadOnly              bool
	SubmitDisabled, AllowEmpty      bool
	DisableEnterSubmit              bool
	MinLines, MaxLines              int
	Actions                         func()
	OnChange                        func()
	OnSubmit                        func(string)
}

// Composer builds a Coss-style plain-text draft editor. Enter submits and
// Shift+Enter inserts a newline; the caller clears the draft in OnSubmit.
// Coss's rich tokens, suggestion menus and adaptive layout are not provided.
func Composer(c *ui.Context, draft *string, o ComposerOptions) ui.Element {
	if draft == nil {
		panic("form: Composer requires a caller-owned draft")
	}
	t := theme.Colors(c)
	minLines, maxLines := o.MinLines, o.MaxLines
	if minLines <= 0 {
		minLines = 2
	}
	if maxLines < minLines {
		maxLines = max(minLines, 8)
	}
	label := o.Label
	if label == "" {
		label = o.Placeholder
	}
	canSubmit := func(text string) bool {
		return !o.Disabled && !o.ReadOnly && !o.SubmitDisabled && o.OnSubmit != nil &&
			(o.AllowEmpty || strings.TrimSpace(text) != "")
	}
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4).
		Padding(theme.Space8).Radius(theme.RadiusScale().XL)
	root = fieldSurface(root, c, ControlPrimary, false, o.Disabled)
	root.Radius(theme.RadiusScale().XL)
	root.Children(func() {
		editor := ui.TextAreaBase(c, draft).FillWidth().MinWidth(0).
			Lines(minLines, maxLines).Padding(theme.Space4).Label(label).
			Placeholder(o.Placeholder).ReadOnly(o.ReadOnly)
		fieldText(editor, c)
		if o.OnChange != nil {
			editor.OnChange(o.OnChange)
		}
		if !o.DisableEnterSubmit {
			editor.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyEnter || k.Mods != 0 {
					return false
				}
				text := k.Text()
				if canSubmit(text) {
					*draft = text
					o.OnSubmit(text)
					if *draft != text {
						k.SetText(*draft)
					}
				}
				return true
			})
		}
		ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			if o.Actions != nil {
				o.Actions()
			}
			ui.Spacer(c)
			submitLabel := o.SubmitLabel
			if submitLabel == "" {
				submitLabel = "Submit"
			}
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: submitLabel, Icon: lucide.ArrowUp(), IconOnly: true,
				Disabled: !canSubmit(*draft),
				OnClick: func() {
					if canSubmit(*draft) {
						o.OnSubmit(*draft)
					}
				},
			}).TextColor(t.PrimaryForeground).KeepFocus()
		})
	})
	return root
}
