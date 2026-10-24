package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// componentPageIndex returns every registered page by ID.
func componentPageIndex() map[string]Page {
	pages := map[string]Page{}
	for _, page := range Pages() {
		pages[page.ID] = page
	}
	return pages
}

func TestFormsAndSelectionComponentPagesExist(t *testing.T) {
	pages := componentPageIndex()

	forms := []string{
		"form", "fieldset", "label", "description", "fielderror", "errormessage",
		"input", "textarea", "textfield", "inputgroup", "inputotp", "searchfield",
		"numberfield", "checkbox", "checkboxgroup", "radio", "radiogroup",
		"switch", "switchgroup", "controlfield", "slider",
	}
	for _, id := range forms {
		page, ok := pages[id]
		if !ok {
			t.Fatalf("missing forms page %q", id)
		}
		if page.Group != "Forms" {
			t.Fatalf("page %q group = %q, want Forms", id, page.Group)
		}
		if page.Title == "" || page.Description == "" {
			t.Fatalf("page %q has no title or description", id)
		}
	}

	for _, id := range []string{"select", "combobox", "autocomplete"} {
		page, ok := pages[id]
		if !ok {
			t.Fatalf("missing forms page %q", id)
		}
		if page.Group != "Forms" {
			t.Fatalf("page %q group = %q, want Forms", id, page.Group)
		}
		if page.Order < 200 {
			t.Fatalf("page %q order = %d, want 200 or later", id, page.Order)
		}
	}

	listbox, ok := pages["listbox"]
	if !ok {
		t.Fatal("missing listbox page")
	}
	if listbox.Group != "Collections" {
		t.Fatalf("listbox group = %q, want Collections", listbox.Group)
	}
}

func TestFormsComponentPagesRenderTheirSections(t *testing.T) {
	pages := componentPageIndex()
	cases := map[string]struct{ body, token string }{
		"form":     {"Submit form", "Danger"},
		"input":    {"Username", "FieldBackground"},
		"checkbox": {"Parent folder", "AccentForeground"},
		"select":   {"Current fruit:", "AccentSoft"},
		"listbox":  {"Fruit: banana", "Foreground"},
	}
	for id, want := range cases {
		page, ok := pages[id]
		if !ok {
			t.Fatalf("missing page %q", id)
		}
		state := New()
		state.Selected = id
		tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
		tt.Frame()
		for _, text := range []string{
			page.Title, "Example", "Usage", "API", "Theme tokens", want.body, want.token,
		} {
			if !tt.HasText(text) {
				t.Fatalf("%s page does not show %q", id, text)
			}
		}
	}
}
