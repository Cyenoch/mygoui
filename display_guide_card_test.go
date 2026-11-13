package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestGuideCardDismissCallsTheCaller(t *testing.T) {
	dismissed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.GuideCard(c, mygoui.GuideCardOptions{
			Title: "Pin your tabs", Description: "Pinned tabs stay open.",
			Cover:     redImage(t),
			OnDismiss: func() { dismissed++ },
		})
	}, 320, 480)
	if !tt.HasText("Pin your tabs") || !tt.HasText("Pinned tabs stay open.") {
		t.Fatalf("card shows %v", tt.Texts())
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if dismissed != 1 {
		t.Fatalf("dismiss fired %d times", dismissed)
	}
}
