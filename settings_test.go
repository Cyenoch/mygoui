package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// An expander stays open when a row is inserted above it.
func TestSettingsExpanderKeepsItsStateWhenARowIsInsertedAbove(t *testing.T) {
	rows := []mygoui.SettingsRow{
		{Title: "More", Rows: []mygoui.SettingsRow{{Title: "Inner"}}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{Rows: rows})
	}, 480, 480)
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if !tt.HasText("Inner") {
		t.Fatalf("the expander did not open: %q", tt.Texts())
	}
	rows = append([]mygoui.SettingsRow{{Title: "First"}}, rows...)
	tt.Frame()
	tt.Advance(time.Second)
	if !tt.HasText("Inner") {
		t.Fatalf("the open state moved off its row: %q", tt.Texts())
	}
}
