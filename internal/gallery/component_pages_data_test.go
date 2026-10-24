package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// componentDataPage finds one registered page by ID.
func componentDataPage(t *testing.T, id string) Page {
	t.Helper()
	for _, page := range Pages() {
		if page.ID == id {
			return page
		}
	}
	t.Fatalf("page %q is not registered", id)
	return Page{}
}

func TestComponentDataPagesRenderTitles(t *testing.T) {
	wanted := []struct{ id, title, group string }{
		{"calendar", "Calendar", "Dates"},
		{"rangecalendar", "RangeCalendar", "Dates"},
		{"colorpicker", "ColorPicker", "Color"},
		{"table", "Table", "Data"},
	}
	for _, want := range wanted {
		page := componentDataPage(t, want.id)
		if page.Title != want.title || page.Group != want.group {
			t.Fatalf("page %q = %q/%q, want %q/%q", want.id, page.Title, page.Group, want.title, want.group)
		}
		state := New()
		state.Selected = page.ID
		tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
		tt.Frame()
		for _, text := range []string{page.Title, "Example", "Usage", "API", "Theme tokens"} {
			if !tt.HasText(text) {
				t.Fatalf("%s page does not show %q", page.ID, text)
			}
		}
	}
}

func TestComponentDataPageShowsUsageAPIAndTokens(t *testing.T) {
	state := New()
	state.Selected = "calendar"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	for _, want := range []string{
		"Usage", "API", "Theme tokens",
		"date.Calendar(c, date.CalendarOptions{", "DisabledDates", "AccentSoftForeground",
	} {
		if !tt.HasText(want) {
			t.Fatalf("calendar page does not show %q", want)
		}
	}
}
