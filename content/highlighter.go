package content

import (
	"strings"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// HighlighterOptions configures a read-only code block. Language is a name
// or common alias ("go", "ts", "sh"); an unknown or empty one renders the
// code plain. Title replaces the language name in the header, as a file
// name. Wrap wraps long lines instead of scrolling them sideways. NoCopy
// leaves out the copy button.
type HighlighterOptions struct {
	Code     string
	Language string
	Title    string
	Wrap     bool
	NoCopy   bool
}

// codeFontSize and codeLineHeight are the monospace metrics of code blocks
// and snippets, in DIPs before the text scale.
const (
	codeFontSize   = 13
	codeLineHeight = 20
)

type highlightKey struct {
	code, language string
	colors         syntaxColors
}

// highlighted keeps a block's spans until its code, language or theme
// colours change, so a build does not tokenize again.
type highlighted struct {
	key   highlightKey
	spans []ui.Span
	name  string
}

type highlightStateKey struct{}

// Highlighter shows code in the monospace face, colouring comments,
// strings, numbers and keywords of the languages it knows. The code is
// selectable; the copy button copies all of it. Blank lines around the code
// are trimmed.
func Highlighter(c *ui.Context, o HighlighterOptions) ui.Element {
	t := theme.Colors(c)
	code := strings.TrimRight(strings.TrimLeft(o.Code, "\r\n"), " \t\r\n")
	root := ui.Column(c).AlignSelf(ui.Stretch).MinWidth(0).Clip().Radius(theme.RadiusScale().XL).
		Background(t.Code).TextColor(t.CodeForeground).
		BorderWidth(theme.Hairline).BorderColor(t.Border).Role(ui.RoleGroup)
	hl := ui.Local(root, highlightStateKey{}, func() highlighted { return highlighted{} })
	key := highlightKey{code: code, language: o.Language, colors: syntaxColorsOf(t)}
	if hl.spans == nil || hl.key != key {
		*hl = highlighted{key: key, spans: []ui.Span{{Text: code}}}
		if s, ok := syntaxFor(o.Language); ok {
			hl.spans, hl.name = highlight(code, s, key.colors), s.name
		}
	}
	title := o.Title
	if title == "" {
		title = hl.name
	}
	if title == "" {
		title = o.Language
	}
	if title != "" {
		root.Label(title)
	}
	copyButton := func() ui.Element {
		return CopyButton(c, CopyButtonOptions{Text: code, Size: theme.ButtonSmall})
	}
	return root.Children(func() {
		if title != "" {
			highlighterHeader(c, t, title, !o.NoCopy, copyButton)
		}
		floats := title == "" && !o.NoCopy
		highlighterBody(c, hl.spans, o.Wrap, floats)
		if floats {
			// Code is laid out left to right, so its free corner is the right.
			inset := theme.Scale(c, theme.Space8)
			copyButton().Absolute().Top(inset).Right(inset)
		}
	})
}

func highlighterHeader(c *ui.Context, t theme.Palette, title string, copies bool, copyButton func() ui.Element) {
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	ui.Row(c).FillWidth().MinWidth(0).MinHeight(theme.Scale(c, 40)).
		AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(theme.Scale(c, theme.Space8)).
		Padding(0, theme.Scale(c, 10), 0, theme.Scale(c, theme.Space16)).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).
		Children(func() {
			style.Text(ui.Text(c, title), c, ts).MinWidth(0).SingleLine().Unselectable()
			if copies {
				copyButton()
			}
		})
}

// highlighterBody builds the code. reserve keeps wrapped lines clear of a
// copy button floating over the top right corner.
func highlighterBody(c *ui.Context, spans []ui.Span, wrap, reserve bool) {
	var outer ui.Element
	if wrap {
		outer = ui.Column(c)
	} else {
		outer = ui.ScrollHorizontal(c)
	}
	// Code reads left to right in a right-to-left interface too.
	outer.FillWidth().MinWidth(0).Direction(ui.LeftToRight).Children(func() {
		pad := theme.Scale(c, theme.Space16)
		right := pad
		if reserve {
			right += theme.Scale(c, theme.ButtonHeight(theme.ButtonSmall))
		}
		inner := ui.Column(c).Padding(pad, right, pad, pad)
		if !wrap {
			inner.Shrink(0)
		}
		inner.Children(func() {
			text := ui.RichText(c, spans...).Font(theme.MonoFontFamily()).
				FontSize(theme.Scale(c, codeFontSize)).
				FixedLineHeight(theme.Scale(c, codeLineHeight)).TabSize(4).Selectable()
			if !wrap {
				text.NoWrap()
			}
		})
	})
}
