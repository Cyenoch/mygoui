package data_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestTableNativeColumnResizeKeepsMinimumAndSavedOrder(t *testing.T) {
	state := data.TableState{}
	state.List.Columns.Order = []string{"status", "name"}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		data.Table(c, &state, data.TableOptions{Rows: []data.Row{{ID: "one", Cells: []string{"One", "Paid"}}}, Columns: []data.Column{
			{ID: "name", Title: "Name", Width: 200, MinWidth: 120},
			{ID: "status", Title: "Status", Width: 120, MinWidth: 80},
		}}).Size(400, 180)
	}, 420, 220)
	name, ok := tt.Find("Name")
	if !ok {
		t.Fatal("name header missing")
	}
	status, ok := tt.Find("Status")
	if !ok || status.X >= name.X {
		t.Fatal("native persisted column order was ignored")
	}
	// The table header text has intrinsic width. Its column's right edge is
	// the saved first-column width (120), independent of text measurement.
	tt.Press(119, status.Y+status.H/2)
	tt.Move(40, status.Y+status.H/2)
	tt.Release(40, status.Y+status.H/2)
	if width := state.List.Columns.Widths["status"]; width != 80 {
		t.Fatalf("native resize ignored column minimum: %v", width)
	}
	tt.SetPreferences(ui.Preferences{TextScale: 2, ReduceMotion: true})
	tt.SetDirection(ui.RightToLeft)
	if !tt.HasText("One") || !tt.HasText("Paid") {
		t.Fatal("large RTL table lost readable cells")
	}
}
