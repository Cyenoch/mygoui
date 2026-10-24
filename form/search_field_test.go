package form

import (
	"testing"

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
