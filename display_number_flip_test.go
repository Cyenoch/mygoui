package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func countNamed(tt *ui.Tester, name string) int {
	n := 0
	for _, a := range tt.AccessAll() {
		if a.Name == name {
			n++
		}
	}
	return n
}

// TestNumberFlipRestsAndReadsOnce pins the first frame and the tree: the
// value shows without moving, and assistive technology gets it as one text,
// not a node per glyph.
func TestNumberFlipRestsAndReadsOnce(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.NumberFlip(c, mygoui.NumberFlipOptions{Value: "12,480"})
	}, 240, 80)
	findRect(t, tt, "12,480")
	if tt.FramePending() {
		t.Fatal("a value at rest requests frames")
	}
	if n := countNamed(tt, "12,480"); n != 1 {
		t.Fatalf("the value is read %d times, want once", n)
	}
}

// TestNumberFlipRollsToTheNewValue pins the change: it animates through a
// state that is neither value, then rests on the new one with no frame
// pending, and reads as the new value alone.
func TestNumberFlipRollsToTheNewValue(t *testing.T) {
	value := "999"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.NumberFlip(c, mygoui.NumberFlipOptions{Value: value})
	}, 240, 80)
	value = "1,000"
	tt.Frame()
	if !tt.FramePending() {
		t.Fatal("a change requests no frames")
	}
	if _, ok := tt.Find("9"); !ok {
		t.Fatal("the old digits are gone on the first frame of the change")
	}
	tt.Advance(120 * time.Millisecond)
	if !tt.FramePending() {
		t.Fatal("the roll stopped early")
	}
	tt.Advance(time.Second)
	findRect(t, tt, "1,000")
	if tt.FramePending() {
		t.Fatal("a settled value requests frames")
	}
	if n := countNamed(tt, "1,000"); n != 1 || countNamed(tt, "999") != 0 {
		t.Fatal("assistive technology does not read the new value alone")
	}
}

// TestNumberFlipReducedMotionKeepsGlyphsInPlace pins reduced motion: digits
// fade rather than slide, so no glyph leaves its line while changing.
func TestNumberFlipReducedMotionKeepsGlyphsInPlace(t *testing.T) {
	value := "10"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.NumberFlip(c, mygoui.NumberFlipOptions{Value: value})
	}, 240, 80)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	rest := findRect(t, tt, "10")
	value = "20"
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	for _, s := range []string{"2", "1"} {
		if r, ok := tt.Find(s); ok && r.Y != rest.Y {
			t.Fatalf("glyph %q slid to y=%v from %v", s, r.Y, rest.Y)
		}
	}
	tt.Advance(time.Second)
	findRect(t, tt, "20")
}
