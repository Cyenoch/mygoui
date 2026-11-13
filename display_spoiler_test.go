package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func spoilerLines(c *ui.Context, n int) {
	for i := range n {
		ui.Box(c.Key(i)).Height(20).FillWidth()
	}
	ui.Text(c, "End")
}

// TestSpoilerClampsAndTogglesTheCallerFlag pins the disclosure: tall
// content is cut at MaxHeight with a control under it, and the pointer and
// the keyboard toggle the caller's flag, which shows the rest.
func TestSpoilerClampsAndTogglesTheCallerFlag(t *testing.T) {
	expanded := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Spoiler(c, mygoui.SpoilerOptions{Expanded: &expanded, MaxHeight: 60}, func() { spoilerLines(c, 8) })
	}, 320, 400)
	more, less := locale.M().ShowMore, locale.M().ShowLess
	if control := findRect(t, tt, more); control.Y < 60 || control.Y > 72 {
		t.Fatalf("the control is at %+v, want just under the 60 DIP clamp", control)
	}
	if a, ok := tt.AccessByRole(ui.RoleButton, more); !ok || a.Expanded {
		t.Fatalf("the collapsed control reads as %+v (found=%v)", a, ok)
	}
	if err := tt.Click(more); err != nil {
		t.Fatal(err)
	}
	if !expanded {
		t.Fatal("the control did not set the caller's flag")
	}
	tt.Advance(300 * time.Millisecond)
	control, end := findRect(t, tt, less), findRect(t, tt, "End")
	if control.Y < end.Y+end.H {
		t.Fatalf("expanded content is still cut: control=%+v end=%+v", control, end)
	}
	if a, _ := tt.AccessByRole(ui.RoleButton, less); !a.Expanded {
		t.Fatal("the expanded control does not read as expanded")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if expanded {
		t.Fatal("Enter on the focused control did not collapse the content")
	}
	tt.Advance(300 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("a settled spoiler keeps asking for frames")
	}
}

// TestSpoilerHidesItsControlWhenContentFits pins the fit: short content has
// no control, and one appears once the content outgrows the clamp.
func TestSpoilerHidesItsControlWhenContentFits(t *testing.T) {
	lines := 1
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Spoiler(c, mygoui.SpoilerOptions{MaxHeight: 60}, func() { spoilerLines(c, lines) })
	}, 320, 400)
	if tt.HasText(locale.M().ShowMore) {
		t.Fatal("content that fits shows a control")
	}
	lines = 8
	tt.Frame()
	if !tt.HasText(locale.M().ShowMore) {
		t.Fatal("content that outgrew the clamp shows no control")
	}
	// Without a caller flag the element keeps its own.
	if err := tt.Click(locale.M().ShowMore); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(locale.M().ShowLess) {
		t.Fatal("the control did not expand the content")
	}
}

// Focus that tabs into content below the cut opens the spoiler, so the
// keyboard never rests on a control that is not shown.
func TestSpoilerOpensWhenFocusReachesHiddenContent(t *testing.T) {
	expanded := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Spoiler(c, mygoui.SpoilerOptions{Expanded: &expanded, MaxHeight: 40}, func() {
			spoilerLines(c, 6)
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Hidden action"})
		})
	}, 320, 400)
	tt.Key(0, ui.KeyTab)
	if !expanded {
		t.Fatal("focus went under the cut and the spoiler stayed closed")
	}
}
