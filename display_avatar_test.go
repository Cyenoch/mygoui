package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestAvatarGroupCountFollowsTheGroupSize pins the overflow count to the
// group: HeroUI draws the count as an avatar-sized circle whose type follows
// the same size role, so a small group counts at 32 DIP and text-xs, and a
// large one at 48 DIP and text-base.
func TestAvatarGroupCountFollowsTheGroupSize(t *testing.T) {
	for _, tc := range []struct {
		name   string
		size   mygoui.AvatarSize
		line   float32
		circle int
	}{
		{"small", mygoui.AvatarSmall, 16, 790},
		{"large", mygoui.AvatarLarge, 24, 1790},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var pal theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				pal = theme.Colors(c)
				mygoui.AvatarGroup(c, mygoui.AvatarGroupOptions{
					Size:  tc.size,
					Total: 4,
					Avatars: []mygoui.AvatarOptions{
						{Name: "Ada Lovelace", Soft: true, Color: mygoui.ColorAccent},
					},
				})
			}, 200, 100)
			count, ok := tt.Find("+3")
			if !ok {
				t.Fatalf("overflow count is missing: %v", tt.Texts())
			}
			if count.H != tc.line {
				t.Fatalf("count line = %v, want %v", count.H, tc.line)
			}
			// The count is the only default-filled shape here.
			if got := inkCount(tt.Image(), pal.Default); got < tc.circle*3/4 || got > tc.circle*5/4 {
				t.Fatalf("count circle = %d px, want about %d", got, tc.circle)
			}
		})
	}
}
