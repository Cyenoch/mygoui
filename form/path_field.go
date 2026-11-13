package form

import (
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PathFieldOptions configures a path field. OnBrowse runs when the trailing
// button is activated; the field opens no dialog itself, so the caller shows
// one (see PickFiles) and writes the chosen path to the value. A nil OnBrowse
// disables the button. BrowseLabel replaces the localized "Browse".
type PathFieldOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	BrowseLabel               string
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth                 bool
	OnBrowse                  func()
	OnChange, OnSubmit        func()
}

// PathField is a native editor for a file-system path with a trailing browse
// button. A path too long for the field shows its start and end around an
// ellipsis, and the whole path in a tooltip, until the editor takes the
// focus, which shows the editable text. A file or folder dropped on the field
// sets its path, the first one when several are dropped, and reports it to
// OnChange; a read-only or disabled field takes no drop.
func PathField(c *ui.Context, value *string, o PathFieldOptions) ui.Element {
	// An error message implies the invalid state.
	o.Invalid = o.Invalid || o.Error != ""
	value = inputTextAnchor(value)
	fieldOptions := InputOptions{
		Label: o.Label, Description: o.Description, Error: o.Error,
		Required: o.Required, Disabled: o.Disabled, Invalid: o.Invalid, FullWidth: o.FullWidth,
	}
	if !inputNeedsField(fieldOptions) {
		return pathShell(c, value, o)
	}
	return inputField(c, fieldOptions, true, func() { pathShell(c, value, o) })
}

// pathShell builds the field face around the editor and the browse button.
func pathShell(c *ui.Context, value *string, o PathFieldOptions) ui.Element {
	height := theme.Scale(c, 32)
	shell := ui.Grid(c).ColumnTracks(ui.Fr(1), ui.FitContent()).
		AlignItems(ui.Center).MinHeight(height).Shrink(0)
	if o.FullWidth {
		shell.FillWidth()
	} else {
		shell.MinWidth(theme.Scale(c, 280))
	}
	// A drop sets the path before the editor is built from it.
	over := false
	if !o.ReadOnly && !o.Disabled {
		if files := shell.DroppedFiles(); len(files) > 0 {
			*value = files[0]
			if o.OnChange != nil {
				o.OnChange()
			}
		}
		over = shell.FileDragOver()
	}
	var editor ui.Element
	shell.Children(func() {
		ui.Box(c).MinWidth(0).JustifySelf(ui.Stretch).Children(func() {
			editor = pathEditor(c, value, o)
		})
		pathBrowseButton(c, o, height)
	})
	return formstyle.SurfaceFocused(shell, c, editor.Focused() || over, o.Invalid, o.Disabled)
}

// pathEditor is the native editor. The editor cuts overflowing text at its
// end only, so while it rests with a value its text is hidden under a label
// that cuts the same text in the middle; the editor still owns focus, input
// and the accessible value.
func pathEditor(c *ui.Context, value *string, o PathFieldOptions) ui.Element {
	padX, padY := float32(theme.Space(3)), theme.Scale(c, 5)
	e := ui.TextInputBase(c, value).FillWidth().MinWidth(0).PaddingX(padX).PaddingY(padY)
	e = fieldText(e, c)
	if o.Placeholder != "" {
		e.Placeholder(o.Placeholder)
	}
	if name := inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder); name != "" {
		e.Label(name)
	}
	if o.ReadOnly {
		e.ReadOnly(true)
	}
	if o.OnChange != nil {
		e.OnChange(o.OnChange)
	}
	if o.OnSubmit != nil {
		e.OnSubmit(o.OnSubmit)
	}
	if e.Focused() || *value == "" {
		return e
	}
	e.TextColor(ui.Transparent)
	if pathTruncated(c, e, *value, padX) {
		e.Tooltip(*value)
	}
	resting := fieldText(ui.Text(c, *value), c).SingleLine().TruncateMiddle().
		Absolute().Start(padX).End(padX).Top(padY).IgnoreInput(true)
	// The editor already reads the path to assistive technology.
	resting.Role(ui.RoleNone)
	return e
}

// pathBrowseButton is the trailing action, divided from the editor by a
// hairline as part of the same field.
func pathBrowseButton(c *ui.Context, o PathFieldOptions, height float32) {
	t := theme.Colors(c)
	label := o.BrowseLabel
	if label == "" {
		label = locale.M().Browse
	}
	b := ui.ButtonBase(c).Height(height).Center().Shrink(0).PaddingX(theme.Space(3)).
		Radius(0).FocusRing(false)
	b = style.Disabled(b, o.Disabled || o.OnBrowse == nil)
	rest := ui.Transparent
	if b.FocusVisible() {
		rest = t.Accent
	}
	b = style.Background(b, 0, rest, t.Accent, t.Accent)
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		x := r.X
		if c.Direction() == ui.RightToLeft {
			x = r.X + r.W - theme.Hairline
		}
		p.Fill(ui.Rect{X: x, Y: r.Y, W: theme.Hairline, H: r.H}, t.Input, 0)
	})
	b.Children(func() {
		style.Text(ui.Text(c, label), c, theme.SM()).FontWeight(theme.WeightMedium).
			TextColor(t.Foreground).SingleLine()
	})
	if o.OnBrowse != nil && !o.Disabled {
		b.OnClick(o.OnBrowse)
	}
}

// pathTruncated reports whether path, as the editor showed it in the last
// frame, did not fit between its paddings.
func pathTruncated(c *ui.Context, e ui.Element, path string, padX float32) bool {
	room := e.Bounds().W - 2*padX
	if room <= 0 {
		return false
	}
	width, _ := c.MeasureText(0, ui.Span{Text: path, Size: theme.Scale(c, theme.SM().Size)})
	return width > room
}
