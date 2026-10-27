package navigation

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCommandFiltersAndActivatesFromKeyboardAndPointer(t *testing.T) {
	var state CommandState
	var selected []string
	items := []CommandItem{
		{Key: "archive", Label: "Archive", Keywords: "store"},
		{Key: "delete", Label: "Delete", Disabled: true},
		{Key: "duplicate", Label: "Duplicate"},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			Command(c, CommandOptions{State: &state, Label: "Commands", Items: items,
				OnSelect: func(key string) { selected = append(selected, key) }})
		})
	}, 400, 320)
	if err := tt.Click("Commands"); err != nil {
		t.Fatal(err)
	}
	tt.Type("sto")
	if !tt.HasText("Archive") || tt.HasText("Duplicate") {
		t.Fatalf("filtered items: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if len(selected) != 1 || selected[0] != "archive" {
		t.Fatalf("keyboard selection: %v", selected)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("dupl")
	if err := tt.Click("Duplicate"); err != nil {
		t.Fatal(err)
	}
	if len(selected) != 2 || selected[1] != "duplicate" {
		t.Fatalf("pointer selection: %v", selected)
	}
}

func TestCommandArrowsSkipDisabledActions(t *testing.T) {
	var state CommandState
	var selected string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, CommandOptions{
			State: &state, Label: "Commands",
			Items: []CommandItem{
				{Key: "one", Label: "One"},
				{Key: "locked", Label: "Locked", Disabled: true},
				{Key: "three", Label: "Three"},
			},
			OnSelect: func(key string) { selected = key },
		})
	}, 400, 240)
	if err := tt.Click("Commands"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	if state.List.HighlightKey() != "three" || !tt.Focused("Commands") {
		t.Fatalf("highlight %q, editor focused %v", state.List.HighlightKey(), tt.Focused("Commands"))
	}
	tt.Key(0, ui.KeyEnter)
	if selected != "three" {
		t.Fatalf("activated %q", selected)
	}
}

func TestCommandFiltersGroupedActions(t *testing.T) {
	var state CommandState
	var selected string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, CommandOptions{
			State: &state, Label: "Commands",
			Groups: []CommandGroup{
				{Label: "File", Items: []CommandItem{{Key: "save", Label: "Save"}}},
				{Label: "View", Items: []CommandItem{{Key: "zoom", Label: "Zoom", Keywords: "magnify"}}},
			},
			OnSelect: func(key string) { selected = key },
		})
	}, 400, 240)
	if err := tt.Click("Commands"); err != nil {
		t.Fatal(err)
	}
	tt.Type("magnify")
	if !tt.HasText("View") || !tt.HasText("Zoom") || tt.HasText("File") || tt.HasText("Save") {
		t.Fatalf("grouped filter: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if selected != "zoom" {
		t.Fatalf("grouped activation selected %q", selected)
	}
}

func TestCommandPointerAndKeyboardShareHighlight(t *testing.T) {
	var state CommandState
	var colors theme.Palette
	var selected []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			Command(c, CommandOptions{
				State: &state, Label: "Commands",
				Items:    []CommandItem{{Key: "archive", Label: "Archive"}, {Key: "duplicate", Label: "Duplicate"}},
				OnSelect: func(key string) { selected = append(selected, key) },
			})
		})
	}, 400, 240)
	if err := tt.Click("Commands"); err != nil {
		t.Fatal(err)
	}
	archive, _ := tt.Find("Archive")
	duplicate, _ := tt.Find("Duplicate")
	tt.Move(duplicate.X+duplicate.W/2, duplicate.Y+duplicate.H/2)
	tt.Frame()
	if key := state.List.HighlightKey(); key != "duplicate" {
		t.Fatalf("pointer highlighted %q", key)
	}
	img := tt.Image()
	wantSegmentPixel(t, img.RGBAAt(int(duplicate.X)+4, int(duplicate.Y+duplicate.H/2)),
		colors.Accent.Over(colors.Popover), "hovered command")
	wantSegmentPixel(t, img.RGBAAt(int(archive.X)+4, int(archive.Y+archive.H/2)),
		colors.Popover, "previous command")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyUp)
	if key := state.List.HighlightKey(); key != "archive" {
		t.Fatalf("keyboard highlighted %q", key)
	}
	img = tt.Image()
	wantSegmentPixel(t, img.RGBAAt(int(archive.X)+4, int(archive.Y+archive.H/2)),
		colors.Accent.Over(colors.Popover), "keyboard command")
	wantSegmentPixel(t, img.RGBAAt(int(duplicate.X)+4, int(duplicate.Y+duplicate.H/2)),
		colors.Popover, "stationary pointer command")
	tt.Key(0, ui.KeyEnter)
	if len(selected) != 2 || selected[0] != "duplicate" || selected[1] != "archive" {
		t.Fatalf("activated commands: %v", selected)
	}
}

func TestCommandKeyUsesEditorTextBeforeBoundValueCatchesUp(t *testing.T) {
	var state CommandState
	options := CommandOptions{
		State: &state, Label: "Commands",
		Items: []CommandItem{
			{Key: "archive", Label: "Archive"},
			{Key: "duplicate", Label: "Duplicate"},
			{Key: "duplex", Label: "Duplex"},
		},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, options)
	}, 400, 200)
	if err := tt.Click("Commands"); err != nil {
		t.Fatal(err)
	}
	if state.Query != "" || state.List.HighlightKey() != "archive" {
		t.Fatalf("initial query %q, highlight %q", state.Query, state.List.HighlightKey())
	}
	if got := commandKeySelection(options, &state, "dup", ui.KeyDown); got != "duplex" {
		t.Fatalf("Down after pending text highlighted %q", got)
	}
	if got := commandKeySelection(options, &state, "dup", ui.KeyEnter); got != "duplex" {
		t.Fatalf("Enter after pending text would activate %q", got)
	}
	if got := commandKeySelection(options, &state, "missing", ui.KeyEnter); got != "" {
		t.Fatalf("Enter with no match would activate %q", got)
	}
}

func TestCommandPendingHighlightSurvivesRepeatedActivation(t *testing.T) {
	state := CommandState{Query: "arc"}
	options := CommandOptions{
		State: &state, Label: "Commands",
		Items: []CommandItem{
			{Key: "archive", Label: "Archive"},
			{Key: "duplicate", Label: "Duplicate"},
			{Key: "duplex", Label: "Duplex"},
		},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, options)
	}, 400, 200)
	if err := tt.Click("Commands"); err != nil {
		t.Fatal(err)
	}
	if got := commandKeySelection(options, &state, "dup", ui.KeyDown); got != "duplex" {
		t.Fatalf("Down after pending text highlighted %q", got)
	}
	state.Query = "dup"
	for range 2 {
		if got := commandKeySelection(options, &state, "dup", ui.KeyEnter); got != "duplex" {
			t.Fatalf("repeated Enter would activate %q", got)
		}
	}
}

func TestCommandSelectionCanClearQuery(t *testing.T) {
	var state CommandState
	var queryAtSelection string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, CommandOptions{
			State: &state, Label: "Commands",
			Items: []CommandItem{{Key: "duplicate", Label: "Duplicate"}},
			OnSelect: func(string) {
				queryAtSelection = state.Query
				state.Query = ""
			},
		})
	}, 400, 200)
	if err := tt.Click("Commands"); err != nil {
		t.Fatal(err)
	}
	tt.Type("dupl")
	tt.Key(0, ui.KeyEnter)
	if queryAtSelection != "dupl" || state.Query != "" {
		t.Fatalf("query at selection %q, after selection %q", queryAtSelection, state.Query)
	}
	if !tt.HasText("Duplicate") || !tt.Focused("Commands") {
		t.Fatalf("clearing query lost the command list or input focus: %v", tt.Texts())
	}
}
