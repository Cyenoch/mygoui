package mygoui_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// wantPixel fails unless the pixel at x, y is c.
func wantPixel(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	r, g, b := pixelAt(img, x, y)
	if !near(r, c.R) || !near(g, c.G) || !near(b, c.B) {
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d", what, x, y, r, g, b, c.R, c.G, c.B)
	}
}

// TestToggleButtonGroupFillStaysInsideRoundedBox checks that a joined
// toggle's selected, hovered and pressed fills never square off the rounded
// outer corner of the group, while a mid-group button's square corners are
// filled.
func TestToggleButtonGroupFillStaysInsideRoundedBox(t *testing.T) {
	var colors theme.Palette
	var one, two, three bool
	one, two, three = true, true, true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{},
				mygoui.ToggleButtonOptions{Label: "One", Variant: mygoui.ToggleOutline, Selected: &one},
				mygoui.ToggleButtonOptions{Label: "Two", Variant: mygoui.ToggleOutline, Selected: &two},
				mygoui.ToggleButtonOptions{Label: "Three", Variant: mygoui.ToggleOutline, Selected: &three},
			)
		})
	}, 360, 160)

	// A point in the button's rounded corner, outside its fill, and one on
	// its fill, left of the text.
	outside := func(r ui.Rect, left bool) (int, int) {
		if left {
			return int(r.X) + 1, int(r.Y) + 1
		}
		return int(r.X+r.W) - 2, int(r.Y) + 1
	}
	inside := func(r ui.Rect) (int, int) {
		return int(r.X) + 6, int(r.Y) + int(r.H)/2
	}

	for _, tc := range []struct {
		label     string
		left      bool
		cornerOut bool
	}{
		{"One", true, true},
		{"Three", false, true},
		{"Two", false, false}, // mid-group corners are square
	} {
		rect, ok := tt.Find(tc.label)
		if !ok {
			t.Fatalf("%q missing", tc.label)
		}
		// Each translucent state face composites once against the surface.
		selected := colors.Input.Alpha(0.64).Over(colors.Background)
		hovered := colors.Input.Over(colors.Background)
		tt.Frame()
		img := tt.Image()
		x, y := outside(rect, tc.left)
		if tc.cornerOut {
			// Selected, not hovered, paints the rounded fill.
			wantPixel(t, img, x, y, colors.Background, tc.label+" selected corner")
		}
		ix, iy := inside(rect)
		wantPixel(t, img, ix, iy, selected, tc.label+" selected fill")

		tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
		tt.Frame()
		img = tt.Image()
		if tc.cornerOut {
			wantPixel(t, img, x, y, colors.Background, tc.label+" hover corner")
		} else {
			// The square corner of a mid-group button takes the fill.
			wantPixel(t, img, int(rect.X)+1, int(rect.Y)+1, hovered, tc.label+" square corner")
		}
		wantPixel(t, img, ix, iy, hovered, tc.label+" hover fill")

		tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
		tt.Frame()
		img = tt.Image()
		if tc.cornerOut {
			wantPixel(t, img, x, y, colors.Background, tc.label+" pressed corner")
		}
		wantPixel(t, img, ix, iy, hovered, tc.label+" pressed fill")
		tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	}

	// The fill layer must not swallow the button's clicks.
	if err := tt.Click("One"); err != nil {
		t.Fatal(err)
	}
	if !one {
		t.Fatal("a click through the fill layer did not toggle the button")
	}
}

// TestButtonGroupFillStaysInsideRoundedBox checks the same contract on the
// primary ButtonGroup path.
func TestButtonGroupFillStaysInsideRoundedBox(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{},
				mygoui.ButtonOptions{Label: "First", OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Second", OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Third", OnClick: func() {}},
			)
		})
	}, 360, 160)

	for _, tc := range []struct {
		label string
		left  bool
	}{{"First", true}, {"Third", false}} {
		rect, ok := tt.Find(tc.label)
		if !ok {
			t.Fatalf("%q missing", tc.label)
		}
		var x, y int
		if tc.left {
			x, y = int(rect.X)+1, int(rect.Y)+1
		} else {
			x, y = int(rect.X+rect.W)-2, int(rect.Y)+1
		}
		ix, iy := int(rect.X)+6, int(rect.Y)+int(rect.H)/2

		tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
		tt.Frame()
		wantPixel(t, tt.Image(), x, y, colors.Background, tc.label+" hover corner")
		wantPixel(t, tt.Image(), ix, iy, colors.Primary.Alpha(0.9).Over(colors.Background), tc.label+" hover fill")

		tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
		tt.Frame()
		wantPixel(t, tt.Image(), x, y, colors.Background, tc.label+" pressed corner")
		wantPixel(t, tt.Image(), ix, iy, colors.Primary.Alpha(0.9).Over(colors.Background), tc.label+" pressed fill")
		tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	}
}

// textSpan returns the horizontal extent of the label ink in r, the glyphs
// on the button's face. Ink is what differs from the face color; the page
// background before the rounded corners antialiases differently per
// platform and is not ink, so the scan stays between the row's face
// pixels.
func textSpan(img *image.RGBA, r ui.Rect, face ui.Color) (int, int, bool) {
	painted := func(x, y int) bool {
		c := img.RGBAAt(x, y)
		return near(c.R, face.R) && near(c.G, face.G) && near(c.B, face.B)
	}
	x0, x1 := -1, -1
	for y := int(r.Y) + int(r.H)/4; y < int(r.Y+r.H)-int(r.H)/4; y++ {
		left, right := -1, -1
		for x := int(r.X) + 1; x < int(r.X+r.W)-1; x++ {
			if painted(x, y) {
				if left < 0 {
					left = x
				}
				right = x
			}
		}
		if left < 0 {
			continue
		}
		for x := left; x <= right; x++ {
			if painted(x, y) {
				continue
			}
			if x0 < 0 || x < x0 {
				x0 = x
			}
			if x > x1 {
				x1 = x
			}
		}
	}
	return x0, x1, x0 >= 0
}

// darkSpan returns the horizontal extent of dark pixels in r, an icon on a
// light face.
func darkSpan(img *image.RGBA, r ui.Rect) (int, int, bool) {
	x0, x1 := -1, -1
	for y := int(r.Y) + int(r.H)/4; y < int(r.Y+r.H)-int(r.H)/4; y++ {
		for x := int(r.X) + 1; x < int(r.X+r.W)-1; x++ {
			c := img.RGBAAt(x, y)
			// The icon's ink on the light selected face.
			if c.R < 150 && c.G < 180 {
				if x0 < 0 {
					x0 = x
				}
				x1 = x
			}
		}
	}
	return x0, x1, x0 >= 0
}

// TestFullWidthToggleGroupFillsItsShare checks a full-width joined toggle
// set: equal shares, adjacent boxes, the state fill across the whole
// segment and no fill outside its rounded end.
func TestFullWidthToggleGroupFillsItsShare(t *testing.T) {
	var colors theme.Palette
	var one, two, three bool
	one, two, three = true, true, true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{FullWidth: true},
				mygoui.ToggleButtonOptions{Label: "One", Variant: mygoui.ToggleOutline, Selected: &one},
				mygoui.ToggleButtonOptions{Label: "Two", Variant: mygoui.ToggleOutline, Selected: &two},
				mygoui.ToggleButtonOptions{Label: "Three", Variant: mygoui.ToggleOutline, Selected: &three},
			)
		})
	}, 360, 160)

	rects := make([]ui.Rect, 0, 3)
	for _, label := range []string{"One", "Two", "Three"} {
		rect, ok := tt.Find(label)
		if !ok {
			t.Fatalf("%q missing", label)
		}
		rects = append(rects, rect)
	}
	want := (360 - 2*24) / 3
	for i, r := range rects {
		if d := r.W - float32(want); d < -1 || d > 1 {
			t.Fatalf("button %d width = %v, want %d (rects %+v)", i, r.W, want, rects)
		}
	}
	if float32(rects[0].X+rects[0].W) < rects[1].X-0.5 || float32(rects[0].X+rects[0].W) > rects[1].X+0.5 {
		t.Fatalf("buttons are not adjacent: %+v", rects)
	}
	if rects[0].X < 23 || rects[0].X > 25 || rects[2].X+rects[2].W < 335 || rects[2].X+rects[2].W > 337 {
		t.Fatalf("group does not span its column: %+v", rects)
	}

	first := rects[0]
	tt.Move(first.X+first.W/2, first.Y+first.H/2)
	tt.Frame()
	img := tt.Image()
	mid := int(first.Y + first.H/2)
	wantPixel(t, img, int(first.X)+1, int(first.Y)+1, colors.Background, "full-width outer corner")
	wantPixel(t, img, int(first.X+first.W)-3, mid, colors.Input.Over(colors.Background), "fill at the segment's end")
	wantPixel(t, img, int(first.X)+6, mid, colors.Input.Over(colors.Background), "fill at the segment's start")
}

// TestFullWidthButtonGroupAlignsLabels checks that the fill layer does not
// change label placement in a full-width group: Button.Align still centers,
// starts and ends the label.
func TestFullWidthButtonGroupAlignsLabels(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{FullWidth: true},
				mygoui.ButtonOptions{Label: "Center", OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Start", Align: mygoui.ButtonAlignStart, OnClick: func() {}},
				mygoui.ButtonOptions{Label: "End", Align: mygoui.ButtonAlignEnd, OnClick: func() {}},
			)
		})
	}, 480, 160)

	img := tt.Image()
	rects := make([]ui.Rect, 0, 3)
	for _, label := range []string{"Center", "Start", "End"} {
		rect, ok := tt.Find(label)
		if !ok {
			t.Fatalf("%q missing", label)
		}
		rects = append(rects, rect)
		x0, x1, ok := textSpan(img, rect, colors.Primary)
		if !ok {
			t.Fatalf("%q has no label pixels in %+v", label, rect)
		}
		center := float32(x0+x1) / 2
		switch label {
		case "Center":
			if d := center - (rect.X + rect.W/2); d < -2 || d > 2 {
				t.Fatalf("centered label at %v, want %v (span %d..%d, rect %+v)", center, rect.X+rect.W/2, x0, x1, rect)
			}
		case "Start":
			if x0 < int(rect.X)+11 || x0 > int(rect.X)+13 {
				t.Fatalf("start label begins at %d, want near %v (rect %+v)", x0, rect.X+11, rect)
			}
		case "End":
			if x1 < int(rect.X+rect.W)-13 || x1 > int(rect.X+rect.W)-11 {
				t.Fatalf("end label ends at %d, want near %v (rect %+v)", x1, rect.X+rect.W-11, rect)
			}
		}
	}

	want := (480 - 2*24) / 3
	for i, r := range rects {
		if d := r.W - float32(want); d < -1 || d > 1 {
			t.Fatalf("button %d width = %v, want %d (rects %+v)", i, r.W, want, rects)
		}
	}
	first := rects[0]
	tt.Move(first.X+first.W/2, first.Y+first.H/2)
	tt.Frame()
	img = tt.Image()
	mid := int(first.Y + first.H/2)
	wantPixel(t, img, int(first.X)+1, int(first.Y)+1, colors.Background, "full-width button corner")
	wantPixel(t, img, int(first.X+first.W)-3, mid, colors.Primary.Alpha(0.9).Over(colors.Background), "full-width button fill")
}

// TestIconOnlyGroupedTogglesFillTheirSquare checks grouped icon-only
// toggles: a square box, the icon centered and the state fill across the
// square.
func TestIconOnlyGroupedTogglesFillTheirSquare(t *testing.T) {
	var colors theme.Palette
	var solo, one, two, three bool
	solo, one, two, three = true, true, true, true
	icon := lucide.ChevronDown()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
				Label: "Solo", Icon: icon, IconOnly: true, Variant: mygoui.ToggleOutline, Selected: &solo,
			})
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{},
				mygoui.ToggleButtonOptions{Label: "One", Icon: icon, IconOnly: true, Variant: mygoui.ToggleOutline, Selected: &one},
				mygoui.ToggleButtonOptions{Label: "Two", Icon: icon, IconOnly: true, Variant: mygoui.ToggleOutline, Selected: &two},
				mygoui.ToggleButtonOptions{Label: "Three", Icon: icon, IconOnly: true, Variant: mygoui.ToggleOutline, Selected: &three},
			)
		})
	}, 360, 200)

	rects := make([]ui.Rect, 0, 3)
	for _, label := range []string{"One", "Two", "Three"} {
		rect, ok := tt.Find(label)
		if !ok {
			t.Fatalf("%q missing", label)
		}
		if rect.W != 32 || rect.H != 32 {
			t.Fatalf("%q box = %+v, want a 32 DIP square", label, rect)
		}
		rects = append(rects, rect)
	}

	first := rects[0]
	mid := int(first.Y + first.H/2)
	// Compare with a standalone outline toggle before hovering, so both
	// faces and their antialiasing match.
	x0, x1, ok := darkSpan(tt.Image(), first)
	if !ok {
		t.Fatal("the icon is not painted")
	}
	soloRect, ok := tt.Find("Solo")
	if !ok {
		t.Fatal(`"Solo" missing`)
	}
	sx0, sx1, ok := darkSpan(tt.Image(), soloRect)
	if !ok {
		t.Fatal("the standalone icon is not painted")
	}
	groupX0, groupX1 := x0-int(first.X), x1-int(first.X)
	soloX0, soloX1 := sx0-int(soloRect.X), sx1-int(soloRect.X)
	// Joining removes an interior outline border; half-pixel centering can
	// shift the rasterized glyph by one pixel without changing its alignment.
	if groupX0-soloX0 < -1 || groupX0-soloX0 > 1 || groupX1-soloX1 < -1 || groupX1-soloX1 > 1 {
		t.Fatalf("grouped icon span %d..%d in %+v differs from standalone %d..%d in %+v",
			groupX0, groupX1, first, soloX0, soloX1, soloRect)
	}

	tt.Move(first.X+first.W/2, first.Y+first.H/2)
	tt.Frame()
	img := tt.Image()
	wantPixel(t, img, int(first.X)+1, int(first.Y)+1, colors.Background, "icon-only outer corner")
	wantPixel(t, img, int(first.X+first.W)-3, mid, colors.Input.Over(colors.Background), "icon-only fill")

	mid2 := rects[1]
	tt.Move(mid2.X+mid2.W/2, mid2.Y+mid2.H/2)
	tt.Frame()
	img = tt.Image()
	wantPixel(t, img, int(mid2.X)+1, int(mid2.Y)+1, colors.Input.Over(colors.Background), "icon-only square corner")
}

// TestGroupedFocusRingIsInset checks that joined buttons draw the focus
// ring inside their box, where the next segment cannot cover it, and that
// the state fill inside the ring is unchanged.
func TestGroupedFocusRingIsInset(t *testing.T) {
	var colors theme.Palette
	var one, two, three bool
	one, two, three = true, true, true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{},
				mygoui.ToggleButtonOptions{Label: "One", Variant: mygoui.ToggleOutline, Selected: &one},
				mygoui.ToggleButtonOptions{Label: "Two", Variant: mygoui.ToggleOutline, Selected: &two},
				mygoui.ToggleButtonOptions{Label: "Three", Variant: mygoui.ToggleOutline, Selected: &three},
			)
		})
	}, 360, 160)

	tt.Key(0, ui.KeyTab)
	if !tt.Focused("One") {
		t.Fatal("Tab did not focus the first toggle")
	}
	rect, ok := tt.Find("One")
	if !ok {
		t.Fatal(`"One" missing`)
	}
	tt.Frame()
	img := tt.Image()
	mid := int(rect.Y + rect.H/2)
	selected := colors.Input.Alpha(0.64).Over(colors.Background)
	wantPixel(t, img, int(rect.X)+2, mid, colors.Focus, "inset ring")
	wantPixel(t, img, int(rect.X)-2, mid, colors.Background, "outward ring")
	wantPixel(t, img, int(rect.X+rect.W)-2, mid, colors.Focus, "shared-edge ring")
	wantPixel(t, img, int(rect.X)+6, mid, selected, "fill inside the ring")

	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Two") {
		t.Fatal("Tab did not move the focus to the second toggle")
	}
	tt.Frame()
	img = tt.Image()
	second, ok := tt.Find("Two")
	if !ok {
		t.Fatal(`"Two" missing`)
	}
	wantPixel(t, img, int(second.X)+2, mid, colors.Focus, "second inset ring")
	wantPixel(t, img, int(second.X)-2, mid, selected, "previous button has no ring")
}
