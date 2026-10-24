package gallery

import (
	"fmt"
	"strings"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type selectionState struct {
	ready bool

	fruitList   list.ListBoxState
	tagList     list.ListBoxState
	sectionList list.ListBoxState
	disabled    list.ListBoxState
	empty       list.ListBoxState

	fruitSelect   form.SelectState
	tagSelect     form.SelectState
	clearSelect   form.SelectState
	sectionSelect form.SelectState
	variantSelect form.SelectState

	combo       form.ComboBoxState
	comboMulti  form.ComboBoxState
	comboCustom form.ComboBoxState
	comboFlat   form.ComboBoxState

	auto      form.AutocompleteState
	autoEmpty form.AutocompleteState
	autoLoad  form.AutocompleteState
	autoFlat  form.AutocompleteState

	listNote   string
	selectNote string
	comboNote  string
	autoNote   string
}

func (st *selectionState) initState() {
	if st.ready {
		return
	}
	st.ready = true

	st.fruitList.SelectedKey = "banana"
	st.tagList.Selected.Add("go")
	st.tagList.Selected.Add("typescript")
	st.sectionList.SelectedKey = "go"

	st.fruitSelect.List.SelectedKey = "apple"
	st.tagSelect.List.Selected.Add("go")
	st.tagSelect.List.Selected.Add("rust")
	st.clearSelect.List.SelectedKey = "typescript"
	st.sectionSelect.List.SelectedKey = "go"

	st.combo.List.SelectedKey = "go"
	st.combo.Text = "Go"
	st.comboMulti.List.Selected.Add("go")
	st.comboMulti.List.Selected.Add("rust")
	st.auto.List.SelectedKey = "paris"

	st.listNote = "banana"
	st.selectNote = "apple"
	st.comboNote = "Go"
	st.autoNote = "Paris"
}

// selectionFruits, selectionTags and the other panels are immutable item
// definitions shared by the specimens; all mutable state lives in
// selectionState.
var selectionFruits = []list.ListBoxItem{
	{Key: "apple", Label: "Apple", Description: "Crisp and sweet."},
	{Key: "banana", Label: "Banana"},
	{Key: "cherry", Label: "Cherry", Description: "Two varieties in season."},
	{Key: "dragonfruit", Label: "Dragon fruit", Description: "Unavailable this week.", Disabled: true},
}

var selectionTags = []list.ListBoxItem{
	{Key: "go", Label: "Go"},
	{Key: "typescript", Label: "TypeScript"},
	{Key: "rust", Label: "Rust"},
	{Key: "zig", Label: "Zig"},
	{Key: "elixir", Label: "Elixir"},
}

var selectionFrameworks = []list.ListBoxItem{
	{Key: "go", Label: "Go", Description: "Standard library first."},
	{Key: "rust", Label: "Rust"},
	{Key: "typescript", Label: "TypeScript"},
	{Key: "react", Label: "React"},
	{Key: "vue", Label: "Vue"},
	{Key: "svelte", Label: "Svelte"},
}

var selectionSections = []list.ListBoxSection{
	{
		Key:   "frontend",
		Label: "Frontend",
		Items: []list.ListBoxItem{
			{Key: "react", Label: "React"},
			{Key: "vue", Label: "Vue"},
			{Key: "svelte", Label: "Svelte"},
		},
	},
	{
		Key:   "backend",
		Label: "Backend",
		Items: []list.ListBoxItem{
			{Key: "go", Label: "Go"},
			{Key: "rust", Label: "Rust"},
			{Key: "elixir", Label: "Elixir", Description: "Not on this plan.", Disabled: true},
		},
	},
	{
		Key:   "tools",
		Label: "Tools",
		Items: []list.ListBoxItem{
			{Key: "vite", Label: "Vite"},
			{Key: "biome", Label: "Biome"},
		},
	},
}

var selectionCities = []list.ListBoxItem{
	{Key: "paris", Label: "Paris", Description: "France"},
	{Key: "tokyo", Label: "Tokyo", Description: "Japan"},
	{Key: "lisbon", Label: "Lisbon", Description: "Portugal"},
	{Key: "nairobi", Label: "Nairobi", Description: "Kenya"},
	{Key: "bogota", Label: "Bogotá", Description: "Colombia"},
}

// selectionNote renders the live value reported by a specimen's state.
func selectionNote(c *ui.Context, prefix, value string) {
	if value == "" {
		value = "none"
	}
	mygoui.Chip(c, mygoui.ChipOptions{
		Label:   prefix + value,
		Color:   mygoui.ColorAccent,
		Variant: mygoui.ChipSoft,
		Size:    mygoui.ChipSmall,
	})
}

// selectionExampleDemo adapts one single-variant selection specimen to a
// named ComponentExample demo. The fixed column keeps triggers at a
// readable width inside the example panel.
func selectionExampleDemo(body func(c *ui.Context, st *selectionState)) func(c *ui.Context, s *State) {
	return func(c *ui.Context, s *State) {
		st := selectionPageState(s)
		ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 340)).
			Gap(theme.Space8).Children(func() { body(c, st) })
	}
}

// selectionSelectExampleDefault shows a single-choice trigger.
func selectionSelectExampleDefault(c *ui.Context, st *selectionState) {
	form.Select(c, form.SelectOptions{
		Label:       "Fruit",
		Placeholder: "Choose a fruit",
		Items:       selectionFruits,
		State:       &st.fruitSelect,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.selectNote = strings.Join(keys, ", ")
		},
	})
	selectionNote(c, "Current fruit: ", st.fruitSelect.List.SelectedKey)
}

// selectionSelectExampleSecondary shows the flat secondary trigger face.
func selectionSelectExampleSecondary(c *ui.Context, st *selectionState) {
	form.Select(c, form.SelectOptions{
		Label:       "Fruit",
		Placeholder: "Choose a fruit",
		Items:       selectionFruits,
		State:       &st.variantSelect,
		Variant:     form.ControlSecondary,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.selectNote = strings.Join(keys, ", ")
		},
	})
	selectionNote(c, "Current fruit: ", st.variantSelect.List.SelectedKey)
}

// selectionSelectExampleMultiple shows removable multiple choices.
func selectionSelectExampleMultiple(c *ui.Context, st *selectionState) {
	form.Select(c, form.SelectOptions{
		Label:       "Tags",
		Placeholder: "Choose tags",
		Items:       selectionTags,
		Selection:   list.SelectionMultiple,
		State:       &st.tagSelect,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.selectNote = strings.Join(keys, ", ")
		},
	})
	selectionNote(c, "Tags: ",
		fmt.Sprintf("%d selected", st.tagSelect.List.SelectionCount(list.SelectionMultiple)))
}

// selectionSelectExampleClearable shows the clear action on a chosen value.
func selectionSelectExampleClearable(c *ui.Context, st *selectionState) {
	form.Select(c, form.SelectOptions{
		Label:       "Framework",
		Placeholder: "Choose a framework",
		Items:       selectionFrameworks,
		Clearable:   true,
		State:       &st.clearSelect,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.selectNote = strings.Join(keys, ", ")
		},
	})
}

// selectionSelectExampleSections shows grouped choices in the popup.
func selectionSelectExampleSections(c *ui.Context, st *selectionState) {
	form.Select(c, form.SelectOptions{
		Label:       "Service",
		Placeholder: "Choose a service",
		Sections:    selectionSections,
		State:       &st.sectionSelect,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.selectNote = strings.Join(keys, ", ")
		},
	})
}

// selectionComboBoxExampleDefault shows the filtering single-choice field.
func selectionComboBoxExampleDefault(c *ui.Context, st *selectionState) {
	form.ComboBox(c, form.ComboBoxOptions{
		Label:       "Language",
		Placeholder: "Pick a language",
		Items:       selectionFrameworks,
		Clearable:   true,
		State:       &st.combo,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.comboNote = strings.Join(keys, ", ")
		},
	})
}

// selectionComboBoxExampleSecondary shows the flat secondary field face.
func selectionComboBoxExampleSecondary(c *ui.Context, st *selectionState) {
	form.ComboBox(c, form.ComboBoxOptions{
		Label:       "Language",
		Placeholder: "Pick a language",
		Items:       selectionFrameworks,
		Clearable:   true,
		State:       &st.comboFlat,
		Variant:     form.ControlSecondary,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.comboNote = strings.Join(keys, ", ")
		},
	})
}

// selectionComboBoxExampleMultiple shows removable chips for each choice.
func selectionComboBoxExampleMultiple(c *ui.Context, st *selectionState) {
	form.ComboBox(c, form.ComboBoxOptions{
		Label:       "Stack",
		Placeholder: "Add technologies",
		Items:       selectionFrameworks,
		Selection:   list.SelectionMultiple,
		Clearable:   true,
		State:       &st.comboMulti,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.comboNote = strings.Join(keys, ", ")
		},
	})
}

// selectionComboBoxExampleCustom keeps typed text that matches no option.
func selectionComboBoxExampleCustom(c *ui.Context, st *selectionState) {
	form.ComboBox(c, form.ComboBoxOptions{
		Label:             "Tool",
		Description:       "Type any name and press Enter to keep it.",
		Placeholder:       "Type a tool",
		Items:             selectionFrameworks,
		AllowsCustomValue: true,
		State:             &st.comboCustom,
		FullWidth:         true,
		MaxHeight:         theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.comboNote = strings.Join(keys, ", ")
		},
	})
	selectionNote(c, "Custom value: ", st.comboCustom.Custom)
}

// selectionAutocompleteExampleFiltering shows the searchable popup list.
func selectionAutocompleteExampleFiltering(c *ui.Context, st *selectionState) {
	form.Autocomplete(c, form.AutocompleteOptions{
		Label:             "Destination",
		Placeholder:       "Choose a city",
		SearchPlaceholder: "Filter cities",
		SearchLabel:       "Filter destinations",
		Items:             selectionCities,
		Clearable:         true,
		State:             &st.auto,
		FullWidth:         true,
		MaxHeight:         theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.autoNote = strings.Join(keys, ", ")
		},
	})
}

// selectionAutocompleteExampleSecondary shows the flat secondary trigger.
func selectionAutocompleteExampleSecondary(c *ui.Context, st *selectionState) {
	form.Autocomplete(c, form.AutocompleteOptions{
		Label:             "Destination",
		Placeholder:       "Choose a city",
		SearchPlaceholder: "Filter cities",
		SearchLabel:       "Filter destinations",
		Items:             selectionCities,
		Clearable:         true,
		State:             &st.autoFlat,
		Variant:           form.ControlSecondary,
		FullWidth:         true,
		MaxHeight:         theme.Scale(c, 220),
		OnChange: func(keys []string) {
			st.autoNote = strings.Join(keys, ", ")
		},
	})
}

// selectionAutocompleteExampleEmpty shows the message for no matches.
func selectionAutocompleteExampleEmpty(c *ui.Context, st *selectionState) {
	form.Autocomplete(c, form.AutocompleteOptions{
		Label:       "Saved destination",
		Placeholder: "No destinations",
		Empty:       "No destinations yet",
		State:       &st.autoEmpty,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 160),
	})
}

// selectionAutocompleteExampleLoading shows the caller-owned loading state.
func selectionAutocompleteExampleLoading(c *ui.Context, st *selectionState) {
	form.Autocomplete(c, form.AutocompleteOptions{
		Label:       "Remote search",
		Placeholder: "Open to load",
		Empty:       "No matches",
		Loading:     true,
		State:       &st.autoLoad,
		FullWidth:   true,
		MaxHeight:   theme.Scale(c, 160),
		OnQueryChange: func(query string) {
			st.autoNote = "Query: " + query
		},
	})
}

// selectionListBoxExampleSingle shows single selection with a disabled row.
func selectionListBoxExampleSingle(c *ui.Context, st *selectionState) {
	list.ListBox(c, list.ListBoxOptions{
		Label:     "Fruit",
		Items:     selectionFruits,
		Selection: list.SelectionSingle,
		State:     &st.fruitList,
		MaxHeight: theme.Scale(c, 220),
		OnSelectionChange: func(keys []string) {
			st.listNote = strings.Join(keys, ", ")
		},
	})
	selectionNote(c, "Fruit: ", st.fruitList.SelectedKey)
}

// selectionListBoxExampleMultiple shows multiple selection.
func selectionListBoxExampleMultiple(c *ui.Context, st *selectionState) {
	list.ListBox(c, list.ListBoxOptions{
		Label:     "Tags",
		Items:     selectionTags,
		Selection: list.SelectionMultiple,
		State:     &st.tagList,
		MaxHeight: theme.Scale(c, 220),
		OnSelectionChange: func(keys []string) {
			st.listNote = fmt.Sprintf("%d tags", len(keys))
		},
	})
	selectionNote(c, "Tags: ",
		fmt.Sprintf("%d selected", st.tagList.SelectionCount(list.SelectionMultiple)))
}

// selectionListBoxExampleSections shows labelled groups in order.
func selectionListBoxExampleSections(c *ui.Context, st *selectionState) {
	list.ListBox(c, list.ListBoxOptions{
		Label:     "Services",
		Sections:  selectionSections,
		Selection: list.SelectionSingle,
		State:     &st.sectionList,
		MaxHeight: theme.Scale(c, 240),
		OnSelectionChange: func(keys []string) {
			st.listNote = strings.Join(keys, ", ")
		},
	})
}

// selectionListBoxExampleDisabled dims the list and stops activation.
func selectionListBoxExampleDisabled(c *ui.Context, st *selectionState) {
	list.ListBox(c, list.ListBoxOptions{
		Label:     "Archived labels",
		Items:     selectionFruits[:2],
		Selection: list.SelectionSingle,
		State:     &st.disabled,
		Disabled:  true,
		MaxHeight: theme.Scale(c, 160),
	})
}

// selectionListBoxExampleEmpty states why the list has no rows.
func selectionListBoxExampleEmpty(c *ui.Context, st *selectionState) {
	list.ListBox(c, list.ListBoxOptions{
		Label:     "Saved filters",
		Selection: list.SelectionMultiple,
		State:     &st.empty,
		MaxHeight: theme.Scale(c, 160),
		Empty:     "No saved filters yet",
	})
}
