package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestSkeletonMotionIsOptIn pins the frame scheduling of each placeholder
// mode: only a chosen animation keeps asking for paint frames.
func TestSkeletonMotionIsOptIn(t *testing.T) {
	cases := []struct {
		name string
		anim mygoui.SkeletonAnimation
		want bool
	}{
		{"none", mygoui.SkeletonNone, false},
		{"shimmer", mygoui.SkeletonShimmer, true},
		{"pulse", mygoui.SkeletonPulse, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				mygoui.Skeleton(c, mygoui.SkeletonOptions{
					Width: 160, Height: 12, Animation: tc.anim,
				})
			}, 240, 60)
			tt.Frame()
			if got := tt.FramePending(); got != tc.want {
				t.Fatalf("%s asks for frames = %v, want %v", tc.name, got, tc.want)
			}
			if !tc.want {
				return
			}
			tt.Advance(100 * time.Millisecond)
			if !tt.FramePending() {
				t.Fatalf("%s stopped on its own", tc.name)
			}
		})
	}
}

// TestSkeletonGroupMotionIsOptIn is the group counterpart: one synchronized
// animation, and none at all when the group asks for none.
func TestSkeletonGroupMotionIsOptIn(t *testing.T) {
	for _, tc := range []struct {
		name string
		anim mygoui.SkeletonAnimation
		want bool
	}{
		{"none", mygoui.SkeletonNone, false},
		{"shimmer", mygoui.SkeletonShimmer, true},
		{"pulse", mygoui.SkeletonPulse, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				mygoui.SkeletonGroup(c, mygoui.SkeletonGroupOptions{Animation: tc.anim}, func() {
					ui.Column(c).Gap(theme.Space8).Children(func() {
						for range 3 {
							mygoui.Skeleton(c, mygoui.SkeletonOptions{
								Height: 12, Animation: mygoui.SkeletonNone,
							})
						}
					})
				})
			}, 240, 100)
			tt.Frame()
			if got := tt.FramePending(); got != tc.want {
				t.Fatalf("%s group asks for frames = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// TestSkeletonStillsForHostPreferences proves reduced motion and high
// contrast settle placeholders, group or not.
func TestSkeletonStillsForHostPreferences(t *testing.T) {
	for _, tc := range []struct {
		name string
		pref ui.Preferences
	}{
		{"reduced motion", ui.Preferences{ReduceMotion: true}},
		{"high contrast", ui.Preferences{HighContrast: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 160, Height: 12})
				mygoui.SkeletonGroup(c, mygoui.SkeletonGroupOptions{
					Animation: mygoui.SkeletonPulse,
				}, func() {
					mygoui.Skeleton(c, mygoui.SkeletonOptions{
						Height: 12, Animation: mygoui.SkeletonNone,
					})
				})
			}, 240, 100)
			tt.SetPreferences(tc.pref)
			tt.Frame()
			if tt.FramePending() {
				t.Fatalf("%s still animates placeholders", tc.name)
			}
		})
	}
}
