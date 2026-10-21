package calendar_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/calendar"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCivilDatesRejectNormalizationAndPreserveLocation(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		valid      bool
	}{
		{"leap day", "2024-02-29", true}, {"nonleap day", "2025-02-29", false}, {"month overflow", "2026-13-01", false}, {"day overflow", "2026-04-31", false}, {"ambiguous format", "02/03/2026", false}, {"zero year", "0000-01-01", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, err := calendar.Parse(tc.text)
			if (err == nil) != tc.valid || d.Valid() != tc.valid {
				t.Fatalf("Parse(%q)=%v,%v", tc.text, d, err)
			}
		})
	}
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	state := calendar.State{Date: calendar.Date{Year: 2026, Month: time.March, Day: 7}}
	for _, day := range []int{8, 9, 10} {
		d := calendar.Date{Year: 2026, Month: time.March, Day: day}
		if !state.Select(d, calendar.Options{Location: location}) {
			t.Fatal("civil selection failed")
		}
		instant := state.Date.In(location)
		if instant.Location() != location || calendar.FromTime(instant) != d || instant.Hour() != 12 {
			t.Fatalf("DST civil date drift: %v", instant)
		}
	}
	a := calendar.Date{Year: 2026, Month: time.March, Day: 7}.In(location)
	b := calendar.Date{Year: 2026, Month: time.March, Day: 8}.In(location)
	if b.Sub(a) != 23*time.Hour {
		t.Fatalf("fixture did not cross DST: %v", b.Sub(a))
	}
}

func TestCalendarRangeOrderingAndDisabledInterior(t *testing.T) {
	day := func(n int) calendar.Date { return calendar.Date{Year: 2026, Month: time.October, Day: n} }
	state := calendar.State{}
	o := calendar.Options{Mode: calendar.RangeMode}
	state.Select(day(20), o)
	state.Select(day(12), o)
	if state.Range.Start != day(12) || state.Range.End != day(20) {
		t.Fatalf("reversed range=%v", state.Range)
	}
	state.Range = calendar.Range{}
	o.Disabled = func(d calendar.Date) bool { return d == day(16) }
	state.Select(day(12), o)
	if state.Select(day(20), o) || state.Range.End.Valid() {
		t.Fatal("range crossed unavailable interior")
	}
	if state.Select(day(16), o) {
		t.Fatal("disabled endpoint selected")
	}
	o.Mode = calendar.Multiple
	state.Select(day(20), o)
	state.Select(day(12), o)
	state.Select(day(20), o)
	if len(state.Dates) != 1 || state.Dates[0] != day(12) {
		t.Fatalf("multiple toggle=%v", state.Dates)
	}
}

func TestCalendarNativeKeyboardSkipsDisabledAndClampsMonth(t *testing.T) {
	state := calendar.State{Date: calendar.Date{Year: 2024, Month: time.January, Day: 31}}
	monday := time.Monday
	changes := 0
	o := calendar.Options{Label: "Schedule", Locale: "en", FirstWeekday: &monday, OnChange: func() { changes++ }}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); calendar.Calendar(c, &state, o) }, 380, 430)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyPageDown)
	tt.Key(0, ui.KeyEnter)
	if state.Date != (calendar.Date{Year: 2024, Month: time.February, Day: 29}) || changes != 1 {
		t.Fatalf("PageDown leap clamp=%v changes=%d", state.Date, changes)
	}
	o.Disabled = func(d calendar.Date) bool { return d.Day == 1 && d.Month == time.March }
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if state.Date != (calendar.Date{Year: 2024, Month: time.March, Day: 2}) || changes != 2 {
		t.Fatalf("disabled skip=%v changes=%d", state.Date, changes)
	}
	o.ReadOnly = true
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if state.Date.Day != 2 || changes != 2 {
		t.Fatal("read-only keyboard selected")
	}
}

func TestCalendarMonthYearNativeChoicesChangeVisibleDates(t *testing.T) {
	state := calendar.State{Date: calendar.Date{Year: 2026, Month: time.October, Day: 20}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		calendar.Calendar(c, &state, calendar.Options{Label: "Birth date", Locale: "en", MonthChoice: true, StartYear: 1900, EndYear: 2026})
	}, 400, 430)
	if err := tt.Click("Choose year"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("2000"); err != nil {
		t.Fatal(err)
	}
	if state.Month.Year != 2000 {
		t.Fatal("year choice did not change actual calendar")
	}
	if err := tt.Click("Choose month"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("February"); err != nil {
		t.Fatal(err)
	}
	if state.Month.Month != time.February {
		t.Fatal("month choice did not change actual calendar")
	}
	tt.Frame()
	if err := tt.Click("Tuesday, February 29, 2000"); err != nil {
		t.Fatal(err)
	}
	if state.Date != (calendar.Date{Year: 2000, Month: time.February, Day: 29}) {
		t.Fatalf("leap day choice=%v", state.Date)
	}
}
