package date_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/egoist/mygo/ui"
)

func TestCalendarGridClickFocusesKeyboard(t *testing.T) {
	today := time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC)
	for _, width := range []int{360, 1180} {
		for _, variant := range []string{"single", "range"} {
			t.Run(fmt.Sprintf("%s/%d", variant, width), func(t *testing.T) {
				single := &date.CalendarState{}
				ranges := &date.RangeCalendarState{}
				changes, builds := 0, 0
				disabled := func(day time.Time) bool {
					return day.Year() == 2026 && day.Month() == time.November && day.Day() == 15
				}
				tt := start(t, func(c *ui.Context) {
					builds++
					if variant == "single" {
						date.Calendar(c, date.CalendarOptions{
							Label: "Days", Today: today, DayLabel: label, DisabledDates: disabled,
							OnChange: func(time.Time) { changes++ },
						}, single)
					} else {
						date.RangeCalendar(c, date.RangeCalendarOptions{
							Label: "Days", Today: today, DayLabel: label, DisabledDates: disabled,
							OnChange: func(date.DateRange) { changes++ },
						}, ranges)
					}
				}, width, 480)
				tt.SetPreferences(ui.Preferences{ReduceMotion: true})
				if err := tt.Click("Days"); err != nil {
					t.Fatal(err)
				}
				if !tt.Focused("Days") {
					t.Fatal("clicking the day grid did not focus it")
				}
				before := single.Value
				if variant == "range" {
					before = ranges.Value.Start
				}
				if before.IsZero() {
					t.Fatal("clicking the day grid did not select a day")
				}
				tt.Key(0, ui.KeyRight)
				tt.Key(0, ui.KeyEnter)
				want := before.AddDate(0, 0, 1)
				if variant == "single" {
					if !single.Value.Equal(want) {
						t.Fatalf("keyboard selected %v, want %v", single.Value, want)
					}
				} else if !ranges.Value.End.Equal(want) {
					t.Fatalf("keyboard selected range %+v, want end %v", ranges.Value, want)
				}

				tt.SetPreferences(ui.Preferences{ReduceMotion: false})
				tt.SetPreferences(ui.Preferences{ReduceMotion: true})
				if !tt.Focused("Days") {
					t.Fatal("changing Reduce Motion lost grid focus")
				}
				tt.Key(0, ui.KeyPageDown)
				month := single.Month
				if variant == "range" {
					month = ranges.Month
				}
				if month.Month() != time.November || !tt.Focused("Days") {
					t.Fatalf("after paging: month %v, focus %v", month, tt.Focused("Days"))
				}

				oldSingle, oldRange, oldChanges := single.Value, ranges.Value, changes
				tt.Key(0, ui.KeyEnter)
				if err := tt.Click("2026-11-15"); err != nil {
					t.Fatal(err)
				}
				if single.Value != oldSingle || ranges.Value != oldRange || changes != oldChanges ||
					!tt.Focused("Days") {
					t.Fatalf("disabled date changed selection or focus: single %v, range %+v, changes %d, focus %v",
						single.Value, ranges.Value, changes, tt.Focused("Days"))
				}
				tt.Advance(time.Second)
				beforeFrame := builds
				tt.Frame()
				if builds > beforeFrame+1 || tt.FramePending() {
					t.Fatalf("idle frame: %d builds, frame pending %v", builds-beforeFrame, tt.FramePending())
				}
			})
		}
	}
}
