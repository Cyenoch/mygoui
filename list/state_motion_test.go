package list_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// A state change of a list passes through intermediate frames: these tests
// compare a region of the screen before, partway and after it.

const quarter = 60 * time.Millisecond

func snapshot(tt *ui.Tester, r ui.Rect) []uint8 {
	img := tt.Image()
	var out []uint8
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			c := img.RGBAAt(x, y)
			out = append(out, c.R, c.G, c.B)
		}
	}
	return out
}

func distance(a, b []uint8) int {
	d := 0
	for i := range a {
		if a[i] > b[i] {
			d += int(a[i] - b[i])
		} else {
			d += int(b[i] - a[i])
		}
	}
	return d
}

// wantPassesThrough fails unless mid differs from both ends, which differ
// from each other.
func wantPassesThrough(t *testing.T, what string, start, mid, end []uint8) {
	t.Helper()
	if distance(start, end) == 0 {
		t.Fatalf("%s: the change is not visible", what)
	}
	if distance(start, mid) == 0 || distance(mid, end) == 0 {
		t.Fatalf("%s: the middle frame equals an end (start-mid %d, mid-end %d)",
			what, distance(start, mid), distance(mid, end))
	}
}

func rowRegion(tt *ui.Tester, label string, width float32) ui.Rect {
	r, _ := tt.Find(label)
	return ui.Rect{X: 0, Y: r.Y, W: width, H: r.H}
}

func TestTreeArrowTurnsWithTheBranch(t *testing.T) {
	var state list.TreeState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.Tree(c, list.TreeOptions{Label: "Files", Nodes: treeNodes(), State: &state})
	}, 320, 320)
	arrow := rowRegion(tt, "src", 30)
	start := snapshot(tt, arrow)
	state.SetExpanded("src", true)
	tt.Frame()
	tt.Advance(quarter)
	mid := snapshot(tt, arrow)
	tt.Advance(400 * time.Millisecond)
	wantPassesThrough(t, "tree arrow", start, mid, snapshot(tt, arrow))
}

func TestTreeCheckDrawsIn(t *testing.T) {
	var state list.TreeState
	build := func(state *list.TreeState) *ui.Tester {
		return ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			list.Tree(c, list.TreeOptions{Label: "Files", Nodes: treeNodes(), Checkable: true, State: state})
		}, 320, 320)
	}
	tt := build(&state)
	box := rowRegion(tt, "go.mod", 44)
	start := snapshot(tt, box)
	state.SetChecked("go.mod", true)
	tt.Frame()
	tt.Advance(quarter)
	mid := snapshot(tt, box)
	tt.Advance(400 * time.Millisecond)
	end := snapshot(tt, box)
	wantPassesThrough(t, "tree check", start, mid, end)

	// A node built checked shows its tick on the first frame.
	var preset list.TreeState
	preset.SetChecked("go.mod", true)
	first := build(&preset)
	if distance(snapshot(first, box), end) != 0 {
		t.Fatal("a node built checked animates in")
	}
}

func TestListBoxCheckFadesIn(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Fruit", Items: items(), Selection: list.SelectionSingle, State: &state,
		})
	}, 320, 300)
	r, _ := tt.Find("Apple")
	end := ui.Rect{X: 0, Y: r.Y, W: 320, H: r.H}
	start := snapshot(tt, end)
	state.SetSelected(list.SelectionSingle, "apple", true)
	tt.Frame()
	tt.Advance(quarter / 2)
	mid := snapshot(tt, end)
	tt.Advance(400 * time.Millisecond)
	wantPassesThrough(t, "list box check", start, mid, snapshot(tt, end))
}
