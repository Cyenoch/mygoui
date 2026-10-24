package gallery

// selectionPageState keeps the shared selection specimen state and fills its
// defaults once, so every selection page keeps its chosen values.
func selectionPageState(s *State) *selectionState {
	st := StateOf[selectionState](s, "selection")
	st.initState()
	return st
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "select", Title: "Select", Group: "Forms", Order: 200,
		Description: "A trigger that opens a virtualized list of single or multiple choices.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "A single-choice trigger with a chosen fruit.", Demo: selectionExampleDemo(selectionSelectExampleDefault)},
			{Title: "Secondary", Description: "The flat secondary trigger face.", Demo: selectionExampleDemo(selectionSelectExampleSecondary)},
			{Title: "Multiple", Description: "Several choices show as removable values.", Demo: selectionExampleDemo(selectionSelectExampleMultiple)},
			{Title: "Clearable", Description: "The clear action empties the trigger.", Demo: selectionExampleDemo(selectionSelectExampleClearable)},
			{Title: "Sections", Description: "Grouped choices keep their headers in the popup.", Demo: selectionExampleDemo(selectionSelectExampleSections)},
		},
		Code: `form.Select(c, form.SelectOptions{
    Label: "Fruit", Placeholder: "Choose a fruit",
    Items: fruits, State: &fruit, FullWidth: true,
    OnChange: func(keys []string) { report(keys) },
})`,
		Params: []ComponentParam{
			{"Placeholder", "string", "", "Hint shown while nothing is chosen."},
			{"Variant", "form.ControlVariant", "ControlPrimary", "Primary field face or the flat secondary face."},
			{"Selection", "list.SelectionMode", "SelectionSingle", "Single or multiple choice."},
			{"Clearable", "bool", "false", "Shows the clear action when something is chosen."},
			{"MaxHeight", "float32", "0", "Caps the popup list's height."},
			{"Disabled", "bool", "false", "Stops opening and selecting."},
			{"Invalid", "bool", "false", "Paints the danger trigger and error."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldPlaceholder", "AccentSoft", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "combobox", Title: "ComboBox", Group: "Forms", Order: 210,
		Description: "A text input that filters its options, with chips and custom values.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "Typing filters the options of a single choice.", Demo: selectionExampleDemo(selectionComboBoxExampleDefault)},
			{Title: "Secondary", Description: "The flat secondary field face.", Demo: selectionExampleDemo(selectionComboBoxExampleSecondary)},
			{Title: "Multiple", Description: "Each choice becomes a removable chip.", Demo: selectionExampleDemo(selectionComboBoxExampleMultiple)},
			{Title: "Custom value", Description: "Typed text that matches no option survives.", Demo: selectionExampleDemo(selectionComboBoxExampleCustom)},
		},
		Code: `form.ComboBox(c, form.ComboBoxOptions{
    Label: "Language", Placeholder: "Pick a language",
    Items: languages, Clearable: true, State: &language, FullWidth: true,
    OnChange: func(keys []string) { report(keys) },
})`,
		Params: []ComponentParam{
			{"Variant", "form.ControlVariant", "ControlPrimary", "Primary field face or the flat secondary face."},
			{"Selection", "list.SelectionMode", "SelectionSingle", "Single choice or removable multiple chips."},
			{"AllowsCustomValue", "bool", "false", "Keeps typed text that matches no option."},
			{"Clearable", "bool", "false", "Shows the clear action."},
			{"Empty", "string", "localized", "Message when the filter matches nothing."},
			{"MaxHeight", "float32", "0", "Caps the popup list's height."},
			{"OnQueryChange", "func(string)", "nil", "Reports each filter query for caller-owned async work."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldPlaceholder", "DefaultHover", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "autocomplete", Title: "Autocomplete", Group: "Forms", Order: 220,
		Description: "A field-styled trigger that opens a searchable list of options.",
		Examples: []ComponentExample{
			{Title: "Filtering", Description: "The popup search narrows the option list.", Demo: selectionExampleDemo(selectionAutocompleteExampleFiltering)},
			{Title: "Secondary", Description: "The flat secondary trigger face.", Demo: selectionExampleDemo(selectionAutocompleteExampleSecondary)},
			{Title: "Empty", Description: "The list states why it has no matches.", Demo: selectionExampleDemo(selectionAutocompleteExampleEmpty)},
			{Title: "Loading", Description: "The caller-owned loading state is shown.", Demo: selectionExampleDemo(selectionAutocompleteExampleLoading)},
		},
		Code: `form.Autocomplete(c, form.AutocompleteOptions{
    Label: "Destination", Placeholder: "Choose a city",
    SearchPlaceholder: "Filter cities", Items: cities,
    State: &city, FullWidth: true,
    OnChange: func(keys []string) { report(keys) },
})`,
		Params: []ComponentParam{
			{"Variant", "form.ControlVariant", "ControlPrimary", "Primary field face or the flat secondary face."},
			{"SearchPlaceholder", "string", "", "Hint in the popup's search field."},
			{"SearchLabel", "string", "", "Accessible name for the popup search."},
			{"Empty", "string", "localized", "Message when nothing matches."},
			{"Loading", "bool", "false", "Shows the caller-owned loading state."},
			{"Clearable", "bool", "false", "Shows the clear action."},
			{"MaxHeight", "float32", "0", "Caps the popup list's height."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldPlaceholder", "DefaultHover", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "listbox", Title: "ListBox", Group: "Collections", Order: 10,
		Description: "A virtualized collection of selectable rows, with sections and empty state.",
		Examples: []ComponentExample{
			{Title: "Single", Description: "One row is selected; disabled rows never activate.", Demo: selectionExampleDemo(selectionListBoxExampleSingle)},
			{Title: "Multiple", Description: "Several rows are selected at once.", Demo: selectionExampleDemo(selectionListBoxExampleMultiple)},
			{Title: "Sections", Description: "Labelled groups appear in order.", Demo: selectionExampleDemo(selectionListBoxExampleSections)},
			{Title: "Disabled", Description: "The list dims and stops activating rows.", Demo: selectionExampleDemo(selectionListBoxExampleDisabled)},
			{Title: "Empty", Description: "An empty list states why it has no rows.", Demo: selectionExampleDemo(selectionListBoxExampleEmpty)},
		},
		Code: `list.ListBox(c, list.ListBoxOptions{
    Label: "Fruit", Items: fruits, Selection: list.SelectionSingle,
    State: &fruit, MaxHeight: theme.Scale(c, 220),
    OnSelectionChange: func(keys []string) { report(keys) },
})`,
		Params: []ComponentParam{
			{"Items", "[]list.ListBoxItem", "nil", "Ungrouped rows shown before sections."},
			{"Sections", "[]list.ListBoxSection", "nil", "Labelled groups shown in order."},
			{"Selection", "list.SelectionMode", "SelectionNone", "None, single or multiple selection."},
			{"Variant", "list.ListBoxVariant", "ListBoxDefault", "Default or danger row color."},
			{"Disabled", "bool", "false", "Dims the list and stops activation."},
			{"Empty", "string", "", "Message shown when there are no rows."},
		},
		Tokens: []string{"Default", "Foreground", "Muted", "Accent", "Danger"},
	})
}
