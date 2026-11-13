package content

import (
	"strings"
	"unicode/utf8"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// jsonMatches are the nodes of a document whose key or scalar holds a query,
// found again only when the document or the query changes.
type jsonMatches struct {
	doc     *jsonDoc
	query   string // trimmed
	nodes   []int
	current int // index in nodes; -1 for none
}

// syncMatches finds the matches of Query. A new query chooses the first and
// shows it; a new document under the same query keeps the match the user is
// on, clamped to the new count, and leaves the selection and the scroll be.
func (s *JSONViewerState) syncMatches() {
	f := &s.found
	query := strings.TrimSpace(s.Query)
	if s.doc == nil || f.doc == s.doc && f.query == query {
		return
	}
	queryChanged := f.query != query
	was := f.current
	f.doc, f.query, f.nodes, f.current = s.doc, query, f.nodes[:0], -1
	if query == "" {
		return
	}
	lower := strings.ToLower(query)
	for i := range s.doc.nodes {
		if s.doc.matches(i, lower) {
			f.nodes = append(f.nodes, i)
		}
	}
	if len(f.nodes) == 0 {
		return
	}
	if queryChanged {
		f.current = 0
		s.reveal(f.nodes[0])
		return
	}
	f.current = min(max(was, 0), len(f.nodes)-1)
}

// matches reports whether node i's key, or its text when it is a scalar,
// holds lower.
func (d *jsonDoc) matches(i int, lower string) bool {
	n := &d.nodes[i]
	if containsFold(n.key, lower) {
		return true
	}
	return n.kind != jsonObject && n.kind != jsonArray && containsFold(d.source[n.from:n.to], lower)
}

// step chooses the match d places after the current one, wrapping around.
func (s *JSONViewerState) step(d int) {
	f := &s.found
	if len(f.nodes) == 0 {
		return
	}
	f.current = (f.current + d + len(f.nodes)) % len(f.nodes)
	s.reveal(f.nodes[f.current])
}

// reveal opens the nodes above node i, chooses its row and scrolls it into
// the middle of the tree.
func (s *JSONViewerState) reveal(i int) {
	nodes := s.doc.nodes
	for p := nodes[i].parent; p >= 0; p = nodes[p].parent {
		s.outline.Open.Add(p)
	}
	// Nodes are in document order, so the row of i is the number of nodes
	// before it that show.
	shown := make([]bool, i+1)
	row := 0
	for j := 0; j <= i; j++ {
		p := nodes[j].parent
		shown[j] = p < 0 || shown[p] && s.outline.Open.Has(p)
		if shown[j] && j < i {
			row++
		}
	}
	s.selected = row
	// A long string wraps over many lines once chosen: its top shows, with
	// its key.
	align := ui.Align(ui.Center)
	if n := &nodes[i]; n.kind == jsonString && n.to-n.from > 200 {
		align = ui.Start
	}
	s.outline.List.ScrollTo(row, align)
}

// jsonToolbar is the search field, the number of matches and the buttons
// stepping through them.
func jsonToolbar(c *ui.Context, t theme.Palette, st *JSONViewerState) {
	m := locale.M()
	f := &st.found
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(2))).
		Padding(theme.Scale(c, theme.Space(1.5)), theme.Scale(c, theme.Space(2))).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).
		Children(func() {
			searchField(c, t, &st.Query, &st.search, func(k ui.TextKey) bool {
				if k.Key != ui.KeyEnter {
					return false
				}
				if k.Mods&ui.Shift != 0 {
					st.step(-1)
				} else {
					st.step(1)
				}
				return true
			})
			// A button keeps its own alignment, so the row around them
			// centers them against the field.
			ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(1))).Children(func() {
				if f.query != "" {
					ts := theme.Code()
					ts.Size, ts.LineHeight = codeFontSize-1, codeLineHeight
					at := 0
					if f.current >= 0 {
						at = f.current + 1
					}
					style.Text(ui.Text(c.Key("count"), groupDigits(at)+" / "+groupDigits(len(f.nodes))), c, ts).
						TextColor(t.Muted).Shrink(0).MarginEnd(theme.Scale(c, theme.Space(1))).Unselectable().
						Role(ui.RoleStatus)
				}
				step := func(label string, icon *ui.SVG, d int) {
					mygoui.Button(c.Key(label), mygoui.ButtonOptions{
						Label: label, Icon: icon, IconOnly: true, Variant: mygoui.ButtonGhost,
						Size: theme.ButtonSmall, Disabled: len(f.nodes) == 0, OnClick: func() { st.step(d) },
					})
				}
				step(m.Previous, lucide.ChevronUp(), -1)
				step(m.Next, lucide.ChevronDown(), 1)
			})
		})
}

// jsonPreviewLimit is about how many bytes of members a closed container's
// preview shows before it ends in an ellipsis.
const jsonPreviewLimit = 72

// jsonBriefRunes is how many characters of a string a preview shows.
const jsonBriefRunes = 24

// preview writes a non-empty container's members on one line, as
// {name: "Ada", tags: […]}.
func (d *jsonDoc) preview(i int) string {
	n := &d.nodes[i]
	open, end := byte('{'), byte('}')
	if n.kind == jsonArray {
		open, end = '[', ']'
	}
	var b strings.Builder
	b.WriteByte(open)
	for k, kid := range n.kids {
		if k > 0 {
			if b.Len() >= jsonPreviewLimit {
				b.WriteString(", …")
				break
			}
			b.WriteString(", ")
		}
		if n.kind == jsonObject {
			b.WriteString(d.nodes[kid].key)
			b.WriteString(": ")
		}
		b.WriteString(d.brief(kid))
	}
	b.WriteByte(end)
	return b.String()
}

// brief is a member in a preview: a scalar's text, a long string cut, and
// a container as {…} or [], whichever it holds.
func (d *jsonDoc) brief(i int) string {
	n := &d.nodes[i]
	switch n.kind {
	case jsonObject, jsonArray:
		open, end := "{", "}"
		if n.kind == jsonArray {
			open, end = "[", "]"
		}
		if len(n.kids) == 0 {
			return open + end
		}
		return open + "…" + end
	case jsonString:
		text := d.source[n.from:n.to]
		at := 0
		for range jsonBriefRunes + 1 {
			if at >= len(text)-1 {
				return text
			}
			_, w := utf8.DecodeRuneInString(text[at:])
			at += w
		}
		return text[:at] + `…"`
	}
	return d.source[n.from:n.to]
}
