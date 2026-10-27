package date_test

import (
	"testing"
	"time"

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

func TestDatePickerTriggerHoverAndFocus(t *testing.T) {
	var colors theme.Palette
	open := false
	calendar := &date.CalendarState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(24).Background(colors.Background).Children(func() {
			date.DatePicker(c, date.DatePickerOptions{
				Label: "Release", OpenLabel: "Open release calendar", Open: &open,
				Today: march, DayLabel: label, State: calendar,
			})
		})
	}, 420, 520)
	trigger, ok := tt.Find("Open release calendar")
	if !ok {
		t.Fatal("calendar trigger missing")
	}
	x, y := int(trigger.X+4), int(trigger.Y+4)
	wantDateHoverPixel(t, tt.Image(), x, y, colors.Background, "idle trigger")
	tt.Move(trigger.X+trigger.W/2, trigger.Y+trigger.H/2)
	wantDateHoverPixel(t, tt.Image(), x, y, colors.Accent.Over(colors.Background), "hovered trigger")
	if err := tt.Click("Open release calendar"); err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("calendar did not open")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(250 * time.Millisecond)
	if open || !tt.Focused("Open release calendar") {
		t.Fatal("Escape did not close the popup and restore trigger focus")
	}
	if err := tt.Click("Open release calendar"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.HasText("March 2026") {
		t.Fatal("calendar did not reopen after its exit transition")
	}
}

func TestDateRangePickerCompletesAndRestoresFocus(t *testing.T) {
	var value date.DateRange
	var changes int
	open := false
	state := &date.RangeCalendarState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Padding(20).Children(func() {
			date.DateRangePicker(c, date.DateRangePickerOptions{
				Label: "Dates", Open: &open, Value: &value, State: state,
				Today: march, DayLabel: label, OnChange: func(date.DateRange) { changes++ },
			})
		})
	}, 420, 560)
	if err := tt.Click("Open calendar"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("2026-03-18"); err != nil {
		t.Fatal(err)
	}
	if !open || !value.Empty() || changes != 0 {
		t.Fatalf("first endpoint committed prematurely: open=%v value=%+v changes=%d", open, value, changes)
	}
	if err := tt.Click("2026-03-12"); err != nil {
		t.Fatal(err)
	}
	if open || value.Start.Day() != 12 || value.End.Day() != 18 || changes != 1 {
		t.Fatalf("completed range: open=%v value=%+v changes=%d", open, value, changes)
	}
	tt.Advance(250 * time.Millisecond)
	if !tt.Focused("Open calendar") {
		t.Fatal("range popup did not return focus to its trigger")
	}
}
