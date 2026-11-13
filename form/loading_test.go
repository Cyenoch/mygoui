package form_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// A loading field is where a query is typed while results load: it keeps
// focus and takes text across the flag flipping, and reports busy.
func TestLoadingInputStaysEditableAndReportsBusy(t *testing.T) {
	var value string
	loading := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Input(c, &value, form.InputOptions{Label: "Query", Loading: loading})
		})
	}, 360, 120)
	if err := tt.Click("Query"); err != nil {
		t.Fatal(err)
	}
	tt.Type("a")
	if tt.FramePending() {
		t.Fatal("idle input requested a frame")
	}
	loading = true
	tt.Frame()
	tt.Advance(200 * time.Millisecond)
	tt.Type("b")
	if value != "ab" {
		t.Fatalf("typing while loading committed %q", value)
	}
	if a, ok := tt.AccessByRole(ui.RoleTextField, "Query"); !ok || !a.Busy {
		t.Fatalf("loading input busy=%v found=%v", a.Busy, ok)
	}
	if !tt.FramePending() {
		t.Fatal("loading input does not turn its spinner")
	}
	loading = false
	tt.Frame()
	tt.Advance(400 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("spinner kept turning after loading ended")
	}
}

// A bare Input is one element: the caller's layout modifiers size the field
// that is painted and clicked, loading or not.
func TestLoadingInputKeepsCallerLayout(t *testing.T) {
	var value string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Row(c).Width(600).Children(func() {
			form.Input(c, &value, form.InputOptions{AriaLabel: "Query", Loading: true}).Grow(1)
		})
	}, 640, 80)
	tt.Advance(200 * time.Millisecond)
	a, ok := tt.AccessByRole(ui.RoleTextField, "Query")
	if !ok || a.Bounds.W < 590 {
		t.Fatalf("grown input is %v wide, found=%v", a.Bounds.W, ok)
	}
}

func TestLoadingSearchFieldKeepsSearchEditable(t *testing.T) {
	var value string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.SearchField(c, &value, form.SearchFieldOptions{Label: "Find", Loading: true})
		})
	}, 360, 120)
	tt.Advance(200 * time.Millisecond)
	if err := tt.Click("Find"); err != nil {
		t.Fatal(err)
	}
	tt.Type("go")
	if value != "go" {
		t.Fatalf("search committed %q", value)
	}
	if _, ok := tt.Access(locale.M().Loading); !ok {
		t.Fatal("no spinner in place of the search icon")
	}
}

// An empty loading popup says so instead of "No results"; with rows it keeps
// them.
func TestLoadingPopupRowsFollowOptions(t *testing.T) {
	fields := map[string]func() func(c *ui.Context, items []list.ListBoxItem){
		"select": func() func(*ui.Context, []list.ListBoxItem) {
			state := &form.SelectState{}
			return func(c *ui.Context, items []list.ListBoxItem) {
				form.Select(c, form.SelectOptions{Label: "Fruit", Loading: true, Items: items, State: state})
			}
		},
		"combo box": func() func(*ui.Context, []list.ListBoxItem) {
			state := &form.ComboBoxState{}
			return func(c *ui.Context, items []list.ListBoxItem) {
				form.ComboBox(c, form.ComboBoxOptions{Label: "Fruit", Loading: true, Items: items, State: state})
			}
		},
		"autocomplete": func() func(*ui.Context, []list.ListBoxItem) {
			state := &form.AutocompleteState{}
			return func(c *ui.Context, items []list.ListBoxItem) {
				form.Autocomplete(c, form.AutocompleteOptions{Label: "Fruit", Loading: true, Items: items, State: state})
			}
		},
	}
	// The field's own spinner is named Loading too, so the popup's row is
	// the extra mention an empty popup adds.
	mentions := func(t *testing.T, newField func() func(*ui.Context, []list.ListBoxItem), items []list.ListBoxItem) int {
		field := newField()
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			ui.Column(c).Fill().Padding(16).Children(func() { field(c, items) })
		}, 360, 400)
		if err := tt.Click("Fruit"); err != nil {
			t.Fatal(err)
		}
		tt.Key(0, ui.KeyDown) // a select opens on the arrow
		tt.Advance(200 * time.Millisecond)
		if tt.HasText(locale.M().NoResults) {
			t.Fatal("loading popup claims no results")
		}
		if items != nil && !tt.HasText("Alpha") {
			t.Fatal("rows vanished while loading")
		}
		n := 0
		for _, text := range tt.Texts() {
			if text == locale.M().Loading {
				n++
			}
		}
		return n
	}
	for name, newField := range fields {
		t.Run(name, func(t *testing.T) {
			empty, withRows := mentions(t, newField, nil), mentions(t, newField, comboItems())
			if empty <= withRows {
				t.Fatalf("empty popup mentions Loading %d times, popup with rows %d", empty, withRows)
			}
		})
	}
}
