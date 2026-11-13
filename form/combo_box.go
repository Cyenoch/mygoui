package form

import (
	"strings"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/busy"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ComboBoxOptions configures a combo box: a text input filtering a
// virtualized list of options.
//
// SelectionMultiple shows chosen keys as removable chips before the input.
// AllowsCustomValue keeps text that matches no item as State.Custom. State is
// required and caller-owned.
type ComboBoxOptions struct {
	Label       string
	Description string
	Error       string
	Placeholder string
	Required    bool
	Disabled    bool
	Invalid     bool
	FullWidth   bool
	Clearable   bool
	// Loading crossfades the chevron to a spinner. The popup shows a Loading
	// row when it has no options and keeps its rows otherwise. The input
	// stays editable.
	Loading   bool
	Selection list.SelectionMode
	Variant   ControlVariant
	Items     []list.ListBoxItem
	Sections  []list.ListBoxSection
	State     *ComboBoxState
	MaxHeight float32
	Placement mygoui.Placement

	// AllowsCustomValue keeps text that matches no option.
	AllowsCustomValue bool
	// Filter overrides the default case-insensitive substring match on Label.
	// It must not perform side effects. It runs again when the query or the
	// options change; bump State.Revision when its results change for the
	// same query.
	Filter func(query string, item *list.ListBoxItem) bool
	// Empty is shown when nothing matches; zero uses the localized message.
	Empty string
	// OnChange reports the chosen keys in item order.
	OnChange func(keys []string)
	// OnQueryChange reports each text change for caller-owned async work.
	OnQueryChange func(query string)
}

// ComboBoxState is the caller-owned open flag, text, selection and focus
// handles of one ComboBox, and a revision for option changes a frame cannot
// see.
//
// Revision must change when an option or a section is edited in place, or
// when Filter's results change for the same query. Options appearing or
// disappearing through a new slice or a different length are noticed
// without it.
type ComboBoxState struct {
	// List keeps the popup list's place, highlight and selection.
	List list.ListBoxState
	// Open reports whether the popup shows.
	Open bool
	// Text is the input's text: the value in single mode, a filter query in
	// multiple mode.
	Text string
	// Custom is the text kept by AllowsCustomValue, when it matched no item.
	Custom string
	// Input is the text input.
	Input ui.Handle
	// Revision tells the combo box that options changed in place. Increment
	// it when an option's key, label or other field, or a Filter result,
	// changes without the slice or its length changing.
	Revision uint64

	dirty      bool
	escaped    bool
	wasActive  bool
	placeCaret bool
	pending    pendingFilter
	filtered   filteredOptions
}

// SelectedKeys returns the chosen keys, single choice first.
func (s *ComboBoxState) SelectedKeys() []string {
	keys := s.List.SelectedKeys(list.SelectionMultiple)
	if k := s.List.SelectedKey; k != "" {
		return append([]string{k}, keys...)
	}
	return keys
}

// Clear removes every choice, the custom value and the text.
func (s *ComboBoxState) Clear() {
	s.List.ClearSelection()
	s.Custom = ""
	s.Text = ""
	s.dirty = false
}

// listChosen returns chosen keys and labels in item order. A row with an
// empty key is never chosen, so without a choice there is nothing to walk.
func listChosen(items []list.ListBoxItem, sections []list.ListBoxSection, ls *list.ListBoxState) (keys, labels []string) {
	if ls.SelectedKey == "" && ls.Selected.Len() == 0 {
		return nil, nil
	}
	// Keys are unique, so without multiple choices one row at most matches.
	single := ls.Selected.Len() == 0
	eachSelectItem(items, sections, func(it *list.ListBoxItem) bool {
		key := itemKey(it)
		if key != "" && ls.IsSelected(key) {
			keys = append(keys, key)
			labels = append(labels, it.Label)
			return !single
		}
		return true
	})
	return keys, labels
}

func listNotify(items []list.ListBoxItem, sections []list.ListBoxSection, ls *list.ListBoxState, fn func([]string)) {
	if fn == nil {
		return
	}
	keys, _ := listChosen(items, sections, ls)
	fn(keys)
}

// comboMode resolves an unset mode to single choice.
func comboMode(m list.SelectionMode) list.SelectionMode {
	if m == list.SelectionNone {
		return list.SelectionSingle
	}
	return m
}

// optionSource identifies the options a result was derived from: the Items
// and each section's rows by address and length, and the caller's Revision.
// It keeps the slices themselves, not only their addresses, so the collector
// cannot give a new slice the address of one a cached result describes.
type optionSource struct {
	items        []list.ListBoxItem
	sections     []list.ListBoxSection
	sectionItems [][]list.ListBoxItem
	revision     uint64
}

// same reports whether items, sections and revision are the remembered ones.
func (src *optionSource) same(items []list.ListBoxItem, sections []list.ListBoxSection, revision uint64) bool {
	if revision != src.revision || !sameSlice(items, src.items) ||
		!sameSlice(sections, src.sections) {
		return false
	}
	for i := range sections {
		if !sameSlice(sections[i].Items, src.sectionItems[i]) {
			return false
		}
	}
	return true
}

func (src *optionSource) set(items []list.ListBoxItem, sections []list.ListBoxSection, revision uint64) {
	src.items, src.sections, src.revision = items, sections, revision
	src.sectionItems = src.sectionItems[:0]
	for i := range sections {
		src.sectionItems = append(src.sectionItems, sections[i].Items)
	}
}

// sameSlice reports whether a and b are the same rows of the same array.
func sameSlice[T any](a, b []T) bool {
	return len(a) == len(b) && (len(a) == 0 || &a[0] == &b[0])
}

// filteredOptions is a popup's filtered rows, rebuilt only when the query,
// the options, whether a custom filter is set, or the revision changes.
type filteredOptions struct {
	items        []list.ListBoxItem
	sectionItems []list.ListBoxItem
	sections     []list.ListBoxSection

	source optionSource
	query  string
	custom bool
	valid  bool
}

// update refilters when its inputs changed. A non-nil indicator replaces
// every row's Options.Indicator.
func (f *filteredOptions) update(items []list.ListBoxItem, sections []list.ListBoxSection,
	query string, filter func(string, *list.ListBoxItem) bool, revision uint64, indicator *bool) {
	if f.valid && query == f.query && (filter != nil) == f.custom &&
		f.source.same(items, sections, revision) {
		return
	}
	f.valid, f.query, f.custom = true, query, filter != nil
	f.source.set(items, sections, revision)
	f.filter(items, sections, filterMatch(query, filter), indicator)
}

// filter rebuilds the rows in fresh slices: a copy of the state, or a list
// built from the previous result, may still hold the old ones.
func (f *filteredOptions) filter(items []list.ListBoxItem, sections []list.ListBoxSection,
	match func(*list.ListBoxItem) bool, indicator *bool) {
	f.items = make([]list.ListBoxItem, 0, len(items))
	f.sectionItems = make([]list.ListBoxItem, 0, sectionItemCount(sections))
	f.sections = make([]list.ListBoxSection, 0, len(sections))
	keep := func(dst []list.ListBoxItem, it *list.ListBoxItem) []list.ListBoxItem {
		if !match(it) {
			return dst
		}
		row := *it
		if indicator != nil {
			row.Options.Indicator = indicator
		}
		return append(dst, row)
	}
	for i := range items {
		f.items = keep(f.items, &items[i])
	}
	for si := range sections {
		sec := &sections[si]
		start := len(f.sectionItems)
		for ii := range sec.Items {
			f.sectionItems = keep(f.sectionItems, &sec.Items[ii])
		}
		if len(f.sectionItems) == start {
			continue
		}
		f.sections = append(f.sections, list.ListBoxSection{
			Key:     sec.Key,
			Label:   sec.Label,
			Options: sec.Options,
			Items:   f.sectionItems[start:],
		})
	}
}

// ComboBox renders a combo box. The input keeps the tab stop and the
// keyboard; the popup list highlights the row the arrows are on, and Enter or
// a click chooses it. The list API does not expose the highlighted row as an
// active descendant for assistive technology.
func ComboBox(c *ui.Context, o ComboBoxOptions) ui.Element {
	// An error message implies the invalid state.
	o.Invalid = o.Invalid || o.Error != ""
	s := o.State
	if s == nil {
		panic("form: ComboBox requires a caller-owned State")
	}
	mode := comboMode(o.Selection)
	invalid := o.Invalid
	errText := formstyle.ErrorText(invalid, o.Error)
	if o.Disabled {
		s.Open = false
	}
	// Losing focus closes the popup. A dismissed popup returns focus to the
	// input on the next build; like Escape, that must not reopen it. Focus
	// that stays away for a build clears the flag.
	if !s.Input.Focused(c) && !s.List.FocusWithin(c) {
		if s.wasActive {
			s.Open = false
			s.escaped = true
		} else {
			s.escaped = false
		}
	}

	root := ui.Column(c).MinWidth(0).Gap(theme.Space(1))
	if o.FullWidth {
		root.FillWidth()
	}
	root.Children(func() {
		if o.Label != "" {
			fieldLabel(c, o.Label, o.Required, invalid, o.Disabled, &s.Input)
		}

		field := ui.Row(c).MinHeight(theme.Scale(c, 32)).
			Padding(0, theme.Space(3), 0, theme.Space(3)).Gap(theme.Space(2)).
			AlignItems(ui.Center).Shrink(0).FocusRing(false).Clip()
		if o.FullWidth {
			field.FillWidth()
		}
		if mode == list.SelectionMultiple {
			field.Wrap().Gap(theme.Space(1)).PaddingY(theme.Space(1))
		}
		name := o.Label
		if name == "" {
			name = o.Placeholder
		}
		field.Label(name).Description(o.Description).Error(errText)
		field = fieldSurface(field, c, o.Variant, invalid, o.Disabled)
		shown := busy.Swap(field, o.Loading)

		field.OnClick(func() {
			if o.Disabled {
				return
			}
			s.escaped = false
			if !s.Open {
				comboOpen(s)
			}
			s.Input.Focus(c)
		})
		// Multiple chooses chips in item order, which needs one scan; a
		// single choice keeps its label in Text, so it skips the walk.
		var keys, labels []string
		multiple := mode == list.SelectionMultiple
		if multiple {
			keys, labels = listChosen(o.Items, o.Sections, &s.List)
		}
		var input ui.Element
		field.Children(func() {
			if multiple {
				for i, key := range keys {
					k, label := key, labels[i]
					comboChip(c, label, func() { comboRemove(c, o, s, k) })
				}
			}
			c.Key("combobox-input")
			input = ui.TextInputBase(c, &s.Text).Grow(1).MinWidth(0).
				Shrink(1).FocusRing(false).Placeholder(o.Placeholder).
				Label(name).Description(o.Description).Error(errText).Busy(o.Loading)
			if multiple {
				input.MinWidth(theme.Scale(c, 48))
			}
			input = fieldText(input, c)
			input.Role(ui.RoleComboBox).Expanded(s.Open)
			input.Bind(&s.Input)
			input.OnTextKey(func(k ui.TextKey) bool { return comboTextKey(c, o, s, k, mode) })
			if s.placeCaret {
				s.placeCaret = false
				end := utf8.RuneCountInString(s.Text)
				input.SetTextSelection(end, end)
			}
			if o.Clearable {
				chosen := s.List.SelectedKey != ""
				if multiple {
					chosen = len(keys) > 0
				}
				if chosen || s.Custom != "" || s.Text != "" {
					selectClearButton(c, func() { comboClear(c, o, s) })
					if !s.Open && !o.Disabled {
						input.OnShortcut(0, ui.KeyEscape, func() { comboClear(c, o, s) })
					}
				} else {
					ui.Box(c).Width(theme.Scale(c, 20)).Height(theme.Scale(c, 20)).Shrink(0)
				}
			}
			// Escape closes the popup before a surrounding shortcut sees it.
			if s.Open && !o.Disabled {
				input.OnShortcut(0, ui.KeyEscape, func() { s.Open, s.escaped = false, true })
			}
			selectIndicator(c, shown)
		})
		if input.Valid() {
			focused := input.Focused() || s.Input.FocusWithin(c)
			active := focused || s.List.FocusWithin(c)
			if focused && !s.wasActive && !s.Open && !o.Disabled && !s.escaped {
				comboOpen(s)
			}
			s.wasActive = active
		}
		if input.Changed() {
			s.dirty = true
			s.Open = true
			s.escaped = false
			s.Custom = ""
			s.List.Highlight = 0
			if o.OnQueryChange != nil {
				o.OnQueryChange(s.Text)
			}
		}
		if s.Open {
			query := ""
			if s.dirty {
				query = s.Text
			}
			s.filtered.update(o.Items, o.Sections, query, o.Filter, s.Revision, nil)
			s.pending.apply(&s.List, s.filtered.items, s.filtered.sections)
		} else {
			s.pending.active = false
		}
		style.Disabled(field, o.Disabled)

		comboPopup(c, o, s, field, mode)
		if o.Description != "" && !invalid {
			Description(c, o.Description)
		}
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return root
}

// comboLabel returns the label of key.
func comboLabel(o ComboBoxOptions, key string) string {
	return labelOf(o.Items, o.Sections, key)
}

// labelOf finds the label of key among items or fall back to the key.
func labelOf(items []list.ListBoxItem, sections []list.ListBoxSection, key string) string {
	label := key
	eachSelectItem(items, sections, func(it *list.ListBoxItem) bool {
		if itemKey(it) == key {
			label = it.Label
			return false
		}
		return true
	})
	return label
}

// comboTextKey handles keys the native editor would otherwise consume.
func comboTextKey(c *ui.Context, o ComboBoxOptions, s *ComboBoxState, k ui.TextKey, mode list.SelectionMode) bool {
	if o.Disabled {
		return false
	}
	// Alt+Down opens the popup and Alt+Up closes it, keeping the text.
	if k.Mods == ui.Alt && (k.Key == ui.KeyDown || k.Key == ui.KeyUp) {
		if k.Key == ui.KeyDown && !s.Open {
			comboOpen(s)
			s.escaped = false
		} else if k.Key == ui.KeyUp && s.Open {
			s.Open, s.escaped = false, true
		}
		return true
	}
	if k.Mods != 0 {
		return false
	}
	switch k.Key {
	case ui.KeyDown, ui.KeyUp, ui.KeyEnter:
		if text := k.Text(); text != s.Text {
			s.Text = text
			s.dirty, s.Open, s.escaped = true, true, false
			s.Custom = ""
			if o.OnQueryChange != nil {
				o.OnQueryChange(text)
			}
			k.SetText(s.Text)
			s.pending.reset(s.Text, o.Filter, o.Items, o.Sections)
		}
	}
	switch k.Key {
	case ui.KeyDown, ui.KeyUp:
		if !s.Open {
			comboOpen(s)
			s.escaped = false
		} else if s.pending.active && k.Key == ui.KeyDown {
			s.pending.move(1)
		} else if s.pending.active {
			s.pending.move(-1)
		} else if k.Key == ui.KeyDown {
			s.List.MoveHighlight(1)
		} else {
			s.List.MoveHighlight(-1)
		}
		return true
	case ui.KeyEnter:
		if s.Open {
			key := s.List.HighlightKey()
			if s.pending.active {
				key = s.pending.key()
			}
			if key != "" {
				s.Text = k.Text()
				comboSetSelection(s, key, mode)
				comboCommit(c, o, s, key, mode, false)
				s.pending.active = false
				k.SetText(s.Text)
				listNotify(o.Items, o.Sections, &s.List, o.OnChange)
				return true
			}
		}
		s.Text = k.Text()
		changed := comboCustom(s, o)
		k.SetText(s.Text)
		if changed {
			listNotify(o.Items, o.Sections, &s.List, o.OnChange)
		}
		return true
	}
	return false
}

func comboOpen(s *ComboBoxState) {
	s.Open = true
	if s.List.SelectedKey == "" && s.Text != "" {
		s.dirty = true
	}
	s.List.Highlight = 0
}

// comboSetSelection applies a keyboard choice.
func comboSetSelection(s *ComboBoxState, key string, mode list.SelectionMode) {
	switch mode {
	case list.SelectionSingle:
		s.List.SelectedKey = key
	default:
		s.List.Toggle(mode, key)
	}
}

// comboCommit finalizes a choice the popup list already applied, or one
// applied from the keyboard. It must not mutate the selection again.
func comboCommit(c *ui.Context, o ComboBoxOptions, s *ComboBoxState, key string, mode list.SelectionMode, notify bool) {
	switch mode {
	case list.SelectionSingle:
		s.Text = comboLabel(o, key)
		s.Custom = ""
		s.dirty = false
		s.Open = false
		s.escaped = true
		s.placeCaret = true
		s.Input.Focus(c)
	case list.SelectionMultiple:
		s.Text = ""
		s.Custom = ""
		s.dirty = false
		s.List.Highlight = 0
		s.Input.Focus(c)
	}
	if notify {
		listNotify(o.Items, o.Sections, &s.List, o.OnChange)
	}
}

// comboCustom keeps the text when it matches no item and custom values are
// allowed.
func comboCustom(s *ComboBoxState, o ComboBoxOptions) bool {
	if !o.AllowsCustomValue {
		s.Open = false
		s.escaped = true
		return false
	}
	text := strings.TrimSpace(s.Text)
	if text == "" {
		return false
	}
	s.Custom = text
	s.List.SelectedKey = ""
	s.Open = false
	s.escaped = true
	s.dirty = false
	return true
}

// comboClear clears the choice and the text.
func comboClear(c *ui.Context, o ComboBoxOptions, s *ComboBoxState) {
	empty := s.Text == "" && s.Custom == "" &&
		s.List.SelectionCount(comboMode(o.Selection)) == 0
	if o.Disabled || empty {
		return
	}
	s.Clear()
	s.Open = false
	s.escaped = true
	s.Input.Focus(c)
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
}

// comboRemove takes one chosen key out.
func comboRemove(c *ui.Context, o ComboBoxOptions, s *ComboBoxState, key string) {
	if o.Disabled {
		return
	}
	s.List.SetSelected(list.SelectionMultiple, key, false)
	if s.List.SelectedKey == key {
		s.List.SelectedKey = ""
	}
	s.Input.Focus(c)
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
}

// comboChip draws one chosen value inside the field.
func comboChip(c *ui.Context, label string, remove func()) {
	t := theme.Colors(c)
	radius := theme.RadiusScale().XL
	chip := ui.Row(c).Shrink(0).Center().Height(theme.Scale(c, 20)).
		Gap(theme.Space(1)).Padding(0, theme.Space(1.5), 0, theme.Space(1.5)).
		Radius(radius).Background(t.Accent)
	chip.Children(func() {
		style.Text(ui.Text(c, label), c, theme.XS()).
			TextColor(t.AccentForeground).SingleLine()
		btn := ui.ButtonBase(c).Width(theme.Scale(c, 12)).Height(theme.Scale(c, 12)).
			Shrink(0).Center().Padding(0).Label(locale.M().Remove + " " + label).
			Cursor(ui.CursorPointer).FocusRing(true)
		btn.OnClick(remove)
		btn.Children(func() {
			ui.Icon(c, selectClearIcon).Size(theme.Scale(c, 10), theme.Scale(c, 10)).
				TextColor(t.Muted)
		})
	})
}

// comboPopup builds the anchored option list while it is open.
func comboPopup(c *ui.Context, o ComboBoxOptions, s *ComboBoxState, field ui.Element, mode list.SelectionMode) {
	maxHeight := o.MaxHeight
	if maxHeight == 0 {
		maxHeight = 320
	}
	radius := theme.RadiusScale().LG
	empty := o.Empty
	if empty == "" {
		empty = locale.M().NoResults
	}
	anchor := field
	ui.PopoverBase(c, field, &s.Open, func(panel ui.Element) {
		panel.MinWidth(anchor.Bounds().W)
		panel = anchored.Surface(c, panel.Padding(0).Clip().KeepFocus(), radius)
		anchored.Apply(panel, anchor, o.Placement, mygoui.SideBottom)
		loadingList(c, list.ListBoxOptions{
			Label:             o.Label,
			Items:             s.filtered.items,
			Sections:          s.filtered.sections,
			Selection:         mode,
			State:             &s.List,
			Padding:           theme.Space(1),
			MaxHeight:         maxHeight,
			Empty:             empty,
			NoFocusOnClick:    true,
			KeyboardHighlight: s.Input.Focused(c),
			Editor:            &s.Input,
			IndicatorStart:    true,
			OnAction: func(key string) {
				comboCommit(c, o, s, key, mode, true)
			},
		}, o.Loading, len(s.filtered.items)+sectionItemCount(s.filtered.sections))
	})
}
