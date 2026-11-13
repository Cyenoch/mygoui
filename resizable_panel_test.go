package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func resizablePanelTester(o *mygoui.ResizablePanelOptions) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Row(c).Fill().Children(func() {
			if o.Side == mygoui.PanelStart {
				ui.Box(c).Grow(1)
			}
			mygoui.ResizablePanel(c, *o, func() { ui.Text(c, "Files") })
		})
	}, 600, 300)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	tt.Frame()
	return tt
}

func dragHandle(t *testing.T, tt *ui.Tester, dx float32) {
	t.Helper()
	r := findRect(t, tt, "Resize")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x+dx/2, y)
	tt.Move(x+dx, y)
	tt.Release(x+dx, y)
	// Keep the next press from reading as a double click.
	tt.Advance(time.Second)
}

func TestResizablePanelDragClampsToBounds(t *testing.T) {
	for _, tc := range []struct {
		name string
		side mygoui.PanelSide
		rtl  bool
		grow float32
	}{
		{name: "end", side: mygoui.PanelEnd, grow: 1},
		{name: "start", side: mygoui.PanelStart, grow: -1},
		{name: "end rtl", side: mygoui.PanelEnd, rtl: true, grow: -1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			size := float32(200)
			var reported []float32
			o := mygoui.ResizablePanelOptions{
				Side: tc.side, Size: &size, Min: 120, Max: 320,
				OnResize: func(v float32) { reported = append(reported, v) },
			}
			tt := resizablePanelTester(&o)
			if tc.rtl {
				tt.SetDirection(ui.RightToLeft)
				tt.Frame()
			}
			dragHandle(t, tt, 40*tc.grow)
			if size != 240 {
				t.Fatalf("size after a 40 DIP drag = %v, want 240", size)
			}
			dragHandle(t, tt, 300*tc.grow)
			if size != 320 {
				t.Fatalf("size past the maximum = %v, want 320", size)
			}
			dragHandle(t, tt, -500*tc.grow)
			if size != 120 {
				t.Fatalf("size past the minimum = %v, want 120", size)
			}
			if len(reported) == 0 || reported[len(reported)-1] != 120 {
				t.Fatalf("OnResize reported %v", reported)
			}
			if w := findRect(t, tt, "Resize"); w.H != 300 {
				t.Fatalf("handle = %+v, want the panel's height", w)
			}
		})
	}
}

func TestResizablePanelKeyboard(t *testing.T) {
	size, collapsed := float32(200), false
	o := mygoui.ResizablePanelOptions{Size: &size, Min: 120, Max: 320, Collapsed: &collapsed}
	tt := resizablePanelTester(&o)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Resize") {
		t.Fatal("Tab did not reach the handle")
	}
	tt.Key(0, ui.KeyRight)
	if size != 216 {
		t.Fatalf("Right: size = %v, want 216", size)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	if size != 184 {
		t.Fatalf("Left twice: size = %v, want 184", size)
	}
	tt.Key(0, ui.KeyEnd)
	if size != 320 {
		t.Fatalf("End: size = %v, want 320", size)
	}
	tt.Key(0, ui.KeyRight)
	if size != 320 {
		t.Fatalf("Right at the maximum: size = %v", size)
	}
	tt.Key(0, ui.KeyHome)
	if size != 120 {
		t.Fatalf("Home: size = %v, want 120", size)
	}
	tt.Key(0, ui.KeyEnter)
	if !collapsed || tt.HasText("Files") {
		t.Fatalf("Enter: collapsed=%v texts=%v", collapsed, tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	if collapsed || size != 120 || !tt.HasText("Files") {
		t.Fatalf("Right on a collapsed panel: collapsed=%v size=%v", collapsed, size)
	}
}

func TestResizablePanelDragCollapses(t *testing.T) {
	size, collapsed := float32(200), false
	var changes []bool
	o := mygoui.ResizablePanelOptions{
		Size: &size, Min: 120, Collapsed: &collapsed,
		OnCollapseChange: func(v bool) { changes = append(changes, v) },
	}
	tt := resizablePanelTester(&o)
	dragHandle(t, tt, -100)
	if collapsed || size != 120 {
		t.Fatalf("drag to 100: collapsed=%v size=%v, want the minimum", collapsed, size)
	}
	size = 200
	tt.Frame()
	dragHandle(t, tt, -150)
	if !collapsed || size != 200 || len(changes) != 1 {
		t.Fatalf("drag below half the minimum: collapsed=%v size=%v changes=%v", collapsed, size, changes)
	}
	if r := findRect(t, tt, "Resize"); r.X != 0 || r.W != 4 {
		t.Fatalf("collapsed handle = %+v, want it on the window edge", r)
	}
	dragHandle(t, tt, 150)
	if collapsed || size != 150 {
		t.Fatalf("drag back out: collapsed=%v size=%v", collapsed, size)
	}
	r := findRect(t, tt, "Resize")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.ClickAt(x, y)
	tt.Press(x, y)
	tt.Move(x+1, y)
	tt.Release(x+1, y)
	if !collapsed || size != 150 {
		t.Fatalf("double click: collapsed=%v size=%v", collapsed, size)
	}
}

// The handle is a Tab stop where it is on screen: before the content of a
// panel resized from its start edge, after it otherwise. Either way it takes
// the press on its own strip, over the content.
func TestResizablePanelHandleFollowsItsEdgeInTabOrder(t *testing.T) {
	for _, tc := range []struct {
		name          string
		side          mygoui.PanelSide
		first, second string
	}{
		{"end edge", mygoui.PanelEnd, "Open file", "Resize"},
		{"start edge", mygoui.PanelStart, "Resize", "Open file"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			size := float32(200)
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Row(c).Fill().Children(func() {
					if tc.side == mygoui.PanelStart {
						ui.Box(c).Grow(1)
					}
					mygoui.ResizablePanel(c, mygoui.ResizablePanelOptions{Side: tc.side, Size: &size, Min: 120, Max: 320}, func() {
						ui.Column(c).Fill().Background(theme.Colors(c).Background).Children(func() {
							ui.Button(c, "Open file")
						})
					})
				})
			}, 600, 300)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			tt.Frame()
			tt.Key(0, ui.KeyTab)
			if !tt.Focused(tc.first) {
				t.Fatalf("Tab did not reach %q first", tc.first)
			}
			tt.Key(0, ui.KeyTab)
			if !tt.Focused(tc.second) {
				t.Fatalf("Tab did not go on to %q", tc.second)
			}
			grow := float32(40)
			if tc.side == mygoui.PanelStart {
				grow = -40
			}
			dragHandle(t, tt, grow)
			if size != 240 {
				t.Fatalf("dragging the handle over the content: size = %v, want 240", size)
			}
		})
	}
}
