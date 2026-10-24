package form

import (
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SelectOptions configures a HeroUI select. SelectionSingle is the default;
// SelectionMultiple keeps the popup open while keys are toggled.
//
// State is required and caller-owned. OnChange reports the chosen keys in item
// order after a user change.
type SelectOptions struct {
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
	State       *SelectState
	MaxHeight   float32
	Placement   mygoui.Placement
	OnChange    func(keys []string)
}

// SelectState is the caller-owned open flag, focus handle and list selection
// of one Select.
type SelectState struct {
	// List keeps the popup list's place, highlight and selection.
	List list.ListBoxState
	// Open reports whether the popup shows.
	Open bool
	// Trigger focuses the trigger, as it returns after a choice or a clear.
	Trigger ui.Handle

	focusList bool
}

// SelectedKeys returns the chosen keys, single choice first.
func (s *SelectState) SelectedKeys() []string {
	keys := s.List.SelectedKeys(list.SelectionMultiple)
	if k := s.List.SelectedKey; k != "" {
		return append([]string{k}, keys...)
	}
	return keys
}

// Clear removes every choice.
func (s *SelectState) Clear() { s.List.ClearSelection() }

// selectMode resolves an unset mode to HeroUI's single-choice select.
func selectMode(m list.SelectionMode) list.SelectionMode {
	if m == list.SelectionNone {
		return list.SelectionSingle
	}
	return m
}

// eachSelectItem visits items in display order, stopping when fn returns false.
func eachSelectItem(items []list.ListBoxItem, sections []list.ListBoxSection, fn func(*list.ListBoxItem) bool) {
	for i := range items {
		if !fn(&items[i]) {
			return
		}
	}
	for si := range sections {
		for ii := range sections[si].Items {
			if !fn(&sections[si].Items[ii]) {
				return
			}
		}
	}
}

// selectChosen returns the chosen keys and labels in item order.
func selectChosen(items []list.ListBoxItem, sections []list.ListBoxSection, s *SelectState) (keys, labels []string) {
	eachSelectItem(items, sections, func(it *list.ListBoxItem) bool {
		key := itemKey(it)
		if s.List.IsSelected(key) {
			keys = append(keys, key)
			labels = append(labels, it.Label)
		}
		return true
	})
	return keys, labels
}

func selectNotify(o SelectOptions, s *SelectState) {
	if o.OnChange == nil {
		return
	}
	keys, _ := selectChosen(o.Items, o.Sections, s)
	o.OnChange(keys)
}

// fieldMessage returns the visible and accessible error for a field.
func fieldMessage(invalid bool, text string) string {
	if !invalid {
		return ""
	}
	if text != "" {
		return text
	}
	return locale.M().Invalid
}

// itemKey is a list item's selection identity.
func itemKey(it *list.ListBoxItem) string {
	if it.Key != "" {
		return it.Key
	}
	return it.Label
}

// Select renders a HeroUI select: a field-styled trigger opening a virtualized
// list box. The trigger keeps the tab stop; the list receives the keyboard when
// the popup opens, so its own navigation, typeahead and selection apply.
func Select(c *ui.Context, o SelectOptions) ui.Element {
	s := o.State
	if s == nil {
		panic("form: Select requires a caller-owned State")
	}
	mode := selectMode(o.Selection)
	t := theme.Colors(c)
	invalid := o.Invalid || o.Error != ""
	errText := fieldMessage(invalid, o.Error)

	root := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	if o.FullWidth {
		root.FillWidth()
	}
	root.Children(func() {
		if o.Label != "" {
			fieldLabel(c, o.Label, o.Required, invalid, o.Disabled, &s.Trigger)
		}

		keys, labels := selectChosen(o.Items, o.Sections, s)
		value := strings.Join(labels, ", ")
		placeholder := value == ""

		trigger := ui.ButtonBase(c).MinHeight(theme.Scale(c, theme.HeightInput)).
			Padding(0, theme.Space12, 0, theme.Space12).Gap(theme.Space8).
			AlignItems(ui.Center).Shrink(0).FocusRing(false).
			Role(ui.RolePopUpButton).Expanded(s.Open).Value(value)
		if o.FullWidth {
			trigger.FillWidth()
		}
		name := o.Label
		if name == "" {
			name = o.Placeholder
		}
		trigger.Label(name).Description(o.Description).Error(errText)
		trigger = fieldSurface(trigger, c, o.Variant, invalid, o.Disabled)
		trigger.Bind(&s.Trigger)

		trigger.OnClick(func() {
			if o.Disabled {
				return
			}
			s.Open = !s.Open
			if s.Open {
				s.focusList = true
				if key := selectPreferred(o, s, keys); key != "" {
					s.List.SetHighlight(key)
				}
			}
		})
		trigger.OnShortcut(0, ui.KeyDown, func() { selectOpen(o, s, keys, true) })
		trigger.OnShortcut(0, ui.KeyUp, func() { selectOpen(o, s, keys, false) })
		trigger.OnShortcut(0, ui.KeyEnter, func() { selectOpen(o, s, keys, true) })
		trigger.OnShortcut(0, ui.KeySpace, func() { selectOpen(o, s, keys, true) })
		if o.Clearable {
			trigger.OnShortcut(0, ui.KeyBackspace, func() { selectClear(c, o, s) })
			trigger.OnShortcut(0, ui.KeyDelete, func() { selectClear(c, o, s) })
		}
		if mode == list.SelectionSingle && !o.Disabled {
			for i := range selectTypeKeys {
				r := selectTypeRunes[i]
				trigger.OnShortcut(0, selectTypeKeys[i], func() { selectTypeahead(o, s, r) })
			}
		}

		trigger.Children(func() {
			if !placeholder {
				style.Text(ui.Text(c, value), c, theme.SM()).
					TextColor(t.FieldForeground).SingleLine().Grow(1).Shrink(1)
			} else {
				style.Text(ui.Text(c, o.Placeholder), c, theme.SM()).
					TextColor(t.FieldPlaceholder).SingleLine().Grow(1).Shrink(1)
			}
			if o.Clearable {
				if len(keys) > 0 && !o.Disabled {
					selectClearButton(c, func() { selectClear(c, o, s) })
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
		selectPopup(c, o, s, trigger, mode)
		if o.Description != "" && !invalid {
			Description(c, o.Description)
		}
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return root
}

// selectPreferred returns the key the popup should highlight when it opens.
func selectPreferred(o SelectOptions, s *SelectState, keys []string) string {
	if len(keys) > 0 {
		return keys[0]
	}
	key := ""
	eachSelectItem(o.Items, o.Sections, func(it *list.ListBoxItem) bool {
		if it.Disabled {
			return true
		}
		key = itemKey(it)
		return false
	})
	return key
}

// selectOpen opens the popup from the trigger.
func selectOpen(o SelectOptions, s *SelectState, keys []string, first bool) {
	if o.Disabled || s.Open {
		return
	}
	s.Open = true
	s.focusList = true
	key := selectPreferred(o, s, keys)
	if !first {
		key = selectLast(o)
	}
	if key != "" {
		s.List.SetHighlight(key)
	}
}

// selectLast returns the last selectable item's key.
func selectLast(o SelectOptions) string {
	key := ""
	eachSelectItem(o.Items, o.Sections, func(it *list.ListBoxItem) bool {
		if !it.Disabled {
			key = itemKey(it)
		}
		return true
	})
	return key
}

// selectClear clears the choice and reports it.
func selectClear(c *ui.Context, o SelectOptions, s *SelectState) {
	if o.Disabled || s.List.SelectionCount(selectMode(o.Selection)) == 0 {
		return
	}
	s.List.ClearSelection()
	s.Trigger.Focus(c)
	selectNotify(o, s)
}

// selectTypeahead chooses the next item whose label starts with r.
func selectTypeahead(o SelectOptions, s *SelectState, r rune) {
	if o.Disabled {
		return
	}
	found := ""
	letter := string(r)
	eachSelectItem(o.Items, o.Sections, func(it *list.ListBoxItem) bool {
		if it.Disabled || it.Label == "" {
			return true
		}
		if strings.HasPrefix(strings.ToLower(it.Label), letter) {
			found = itemKey(it)
			return false
		}
		return true
	})
	if found == "" || s.List.SelectedKey == found {
		return
	}
	s.List.SelectedKey = found
	selectNotify(o, s)
}

// selectPopup builds the anchored popup while it is open.
func selectPopup(c *ui.Context, o SelectOptions, s *SelectState, trigger ui.Element, mode list.SelectionMode) {
	maxHeight := o.MaxHeight
	if maxHeight == 0 {
		maxHeight = 320
	}
	t := theme.Colors(c)
	radius := theme.RadiusScale().XXXL
	if radius > 32 {
		radius = 32
	}
	anchor := trigger
	ui.PopoverBase(c, trigger, &s.Open, func(panel ui.Element) {
		panel.MinWidth(anchor.Bounds().W)
		panel.Padding(0).Radius(radius).Background(t.Overlay).Clip()
		panel = style.Shadow(panel, c, theme.OverlayShadows(c.Theme().Dark))
		anchored.Apply(panel, anchor, o.Placement, mygoui.SideBottom)
		list.ListBox(c, list.ListBoxOptions{
			Label:     o.Label,
			Items:     o.Items,
			Sections:  o.Sections,
			Selection: mode,
			State:     &s.List,
			Padding:   theme.Space6,
			MaxHeight: maxHeight,
			Empty:     locale.M().NoResults,
			OnSelectionChange: func(keys []string) {
				selectNotify(o, s)
				if mode == list.SelectionSingle {
					s.Open = false
					s.Trigger.Focus(c)
				}
			},
			OnAction: func(key string) {
				if mode == list.SelectionSingle {
					s.Open = false
					s.Trigger.Focus(c)
				}
			},
		})
		if s.focusList {
			s.focusList = false
			s.List.Focus(c)
		}
	})
}

func selectClearButton(c *ui.Context, clear func()) {
	t := theme.Colors(c)
	radius := theme.RadiusScale().XL
	btn := ui.Row(c).Width(theme.Scale(c, 20)).Height(theme.Scale(c, 20)).
		Shrink(0).Center().Radius(radius).FocusRing(true).
		Role(ui.RoleButton).Label(locale.M().Clear).Cursor(ui.CursorPointer)
	style.Background(btn, radius, ui.Transparent, t.DefaultHover, t.DefaultHover)
	btn.OnClick(clear)
	btn.Children(func() {
		ui.Icon(c, selectClearIcon).Size(theme.Scale(c, 14), theme.Scale(c, 14)).
			TextColor(t.Muted)
	})
}

var (
	selectChevron   = lucide.ChevronDown()
	selectClearIcon = lucide.X()
)

var selectTypeKeys = [...]ui.Key{
	ui.KeyA, ui.KeyB, ui.KeyC, ui.KeyD, ui.KeyE, ui.KeyF, ui.KeyG, ui.KeyH,
	ui.KeyI, ui.KeyJ, ui.KeyK, ui.KeyL, ui.KeyM, ui.KeyN, ui.KeyO, ui.KeyP,
	ui.KeyQ, ui.KeyR, ui.KeyS, ui.KeyT, ui.KeyU, ui.KeyV, ui.KeyW, ui.KeyX,
	ui.KeyY, ui.KeyZ,
}

var selectTypeRunes = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}
