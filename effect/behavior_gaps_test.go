package effect_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/effect"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// A shimmering region is busy, under reduced motion too, where it rests.
func TestShimmerMarksItsRegionBusy(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			region := ui.Column(c).Size(120, 40).Role(ui.RoleGroup).Label("Messages")
			effect.Shimmer(c, region, effect.ShimmerOptions{})
		}, 200, 80)
		tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: reduce})
		tt.Frame()
		if a, ok := tt.AccessByRole(ui.RoleGroup, "Messages"); !ok || !a.Busy {
			t.Fatalf("reduce motion %v: region = %+v, %v, want busy", reduce, a, ok)
		}
	}
}

// The fallback keeps content inside the masked element, as the shader's
// layer does: nothing overflowing it shows unshaded past an edge.
func TestMaskFallbackKeepsOverflowInside(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Box(c).Fill().Background(behind).Children(func() {
			e := ui.Box(c).Absolute().Left(0).Top(0).Size(40, 60)
			e.Children(func() {
				ui.Box(c).Size(40, 100).Shrink(0).Background(content)
			})
			effect.Mask(c, e, effect.MaskOptions{Edges: effect.EdgeBottom, Fallback: true})
		})
	}, 40, 100)
	tt.Frame()
	if got := tt.Image().RGBAAt(20, 80); got != opaque(behind) {
		t.Fatalf("below the element = %v, want the background %v", got, opaque(behind))
	}
}
