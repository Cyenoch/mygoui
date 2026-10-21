package navigation

import (
	"image/color"
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCommandSearchCornerFollowsActualOuterInset(t *testing.T) {
	for _, dark := range []bool{false, true} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Fill().Padding(20).Children(func() {
				Command(c, &CommandState{}, CommandOptions{Label: "Search commands", Height: 90}, []CommandItem{{ID: "new", Label: "New project", OnSelect: func() {}}}).FillWidth().Label("Command shell")
			})
		}, 420, 200)
		tt.SetDark(dark)
		outer, _ := tt.Find("Command shell")
		inner, _ := tt.Find("Command search surface")
		inset := inner.X - outer.X
		radius := theme.Radii().XXL - inset
		// Inside the input's expected arc at 45 degrees, but outside the same-radius
		// copy that would leave a visibly widening gap against the outer corner.
		step := radius*(1-float32(math.Sqrt(0.5))) + 1
		x, y := int(inner.X+step), int(inner.Y+step)
		base := tt.Image().RGBAAt(int(outer.X+outer.W/2), int(outer.Y+2))
		if got := tt.Image().RGBAAt(x, y); distanceRGBA(got, base) < 8 {
			t.Fatalf("dark=%v nested search radius left an oversized corner gap: pixel=%v shell=%v inset=%g", dark, got, base, inset)
		}
	}
}

func distanceRGBA(a, b color.RGBA) int {
	return int(math.Abs(float64(a.R)-float64(b.R))) + int(math.Abs(float64(a.G)-float64(b.G))) + int(math.Abs(float64(a.B)-float64(b.B)))
}
