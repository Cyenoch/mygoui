package mygoui_test

import (
	"bytes"
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func badgeOffscreenImage(t *testing.T, withBadge bool) *image.RGBA {
	t.Helper()
	var scroll ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(20).Background(theme.Colors(c).Background).Children(func() {
			ui.Scroll(c).Size(140, 100).TrackScroll(&scroll).Children(func() {
				ui.Column(c).FillWidth().Children(func() {
					ui.Box(c).Height(260).Shrink(0)
					body := func() {
						ui.Box(c).Size(80, 40).Background(ui.RGB(40, 100, 180))
					}
					if withBadge {
						mygoui.BadgeAnchor(c, mygoui.BadgeAnchorOptions{
							Badge: mygoui.BadgeOptions{Label: "Floating"},
						}, body)
					} else {
						ui.Box(c).Shrink(0).AlignSelf(ui.Start).Children(body)
					}
					ui.Box(c).Height(300).Shrink(0)
				})
			})
		})
	}, 220, 180)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	scroll.Y = 120
	tt.Frame()
	tt.Advance(2 * time.Second)
	if scroll.Y != 120 {
		t.Fatalf("viewport did not scroll to the requested offset: %+v", scroll)
	}
	return tt.Image()
}

func TestBadgeAnchorClipsWhenScrolledOffscreen(t *testing.T) {
	plain, badged := badgeOffscreenImage(t, false), badgeOffscreenImage(t, true)
	if !bytes.Equal(plain.Pix, badged.Pix) {
		t.Fatal("an offscreen badge painted outside its scrolling viewport")
	}
}

func TestBadgeAnchorProtrudesAndKeepsContentClickable(t *testing.T) {
	placements := []struct {
		name      string
		placement mygoui.BadgePlacement
		right     bool
		bottom    bool
	}{
		{"top right", mygoui.BadgeTopRight, true, false},
		{"top left", mygoui.BadgeTopLeft, false, false},
		{"bottom right", mygoui.BadgeBottomRight, true, true},
		{"bottom left", mygoui.BadgeBottomLeft, false, true},
	}
	for _, placement := range placements {
		t.Run(placement.name, func(t *testing.T) {
			calls := 0
			background := ui.RGB(255, 255, 255)
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Column(c).Fill().Padding(30).Background(background).Children(func() {
					mygoui.BadgeAnchor(c, mygoui.BadgeAnchorOptions{
						Placement: placement.placement,
						Badge:     mygoui.BadgeOptions{Label: "New", Icon: lucide.CircleCheck()},
					}, func() {
						ui.ButtonBase(c).Size(80, 40).Background(ui.RGB(40, 100, 180)).
							Label("Anchor").OnClick(func() { calls++ })
					})
				})
			}, 220, 150)
			anchor, ok := tt.Find("Anchor")
			if !ok {
				t.Fatal("anchor is missing")
			}
			badge, ok := tt.Find("New")
			if !ok {
				t.Fatal("badge is missing")
			}
			if placement.right && badge.X+badge.W <= anchor.X+anchor.W ||
				!placement.right && badge.X >= anchor.X ||
				placement.bottom && badge.Y+badge.H <= anchor.Y+anchor.H ||
				!placement.bottom && badge.Y >= anchor.Y {
				t.Fatalf("badge does not protrude from anchor: anchor=%+v badge=%+v", anchor, badge)
			}
			x := max(anchor.X+4, min(badge.X+badge.W/2, anchor.X+anchor.W-4))
			y := anchor.Y - 2
			if placement.bottom {
				y = anchor.Y + anchor.H + 1
			}
			pixel := tt.Image().RGBAAt(int(x), int(y))
			if pixel.R == background.R && pixel.G == background.G && pixel.B == background.B {
				t.Fatalf("badge protrusion did not paint at (%v, %v)", x, y)
			}

			// Both the icon and the text overlap the anchor and must pass input through.
			y = badge.Y + badge.H/2
			tt.ClickAt(badge.X+10, y)
			tt.ClickAt(badge.X+badge.W-7, y)
			if calls != 2 || !tt.Focused("Anchor") {
				t.Fatalf("click through badge: calls=%d focused=%v", calls, tt.Focused("Anchor"))
			}
		})
	}
}
