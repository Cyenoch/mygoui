package mygoui_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestEmptyStateDefaultsToTheLocalePlaceholder pins the default copy: an
// empty state with nothing to say still names what is missing, while caller
// content replaces it.
func TestEmptyStateDefaultsToTheLocalePlaceholder(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.EmptyState(c, mygoui.EmptyStateOptions{})
	}, 240, 80)
	placeholder := locale.M().Empty
	if placeholder == "" {
		t.Fatal("the locale publishes no empty placeholder")
	}
	if !tt.HasText(placeholder) {
		t.Fatalf("empty state shows %v, want %q", tt.Texts(), placeholder)
	}

	tt = ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.EmptyState(c, mygoui.EmptyStateOptions{
			Title: "No saved searches", Description: "Save one to see it here.",
		})
	}, 240, 80)
	if tt.HasText(placeholder) {
		t.Fatalf("caller copy kept the locale placeholder: %v", tt.Texts())
	}
	if !tt.HasText("No saved searches") {
		t.Fatalf("empty state shows %v, want the caller title", tt.Texts())
	}
}

func TestEmptyStateCentersRecoveryButton(t *testing.T) {
	for _, width := range []int{280, 640} {
		for _, scale := range []float32{1, 1.3} {
			t.Run(fmt.Sprintf("%d/%.1f", width, scale), func(t *testing.T) {
				clicks := 0
				tt := ui.NewTester(func(c *ui.Context) {
					theme.ApplyMode(c, false)
					mygoui.EmptyState(c, mygoui.EmptyStateOptions{
						Title: "No saved searches", Icon: lucide.Search(),
						Description: "Searches you save appear here.",
						Action: func() {
							mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Create search", Icon: lucide.Plus(),
								Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
								OnClick: func() { clicks++ },
							})
						},
					}).Label("Empty result")
				}, width, 500)
				tt.SetPreferences(ui.Preferences{TextScale: scale, ReduceMotion: true})
				tt.Frame()
				root, ok := tt.Find("Empty result")
				if !ok {
					t.Fatal("empty result is missing")
				}
				for _, label := range []string{"No saved searches", "Create search"} {
					r, ok := tt.Find(label)
					if !ok || math.Abs(float64(r.X+r.W/2-root.X-root.W/2)) > 1 {
						t.Fatalf("%q is not centered: %+v in %+v", label, r, root)
					}
				}
				if err := tt.Click("Create search"); err != nil {
					t.Fatal(err)
				}
				if clicks != 1 {
					t.Fatalf("recovery action fired %d times", clicks)
				}
			})
		}
	}
}
