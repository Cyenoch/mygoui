package formstyle

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// fieldFixture renders one field on the page background and reports its
// committed box, so a test can sample and click it without depending on an
// accessible name.
type fieldFixture struct {
	tester *ui.Tester
	pal    theme.Palette
	rect   ui.Rect
}

func newFieldFixture(t *testing.T, secondary, invalid bool) *fieldFixture {
	t.Helper()
	f := &fieldFixture{}
	var value string
	var handle ui.Handle
	frames := 0
	f.tester = ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		f.pal = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(32).
			Background(f.pal.Background).Children(func() {
			e := ui.TextInputBase(c, &value).MinWidth(200).MinHeight(36).PaddingX(12)
			e.Label("field")
			e = Surface(e, c, secondary, invalid, false)
			Text(e, c).Bind(&handle)
		})
		if frames > 0 {
			f.rect = handle.Bounds(c)
		}
		frames++
	}, 320, 140)
	f.tester.Frame()
	f.tester.Frame()
	if f.rect.W <= 0 || f.rect.H <= 0 {
		t.Fatalf("field box is empty (%+v)", f.rect)
	}
	return f
}

// focus puts the pointer focus on the field, as a click does.
func (f *fieldFixture) focus() {
	f.tester.ClickAt(f.rect.X+f.rect.W/2, f.rect.Y+f.rect.H/2)
	f.tester.Frame()
}

// pixelNear reports whether img paints c at x, y, allowing the antialiasing
// other paint tests allow.
func pixelNear(img *image.RGBA, x, y int, c ui.Color) bool {
	got := img.RGBAAt(x, y)
	near := func(a, b uint8) bool {
		d := int(a) - int(b)
		return d >= -3 && d <= 3
	}
	return near(got.R, c.R) && near(got.G, c.G) && near(got.B, c.B)
}

// TestSecondaryFieldKeepsItsFaceWhenFocused pins input--secondary: the fill
// stays on the Default role under the pointer and while focused, where a
// primary field moves to the field fill.
func TestSecondaryFieldKeepsItsFaceWhenFocused(t *testing.T) {
	f := newFieldFixture(t, true, false)
	x, y := int(f.rect.X)+8, int(f.rect.Y+f.rect.H/2)
	if !pixelNear(f.tester.Image(), x, y, f.pal.Default) {
		t.Fatalf("secondary field at rest = %v, want the Default face %v",
			f.tester.Image().RGBAAt(x, y), f.pal.Default)
	}
	f.focus()
	if !pixelNear(f.tester.Image(), x, y, f.pal.Default) {
		t.Fatalf("focused secondary field = %v, want the Default face %v",
			f.tester.Image().RGBAAt(x, y), f.pal.Default)
	}
}

// TestFieldRingHugsTheBox pins focus-field-ring: a focused field shows one
// two-DIP ring against its box, in Focus, or in Danger while invalid.
func TestFieldRingHugsTheBox(t *testing.T) {
	for _, tc := range []struct {
		name    string
		invalid bool
		want    func(theme.Palette) ui.Color
	}{
		{"valid", false, func(p theme.Palette) ui.Color { return p.Focus }},
		{"invalid", true, func(p theme.Palette) ui.Color { return p.Danger }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFieldFixture(t, false, tc.invalid)
			f.focus()
			img := f.tester.Image()
			want := tc.want(f.pal)
			y := int(f.rect.Y + f.rect.H/2)
			// The ring starts at the box edge and is two DIPs wide.
			ring := false
			for dx := -3; dx <= -1; dx++ {
				if pixelNear(img, int(f.rect.X)+dx, y, want) {
					ring = true
				}
			}
			if !ring {
				t.Fatalf("no %v ring just outside the field at x %v..%v",
					want, int(f.rect.X)-3, int(f.rect.X)-1)
			}
			if pixelNear(img, int(f.rect.X)-5, y, want) {
				t.Fatal("the field ring extends past its two DIPs")
			}
		})
	}
}

// TestTriggerRingWaitsForTheKeyboard pins the button-like trigger contract:
// RingOnKeyboardFocus rings the shell for Tab, not for a pointer press, as
// select and combo-box triggers follow :focus-visible.
func TestTriggerRingWaitsForTheKeyboard(t *testing.T) {
	var pal theme.Palette
	var handle ui.Handle
	var rect ui.Rect
	frames := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(32).
			Background(pal.Background).Children(func() {
			e := ui.ButtonBase(c).FillWidth().MinHeight(36).PaddingX(12).Label("trigger")
			e = SurfaceMode(e, c, RingOnKeyboardFocus, false, false, false)
			e.Bind(&handle).Children(func() { ui.Text(c, "Open") })
		})
		if frames > 0 {
			rect = handle.Bounds(c)
		}
		frames++
	}, 320, 140)
	tt.Frame()
	tt.Frame()
	if rect.W <= 0 {
		t.Fatalf("trigger box is empty (%+v)", rect)
	}
	y := int(rect.Y + rect.H/2)
	ringed := func() bool {
		img := tt.Image()
		for dx := -3; dx <= -1; dx++ {
			if pixelNear(img, int(rect.X)+dx, y, pal.Focus) {
				return true
			}
		}
		return false
	}

	tt.ClickAt(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Frame()
	tt.Move(2, 2)
	tt.Frame()
	if ringed() {
		t.Fatal("a pointer press rings the trigger; upstream rings :focus-visible only")
	}

	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if !ringed() {
		t.Fatal("the keyboard does not ring the trigger")
	}
}
