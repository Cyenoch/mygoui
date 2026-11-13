package mygoui

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// HighlightOptions configures marked text. Every occurrence of each query
// in Text is marked; queries are trimmed and blank ones ignored. IgnoreCase
// matches letters whatever their case.
type HighlightOptions struct {
	Text       string
	Queries    []string
	IgnoreCase bool
}

// Highlight shows Text as one paragraph with the matches of Queries marked,
// as the hits of a search. Matches are found rune by rune, so text in any
// script keeps its characters whole, and overlapping or touching matches
// share one mark. The matches are found again only when the options change.
// Style the returned paragraph as any text.
func Highlight(c *ui.Context, o HighlightOptions) ui.Element {
	t := theme.Colors(c)
	mark := statusColors(&t, ColorWarning)
	root := style.Text(ui.Text(c, ""), c, theme.SM()).TextColor(t.Foreground).MinWidth(0)
	look := highlightLook{fill: mark.Soft, ink: mark.SoftForeground, radius: theme.RadiusScale().SM}
	spans := ui.Local(root, "spans", func() highlightSpans { return highlightSpans{} })
	if !spans.current(o, look) {
		spans.rebuild(o, look)
	}
	return root.Children(func() { ui.RichText(c, spans.spans...) })
}

type highlightLook struct {
	fill, ink ui.Color
	radius    float32
}

// highlightSpans keeps the paragraph's spans with the options that made them.
type highlightSpans struct {
	built      bool
	text       string
	queries    []string
	ignoreCase bool
	look       highlightLook
	spans      []ui.Span
}

func (s *highlightSpans) current(o HighlightOptions, look highlightLook) bool {
	return s.built && s.text == o.Text && s.ignoreCase == o.IgnoreCase && s.look == look &&
		slices.Equal(s.queries, o.Queries)
}

func (s *highlightSpans) rebuild(o HighlightOptions, look highlightLook) {
	s.built, s.text, s.ignoreCase, s.look = true, o.Text, o.IgnoreCase, look
	// The caller may reuse its slice for the next search.
	s.queries = append(s.queries[:0], o.Queries...)
	s.spans = s.spans[:0]
	at := 0
	plain := func(end int) {
		if end > at {
			s.spans = append(s.spans, ui.Span{Text: o.Text[at:end]})
		}
	}
	for _, r := range highlightRanges(o.Text, o.Queries, o.IgnoreCase) {
		plain(r.start)
		s.spans = append(s.spans, ui.Span{
			Text: o.Text[r.start:r.end], Background: look.fill, Color: look.ink, Radius: look.radius,
		})
		at = r.end
	}
	plain(len(o.Text))
}

// highlightRange is a marked run of text, in bytes.
type highlightRange struct{ start, end int }

// highlightRanges returns the runs of text matched by a query, in order,
// with the runs that overlap or touch joined.
func highlightRanges(text string, queries []string, ignoreCase bool) []highlightRange {
	wanted := make([]string, 0, len(queries))
	for _, q := range queries {
		if q = strings.TrimSpace(q); q != "" {
			wanted = append(wanted, q)
		}
	}
	if len(wanted) == 0 {
		return nil
	}
	var ranges []highlightRange
	for i := 0; i < len(text); {
		end := i
		for _, q := range wanted {
			if n, ok := highlightPrefix(text[i:], q, ignoreCase); ok {
				end = max(end, i+n)
			}
		}
		if last := len(ranges) - 1; end > i && last >= 0 && ranges[last].end >= i {
			ranges[last].end = max(ranges[last].end, end)
		} else if end > i {
			ranges = append(ranges, highlightRange{i, end})
		}
		_, size := utf8.DecodeRuneInString(text[i:])
		i += size
	}
	return ranges
}

// highlightPrefix reports how many bytes of text the query matches at its
// start. Folding compares runes, since a letter and its other case may
// differ in length.
func highlightPrefix(text, query string, ignoreCase bool) (int, bool) {
	if !ignoreCase {
		return len(query), strings.HasPrefix(text, query)
	}
	n := 0
	for _, want := range query {
		got, size := utf8.DecodeRuneInString(text[n:])
		if size == 0 || !sameLetter(got, want) {
			return 0, false
		}
		n += size
	}
	return n, true
}

func sameLetter(a, b rune) bool {
	if a == b {
		return true
	}
	for r := unicode.SimpleFold(a); r != a; r = unicode.SimpleFold(r) {
		if r == b {
			return true
		}
	}
	return false
}
