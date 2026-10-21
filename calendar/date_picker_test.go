package calendar_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/calendar"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDatePickerNativeEditorInvalidClearAndPopover(t *testing.T) {
	state := calendar.DatePickerState{Calendar: calendar.State{Date: calendar.Date{Year: 2026, Month: time.October, Day: 20}}}
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		calendar.DatePicker(c, &state, calendar.DatePickerOptions{Label: "Appointment", Calendar: calendar.Options{Locale: "en"}, Clearable: true, OnChange: func() { changes++ }})
	}, 400, 520)
	if err := tt.Click("Appointment"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Appointment") {
		t.Fatal("picker caption did not focus its editable input")
	}
	tt.Command("selectAll")
	tt.Type("2025-02-29")
	if state.Invalid == "" || state.Calendar.Date.Day != 20 || changes != 0 {
		t.Fatalf("invalid draft replaced valid value: %+v changes=%d", state, changes)
	}
	tt.Command("selectAll")
	tt.Type("2024-02-29")
	if state.Invalid != "" || state.Calendar.Date != (calendar.Date{Year: 2024, Month: time.February, Day: 29}) || changes != 1 {
		t.Fatalf("valid editor selection: %+v changes=%d", state, changes)
	}
	if err := tt.Click("Open calendar"); err != nil {
		t.Fatal(err)
	}
	if !state.Open {
		t.Fatal("calendar popover did not open")
	}
	if err := tt.Click("Wednesday, February 28, 2024"); err != nil {
		t.Fatal(err)
	}
	if state.Open || state.Text != "2024-02-28" || state.Calendar.Date.Day != 28 || changes != 2 {
		t.Fatalf("calendar and editor diverged: %+v changes=%d", state, changes)
	}
	if err := tt.Click("Clear date"); err != nil {
		t.Fatal(err)
	}
	if state.Calendar.Date.Valid() || state.Text != "" || state.Invalid != "" || changes != 3 {
		t.Fatalf("clear=%+v changes=%d", state, changes)
	}
}

func TestDatePickerRangePresetsAndOrdering(t *testing.T) {
	start := calendar.Date{Year: 2026, Month: time.October, Day: 12}
	end := calendar.Date{Year: 2026, Month: time.October, Day: 20}
	state := calendar.DatePickerState{Calendar: calendar.State{Range: calendar.Range{Start: start, End: end}}}
	o := calendar.DatePickerOptions{Label: "Trip", Range: true, Calendar: calendar.Options{Locale: "en"}, Presets: []calendar.Preset{{Label: "Sample week", Range: calendar.Range{Start: start, End: calendar.Date{Year: 2026, Month: time.October, Day: 18}}}}}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); calendar.DatePicker(c, &state, o) }, 480, 550)
	tt.Click("Trip end")
	tt.Command("selectAll")
	tt.Type("2026-10-01")
	if state.EndInvalid == "" || state.Calendar.Range.End != end {
		t.Fatal("reversed typed range replaced valid selection")
	}
	tt.Click("Open calendar")
	if err := tt.Click("Sample week"); err != nil {
		t.Fatal(err)
	}
	if state.Open || state.Calendar.Range.End.Day != 18 || state.EndText != "2026-10-18" || state.EndInvalid != "" {
		t.Fatalf("preset did not update both editors: %+v", state)
	}
}
