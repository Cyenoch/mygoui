package form

import (
	"unicode/utf8"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AutocompleteOptions configures an editable input with a filtered,
// virtualized popup list.
//
// Loading is the caller's async indicator; the component never performs I/O.
// State is required and caller-owned.
type AutocompleteOptions struct {
	Label       string
	Description string
	Error       string
	Placeholder string
	// SearchPlaceholder and SearchLabel are fallbacks for the editor when
	// Placeholder or Label is empty. There is no separate search field.
	SearchPlaceholder string
	SearchLabel       string
	Empty             string
	Required          bool
	Disabled          bool
	Invalid           bool
	FullWidth         bool
	Clearable         bool
	Loading           bool
	Selection         list.SelectionMode
	Variant           ControlVariant
	Items             []list.ListBoxItem
	Sections          []list.ListBoxSection
	State             *AutocompleteState
	MaxHeight         float32
	Placement         mygoui.Placement

	// Filter overrides the default case-insensitive substring match on Label.
	// It must not perform side effects. It runs again when the query or the
	// options change; bump State.Revision when its results change for the
	// same query.
	Filter func(query string, item *list.ListBoxItem) bool
	// OnChange reports the chosen keys in item order.
	OnChange func(keys []string)
	// OnQueryChange reports user edits and text changes made by choosing or
	// clearing an option, including resets after multiple selection.
	OnQueryChange func(query string)
}

// AutocompleteState is the caller-owned open flag, query, selection and focus
// handles of one Autocomplete, and a revision for option changes a frame
// cannot see.
//
// Revision must change when an option or a section is edited in place, or
// when Filter's results change for the same query. Options appearing or
// disappearing through a new slice or a different length are noticed
// without it.
type AutocompleteState struct {
	// List keeps the popup list's place, highlight and selection.
	List list.ListBoxState
	// Open reports whether the popup shows.
	Open bool
	// Query is the editor text. In single mode a choice replaces it with the
	// chosen label; in multiple mode it is the current filter text.
	Query string
	// Trigger identifies the editable input.
	Trigger ui.Handle
	// Search identifies the same input.
	Search ui.Handle
	// Revision tells the autocomplete that options changed in place.
	// Increment it when an option's key, label or other field, or a Filter
	// result, changes without the slice or its length changing.
	Revision uint64

	dirty           bool
	escaped         bool
	wasActive       bool
	placeCaret      bool
	displayedKey    string
	displayed       string
	displayedSource optionSource
	pending         pendingFilter
	filtered        filteredOptions
}

// SelectedKeys returns the chosen keys, single choice first.
func (s *AutocompleteState) SelectedKeys() []string {
	keys := s.List.SelectedKeys(list.SelectionMultiple)
	if k := s.List.SelectedKey; k != "" {
		return append([]string{k}, keys...)
	}
	return keys
}

// Clear removes every choice and the search text.
func (s *AutocompleteState) Clear() {
	s.List.ClearSelection()
	s.Query = ""
	s.dirty = false
}

// Autocomplete renders a native editable input with an anchored list of
// suggestions. The input retains focus while the popup is navigated. The list
// API does not expose the highlighted row as an active descendant for
// assistive technology.
func Autocomplete(c *ui.Context, o AutocompleteOptions) ui.Element {
	// An error message implies the invalid state.
	o.Invalid = o.Invalid || o.Error != ""
	s := o.State
	if s == nil {
		panic("form: Autocomplete requires a caller-owned State")
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
	if !s.Trigger.Focused(c) && !s.List.FocusWithin(c) {
		if s.wasActive {
			s.Open = false
			s.escaped = true
		} else {
			s.escaped = false
		}
	}
	if mode == list.SelectionSingle {
		key := s.List.SelectedKey
		label := s.displayed
		if key != s.displayedKey || !s.displayedSource.same(o.Items, o.Sections, s.Revision) {
			label = ""
			if key != "" {
				label = labelOf(o.Items, o.Sections, key)
			}
			s.displayedSource.set(o.Items, o.Sections, s.Revision)
		}
		if (key != s.displayedKey || label != s.displayed) && !s.dirty && s.Query == s.displayed {
			s.Query = label
		}
		s.displayedKey, s.displayed = key, label
	}

	root := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	if o.FullWidth {
		root.FillWidth()
	}
	root.Children(func() {
		if o.Label != "" {
			fieldLabel(c, o.Label, o.Required, invalid, o.Disabled, &s.Trigger)
		}
		var keys, labels []string
		if mode == list.SelectionMultiple {
			keys, labels = listChosen(o.Items, o.Sections, &s.List)
		}
		placeholder := o.Placeholder
		if placeholder == "" {
			placeholder = o.SearchPlaceholder
		}
		name := o.Label
		if name == "" {
			name = o.SearchLabel
		}
		if name == "" {
			name = placeholder
		}
		field := ui.Row(c).MinHeight(theme.Scale(c, 32)).
			Padding(0, theme.Space12, 0, theme.Space12).Gap(theme.Space8).
			AlignItems(ui.Center).Shrink(0).FocusRing(false).Clip()
		if o.FullWidth {
			field.FillWidth()
		} else {
			field.MinWidth(theme.Scale(c, inputDefaultWidth))
		}
		if mode == list.SelectionMultiple {
			field.Wrap().Gap(theme.Space4).PaddingY(theme.Space4)
		}
		field = fieldSurface(field, c, o.Variant, invalid, o.Disabled)
		field.OnClick(func() {
			if o.Disabled {
				return
			}
			s.escaped = false
			s.Open = true
			s.Trigger.Focus(c)
		})
		var input ui.Element
		field.Children(func() {
			if mode == list.SelectionMultiple {
				for i, key := range keys {
					k, label := key, labels[i]
					comboChip(c, label, func() { autoRemove(c, o, s, k) })
				}
			}
			c.Key("autocomplete-input")
			input = ui.TextInputBase(c, &s.Query).Grow(1).MinWidth(0).Shrink(1).
				FocusRing(false).Placeholder(placeholder).Label(name).
				Description(o.Description).Error(errText)
			if mode == list.SelectionMultiple {
				input.MinWidth(theme.Scale(c, 48))
			}
			input = fieldText(input, c)
			input.Role(ui.RoleComboBox).Expanded(s.Open).Bind(&s.Trigger).Bind(&s.Search)
			input.OnTextKey(func(k ui.TextKey) bool { return autoTextKey(c, o, s, k, mode) })
			if s.placeCaret {
				s.placeCaret = false
				end := utf8.RuneCountInString(s.Query)
				input.SetTextSelection(end, end)
			}
			if o.Clearable {
				if (s.Query != "" || s.List.SelectionCount(mode) > 0) && !o.Disabled {
					selectClearButton(c, func() { autoClear(c, o, s) })
				} else {
					ui.Box(c).Width(theme.Scale(c, 20)).Height(theme.Scale(c, 20)).Shrink(0)
				}
			}
		})
		focused := s.Trigger.Focused(c)
		active := focused || s.List.FocusWithin(c)
		if focused && !s.wasActive && !s.Open && !o.Disabled && !s.escaped {
			s.Open = true
		}
		s.wasActive = active
		if input.Changed() {
			s.dirty = true
			s.Open = true
			s.escaped = false
			s.List.Highlight = 0
			if o.OnQueryChange != nil {
				o.OnQueryChange(s.Query)
			}
		}
		if s.Open {
			query := s.Query
			if mode == list.SelectionSingle && !s.dirty && s.List.SelectedKey != "" {
				query = ""
			}
			s.filtered.update(o.Items, o.Sections, query, o.Filter, s.Revision, &autoNoIndicator)
			s.pending.apply(&s.List, s.filtered.items, s.filtered.sections)
		} else {
			s.pending.active = false
		}
		autoPopup(c, o, s, field, mode)
		if o.Description != "" && !invalid {
			Description(c, o.Description)
		}
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return root
}

// autoClear clears the choice and the query.
func autoClear(c *ui.Context, o AutocompleteOptions, s *AutocompleteState) {
	if o.Disabled || (s.List.SelectionCount(comboMode(o.Selection)) == 0 && s.Query == "") {
		return
	}
	hadSelection := s.List.SelectionCount(comboMode(o.Selection)) > 0
	hadQuery := s.Query != ""
	s.Clear()
	s.dirty = false
	s.Open = false
	s.escaped = true
	s.Trigger.Focus(c)
	if hadQuery && o.OnQueryChange != nil {
		o.OnQueryChange("")
	}
	if hadSelection {
		listNotify(o.Items, o.Sections, &s.List, o.OnChange)
	}
}

// autoRemove takes one chosen key out.
func autoRemove(c *ui.Context, o AutocompleteOptions, s *AutocompleteState, key string) {
	if o.Disabled {
		return
	}
	s.List.SetSelected(list.SelectionMultiple, key, false)
	if s.List.SelectedKey == key {
		s.List.SelectedKey = ""
	}
	s.Trigger.Focus(c)
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
}

// autoSetSelection applies a keyboard choice.
func autoSetSelection(s *AutocompleteState, key string, mode list.SelectionMode) {
	if mode == list.SelectionSingle {
		s.List.SelectedKey = key
	} else {
		s.List.Toggle(mode, key)
	}
}

// autoFinalize closes the popup for a choice the list already applied. It
// must not mutate the selection again.
func autoFinalize(s *AutocompleteState, o AutocompleteOptions, key string, mode list.SelectionMode) {
	if mode == list.SelectionSingle {
		s.Open = false
		s.escaped = true
		s.Query = labelOf(o.Items, o.Sections, key)
	} else {
		s.Query = ""
		s.List.Highlight = 0
	}
	s.dirty = false
	s.placeCaret = true
}

func autoNotify(o AutocompleteOptions, s *AutocompleteState, previous string) {
	if s.Query != previous && o.OnQueryChange != nil {
		o.OnQueryChange(s.Query)
	}
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
}

// autoPopup builds the suggestions without moving focus away from the editor.
func autoPopup(c *ui.Context, o AutocompleteOptions, s *AutocompleteState, field ui.Element, mode list.SelectionMode) {
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
		if w := anchor.Bounds().W; w > 0 {
			panel.MinWidth(w)
		} else {
			panel.MinWidth(theme.Scale(c, 200))
		}
		panel = anchored.Surface(c, panel.Padding(0).Clip().KeepFocus(), radius)
		anchored.Apply(panel, anchor, o.Placement, mygoui.SideBottom)
		ui.Column(c).MinWidth(0).Children(func() {
			listEmpty := empty
			if o.Loading {
				listEmpty = ""
			}
			list.ListBox(c, list.ListBoxOptions{
				Label:             o.Label,
				Items:             s.filtered.items,
				Sections:          s.filtered.sections,
				Selection:         mode,
				State:             &s.List,
				Padding:           theme.Space4,
				MaxHeight:         maxHeight,
				Empty:             listEmpty,
				NoFocusOnClick:    true,
				KeyboardHighlight: s.Trigger.Focused(c),
				OnAction: func(key string) {
					previous := s.Query
					autoFinalize(s, o, key, mode)
					s.Trigger.Focus(c)
					autoNotify(o, s, previous)
				},
			})
			if o.Loading {
				loadingRow(c)
			}
		})
	})
}

// autoTextKey handles keys the native editor would otherwise consume.
func autoTextKey(c *ui.Context, o AutocompleteOptions, s *AutocompleteState, k ui.TextKey, mode list.SelectionMode) bool {
	if o.Disabled || k.Mods != 0 {
		return false
	}
	switch k.Key {
	case ui.KeyDown, ui.KeyUp, ui.KeyEnter, ui.KeyEscape:
		if text := k.Text(); text != s.Query {
			s.Query = text
			s.dirty, s.Open, s.escaped = true, true, false
			if o.OnQueryChange != nil {
				o.OnQueryChange(text)
			}
			k.SetText(s.Query)
			s.pending.reset(s.Query, o.Filter, o.Items, o.Sections)
		}
	}
	switch k.Key {
	case ui.KeyDown, ui.KeyUp:
		if !s.Open {
			s.Open = true
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
				previous := k.Text()
				autoSetSelection(s, key, mode)
				autoFinalize(s, o, key, mode)
				s.pending.active = false
				k.SetText(s.Query)
				autoNotify(o, s, previous)
				return true
			}
		}
	case ui.KeyEscape:
		if s.Open {
			s.Open = false
			s.escaped = true
			return true
		}
	}
	return false
}

// loadingRow is the caller-owned async indicator.
func loadingRow(c *ui.Context) {
	t := theme.Colors(c)
	ui.Row(c).Gap(theme.Space8).AlignItems(ui.Center).
		Padding(theme.Space8, theme.Space12, theme.Space8, theme.Space12).Children(func() {
		mygoui.Spinner(c, mygoui.SpinnerOptions{Size: mygoui.SpinnerSmall})
		style.Text(ui.Text(c, locale.M().Loading), c, theme.SM()).TextColor(t.Muted)
	})
}

var autoNoIndicator = false
