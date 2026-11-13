package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/egoist/mygo/ui"
)

func TestSearchFieldClearsWithButtonAndEscape(t *testing.T) {
	var value string
	submitted, cleared := 0, 0
	tester := formTextTester(t, func(c *ui.Context) {
		SearchField(c, &value, SearchFieldOptions{
			Label:       "Find",
			Placeholder: "Search records",
			OnSubmit:    func() { submitted++ },
			OnClear:     func() { cleared++ },
		})
	})
	formTextClick(t, tester, "Find")
	tester.Type("term")
	if value != "term" {
		t.Fatalf("search kept %q", value)
	}
	formTextClick(t, tester, "Clear")
	if value != "" || cleared != 1 {
		t.Fatalf("clear button: value %q, cleared %d", value, cleared)
	}

	formTextClick(t, tester, "Find")
	tester.Type("again")
	tester.Key(0, ui.KeyEscape)
	if value != "" || cleared != 2 {
		t.Fatalf("Escape: value %q, cleared %d", value, cleared)
	}

	formTextClick(t, tester, "Find")
	tester.Type("go")
	tester.Key(0, ui.KeyEnter)
	if submitted != 1 {
		t.Fatalf("submitted %d times", submitted)
	}
}

func TestSearchFieldEscapeClearsBeforeAncestor(t *testing.T) {
	value := ""
	escapes := 0
	tt := formTextTester(t, func(c *ui.Context) {
		ui.Column(c).Children(func() {
			SearchField(c, &value, SearchFieldOptions{Label: "Find"})
		}).OnShortcut(0, ui.KeyEscape, func() { escapes++ })
	})
	formTextClick(t, tt, "Find")
	tt.Type("term")
	tt.Key(0, ui.KeyEscape)
	if value != "" || escapes != 0 {
		t.Fatalf("Escape: value %q, ancestor ran %d times", value, escapes)
	}
	tt.Key(0, ui.KeyEscape)
	if escapes != 1 {
		t.Fatalf("Escape on the empty field reached the ancestor %d times, want 1", escapes)
	}
}

func TestSearchFieldHiddenClearIsNoTabStop(t *testing.T) {
	value, after := "", ""
	tt := formTextTester(t, func(c *ui.Context) {
		SearchField(c, &value, SearchFieldOptions{Label: "Find"})
		Input(c, &after, InputOptions{Label: "After"})
	})
	formTextClick(t, tt, "Find")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab stopped on the hidden clear button")
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, locale.M().Clear); ok {
		t.Fatal("the hidden clear button is exposed")
	}
}
