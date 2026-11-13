package overlay

import (
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func shortcutsTester(state *ShortcutsDialogState, groups []ShortcutGroup, w, h int) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { state.Open = true })
		ShortcutsDialog(c, state, ShortcutsDialogOptions{Title: "Keyboard shortcuts", Groups: groups})
	}, w, h)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return tt
}

func TestShortcutsDialogSearchNarrowsRowsAndGroups(t *testing.T) {
	var state ShortcutsDialogState
	tt := shortcutsTester(&state, []ShortcutGroup{
		{Title: "General", Items: []Shortcut{
			{Label: "Open command menu", Keys: []string{"⌘", "K"}},
			{Label: "Toggle sidebar", Keys: []string{"⌘", "B"}},
		}},
		{Title: "Editing", Items: []Shortcut{{Label: "Über-paste", Keys: []string{"⌘", "V"}}}},
	}, 800, 600)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.Focused("Search") {
		t.Fatal("the search field did not take the focus")
	}
	if chord, ok := tt.Access("⌘ K"); !ok || chord.Bounds.W <= 0 {
		t.Fatalf("the chord is not exposed: %+v (found %v)", chord, ok)
	}

	// Typed text reaches the field without a click and ignores case.
	tt.Type("SIDE")
	tt.Frame()
	if !tt.HasText("Toggle sidebar") || tt.HasText("Open command menu") || tt.HasText("Editing") {
		t.Fatalf("query %q shows %v", state.Query, tt.Texts())
	}
	state.Query = "über"
	tt.Frame()
	if !tt.HasText("Über-paste") || tt.HasText("General") {
		t.Fatalf("query %q shows %v", state.Query, tt.Texts())
	}
	state.Query = "zzz"
	tt.Frame()
	if !tt.HasText("No results") || tt.HasText("Editing") {
		t.Fatalf("query %q shows %v", state.Query, tt.Texts())
	}

	// A reopened sheet starts from the whole list again.
	tt.Key(0, ui.KeyEscape) // clears the field
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Open {
		t.Fatal("Escape did not close the dialog")
	}
	state.Query = "stale"
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state.Query != "" || !tt.HasText("Open command menu") {
		t.Fatalf("reopened with query %q", state.Query)
	}
}

func TestShortcutsDialogKeepsSearchInPlaceWhileRowsScroll(t *testing.T) {
	state := ShortcutsDialogState{Open: true}
	items := make([]Shortcut, 60)
	for i := range items {
		items[i] = Shortcut{Label: "Action " + string(rune('A'+i%26)) + string(rune('a'+i/26)), Keys: []string{"F1"}}
	}
	tt := shortcutsTester(&state, []ShortcutGroup{{Title: "All", Items: items}}, 800, 480)
	tt.Frame()
	before, ok := tt.Find("Search")
	first, _ := tt.Find(items[0].Label)
	if !ok {
		t.Fatal("missing search field")
	}
	tt.Wheel(first.X+10, first.Y+10, 0, 400)
	tt.Frame()
	after, _ := tt.Find("Search")
	moved, _ := tt.Find(items[0].Label)
	if after != before || after.Y < 0 || after.Y+after.H > 480 {
		t.Fatalf("search field moved from %+v to %+v", before, after)
	}
	if moved.Y >= first.Y {
		t.Fatalf("rows did not scroll: first row at %+v, then %+v", first, moved)
	}
}

// Users look a shortcut up by the keys they remember as often as by its name.
func TestShortcutsDialogSearchFindsKeysAndBuildsChordsFromHotkeys(t *testing.T) {
	hotkey := form.Hotkey{Mods: ui.Ctrl | ui.Shift, Key: ui.KeyP}
	var state ShortcutsDialogState
	tt := shortcutsTester(&state, []ShortcutGroup{{Title: "General", Items: []Shortcut{
		{Label: "Open command menu", Keys: []string{"⌘", "K"}},
		{Label: "Print", Hotkey: hotkey},
		{Label: "Quit", Keys: []string{"Esc"}},
	}}}, 400, 600)
	state.Open = true
	tt.Frame()
	if _, ok := tt.Access(strings.Join(hotkey.Keys(), " ")); !ok {
		t.Fatal("a Hotkey is not shown as its platform's key caps")
	}
	for query, want := range map[string]string{
		"⌘k":     "Open command menu",
		"cmd k":  "Open command menu",
		"escape": "",
		"esc":    "Quit",
		"shift":  "Print",
	} {
		state.Query = query
		tt.Frame()
		for _, label := range []string{"Open command menu", "Print", "Quit"} {
			if shown := tt.HasText(label); shown != (label == want) {
				t.Errorf("query %q: %q shown=%v", query, label, shown)
			}
		}
	}
}
