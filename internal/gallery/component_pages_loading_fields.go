package gallery

import (
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// fieldLoadingState drives the Loading example of each field: one switch sets
// the flag, and the popup fields offer their options only once it is off.
type fieldLoadingState struct {
	loading bool
	text    string
	search  string
	sel     form.SelectState
	combo   form.ComboBoxState
	auto    form.AutocompleteState
}

// fieldLoadingDemo adapts one field body to the Loading example: the switch
// above the field flips the flag the field shows.
func fieldLoadingDemo(body func(c *ui.Context, st *fieldLoadingState, items []list.ListBoxItem)) func(c *ui.Context, s *State) {
	return func(c *ui.Context, s *State) {
		st := StateOf[fieldLoadingState](s, "field-loading")
		items := selectionFruits
		if st.loading {
			items = nil
		}
		ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 340)).Gap(theme.Space12).Children(func() {
			form.Switch(c, form.SwitchOptions{On: &st.loading, Label: "Loading"})
			body(c, st, items)
		})
	}
}

var (
	inputLoadingDemo = fieldLoadingDemo(func(c *ui.Context, st *fieldLoadingState, _ []list.ListBoxItem) {
		form.Input(c, &st.text, form.InputOptions{
			Label: "Username", Placeholder: "Checking availability", FullWidth: true, Loading: st.loading,
		})
	})
	textFieldLoadingDemo = fieldLoadingDemo(func(c *ui.Context, st *fieldLoadingState, _ []list.ListBoxItem) {
		form.TextField(c, form.TextFieldOptions{
			Value: &st.text, Label: "Username", Placeholder: "Checking availability", FullWidth: true, Loading: st.loading,
		})
	})
	searchFieldLoadingDemo = fieldLoadingDemo(func(c *ui.Context, st *fieldLoadingState, _ []list.ListBoxItem) {
		form.SearchField(c, &st.search, form.SearchFieldOptions{
			Label: "Search", Placeholder: "Search by name", FullWidth: true, Loading: st.loading,
		})
	})
	selectLoadingDemo = fieldLoadingDemo(func(c *ui.Context, st *fieldLoadingState, items []list.ListBoxItem) {
		form.Select(c, form.SelectOptions{
			Label: "Fruit", Placeholder: "Choose a fruit", Items: items, State: &st.sel,
			FullWidth: true, Loading: st.loading,
		})
	})
	comboBoxLoadingDemo = fieldLoadingDemo(func(c *ui.Context, st *fieldLoadingState, items []list.ListBoxItem) {
		form.ComboBox(c, form.ComboBoxOptions{
			Label: "Fruit", Placeholder: "Pick a fruit", Items: items, State: &st.combo,
			FullWidth: true, Loading: st.loading,
		})
	})
	autocompleteLoadingDemo = fieldLoadingDemo(func(c *ui.Context, st *fieldLoadingState, items []list.ListBoxItem) {
		form.Autocomplete(c, form.AutocompleteOptions{
			Label: "Fruit", Placeholder: "Search fruit", Items: items, State: &st.auto,
			FullWidth: true, Loading: st.loading,
		})
	})
)
