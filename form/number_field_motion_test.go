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

// numberFieldContrast reports how far the pixel at x, y is from the field
// face, in luminance, so an icon reads in light and dark modes alike.
func numberFieldContrast(img *image.RGBA, x, y int, face ui.Color) int {
	got := img.RGBAAt(x, y)
	lum := func(r, g, b uint8) int { return (299*int(r) + 587*int(g) + 114*int(b)) / 1000 }
	return int(math.Abs(float64(lum(got.R, got.G, got.B) - lum(face.R, face.G, face.B))))
}

// TestNumberFieldStepperChrome checks the two steppers share the editor's row,
// mirror in a right-to-left layout, and paint a legible minus and plus icon in
// both appearances.
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

			if dec.Y != inc.Y {
				t.Fatalf("steppers sit on different rows: decrement y=%v, increment y=%v", dec.Y, inc.Y)
			}
			// HeroUI mirrors the group: the decrement button starts the
			// inline axis, which is the right in a right-to-left layout.
			if mode.rtl != (dec.X > inc.X) {
				t.Fatalf("decrement x=%v, increment x=%v in rtl=%v", dec.X, inc.X, mode.rtl)
			}

			img := f.tester.Image()
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

// TestNumberFieldStepperPressScale checks the pressed stepper shrinks to 97%,
// settles without asking for more frames and returns to full size; reduced
// motion applies the pose at once.
func TestNumberFieldStepperPressScale(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		name := "motion"
		if reduce {
			name = "reduced motion"
		}
		t.Run(name, func(t *testing.T) {
			f := newNumberFieldFixture(t, false, reduce, false)
			inc := f.numberFieldButton(t, "Increase")
			x, y := inc.X+inc.W/2, inc.Y+inc.H/2

			f.tester.Press(x, y)
			if !reduce {
				if early := f.numberFieldButton(t, "Increase"); early.W < 39.8 {
					t.Fatalf("pressed stepper jumped to %v instead of starting the pose", early.W)
				}
				f.tester.Advance(200 * time.Millisecond)
			}
			pressed := f.numberFieldButton(t, "Increase")
			if want := float32(40 * 0.97); math.Abs(float64(pressed.W-want)) > 0.2 {
				t.Fatalf("pressed stepper is %v wide, want %v", pressed.W, want)
			}
			if f.tester.FramePending() {
				t.Fatal("pressed stepper still asks for frames after settling")
			}

			f.tester.Release(x, y)
			if reduce {
				return
			}
			f.tester.Advance(200 * time.Millisecond)
			if back := f.numberFieldButton(t, "Increase"); math.Abs(float64(back.W)-40) > 0.2 {
				t.Fatalf("released stepper is %v wide, want 40", back.W)
			}
		})
	}
}
