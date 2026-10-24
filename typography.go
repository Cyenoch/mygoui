package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TypographyOptions selects a HeroUI semantic role and emphasis.
type TypographyOptions struct {
	Type       theme.SemanticType
	Muted      bool
	Truncate   bool
	Align      ui.Align
	Weight     int
	Selectable bool
	Level      int
}

// Typography renders a semantic HeroUI type role. Headings expose a heading
// role and level to assistive technology; body and code remain ordinary text.
func Typography(c *ui.Context, text string, o TypographyOptions) ui.Element {
	ts := o.Type.Style()
	if o.Weight != 0 {
		ts.Weight = o.Weight
	}
	e := ui.Text(c, text)
	e = style.Text(e, c, ts)
	if o.Muted {
		e.TextColor(theme.Colors(c).Muted)
	} else {
		e.TextColor(theme.Colors(c).Foreground)
	}
	if o.Truncate {
		e.SingleLine()
	}
	if o.Selectable {
		e.Selectable()
	}
	e.TextAlign(o.Align)
	if o.Type >= theme.TypeH1 && o.Type <= theme.TypeH6 {
		level := o.Level
		if level == 0 {
			level = int(o.Type-theme.TypeH1) + 1
		}
		e.Role(ui.RoleHeading).Level(level)
	} else if o.Type == theme.TypeCode {
		// HeroUI's code role is an inline chip on the default fill.
		e.Background(theme.Colors(c).Default).Radius(theme.RadiusScale().MD).
			PaddingX(theme.Space6).PaddingY(theme.Space2).AlignSelf(ui.Start)
	}
	return e
}

// Heading renders a HeroUI heading level 1–6.
func Heading(c *ui.Context, level int, text string) ui.Element {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	return Typography(c, text, TypographyOptions{
		Type:  theme.SemanticType(int(theme.TypeH1) + level - 1),
		Level: level,
	})
}

// Paragraph renders HeroUI body copy at base, small or extra-small size.
func Paragraph(c *ui.Context, text string, small bool) ui.Element {
	t := theme.TypeBody
	if small {
		t = theme.TypeBodySM
	}
	return Typography(c, text, TypographyOptions{Type: t, Selectable: true})
}

// Code renders inline code with the HeroUI code role.
func Code(c *ui.Context, text string) ui.Element {
	return Typography(c, text, TypographyOptions{Type: theme.TypeCode, Selectable: true})
}

// Prose arranges caller-composed typography with HeroUI's prose rhythm.
func Prose(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).
		TextColor(theme.Colors(c).Foreground).Children(body)
}
