package form

import (
	"strings"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AutocompleteOptions configures a HeroUI autocomplete: a field-styled trigger
// opening a searchable, virtualized list.
//
// Loading is the caller's async indicator; the component never performs I/O.
// State is required and caller-owned.
type AutocompleteOptions struct {
	Label             string
	Description       string
	Error             string
	Placeholder       string
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
	Filter func(query string, item *list.ListBoxItem) bool
	// OnChange reports the chosen keys in item order.
	OnChange func(keys []string)
	// OnQueryChange reports each search text change for caller-owned async work.
	OnQueryChange func(query string)
}

// AutocompleteState is the caller-owned open flag, query, selection and focus
// handles of one Autocomplete.
type AutocompleteState struct {
	// List keeps the popup list's place, highlight and selection.
	List list.ListBoxState
	// Open reports whether the popup shows.
	Open bool
	// Query is the search text in the popup.
	Query string
	// Trigger focuses the trigger field.
	Trigger ui.Handle
	// Search focuses the popup's search input.
	Search ui.Handle

	focusSearch  bool
	items        []list.ListBoxItem
	sectionItems []list.ListBoxItem
	sections     []list.ListBoxSection
	filtered     bool
	filterQuery  string
	filterSource listSource
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
}

// applyFilter rebuilds the filtered items and sections in caller-owned scratch
// capacity. An unchanged query over the same option slices reuses the last
// result, so a rebuild for focus or a highlight does not rescan every option.
func (s *AutocompleteState) applyFilter(o AutocompleteOptions, query string) {
	src := sourceOf(o.Items, o.Sections)
	if s.filtered && query == s.filterQuery && src == s.filterSource {
		return
	}
	s.filtered, s.filterQuery, s.filterSource = true, query, src

	totalSections := src.sectionItems
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

// Autocomplete renders a HeroUI autocomplete. The trigger is a field-styled
// button; opening it moves the keyboard into a search field above a virtualized
// list box.
func Autocomplete(c *ui.Context, o AutocompleteOptions) ui.Element {
	s := o.State
	if s == nil {
		panic("form: Autocomplete requires a caller-owned State")
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
			fieldLabel(c, o.Label, o.Required, invalid, o.Disabled, &s.Trigger)
		}
		keys, labels := listChosen(o.Items, o.Sections, &s.List)
		value := strings.Join(labels, ", ")
		placeholder := value == ""
		name := o.Label
		if name == "" {
			name = o.Placeholder
		}

		trigger := ui.ButtonBase(c).MinHeight(theme.Scale(c, theme.HeightInput)).
			Padding(0, theme.Space12, 0, theme.Space12).Gap(theme.Space8).
			AlignItems(ui.Center).Shrink(0).FocusRing(false).
			Role(ui.RolePopUpButton).Expanded(s.Open).Value(value)
		if o.FullWidth {
			trigger.FillWidth()
		}
		trigger.Label(name).Description(o.Description).Error(errText)
		trigger = triggerSurface(trigger, c, o.Variant, invalid, o.Disabled)
		trigger.Bind(&s.Trigger)

		trigger.OnClick(func() {
			if o.Disabled {
				return
			}
			s.Open = !s.Open
			if s.Open {
				autoOpen(o, s)
			}
		})
		trigger.OnShortcut(0, ui.KeyDown, func() { autoOpen(o, s) })
		trigger.OnShortcut(0, ui.KeyUp, func() { autoOpen(o, s) })
		trigger.OnShortcut(0, ui.KeyEnter, func() { autoOpen(o, s) })
		trigger.OnShortcut(0, ui.KeySpace, func() { autoOpen(o, s) })
		if o.Clearable {
			trigger.OnShortcut(0, ui.KeyBackspace, func() { autoClear(c, o, s) })
			trigger.OnShortcut(0, ui.KeyDelete, func() { autoClear(c, o, s) })
		}

		trigger.Children(func() {
			if mode == list.SelectionMultiple {
				for i, key := range keys {
					k, label := key, labels[i]
					comboChip(c, label, func() { autoRemove(c, o, s, k) })
				}
			} else if !placeholder {
				style.Text(ui.Text(c, value), c, theme.SM()).
					TextColor(t.FieldForeground).SingleLine().Grow(1).Shrink(1)
			}
			if placeholder {
				style.Text(ui.Text(c, o.Placeholder), c, theme.SM()).
					TextColor(t.FieldPlaceholder).SingleLine().Grow(1).Shrink(1)
			}
			if o.Clearable {
				if len(keys) > 0 && !o.Disabled {
					selectClearButton(c, func() { autoClear(c, o, s) })
				} else {
					ui.Box(c).Width(theme.Scale(c, 20)).Height(theme.Scale(c, 20)).Shrink(0)
				}
			}
			selectIndicator(c, s.Open)
		})
		style.Disabled(trigger, o.Disabled)

		autoPopup(c, o, s, trigger, mode)
		if o.Description != "" && !invalid {
			Description(c, o.Description)
		}
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return root
}

// autoOpen opens the popup and moves the keyboard to the search field.
func autoOpen(o AutocompleteOptions, s *AutocompleteState) {
	if o.Disabled || s.Open {
		return
	}
	s.Open = true
	s.focusSearch = true
	s.Query = ""
	s.List.Highlight = 0
	if len(s.List.SelectedKey) > 0 {
		s.List.SetHighlight(s.List.SelectedKey)
	}
}

// autoClear clears the choice and the query.
func autoClear(c *ui.Context, o AutocompleteOptions, s *AutocompleteState) {
	if o.Disabled || (s.List.SelectionCount(comboMode(o.Selection)) == 0 && s.Query == "") {
		return
	}
	s.Clear()
	s.Open = false
	s.Trigger.Focus(c)
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
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

// autoChoose applies a keyboard choice; the popup list has not selected it.
func autoChoose(c *ui.Context, o AutocompleteOptions, s *AutocompleteState, key string, mode list.SelectionMode) {
	if mode == list.SelectionSingle {
		s.List.SelectedKey = key
	} else {
		s.List.Toggle(mode, key)
	}
	autoFinalize(c, o, s, key, mode)
}

// autoFinalize closes the popup for a choice the list already applied. It
// must not mutate the selection again.
func autoFinalize(c *ui.Context, o AutocompleteOptions, s *AutocompleteState, key string, mode list.SelectionMode) {
	if mode == list.SelectionSingle {
		s.Open = false
		s.Query = ""
		s.Trigger.Focus(c)
	} else {
		s.Query = ""
		s.List.Highlight = 0
		s.Search.Focus(c)
	}
	listNotify(o.Items, o.Sections, &s.List, o.OnChange)
}

// autoPopup builds the searchable popup while it is open.
func autoPopup(c *ui.Context, o AutocompleteOptions, s *AutocompleteState, trigger ui.Element, mode list.SelectionMode) {
	if !s.Open {
		return
	}
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
	searchName := o.SearchLabel
	if searchName == "" {
		searchName = locale.M().Search
	}
	searchPlaceholder := o.SearchPlaceholder
	if searchPlaceholder == "" {
		searchPlaceholder = searchName
	}
	query := s.Query
	s.applyFilter(o, query)
	anchor := trigger
	ui.PopoverBase(c, trigger, &s.Open, func(panel ui.Element) {
		if w := anchor.Bounds().W; w > 0 {
			panel.MinWidth(w).MaxWidth(w)
		} else {
			panel.MinWidth(theme.Scale(c, 200))
		}
		panel.Padding(0).Radius(radius).Background(t.Overlay).Clip().KeepFocus()
		panel = style.Shadow(panel, c, theme.OverlayShadows(c.Theme().Dark))
		panel = popoverMotion(c, panel, o.Placement, mygoui.SideBottom, 250*time.Millisecond)
		anchored.Apply(panel, anchor, o.Placement, mygoui.SideBottom)
		ui.Column(c).MinWidth(0).Padding(theme.Space8, 0, 0, 0).Children(func() {
			searchRadius := theme.RadiusScale().Field
			search := ui.Row(c).MinHeight(theme.Scale(c, theme.HeightInput)).
				Margin(0, theme.Space8, 0, theme.Space8).
				Padding(0, theme.Space8, 0, theme.Space8).Gap(theme.Space6).
				AlignItems(ui.Center).Radius(searchRadius).Clip().
				Border(theme.Hairline, t.FieldBorder).Background(t.Default)
			var input ui.Element
			search.Children(func() {
				ui.Icon(c, autoSearchIcon).Size(theme.Scale(c, 14), theme.Scale(c, 14)).
					Shrink(0).TextColor(t.Muted)
				input = ui.TextInputBase(c, &s.Query).Grow(1).MinWidth(0).Shrink(1).
					FocusRing(false).Placeholder(searchPlaceholder).
					Label(searchName).Description(o.Description)
				input = fieldText(input, c)
				input.Bind(&s.Search)
				if s.Query != "" {
					selectClearButton(c, func() {
						s.Query = ""
						s.List.Highlight = 0
						s.Search.Focus(c)
						if o.OnQueryChange != nil {
							o.OnQueryChange("")
						}
					})
				}
			})
			if input.Changed() {
				s.List.Highlight = 0
				if o.OnQueryChange != nil {
					o.OnQueryChange(s.Query)
				}
			}
			searchShortcuts(c, o, s, search)
			if s.focusSearch {
				s.focusSearch = false
				s.Search.Focus(c)
			}

			listEmpty := empty
			if o.Loading {
				listEmpty = ""
			}
			list.ListBox(c, list.ListBoxOptions{
				Label:          searchName,
				Items:          s.items,
				Sections:       s.sections,
				Selection:      mode,
				State:          &s.List,
				Padding:        theme.Space6,
				MaxHeight:      maxHeight,
				Empty:          listEmpty,
				NoFocusOnClick: true,
				OnAction: func(key string) {
					autoFinalize(c, o, s, key, mode)
				},
			})
			if o.Loading {
				loadingRow(c)
			}
		})
	})
}

// searchShortcuts drives the list from the popup's search input.
func searchShortcuts(c *ui.Context, o AutocompleteOptions, s *AutocompleteState, search ui.Element) {
	search.OnShortcut(0, ui.KeyDown, func() { s.List.MoveHighlight(1) })
	search.OnShortcut(0, ui.KeyUp, func() { s.List.MoveHighlight(-1) })
	search.OnShortcut(0, ui.KeyEnter, func() {
		if key := s.List.HighlightKey(); key != "" {
			autoChoose(c, o, s, key, comboMode(o.Selection))
		}
	})
	search.OnShortcut(0, ui.KeyEscape, func() {
		s.Open = false
		s.Trigger.Focus(c)
	})
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

var autoSearchIcon = lucide.Search()
