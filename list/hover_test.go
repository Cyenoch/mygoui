package list_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func listHoverPixel(img *image.RGBA, x, y int) (uint8, uint8, uint8) {
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B
}

func listHoverNear(got, want uint8) bool {
	d := int(got) - int(want)
	return d >= -3 && d <= 3
}

func wantListHoverPixel(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	r, g, b := listHoverPixel(img, x, y)
	if !listHoverNear(r, c.R) || !listHoverNear(g, c.G) || !listHoverNear(b, c.B) {
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d", what, x, y, r, g, b, c.R, c.G, c.B)
	}
}

// TestListBoxRowHoverPaint checks a row's Default hover fill follows the
// pointer and reverts. The list viewport paints nothing itself, so the page
// background shows between rows.
func TestListBoxRowHoverPaint(t *testing.T) {
	var colors theme.Palette
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			list.ListBox(c, list.ListBoxOptions{
				Label: "Fruit", State: &state,
				Items: []list.ListBoxItem{
					{Key: "apple", Label: "Apple"},
					{Key: "pear", Label: "Pear"},
				},
			})
		})
	}, 420, 220)

	row, ok := tt.Find("Apple")
	if !ok {
		t.Fatal("list row missing")
	}
	// Sample the row's left padding, off its label ink.
	x, y := int(row.X)+4, int(row.Y+row.H/2)
	wantListHoverPixel(t, tt.Image(), x, y, colors.Background, "list row normal fill")

	tt.Move(row.X+row.W/2, row.Y+row.H/2)
	tt.Frame()
	wantListHoverPixel(t, tt.Image(), x, y, colors.Accent.Over(colors.Background), "list row hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantListHoverPixel(t, tt.Image(), x, y, colors.Background, "list row fill after leaving")
}

func TestListBoxEditorHighlightPaintsOnceAndMoves(t *testing.T) {
	var state list.ListBoxState
	var colors theme.Palette
	var query string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(20).Background(colors.Background).Children(func() {
			input := ui.TextInputBase(c, &query).Label("Search").Height(32)
			input.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyDown {
					return false
				}
				state.MoveHighlight(1)
				return true
			})
			list.ListBox(c, list.ListBoxOptions{
				State: &state, KeyboardHighlight: input.Focused(), NoFocusOnClick: true,
				Items: []list.ListBoxItem{{Label: "Apple"}, {Label: "Pear"}},
			})
		})
	}, 320, 200)
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	tt.Move(2, 2)
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	if !tt.Focused("Search") || state.HighlightKey() != "Pear" {
		t.Fatalf("editor focus %v, highlight %q", tt.Focused("Search"), state.HighlightKey())
	}
	apple, _ := tt.Find("Apple")
	pear, _ := tt.Find("Pear")
	img := tt.Image()
	wantListHoverPixel(t, img, int(pear.X)+4, int(pear.Y+pear.H/2),
		colors.Accent.Over(colors.Background), "keyboard highlight")
	wantListHoverPixel(t, img, int(apple.X)+4, int(apple.Y+apple.H/2),
		colors.Background, "previous highlight")
}
