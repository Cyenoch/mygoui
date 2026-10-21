package mygoui_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSkeletonShimmerPaintsWithoutRebuildingOrChangingBounds(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(map[bool]string{false: "light", true: "dark"}[dark], func(t *testing.T) {
			builds := 0
			tt := ui.NewTester(func(c *ui.Context) {
				builds++
				theme.Apply(c)
				ui.Column(c).Fill().Padding(20).Children(func() {
					mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 180, Height: 24}).Label("Loading title")
					mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 44, Height: 44, Round: true}).Label("Loading avatar")
				})
			}, 240, 130)
			tt.SetDark(dark)
			before := bytes.Clone(tt.Image().Pix)
			bounds, _ := tt.Find("Loading title")
			buildsBefore := builds
			tt.Advance(800 * time.Millisecond)
			if bytes.Equal(before, tt.Image().Pix) {
				t.Fatal("skeleton shimmer did not move")
			}
			if builds != buildsBefore {
				t.Fatalf("paint-only shimmer rebuilt view %d times", builds-buildsBefore)
			}
			after, _ := tt.Find("Loading title")
			if bounds != after {
				t.Fatal("shimmer changed reserved content geometry")
			}
			// The circular avatar's bounding corner remains outside the placeholder.
			r, _ := tt.Find("Loading avatar")
			if a, b := tt.Image().RGBAAt(int(r.X), int(r.Y)), tt.Image().RGBAAt(0, 0); a != b {
				t.Fatalf("shimmer leaked beyond the avatar clip: corner=%v outside=%v", a, b)
			}
		})
	}
}

func TestSkeletonShimmerStopsForAccessibilityAndRemoval(t *testing.T) {
	visible := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if visible {
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 180, Height: 24})
		}
	}, 240, 90)
	for _, pref := range []ui.Preferences{{TextScale: 1, ReduceMotion: true}, {TextScale: 1, HighContrast: true}} {
		tt.SetPreferences(pref)
		before := bytes.Clone(tt.Image().Pix)
		tt.Advance(time.Second)
		if !bytes.Equal(before, tt.Image().Pix) || tt.FramePending() {
			t.Fatal("accessibility preference retained moving skeleton")
		}
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	visible = false
	tt.Frame()
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Fatal("removed skeleton retained a paint loop")
	}
}
