package form

import (
	"strings"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ComboBoxOptions configures a HeroUI combo box: a text input filtering a
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
		if ls.IsSelected(itemKey(it)) {
			keys = append(keys, itemKey(it))
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

// applyFilter rebuilds the filtered items and sections in caller-owned scratch
// capacity, so a frame with the same size does not allocate.
func (s *ComboBoxState) applyFilter(o ComboBoxOptions, query string) {
	totalSections := 0
	for si := range o.Sections {
		totalSections += len(o.Sections[si].Items)
	}
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
	match := func(it *list.ListBoxItem) bool {
		if query == "" {
			return true
		}
		if o.Filter != nil {
			return o.Filter(query, it)
		}
		return strings.Contains(strings.ToLower(it.Label), strings.ToLower(query))
	}
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

// ComboBox renders a HeroUI combo box. The input keeps the tab stop and the
// keyboard; the popup list highlights the row the arrows are on, and Enter or
// a click chooses it.
func ComboBox(c *ui.Context, o ComboBoxOptions) ui.Element {
	s := o.State
	if s == nil {
		panic("form: ComboBox requires a caller-owned State")
	}
	mode := comboMode(o.Selection)
	t := theme.Colors(c)
	invalid := o.Invalid || o.Error != ""
	errText := fieldMessage(invalid, o.Error)

	root := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	if o.FullWidth {
		root.FillWidth()
	}
	root.Children(func() {
		if o.Label != "" {
			fieldLabel(c, o.Label, o.Required, invalid, o.Disabled, &s.Input)
		}

		field := ui.Row(c).MinHeight(theme.Scale(c, theme.HeightInput)).
			Padding(0, theme.Space12, 0, theme.Space12).Gap(theme.Space8).
			AlignItems(ui.Center).Shrink(0).FocusRing(false).Clip()
		if o.FullWidth {
			field.FillWidth()
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
		keys, _ := listChosen(o.Items, o.Sections, &s.List)
		var input ui.Element
		field.Children(func() {
			if mode == list.SelectionMultiple {
				for _, key := range keys {
					k := key
					comboChip(c, comboLabel(o, k), func() { comboRemove(c, o, s, k) })
				}
			}
			input = ui.TextInputBase(c, &s.Text).Grow(1).MinWidth(0).
				Shrink(1).FocusRing(false).Placeholder(o.Placeholder).
				Label(name).Description(o.Description).Error(errText)
			input = fieldText(input, c)
			input.Role(ui.RoleComboBox).Expanded(s.Open)
			input.Bind(&s.Input)
			if o.Clearable {
				if len(keys) > 0 || s.Custom != "" || s.Text != "" {
					selectClearButton(c, func() { comboClear(c, o, s) })
				} else {
					ui.Box(c).Width(theme.Scale(c, 20)).Height(theme.Scale(c, 20)).Shrink(0)
				}
			}
			indicator := ui.Icon(c, selectChevron).Size(theme.Scale(c, 16), theme.Scale(c, 16)).
				Shrink(0).TextColor(t.FieldPlaceholder)
			if s.Open {
				indicator.Rotate(180)
			}
		})
		if input.Valid() {
			focused := input.Focused() || s.Input.FocusWithin(c)
			active := focused || s.List.FocusWithin(c)
			if active && !s.Open && !o.Disabled && !s.escaped {
				comboOpen(s)
			}
			if s.wasActive && !active {
				s.Open = false
				s.escaped = false
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
		}
		style.Disabled(field, o.Disabled)

		comboKeys(c, o, s, field, mode)

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

// comboKeys binds the input's keyboard behavior: the arrows open and move,
// Enter chooses, Escape closes and Backspace clears.
func comboKeys(c *ui.Context, o ComboBoxOptions, s *ComboBoxState, field ui.Element, mode list.SelectionMode) {
	field.OnShortcut(0, ui.KeyDown, func() {
		if o.Disabled {
			return
		}
		if !s.Open {
			comboOpen(s)
		}
		s.List.MoveHighlight(1)
	})
	field.OnShortcut(0, ui.KeyUp, func() {
		if o.Disabled {
			return
		}
		if !s.Open {
			comboOpen(s)
		}
		s.List.MoveHighlight(-1)
	})
	field.OnShortcut(0, ui.KeyEnter, func() {
		if o.Disabled {
			return
		}
		if s.Open {
			if key := s.List.HighlightKey(); key != "" {
				comboChoose(c, o, s, key, mode)
				return
			}
		}
		comboCustom(c, o, s)
	})
	field.OnShortcut(0, ui.KeyEscape, func() {
		if s.Open {
			// Keep the focus from reopening the popup it just closed;
			// a click or typing clears escaped.
			s.Open = false
			s.escaped = true
			s.Input.Focus(c)
		}
	})
}

func comboOpen(s *ComboBoxState) {
	s.Open = true
	s.dirty = true
	s.List.Highlight = 0
}

// comboChoose applies a keyboard choice; the popup list has not selected it.
func comboChoose(c *ui.Context, o ComboBoxOptions, s *ComboBoxState, key string, mode list.SelectionMode) {
	switch mode {
	case list.SelectionSingle:
		s.List.SelectedKey = key
	default:
		s.List.Toggle(mode, key)
	}
	comboCommit(c, o, s, key, mode)
}

// comboCommit finalizes a choice the popup list already applied, or one
// comboChoose just applied. It must not mutate the selection again.
func comboCommit(c *ui.Context, o ComboBoxOptions, s *ComboBoxState, key string, mode list.SelectionMode) {
	switch mode {
	case list.SelectionSingle:
		s.Text = comboLabel(o, key)
		s.Custom = ""
		s.dirty = false
		s.Open = false
		s.Input.Focus(c)
	case list.SelectionMultiple:
		s.Text = ""
		s.Custom = ""
		s.dirty = false
		s.List.Highlight = 0
		s.Input.Focus(c)
	}
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
}

// comboCustom keeps the text when it matches no item and custom values are
// allowed.
func comboCustom(c *ui.Context, o ComboBoxOptions, s *ComboBoxState) {
	if !o.AllowsCustomValue {
		s.Open = false
		return
	}
	text := strings.TrimSpace(s.Text)
	if text == "" {
		return
	}
	s.Custom = text
	s.List.SelectedKey = ""
	s.Open = false
	s.dirty = false
	s.Input.Focus(c)
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
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
		Radius(radius).Background(t.Default)
	chip.Children(func() {
		style.Text(ui.Text(c, label), c, theme.XS()).
			TextColor(t.DefaultForeground).SingleLine()
		btn := ui.Box(c).Width(theme.Scale(c, 12)).Height(theme.Scale(c, 12)).
			Shrink(0).Center().Role(ui.RoleButton).Label(locale.M().Remove).
			Cursor(ui.CursorPointer).FocusRing(false)
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
	radius := theme.RadiusScale().XXXL
	if radius > 32 {
		radius = 32
	}
	empty := o.Empty
	if empty == "" {
		empty = locale.M().NoResults
	}
	anchor := field
	ui.PopoverBase(c, field, &s.Open, func(panel ui.Element) {
		panel.MinWidth(anchor.Bounds().W)
		panel.Padding(0).Radius(radius).Background(t.Overlay).Clip().KeepFocus()
		panel = style.Shadow(panel, c, theme.OverlayShadows(c.Theme().Dark))
		anchored.Apply(panel, anchor, o.Placement, mygoui.SideBottom)
		list.ListBox(c, list.ListBoxOptions{
			Label:          o.Label,
			Items:          s.items,
			Sections:       s.sections,
			Selection:      mode,
			State:          &s.List,
			Padding:        theme.Space6,
			MaxHeight:      maxHeight,
			Empty:          empty,
			NoFocusOnClick: true,
			OnAction: func(key string) {
				comboCommit(c, o, s, key, mode)
			},
		})
	})
}
