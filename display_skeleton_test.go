package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestSkeletonMotionIsOptIn pins the frame scheduling of each placeholder
// mode, standalone and grouped: only a chosen animation keeps asking for
// paint frames, and a group drives all of its placeholders from one
// animation.
func TestSkeletonMotionIsOptIn(t *testing.T) {
	cases := []struct {
		name      string
		group     bool
		anim      mygoui.SkeletonAnimation
		wantFrame bool
	}{
		{"single none", false, mygoui.SkeletonNone, false},
		{"single shimmer", false, mygoui.SkeletonShimmer, true},
		{"single pulse", false, mygoui.SkeletonPulse, true},
		{"group none", true, mygoui.SkeletonNone, false},
		{"group shimmer", true, mygoui.SkeletonShimmer, true},
		{"group pulse", true, mygoui.SkeletonPulse, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				if !tc.group {
					mygoui.Skeleton(c, mygoui.SkeletonOptions{
						Width: 160, Height: 12, Animation: tc.anim,
					})
					return
				}
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
			if got := tt.FramePending(); got != tc.wantFrame {
				t.Fatalf("%s asks for frames = %v, want %v", tc.name, got, tc.wantFrame)
			}
			if !tc.wantFrame {
				return
			}
			tt.Advance(100 * time.Millisecond)
			if !tt.FramePending() {
				t.Fatalf("%s stopped on its own", tc.name)
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
