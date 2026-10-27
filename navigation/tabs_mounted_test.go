package navigation

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestTabsKeepMountedPreservesLocalStateWithoutHiddenTargets(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Tabs(c, TabsOptions{
			Selected: &selected, KeepMounted: true,
			Items: []TabOptions{
				{Label: "First", Content: func() {
					owner := ui.Column(c)
					count := ui.Local(owner, "count", func() int { return 0 })
					owner.Children(func() {
						ui.Button(c, fmt.Sprintf("Count %d", *count)).OnClick(func() { *count++ })
					})
				}},
				{Label: "Second", Content: func() { ui.Button(c, "Second action") }},
			},
		})
	}, 400, 200)
	if err := tt.Click("Count 0"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Second"); err != nil {
		t.Fatal(err)
	}
	if _, found := tt.Find("Count 1"); found {
		t.Fatal("hidden panel still exposes an input target")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Second action") {
		t.Fatal("hidden controls intercepted keyboard traversal")
	}
	if err := tt.Click("First"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Count 1") {
		t.Fatal("native local state was discarded when its panel was hidden")
	}
	if _, found := tt.Find("Second action"); found {
		t.Fatal("inactive panel still takes layout space")
	}
}
