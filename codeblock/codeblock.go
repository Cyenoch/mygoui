// Package codeblock renders native code, output and diffs. Markdown fences keep
// MyGo's maintained syntax, span and selection caches through Markdown.
package codeblock

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/transfer"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/egoist/mygo/ui/markdown"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
)

type Variant uint8

const (
	Default Variant = iota
	Destructive
)

type Size uint8

const (
	SizeDefault Size = iota
	SizeSmall
)

type LineVariant uint8

const (
	LineDefault LineVariant = iota
	Added
	Removed
	Highlighted
)

// Line contains source text without its newline or a decorative diff marker.
// In Options.Lines, joining these texts with newlines must reproduce View's
// text argument (apart from its optional final newline).
type Line struct {
	Text    string
	Variant LineVariant
}

// Options describes genuine caller metadata; an empty title invents no filename
// or language. OnDownload is the only way to expose a Download action.
type Options struct {
	Title, Footer string
	// Language is an explicit native language name or alias, such as "go".
	// Empty and unknown names render plain text; Title never selects syntax.
	Language          string
	Variant           Variant
	Size              Size
	Wrap, LineNumbers bool
	StartLine         int
	MaxHeight         float32
	// Embedded lets a containing source card own the border and corner clip.
	Embedded           bool
	Open               bool
	HideCopy, HideWrap bool
	FloatingActions    bool
	OnDownload         func(string)
	Lines              []Line
}

// State belongs to one block in one window. The zero value is ready for use.
// Wrap initializes once from Options.Wrap, then reflects the user's choice.
// ContentHandle and Scroll let a caller focus or inspect the content viewport.
type State struct {
	Wrap          bool
	ContentHandle ui.Handle
	Scroll        ui.ScrollState
	initialized   bool
	copiedAt      time.Time
	paragraphs    []uint64
	syntax        syntaxState
}

// syntaxState retains the native tokenizer's state across complete source lines.
// Only changed source retokenizes; a theme change recolors the retained tokens.
type syntaxState struct {
	source, language string
	tokenizer        editor.Tokenizer
	theme            ui.Theme
	colors           *editor.Theme
	lines            []syntaxLine
}

type syntaxLine struct {
	text   string
	end    uint64
	tokens []editor.Token
	spans  []ui.Span
}

func (s *syntaxState) update(t *ui.Theme, text, language string) {
	languageChanged := s.language != language
	if languageChanged {
		s.language, s.tokenizer = language, nil
		if native := editor.LanguageNamed(language); native != nil {
			s.tokenizer = native.Tokenizer
		}
	}
	themeChanged := s.colors == nil || s.theme != *t
	if themeChanged {
		s.theme, s.colors = *t, editor.DefaultTheme(t)
	}
	if !languageChanged && !themeChanged && s.source == text && s.lines != nil {
		return
	}
	changed := languageChanged
	state := uint64(0)
	n := 0
	for line := range strings.SplitSeq(strings.TrimSuffix(text, "\n"), "\n") {
		if n == len(s.lines) {
			s.lines = append(s.lines, syntaxLine{})
			changed = true
		}
		cached := &s.lines[n]
		changed = changed || cached.text != line
		if changed {
			cached.text, cached.end = line, 0
			cached.tokens = cached.tokens[:0]
			if s.tokenizer != nil {
				cached.tokens, cached.end = s.tokenizer.Tokenize(line, state, cached.tokens)
			}
		}
		if changed || themeChanged {
			cached.spans = cached.spans[:0]
			start := 0
			for _, token := range cached.tokens {
				end := min(token.End, len(line))
				if end > start {
					style := s.colors.Tokens[token.Kind]
					span := ui.Span{Text: line[start:end], Color: s.colors.TokenColor(token.Kind), Italic: style.Italic}
					if style.Bold {
						span.Weight = 600
					}
					cached.spans = append(cached.spans, span)
				}
				start = max(start, end)
			}
			if start < len(line) && s.tokenizer != nil {
				cached.spans = append(cached.spans, ui.Span{Text: line[start:], Color: s.colors.Text})
			}
		}
		state = cached.end
		n++
	}
	clear(s.lines[n:])
	s.source, s.lines = text, s.lines[:n]
}

func (s *State) initialize(o Options) {
	if !s.initialized {
		s.Wrap, s.initialized = o.Wrap, true
	}
}

func metrics(c *ui.Context, o Options) (padding, font, radius float32) {
	padding, font, radius = theme.Space16, theme.Space14, theme.RadiusBubble
	if o.Size == SizeSmall {
		padding, font, radius = theme.Space12, theme.Space12, theme.RadiusPanel
	}
	return padding, theme.Scale(c, font), radius
}

func blockColors(c *ui.Context, o Options) (fill, foreground, border ui.Color) {
	t := theme.Colors(c)
	fill, foreground = t.Muted, t.Foreground
	border = t.Foreground.Alpha(0.05)
	if c.Theme().Dark {
		border = t.Foreground.Alpha(0.10)
	}
	if o.Variant == Destructive {
		fill, foreground, border = t.Destructive.Alpha(0.05), t.Destructive, t.Destructive.Alpha(0.15)
		if c.Theme().Dark {
			fill, border = t.Destructive.Alpha(0.10), t.Destructive.Alpha(0.20)
		}
	}
	return fill.Over(t.Background), foreground, border
}

func surface(c *ui.Context, o Options) ui.Element {
	fill, foreground, border := blockColors(c, o)
	_, font, radius := metrics(c, o)
	box := ui.Column(c).FillWidth().MinWidth(0).Background(fill).
		TextColor(foreground).FontSize(font).Label(locales.ComponentCodeBlock())
	if !o.Embedded {
		box.Radius(radius).Clip().Border(1, border)
	}
	return box
}

func header(c *ui.Context, o Options) ui.Element {
	pad, _, _ := metrics(c, o)
	return ui.Row(c).FillWidth().MinWidth(0).Shrink(0).Gap(8).Padding(pad, pad, 0, pad).
		TextColor(theme.Colors(c).MutedForeground)
}

func title(c *ui.Context, text string, o Options) ui.Element {
	_, font, _ := metrics(c, o)
	color := theme.Colors(c).Foreground
	if o.Variant == Destructive {
		color = theme.Colors(c).Destructive
	}
	return ui.Text(c, text).Font(theme.FontFamily(c)).FontSize(font).FontWeight(500).
		MinWidth(0).Grow(1).SingleLine().TextColor(color).Unselectable()
}

func footer(c *ui.Context, text string, o Options) ui.Element {
	pad, _, _ := metrics(c, o)
	return ui.Row(c).MinWidth(0).Shrink(0).Gap(8).Padding(0, pad, pad, pad).
		TextColor(theme.Colors(c).MutedForeground).Children(func() {
		ui.Text(c, text).Unselectable()
	})
}

// content uses native focus and scrolling; wrapping changes no source text.
func content(c *ui.Context, state *State, o Options) ui.Element {
	state.initialize(o)
	var content ui.Element
	axes := ui.Vertical
	if state.Wrap {
		content = ui.Scroll(c)
	} else {
		content = ui.ScrollBoth(c)
		axes |= ui.Horizontal
	}
	pad, font, _ := metrics(c, o)
	top, bottom := pad, pad
	if o.Title != "" {
		top = pad / 2
	}
	if o.Footer != "" {
		bottom = 8
	}
	if o.MaxHeight > 0 && o.MaxHeight <= top+bottom {
		top, bottom = 0, 0
	}
	content.Direction(ui.LeftToRight).MinWidth(0).MinHeight(0).Focusable().Bind(&state.ContentHandle).
		TrackScroll(&state.Scroll).Font(theme.Mono).FontSize(font).FixedLineHeight(font*1.625).
		Padding(top, 0, bottom, 0).Margin(0, pad, pad/2, pad).Label(locales.ComponentCodeContent())
	if o.MaxHeight > 0 {
		content.MaxHeight(o.MaxHeight)
	}
	fill, _, _ := blockColors(c, o)
	mygoui.ScrollEdges(c, content, &state.Scroll, mygoui.ScrollEdgesOptions{Surface: fill, Axes: axes})
	return content
}

func renderLine(c *ui.Context, line Line, spans []ui.Span, number int, o Options) (ui.Element, uint64) {
	_, font, _ := metrics(c, o)
	t := theme.Colors(c)
	row := ui.Row(c).MinWidth(0).AlignItems(ui.Start).MinHeight(font * 1.625)
	if line.Variant == Added || line.Variant == Removed {
		row.PaddingStart(theme.Space12)
	}
	if line.Variant == Added || line.Variant == Highlighted {
		row.Background(t.Foreground.Alpha(0.05))
	}
	if line.Variant == Removed {
		row.Background(t.Destructive.Alpha(0.10))
	}
	var paragraph uint64
	row.Children(func() {
		if line.Variant == Added || line.Variant == Removed {
			marker, color, label := "+", t.MutedForeground, locales.ComponentAddedLine()
			if line.Variant == Removed {
				marker, color, label = "−", t.Destructive, locales.ComponentRemovedLine()
			}
			ui.Text(c, marker).Absolute().Left(4).Top(0).TextColor(color).Label(label).Unselectable()
		}
		if o.LineNumbers {
			ui.Text(c, strconv.Itoa(number)).MinWidth(24).Margin(0, 12, 0, 0).
				TextAlign(ui.End).TextColor(t.MutedForeground).Unselectable()
		}
		value := line.Text
		if value == "" {
			value = "\n"
		}
		var el ui.Element
		if len(spans) > 0 {
			el = ui.RichText(c.Key("source"), spans...).MinWidth(0)
		} else {
			el = ui.Text(c.Key("source"), value).MinWidth(0)
		}
		if line.Text == "" {
			el.Height(font * 1.625).Clip()
		}
		if o.Wrap {
			el.Grow(1)
		} else {
			el.NoWrap()
		}
		paragraph = el.ID()
	})
	return row, paragraph
}

func copyButton(c *ui.Context, value func() string, state *State, o Options) ui.Element {
	copied := !state.copiedAt.IsZero() && c.Now().Sub(state.copiedAt) < 2*time.Second
	label, icon := locales.ComponentCopyCode(), lucide.Copy()
	if copied {
		label, icon = locales.ComponentCopied(), lucide.Check()
	}
	return mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Icon: icon, Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconXS,
		Disabled: o.Open,
		OnClick: func() {
			c.Services().WriteClipboard(value())
			state.copiedAt = c.Now()
			c.Announce(locales.ComponentCopied())
			c.After(2 * time.Second)
		},
	})
}

func wrapButton(c *ui.Context, state *State, o Options) ui.Element {
	state.initialize(o)
	label := locales.ComponentWrapLines()
	if state.Wrap {
		label = locales.ComponentDisableLineWrap()
	}
	return mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Icon: lucide.WrapText(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconXS,
		OnClick: func() { state.Wrap = !state.Wrap },
	}).Checked(state.Wrap)
}

func actions(c *ui.Context, value func() string, state *State, o Options, floating bool) ui.Element {
	a := ui.Row(c).Gap(theme.Space4).Shrink(0).Unselectable()
	if floating {
		pad, _, _ := metrics(c, o)
		a.Absolute().Top(pad - 4).End(pad - 4)
	} else {
		a.Margin(-4, -6, -4, 0)
	}
	a.Children(func() {
		if !o.HideWrap {
			button := wrapButton(c.Key("wrap"), state, o)
			if floating {
				button.Background(theme.Colors(c).Background)
			}
		}
		if o.OnDownload != nil {
			mygoui.Button(c.Key("download"), mygoui.ButtonOptions{
				Label: locales.ComponentDownload(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonXS,
				Disabled: o.Open, OnClick: func() { o.OnDownload(value()) },
			})
		}
		if !o.HideCopy {
			button := copyButton(c.Key("copy"), value, state, o)
			if floating {
				button.Background(theme.Colors(c).Background)
			}
		}
	})
	return a
}

// View renders native source paragraphs using MyGo's language tokenizers and
// syntax theme. Numbered or decorated lines keep gutters outside selection.
// A nil state uses native local state attached to the surface's stable key.
func View(c *ui.Context, text string, state *State, o Options) ui.Element {
	box := surface(c, o)
	if state == nil {
		state = *ui.Local(box, "codeblock", func() *State { return &State{} })
	}
	state.initialize(o)
	floating := o.FloatingActions || o.Title == ""
	showActions := box.Hovered() || box.FocusWithin()
	return box.Children(func() {
		if o.Title != "" {
			header(c, o).Children(func() {
				title(c, o.Title, o)
				if !floating {
					actions(c, func() string { return text }, state, o, false)
				}
			})
		}
		viewport := content(c, state, o).Selectable()
		if o.Language != "" || o.LineNumbers || o.Lines != nil || state.syntax.lines != nil {
			viewport.CopyData(func(parts []ui.SelectedText) transfer.Data {
				return transfer.TextData(selectedSource(parts, state.paragraphs, state.syntax.lines, strings.HasSuffix(text, "\n")))
			})
		}
		viewport.Children(func() { source(c, text, state, o) })
		if o.Footer != "" {
			footer(c, o.Footer, o)
		}
		if floating {
			// Keep controls mounted for Tab navigation, revealing on focus as well
			// as hover; opacity never disables their native keyboard interaction.
			layer := actions(c, func() string { return text }, state, o, true)
			if !showActions {
				layer.Opacity(0)
			}
		}
	})
}

func source(c *ui.Context, text string, state *State, o Options) {
	if o.Language == "" && !o.LineNumbers && o.Lines == nil && state.syntax.lines == nil {
		col := ui.Column(c)
		if !state.Wrap {
			col.Shrink(0)
		}
		col.Children(func() {
			el := ui.Text(c, text)
			if !state.Wrap {
				el.NoWrap()
			}
		})
		return
	}
	state.syntax.update(c.Theme(), text, o.Language)
	state.paragraphs = state.paragraphs[:0]
	col := ui.Column(c).MinWidth(0)
	if !state.Wrap {
		col.Shrink(0)
	}
	col.Children(func() {
		start := o.StartLine
		if start == 0 {
			start = 1
		}
		o.Wrap = state.Wrap
		for i, cached := range state.syntax.lines {
			line := Line{Text: cached.text}
			if i < len(o.Lines) {
				line.Variant = o.Lines[i].Variant
			}
			_, id := renderLine(c.Key(i), line, cached.spans, start+i, o)
			state.paragraphs = append(state.paragraphs, id)
		}
	})
}

func selectedSource(parts []ui.SelectedText, ids []uint64, lines []syntaxLine, finalNewline bool) string {
	var out strings.Builder
	lineIndex := 0
	for i, part := range parts {
		for lineIndex < len(ids) && ids[lineIndex] != part.ID {
			lineIndex++
		}
		if lineIndex >= len(lines) {
			continue
		}
		if i > 0 {
			out.WriteByte('\n')
		}
		line := lines[lineIndex].text
		start, end := runeByte(line, part.Start), runeByte(line, part.End)
		out.WriteString(line[start:end])
		if lineIndex == len(lines)-1 && finalNewline && part.End >= utf8.RuneCountInString(line) {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func runeByte(text string, runeIndex int) int {
	for i := range text {
		if runeIndex == 0 {
			return i
		}
		runeIndex--
	}
	return len(text)
}

// Markdown adapts every native fenced block without reconstructing its source,
// tokens, spans or selection bookkeeping. Info supplies no inferred filename.
func Markdown(c *ui.Context, block *markdown.CodeBlockView, options *markdown.Options) {
	o := Options{Title: block.Language, Open: block.Open, Wrap: options.WrapCode}
	box := surface(c, o)
	state := *ui.Local(box, "codeblock", func() *State { return &State{} })
	state.initialize(o)
	showActions := box.Hovered() || box.FocusWithin()
	box.Children(func() {
		if o.Title != "" {
			header(c, o).Children(func() {
				title(c, o.Title, o)
				actions(c, block.Code, state, o, false)
			})
		}
		top := float32(16)
		if o.Title != "" {
			top = 8
		}
		// Body owns native cached text and its .25em top inset. Supply only
		// the remaining source inset outside that maintained body.
		if options.Theme != nil {
			top = max(0, top-options.Theme.FontSize*0.25)
		}
		initialWrap := options.WrapCode
		options.WrapCode = state.Wrap
		_, font, _ := metrics(c, o)
		// Body owns the maintained scroller, cached spans and Markdown selection.
		// Keep its ancestor selection scope: Focusable here would stop native
		// selection inheritance at the viewport (including streaming appends).
		viewport := block.Body(c).FontSize(font).FixedLineHeight(font*1.625).
			Padding(top, 0, 0, 0).Bind(&state.ContentHandle).Label("Code content")
		options.WrapCode = initialWrap
		if !state.Wrap {
			viewport.TrackScroll(&state.Scroll)
			fill, _, _ := blockColors(c, o)
			mygoui.ScrollEdges(c, viewport, &state.Scroll, mygoui.ScrollEdgesOptions{Surface: fill, Axes: ui.Horizontal})
		} else {
			state.Scroll.X, state.Scroll.Y, state.Scroll.MaxX, state.Scroll.MaxY = 0, 0, 0, 0
		}
		if o.Title == "" {
			layer := actions(c, block.Code, state, o, true)
			if !showActions {
				layer.Opacity(0)
			}
		}
	})
}
