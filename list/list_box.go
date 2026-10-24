// Package list implements HeroUI v3's ListBox: a virtualized collection of
// selectable rows grouped in labelled sections.
//
// The caller owns the items, the selection and the scroll state. Rows are
// built lazily by MyGo's List, so long collections stay cheap; keyboard
// navigation, typeahead and selection run on the native list state rather
// than on a second widget layer.
package list

import (
	"time"
	"unicode"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SelectionMode is HeroUI's list-box selection behavior.
type SelectionMode uint8

const (
	// SelectionNone renders actions only: rows never hold a choice.
	SelectionNone SelectionMode = iota
	// SelectionSingle keeps at most one chosen key.
	SelectionSingle
	// SelectionMultiple keeps any number of chosen keys.
	SelectionMultiple
)

// ListBoxVariant is HeroUI's list-box color treatment.
type ListBoxVariant uint8

const (
	ListBoxDefault ListBoxVariant = iota
	ListBoxDanger
)

// ListBoxItemOptions are the per-row presentation overrides. The zero value
// follows the list: the check indicator shows with any selection mode and the
// row uses the list's variant.
type ListBoxItemOptions struct {
	// Indicator overrides whether the selection check is reserved for the row.
	Indicator *bool
	// Variant overrides the list's color for this row.
	Variant ListBoxVariant
	// Height overrides HeroUI's 36 DIP minimum row height.
	Height float32
}

// ListBoxItem is one row of a ListBox. Key is the stable identity used by the
// selection; it falls back to Label, and must be unique in the list.
//
// Start and End build leading and trailing content. They run for every built
// row and must not perform I/O.
type ListBoxItem struct {
	Key         string
	Label       string
	Description string
	Disabled    bool
	Start       func()
	End         func()
	Options     ListBoxItemOptions
}

func (it ListBoxItem) key() string {
	if it.Key != "" {
		return it.Key
	}
	return it.Label
}

// ListBoxSectionOptions configures a section header.
type ListBoxSectionOptions struct {
	// HideLabel omits the header while keeping the section boundary.
	HideLabel bool
	// End builds trailing header content, such as a count.
	End func()
}

// ListBoxSection is a labelled group of rows. Items keep their order.
type ListBoxSection struct {
	Key     string
	Label   string
	Items   []ListBoxItem
	Options ListBoxSectionOptions
}

// ListBoxOptions configures a ListBox. State is required and caller-owned.
//
// Items render before Sections. Height fixes the scroll viewport; without it
// the list sizes to its rows, and MaxHeight caps it.
type ListBoxOptions struct {
	// Label names the list for assistive technology.
	Label string
	// Items are ungrouped rows shown before Sections.
	Items []ListBoxItem
	// Sections are labelled groups shown in order.
	Sections []ListBoxSection
	// Selection selects HeroUI's selection mode. The zero value allows no
	// choice, only OnAction.
	Selection SelectionMode
	// Variant colors the list's rows.
	Variant ListBoxVariant
	// Disabled dims the list and stops every activation.
	Disabled bool
	// State is required and owned by the caller.
	State *ListBoxState
	// Height fixes the visible rows' viewport in DIPs.
	Height float32
	// MaxHeight caps a content-sized viewport in DIPs.
	MaxHeight float32
	// Padding is the content inset; zero uses HeroUI's 4 DIP base, scaled
	// with the host's text scale.
	Padding float32
	// Empty is shown when the list has no rows.
	Empty string
	// NoFocusOnClick leaves the keyboard focus where it is, as a text
	// input's popup needs.
	NoFocusOnClick bool
	// OnAction reports a row activated by a click, Enter or Space.
	OnAction func(key string)
	// OnSelectionChange reports the chosen keys after a user change.
	OnSelectionChange func(keys []string)
}

// ListBoxState is the caller-owned selection, highlight and scroll state of one
// ListBox. Keep one per list and do not copy it.
type ListBoxState struct {
	// List keeps the virtual list's place, row heights and focus.
	List ui.ListState
	// SelectedKey is the chosen key in SelectionSingle mode; empty for none.
	SelectedKey string
	// Selected holds the chosen keys in SelectionMultiple mode.
	Selected ui.Selection[string]
	// Highlight is the flattened row under the keyboard or pointer.
	Highlight int

	rows       []listBoxRow
	pending    string
	hasPending bool
}

// listBoxRow is a flattened row, kept between frames for keyboard handling.
type listBoxRow struct {
	key      string
	label    string
	header   bool
	disabled bool
	section  int // -1 for an ungrouped item
	item     int // -1 for a header
}

// Focus gives the list the keyboard focus and reports whether it was built
// in the current pass. Call it after building the list.
func (s *ListBoxState) Focus(c *ui.Context) bool { return s.List.Focus(c) }

// Focused reports whether the list holds the keyboard focus.
func (s *ListBoxState) Focused(c *ui.Context) bool { return s.List.Focused(c) }

// FocusWithin reports whether the list or a row inside it holds the focus.
func (s *ListBoxState) FocusWithin(c *ui.Context) bool { return s.List.FocusWithin(c) }

// IsSelected reports whether key is chosen in either mode.
func (s *ListBoxState) IsSelected(key string) bool {
	return s.SelectedKey == key || s.Selected.Has(key)
}

// SetSelected chooses or unchoses key for mode.
func (s *ListBoxState) SetSelected(mode SelectionMode, key string, on bool) {
	if key == "" {
		return
	}
	switch mode {
	case SelectionSingle:
		if on {
			s.SelectedKey = key
		} else if s.SelectedKey == key {
			s.SelectedKey = ""
		}
	case SelectionMultiple:
		if on {
			s.Selected.Add(key)
		} else {
			s.Selected.Remove(key)
		}
	}
}

// Toggle flips the choice of key for mode. Single mode always ends with key
// chosen; choosing the chosen row again keeps it, as HeroUI's select does.
func (s *ListBoxState) Toggle(mode SelectionMode, key string) {
	switch mode {
	case SelectionSingle:
		s.SelectedKey = key
	case SelectionMultiple:
		if s.Selected.Has(key) {
			s.Selected.Remove(key)
		} else {
			s.Selected.Add(key)
		}
	}
}

// ClearSelection chooses nothing in either mode.
func (s *ListBoxState) ClearSelection() {
	s.SelectedKey = ""
	s.Selected.Clear()
}

// SelectionCount reports how many keys mode holds.
func (s *ListBoxState) SelectionCount(mode SelectionMode) int {
	switch mode {
	case SelectionSingle:
		if s.SelectedKey != "" {
			return 1
		}
	case SelectionMultiple:
		return s.Selected.Len()
	}
	return 0
}

// SelectedKeys returns the chosen keys for mode. Multiple keys come in no
// order; components displaying values order them by their items.
func (s *ListBoxState) SelectedKeys(mode SelectionMode) []string {
	switch mode {
	case SelectionSingle:
		if s.SelectedKey != "" {
			return []string{s.SelectedKey}
		}
	case SelectionMultiple:
		if n := s.Selected.Len(); n > 0 {
			keys := make([]string, 0, n)
			for k := range s.Selected.All() {
				keys = append(keys, k)
			}
			return keys
		}
	}
	return nil
}

// SetHighlight moves the keyboard highlight to key, scrolling it into view.
// A key not built yet is applied on the next build.
func (s *ListBoxState) SetHighlight(key string) {
	for i := range s.rows {
		if s.rows[i].key == key && !s.rows[i].header && !s.rows[i].disabled {
			s.Highlight = i
			s.List.ScrollIntoView(i)
			return
		}
	}
	s.pending, s.hasPending = key, true
}

// HighlightKey returns the highlighted item's key, or "" when the highlight is
// on no item.
func (s *ListBoxState) HighlightKey() string {
	if s.Highlight < 0 || s.Highlight >= len(s.rows) {
		return ""
	}
	r := s.rows[s.Highlight]
	if r.header || r.disabled {
		return ""
	}
	return r.key
}

// MoveHighlight moves the keyboard highlight by delta rows, skipping headers
// and disabled rows, and returns the new key.
func (s *ListBoxState) MoveHighlight(delta int) string {
	if len(s.rows) == 0 {
		return ""
	}
	step := 1
	if delta < 0 {
		step, delta = -1, -delta
	}
	i := s.Highlight
	if i < 0 || i >= len(s.rows) {
		i = 0
		if step < 0 {
			i = len(s.rows) - 1
		}
		if !s.selectable(i) {
			if j, ok := s.step(i, step); ok {
				i = j
			}
		}
		delta--
	}
	for ; delta > 0; delta-- {
		j, ok := s.step(i, step)
		if !ok {
			break
		}
		i = j
	}
	s.Highlight = i
	s.List.ScrollIntoView(i)
	return s.HighlightKey()
}

// HighlightEdge moves the highlight to the first or the last selectable row.
func (s *ListBoxState) HighlightEdge(last bool) string {
	n := len(s.rows)
	if n == 0 {
		return ""
	}
	step := 1
	if last {
		step = -1
	}
	start := 0
	if last {
		start = n - 1
	}
	for i := start; i >= 0 && i < n; i += step {
		if s.selectable(i) {
			s.Highlight = i
			s.List.ScrollIntoView(i)
			return s.rows[i].key
		}
	}
	return ""
}

func (s *ListBoxState) selectable(i int) bool {
	return i >= 0 && i < len(s.rows) && !s.rows[i].header && !s.rows[i].disabled
}

func (s *ListBoxState) step(from, by int) (int, bool) {
	for i := from + by; i >= 0 && i < len(s.rows); i += by {
		if s.selectable(i) {
			return i, true
		}
	}
	return from, false
}

// buildRows flattens the caller's items into the scratch rows slice.
func (s *ListBoxState) buildRows(o ListBoxOptions) {
	s.rows = s.rows[:0]
	for i := range o.Items {
		it := &o.Items[i]
		s.rows = append(s.rows, listBoxRow{
			key: it.key(), label: it.Label, disabled: it.Disabled, section: -1, item: i,
		})
	}
	for si := range o.Sections {
		sec := &o.Sections[si]
		if sec.Label != "" && !sec.Options.HideLabel {
			key := sec.Key
			if key == "" {
				key = "section:" + sec.Label
			}
			s.rows = append(s.rows, listBoxRow{
				key: key, label: sec.Label, header: true, section: si, item: -1,
			})
		}
		for ii := range sec.Items {
			it := &sec.Items[ii]
			s.rows = append(s.rows, listBoxRow{
				key: it.key(), label: it.Label, disabled: it.Disabled, section: si, item: ii,
			})
		}
	}
	s.resolveHighlight()
}

func (s *ListBoxState) resolveHighlight() {
	if s.hasPending {
		for i := range s.rows {
			if s.rows[i].key == s.pending && s.selectable(i) {
				s.Highlight = i
				s.List.ScrollIntoView(i)
				break
			}
		}
		s.pending, s.hasPending = "", false
	}
	if len(s.rows) == 0 {
		s.Highlight = 0
		return
	}
	if s.Highlight < 0 {
		s.Highlight = 0
	}
	if s.Highlight >= len(s.rows) {
		s.Highlight = len(s.rows) - 1
	}
	if s.selectable(s.Highlight) {
		return
	}
	if i, ok := s.step(s.Highlight, 1); ok {
		s.Highlight = i
		return
	}
	if i, ok := s.step(s.Highlight, -1); ok {
		s.Highlight = i
	}
}

func (s *ListBoxState) keyAt(i int) any {
	if i < 0 || i >= len(s.rows) {
		return i
	}
	return s.rows[i].key
}

func (s *ListBoxState) labelAt(i int) string {
	if i < 0 || i >= len(s.rows) {
		return ""
	}
	return s.rows[i].label
}

// rowVariant resolves a row's variant from the item and the list.
func rowVariant(list ListBoxVariant, it *ListBoxItem) ListBoxVariant {
	if it.Options.Variant != ListBoxDefault {
		return it.Options.Variant
	}
	return list
}

// indicator reports whether the check indicator is reserved for a row.
func indicator(mode SelectionMode, it *ListBoxItem) bool {
	if it.Options.Indicator != nil {
		return *it.Options.Indicator
	}
	return mode != SelectionNone
}

var typeKeys = [...]ui.Key{
	ui.KeyA, ui.KeyB, ui.KeyC, ui.KeyD, ui.KeyE, ui.KeyF, ui.KeyG, ui.KeyH,
	ui.KeyI, ui.KeyJ, ui.KeyK, ui.KeyL, ui.KeyM, ui.KeyN, ui.KeyO, ui.KeyP,
	ui.KeyQ, ui.KeyR, ui.KeyS, ui.KeyT, ui.KeyU, ui.KeyV, ui.KeyW, ui.KeyX,
	ui.KeyY, ui.KeyZ,
}

var typeRunes = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}

// firstFold returns the first letter of s, folded to lower case, or 0.
func firstFold(s string) rune {
	for _, r := range s {
		return unicode.ToLower(r)
	}
	return 0
}

// typeahead moves the highlight to the next row whose label starts with r.
func (s *ListBoxState) typeahead(r rune) string {
	n := len(s.rows)
	if n == 0 {
		return ""
	}
	from := s.Highlight
	if s.selectable(from) {
		from++
	}
	for d := range n {
		i := (from + d) % n
		if !s.selectable(i) {
			continue
		}
		if firstFold(s.rows[i].label) != r {
			continue
		}
		s.Highlight = i
		s.List.ScrollIntoView(i)
		return s.rows[i].key
	}
	return ""
}

// ListBox renders a HeroUI list box. Rows are virtualized through the native
// List, and hover, highlight and check paints read live pointer and focus
// state.
func ListBox(c *ui.Context, o ListBoxOptions) ui.Element {
	s := o.State
	if s == nil {
		panic("list: ListBox requires a caller-owned State")
	}
	s.buildRows(o)
	if s.List.Key == nil {
		s.List.Key = s.keyAt
		s.List.Label = s.labelAt
		// Headers scroll with their section, as HeroUI's do. The native
		// ListState.Header would pin them to the viewport top, where they
		// would overlap the rows scrolling under them.
		s.List.Header = nil
	}
	n := len(s.rows)
	padding := o.Padding
	if padding == 0 {
		padding = theme.Scale(c, theme.Space4)
	}
	t := theme.Colors(c)
	mode := o.Selection
	st := rowStyle{
		palette: t,
		radius:  theme.RadiusScale().XXL,
		padX:    theme.Scale(c, theme.Space8),
		padY:    theme.Scale(c, theme.Space6),
		gap:     theme.Scale(c, theme.Space12),
		// HeroUI reserves pe-7 for the absolutely placed indicator and puts
		// the 16 DIP mark 8 DIP from the row's end.
		endPad:   theme.Scale(c, theme.Space(7)),
		mark:     theme.Scale(c, theme.Space16),
		check:    theme.Scale(c, theme.Space(2.5)),
		blockGap: theme.Scale(c, theme.Space4),
	}
	// The list's focus is the same for every row built in this pass; resolve
	// it once, when the first row is built and the list is bound.
	var focused, focusKnown bool
	list := ui.List(c, &s.List, n, func(i int) {
		if !focusKnown {
			focused = c.Resolve(s.List.Handle).FocusVisible()
			focusKnown = true
		}
		buildRow(c, s, o, i, st, focused)
	}).Padding(padding).Gap(0).MinWidth(0).FillWidth().
		Focusable().FocusRing(false).Label(o.Label)
	if o.Height > 0 {
		list.Height(o.Height)
	}
	if o.MaxHeight > 0 {
		list.MaxHeight(o.MaxHeight)
	}
	if o.Disabled {
		list = style.Disabled(list, true)
	}
	if n == 0 && o.Empty != "" {
		list.Children(func() {
			style.Text(ui.Text(c, o.Empty), c, theme.SM()).
				TextColor(t.Muted).Padding(theme.Scale(c, theme.Space8))
		})
	}
	if !o.Disabled {
		bindKeys(c, s, list, o, mode)
	}
	return list
}

// rowStyle is the measured form of a list's options, shared by the rows built
// in one pass. Spacing and icons follow HeroUI's rem metrics, so they scale
// with the host's text scale.
type rowStyle struct {
	palette             theme.Palette
	radius              float32
	padX, padY, gap     float32
	endPad, mark, check float32
	// blockGap separates a block: an ungrouped row, or a whole section. A
	// section's own rows sit flush, as HeroUI's gap-0 section does.
	blockGap float32
}

func buildRow(c *ui.Context, s *ListBoxState, o ListBoxOptions, i int, st rowStyle, focused bool) {
	r := s.rows[i]
	gap := st.blockGapAfter(s, i)
	if r.header {
		sec := &o.Sections[r.section]
		header := ui.Row(c).FillWidth().Shrink(0).
			Gap(theme.Scale(c, theme.Space8)).AlignItems(ui.Center).
			Margin(0, 0, gap, 0).
			Padding(theme.Scale(c, theme.Space(1.5)), theme.Scale(c, theme.Space8),
				theme.Scale(c, theme.Space4), theme.Scale(c, theme.Space8))
		header.Children(func() {
			style.Text(ui.Text(c, sec.Label), c, theme.XS()).
				FontWeight(theme.WeightMedium).TextColor(st.palette.Muted).
				Role(ui.RoleHeading).Level(3).Grow(1)
			if sec.Options.End != nil {
				sec.Options.End()
			}
		})
		return
	}
	var it *ListBoxItem
	if r.section < 0 {
		it = &o.Items[r.item]
	} else {
		it = &o.Sections[r.section].Items[r.item]
	}
	variant := rowVariant(o.Variant, it)
	fg := st.palette.Foreground
	if variant == ListBoxDanger {
		fg = st.palette.Danger
	}
	height := it.Options.Height
	if height == 0 {
		height = 36
	}
	showIndicator := indicator(o.Selection, it)
	selected := s.IsSelected(r.key)
	row := ui.Row(c).FillWidth().Shrink(0).Gap(st.gap).AlignItems(ui.Center).
		Margin(0, 0, gap, 0).Padding(st.padY, st.padX).
		MinHeight(theme.Scale(c, height)).Radius(st.radius).Clip().
		Label(it.Label).TextColor(fg)
	if showIndicator {
		row.PaddingEnd(st.endPad)
	}
	highlighted := s.Highlight == i
	// Selected sets the toolkit's accent face; HeroUI's selected row changes
	// only its indicator, so the fill is cleared while the chosen state stays
	// for assistive technology.
	row.Selected(selected).TextColor(fg).Background(ui.Transparent)
	// Snapshot the row's hover in the build: a retained copy of a leaving
	// row paints in later passes, where row has expired.
	hovered := row.Hovered()
	row.Draw(func(p *ui.Painter, box ui.Rect) {
		if hovered {
			p.Fill(box, st.palette.Default, st.radius)
		}
		if focused && highlighted {
			p.FocusRing(box, [4]float32{st.radius, st.radius, st.radius, st.radius})
		}
	})
	enabled := !o.Disabled && !it.Disabled
	if enabled {
		row.Cursor(ui.CursorPointer)
	}
	row.Children(func() {
		if it.Start != nil {
			it.Start()
		}
		column := ui.Column(c).Grow(1).MinWidth(0).Gap(0).Justify(ui.Center)
		column.Children(func() {
			if it.Description == "" {
				style.Text(ui.Text(c, it.Label), c, theme.SM()).
					TextColor(fg).SingleLine()
			} else {
				style.Text(ui.Text(c, it.Label), c, theme.SM()).
					TextColor(fg).SingleLine()
				style.Text(ui.Text(c, it.Description), c, theme.XS()).
					TextColor(st.palette.Muted).SingleLine()
			}
		})
		if it.End != nil {
			it.End()
		}
		if showIndicator {
			indicatorMark(c, st, fg, selected)
		}
	})
	if enabled {
		row.OnClick(func() {
			if !o.NoFocusOnClick {
				s.List.Focus(c)
			}
			activate(c, s, o, r.key, o.Selection)
		})
	}
	style.Disabled(row, !enabled)
}

// blockGapAfter returns the space below row i: none inside a section, and the
// list's block gap after an ungrouped row or a whole section. The last row
// leaves the end to the list's padding.
func (st rowStyle) blockGapAfter(s *ListBoxState, i int) float32 {
	if i+1 >= len(s.rows) {
		return 0
	}
	r, next := s.rows[i], s.rows[i+1]
	if r.section >= 0 && next.section == r.section && !next.header {
		return 0
	}
	return st.blockGap
}

// indicatorMark is HeroUI's absolutely placed check slot: the row reserves
// its 28 DIP end padding, so the mark never moves the label. The mark scales
// and fades in as the selection changes, as HeroUI's check draws in.
func indicatorMark(c *ui.Context, st rowStyle, fg ui.Color, selected bool) {
	mark := ui.Box(c).Absolute().End(theme.Scale(c, theme.Space8)).Top(0).Bottom(0).
		Width(st.mark).Height(st.mark).
		Margin(ui.Auto, 0, ui.Auto, 0).Center().
		Role(ui.RoleNone)
	mark.Children(func() {
		if theme.ReduceMotion(c) {
			if selected {
				ui.Icon(c, checkIcon).Size(st.check, st.check).TextColor(fg)
			}
			return
		}
		check := ui.Icon(c, checkIcon).Size(st.check, st.check).TextColor(fg)
		pose := ui.Pose{Scale: 0.6, Fade: 1}
		if selected {
			pose = ui.Pose{}
		}
		check.AnimatePose(pose, ui.Animation{Duration: 200 * time.Millisecond, Ease: ui.EaseOut})
	})
}

// notify reports the keys chosen for mode.
func notify(o ListBoxOptions, s *ListBoxState) {
	if o.OnSelectionChange != nil {
		o.OnSelectionChange(s.SelectedKeys(o.Selection))
	}
}

// activate chooses a row for the list's mode and reports the action.
func activate(c *ui.Context, s *ListBoxState, o ListBoxOptions, key string, mode SelectionMode) {
	switch mode {
	case SelectionSingle:
		changed := s.SelectedKey != key
		s.SelectedKey = key
		if changed {
			notify(o, s)
		}
	case SelectionMultiple:
		s.Toggle(mode, key)
		notify(o, s)
	}
	if o.OnAction != nil {
		o.OnAction(key)
	}
}

// bindKeys registers the list's keyboard behavior on its focusable element.
func bindKeys(c *ui.Context, s *ListBoxState, list ui.Element, o ListBoxOptions, mode SelectionMode) {
	list.OnShortcut(0, ui.KeyDown, func() { s.MoveHighlight(1) })
	list.OnShortcut(0, ui.KeyUp, func() { s.MoveHighlight(-1) })
	list.OnShortcut(0, ui.KeyHome, func() { s.HighlightEdge(false) })
	list.OnShortcut(0, ui.KeyEnd, func() { s.HighlightEdge(true) })
	list.OnShortcut(0, ui.KeySpace, func() {
		if mode == SelectionNone {
			if key := s.HighlightKey(); key != "" && o.OnAction != nil {
				o.OnAction(key)
			}
			return
		}
		if key := s.HighlightKey(); key != "" {
			activate(c, s, o, key, mode)
		}
	})
	list.OnShortcut(0, ui.KeyEnter, func() {
		key := s.HighlightKey()
		if key == "" {
			return
		}
		if mode != SelectionNone {
			switch mode {
			case SelectionSingle:
				if s.SelectedKey != key {
					s.SelectedKey = key
					notify(o, s)
				}
			case SelectionMultiple:
				s.Toggle(mode, key)
				notify(o, s)
			}
		}
		if o.OnAction != nil {
			o.OnAction(key)
		}
	})
	for i := range typeKeys {
		r := typeRunes[i]
		list.OnShortcut(0, typeKeys[i], func() { s.typeahead(r) })
	}
}

// checkIcon is HeroUI's checkmark, shared and immutable.
var checkIcon = lucide.Check()
