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

// TestButtonPressScaleAnimatesAndReturns checks HeroUI's pressed scale: the
// button shrinks over the transition, settles at 0.97 and eases back to its
// resting size. A full-width button gives the left edge room to move.
func TestButtonPressScaleAnimatesAndReturns(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Scale", FullWidth: true, OnClick: func() {}})
		})
	}, 480, 160)
	tt.Frame()
	rect, ok := tt.Find("Scale")
	if !ok {
		t.Fatal("button missing")
	}
	y := int(rect.Y + rect.H/2)
	left := func(dx int) int { return int(rect.X) + dx }

	wantPixel(t, tt.Image(), left(1), y, colors.Accent, "resting face")
	if tt.FramePending() {
		t.Fatal("a resting button asks for frames")
	}

	tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
	if !tt.FramePending() {
		t.Fatal("a press did not animate the scale")
	}
	// 50 ms into a 250 ms EaseOut: 3.1 DIP of a 6.5 DIP inset, so x+4 is
	// still on the face while the pressed size is not reached yet.
	tt.Advance(60 * time.Millisecond)
	wantPixel(t, tt.Image(), left(4), y, colors.AccentHover, "half-pressed face")
	tt.Advance(400 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("the press scale did not settle")
	}
	wantPixel(t, tt.Image(), left(4), y, colors.Background, "settled pressed inset")
	wantPixel(t, tt.Image(), left(1), y, colors.Background, "pressed left edge")

	tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	if !tt.FramePending() {
		t.Fatal("the release did not animate back")
	}
	tt.Advance(400 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("the release scale did not settle")
	}
	wantPixel(t, tt.Image(), left(1), y, colors.AccentHover, "face after release")
	wantPixel(t, tt.Image(), left(4), y, colors.AccentHover, "restored left edge")
}

// TestButtonPressScaleHonorsReducedMotion checks the pressed size snaps when
// the desktop asks for less motion, without animation frames.
func TestButtonPressScaleHonorsReducedMotion(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Scale", FullWidth: true})
		})
	}, 480, 160)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	rect, ok := tt.Find("Scale")
	if !ok {
		t.Fatal("button missing")
	}
	y := int(rect.Y + rect.H/2)
	left := func(dx int) int { return int(rect.X) + dx }

	tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
	if tt.FramePending() {
		t.Fatal("reduced motion still animates the press scale")
	}
	wantPixel(t, tt.Image(), left(4), y, colors.Background, "reduced-motion pressed inset")
	tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	if tt.FramePending() {
		t.Fatal("reduced motion still animates the release")
	}
	wantPixel(t, tt.Image(), left(1), y, colors.AccentHover, "reduced-motion restored face")
}

// TestButtonPendingKeepsFaceAndBlocksActivation checks HeroUI's pending
// contract: full face, no hover or press paint, no scale, no activation.
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
	wantPixel(t, tt.Image(), left(1), y, colors.Accent, "pending hover face")
	tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Advance(300 * time.Millisecond)
	wantPixel(t, tt.Image(), left(1), y, colors.Accent, "pending pressed face")
	wantPixel(t, tt.Image(), left(4), y, colors.Accent, "pending does not scale")
	tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	if calls != 0 {
		t.Fatalf("pending button activated %d times", calls)
	}
}

// TestDisabledTogglePaintsRestFace checks a disabled control reads its rest
// face while hovered: the disabled flag is set before the face samples the
// pointer, matching HeroUI's pointer-events:none.
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
	x, y := int(rect.X)+8, int(rect.Y+rect.H/2)
	wantPixelWith(t, tt.Image(), x, y,
		colors.Default.Alpha(0.5).Over(colors.Background), "disabled rest face")
	tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Frame()
	wantPixelWith(t, tt.Image(), x, y,
		colors.Default.Alpha(0.5).Over(colors.Background), "disabled hover face")
}

// TestButtonGroupSeparatorsKeepLayoutAndPaintSeams checks the shared-edge
// divider takes no layout width and paints HeroUI's 15% current-color line,
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
	wantPixel(t, tt.Image(), int(second.X)-12, y, colors.Default, "plain face")
	found := false
	for x := int(second.X) - 2; x <= int(second.X)+1; x++ {
		r, g, b := pixelAt(tt.Image(), x, y)
		if absDiff(r, colors.Default.R) > 8 || absDiff(g, colors.Default.G) > 8 ||
			absDiff(b, colors.Default.B) > 8 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("no seam near x=%v on the shared edge", second.X)
	}

	tt.Press(second.X+second.W/2, second.Y+second.H/2)
	if tt.FramePending() {
		t.Fatal("a joined button animated its press scale")
	}
	wantPixel(t, tt.Image(), int(second.X)+1, y, colors.DefaultHover, "joined pressed face")
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
	wantPixel(t, tt.Image(), int(first.X)+1, int(first.Y)+1, colors.Default, "RTL square end corner")
	wantPixel(t, tt.Image(), int(first.X+first.W)-2, int(first.Y)+1, colors.Background, "RTL rounded start corner")
	if gap := first.X - (second.X + second.W); gap < -0.5 || gap > 0.5 {
		t.Fatalf("RTL separator took layout width: gap %v", gap)
	}
	seam := false
	x, y := int(second.X+second.W), int(second.Y+second.H/2)
	for sx := x - 1; sx <= x+2; sx++ {
		r, g, b := pixelAt(tt.Image(), sx, y)
		if absDiff(r, colors.Default.R) > 8 || absDiff(g, colors.Default.G) > 8 ||
			absDiff(b, colors.Default.B) > 8 {
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
	wantPixel(t, tt.Image(), x, y, colors.Default, "tag rest face")

	if err := tt.Click("Go"); err != nil {
		t.Fatal(err)
	}
	// The selection transition starts from the face on screen.
	wantPixel(t, tt.Image(), x, y, colors.Default, "tag face at transition start")
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
// HeroUI's shared-edge divider without taking layout width.
func TestToggleButtonGroupSeparatorsPaintSeams(t *testing.T) {
	var colors theme.Palette
	off := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{Separators: true},
				mygoui.ToggleButtonOptions{Label: "Left", Selected: &off},
				mygoui.ToggleButtonOptions{Label: "Right", Selected: &off},
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
	for x := int(right.X) - 2; x <= int(right.X)+1; x++ {
		r, g, b := pixelAt(tt.Image(), x, y)
		if absDiff(r, colors.Default.R) > 8 || absDiff(g, colors.Default.G) > 8 ||
			absDiff(b, colors.Default.B) > 8 {
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

// TestCloseButtonFadesAndScales checks the close button's press scale runs
// and settles independently of its hover fade.
func TestCloseButtonFadesAndScales(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Padding(24).Children(func() {
			mygoui.CloseButton(c, mygoui.CloseButtonOptions{Label: "Dismiss", OnClick: func() {}})
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
	tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Advance(400 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("the close-button release did not settle")
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

// TestGroupedButtonsPaintThroughExitTransition checks a joined group with
// seams inside an exit copy paints without expired elements.
func TestGroupedButtonsPaintThroughExitTransition(t *testing.T) {
	show := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			ui.Button(c, "Remove").OnClick(func() { show = false })
			if show {
				ui.Row(c).Key("leaving").Transition(ui.ElementTransition{
					Duration: 250 * time.Millisecond,
					Exit:     &ui.Motion{Opacity: 0},
				}).Children(func() {
					mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Separators: true},
						mygoui.ButtonOptions{Label: "Group A", OnClick: func() {}},
						mygoui.ButtonOptions{Label: "Group B", OnClick: func() {}},
					)
				})
			}
		})
	}, 640, 480)
	tt.Frame()
	if err := tt.Click("Remove"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if tt.HasText("Group A") {
		t.Fatal("leaving group is still built")
	}
}
