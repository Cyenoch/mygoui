// Package list implements HeroUI v3's ListBox: a virtualized collection of
// selectable rows grouped in labelled sections.
//
// The caller owns the items, the selection and the scroll state. Rows are
// built lazily by MyGo's List, so long collections stay cheap; keyboard
// navigation, typeahead and selection run on the native list state rather
// than on a second widget layer.
package list

import (
	"strings"

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
	// Padding is the content inset; zero uses HeroUI's 4 DIP base.
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
		label := s.rows[i].label
		if label == "" {
			continue
		}
		first := []rune(strings.ToLower(label))[0]
		if first == r {
			s.Highlight = i
			s.List.ScrollIntoView(i)
			return s.rows[i].key
		}
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
		s.List.Header = func(i int) bool {
			return i >= 0 && i < len(s.rows) && s.rows[i].header
		}
	}
	n := len(s.rows)
	padding := o.Padding
	if padding == 0 {
		padding = theme.Space4
	}
	t := theme.Colors(c)
	radius := theme.RadiusScale().XXL
	mode := o.Selection

	list := ui.List(c, &s.List, n, func(i int) {
		buildRow(c, s, o, i, radius, t)
	}).Padding(padding).Gap(theme.Space4).MinWidth(0).FillWidth().
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
				TextColor(t.Muted).Padding(theme.Space8)
		})
	}
	if !o.Disabled {
		bindKeys(c, s, list, o, mode)
	}
	return list
}

func buildRow(c *ui.Context, s *ListBoxState, o ListBoxOptions, i int, radius float32, t theme.Palette) {
	r := s.rows[i]
	if r.header {
		sec := &o.Sections[r.section]
		header := ui.Row(c).FillWidth().Shrink(0).Gap(theme.Space8).AlignItems(ui.Center).
			Padding(theme.Space6, theme.Space8, theme.Space4, theme.Space8)
		header.Children(func() {
			style.Text(ui.Text(c, sec.Label), c, theme.XS()).
				FontWeight(theme.WeightMedium).TextColor(t.Muted).
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
	fg := t.Foreground
	if variant == ListBoxDanger {
		fg = t.Danger
	}
	height := it.Options.Height
	if height == 0 {
		height = 36
	}
	row := ui.Row(c).FillWidth().Shrink(0).Gap(theme.Space12).AlignItems(ui.Center).
		Padding(theme.Space6, theme.Space8, theme.Space6, theme.Space8).
		MinHeight(theme.Scale(c, height)).Radius(radius).
		Label(it.Label).TextColor(fg)
	selected := s.IsSelected(r.key)
	highlighted := s.Highlight == i
	focused := false
	if el := c.Resolve(s.List.Handle); el.Valid() {
		focused = el.FocusVisible()
	}
	row.Selected(selected).TextColor(fg)
	row.Background(ui.Transparent).Radius(radius).Clip()
	// Snapshot the row's hover in the build: a retained copy of a leaving
	// row paints in later passes, where row has expired.
	hovered := row.Hovered()
	row.Draw(func(p *ui.Painter, box ui.Rect) {
		if hovered {
			p.Fill(box, t.Default, radius)
		}
		if focused && highlighted {
			p.FocusRing(ui.Rect{X: box.X + 3, Y: box.Y + 3, W: box.W - 6, H: box.H - 6},
				[4]float32{radius, radius, radius, radius})
		}
	})
	showIndicator := indicator(o.Selection, it)
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
					TextColor(t.Muted).SingleLine()
			}
		})
		if it.End != nil {
			it.End()
		}
		if showIndicator {
			mark := ui.Box(c).Width(theme.Scale(c, 16)).Height(theme.Scale(c, 16)).
				Shrink(0).Center().Role(ui.RoleNone)
			if selected {
				mark.Children(func() {
					ui.Icon(c, checkIcon).Size(theme.Scale(c, 10), theme.Scale(c, 10)).
						TextColor(fg)
				})
			}
		}
	})
	if !o.Disabled && !it.Disabled {
		row.OnClick(func() {
			if !o.NoFocusOnClick {
				s.List.Focus(c)
			}
			activate(c, s, o, r.key, o.Selection)
		})
	}
	style.Disabled(row, o.Disabled || it.Disabled)
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
