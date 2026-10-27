package date_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/egoist/mygo/ui"
)

func TestCalendarsRenderOnlyWeeksContainingTheMonth(t *testing.T) {
	eastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name         string
		month        time.Time
		firstWeekday time.Weekday
		weeks        int
		first, last  string
	}{
		{"four weeks", time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
			time.Sunday, 4, "2026-02-01", "2026-02-28"},
		{"five weeks", time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
			time.Sunday, 5, "2026-09-27", "2026-10-31"},
		{"six weeks", time.Date(2026, time.May, 15, 12, 0, 0, 0, time.UTC),
			time.Sunday, 6, "2026-04-26", "2026-06-06"},
		{"four-week month with monday start", time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
			time.Monday, 5, "2026-01-26", "2026-03-01"},
		{"six weeks with monday start", time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC),
			time.Monday, 6, "2026-02-23", "2026-04-05"},
		{"leap february", time.Date(2024, time.February, 15, 12, 0, 0, 0, time.UTC),
			time.Sunday, 5, "2024-01-28", "2024-03-02"},
		{"daylight saving transition", time.Date(2026, time.March, 15, 12, 0, 0, 0, eastern),
			time.Sunday, 5, "2026-03-01", "2026-04-04"},
	}
	for _, variant := range []string{"single", "range"} {
		for _, tc := range cases {
			t.Run(variant+"/"+tc.name, func(t *testing.T) {
				single := &date.CalendarState{}
				ranges := &date.RangeCalendarState{}
				tt := start(t, func(c *ui.Context) {
					switch variant {
					case "single":
						date.Calendar(c, date.CalendarOptions{
							Label: "Days", Today: tc.month, FirstWeekday: tc.firstWeekday, DayLabel: label,
						}, single)
					case "range":
						date.RangeCalendar(c, date.RangeCalendarOptions{
							Label: "Days", Today: tc.month, FirstWeekday: tc.firstWeekday, DayLabel: label,
						}, ranges)
					}
				}, 360, 480)
				tt.SetPreferences(ui.Preferences{ReduceMotion: true})

				firstDate, err := time.Parse("2006-01-02", tc.first)
				if err != nil {
					t.Fatal(err)
				}
				for i := range tc.weeks * 7 {
					day := label(firstDate.AddDate(0, 0, i))
					if _, ok := tt.Find(day); !ok {
						t.Fatalf("visible day %s missing", day)
					}
				}
				for _, day := range []string{
					label(firstDate.AddDate(0, 0, -1)),
					label(firstDate.AddDate(0, 0, tc.weeks*7)),
				} {
					if _, ok := tt.Find(day); ok {
						t.Fatalf("day %s appears outside the visible weeks", day)
					}
				}
				grid, ok := tt.Find("Days")
				if !ok {
					t.Fatal("day grid missing")
				}
				first, ok := tt.Find(tc.first)
				if !ok {
					t.Fatalf("first day %s missing", tc.first)
				}
				last, ok := tt.Find(tc.last)
				if !ok {
					t.Fatalf("last day %s missing", tc.last)
				}
				if last.Y-first.Y != float32(tc.weeks-1)*first.H ||
					last.Y+last.H != grid.Y+grid.H {
					t.Fatalf("grid %+v, first %+v, last %+v: want %d complete rows",
						grid, first, last, tc.weeks)
				}
			})
		}
	}
}

func TestCalendarDefaultCaptionAndWeekdays(t *testing.T) {
	for _, variant := range []string{"single", "range"} {
		t.Run(variant, func(t *testing.T) {
			single := &date.CalendarState{}
			ranges := &date.RangeCalendarState{}
			var formatter func(time.Weekday) string
			month := time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC)
			tt := start(t, func(c *ui.Context) {
				switch variant {
				case "single":
					date.Calendar(c, date.CalendarOptions{
						Label: "Days", Today: month, FirstWeekday: time.Monday, WeekdayLabel: formatter,
					}, single)
				case "range":
					date.RangeCalendar(c, date.RangeCalendarOptions{
						Label: "Days", Today: month, FirstWeekday: time.Monday, WeekdayLabel: formatter,
					}, ranges)
				}
			}, 360, 480)

			assertWeekdayOrder(t, tt, []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"})
			if _, ok := tt.Find("October 2026, Choose year"); ok {
				t.Fatal("default heading offers the year picker")
			}
			if err := tt.Click("October 2026"); err != nil {
				t.Fatal(err)
			}
			if single.YearOpen || ranges.YearOpen {
				t.Fatal("static heading opened a year picker")
			}
			formatter = func(day time.Weekday) string { return day.String() }
			tt.Frame()
			assertWeekdayOrder(t, tt, []string{
				"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
			})
		})
	}
}

func assertWeekdayOrder(t *testing.T, tt *ui.Tester, names []string) {
	t.Helper()
	var previous float32
	for i, name := range names {
		bounds, ok := tt.Find(name)
		if !ok || i > 0 && bounds.X <= previous {
			t.Fatalf("weekday %q: bounds %+v, found %v after x=%v", name, bounds, ok, previous)
		}
		previous = bounds.X
	}
}

func TestYearPickerOptInAcrossDateControls(t *testing.T) {
	for _, variant := range []string{"single", "range", "date picker", "range picker"} {
		for _, hidden := range []bool{false, true} {
			name := variant
			if hidden {
				name += "/hidden override"
			}
			t.Run(name, func(t *testing.T) {
				single := &date.CalendarState{}
				ranges := &date.RangeCalendarState{}
				open := false
				month := time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC)
				tt := start(t, func(c *ui.Context) {
					switch variant {
					case "single":
						date.Calendar(c, date.CalendarOptions{
							Today: month, ShowYearPicker: true, HideYearPicker: hidden,
						}, single)
					case "range":
						date.RangeCalendar(c, date.RangeCalendarOptions{
							Today: month, ShowYearPicker: true, HideYearPicker: hidden,
						}, ranges)
					case "date picker":
						date.DatePicker(c, date.DatePickerOptions{
							Label: "Date", Today: month, ShowYearPicker: true, HideYearPicker: hidden,
							State: single, Open: &open,
						})
					case "range picker":
						date.DateRangePicker(c, date.DateRangePickerOptions{
							Label: "Dates", Today: month, ShowYearPicker: true, HideYearPicker: hidden,
							State: ranges, Open: &open,
						})
					}
				}, 420, 620)
				if variant == "date picker" || variant == "range picker" {
					if err := tt.Click("Open calendar"); err != nil {
						t.Fatal(err)
					}
				}
				heading := "October 2026, Choose year"
				if hidden {
					if _, ok := tt.Find(heading); ok {
						t.Fatal("HideYearPicker did not override ShowYearPicker")
					}
					heading = "October 2026"
				}
				if err := tt.Click(heading); err != nil {
					t.Fatal(err)
				}
				if got := single.YearOpen || ranges.YearOpen; got == hidden {
					t.Fatalf("year picker open %v with hide override %v", got, hidden)
				}
				if !hidden {
					tt.Key(0, ui.KeyEscape)
					focused := tt.Focused(heading) || tt.Focused("Thursday, October 1, 2026")
					if single.YearOpen || ranges.YearOpen || !focused {
						t.Fatalf("Escape: year picker open %v, focus returned to calendar %v",
							single.YearOpen || ranges.YearOpen, focused)
					}
					if (variant == "date picker" || variant == "range picker") && !open {
						t.Fatal("closing the year picker also closed the calendar popup")
					}
				}
			})
		}
	}
}

func TestYearPickerScrollsWithinMonthHeight(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Label: "Days", Today: time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
			Min:      time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC),
			Max:      time.Date(2100, time.December, 31, 0, 0, 0, 0, time.UTC),
			DayLabel: label, ShowYearPicker: true,
		}, st)
	}, 360, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	days, ok := tt.Find("Days")
	if !ok {
		t.Fatal("day grid missing")
	}
	firstDay, ok := tt.Find("2026-02-01")
	if !ok {
		t.Fatal("first day missing")
	}
	if err := tt.Click("February 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	year, ok := tt.Find("2026")
	if !ok || year.Y < days.Y || year.Y+year.H > days.Y+days.H {
		t.Fatalf("selected year %+v is outside the four-week viewport %+v", year, days)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Frame()
	year, ok = tt.Find("2100")
	if !ok || year.Y < days.Y || year.Y+year.H > days.Y+days.H+firstDay.H {
		t.Fatalf("last year %+v is outside the five-week viewport starting at %+v", year, days)
	}
	tt.Key(0, ui.KeyEnter)
	if st.YearOpen || st.Month.Year() != 2100 || !tt.Focused("February 2100, Choose year") {
		t.Fatalf("year selection: month %v, year picker open %v, heading focused %v",
			st.Month, st.YearOpen, tt.Focused("February 2100, Choose year"))
	}
	if _, ok := tt.Find("2100-03-06"); !ok {
		t.Fatal("February 2100 did not render its trailing week")
	}
}

func TestCalendarFocusAndSelectionAcrossDifferentWeekCounts(t *testing.T) {
	st := &date.CalendarState{}
	var changes []time.Time
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Label: "Days", Today: jan2026, DayLabel: label,
			OnChange: func(day time.Time) { changes = append(changes, day) },
		}, st)
	}, 360, 480)
	if err := tt.Click("2026-01-31"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if st.Month.Month() != time.February || st.Cursor.Day() != 1 || st.Value.Day() != 31 || !tt.Focused("Days") {
		t.Fatalf("crossed into February: month %v, cursor %v, value %v, focused %v",
			st.Month, st.Cursor, st.Value, tt.Focused("Days"))
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if st.Value.Month() != time.February || st.Value.Day() != 28 || len(changes) != 2 {
		t.Fatalf("February selection %v, callbacks %v", st.Value, changes)
	}
	for _, month := range []time.Month{time.March, time.April, time.May} {
		tt.Key(0, ui.KeyPageDown)
		if st.Month.Month() != month || st.Cursor.Day() != 28 || !tt.Focused("Days") {
			t.Fatalf("paged to %v: month %v, cursor %v, focused %v", month, st.Month, st.Cursor, tt.Focused("Days"))
		}
		if st.Value.Month() != time.February || len(changes) != 2 {
			t.Fatalf("paging changed the selection %v or callbacks %v", st.Value, changes)
		}
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if st.Value.Month() != time.May || st.Value.Day() != 31 || len(changes) != 3 || !tt.Focused("Days") {
		t.Fatalf("sixth-row selection %v, callbacks %v, focused %v", st.Value, changes, tt.Focused("Days"))
	}
}

func TestCalendarSelectsVisibleOutsideDay(t *testing.T) {
	st := &date.CalendarState{}
	var selected time.Time
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Label: "Days", Today: time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
			DayLabel: label, OnChange: func(day time.Time) { selected = day },
		}, st)
	}, 360, 480)
	if err := tt.Click("2026-09-27"); err != nil {
		t.Fatal(err)
	}
	if st.Month.Month() != time.September || st.Value.Month() != time.September ||
		st.Value.Day() != 27 || !selected.Equal(st.Value) || !tt.Focused("Days") {
		t.Fatalf("outside-day selection %v, month %v, callback %v, focused %v",
			st.Value, st.Month, selected, tt.Focused("Days"))
	}
}

func TestRangeCalendarSelectionAcrossDifferentWeekCounts(t *testing.T) {
	st := &date.RangeCalendarState{}
	var changes []date.DateRange
	tt := start(t, func(c *ui.Context) {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Label: "Days", Today: time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC), DayLabel: label,
			DisabledDates: func(day time.Time) bool { return day.Month() == time.May && day.Day() == 1 },
			OnChange:      func(value date.DateRange) { changes = append(changes, value) },
		}, st)
	}, 360, 480)
	if err := tt.Click("2026-02-28"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyPageDown)
	tt.Key(0, ui.KeyPageDown)
	if st.Month.Month() != time.May || st.Cursor.Day() != 1 || !tt.Focused("Days") {
		t.Fatalf("cross-month cursor %v, month %v, focused %v", st.Cursor, st.Month, tt.Focused("Days"))
	}
	tt.Key(0, ui.KeyEnter)
	if st.Anchor.Month() != time.February || !st.Value.End.IsZero() || len(changes) != 0 {
		t.Fatalf("disabled endpoint changed pending range: anchor %v, value %+v, callbacks %v",
			st.Anchor, st.Value, changes)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if st.Value.Start.Month() != time.February || st.Value.Start.Day() != 28 ||
		st.Value.End.Month() != time.May || st.Value.End.Day() != 2 ||
		!st.Anchor.IsZero() || len(changes) != 1 || !tt.Focused("Days") {
		t.Fatalf("completed range %+v, anchor %v, callbacks %v, focused %v",
			st.Value, st.Anchor, changes, tt.Focused("Days"))
	}
}
