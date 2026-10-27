package gallery

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDefaultCalendarIsCenteredAndIndependent(t *testing.T) {
	page, ok := componentPage("calendar")
	if !ok || len(page.Examples) == 0 || page.Examples[0].Title != "Default" {
		t.Fatal("Calendar does not start with the default example")
	}
	for _, width := range []int{360, 1180} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			s := New()
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Fill().Children(func() { page.Examples[0].Demo(c, s) })
			}, width, 700)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			st := StateOf[datesState](s, "dates")
			today := datesStartOfDay(tt.Now())
			if !st.CalDefault.Value.Equal(today) {
				t.Fatalf("default selection = %v, want today %v", st.CalDefault.Value, today)
			}
			grid, ok := tt.Find("Calendar")
			if !ok || math.Abs(float64(grid.X+grid.W/2-float32(width)/2)) > 1 {
				t.Fatalf("default calendar is not centered: %+v", grid)
			}
			if _, ok := tt.Find(today.Format("January 2006") + ", Choose year"); ok {
				t.Fatal("default calendar unexpectedly offers a year picker")
			}
			if err := tt.Click(today.Format("January 2006")); err != nil {
				t.Fatal(err)
			}
			if st.CalDefault.YearOpen {
				t.Fatal("default calendar's static heading opened a year picker")
			}

			saturday := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
			for saturday.Weekday() != time.Saturday {
				saturday = saturday.AddDate(0, 0, 1)
			}
			if saturday.Equal(today) {
				saturday = saturday.AddDate(0, 0, 7)
			}
			bounded := st.Cal.Value
			if err := tt.Click(saturday.Format("Monday, January 2, 2006")); err != nil {
				t.Fatal(err)
			}
			if !st.CalDefault.Value.Equal(saturday) || !st.Cal.Value.Equal(bounded) {
				t.Fatalf("weekend selection changed the wrong calendar: default=%v bounded=%v",
					st.CalDefault.Value, st.Cal.Value)
			}
			if !tt.Focused("Calendar") {
				t.Fatal("selecting a day did not focus the calendar grid")
			}
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyEnter)
			if want := saturday.AddDate(0, 0, 1); !st.CalDefault.Value.Equal(want) {
				t.Fatalf("keyboard selection = %v, want %v", st.CalDefault.Value, want)
			}
			for range 2 {
				if err := tt.Click("Next"); err != nil {
					t.Fatal(err)
				}
			}
			future := time.Date(today.Year(), today.Month()+2, 15, 0, 0, 0, 0, today.Location())
			if !future.After(st.Max) {
				t.Fatalf("test day %v is not beyond the bounded example's maximum %v", future, st.Max)
			}
			if err := tt.Click(future.Format("Monday, January 2, 2006")); err != nil {
				t.Fatal(err)
			}
			if !st.CalDefault.Value.Equal(future) || !st.Cal.Value.Equal(bounded) {
				t.Fatalf("unbounded selection changed the wrong calendar: default=%v bounded=%v",
					st.CalDefault.Value, st.Cal.Value)
			}
			tt.Frame()
			tt.Frame()
		})
	}
}

func TestCalendarYearPickerExample(t *testing.T) {
	page, ok := componentPage("calendar")
	if !ok {
		t.Fatal("Calendar page missing")
	}
	var demo func(*ui.Context, *State)
	for _, example := range page.Examples {
		if example.Title == "Year picker" {
			demo = example.Demo
		}
	}
	if demo == nil {
		t.Fatal("Calendar year picker example missing")
	}
	s := New()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		demo(c, s)
	}, 700, 700)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	st := StateOf[datesState](s, "dates")
	heading := st.CalYear.Value.Format("January 2006") + ", Choose year"
	if err := tt.Click(heading); err != nil {
		t.Fatal(err)
	}
	if !st.CalYear.YearOpen {
		t.Fatal("year picker did not open")
	}
	nextYear := st.CalYear.Value.Year() + 1
	if err := tt.Click(fmt.Sprint(nextYear)); err != nil {
		t.Fatal(err)
	}
	if st.CalYear.YearOpen || st.CalYear.Month.Year() != nextYear {
		t.Fatalf("year selection month=%v open=%v", st.CalYear.Month, st.CalYear.YearOpen)
	}
	if !st.CalDefault.Value.Equal(st.Today) {
		t.Fatal("year picker changed the default calendar selection")
	}
}

func TestOtherDateExamplesRetainYearPicker(t *testing.T) {
	tests := []struct {
		page, example, trigger string
		view                   func(*datesState) (time.Time, bool)
	}{
		{"rangecalendar", "Open bounds", "", func(st *datesState) (time.Time, bool) {
			month := st.RangeFree.Month
			if month.IsZero() {
				month = st.RangeFree.Value.Start
			}
			return month, st.RangeFree.YearOpen
		}},
		{"datepicker", "Unbounded picker", "Open travel calendar", func(st *datesState) (time.Time, bool) {
			return st.PickerFreeCal.Month, st.PickerFreeCal.YearOpen
		}},
		{"daterangepicker", "Unbounded range", "Open travel calendar", func(st *datesState) (time.Time, bool) {
			return st.RangePickFreeCal.Month, st.RangePickFreeCal.YearOpen
		}},
	}
	for _, tc := range tests {
		t.Run(tc.page, func(t *testing.T) {
			page, ok := componentPage(tc.page)
			if !ok {
				t.Fatal("date page missing")
			}
			var demo func(*ui.Context, *State)
			for _, example := range page.Examples {
				if example.Title == tc.example {
					demo = example.Demo
				}
			}
			if demo == nil {
				t.Fatalf("date example %q missing", tc.example)
			}
			s := New()
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				demo(c, s)
			}, 700, 700)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			if tc.trigger != "" {
				if err := tt.Click(tc.trigger); err != nil {
					t.Fatal(err)
				}
			}
			st := StateOf[datesState](s, "dates")
			month, _ := tc.view(st)
			if err := tt.Click(month.Format("January 2006") + ", Choose year"); err != nil {
				t.Fatal(err)
			}
			if _, open := tc.view(st); !open {
				t.Fatal("year picker did not open")
			}
			tt.Key(0, ui.KeyEscape)
			if _, open := tc.view(st); open {
				t.Fatal("Escape did not close the year picker")
			}
		})
	}
}

func TestBoundedCalendarRemainsInteractive(t *testing.T) {
	s := New()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		calendarBoundedExample(c, datesComponentState(c, s))
	}, 700, 700)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	st := StateOf[datesState](s, "dates")
	before := st.Cal.Value
	var weekend, weekday time.Time
	for day := st.Min; !day.After(st.Max); day = day.AddDate(0, 0, 1) {
		if day.Year() != before.Year() || day.Month() != before.Month() {
			continue
		}
		if datesWeekend(day) {
			weekend = day
		} else if !day.Equal(before) {
			weekday = day
		}
	}
	if weekend.IsZero() || weekday.IsZero() {
		t.Fatal("bounded example has no visible weekend or alternate weekday")
	}
	if err := tt.Click(weekend.Format("Monday, January 2, 2006")); err != nil {
		t.Fatal(err)
	}
	if !st.Cal.Value.Equal(before) {
		t.Fatalf("bounded calendar selected a disabled weekend: %v", st.Cal.Value)
	}
	if err := tt.Click(weekday.Format("Monday, January 2, 2006")); err != nil {
		t.Fatal(err)
	}
	if !st.Cal.Value.Equal(weekday) || !st.CalDefault.Value.Equal(st.Today) {
		t.Fatalf("bounded selection = %v, default selection = %v", st.Cal.Value, st.CalDefault.Value)
	}
}
