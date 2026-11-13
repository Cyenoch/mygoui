package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// The index follows the reader down the page and brings a picked section to
// the top; a window too narrow for it shows the page without one.
func TestPageIndexFollowsAndMovesTheReader(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1440, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	tt.Frame()
	checked := func(label string) bool {
		t.Helper()
		if info, ok := tt.Access(label); ok {
			return info.Checked
		}
		t.Fatalf("the index has no %q", label)
		return false
	}
	if checked("Section: API Reference") {
		t.Fatal("the index marks a section the reader has not reached")
	}
	if err := tt.Click("Section: API Reference"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	content, _ := tt.Find("Component content")
	heading, ok := tt.Find("API Reference")
	if !ok || s.Content.Y <= 0 || heading.Y < content.Y || heading.Y > content.Y+docsIndexReach {
		t.Fatalf("picking a section left its heading at %+v in content %+v", heading, content)
	}
	if !checked("Section: API Reference") {
		t.Fatal("the index does not mark the section the reader picked")
	}
	s.Content.Y = 0
	tt.Frame()
	if checked("Section: API Reference") {
		t.Fatal("the index kept a section marked after the reader scrolled away from it")
	}
	s.Content.Y = s.Content.MaxY
	tt.Frame()
	if !checked("Section: Theme tokens") {
		t.Fatal("the index does not mark the last section at the end of the page")
	}
	tt.SetSize(700, 900)
	tt.Frame()
	if _, ok := tt.Find("Section: Usage"); ok {
		t.Fatal("a narrow window keeps the index beside the content")
	}
}

// An example's Code view is the source of the demo it replaces.
func TestExampleCodeIsItsDemoSource(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	if err := tt.Click("Code"); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, text := range tt.Texts() {
		found = found || strings.Contains(text, "ButtonDestructiveOutline")
	}
	if !found {
		t.Fatal("the lead example's Code view does not show its demo's source")
	}
}

// The reference lists every field of the option types a page's snippet
// names, with the page's own defaults and descriptions where it gives them.
func TestAPIReferenceMergesSourceAndPageRows(t *testing.T) {
	page, _ := componentPage("button")
	sections := page.apiSections()
	if len(sections) != 1 || sections[0].Name != "mygoui.ButtonOptions" || sections[0].Doc == "" {
		t.Fatalf("button reference = %+v", sections)
	}
	rows := map[string]ComponentParam{}
	for _, row := range sections[0].Rows {
		rows[row.Name] = row
	}
	if got := rows["Variant"]; got.Type != "mygoui.ButtonVariant" || got.Default == "" || got.Description == "" {
		t.Fatalf("Variant row = %+v", got)
	}
	if len(rows) != len(apiTypes["mygoui.ButtonOptions"].Fields) {
		t.Fatalf("reference lists %d of %d fields", len(rows), len(apiTypes["mygoui.ButtonOptions"].Fields))
	}
	described := map[string]bool{}
	for _, param := range page.Params {
		described[param.Name] = true
	}
	for name := range rows {
		if !described[name] {
			return
		}
	}
	t.Fatal("no field came from the source alone; the reference only repeats the page's rows")
}

// The landing page owns the window: no sidebar beside it, yet its title bar
// and the search shortcut still lead into the catalog.
func TestLandingPageNavigatesWithoutSidebar(t *testing.T) {
	s := New()
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1440, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	tt.Frame()
	if navHas(tt, "Typography") {
		t.Fatal("the landing page shows the sidebar")
	}
	tt.Key(ui.Cmd, ui.KeyK)
	if !tt.Focused("Search components") {
		t.Fatal("the search shortcut does not reach search from the landing page")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if err := tt.Click("Components"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if s.Selected != "button" || !navHas(tt, "Typography") {
		t.Fatalf("Components led to %q; sidebar shown: %v", s.Selected, navHas(tt, "Typography"))
	}
}
