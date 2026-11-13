package content

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// twoChanges is a 60-line text and the same text with lines 5 and 50
// changed.
func twoChanges() (old, next string) {
	var o, n strings.Builder
	for i := range 60 {
		fmt.Fprintf(&o, "line %d\n", i)
		if i == 5 || i == 50 {
			fmt.Fprintf(&n, "changed %d\n", i)
			continue
		}
		fmt.Fprintf(&n, "line %d\n", i)
	}
	return o.String(), n.String()
}

// searchViewers builds each viewer with a toolbar inside a column that
// counts the Escapes reaching it, as a dialog would.
func searchViewers(query *string, escapes *int) map[string]func(c *ui.Context) {
	frame := func(c *ui.Context, viewer func()) {
		ui.Column(c).Fill().Padding(16).Children(viewer).
			OnShortcut(0, ui.KeyEscape, func() { *escapes++ })
	}
	var log LogViewerState
	var json JSONViewerState
	return map[string]func(c *ui.Context){
		"log": func(c *ui.Context) {
			frame(c, func() {
				LogViewer(c, LogViewerOptions{
					Line: func(i int) string { return fmt.Sprintf("entry %d", i) }, Count: 20,
					Height: 200, Toolbar: true, State: &log,
				})
				*query = log.Query
			})
		},
		"json": func(c *ui.Context) {
			frame(c, func() {
				_, _ = JSONViewer(c, JSONViewerOptions{Source: jsonSample, Toolbar: true, State: &json})
				*query = json.Query
			})
		},
	}
}

// Inside a dialog, Escape in a search field holding text clears the text
// and leaves the dialog open; the next Escape reaches the dialog.
func TestViewerSearchEscapeClearsBeforeTheDialog(t *testing.T) {
	var query string
	var escapes int
	for name, viewer := range searchViewers(&query, &escapes) {
		t.Run(name, func(t *testing.T) {
			escapes = 0
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); viewer(c) }, 480, 480)
			if err := tt.Click("Search"); err != nil {
				t.Fatal(err)
			}
			tt.Type("a")
			tt.Frame()
			tt.Key(0, ui.KeyEscape)
			tt.Frame()
			if query != "" || escapes != 0 {
				t.Fatalf("Escape left query %q and reached the dialog %d times", query, escapes)
			}
			tt.Key(0, ui.KeyEscape)
			if escapes != 1 {
				t.Fatalf("Escape on the cleared field reached the dialog %d times, want 1", escapes)
			}
		})
	}
}

// Clear empties the query and leaves the focus in the search field.
func TestViewerSearchClearKeepsTheFocus(t *testing.T) {
	var query string
	var escapes int
	for name, viewer := range searchViewers(&query, &escapes) {
		t.Run(name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); viewer(c) }, 480, 480)
			if err := tt.Click("Search"); err != nil {
				t.Fatal(err)
			}
			tt.Type("a")
			tt.Frame()
			if err := tt.Click("Clear"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if query != "" || !tt.Focused("Search") {
				t.Fatalf("Clear left query %q, search focused %v", query, tt.Focused("Search"))
			}
		})
	}
}

// A window's Find command focuses a viewer's search field through its
// state.
func TestViewerSearchFocusesForTheWindowsFind(t *testing.T) {
	var log LogViewerState
	var json JSONViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Gap(8).Children(func() {
			ui.Button(c, "Find log").OnClick(log.FocusSearch)
			ui.Button(c, "Find JSON").OnClick(json.FocusSearch)
			LogViewer(c, LogViewerOptions{Lines: []string{"a", "b"}, Height: 80, Toolbar: true, State: &log, Label: "Log"})
			_, _ = JSONViewer(c, JSONViewerOptions{Source: jsonSample, Toolbar: true, State: &json})
		})
	}, 480, 640)
	for _, button := range []string{"Find log", "Find JSON"} {
		if err := tt.Click(button); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if !tt.Focused("Search") {
			t.Fatalf("%s did not focus a search field", button)
		}
	}
}

// The match, filter and change counts are statuses: stepping and filtering
// are read out.
func TestViewerCountsAreAnnounced(t *testing.T) {
	old, next := twoChanges()
	var log LogViewerState
	var json JSONViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Gap(8).Children(func() {
			_, _ = JSONViewer(c, JSONViewerOptions{Source: jsonSample, Toolbar: true, State: &json})
			LogViewer(c, LogViewerOptions{
				Line: func(i int) string { return fmt.Sprintf("entry %d", i) }, Count: 30,
				Height: 80, Toolbar: true, State: &log,
			})
		})
	}, 640, 600)
	tt.UseScreenReader()
	json.Query = "a"
	log.Query = "entry 1"
	tt.Frame()
	json.step(1)
	log.Query = "entry 15"
	tt.Frame()
	got := tt.Announcements()
	for _, want := range []string{"2 / 4", "1 / 30"} {
		if !slices.Contains(got, want) {
			t.Fatalf("%q was not announced: %q", want, got)
		}
	}
	diff := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			DiffViewer(c, DiffViewerOptions{Old: old, New: next, Title: "f.txt", Height: 120})
		})
	}, 640, 300)
	diff.UseScreenReader()
	if err := diff.Click("Next"); err != nil {
		t.Fatal(err)
	}
	diff.Frame()
	if got := diff.Announcements(); !slices.Contains(got, "1 / 2") {
		t.Fatalf("the change position was not announced: %q", got)
	}
}

// A document that does not parse is reported at once.
func TestJSONViewerAnnouncesAnInvalidDocument(t *testing.T) {
	var err error
	o := JSONViewerOptions{Source: jsonSample}
	tt := jsonTester(&o, &err)
	tt.UseScreenReader()
	o.Source = `{"a": `
	tt.Frame()
	if err == nil || !slices.Contains(tt.Alerts(), err.Error()) {
		t.Fatalf("error %v not alerted: %q", err, tt.Alerts())
	}
}

// The lines take the focus: Down chooses a line, whose level is read with
// it, and the context-menu key opens its menu.
func TestLogLinesAreReachableByKeyboard(t *testing.T) {
	lines := []string{"boot", "ready", "disk failed"}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{Lines: lines, Height: 120, Level: func(i int) LogLevel {
				if i == 2 {
					return LogError
				}
				return LogPlain
			}})
		})
	}, 480, 240)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Log") {
		t.Fatal("Tab did not reach the lines")
	}
	for range 3 {
		tt.Key(0, ui.KeyDown)
	}
	tt.Frame()
	if line, ok := tt.AccessByRole(ui.RoleListItem, "disk failed, Error"); !ok || !line.Selected {
		t.Fatalf("the error line, read with its level, is not chosen: %+v, %v", line, ok)
	}
	tt.Key(ui.Shift, ui.KeyF10)
	if err := tt.ChooseMenuItem("Copy line"); err != nil {
		t.Fatalf("no line menu from the keyboard: %v (menu %q)", err, tt.Menu())
	}
	if got := tt.Clipboard(); got != "disk failed" {
		t.Fatalf("clipboard = %q", got)
	}
}

// Jump to latest goes once the end shows and leaves the focus on the lines.
func TestLogJumpToLatestLeavesTheFocusOnTheLines(t *testing.T) {
	var st LogViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{
				Line: func(i int) string { return fmt.Sprintf("entry %d", i) }, Count: 200,
				Height: 160, Follow: true, State: &st,
			})
		})
	}, 480, 240)
	st.List.ScrollTo(0, ui.Start)
	tt.Frame()
	if err := tt.Click("Jump to latest"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(theme.Motion().State * 2)
	tt.Frame()
	if !st.List.AtEnd() || !tt.Focused("Log") {
		t.Fatalf("at end %v, lines focused %v", st.List.AtEnd(), tt.Focused("Log"))
	}
}

// A code block or snippet without a name is not an unnamed group, and a
// named one names its copy button after itself.
func TestCodeBlocksNameTheirGroupAndCopyButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Gap(8).Children(func() {
			Highlighter(c, HighlighterOptions{Code: "x := 1"})
			Snippet(c, SnippetOptions{Text: "npm install"})
			Highlighter(c, HighlighterOptions{Code: "y := 2", Title: "main.go"})
			Snippet(c, SnippetOptions{Text: "go test", Prefix: "$", Label: "Test"})
		})
	}, 480, 480)
	for _, n := range tt.AccessAll() {
		if n.Role == ui.RoleGroup && n.Name == "" {
			t.Fatalf("unnamed group: %+v", n)
		}
	}
	for button, want := range map[string]string{"Copy main.go": "y := 2", "Copy Test": "go test"} {
		if err := tt.Click(button); err != nil {
			t.Fatal(err)
		}
		if got := tt.Clipboard(); got != want {
			t.Fatalf("%s copied %q", button, got)
		}
	}
}

// A diff without a Title still steps through its changes and switches
// layout.
func TestDiffViewerWithoutTitleKeepsItsControls(t *testing.T) {
	old, next := twoChanges()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			DiffViewer(c, DiffViewerOptions{Old: old, New: next, Height: 200})
		})
	}, 640, 360)
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyF7)
	tt.Frame()
	if !tt.HasText("2 / 2") {
		t.Fatalf("Next and F7 did not step: %q", tt.Texts())
	}
	if _, ok := tt.AccessByRole(ui.RoleToggleButton, "Side by side"); !ok {
		t.Fatal("no layout toggle")
	}
	if err := tt.Click("Copy"); err != nil || tt.Clipboard() != next {
		t.Fatalf("copy: %v", err)
	}
}
