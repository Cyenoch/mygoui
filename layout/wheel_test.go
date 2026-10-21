package layout_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/theme"
)

func TestScrollAreaWheelReachesContentWithoutMovingHeader(t *testing.T) {
	var state ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			ui.Text(c, "Fixed heading").Shrink(0)
			layout.ScrollArea(c, &state, layout.ScrollAreaOptions{Label: "Documents"}, func() {
				for range 30 {
					ui.Text(c, "Document").Height(32).Shrink(0)
				}
			}).Grow(1).MinHeight(0)
		})
	}, 320, 240)
	header, _ := tt.Find("Fixed heading")
	tt.Wheel(160, 120, 0, 120)
	tt.Advance(time.Second)
	if state.Y <= 0 || state.Y > state.MaxY {
		t.Fatalf("wheel did not reach bounded native content: %+v", state)
	}
	after, _ := tt.Find("Fixed heading")
	if header != after {
		t.Fatal("wheel scrolling moved the fixed heading")
	}
}
