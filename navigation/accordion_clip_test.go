package navigation

import (
	"image/color"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAccordionClipsFirstAndLastHighlightToRoundedGroup(t *testing.T) {
	for _, borderless := range []bool{false, true} {
		t.Run(map[bool]string{false: "bordered", true: "borderless"}[borderless], func(t *testing.T) {
			var state AccordionState[string]
			var outside color.RGBA
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Fill().Padding(24).Background(theme.Colors(c).Background).Children(func() {
					Accordion(c, &state, AccordionOptions{Borderless: borderless}, []AccordionItem[string]{
						{Value: "first", Label: "First section", Content: func() { ui.Text(c, "First content") }},
						{Value: "middle", Label: "Middle section", Content: func() { ui.Text(c, "Middle content") }},
						{Value: "last", Label: "Last section", Content: func() { ui.Text(c, "Last content") }},
					}).FillWidth().Label("Rounded accordion")
				})
			}, 420, 300)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			outside = tt.Image().RGBAAt(1, 1)
			check := func(label string, bottom bool) {
				t.Helper()
				item, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing %s", label)
				}
				tt.Move(item.X+item.W/2, item.Y+item.H/2)
				tt.Frame()
				r, ok := tt.Find("Rounded accordion")
				if !ok {
					t.Fatal("missing group")
				}
				y := int(r.Y + 1)
				if bottom {
					y = int(r.Y + r.H - 2)
				}
				for _, x := range []int{int(r.X + 1), int(r.X + r.W - 2)} {
					if got := tt.Image().RGBAAt(x, y); got != outside {
						t.Fatalf("%s highlight leaked around %v,%v: got=%v outside=%v", label, x, y, got, outside)
					}
				}
				midY := int(item.Y + item.H/2)
				if got := tt.Image().RGBAAt(int(item.X+4), midY); got == outside {
					t.Fatalf("%s lost actual hover highlight", label)
				}
			}
			check("First section", false)
			check("Last section", true)
			if err := tt.Click("First section"); err != nil {
				t.Fatal(err)
			}
			tt.Advance(time.Second)
			check("First section", false)
			if !tt.HasText("First content") {
				t.Fatal("clip removed expanded content")
			}
		})
	}
}
