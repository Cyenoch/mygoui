package content

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SnippetOptions configures a one-line snippet. Prefix shows before the
// text, as "$" before a command, and is left out of what is copied. Label
// names the snippet for assistive technology. NoCopy leaves out the button.
type SnippetOptions struct {
	Text   string
	Prefix string
	Label  string
	NoCopy bool
}

// Snippet shows one line of monospace text with a copy button. Text that
// does not fit ends with an ellipsis; the whole text is still copied.
func Snippet(c *ui.Context, o SnippetOptions) ui.Element {
	t := theme.Colors(c)
	ts := theme.Code()
	ts.Size = codeFontSize
	// The copy button sits closer to the edge than text does.
	end := theme.Scale(c, 10)
	if o.NoCopy {
		end = theme.Scale(c, theme.Space16)
	}
	e := ui.Row(c).MinWidth(0).MinHeight(theme.Scale(c, 44)).
		AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space8)).
		Padding(0, end, 0, theme.Scale(c, theme.Space16)).
		Radius(theme.RadiusScale().XL).Clip().Background(t.SurfaceSecondary).
		BorderWidth(theme.Hairline).BorderColor(t.Border).Role(ui.RoleGroup)
	if o.Label != "" {
		e.Label(o.Label)
	}
	return e.Children(func() {
		// Code reads left to right in a right-to-left interface too.
		line := ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center).
			Gap(theme.Scale(c, theme.Space8)).Direction(ui.LeftToRight)
		line.Children(func() {
			if o.Prefix != "" {
				style.Text(ui.Text(c, o.Prefix), c, ts).TextColor(t.Muted).
					Shrink(0).Unselectable()
			}
			style.Text(ui.Text(c, o.Text), c, ts).TextColor(t.Foreground).
				MinWidth(0).Selectable().SingleLine()
		})
		if !o.NoCopy {
			CopyButton(c, CopyButtonOptions{Text: o.Text, Size: theme.ButtonSmall})
		}
	})
}
