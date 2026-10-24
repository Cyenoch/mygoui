package date_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDatePickerNotifiesEscapeDismiss(t *testing.T) {
	open := false
	notified := -1
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		date.DatePicker(c, date.DatePickerOptions{
			Label: "Release", OpenLabel: "Open release calendar", Open: &open,
			Today: march, DayLabel: label, State: &date.CalendarState{},
			OnOpenChange: func(v bool) {
				if !v {
					notified = 1
				}
			},
		})
	}, 360, 520)
	if err := tt.Click("Open release calendar"); err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("picker did not open")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open {
		t.Fatal("Escape left the picker open")
	}
	if notified != 1 {
		t.Fatal("OnOpenChange(false) was not called on Escape")
	}
}
