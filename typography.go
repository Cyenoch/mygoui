package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type TypographyKind uint8

const (
	Paragraph TypographyKind = iota
	Heading1
	Heading2
	Heading3
	Heading4
	Lead
	Large
	Small
	Muted
	InlineCode
	Blockquote
)

// Typography applies native text metrics and heading roles. Layout owns measure
// and flow; text remains selectable and wraps at the available native width.
// Current official typography documentation redirects to /docs/typeset, whose
// rhythm follows surrounding text rather than a CSS class inventory.
func Typography(c *myui.Context, text string, kind TypographyKind) myui.Element {
	t := theme.Colors(c)
	style := theme.Body()
	ink := t.Foreground
	level := 0
	switch kind {
	case Heading1:
		style.Size, style.Weight, style.LineHeight = 36, 700, 1.2
		level = 1
	case Heading2:
		style.Size, style.Weight, style.LineHeight = 30, 600, 1.25
		level = 2
	case Heading3:
		style.Size, style.Weight, style.LineHeight = 24, 600, 1.3
		level = 3
	case Heading4:
		style.Size, style.Weight, style.LineHeight = 16, 500, 1.4
		level = 4
	case Lead:
		style.Size, ink = 20, t.MutedForeground
	case Large:
		style.Size, style.Weight = 18, 600
	case Small:
		style = theme.Small()
	case Muted:
		ink = t.MutedForeground
	case InlineCode:
		style = theme.Code()
	case Blockquote:
		ink = t.MutedForeground
	}
	e := myui.Text(c, text).Font(theme.FontFamily(c)).FontSize(theme.Scale(c, style.Size)).
		LineHeight(style.LineHeight).FontWeight(style.Weight).TextColor(ink).Selectable()
	if level > 0 {
		e.Role(myui.RoleHeading).Level(level)
	}
	if kind == InlineCode {
		e.Font(theme.Mono).Background(t.Muted).Padding(2, theme.Space4).Radius(theme.Radii().SM)
	}
	if kind == Blockquote {
		e.BorderWidth(0, 0, 0, 2).BorderColor(t.Border).PaddingX(theme.Space16).Italic()
	}
	return e
}
