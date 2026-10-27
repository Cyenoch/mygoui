package form

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
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
	Selection   list.SelectionMode
	Variant     ControlVariant
	Items       []list.ListBoxItem
	Sections    []list.ListBoxSection
	State       *ComboBoxState
	MaxHeight   float32
	Placement   mygoui.Placement

	// AllowsCustomValue keeps text that matches no option.
	AllowsCustomValue bool
	// Filter overrides the default case-insensitive substring match on Label.
	// Options are scanned on popup builds; Filter must not perform side effects.
	Filter func(query string, item *list.ListBoxItem) bool
	// Empty is shown when nothing matches; zero uses the localized message.
	Empty string
	// OnChange reports the chosen keys in item order.
	OnChange func(keys []string)
	// OnQueryChange reports each text change for caller-owned async work.
	OnQueryChange func(query string)
}

// ComboBoxState is the caller-owned open flag, text, selection and focus
// handles of one ComboBox.
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

	dirty        bool
	escaped      bool
	wasActive    bool
	placeCaret   bool
	pending      pendingFilter
	items        []list.ListBoxItem
	sectionItems []list.ListBoxItem
	sections     []list.ListBoxSection
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

// listChosen returns chosen keys and labels in item order.
func listChosen(items []list.ListBoxItem, sections []list.ListBoxSection, ls *list.ListBoxState) (keys, labels []string) {
	eachSelectItem(items, sections, func(it *list.ListBoxItem) bool {
		key := itemKey(it)
		if ls.IsSelected(key) {
			keys = append(keys, key)
			labels = append(labels, it.Label)
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

// applyFilter rebuilds filtered options in caller-owned scratch capacity.
func (s *ComboBoxState) applyFilter(o ComboBoxOptions, query string) {
	totalSections := sectionItemCount(o.Sections)
	if cap(s.items) < len(o.Items) {
		s.items = make([]list.ListBoxItem, 0, len(o.Items))
	}
	s.items = s.items[:0]
	if cap(s.sectionItems) < totalSections {
		s.sectionItems = make([]list.ListBoxItem, 0, totalSections)
	}
	s.sectionItems = s.sectionItems[:0]
	if cap(s.sections) < len(o.Sections) {
		s.sections = make([]list.ListBoxSection, 0, len(o.Sections))
	}
	s.sections = s.sections[:0]
	match := filterMatch(query, o.Filter)
	for i := range o.Items {
		if match(&o.Items[i]) {
			s.items = append(s.items, o.Items[i])
		}
	}
	for si := range o.Sections {
		sec := &o.Sections[si]
		start := len(s.sectionItems)
		for ii := range sec.Items {
			if match(&sec.Items[ii]) {
				s.sectionItems = append(s.sectionItems, sec.Items[ii])
			}
		}
		if len(s.sectionItems) == start {
			continue
		}
		s.sections = append(s.sections, list.ListBoxSection{
			Key:     sec.Key,
			Label:   sec.Label,
			Options: sec.Options,
			Items:   s.sectionItems[start:],
		})
	}
}

// ComboBox renders a combo box. The input keeps the tab stop and the
// keyboard; the popup list highlights the row the arrows are on, and Enter or
// a click chooses it. The list API does not expose the highlighted row as an
// active descendant for assistive technology.
func ComboBox(c *ui.Context, o ComboBoxOptions) ui.Element {
	s := o.State
	if s == nil {
		panic("form: ComboBox requires a caller-owned State")
	}
	mode := comboMode(o.Selection)
	invalid := o.Invalid || o.Error != ""
	errText := fieldMessage(invalid, o.Error)
	if o.Disabled {
		s.Open = false
	}
	if s.wasActive && !s.Input.Focused(c) && !s.List.FocusWithin(c) {
		s.Open = false
		s.escaped = false
	}

	root := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	if o.FullWidth {
		root.FillWidth()
	}
	root.Children(func() {
		if o.Label != "" {
			fieldLabel(c, o.Label, o.Required, invalid, o.Disabled, &s.Input)
		}

		field := ui.Row(c).MinHeight(theme.Scale(c, 32)).
			Padding(0, theme.Space12, 0, theme.Space12).Gap(theme.Space8).
			AlignItems(ui.Center).Shrink(0).FocusRing(false).Clip()
		if o.FullWidth {
			field.FillWidth()
		}
		if mode == list.SelectionMultiple {
			field.Wrap().Gap(theme.Space4).PaddingY(theme.Space4)
		}
		name := o.Label
		if name == "" {
			name = o.Placeholder
		}
		field.Label(name).Description(o.Description).Error(errText)
		field = fieldSurface(field, c, o.Variant, invalid, o.Disabled)

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
				Label(name).Description(o.Description).Error(errText)
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
				} else {
					ui.Box(c).Width(theme.Scale(c, 20)).Height(theme.Scale(c, 20)).Shrink(0)
				}
			}
			selectIndicator(c, s.Open)
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
			s.applyFilter(o, query)
			s.pending.apply(&s.List, s.items, s.sections)
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
	if o.Disabled || k.Mods != 0 {
		return false
	}
	switch k.Key {
	case ui.KeyDown, ui.KeyUp, ui.KeyEnter, ui.KeyEscape:
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
	case ui.KeyEscape:
		if s.Open {
			s.Open = false
			s.escaped = true
			return true
		}
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
		Gap(theme.Space4).Padding(0, theme.Space6, 0, theme.Space6).
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
	t := theme.Colors(c)
	radius := theme.RadiusScale().LG
	empty := o.Empty
	if empty == "" {
		empty = locale.M().NoResults
	}
	anchor := field
	ui.PopoverBase(c, field, &s.Open, func(panel ui.Element) {
		panel.MinWidth(anchor.Bounds().W)
		panel.Padding(0).Radius(radius).Background(t.Popover).Border(theme.Hairline, t.Border).Clip().KeepFocus()
		panel = style.Shadow(panel, c, theme.OverlayShadows(c.Theme().Dark))
		panel = popoverMotion(c, panel, o.Placement, mygoui.SideBottom, 150*time.Millisecond)
		placement := o.Placement
		if placement.Offset == 0 {
			placement.Offset = 4
		}
		anchored.Apply(panel, anchor, placement, mygoui.SideBottom)
		list.ListBox(c, list.ListBoxOptions{
			Label:             o.Label,
			Items:             s.items,
			Sections:          s.sections,
			Selection:         mode,
			State:             &s.List,
			Padding:           theme.Space4,
			MaxHeight:         maxHeight,
			Empty:             empty,
			NoFocusOnClick:    true,
			KeyboardHighlight: s.Input.Focused(c),
			IndicatorStart:    true,
			OnAction: func(key string) {
				comboCommit(c, o, s, key, mode, true)
			},
		})
	})
}
