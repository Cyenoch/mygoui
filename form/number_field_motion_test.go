package form

import (
	"image"
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// numberFieldFixture renders one bounded number field and keeps the palette
// it painted with, so tests can compare pixels without depending on the
// control's accessible text.
type numberFieldFixture struct {
	tester *ui.Tester
	pal    theme.Palette
}

func newNumberFieldFixture(t *testing.T, dark, reduceMotion, rtl bool) *numberFieldFixture {
	t.Helper()
	lo, hi := 1.0, 12.0
	state := NumberFieldState{Value: 5}
	f := &numberFieldFixture{}
	f.tester = formTextTester(t, func(c *ui.Context) {
		theme.ApplyMode(c, dark)
		f.pal = theme.Colors(c)
		NumberField(c, &state, NumberFieldOptions{
			Label: "Seats", Min: &lo, Max: &hi, Step: 1,
		})
	})
	if reduceMotion {
		f.tester.SetPreferences(ui.Preferences{ReduceMotion: true})
	}
	if rtl {
		f.tester.SetDirection(ui.RightToLeft)
	}
	return f
}

// numberFieldButton returns the committed box of a stepper.
func (f *numberFieldFixture) numberFieldButton(t *testing.T, label string) ui.Rect {
	t.Helper()
	r, ok := f.tester.Find(label)
	if !ok {
		t.Fatalf("no %q stepper in the field", label)
	}
	return r
}

// numberFieldBlend composites ink at alpha over base, the way the painter
// lays a translucent rule over the field face.
func numberFieldBlend(base, ink ui.Color, alpha float32) ui.Color {
	mix := func(b, i uint8) uint8 {
		return uint8(float32(b)*(1-alpha) + float32(i)*alpha + 0.5)
	}
	return ui.RGBA(mix(base.R, ink.R), mix(base.G, ink.G), mix(base.B, ink.B), 1)
}

// numberFieldNearColor reports whether img paints want at x, y, allowing the
// rounding of a composited translucent pixel.
func numberFieldNearColor(img *image.RGBA, x, y int, want ui.Color) bool {
	got := img.RGBAAt(x, y)
	near := func(a, b uint8) bool {
		d := int(a) - int(b)
		return d >= -5 && d <= 5
	}
	return near(got.R, want.R) && near(got.G, want.G) && near(got.B, want.B)
}

// numberFieldContrast reports how far the pixel at x, y is from the field
// face, in luminance, so an icon reads in light and dark modes alike.
func numberFieldContrast(img *image.RGBA, x, y int, face ui.Color) int {
	got := img.RGBAAt(x, y)
	lum := func(r, g, b uint8) int { return (299*int(r) + 587*int(g) + 114*int(b)) / 1000 }
	return int(math.Abs(float64(lum(got.R, got.G, got.B) - lum(face.R, face.G, face.B))))
}

// TestNumberFieldStepperChrome pins the official group geometry: 40 DIP
// transparent steppers in the same 36 DIP row as the editor, each separated
// from it by a 1 DIP field-placeholder rule, with painted minus and plus
// icons in both appearances and both writing directions.
func TestNumberFieldStepperChrome(t *testing.T) {
	for _, mode := range []struct {
		name string
		dark bool
		rtl  bool
	}{
		{"light", false, false},
		{"dark", true, false},
		{"light-rtl", false, true},
		{"dark-rtl", true, true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			f := newNumberFieldFixture(t, mode.dark, false, mode.rtl)
			dec := f.numberFieldButton(t, "Decrease")
			inc := f.numberFieldButton(t, "Increase")

			for _, tc := range []struct {
				name string
				r    ui.Rect
			}{{"decrement", dec}, {"increment", inc}} {
				if math.Abs(float64(tc.r.W)-40) > 0.01 || math.Abs(float64(tc.r.H)-36) > 0.01 {
					t.Fatalf("%s stepper is %vx%v, want 40x36", tc.name, tc.r.W, tc.r.H)
				}
			}
			if dec.Y != inc.Y {
				t.Fatalf("steppers sit on different rows: decrement y=%v, increment y=%v", dec.Y, inc.Y)
			}
			// HeroUI mirrors the group: the decrement button starts the
			// inline axis, which is the right in a right-to-left layout.
			if mode.rtl != (dec.X > inc.X) {
				t.Fatalf("decrement x=%v, increment x=%v in rtl=%v", dec.X, inc.X, mode.rtl)
			}

			img := f.tester.Image()
			y := int(dec.Y + dec.H/2)
			rule := numberFieldBlend(f.pal.FieldBackground, f.pal.FieldPlaceholder, 0.15)
			// The rule hugs the editor-facing edge, inside each button.
			decRule := int(dec.X+dec.W) - 1
			incRule := int(inc.X)
			if mode.rtl {
				decRule, incRule = int(dec.X), int(inc.X+inc.W)-1
			}
			if !numberFieldNearColor(img, decRule, y, rule) {
				t.Fatalf("decrement rule at x=%d = %v, want %v", decRule, img.RGBAAt(decRule, y), rule)
			}
			if !numberFieldNearColor(img, incRule, y, rule) {
				t.Fatalf("increment rule at x=%d = %v, want %v", incRule, img.RGBAAt(incRule, y), rule)
			}
			// The rest of each stepper face stays transparent over the field:
			// sample away from the editor-facing rule.
			decFace := int(dec.X) + 1
			incFace := int(inc.X+inc.W) - 2
			if mode.rtl {
				decFace, incFace = int(dec.X+dec.W)-2, int(dec.X)+1
			}
			if !numberFieldNearColor(img, decFace, y, f.pal.FieldBackground) {
				t.Fatalf("decrement face at x=%d = %v, want %v", decFace, img.RGBAAt(decFace, y), f.pal.FieldBackground)
			}
			if !numberFieldNearColor(img, incFace, y, f.pal.FieldBackground) {
				t.Fatalf("increment face at x=%d = %v, want %v", incFace, img.RGBAAt(incFace, y), f.pal.FieldBackground)
			}

			for _, tc := range []struct {
				name string
				r    ui.Rect
			}{{"minus", dec}, {"plus", inc}} {
				cx, cy := int(tc.r.X+tc.r.W/2), int(tc.r.Y+tc.r.H/2)
				if got := numberFieldContrast(img, cx, cy, f.pal.FieldBackground); got < 80 {
					t.Fatalf("%s icon center paints %v, only %d from the face %v",
						tc.name, img.RGBAAt(cx, cy), got, f.pal.FieldBackground)
				}
			}
		})
	}
}

// TestNumberFieldStepperPressScale pins the pressed pose: the stepper shrinks
// to 97%, settles without asking for more frames, and returns to full size.
func TestNumberFieldStepperPressScale(t *testing.T) {
	f := newNumberFieldFixture(t, false, false, false)
	inc := f.numberFieldButton(t, "Increase")
	x, y := inc.X+inc.W/2, inc.Y+inc.H/2

	f.tester.Press(x, y)
	if early := f.numberFieldButton(t, "Increase"); early.W < 39.8 {
		t.Fatalf("pressed stepper jumped to %v instead of starting the pose", early.W)
	}
	f.tester.Advance(200 * time.Millisecond)
	pressed := f.numberFieldButton(t, "Increase")
	if want := float32(40 * 0.97); math.Abs(float64(pressed.W-want)) > 0.2 {
		t.Fatalf("pressed stepper is %v wide, want %v", pressed.W, want)
	}
	if f.tester.FramePending() {
		t.Fatal("pressed stepper still asks for frames after settling")
	}

	f.tester.Release(x, y)
	f.tester.Advance(200 * time.Millisecond)
	if back := f.numberFieldButton(t, "Increase"); math.Abs(float64(back.W)-40) > 0.2 {
		t.Fatalf("released stepper is %v wide, want 40", back.W)
	}
}

// TestNumberFieldStepperPressScaleReducedMotion pins the accessibility
// override: the pressed pose applies at once and stops requesting frames.
func TestNumberFieldStepperPressScaleReducedMotion(t *testing.T) {
	f := newNumberFieldFixture(t, false, true, false)
	inc := f.numberFieldButton(t, "Increase")

	f.tester.Press(inc.X+inc.W/2, inc.Y+inc.H/2)
	pressed := f.numberFieldButton(t, "Increase")
	if want := float32(40 * 0.97); math.Abs(float64(pressed.W-want)) > 0.2 {
		t.Fatalf("pressed stepper under reduced motion is %v wide, want %v", pressed.W, want)
	}
	if f.tester.FramePending() {
		t.Fatal("reduced-motion press still asks for frames")
	}
}
