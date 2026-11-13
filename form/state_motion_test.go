package form

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// motionTester renders one control on the page background.
func motionTester(build func(c *ui.Context, pal theme.Palette), reduced bool) (*ui.Tester, *theme.Palette) {
	pal := new(theme.Palette)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		*pal = theme.Colors(c)
		ui.Column(c).Fill().Padding(24).Gap(12).Background(pal.Background).Children(func() { build(c, *pal) })
	}, 360, 200)
	if reduced {
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	}
	return tt, pal
}

// redSum adds the red channel of every pixel, a measure of how much ink a
// control shows against the light page.
func redSum(img *image.RGBA) int {
	n := 0
	for i := 0; i < len(img.Pix); i += 4 {
		n += int(img.Pix[i])
	}
	return n
}

// wantBetween fails unless got lies strictly between a and b.
func wantBetween(t *testing.T, what string, got, a, b int) {
	t.Helper()
	lo, hi := min(a, b), max(a, b)
	if got <= lo || got >= hi {
		t.Fatalf("%s = %d, want between %d and %d", what, got, lo, hi)
	}
}

// markInk counts pixels inside box, away from its edge, that differ from the
// face: the ink of a mark drawn on it.
func markInk(img *image.RGBA, box ui.Rect) int {
	face := int(img.RGBAAt(int(box.X+box.W/2), int(box.Y+3)).R)
	n := 0
	for y := int(box.Y) + 3; y < int(box.Y+box.H)-3; y++ {
		for x := int(box.X) + 3; x < int(box.X+box.W)-3; x++ {
			if d := int(img.RGBAAt(x, y).R) - face; d > 15 || d < -15 {
				n++
			}
		}
	}
	return n
}

// checkboxBox is the 16 DIP box of the row labelled label.
func checkboxBox(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	row, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q row", label)
	}
	return ui.Rect{X: row.X, Y: row.Y + (row.H-16)/2, W: 16, H: 16}
}

// TestCheckboxFadesAndDrawsTheMark pins the contract: checking fades the box
// to the primary face while the check draws in, unchecking fades the mark out
// without drawing it backwards, reduced motion still fades, and a box built
// checked shows at once.
func TestCheckboxFadesAndDrawsTheMark(t *testing.T) {
	for _, reduced := range []bool{false, true} {
		checked := false
		tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
			Checkbox(c, CheckboxOptions{Checked: &checked, Label: "Accept"})
		}, reduced)
		box := checkboxBox(t, tt, "Accept")
		off := redSum(tt.Image())
		if err := tt.Click("Accept"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(40 * time.Millisecond)
		mid, midMark := redSum(tt.Image()), markInk(tt.Image(), box)
		tt.Advance(300 * time.Millisecond)
		on, fullMark := redSum(tt.Image()), markInk(tt.Image(), box)
		wantBetween(t, "box face while checking", mid, off, on)
		wantBetween(t, "mark ink while checking", midMark, 0, fullMark)

		tt.Key(0, ui.KeySpace)
		tt.Advance(40 * time.Millisecond)
		wantBetween(t, "box face while unchecking", redSum(tt.Image()), off, on)
		if markInk(tt.Image(), box) == 0 {
			t.Fatalf("reduced=%v: the mark vanished at once instead of fading", reduced)
		}
		tt.Advance(300 * time.Millisecond)
		if n := markInk(tt.Image(), box); n != 0 {
			t.Fatalf("reduced=%v: %d mark pixels left after the fade", reduced, n)
		}
	}

	on := true
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		Checkbox(c, CheckboxOptions{Checked: &on, Label: "Accept"})
	}, false)
	box := checkboxBox(t, tt, "Accept")
	first := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	if redSum(tt.Image()) != first || markInk(tt.Image(), box) == 0 {
		t.Fatal("a box built checked animated in")
	}
}

// TestCheckboxDashGrowsFromItsCentre pins the mixed mark growing in.
func TestCheckboxDashGrowsFromItsCentre(t *testing.T) {
	checked, mixed := false, false
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		Checkbox(c, CheckboxOptions{Checked: &checked, Indeterminate: &mixed, Label: "All"})
	}, false)
	box := checkboxBox(t, tt, "All")
	mixed = true
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	mid := markInk(tt.Image(), box)
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "dash ink", mid, 0, markInk(tt.Image(), box))
}

// TestRadioDotScalesInAndFadesOut pins the selected dot.
func TestRadioDotScalesInAndFadesOut(t *testing.T) {
	choice := "a"
	state := &RadioGroupState[string]{Selected: &choice}
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		RadioGroup(c, RadioGroupOptions{Label: "Pick"}, func() {
			Radio(c, RadioOptions[string]{State: state, Value: "a", Label: "A"})
			Radio(c, RadioOptions[string]{State: state, Value: "b", Label: "B"})
		})
	}, false)
	rowB, _ := tt.Find("B")
	box := ui.Rect{X: rowB.X, Y: rowB.Y + (rowB.H-16)/2, W: 16, H: 16}
	if err := tt.Click("B"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(40 * time.Millisecond)
	mid := markInk(tt.Image(), box)
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "dot ink while selecting", mid, 0, markInk(tt.Image(), box))
	if err := tt.Click("A"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(40 * time.Millisecond)
	if markInk(tt.Image(), box) == 0 {
		t.Fatal("the dot vanished at once instead of fading")
	}
}

// TestChoiceCardBorderFades pins the card's border and fill following the
// state.
func TestChoiceCardBorderFades(t *testing.T) {
	checked := false
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		CheckboxCard(c, CheckboxCardOptions{Checked: &checked, Title: "Backups", FullWidth: true})
	}, false)
	off := redSum(tt.Image())
	if err := tt.Click("Backups"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(40 * time.Millisecond)
	mid := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "card face while checking", mid, off, redSum(tt.Image()))
}

// TestFieldBorderFadesWhenInvalid pins a field turning invalid fading its
// border rather than snapping.
func TestFieldBorderFadesWhenInvalid(t *testing.T) {
	invalid := false
	value := ""
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		e := ui.TextInputBase(c, &value).Width(200).Height(36).Label("Name")
		formstyle.Surface(e, c, false, invalid, false)
	}, false)
	valid := redSum(tt.Image())
	invalid = true
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	mid := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "border while turning invalid", mid, valid, redSum(tt.Image()))
}

// TestSliderFillEasesToAnApplicationValue pins a value the application sets
// easing, and a value it starts with showing at once.
func TestSliderFillEasesToAnApplicationValue(t *testing.T) {
	value := 0.0
	tt, pal := motionTester(func(c *ui.Context, _ theme.Palette) {
		Slider(c, SliderOptions{State: SliderState{Value: &value}, Label: "Volume", HideOutput: true})
	}, false)
	end := func() int {
		img, best, bestX := tt.Image(), 0, 0
		for y := 0; y < img.Bounds().Dy(); y++ {
			n, last := 0, 0
			for x := 0; x < img.Bounds().Dx(); x++ {
				if px := img.RGBAAt(x, y); px.R == pal.Primary.R && px.G == pal.Primary.G && px.B == pal.Primary.B {
					n++
					last = x
				}
			}
			if n > best {
				best, bestX = n, last
			}
		}
		return bestX
	}
	value = 100
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	mid := end()
	tt.Advance(300 * time.Millisecond)
	full := end()
	if mid <= 0 || mid >= full {
		t.Fatalf("fill end while easing = %d, settled %d", mid, full)
	}

	value = 40
	tt2, pal2 := motionTester(func(c *ui.Context, _ theme.Palette) {
		Slider(c, SliderOptions{State: SliderState{Value: &value}, Label: "Volume", HideOutput: true})
	}, false)
	_ = pal2
	first := redSum(tt2.Image())
	tt2.Advance(300 * time.Millisecond)
	if redSum(tt2.Image()) != first {
		t.Fatal("a slider built at a value animated to it")
	}
}

// TestRateStarsFillAsTheRatingChanges pins the star fill easing.
func TestRateStarsFillAsTheRatingChanges(t *testing.T) {
	value := 0.0
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		Rate(c, &value, RateOptions{AriaLabel: "Rating"})
	}, false)
	empty := redSum(tt.Image())
	value = 3
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	mid := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "stars while filling", mid, empty, redSum(tt.Image()))
}

// TestImageSelectRingFadesIn pins the chosen tile's ring and badge.
func TestImageSelectRingFadesIn(t *testing.T) {
	choice := ""
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		ImageSelect(c, ImageSelectOptions{
			AriaLabel: "Theme", Value: &choice,
			Items: []ImageSelectItem{{Value: "a", Label: "Light"}, {Value: "b", Label: "Dark"}},
		})
	}, false)
	none := redSum(tt.Image())
	if err := tt.Click("Light"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(40 * time.Millisecond)
	mid := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "tile while chosen", mid, none, redSum(tt.Image()))
}

// TestPasswordRevealCrossfades pins the eye icon swapping through a mix.
func TestPasswordRevealCrossfades(t *testing.T) {
	value := "secret"
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		PasswordField(c, &value, PasswordFieldOptions{AriaLabel: "Password", FullWidth: true})
	}, false)
	before := redSum(tt.Image())
	// The toggle is the field's trailing 32 DIP square.
	tt.ClickAt(360-24-16, 24+16)
	tt.Advance(40 * time.Millisecond)
	mid := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	after := redSum(tt.Image())
	if mid == before || mid == after {
		t.Fatalf("icon did not pass through a mix: before %d mid %d after %d", before, mid, after)
	}
}

// TestTagsInputChipFadesOnRemoval pins a removed tag leaving by fading.
func TestTagsInputChipFadesOnRemoval(t *testing.T) {
	tags := []string{"alpha", "beta"}
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		TagsInput(c, &tags, TagsInputOptions{AriaLabel: "Tags", FullWidth: true})
	}, false)
	before := redSum(tt.Image())
	if err := tt.Click("Remove alpha"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(40 * time.Millisecond)
	mid := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "chips while one fades", mid, before, redSum(tt.Image()))
	if tt.HasText("alpha") || !tt.HasText("beta") {
		t.Fatalf("after the fade: texts %q", tt.Texts())
	}
}

// TestHotkeyCapsFadeWhileRecording pins the recorded shortcut dimming
// through a fade when recording starts.
func TestHotkeyCapsFadeWhileRecording(t *testing.T) {
	key := Hotkey{Mods: ui.Ctrl, Key: ui.KeyK}
	tt, _ := motionTester(func(c *ui.Context, _ theme.Palette) {
		HotkeyInput(c, &key, HotkeyInputOptions{AriaLabel: "Shortcut", FullWidth: true})
	}, false)
	tt.ClickAt(60, 24+16)
	// The recording ring appears at once; the caps fade after it.
	tt.Advance(time.Millisecond)
	before := redSum(tt.Image())
	tt.Advance(40 * time.Millisecond)
	mid := redSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	wantBetween(t, "key caps while recording starts", mid, before, redSum(tt.Image()))
}
