package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

func TestCardActionsRemainIndependentAtNarrowLargeText(t *testing.T) {
	calls := 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		mygoui.Card(c, mygoui.CardOptions{Small: true}, func() {
			mygoui.CardHeader(c, mygoui.CardHeaderOptions{
				Title: "Scheduled reports", Description: "Choose a schedule and delivery destination.",
				Action: func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Edit", OnClick: func() { calls++ }})
				},
			})
			mygoui.CardContent(c, func() { mygoui.Typography(c, "Every Monday", mygoui.Paragraph) })
			mygoui.CardFooter(c, func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: func() { calls += 10 }})
			})
		})
	}, 320, 420)
	tt.SetPreferences(myui.Preferences{TextScale: 2})
	if err := tt.Click("Edit"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if calls != 11 {
		t.Fatalf("card action dispatch=%d", calls)
	}
	if !tt.HasText("Scheduled reports") || !tt.HasText("Every Monday") {
		t.Fatal("native card composition lost content")
	}
}
