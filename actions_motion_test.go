package mygoui_test

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// wantPixelWith fails unless the pixel at x, y is within 6 of c, for samples
// that a translucent face composites into a computed color.
func wantPixelWith(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	r, g, b := pixelAt(img, x, y)
	if absDiff(r, c.R) > 6 || absDiff(g, c.G) > 6 || absDiff(b, c.B) > 6 {
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d", what, x, y, r, g, b, c.R, c.G, c.B)
	}
}

func absDiff(a, b uint8) int {
	d := int(a) - int(b)
	if d < 0 {
		return -d
	}
	return d
}

func matchesPixel(img *image.RGBA, x, y int, c ui.Color) bool {
	r, g, b := pixelAt(img, x, y)
	return absDiff(r, c.R) <= 6 && absDiff(g, c.G) <= 6 && absDiff(b, c.B) <= 6
}

// TestButtonPressKeepsGeometryAndActivatesOnRelease checks the native press
// contract with and without reduced motion: the face changes, the target stays
// put, and only a release inside the target activates it.
func TestButtonPressKeepsGeometryAndActivatesOnRelease(t *testing.T) {
	for _, reduced := range []bool{false, true} {
		name := "normal motion"
		if reduced {
			name = "reduced motion"
		}
		t.Run(name, func(t *testing.T) {
			var colors theme.Palette
			calls := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).Fill().Padding(24).Background(colors.Background).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Save", FullWidth: true, OnClick: func() { calls++ },
					})
				})
			}, 480, 160)
			tt.SetPreferences(ui.Preferences{ReduceMotion: reduced, TextScale: 1})
			rect := findRect(t, tt, "Save")
			if rect.H != 32 {
				t.Fatalf("medium button height = %v, want 32 DIP", rect.H)
			}
			x, y := int(rect.X)+4, int(rect.Y+rect.H/2)
			wantPixel(t, tt.Image(), x, y, colors.Primary, "resting face")
			cx, cy := rect.X+rect.W/2, rect.Y+rect.H/2
			tt.Press(cx, cy)
			tt.Advance(300 * time.Millisecond)
			if got := findRect(t, tt, "Save"); got != rect {
				t.Fatalf("pressed geometry = %+v, want %+v", got, rect)
			}
			wantPixel(t, tt.Image(), x, y, colors.Primary.Alpha(0.9).Over(colors.Background), "pressed face")
			if calls != 0 {
				t.Fatal("button activated before release")
			}
			tt.Release(cx, cy)
			if calls != 1 || !tt.Focused("Save") {
				t.Fatalf("released button calls=%d focus=%v", calls, tt.Focused("Save"))
			}
			tt.Press(cx, cy)
			tt.Move(2, 2)
			tt.Release(2, 2)
			if calls != 1 {
				t.Fatalf("release outside activated button: calls=%d", calls)
			}
			wantPixel(t, tt.Image(), x, y, colors.Primary, "restored face")
			if got := findRect(t, tt, "Save"); got != rect {
				t.Fatalf("released geometry = %+v, want %+v", got, rect)
			}
			if tt.FramePending() {
				t.Fatal("settled button asks for animation frames")
			}
		})
	}
}

// TestButtonPendingKeepsFaceAndBlocksActivation checks the native pending
// contract: 50% disabled dimming, no pointer face changes or activation.
func TestButtonPendingKeepsFaceAndBlocksActivation(t *testing.T) {
	var colors theme.Palette
	calls := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Saving", FullWidth: true, Loading: true,
				OnClick: func() { calls++ },
			})
		})
	}, 480, 160)
	tt.Frame()
	rect, ok := tt.Find("Saving")
	if !ok {
		t.Fatal("pending button missing")
	}
	y := int(rect.Y + rect.H/2)
	left := func(dx int) int { return int(rect.X) + dx }

	tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Frame()
	wantPixel(t, tt.Image(), left(1), y, colors.Primary.Alpha(0.5).Over(colors.Background), "pending hover face")
	tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Advance(300 * time.Millisecond)
	wantPixel(t, tt.Image(), left(1), y, colors.Primary.Alpha(0.5).Over(colors.Background), "pending pressed face")
	wantPixel(t, tt.Image(), left(4), y, colors.Primary.Alpha(0.5).Over(colors.Background), "pending edge remains filled")
	tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	if calls != 0 {
		t.Fatalf("pending button activated %d times", calls)
	}
}

// TestDisabledTogglePaintsRestFace checks a disabled control reads its rest
// face while hovered: the disabled flag is set before the face samples the
// pointer, matching Coss's pointer-events:none.
func TestDisabledTogglePaintsRestFace(t *testing.T) {
	var colors theme.Palette
	off := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
				Label: "Locked", Selected: &off, Disabled: true,
			})
		})
	}, 260, 140)
	tt.Frame()
	rect, ok := tt.Find("Locked")
	if !ok {
		t.Fatal("disabled toggle missing")
	}
	x, y := int(rect.X)+4, int(rect.Y+rect.H/2)
	wantPixelWith(t, tt.Image(), x, y,
		colors.Background, "disabled rest face")
	tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Frame()
	wantPixelWith(t, tt.Image(), x, y,
		colors.Background, "disabled hover face")
	if err := tt.Click("Locked"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeySpace)
	if off || tt.Focused("Locked") {
		t.Fatalf("disabled toggle selected=%v focus=%v", off, tt.Focused("Locked"))
	}
}

// TestButtonGroupSeparatorsKeepLayoutAndPaintSeams checks the shared-edge
// divider takes no layout width and paints Coss's current-color line,
// and that joined buttons keep their size while pressed.
func TestButtonGroupSeparatorsKeepLayoutAndPaintSeams(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Separators: true},
				mygoui.ButtonOptions{Label: "First", Variant: mygoui.ButtonSecondary, OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Second", Variant: mygoui.ButtonSecondary, OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Third", Variant: mygoui.ButtonSecondary, OnClick: func() {}},
			)
		})
	}, 480, 160)
	tt.Frame()
	first, ok := tt.Find("First")
	if !ok {
		t.Fatal("first segment missing")
	}
	second, ok := tt.Find("Second")
	if !ok {
		t.Fatal("second segment missing")
	}
	if gap := second.X - (first.X + first.W); gap < -0.5 || gap > 0.5 {
		t.Fatalf("separator took layout width: gap %v", gap)
	}
	y := int(second.Y + second.H/2)
	// A pixel far from the shared edge stays on the plain face; one near
	// it carries the seam's ink. The seam's 1 DIP box can straddle a pixel
	// edge, so match any clear departure from the face.
	img := tt.Image()
	face := colors.Secondary.Over(colors.Background)
	wantPixel(t, img, int(second.X)-4, y, face, "plain face")
	found := false
	for x := int(second.X) - 2; x <= int(second.X)+1; x++ {
		r, g, b := pixelAt(img, x, y)
		if absDiff(r, face.R) > 8 || absDiff(g, face.G) > 8 || absDiff(b, face.B) > 8 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("no seam near x=%v on the shared edge", second.X)
	}

	tt.Press(second.X+second.W/2, second.Y+second.H/2)
	if tt.FramePending() {
		t.Fatal("a settled pressed group schedules animation frames")
	}
	wantPixel(t, tt.Image(), int(second.X)+1, y, colors.Secondary.Alpha(0.8).Over(colors.Background), "joined pressed face")
	tt.Release(second.X+second.W/2, second.Y+second.H/2)
}

// TestButtonGroupRightToLeftRoundsStartAndPaintsSeam checks the group's
// direction-dependent corners and seam side.
func TestButtonGroupRightToLeftRoundsStartAndPaintsSeam(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Separators: true},
				mygoui.ButtonOptions{Label: "First", Variant: mygoui.ButtonSecondary, OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Second", Variant: mygoui.ButtonSecondary, OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Third", Variant: mygoui.ButtonSecondary, OnClick: func() {}},
			)
		})
	}, 480, 160)
	tt.SetDirection(ui.RightToLeft)
	tt.Frame()
	first, ok := tt.Find("First")
	if !ok {
		t.Fatal("first segment missing")
	}
	third, ok := tt.Find("Third")
	if !ok {
		t.Fatal("third segment missing")
	}
	second, ok := tt.Find("Second")
	if !ok {
		t.Fatal("second segment missing")
	}
	if first.X <= third.X {
		t.Fatalf("row did not lay out right to left: first=%v third=%v", first.X, third.X)
	}
	// The start segment is the right one: its start (right) corners round.
	wantPixel(t, tt.Image(), int(first.X)+1, int(first.Y)+1, colors.Secondary.Over(colors.Background), "RTL square end corner")
	wantPixel(t, tt.Image(), int(first.X+first.W)-2, int(first.Y)+1, colors.Background, "RTL rounded start corner")
	if gap := first.X - (second.X + second.W); gap < -0.5 || gap > 0.5 {
		t.Fatalf("RTL separator took layout width: gap %v", gap)
	}
	seam := false
	x, y := int(second.X+second.W), int(second.Y+second.H/2)
	img := tt.Image()
	face := colors.Secondary.Over(colors.Background)
	for sx := x - 1; sx <= x+2; sx++ {
		r, g, b := pixelAt(img, sx, y)
		if absDiff(r, face.R) > 8 || absDiff(g, face.G) > 8 || absDiff(b, face.B) > 8 {
			seam = true
			break
		}
	}
	if !seam {
		t.Fatalf("no RTL seam near x=%v", second.X+second.W)
	}
}

// TestTagSelectionFadesAndRemovalExits checks a tag fades to its selected
// face and a removed tag paints an inert exit copy that fades out.
func TestTagSelectionFadesAndRemovalExits(t *testing.T) {
	var colors theme.Palette
	state := &mygoui.TagGroupState{}
	var tags []mygoui.TagOptions
	tags = []mygoui.TagOptions{
		{Key: "keep", Label: "Keep"},
		{Key: "go", Label: "Go", Removable: true, OnRemove: func() {
			tags = []mygoui.TagOptions{{Key: "keep", Label: "Keep"}}
		}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.TagGroup(c, mygoui.TagGroupOptions{
				Selection: mygoui.TagSelectMultiple, State: state,
			}, tags)
		})
	}, 360, 160)
	tt.Frame()
	goRect, ok := tt.Find("Go")
	if !ok {
		t.Fatal("tag missing")
	}
	x, y := int(goRect.X)+5, int(goRect.Y+goRect.H/2)
	wantPixel(t, tt.Image(), x, y, colors.Default.Over(colors.Background), "tag rest face")

	if err := tt.Click("Go"); err != nil {
		t.Fatal(err)
	}
	// The selection transition starts from the face on screen.
	wantPixel(t, tt.Image(), x, y, colors.Default.Over(colors.Background), "tag face at transition start")
	tt.Advance(200 * time.Millisecond)
	wantPixelWith(t, tt.Image(), x, y,
		colors.AccentSoftHover.Over(colors.Background), "tag selected hover face")
	tt.Move(2, 2)
	tt.Advance(200 * time.Millisecond)
	wantPixelWith(t, tt.Image(), x, y,
		colors.AccentSoft.Over(colors.Background), "tag selected face")

	if err := tt.Click("Remove"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Go") {
		t.Fatal("removed tag still built")
	}
	// The exit copy paints the last frame in place, then fades out.
	if matchesPixel(tt.Image(), x, y, colors.Background) {
		t.Fatal("removed tag left without an exit copy")
	}
	tt.Advance(200 * time.Millisecond)
	if !matchesPixel(tt.Image(), x, y, colors.Background) {
		t.Fatal("the removal exit did not settle")
	}
}

// TestToggleButtonGroupSeparatorsPaintSeams checks the toggle group reaches
// Coss's shared-edge divider without taking layout width.
func TestToggleButtonGroupSeparatorsPaintSeams(t *testing.T) {
	var colors theme.Palette
	off := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{Separators: true},
				mygoui.ToggleButtonOptions{Label: "Left", Variant: mygoui.ToggleOutline, Selected: &off},
				mygoui.ToggleButtonOptions{Label: "Right", Variant: mygoui.ToggleOutline, Selected: &off},
			)
		})
	}, 360, 140)
	tt.Frame()
	left, ok := tt.Find("Left")
	if !ok {
		t.Fatal("left toggle missing")
	}
	right, ok := tt.Find("Right")
	if !ok {
		t.Fatal("right toggle missing")
	}
	if gap := right.X - (left.X + left.W); gap < -0.5 || gap > 0.5 {
		t.Fatalf("toggle separator took layout width: gap %v", gap)
	}
	y := int(right.Y + right.H/2)
	found := false
	img := tt.Image()
	for x := int(right.X) - 2; x <= int(right.X)+1; x++ {
		r, g, b := pixelAt(img, x, y)
		if absDiff(r, colors.Background.R) > 8 || absDiff(g, colors.Background.G) > 8 ||
			absDiff(b, colors.Background.B) > 8 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("no toggle seam near x=%v", right.X)
	}
}

// TestTagSelectionHonorsReducedMotion checks reduced motion settles the
// selection face at once.
func TestTagSelectionHonorsReducedMotion(t *testing.T) {
	var colors theme.Palette
	state := &mygoui.TagGroupState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.TagGroup(c, mygoui.TagGroupOptions{
				Selection: mygoui.TagSelectMultiple, State: state,
			}, []mygoui.TagOptions{{Key: "go", Label: "Go"}})
		})
	}, 360, 160)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	rect, ok := tt.Find("Go")
	if !ok {
		t.Fatal("tag missing")
	}
	x, y := int(rect.X)+5, int(rect.Y+rect.H/2)
	if err := tt.Click("Go"); err != nil {
		t.Fatal(err)
	}
	wantPixelWith(t, tt.Image(), x, y,
		colors.AccentSoftHover.Over(colors.Background), "reduced-motion selected face")
}

// TestCloseButtonHoverSettles checks the close button's hover fade
// schedules frames only until its face has settled.
func TestCloseButtonHoverSettles(t *testing.T) {
	calls := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Padding(24).Children(func() {
			mygoui.CloseButton(c, mygoui.CloseButtonOptions{Label: "Dismiss", OnClick: func() { calls++ }})
		})
	}, 200, 120)
	tt.Frame()
	rect, ok := tt.Find("Dismiss")
	if !ok {
		t.Fatal("close button missing")
	}
	tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
	if !tt.FramePending() {
		t.Fatal("a close-button press did not animate")
	}
	tt.Advance(400 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("the close-button motion did not settle")
	}
	if got := findRect(t, tt, "Dismiss"); got != rect || calls != 0 {
		t.Fatalf("pressed close button bounds=%+v calls=%d, want %+v and 0", got, calls, rect)
	}
	tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Advance(400 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("the close-button release did not settle")
	}
	if calls != 1 {
		t.Fatalf("close-button release activations = %d, want 1", calls)
	}
}

// iconInk returns the darkest channel sum in the trailing icon's box, the
// core of its glyph ink.
func iconInk(img *image.RGBA, r ui.Rect) int {
	min := 255 * 3
	for y := int(r.Y) + 2; y < int(r.Y+r.H)-2; y++ {
		for x := int(r.X+r.W) - 12; x < int(r.X+r.W); x++ {
			if x < 0 || y < 0 {
				continue
			}
			px := img.RGBAAt(x, y)
			sum := int(px.R) + int(px.G) + int(px.B)
			if sum < min {
				min = sum
			}
		}
	}
	return min
}

// TestLinkIconFadesOnHover checks the trailing icon fades from 60% to full
// opacity while the pointer is over the link.
func TestLinkIconFadesOnHover(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.Link(c, mygoui.LinkOptions{
				Label: "Docs", Icon: lucide.ArrowRight(), OnClick: func() {},
			})
		})
	}, 240, 120)
	tt.Frame()
	rect, ok := tt.Find("Docs")
	if !ok {
		t.Fatal("link missing")
	}
	rest := iconInk(tt.Image(), rect)
	tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
	start := iconInk(tt.Image(), rect)
	if d := rest - start; d < -6 || d > 6 {
		t.Fatalf("icon ink changed before the fade ran: rest %d, first frame %d", rest, start)
	}
	tt.Advance(200 * time.Millisecond)
	hovered := iconInk(tt.Image(), rect)
	if hovered >= rest-12 {
		t.Fatalf("icon did not reach full opacity: rest %d, hovered %d", rest, hovered)
	}
}
