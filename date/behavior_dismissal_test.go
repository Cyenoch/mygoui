package date_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/egoist/mygo/ui"
)

// TestDatePickerOutsideClickRestoresTrigger checks the second dismissal path:
// a click outside closes the calendar without choosing, reports the close and
// returns focus to the trigger. Opening puts focus in the day grid.
func TestDatePickerOutsideClickRestoresTrigger(t *testing.T) {
	var value time.Time
	open := false
	var events []bool
	tt := start(t, func(c *ui.Context) {
		date.DatePicker(c, date.DatePickerOptions{
			Label: "Due", Today: march, DayLabel: label, Value: &value, Open: &open,
			OnOpenChange: func(v bool) { events = append(events, v) },
		})
	}, 420, 600)

	if err := tt.Click("Open calendar"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.Focused("Due") {
		t.Fatal("opening did not show the calendar with focus in its grid")
	}
	tt.ClickAt(410, 590)
	tt.Advance(300 * time.Millisecond)
	if open || !value.IsZero() {
		t.Fatalf("outside click: open %v, value %v", open, value)
	}
	if len(events) != 2 || !events[0] || events[1] {
		t.Fatalf("OnOpenChange reported %v, want [true false]", events)
	}
	if !tt.Focused("Open calendar") {
		t.Fatal("focus did not return to the trigger")
	}
}
