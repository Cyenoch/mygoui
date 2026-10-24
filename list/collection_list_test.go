package list_test

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// listBoxPixelNear reports whether the pixel is within tol of c.
func listBoxPixelNear(img *image.RGBA, x, y int, c ui.Color, tol uint8) bool {
	got := img.RGBAAt(x, y)
	near := func(got, want uint8) bool {
		d := int(got) - int(want)
		return d >= -int(tol) && d <= int(tol)
	}
	return near(got.R, c.R) && near(got.G, c.G) && near(got.B, c.B)
}

// listBoxIndicatorInk counts the pixels in a row's indicator slot that differ
// from the page.
func listBoxIndicatorInk(img *image.RGBA, row ui.Rect, bg ui.Color) int {
	n := 0
	for y := int(row.Y); y < int(row.Y+row.H); y++ {
		for x := int(row.X + row.W - 24); x < int(row.X+row.W-8); x++ {
			if !listBoxPixelNear(img, x, y, bg, 1) {
				n++
			}
		}
	}
	return n
}

// TestListBoxChosenRowKeepsSurfaceAndShowsIndicator checks HeroUI's chosen row
// has no fill of its own: the checklist at the row's end is the only change.
func TestListBoxChosenRowKeepsSurfaceAndShowsIndicator(t *testing.T) {
	var state list.ListBoxState
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(20).
			Background(colors.Background).Children(func() {
			list.ListBox(c, list.ListBoxOptions{
				Label: "Fruit", State: &state, Selection: list.SelectionSingle,
				Items: []list.ListBoxItem{
					{Key: "apple", Label: "Apple"},
					{Key: "banana", Label: "Banana"},
				},
				Height: 200,
			})
		})
	}, 360, 260)
	tt.Frame()
	if err := tt.Click("Banana"); err != nil {
		t.Fatal(err)
	}
	tt.Move(2, 2)
	tt.Advance(400 * time.Millisecond)
	img := tt.Image()
	chosen, _ := tt.Find("Banana")
	other, _ := tt.Find("Apple")
	// The row's left padding shows the page, not the toolkit's accent face.
	if !listBoxPixelNear(img, int(chosen.X)+4, int(chosen.Y+chosen.H/2), colors.Background, 2) {
		got := img.RGBAAt(int(chosen.X)+4, int(chosen.Y+chosen.H/2))
		t.Fatalf("chosen row fill = %d,%d,%d, want the page background", got.R, got.G, got.B)
	}
	if listBoxPixelNear(img, int(chosen.X)+4, int(chosen.Y+chosen.H/2), colors.Accent, 6) {
		t.Fatal("chosen row shows the native accent face")
	}
	if ink := listBoxIndicatorInk(img, chosen, colors.Background); ink == 0 {
		t.Fatal("chosen row has no checklist in its indicator slot")
	}
	if ink := listBoxIndicatorInk(img, other, colors.Background); ink != 0 {
		t.Fatalf("unchosen row paints %d indicator pixels", ink)
	}
}

// TestListBoxSectionHeaderKeyboardNavigation checks keyboard movement skips
// section headers and disabled rows, and chooses the next item.
func TestListBoxSectionHeaderKeyboardNavigation(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Food", Selection: list.SelectionSingle, State: &state,
			Sections: []list.ListBoxSection{
				{Label: "Fruit", Items: []list.ListBoxItem{
					{Key: "apple", Label: "Apple"},
					{Key: "cherry", Label: "Cherry", Disabled: true},
				}},
				{Label: "Vegetables", Items: []list.ListBoxItem{
					{Key: "carrot", Label: "Carrot"},
					{Key: "kale", Label: "Kale"},
				}},
			},
			Height: 300,
		})
	}, 320, 320)
	tt.Frame()
	if err := tt.Click("Apple"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	if key := state.HighlightKey(); key != "carrot" {
		t.Fatalf("down from apple highlighted %q, want carrot", key)
	}
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "carrot" {
		t.Fatalf("keyboard choice = %q, want carrot", state.SelectedKey)
	}
	tt.Key(0, ui.KeyEnd)
	if key := state.HighlightKey(); key != "kale" {
		t.Fatalf("End highlighted %q, want kale", key)
	}
	tt.Key(0, ui.KeyHome)
	if key := state.HighlightKey(); key != "apple" {
		t.Fatalf("Home highlighted %q, want apple", key)
	}
}

// TestListBoxSectionRowsSitFlush checks HeroUI's gap-0 section: a header and
// the items under it touch, while the blocks between sections are spaced.
func TestListBoxSectionRowsSitFlush(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Food", State: &state,
			Sections: []list.ListBoxSection{
				{Label: "Fruit", Items: []list.ListBoxItem{
					{Key: "apple", Label: "Apple"},
					{Key: "banana", Label: "Banana"},
				}},
				{Label: "Vegetables", Items: []list.ListBoxItem{
					{Key: "carrot", Label: "Carrot"},
				}},
			},
			Height: 320,
		})
	}, 320, 360)
	tt.Frame()
	header, ok := tt.Find("Fruit")
	if !ok {
		t.Fatalf("section header missing: %v", tt.Texts())
	}
	apple, ok := tt.Find("Apple")
	if !ok {
		t.Fatalf("first item missing: %v", tt.Texts())
	}
	banana, ok := tt.Find("Banana")
	if !ok {
		t.Fatalf("second item missing: %v", tt.Texts())
	}
	if got := apple.Y - (header.Y + header.H); got < -0.6 || got > 0.6 {
		t.Fatalf("header to first item gap = %v, want 0", got)
	}
	if got := banana.Y - (apple.Y + apple.H); got < -0.6 || got > 0.6 {
		t.Fatalf("item gap inside a section = %v, want 0", got)
	}
	veg, ok := tt.Find("Vegetables")
	if !ok {
		t.Fatalf("second section missing: %v", tt.Texts())
	}
	if veg.Y <= banana.Y+banana.H-4 {
		t.Fatalf("sections are not separated: vegetables at %v, banana ends %v", veg.Y, banana.Y+banana.H)
	}
}

// TestListBoxSectionHeaderScrollsAway checks a section header scrolls with
// its rows, as HeroUI's do, instead of pinning over them at the viewport top.
func TestListBoxSectionHeaderScrollsAway(t *testing.T) {
	items := make([]list.ListBoxItem, 40)
	for i := range items {
		items[i] = list.ListBoxItem{Key: "k" + fmt.Sprint(i), Label: fmt.Sprintf("Item %02d", i)}
	}
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label:    "Food",
			State:    &state,
			Sections: []list.ListBoxSection{{Label: "Fruit", Items: items}},
			Height:   200,
		})
	}, 320, 260)
	tt.Frame()
	if _, ok := tt.Find("Fruit"); !ok {
		t.Fatalf("section header missing: %v", tt.Texts())
	}
	tt.Move(160, 120)
	tt.Wheel(160, 120, 0, 200)
	tt.Advance(200 * time.Millisecond)
	if _, ok := tt.Find("Fruit"); ok {
		t.Fatal("section header is pinned after scrolling")
	}
	first, _ := state.List.Visible()
	if first == 0 {
		t.Fatalf("list did not scroll: visible from %d", first)
	}
}

// TestListBoxTypeaheadSkipsHeadersAndDisabled checks typing a letter finds the
// next enabled item, not a section header or a disabled row.
func TestListBoxTypeaheadSkipsHeadersAndDisabled(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Food", Selection: list.SelectionSingle, State: &state,
			Sections: []list.ListBoxSection{
				{Label: "Citrus", Items: []list.ListBoxItem{
					{Key: "apple", Label: "Apple"},
					{Key: "cherry", Label: "Cherry", Disabled: true},
				}},
				{Label: "Greens", Items: []list.ListBoxItem{
					{Key: "carrot", Label: "Carrot"},
				}},
			},
			Height: 300,
		})
	}, 320, 320)
	tt.Frame()
	if err := tt.Click("Apple"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyC)
	if key := state.HighlightKey(); key != "carrot" {
		t.Fatalf("typing c highlighted %q, want carrot", key)
	}
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "carrot" {
		t.Fatalf("typeahead choice = %q, want carrot", state.SelectedKey)
	}
}

// TestListBoxIndicatorEndPadding checks the row reserves HeroUI's pe-7 for the
// indicator, so trailing content sits 28 DIPs from the row's end.
func TestListBoxIndicatorEndPadding(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Files", State: &state, Selection: list.SelectionSingle,
			Items: []list.ListBoxItem{{
				Key: "new", Label: "New file",
				End: func() { ui.Text(c, "TRLR") },
			}},
			Height: 120,
		})
	}, 320, 200)
	tt.Frame()
	row, ok := tt.Find("New file")
	if !ok {
		t.Fatalf("row missing: %v", tt.Texts())
	}
	tail, ok := tt.Find("TRLR")
	if !ok {
		t.Fatalf("trailing content missing: %v", tt.Texts())
	}
	if got, want := (row.X+row.W)-(tail.X+tail.W), float32(28); got < want-0.6 || got > want+0.6 {
		t.Fatalf("trailing content ends %v from the row end, want %v", got, want)
	}
}

// TestListBoxEmptyShowsMessage covers the zero-row path.
func TestListBoxEmptyShowsMessage(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Fruit", State: &state, Selection: list.SelectionSingle,
			Empty: "Nothing here", Height: 160,
		})
	}, 320, 220)
	tt.Frame()
	if !tt.HasText("Nothing here") {
		t.Fatalf("empty message missing: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	if state.HighlightKey() != "" {
		t.Fatalf("empty list highlighted %q", state.HighlightKey())
	}
}
