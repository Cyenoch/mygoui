package form

import (
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestNumberFieldCentersRenderedValueAfterEditing(t *testing.T) {
	state := NumberFieldState{Value: 2}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(20).Children(func() {
			NumberField(c, &state, NumberFieldOptions{AriaLabel: "Quantity editor", FullWidth: true})
		})
	}, 420, 120)

	assertCentered := func() {
		t.Helper()
		editor, ok := tt.Find("Quantity editor")
		if !ok {
			t.Fatal("number editor is missing")
		}
		img := tt.Image()
		dec, decOK := tt.Find("Decrease")
		inc, incOK := tt.Find("Increase")
		if !decOK || !incOK {
			t.Fatal("number steppers are missing")
		}
		slotLeft, slotRight := dec.X+dec.W, inc.X
		left, right := int(slotRight), int(slotLeft)-1
		for y := int(editor.Y) + 3; y < int(editor.Y+editor.H)-3; y++ {
			for x := int(slotLeft) + 3; x < int(slotRight)-3; x++ {
				p := img.RGBAAt(x, y)
				if int(p.R)+int(p.G)+int(p.B) < 450 {
					left, right = min(left, x), max(right, x)
				}
			}
		}
		if right < left {
			t.Fatal("number text produced no visible pixels")
		}
		inkCenter := float64(left+right+1) / 2
		boxCenter := float64(slotLeft+slotRight) / 2
		if math.Abs(inkCenter-boxCenter) > 2.5 {
			t.Fatalf("number ink center %.1f is not centered in editor at %.1f", inkCenter, boxCenter)
		}
		if math.Abs(float64(editor.X-slotLeft)) > 0.5 || math.Abs(float64(editor.X+editor.W-slotRight)) > 0.5 {
			t.Fatalf("editor %+v overlaps its steppers at %v and %v", editor, slotLeft, slotRight)
		}
	}

	assertCentered()
	if err := tt.Click("Quantity editor"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("128")
	if state.Value != 128 || state.Text != "128" {
		t.Fatalf("editing gave value %v and text %q", state.Value, state.Text)
	}
	tt.Key(0, ui.KeyTab)
	assertCentered()
}
