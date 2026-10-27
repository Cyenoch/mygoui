package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func toggleGroupColorDistance(a, b uint8) int {
	d := int(a) - int(b)
	if d < 0 {
		return -d
	}
	return d
}

func TestCOSSToggleGroupGeometryAndPaint(t *testing.T) {
	for _, tc := range []struct {
		name        string
		orientation mygoui.Orientation
		variants    [3]mygoui.ToggleButtonVariant
		detached    bool
		fullWidth   bool
		gap         float32
		joined      bool
	}{
		{name: "default horizontal", gap: 2},
		{name: "default vertical", orientation: mygoui.Vertical, gap: 2},
		{name: "full width", fullWidth: true, gap: 2},
		{name: "outline horizontal", variants: [3]mygoui.ToggleButtonVariant{
			mygoui.ToggleOutline, mygoui.ToggleOutline, mygoui.ToggleOutline,
		}, joined: true},
		{name: "outline vertical", orientation: mygoui.Vertical, variants: [3]mygoui.ToggleButtonVariant{
			mygoui.ToggleOutline, mygoui.ToggleOutline, mygoui.ToggleOutline,
		}, joined: true},
		{name: "detached outline", detached: true, gap: 4, variants: [3]mygoui.ToggleButtonVariant{
			mygoui.ToggleOutline, mygoui.ToggleOutline, mygoui.ToggleOutline,
		}},
		{name: "mixed variants", gap: 2, variants: [3]mygoui.ToggleButtonVariant{
			mygoui.ToggleOutline, mygoui.ToggleDefault, mygoui.ToggleOutline,
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var colors theme.Palette
			selected := [3]bool{true, true, true}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Box(c).Fill().Padding(24).Background(colors.Background).Children(func() {
					mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{
						Orientation: tc.orientation,
						Detached:    tc.detached,
						FullWidth:   tc.fullWidth,
					},
						mygoui.ToggleButtonOptions{Label: "One", IconOnly: true, Selected: &selected[0], Variant: tc.variants[0]},
						mygoui.ToggleButtonOptions{Label: "Two", IconOnly: true, Selected: &selected[1], Variant: tc.variants[1]},
						mygoui.ToggleButtonOptions{Label: "Three", IconOnly: true, Selected: &selected[2], Variant: tc.variants[2]},
					)
				})
			}, 240, 200)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			var rects [3]ui.Rect
			for i, label := range []string{"One", "Two", "Three"} {
				var ok bool
				rects[i], ok = tt.Find(label)
				if !ok {
					t.Fatalf("missing %q", label)
				}
				if i == 0 {
					continue
				}
				gap := rects[i].X - (rects[i-1].X + rects[i-1].W)
				if tc.orientation == mygoui.Vertical {
					gap = rects[i].Y - (rects[i-1].Y + rects[i-1].H)
				}
				if gap != tc.gap {
					t.Fatalf("gap between %q and previous item = %v, want %v", label, gap, tc.gap)
				}
			}
			if tc.fullWidth && (rects[0].X < 23.9 || rects[0].X > 24.1 ||
				rects[2].X+rects[2].W < 215.9 || rects[2].X+rects[2].W > 216.1) {
				t.Fatalf("full-width group does not fill its container: %+v", rects)
			}
			mid := rects[1]
			img := tt.Image()
			corner := img.RGBAAt(int(mid.X+1), int(mid.Y+1))
			background := colors.Background
			isBackground := toggleGroupColorDistance(corner.R, background.R) <= 2 &&
				toggleGroupColorDistance(corner.G, background.G) <= 2 &&
				toggleGroupColorDistance(corner.B, background.B) <= 2
			if isBackground == tc.joined {
				t.Fatalf("middle corner=%v, background=%v, joined=%v", corner, background, tc.joined)
			}
		})
	}
}

func TestCOSSToggleGroupOutlineKeyboardFocus(t *testing.T) {
	var first, second bool
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Box(c).Fill().Padding(24).Background(colors.Background).Children(func() {
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{},
				mygoui.ToggleButtonOptions{Label: "First", Selected: &first, Variant: mygoui.ToggleOutline},
				mygoui.ToggleButtonOptions{Label: "Second", Selected: &second, Variant: mygoui.ToggleOutline},
			)
		})
	}, 240, 120)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	rect, ok := tt.Find("First")
	if !ok {
		t.Fatal("missing first outline toggle")
	}
	x, y := int(rect.X+2), int(rect.Y+rect.H/2)
	before := tt.Image().RGBAAt(x, y)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("First") || tt.Image().RGBAAt(x, y) == before {
		t.Fatal("joined outline toggle did not paint keyboard focus")
	}
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyTab)
	if !first || !tt.Focused("Second") {
		t.Fatalf("outline group keyboard behavior: first=%v second focused=%v", first, tt.Focused("Second"))
	}
}

func TestCOSSToggleGroupHoverFocusAndSelection(t *testing.T) {
	var colors theme.Palette
	var first, disabled, last bool
	firstChanges, disabledChanges, lastChanges := 0, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Box(c).Fill().Padding(24).Background(colors.Background).Children(func() {
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{},
				mygoui.ToggleButtonOptions{Label: "First", Selected: &first, OnChange: func(bool) { firstChanges++ }},
				mygoui.ToggleButtonOptions{Label: "Disabled", Selected: &disabled, Disabled: true,
					OnChange: func(bool) { disabledChanges++ }},
				mygoui.ToggleButtonOptions{Label: "Last", Selected: &last, OnChange: func(bool) { lastChanges++ }},
			)
		})
	}, 360, 160)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Move(0, 0)
	firstRect, ok := tt.Find("First")
	if !ok {
		t.Fatal("missing first toggle")
	}
	unfocused := tt.Image().RGBAAt(int(firstRect.X-2), int(firstRect.Y+firstRect.H/2))
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("First") {
		t.Fatal("Tab did not focus the first toggle")
	}
	focused := tt.Image().RGBAAt(int(firstRect.X-2), int(firstRect.Y+firstRect.H/2))
	if focused == unfocused {
		t.Fatal("keyboard focus did not paint a ring")
	}
	tt.Key(0, ui.KeySpace)
	if !first || firstChanges != 1 {
		t.Fatalf("keyboard selection: selected=%v changes=%d", first, firstChanges)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Last") {
		t.Fatal("Tab did not skip the disabled toggle")
	}
	disabledRect, ok := tt.Find("Disabled")
	if !ok {
		t.Fatal("missing disabled toggle")
	}
	tt.Move(0, 0)
	beforeHover := tt.Image().RGBAAt(int(disabledRect.X+5), int(disabledRect.Y+disabledRect.H/2))
	tt.Move(disabledRect.X+disabledRect.W/2, disabledRect.Y+disabledRect.H/2)
	tt.Frame()
	afterHover := tt.Image().RGBAAt(int(disabledRect.X+5), int(disabledRect.Y+disabledRect.H/2))
	if beforeHover != afterHover {
		t.Fatal("disabled toggle painted a hover face")
	}
	if err := tt.Click("Disabled"); err != nil {
		t.Fatal(err)
	}
	if disabled || disabledChanges != 0 {
		t.Fatal("disabled toggle accepted a click")
	}
	if err := tt.Click("Last"); err != nil {
		t.Fatal(err)
	}
	if !last || lastChanges != 1 {
		t.Fatalf("pointer selection: selected=%v changes=%d", last, lastChanges)
	}
	tt.Move(firstRect.X+firstRect.W/2, firstRect.Y+firstRect.H/2)
	tt.Frame()
	img := tt.Image()
	got := img.RGBAAt(int(firstRect.X+5), int(firstRect.Y+firstRect.H/2))
	want := colors.Input.Over(colors.Background)
	if toggleGroupColorDistance(got.R, want.R) > 3 || toggleGroupColorDistance(got.G, want.G) > 3 ||
		toggleGroupColorDistance(got.B, want.B) > 3 {
		t.Fatalf("selected hover=%v, want %v", got, want)
	}
}
