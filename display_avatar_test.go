package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestAvatarGroupCountFollowsTheGroupSize pins the overflow count to the
// group: Coss draws the count as an avatar-sized circle with a line box half
// its diameter, rather than a fixed typography role.
func TestAvatarGroupCountFollowsTheGroupSize(t *testing.T) {
	for _, tc := range []struct {
		name     string
		size     mygoui.AvatarSize
		diameter float32
	}{
		{"small", mygoui.AvatarSmall, 24},
		{"medium", mygoui.AvatarMedium, 32},
		{"large", mygoui.AvatarLarge, 40},
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
						{Name: "Ada Lovelace", Soft: true, Color: mygoui.ColorSuccess},
					},
				})
			}, 200, 100)
			count, ok := tt.Find("+3")
			if !ok {
				t.Fatalf("overflow count is missing: %v", tt.Texts())
			}
			if count.H != tc.diameter/2 {
				t.Fatalf("count line = %v, want %v", count.H, tc.diameter/2)
			}
			avatar := findRect(t, tt, "Ada Lovelace")
			if avatar.W != tc.diameter || avatar.H != tc.diameter {
				t.Fatalf("avatar = %+v, want %v DIP square", avatar, tc.diameter)
			}
			cx, cy := count.X+count.W/2, count.Y+count.H/2
			img := tt.Image()
			wantPixel(t, img, int(cx), int(cy-tc.diameter/2)+1,
				pal.Default.Over(pal.Background), "count top face")
			wantPixel(t, img, int(cx), int(cy+tc.diameter/2)-2,
				pal.Default.Over(pal.Background), "count bottom face")
			wantPixel(t, img, int(cx+tc.diameter/2)-2, int(cy-tc.diameter/2)+1,
				pal.Background, "count circular corner")
		})
	}
}
