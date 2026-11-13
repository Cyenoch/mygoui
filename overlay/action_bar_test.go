package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestActionBarFloatsAtTheBottomWithoutTakingTheContainer(t *testing.T) {
	selected, pressed, archived := 0, 0, 0
	var bar ui.Rect
	var row ui.Handle
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Row", Handle: &row, OnClick: func() { pressed++; selected = 2 },
			})
			e := ActionBar(c, selected > 0, ActionBarOptions{
				Label: "selected", Count: selected,
				Actions:     []navigation.OverflowItem{{Label: "Archive", OnSelect: func() { archived++ }}},
				OnClose:     func() { selected = 0 },
				ReturnFocus: &row,
			})
			bar = ui.Rect{}
			if e.Valid() {
				bar = e.Bounds()
			}
		})
	}, 480, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if _, ok := tt.AccessByRole(ui.RoleToolbar, "2 selected"); ok || tt.HasText("selected") {
		t.Fatal("the bar shows without a selection")
	}

	// The closed bar's frame covers the container and must not take its input.
	if err := tt.Click("Row"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	tt.Frame()
	if pressed != 1 || !tt.Focused("Row") {
		t.Fatalf("content beneath: pressed=%d focused=%v", pressed, tt.Focused("Row"))
	}
	if _, ok := tt.AccessByRole(ui.RoleToolbar, "2 selected"); !ok {
		t.Fatal("the bar is not a toolbar named by its label")
	}
	if center := bar.X + bar.W/2; center < 238 || center > 242 || bar.Y+bar.H != 320-theme.Space16 {
		t.Fatalf("bar at %+v, want bottom center 16 DIPs above the edge", bar)
	}

	if err := tt.Click("Archive"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if archived != 1 || selected != 0 || tt.HasText("selected") {
		t.Fatalf("archived=%d selected=%d", archived, selected)
	}
	// The closing bar held the focus; it must not leave the window without one.
	if !tt.Focused("Row") {
		t.Fatal("closing the bar dropped the focus instead of returning it")
	}
}

// On a narrow window the actions used to run under the close button and the
// last one vanished; now they fold into a menu and the close button stays.
func TestActionBarFoldsActionsInANarrowWindowAndEscapeCloses(t *testing.T) {
	closed := 0
	var more, closeButton ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			ActionBar(c, true, ActionBarOptions{
				Label: "selected", Count: 12,
				Actions: []navigation.OverflowItem{{Label: "Archive"}, {Label: "Delete"}, {Label: "Move to folder"}},
				OnClose: func() { closed++ },
			})
		})
	}, 280, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	for range 3 {
		tt.Frame()
	}
	more, ok := tt.Find("More")
	if !ok {
		t.Fatal("the actions that do not fit have no menu")
	}
	closeButton, _ = tt.Find("Close")
	bar, _ := tt.AccessByRole(ui.RoleToolbar, "12 selected")
	if more.X+more.W > closeButton.X || closeButton.X+closeButton.W > bar.Bounds.X+bar.Bounds.W {
		t.Fatalf("more %+v, close %+v, bar %+v overlap or overflow", more, closeButton, bar.Bounds)
	}
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("Move to folder") {
		t.Fatal("the folded action is not in the menu")
	}
	tt.Key(0, ui.KeyEscape) // closes the menu
	tt.Frame()
	if closed != 0 {
		t.Fatalf("Escape in the menu closed the bar")
	}
	if !tt.Focused("More") {
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if closed != 1 {
		t.Fatalf("Escape with the focus in the bar ran close %d times", closed)
	}
}
