package navigation_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func navHoverPixel(img *image.RGBA, x, y int) (uint8, uint8, uint8) {
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B
}

func navHoverNear(got, want uint8) bool {
	d := int(got) - int(want)
	return d >= -3 && d <= 3
}

func wantNavHoverPixel(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	r, g, b := navHoverPixel(img, x, y)
	if !navHoverNear(r, c.R) || !navHoverNear(g, c.G) || !navHoverNear(b, c.B) {
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d", what, x, y, r, g, b, c.R, c.G, c.B)
	}
}

// TestTabsHoverDimsLabel checks an unselected tab follows the pointer the
// way upstream's opacity-70 does: the label dims and no face is painted,
// while the selected tab keeps its pill.
func TestTabsHoverDimsLabel(t *testing.T) {
	var colors theme.Palette
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			navigation.Tabs(c, navigation.TabsOptions{
				Selected: &selected,
				Items: []navigation.TabOptions{
					{Label: "Overview", Content: func() {}},
					{Label: "Settings", Content: func() {}},
				},
			})
		})
	}, 440, 200)

	tab, ok := tt.Find("Settings")
	if !ok {
		t.Fatal("unselected tab missing")
	}
	// The strip container paints Default behind an unselected tab.
	x, y := int(tab.X)+4, int(tab.Y+tab.H/2)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Default, "tab normal fill")
	normal := darkestIn(tt.Image(), tab)

	tt.Move(tab.X+tab.W/2, tab.Y+tab.H/2)
	tt.Frame()
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Default, "tab hover hides no face")
	hovered := darkestIn(tt.Image(), tab)
	if hovered <= normal+30 {
		t.Fatalf("hovered label ink sums %d against %d: the label did not dim", hovered, normal)
	}

	tt.Move(2, 2)
	tt.Frame()
	if got := darkestIn(tt.Image(), tab); got > normal+6 {
		t.Fatalf("label ink sums %d after leaving, want about %d", got, normal)
	}
}

// darkestIn returns the smallest color sum of the pixels inside r, away from
// its edges, as the label's ink is the darkest thing there.
func darkestIn(img *image.RGBA, r ui.Rect) int {
	best := 3 * 255
	for y := int(r.Y) + 8; y < int(r.Y+r.H)-8; y++ {
		for x := int(r.X) + 12; x < int(r.X+r.W)-12; x++ {
			c := img.RGBAAt(x, y)
			if sum := int(c.R) + int(c.G) + int(c.B); sum < best {
				best = sum
			}
		}
	}
	return best
}

// TestAccordionHeaderHoverPaint checks a closed section's header paint follows
// the trigger's hover and reverts. An open section suppresses the face.
func TestAccordionHeaderHoverPaint(t *testing.T) {
	var colors theme.Palette
	var expanded []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			navigation.Accordion(c, navigation.AccordionOptions{
				Expanded: &expanded,
				Items: []navigation.AccordionItemOptions{
					{Title: "Shipping", Body: func() {}},
					{Title: "Payment", Body: func() {}},
				},
			})
		})
	}, 440, 300)

	header, ok := tt.Find("Shipping")
	if !ok {
		t.Fatal("accordion header missing")
	}
	// The default variant paints the 3% foreground hover over the page.
	normal := colors.Background
	hover := colors.Foreground.Alpha(0.03).Over(normal)
	// Sample right of the title, away from its ink and the chevron.
	x, y := int(header.X+header.W/2), int(header.Y+header.H/2)
	wantNavHoverPixel(t, tt.Image(), x, y, normal, "accordion header normal fill")

	tt.Move(header.X+header.W/2, header.Y+header.H/2)
	tt.Frame()
	wantNavHoverPixel(t, tt.Image(), x, y, hover, "accordion header hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantNavHoverPixel(t, tt.Image(), x, y, normal, "accordion header fill after leaving")
}

// TestDisclosureHoverIsInert pins the upstream boundary: HeroUI's
// disclosure.css defines no hover face (only accordion headers do), so the
// build-time snapshot audit finds nothing to repaint here. The trigger still
// opens the panel.
func TestDisclosureHoverIsInert(t *testing.T) {
	var colors theme.Palette
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			navigation.Disclosure(c, navigation.DisclosureOptions{
				Title: "Details", Open: &open, Body: func() {
					ui.Text(c, "Body")
				},
			})
		})
	}, 440, 240)

	trigger, ok := tt.Find("Details")
	if !ok {
		t.Fatal("disclosure trigger missing")
	}
	x, y := int(trigger.X+trigger.W/2), int(trigger.Y+trigger.H/2)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Background, "disclosure normal fill")

	tt.Move(trigger.X+trigger.W/2, trigger.Y+trigger.H/2)
	tt.Frame()
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Background, "disclosure hover fill")

	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("the disclosure trigger did not open")
	}
}
