package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestDescriptionsPlacesItemsInRowsOfColumns pins the grid: items fill rows
// of Columns, a span takes that many columns, an item that does not fit
// starts the next row, and caller content stands in for a value.
func TestDescriptionsPlacesItemsInRowsOfColumns(t *testing.T) {
	for _, bordered := range []bool{false, true} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			mygoui.Descriptions(c, mygoui.DescriptionsOptions{
				Title: "Order", Columns: 3, Bordered: bordered,
				Items: []mygoui.DescriptionsItem{
					{Label: "Customer", Value: "Ada"},
					{Label: "Number", Value: "10482"},
					{Label: "Placed", Value: "Today"},
					{Label: "Address", Value: "London", Span: 2},
					{Label: "Notes", Value: "Fragile", Span: 2},
					{Label: "Status", Content: func() { ui.Text(c, "Paid") }},
				},
			})
		}, 720, 400)
		customer, number, placed := findRect(t, tt, "Ada"), findRect(t, tt, "10482"), findRect(t, tt, "Today")
		if customer.Y != number.Y || number.Y != placed.Y || customer.X >= number.X || number.X >= placed.X {
			t.Fatalf("bordered=%v: first row is %+v %+v %+v", bordered, customer, number, placed)
		}
		address, notes, status := findRect(t, tt, "London"), findRect(t, tt, "Fragile"), findRect(t, tt, "Paid")
		if address.Y <= customer.Y || notes.Y <= address.Y {
			t.Fatalf("bordered=%v: a span of 2 after a span of 2 shares its row: %+v %+v", bordered, address, notes)
		}
		if status.Y != notes.Y || status.X <= notes.X {
			t.Fatalf("bordered=%v: the last item is not beside the span before it: %+v %+v", bordered, notes, status)
		}
		if _, ok := tt.AccessByRole(ui.RoleGroup, "Order"); !ok {
			t.Fatalf("bordered=%v: the title does not name the group", bordered)
		}
	}
}

// TestDescriptionsBorderedLabelsShareColumns pins the table layout: labels
// of different lengths in one column start their values on one line.
func TestDescriptionsBorderedLabelsShareColumns(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Descriptions(c, mygoui.DescriptionsOptions{
			Columns: 1, Bordered: true,
			Items: []mygoui.DescriptionsItem{
				{Label: "ID", Value: "10482"},
				{Label: "Shipping address", Value: "London"},
			},
		})
	}, 480, 200)
	if short, long := findRect(t, tt, "10482"), findRect(t, tt, "London"); short.X != long.X {
		t.Fatalf("values start at %v and %v", short.X, long.X)
	}
}

// TestDescriptionsLoadingKeepsLabelsAndHidesValues pins the placeholder: the
// labels stay readable while the values the caller does not have yet are
// not painted, and the group is busy.
func TestDescriptionsLoadingKeepsLabelsAndHidesValues(t *testing.T) {
	o := mygoui.DescriptionsOptions{
		Title: "Order", Loading: true,
		Items: []mygoui.DescriptionsItem{{Label: "Customer", Value: "Ada"}, {Label: "Order", Value: "#1"}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Descriptions(c, o)
	}, 480, 200)
	if g, ok := tt.AccessByRole(ui.RoleGroup, "Order"); !ok || !g.Busy {
		t.Fatalf("a loading record is not busy: %+v", g)
	}
	findRect(t, tt, "Customer")
	if _, ok := tt.Find("Ada"); ok {
		t.Fatal("a value shows while loading")
	}
	o.Loading = false
	tt.Frame()
	findRect(t, tt, "Ada")
}

// TestDescriptionsStayReadableWhenNarrow pins a bordered record in a narrow
// pane: each value sits under its label, on one line, and an empty one
// shows a dash.
func TestDescriptionsStayReadableWhenNarrow(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Descriptions(c, mygoui.DescriptionsOptions{
			Columns: 3, Bordered: true,
			Items: []mygoui.DescriptionsItem{{Label: "Customer", Value: "Ada Lovelace"}, {Label: "Notes"}},
		})
	}, 240, 300)
	value := findRect(t, tt, "Ada Lovelace")
	if value.X > 40 || value.H > 24 {
		t.Fatalf("value %+v: want one line, starting under its label", value)
	}
	findRect(t, tt, "—")
}
