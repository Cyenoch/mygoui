package list_test

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// sortableFixture applies each reported move to its own order, as a caller
// does.
type sortableFixture struct {
	order    []string
	moves    [][2]int
	disabled bool
}

func (f *sortableFixture) tester() *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		items := make([]list.SortableItem, len(f.order))
		for i, id := range f.order {
			items[i] = list.SortableItem{ID: id, Label: id}
		}
		ui.Column(c).Padding(16).FillWidth().Children(func() {
			list.SortableList(c, list.SortableListOptions{
				Label: "Layers", Items: items, Disabled: f.disabled,
				OnMove: func(from, to int) {
					f.moves = append(f.moves, [2]int{from, to})
					id := f.order[from]
					f.order = slices.Insert(slices.Delete(f.order, from, from+1), to, id)
				},
			})
		})
	}, 320, 320)
}

func TestSortableListMovesFromTheKeyboard(t *testing.T) {
	f := &sortableFixture{order: []string{"Alpha", "Bravo", "Charlie"}}
	tt := f.tester()
	if err := tt.Click("Reorder Alpha"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Alt, ui.KeyDown)
	if !slices.Equal(f.order, []string{"Bravo", "Alpha", "Charlie"}) || !tt.Focused("Reorder Alpha") {
		t.Fatalf("after Alt+Down: %q, moves %v, handle focused %v", f.order, f.moves, tt.Focused("Reorder Alpha"))
	}
	// The handle moves with its row, so the same key keeps moving it.
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Key(ui.Alt, ui.KeyUp)
	want := [][2]int{{0, 1}, {1, 2}, {2, 1}}
	if !slices.Equal(f.moves, want) || !slices.Equal(f.order, []string{"Bravo", "Alpha", "Charlie"}) {
		t.Fatalf("moves %v, want %v; order %q", f.moves, want, f.order)
	}
}

func TestSortableListMovesByDraggingTheHandle(t *testing.T) {
	f := &sortableFixture{order: []string{"Alpha", "Bravo", "Charlie"}}
	tt := f.tester()
	from, ok := tt.Find("Reorder Alpha")
	to, ok2 := tt.Find("Reorder Charlie")
	if !ok || !ok2 {
		t.Fatal("handles missing")
	}
	drag := func(x, y0, y1 float32) {
		tt.Press(x, y0)
		for k := 1; k <= 4; k++ {
			tt.Move(x, y0+(y1-y0)*float32(k)/4)
		}
		tt.Release(x, y1)
		tt.Frame()
	}
	x := from.X + from.W/2
	drag(x, from.Y+from.H/2, to.Y+to.H/2)
	if !slices.Equal(f.moves, [][2]int{{0, 2}}) || !slices.Equal(f.order, []string{"Bravo", "Charlie", "Alpha"}) {
		t.Fatalf("drag to the end: moves %v, order %q", f.moves, f.order)
	}

	f.disabled = true
	tt.Frame()
	drag(x, from.Y+from.H/2, to.Y+to.H/2)
	tt.Key(ui.Alt, ui.KeyDown)
	if len(f.moves) != 1 {
		t.Fatalf("a disabled list moved: %v", f.moves)
	}
}
