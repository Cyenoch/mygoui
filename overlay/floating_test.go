package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestFloatingPanelLeavesContentUsableAndClosesFromInside(t *testing.T) {
	open, pressed, dismissed := true, 0, 0
	pos := ui.Point{X: 200, Y: 40}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{Label: "Beneath", OnClick: func() { pressed++ }})
		FloatingPanel(c, &open, FloatingPanelOptions{
			Title: "Inspector", Position: &pos, Width: 200,
			Body:      func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Inside"}) },
			OnDismiss: func() { dismissed++ },
		})
	}, 480, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()

	if err := tt.Click("Beneath"); err != nil {
		t.Fatal(err)
	}
	// The panel takes no focus, so Escape outside it is not its key.
	tt.Key(0, ui.KeyEscape)
	if pressed != 1 || !open || !tt.Focused("Beneath") {
		t.Fatalf("content beneath: pressed %d, panel open %v, focused %v", pressed, open, tt.Focused("Beneath"))
	}
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	if open || dismissed != 1 {
		t.Fatalf("Escape inside the panel: open %v, dismissed %d", open, dismissed)
	}
}

func TestFloatingPanelDragsByTitleWithinTheWindow(t *testing.T) {
	open := true
	pos := ui.Point{X: 100, Y: 40}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		FloatingPanel(c, &open, FloatingPanelOptions{
			Title: "Inspector", Position: &pos, Width: 200,
			Body: func() { ui.Text(c, "Body") },
		})
	}, 480, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	// A point on the title bar, beside the title.
	x, y := pos.X+100, pos.Y+20
	tt.Press(x, y)
	tt.Move(x+30, y+20)
	tt.Frame()
	if pos != (ui.Point{X: 130, Y: 60}) {
		t.Fatalf("dragged position %+v, want {130 60}", pos)
	}
	tt.Move(x+900, y-900)
	tt.Release(x+900, y-900)
	tt.Frame()
	if pos != (ui.Point{X: 280, Y: 0}) {
		t.Fatalf("position past the window %+v, want {280 0}", pos)
	}
}

func TestFloatingSheetResizesWithinLimitsAndStaysNonModal(t *testing.T) {
	open, pressed := true, 0
	height := float32(200)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Beneath", OnClick: func() { pressed++ }})
			FloatingSheet(c, &open, FloatingSheetOptions{
				Title: "Details", Height: &height, MinHeight: 120, MaxHeight: 240,
				Body: func() { ui.Text(c, "Body") },
			})
		})
	}, 480, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()

	if err := tt.Click("Beneath"); err != nil {
		t.Fatal(err)
	}
	if pressed != 1 || !open {
		t.Fatalf("content beneath: pressed %d, sheet open %v", pressed, open)
	}
	grip, ok := tt.Find("Resize")
	if !ok {
		t.Fatal("handle missing")
	}
	if sheet, _ := tt.AccessByRole(ui.RoleDialog, "Details"); sheet.Bounds != (ui.Rect{X: 0, Y: 200, W: 480, H: 200}) {
		t.Fatalf("sheet at %+v, want the bottom 200 DIPs", sheet.Bounds)
	}
	x, y := grip.X+grip.W/2, grip.Y+grip.H/2
	tt.Press(x, y)
	tt.Move(x, y-30)
	tt.Frame()
	if height != 230 {
		t.Fatalf("dragged height %v, want 230", height)
	}
	tt.Move(x, y-300)
	tt.Release(x, y-300)
	tt.Frame()
	if height != 240 {
		t.Fatalf("height past the maximum %v, want 240", height)
	}
	for range 8 {
		tt.Key(0, ui.KeyDown)
	}
	if height != 120 {
		t.Fatalf("height below the minimum %v, want 120", height)
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if open || tt.HasText("Details") {
		t.Fatalf("after close: open %v, shown %v", open, tt.HasText("Details"))
	}
}
