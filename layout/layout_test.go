package layout_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestResizableKeepsIndependentProportionsAndNativeKeyboardBounds(t *testing.T) {
	first, second := layout.ResizableState{Percent: 40}, layout.ResizableState{Percent: 65}
	vertical := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			layout.Resizable(c, &first, layout.ResizableOptions{Vertical: vertical}, func() { ui.Text(c, "First pane") }, func() { ui.Text(c, "Second pane") }).Height(140).FillWidth()
			layout.Resizable(c, &second, layout.ResizableOptions{}, func() { ui.Text(c, "Independent pane") }, func() { ui.Text(c, "Independent detail") }).Height(100).FillWidth()
		})
	}, 400, 260)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Divider") {
		t.Fatal("native splitter is not keyboard focusable")
	}
	tt.Key(0, ui.KeyRight)
	if first.Percent <= 40 || second.Percent != 65 {
		t.Fatal("keyboard splitter changed wrong caller state")
	}
	for range 100 {
		tt.Key(0, ui.KeyRight)
	}
	if first.Percent >= 100 || first.Percent <= 50 {
		t.Fatalf("divider escaped native minimum pane bound: %v", first.Percent)
	}
	saved := first.Percent
	tt.SetSize(180, 260)
	tt.SetSize(400, 260)
	if first.Percent != saved {
		t.Fatal("temporary compression overwrote caller proportion")
	}
	vertical = true
	tt.Frame()
	before := first.Percent
	tt.Key(0, ui.KeyUp)
	if first.Percent >= before {
		t.Fatal("vertical splitter ignored native Up")
	}
}

func TestScrollAreaTracksEachOrientationAndDirectionScope(t *testing.T) {
	for _, test := range []struct {
		name        string
		orientation layout.ScrollOrientation
		x, y        bool
	}{
		{"vertical", layout.ScrollVertical, false, true},
		{"horizontal", layout.ScrollHorizontal, true, false},
		{"both", layout.ScrollBoth, true, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				layout.ScrollArea(c, &state, layout.ScrollAreaOptions{Orientation: test.orientation}, func() { ui.Box(c).Size(600, 600).Shrink(0) }).Size(160, 120)
			}, 180, 140)
			tt.Scroll(80, 60, 70, 70)
			if test.x && state.X <= 0 {
				t.Fatal("horizontal offset did not track wheel")
			}
			if test.y && state.Y <= 0 {
				t.Fatal("vertical offset did not track wheel")
			}
		})
	}
	var inside, outside ui.Direction
	ui.NewTester(func(c *ui.Context) {
		layout.Direction(c, ui.RightToLeft, func() { inside = c.Direction() })
		outside = c.Direction()
	}, 200, 100)
	if inside != ui.RightToLeft || outside != ui.LeftToRight {
		t.Fatal("direction leaked out of its native scope")
	}
}
