package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// JSONViewerOptions configures a JSON tree. Source is the document; pass
// string(data) for bytes, converted where the data changes rather than in
// every build. Title names the viewer in its header, "JSON" when empty.
// Height is the tree's height in DIPs; zero makes it as high as its rows,
// up to 320 DIPs. ExpandDepth is how many levels are open when the viewer
// first shows a document: two when zero, only the root when negative.
// Toolbar adds a search field above the tree, with the number of matches
// and a way to step through them. OnError receives a clipboard failure of a
// copy.
type JSONViewerOptions struct {
	Source      string
	Title       string
	Height      float32
	ExpandDepth int
	Toolbar     bool
	OnError     func(error)
	// State is caller-owned expansion and scroll state; the viewer keeps
	// its own per window when nil.
	State *JSONViewerState
}

// JSONViewerState holds which nodes of a document are open, the row chosen
// and the scroll position. Nodes stay open by their path when Source
// changes.
//
// Query marks the keys and values holding it, whatever their case, and
// chooses the first of them, opening the objects and arrays above it; the
// toolbar's field writes it.
type JSONViewerState struct {
	Query string

	outline  ui.OutlineState[int]
	selected int
	roots    []int
	doc      *jsonDoc
	source   string
	err      error
	loaded   bool
	pending  jsonExpansion
	found    jsonMatches
	search   ui.Handle
}

type jsonExpansion uint8

const (
	jsonExpansionKept jsonExpansion = iota
	jsonExpandAll
	jsonCollapseAll
)

// ExpandAll opens every object and array in the next build.
func (s *JSONViewerState) ExpandAll() { s.pending = jsonExpandAll }

// CollapseAll closes every object and array below the root in the next
// build.
func (s *JSONViewerState) CollapseAll() { s.pending = jsonCollapseAll }

// FocusSearch focuses the toolbar's search field, as Cmd+F (Ctrl+F) does
// inside the viewer, for a Find command of the window.
func (s *JSONViewerState) FocusSearch() { s.search.Focus() }

type jsonKind uint8

const (
	jsonObject jsonKind = iota
	jsonArray
	jsonString
	jsonNumber
	jsonKeyword
)

// jsonNode is one value of a document. A scalar's text, and a container's,
// is source[from:to].
type jsonNode struct {
	key      string // member name in an object
	kind     jsonKind
	parent   int // -1 for the root
	index    int // position among its siblings
	from, to int
	// kids are the node's children in document order; nil for a scalar and
	// for an empty object or array, which shows no arrow.
	kids []int
}

type jsonDoc struct {
	source string
	nodes  []jsonNode
}

// parseJSON reads one JSON value into nodes in document order, keeping the
// order of object members.
func parseJSON(source string) (*jsonDoc, error) {
	dec := json.NewDecoder(strings.NewReader(source))
	dec.UseNumber()
	doc := &jsonDoc{source: source}
	var open []int // the containers being read, innermost last
	key, keyed := "", false
	for {
		from := jsonValueStart(source, int(dec.InputOffset()))
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) && len(doc.nodes) > 0 && len(open) == 0 {
			return doc, nil
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = io.ErrUnexpectedEOF
			}
			return nil, fmt.Errorf("content: invalid JSON at byte %d: %w", dec.InputOffset(), err)
		}
		if len(doc.nodes) > 0 && len(open) == 0 {
			return nil, fmt.Errorf("content: invalid JSON at byte %d: data after the document", from)
		}
		to := int(dec.InputOffset())
		parent := -1
		if len(open) > 0 {
			parent = open[len(open)-1]
		}
		if delim, ok := tok.(json.Delim); ok && (delim == '}' || delim == ']') {
			doc.nodes[parent].to = to
			open = open[:len(open)-1]
			continue
		}
		if name, ok := tok.(string); ok && parent >= 0 && doc.nodes[parent].kind == jsonObject && !keyed {
			key, keyed = name, true
			continue
		}
		node := jsonNode{key: key, parent: parent, from: from, to: to}
		key, keyed = "", false
		at := len(doc.nodes)
		if parent >= 0 {
			node.index = len(doc.nodes[parent].kids)
			doc.nodes[parent].kids = append(doc.nodes[parent].kids, at)
		}
		switch v := tok.(type) {
		case json.Delim:
			node.kind = jsonArray
			if v == '{' {
				node.kind = jsonObject
			}
			open = append(open, at)
		case string:
			node.kind = jsonString
		case json.Number:
			node.kind = jsonNumber
		default: // true, false and null
			node.kind = jsonKeyword
		}
		doc.nodes = append(doc.nodes, node)
	}
}

// jsonValueStart skips the white space and separators before a token.
func jsonValueStart(source string, at int) int {
	for at < len(source) && strings.IndexByte(" \t\r\n,:", source[at]) >= 0 {
		at++
	}
	return at
}

func (d *jsonDoc) children(i int) []int { return d.nodes[i].kids }

// path returns the JSONPath of node i, as $.items[0].name.
func (d *jsonDoc) path(i int) string {
	if i == 0 {
		return "$"
	}
	n := &d.nodes[i]
	parent := d.path(n.parent)
	switch {
	case d.nodes[n.parent].kind == jsonArray:
		return parent + "[" + strconv.Itoa(n.index) + "]"
	case jsonIdentifier(n.key):
		return parent + "." + n.key
	}
	return parent + "[" + strconv.Quote(n.key) + "]"
}

func jsonIdentifier(key string) bool {
	if key == "" || !identStart(key[0]) {
		return false
	}
	for i := 1; i < len(key); i++ {
		if !identStart(key[i]) && !digit(key[i]) {
			return false
		}
	}
	return true
}

// value returns what copying node i writes: a string's characters, and
// the JSON text of any other value.
func (d *jsonDoc) value(i int) string {
	n := &d.nodes[i]
	text := d.source[n.from:n.to]
	if n.kind != jsonString {
		return text
	}
	var s string
	// The decoder accepted this token, so it unquotes.
	_ = json.Unmarshal([]byte(text), &s)
	return s
}

// load parses source when it changed. A document that does not parse keeps
// the last good one's open nodes for the next that does.
func (s *JSONViewerState) load(source string, depth int) {
	if s.loaded && s.source == source {
		return
	}
	s.source, s.loaded = source, true
	doc, err := parseJSON(source)
	s.err = err
	if err != nil {
		return
	}
	previous := s.doc
	s.doc, s.selected = doc, -1
	if previous == nil {
		s.roots = []int{0}
		s.openLevels(depth)
		return
	}
	open := make(map[string]struct{}, s.outline.Open.Len())
	for i := range s.outline.Open.All() {
		open[previous.path(i)] = struct{}{}
	}
	s.outline.Open.Clear()
	for i := range doc.nodes {
		// The root stays in the set while it has no children, so that it
		// is open again once a later document gives it some.
		if i > 0 && doc.nodes[i].kids == nil {
			continue
		}
		if _, ok := open[doc.path(i)]; ok {
			s.outline.Open.Add(i)
		}
	}
}

// openLevels opens the root and the objects and arrays above depth levels.
func (s *JSONViewerState) openLevels(depth int) {
	s.outline.Open.Add(0)
	levels := make([]int, len(s.doc.nodes))
	for i := 1; i < len(levels); i++ {
		levels[i] = levels[s.doc.nodes[i].parent] + 1
		if levels[i] < depth && s.doc.nodes[i].kids != nil {
			s.outline.Open.Add(i)
		}
	}
}

// applyExpansion carries out an ExpandAll or CollapseAll asked since the
// last build.
func (s *JSONViewerState) applyExpansion() {
	pending := s.pending
	s.pending = jsonExpansionKept
	if pending == jsonExpansionKept || s.doc == nil {
		return
	}
	s.outline.Open.Clear()
	s.outline.Open.Add(0)
	if pending == jsonCollapseAll {
		return
	}
	for i := range s.doc.nodes {
		if s.doc.nodes[i].kids != nil {
			s.outline.Open.Add(i)
		}
	}
}

// jsonDefaultMaxHeight bounds a tree without a Height, so that it stays
// virtualized however many nodes are open.
const jsonDefaultMaxHeight = 320

// jsonSelectedMaxLines bounds how many lines the chosen string wraps to;
// copying still takes all of it.
const jsonSelectedMaxLines = 30

type jsonViewerStateKey struct{}

// JSONViewer shows a JSON document as a tree whose objects and arrays open
// and close, building only the rows in view. Members keep the document's
// order. A closed object or array previews what it holds; the row chosen
// wraps its string in full, and Copy (Cmd or Ctrl and C) writes its value.
// A row's context menu copies its value or its path; a double click opens
// or closes it; the header expands or collapses every node and copies the
// document. The tree is MyGo's native outline: the arrows, Left and Right
// open and close nodes.
//
// Source that is not one valid JSON value returns the error, which the
// viewer shows in place of the tree.
func JSONViewer(c *ui.Context, o JSONViewerOptions) (ui.Element, error) {
	t := theme.Colors(c)
	title := o.Title
	if title == "" {
		title = "JSON"
	}
	root := ui.Column(c).AlignSelf(ui.Stretch).MinWidth(0).Clip().Radius(theme.RadiusScale().XL).
		Background(t.Code).TextColor(t.CodeForeground).
		BorderWidth(theme.Hairline).BorderColor(t.Border).Role(ui.RoleGroup).Label(title)
	st := o.State
	if st == nil {
		st = ui.Local(root, jsonViewerStateKey{}, func() JSONViewerState { return JSONViewerState{} })
	}
	depth := o.ExpandDepth
	switch {
	case depth == 0:
		depth = 2
	case depth < 0:
		depth = 1
	}
	st.load(o.Source, depth)
	st.applyExpansion()
	st.syncMatches()
	if st.err == nil {
		// The search field keeps its own copy, so only the tree's copy is
		// the chosen row's.
		searching := st.search.Focused(c)
		root.OnShortcut(ui.Cmd, ui.KeyC, func() {
			if i, ok := st.chosen(); ok && !searching {
				c.WriteClipboardThen(st.doc.value(i), func(err error) {
					if err != nil && o.OnError != nil {
						o.OnError(err)
					}
				})
			}
		})
		if o.Toolbar {
			root.OnShortcut(ui.Cmd, ui.KeyF, func() { st.search.Focus(c) })
		}
	}
	root.Children(func() {
		jsonViewerHeader(c, t, title, st, &o)
		if st.err != nil {
			jsonViewerError(c, t, st.err)
			return
		}
		if o.Toolbar {
			jsonToolbar(c, t, st)
		}
		jsonViewerTree(c, t, title, st, &o)
	})
	return root, st.err
}

func jsonViewerHeader(c *ui.Context, t theme.Palette, title string, st *JSONViewerState, o *JSONViewerOptions) {
	m := locale.M()
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	ui.Row(c).FillWidth().MinWidth(0).MinHeight(theme.Scale(c, 40)).
		AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(1))).
		Padding(0, theme.Scale(c, 10), 0, theme.Scale(c, theme.Space(4))).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).
		Children(func() {
			style.Text(ui.Text(c, title), c, ts).Grow(1).MinWidth(0).SingleLine().Unselectable()
			action := func(label string, icon *ui.SVG, do func()) {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: label, Icon: icon, IconOnly: true,
					Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
					Disabled: st.err != nil, OnClick: do,
				})
			}
			// A button keeps its own alignment, so the row around them
			// centers them in the header.
			ui.Row(c).Shrink(0).Gap(theme.Scale(c, theme.Space(1))).Children(func() {
				action(m.ExpandAll, lucide.ChevronsUpDown(), st.ExpandAll)
				action(m.CollapseAll, lucide.ChevronsDownUp(), st.CollapseAll)
				CopyButton(c, CopyButtonOptions{Text: o.Source, Label: copyLabel(title), Size: theme.ButtonSmall, OnError: o.OnError})
			})
		})
}

func jsonViewerError(c *ui.Context, t theme.Palette, err error) {
	ts := theme.Code()
	ts.Size, ts.LineHeight = codeFontSize, codeLineHeight
	style.Text(ui.Text(c, err.Error()), c, ts).TextColor(t.DestructiveForeground).
		Padding(theme.Scale(c, theme.Space(4))).Selectable().Role(ui.RoleAlert)
}

func jsonViewerTree(c *ui.Context, t theme.Palette, title string, st *JSONViewerState, o *JSONViewerOptions) {
	doc := st.doc
	ts := theme.Code()
	ts.Size, ts.LineHeight = codeFontSize, codeLineHeight
	rows := jsonRows{t: t, colors: syntaxColorsOf(t), ts: ts, st: st, onError: o.OnError}
	if st.found.query != "" {
		rows.queries = []string{st.found.query}
	}
	st.outline.List.Selected = &st.selected
	// The native outline fills its chosen row with the theme's accent, the
	// primary colour here: a quiet wash reads as the place of the keys
	// without burying the syntax colours.
	quiet := *c.Theme()
	quiet.Accent = t.CodeForeground.Alpha(0.08)
	quiet.InactiveAccent = quiet.Accent
	quiet.AccentText = t.CodeForeground
	ui.Themed(c, &quiet).Background(ui.Transparent).AlignSelf(ui.Stretch).MinWidth(0).Children(func() {
		// JSON reads left to right in a right-to-left interface too.
		tree := ui.Outline(c, &st.outline, st.roots, doc.children, func(i int) {
			rows.build(c, i)
		}).FillWidth().MinWidth(0).Padding(theme.Scale(c, theme.Space(2))).
			Direction(ui.LeftToRight).Label(title)
		if o.Height > 0 {
			tree.Height(o.Height)
		} else {
			tree.MaxHeight(jsonDefaultMaxHeight)
		}
	})
}

// chosen is the node of the row chosen, as the last build showed it.
func (s *JSONViewerState) chosen() (int, bool) {
	if s.doc == nil || s.selected < 0 || s.selected >= s.outline.Rows() {
		return 0, false
	}
	return s.outline.Item(s.selected), true
}

// jsonRows builds the rows of one build.
type jsonRows struct {
	t       theme.Palette
	colors  syntaxColors
	ts      theme.TextStyle
	st      *JSONViewerState
	queries []string
	onError func(error)
}

// text is s in the row's type, with the search's matches marked.
func (r *jsonRows) text(c *ui.Context, s string) ui.Element {
	if r.queries == nil {
		return style.Text(ui.Text(c, s), c, r.ts)
	}
	return style.Text(mygoui.Highlight(c, mygoui.HighlightOptions{
		Text: s, Queries: r.queries, IgnoreCase: true,
	}), c, r.ts).TextColor(r.t.CodeForeground)
}

func (r *jsonRows) build(c *ui.Context, i int) {
	st := r.st
	doc := st.doc
	n := &doc.nodes[i]
	t, ts := r.t, r.ts
	chosen, _ := st.chosen()
	wrapped := chosen == i && n.kind == jsonString
	row := ui.Row(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(2)))
	if wrapped {
		row.AlignItems(ui.Start)
	}
	row.ContextMenu(func(m *ui.Menu) {
		msg := locale.M()
		copied := func(err error) {
			if err != nil && r.onError != nil {
				r.onError(err)
			}
		}
		if m.Item(msg.CopyValue).Chosen() {
			c.WriteClipboardThen(doc.value(i), copied)
		}
		if m.Item(msg.CopyPath).Chosen() {
			c.WriteClipboardThen(doc.path(i), copied)
		}
	})
	// A double click opens or closes a container, as its arrow does.
	if len(n.kids) > 0 && row.RowDoubleClicked() {
		if st.outline.Open.Has(i) {
			st.outline.Open.Remove(i)
		} else {
			st.outline.Open.Add(i)
		}
		c.Services().Invalidate()
	}
	row.Children(func() {
		switch {
		case n.parent < 0:
		case doc.nodes[n.parent].kind == jsonArray:
			style.Text(ui.Text(c, strconv.Itoa(n.index)), c, ts).TextColor(t.Muted).Shrink(0)
		default:
			ui.Row(c).Shrink(0).MaxWidthPercent(50).Children(func() {
				r.text(c, n.key).MinWidth(0).SingleLine()
				style.Text(ui.Text(c, ":"), c, ts).TextColor(t.Muted).Shrink(0)
			})
		}
		switch n.kind {
		case jsonObject, jsonArray:
			r.container(c, i)
		default:
			ink := r.colors.str
			switch n.kind {
			case jsonNumber:
				ink = r.colors.number
			case jsonKeyword:
				ink = t.InfoForeground
			}
			value := r.text(c, doc.source[n.from:n.to]).TextColor(ink).Grow(1).Basis(0).MinWidth(0)
			if wrapped {
				value.MaxLines(jsonSelectedMaxLines)
			} else {
				value.SingleLine()
			}
		}
	})
}

// container shows an object or array: its braces when empty, its opening
// brace while open, and what it holds while closed, with the count last.
func (r *jsonRows) container(c *ui.Context, i int) {
	doc := r.st.doc
	n := &doc.nodes[i]
	muted := r.t.Muted
	if len(n.kids) == 0 {
		style.Text(ui.Text(c, doc.brief(i)), c, r.ts).TextColor(muted).Grow(1).Basis(0).MinWidth(0).SingleLine()
		return
	}
	text := doc.preview(i)
	if r.st.outline.Open.Has(i) {
		text = text[:1]
	}
	style.Text(ui.Text(c, text), c, r.ts).TextColor(muted).Grow(1).Basis(0).MinWidth(0).SingleLine()
	count := r.ts
	count.Size--
	style.Text(ui.Text(c, strconv.Itoa(len(n.kids))), c, count).TextColor(muted.Alpha(0.6)).Shrink(0)
}
