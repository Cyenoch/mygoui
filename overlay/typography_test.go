package overlay

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

// Compact metadata must not inherit a document's larger type or a dialog title.
// Compare actual native text geometry with the same compact roles outside it.
func TestAnchoredPreviewKeepsCompactTypographyAndUserScale(t *testing.T) {
	t.Cleanup(func() { theme.SetPreferences(theme.Preferences{}) })
	for _, size := range []float32{16, 20} {
		for _, scale := range []float32{1, 2} {
			for _, hover := range []bool{false, true} {
				t.Run(fmt.Sprintf("size-%g-scale-%g-hover-%t", size, scale, hover), func(t *testing.T) {
					theme.SetPreferences(theme.Preferences{FontSize: size})
					open := true
					state := HoverCardState{Open: true}
					tt := myui.NewTester(func(c *myui.Context) {
						theme.Apply(c)
						myui.Column(c).Padding(20).Children(func() {
							myui.Text(c, "Machine").Label("Reference metadata").FontSize(theme.Scale(c, 13)).FixedLineHeight(theme.Scale(c, 18)).NoWrap().AlignSelf(myui.Start)
						})
						myui.Column(c.Key("large-document")).Absolute().Left(30).Top(160).FontSize(theme.Scale(c, 32)).Children(func() {
							anchor := myui.ButtonBase(c).Size(80, 30).Label("Preview anchor")
							options := PopoverOptions{Title: "Preview title", Width: 320, Placement: Placement{Side: SideRight}}
							body := func() { myui.Text(c, "Machine").Label("Preview metadata").NoWrap().AlignSelf(myui.Start) }
							if hover {
								HoverCard(c, anchor, &state, HoverCardOptions{PopoverOptions: options}, body)
							} else {
								Popover(c, anchor, &open, options, body)
							}
						})
					}, 1000, 700)
					tt.SetPreferences(myui.Preferences{TextScale: scale, ReduceMotion: true})
					tt.Frame()
					reference, ok := tt.Find("Reference metadata")
					if !ok {
						t.Fatal("reference metadata missing")
					}
					preview, ok := tt.Find("Preview metadata")
					if !ok {
						t.Fatal("preview metadata missing")
					}
					if preview.H != reference.H || preview.W != reference.W {
						t.Fatalf("preview inherited document type: reference=%+v preview=%+v", reference, preview)
					}
					t.Logf("interface=%g OS=%g native detail=%gx%g", size, scale, preview.W, preview.H)
				})
			}
		}
	}
}
