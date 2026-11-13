package content

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	mmd "github.com/egoist/mygo/ui/markdown"
)

// MarkdownOptions configures themed Markdown. Compact is the denser scale
// for messages: 14 DIP text, headings close to it and tighter spacing.
// OnLink handles a click on a link; without it MyGo opens web and mail
// links and scrolls to heading anchors.
type MarkdownOptions struct {
	Source  string
	Compact bool
	OnLink  func(url string)
}

// markdownState is what the native renderer needs kept between frames: its
// options by identity, and the theme they point at.
type markdownState struct {
	options mmd.Options
	theme   mmd.Theme
}

type markdownStateKey struct{}

// Markdown renders Source with MyGo's native Markdown view in the theme's
// type scale and colours. Code blocks render through Highlighter. Its text
// is selectable and copies as text, HTML or Markdown.
func Markdown(c *ui.Context, o MarkdownOptions) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0)
	st := *ui.Local(root, markdownStateKey{}, func() *markdownState {
		st := &markdownState{}
		st.options.Theme = &st.theme
		st.options.CodeBlock = func(c *ui.Context, b *mmd.CodeBlockView) {
			Highlighter(c, HighlighterOptions{Code: b.Code(), Language: b.Language})
		}
		return st
	})
	// The renderer compares the theme by value, so refreshing it in place
	// costs nothing while the palette and scale are unchanged.
	st.theme = markdownTheme(c, o.Compact)
	st.options.OnLink = o.OnLink
	return root.Children(func() { mmd.View(c, o.Source, &st.options) })
}

// markdownTheme maps the palette and type roles onto the renderer's theme.
// What it leaves zero keeps MyGo's default.
func markdownTheme(c *ui.Context, compact bool) mmd.Theme {
	t := theme.Colors(c)
	body := theme.Base()
	headings := [6]theme.TextStyle{theme.H1(), theme.H2(), theme.H3(), theme.H4(), theme.H5(), theme.H6()}
	spacing, top, bottom := float32(0.75), float32(1.4), float32(0.55)
	if compact {
		body = theme.SM()
		headings = [6]theme.TextStyle{theme.XL(), theme.LG(), theme.Base(), theme.SM(), theme.SM(), theme.SM()}
		spacing, top, bottom = 0.5, 0.9, 0.4
	}
	th := mmd.Theme{
		FontSize:     theme.Scale(c, body.Size),
		Font:         theme.FontFamily(c),
		CodeFont:     theme.MonoFontFamily(),
		Text:         t.Foreground,
		Muted:        t.Muted,
		Link:         t.Link,
		Strong:       t.Foreground,
		StrongWeight: theme.WeightSemibold,
		// The renderer cannot underline links; the palette's link colour is
		// the foreground, so weight is what sets them apart.
		LinkWeight:     theme.WeightMedium,
		Code:           t.Foreground,
		CodeBackground: t.Default,
		BlockSpacing:   spacing,
		QuoteBar:       t.Border,
		QuoteText:      t.Muted,
		Alerts: [5]ui.Color{
			t.InfoForeground, t.SuccessForeground, t.Primary,
			t.WarningForeground, t.DestructiveForeground,
		},
		Rule:        t.Border,
		TableBorder: t.Border,
		TableHeader: t.MutedBackground,
		Mark:        t.WarningSoft,
		Checkbox:    t.Primary,
	}
	for i, h := range headings {
		th.Headings[i] = mmd.HeadingStyle{
			Size: h.Size / body.Size, Weight: theme.WeightSemibold,
			Color: t.Foreground, Top: top, Bottom: bottom,
		}
	}
	return th
}
