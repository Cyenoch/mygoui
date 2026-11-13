package form

import (
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PasswordFieldOptions configures a password input with a reveal toggle.
// Revealed is the caller-owned reveal flag; nil keeps it per window, starting
// concealed. The characters are concealed again when the focus leaves the
// field, unless KeepRevealed is set.
//
// Strength rates the value as it changes: level 0 to PasswordStrengthLevels
// fills that many segments of a bar under the field, and hint is the text
// beside it. The field sets no policy; the caller decides what makes a
// password strong.
type PasswordFieldOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth                 bool
	Revealed                  *bool
	KeepRevealed              bool
	Strength                  func(value string) (level int, hint string)
	OnChange, OnSubmit        func()
}

// PasswordStrengthLevels is the number of segments of the strength bar.
const PasswordStrengthLevels = 4

type passwordRevealedKey struct{}

// PasswordField is a native password editor with a trailing toggle button
// that shows or conceals the characters. While concealed the text stays out
// of the clipboard and input methods, as the native password mode keeps it.
func PasswordField(c *ui.Context, value *string, o PasswordFieldOptions) ui.Element {
	// An error message implies the invalid state.
	o.Invalid = o.Invalid || o.Error != ""
	value = inputTextAnchor(value)
	fieldOptions := InputOptions{
		Label: o.Label, Description: o.Description, Error: o.Error,
		Required: o.Required, Disabled: o.Disabled, Invalid: o.Invalid, FullWidth: o.FullWidth,
	}
	if !inputNeedsField(fieldOptions) {
		if o.Strength == nil {
			return passwordShell(c, value, o)
		}
		column := ui.Column(c).Gap(theme.Space(1.5))
		if o.FullWidth {
			column.FillWidth()
		}
		return column.Children(func() { passwordShell(c, value, o) })
	}
	return inputField(c, fieldOptions, true, func() { passwordShell(c, value, o) })
}

// passwordShell builds the field face around the editor and its toggle, and
// the strength bar below it.
func passwordShell(c *ui.Context, value *string, o PasswordFieldOptions) ui.Element {
	shell := passwordFace(c, value, o)
	if o.Strength != nil {
		passwordStrength(c, *value, o)
	}
	return shell
}

// passwordFace is the field face around the editor and its toggle.
func passwordFace(c *ui.Context, value *string, o PasswordFieldOptions) ui.Element {
	height := theme.Scale(c, 32)
	shell := ui.Grid(c).ColumnTracks(ui.Fr(1), ui.Fixed(height)).
		AlignItems(ui.Center).MinHeight(height).Shrink(0)
	if o.FullWidth {
		shell.FillWidth()
	} else {
		shell.MinWidth(theme.Scale(c, inputDefaultWidth))
	}
	revealed := o.Revealed
	if revealed == nil {
		revealed = ui.Local(shell, passwordRevealedKey{}, func() bool { return false })
	}
	if shell.FocusLeft() && !o.KeepRevealed {
		*revealed = false
	}
	var editor ui.Element
	shell.Children(func() {
		editor = ui.TextInputBase(c, value).MinWidth(0).JustifySelf(ui.Stretch).
			PaddingStart(theme.Space(3)).PaddingEnd(theme.Space(1)).PaddingY(theme.Scale(c, 5))
		editor = fieldText(editor, c)
		if !*revealed {
			editor.Password()
		}
		if o.Placeholder != "" {
			editor.Placeholder(o.Placeholder)
		}
		if name := inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder); name != "" {
			editor.Label(name)
		}
		if o.ReadOnly {
			editor.ReadOnly(true)
		}
		if o.OnChange != nil {
			editor.OnChange(o.OnChange)
		}
		if o.OnSubmit != nil {
			editor.OnSubmit(o.OnSubmit)
		}
		passwordToggle(c, revealed, height, o.Disabled, "")
	})
	return formstyle.SurfaceFocused(shell, c, editor.Focused(), o.Invalid, o.Disabled)
}

// passwordStrength draws the strength of value as filled segments of a bar,
// warming from the danger of a weak password to the calm of a strong one,
// with the caller's hint under it.
func passwordStrength(c *ui.Context, value string, o PasswordFieldOptions) {
	level, hint := o.Strength(value)
	level = min(max(level, 0), PasswordStrengthLevels)
	t := theme.Colors(c)
	ink := t.Destructive
	switch {
	case level >= 3:
		ink = t.Success
	case level == 2:
		ink = t.Warning
	}
	ui.Column(c).Gap(theme.Space(1)).FillWidth().Children(func() {
		bar := ui.Box(c).FillWidth().Height(theme.Scale(c, 4)).Shrink(0)
		fills := make([]ui.Color, PasswordStrengthLevels)
		for i := range fills {
			target := t.Input
			if i < level {
				target = ink
			}
			fills[i] = motion.Color(bar, i, target)
		}
		gap := theme.Scale(c, theme.Space(1))
		rtl := c.Direction() == ui.RightToLeft
		bar.Draw(func(p *ui.Painter, r ui.Rect) {
			w := (r.W - gap*(PasswordStrengthLevels-1)) / PasswordStrengthLevels
			for i, fill := range fills {
				x := r.X + float32(i)*(w+gap)
				if rtl {
					x = r.X + r.W - w - float32(i)*(w+gap)
				}
				p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, fill, r.H/2)
			}
		})
		style.Text(ui.Text(c, hint), c, theme.HelpText()).TextColor(t.MutedForeground).
			SingleLine().MinHeight(theme.Scale(c, theme.HelpText().LineHeight))
	})
}

// passwordToggle is the reveal button: a toggle button assistive technology
// reads as pressed while the characters show. subject names what it reveals
// when several toggles share a window.
func passwordToggle(c *ui.Context, revealed *bool, size float32, disabled bool, subject string) ui.Element {
	t := theme.Colors(c)
	m := locale.M()
	tip := m.ShowPassword
	if *revealed {
		tip = m.HidePassword
	}
	b := ui.ToggleBase(c, revealed).Size(size, size).Center().Shrink(0).Padding(0).
		FocusRing(false).Label(strings.TrimSpace(m.ShowPassword + " " + subject)).Tooltip(tip).TextColor(t.MutedForeground)
	b = style.Disabled(b, disabled)
	rest := ui.Transparent
	if b.FocusVisible() {
		rest = t.Accent
	}
	b = style.Background(b, 0, rest, t.Accent, t.Accent)
	eyeOff := motion.Progress(b, "revealed", *revealed)
	b.Children(func() {
		// The two icons crossfade in place while the state swaps.
		iconSize := theme.Scale(c, theme.IconMD)
		ui.Box(c).Size(iconSize, iconSize).Shrink(0).Children(func() {
			if eyeOff < 1 {
				ui.Icon(c, lucide.Eye()).Size(iconSize, iconSize).Shrink(0).Opacity(1 - eyeOff)
			}
			if eyeOff > 0 {
				ui.Icon(c, lucide.EyeOff()).Size(iconSize, iconSize).Shrink(0).Absolute().Left(0).Top(0).Opacity(eyeOff)
			}
		})
	})
	return b
}
