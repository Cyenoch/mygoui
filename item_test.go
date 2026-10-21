package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	myui "github.com/egoist/mygo/ui"
)

func TestItemActionsDoNotActivateParent(t *testing.T) {
	opened, removed := 0, 0
	tt := myui.NewTester(func(c *myui.Context) {
		mygoui.Item(c, mygoui.ItemOptions{
			Label: "Report", Description: "Scheduled report", Variant: mygoui.ItemOutline,
			OnClick: func() { opened++ },
			Actions: func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Remove", OnClick: func() { removed++ }})
			},
		})
	}, 420, 160)
	if err := tt.Click("Remove"); err != nil {
		t.Fatal(err)
	}
	if removed != 1 || opened != 0 {
		t.Fatalf("nested action removed=%d opened=%d", removed, opened)
	}
	if err := tt.Click("Report"); err != nil {
		t.Fatal(err)
	}
	if opened != 1 || removed != 1 {
		t.Fatalf("item action removed=%d opened=%d", removed, opened)
	}
}
