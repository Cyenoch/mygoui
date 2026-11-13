package content

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Inside a JSON viewer, Cmd+F reaches the search field; there Enter and
// Shift+Enter step through the matches, Cmd+C leaves the chosen row's value
// alone, and Escape clears the query and keeps the field focused.
func TestJSONViewerSearchFieldKeysStepAndClear(t *testing.T) {
	var st JSONViewerState
	var err error
	o := JSONViewerOptions{Source: jsonSample, State: &st, ExpandDepth: -1, Toolbar: true}
	tt := jsonTester(&o, &err)
	if err := tt.Click("alpha"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyF)
	if !tt.Focused("Search") {
		t.Fatal("Cmd+F did not focus the search field")
	}
	tt.Type("a")
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if !tt.HasText("2 / 4") {
		t.Fatalf("Enter did not step to the next match: %q", tt.Texts())
	}
	tt.Key(ui.Shift, ui.KeyEnter)
	tt.Frame()
	if !tt.HasText("1 / 4") {
		t.Fatalf("Shift+Enter did not step back: %q", tt.Texts())
	}
	tt.SetClipboard("kept")
	tt.Key(ui.Cmd, ui.KeyC)
	if got := tt.Clipboard(); got != "kept" {
		t.Fatalf("Cmd+C in the search field copied the chosen row: %q", got)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if st.Query != "" || !tt.Focused("Search") {
		t.Fatalf("Escape left query %q, search focused %v", st.Query, tt.Focused("Search"))
	}
}

// With the focus inside a scrolling diff, F7 steps forward through the
// changes, wrapping after the last.
func TestDiffViewerF7StepsForward(t *testing.T) {
	var old, next strings.Builder
	for i := range 60 {
		fmt.Fprintf(&old, "line %d\n", i)
		if i == 5 || i == 50 {
			fmt.Fprintf(&next, "changed %d\n", i)
			continue
		}
		fmt.Fprintf(&next, "line %d\n", i)
	}
	o := DiffViewerOptions{Old: old.String(), New: next.String(), Title: "f.txt", Height: 200}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() { DiffViewer(c, o) })
	}, 640, 360)
	tt.Key(0, ui.KeyTab)
	for _, want := range []string{"1 / 2", "2 / 2", "1 / 2"} {
		tt.Key(0, ui.KeyF7)
		tt.Frame()
		if !tt.HasText(want) {
			t.Fatalf("want %q after F7: %q", want, tt.Texts())
		}
	}
}
