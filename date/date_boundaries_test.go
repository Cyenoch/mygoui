package date_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

var jan2026 = time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)

// TestCalendarKeyboardCrossesMonthBoundary checks that moving past the last day
// follows it into the next month, and that paging a month clamps a day that the
// new month does not have.
func TestCalendarKeyboardCrossesMonthBoundary(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Today: jan2026, DayLabel: label}, st)
	}, 340, 460)

	if err := tt.Click("2026-01-31"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if st.Month.Month() != time.February || st.Cursor.Day() != 1 {
		t.Fatalf("after Right from Jan 31: month %v cursor %v, want February 1",
			st.Month, st.Cursor)
	}

	// Paging from the 31st into a short month keeps the day the month has.
	st.Month = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	st.Cursor = time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC)
	tt.Frame()
	tt.Key(0, ui.KeyPageDown)
	if st.Month.Month() != time.February || st.Cursor.Day() != 28 {
		t.Fatalf("after PageDown: month %v cursor %v, want February 28",
			st.Month, st.Cursor)
	}
}

// TestCalendarBoundsStopMonthNavigation checks that the arrows cannot leave a
// calendar whose bounds sit inside one month.
func TestCalendarBoundsStopMonthNavigation(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Today: jan2026,
			Min:   time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			Max:   time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC),
		}, st)
	}, 340, 460)

	before := st.Month
	for _, action := range []string{"Previous", "Next"} {
		if err := tt.Click(action); err != nil {
			t.Fatalf("%s: %v", action, err)
		}
		if !st.Month.Equal(before) {
			t.Fatalf("%s changed the visible month to %v", action, st.Month)
		}
		if !tt.HasText("January 2026") {
			t.Fatalf("%s left the bounded month: %v", action, tt.Texts())
		}
	}
}

// TestCalendarYearPickerHidesMonthNavigation checks that the month arrows stop
// taking presses while the year grid covers the days, and work again after it
// closes.
func TestCalendarYearPickerHidesMonthNavigation(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Today: time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC)}, st)
	}, 340, 460)

	if err := tt.Click("March 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	if !st.YearOpen {
		t.Fatal("the year picker did not open")
	}
	next, ok := tt.Find("Next")
	if !ok {
		t.Fatal("the next-month arrow is missing")
	}
	before := st.Month
	tt.ClickAt(next.X+next.W/2, next.Y+next.H/2)
	if !st.Month.Equal(before) {
		t.Fatalf("a hidden arrow moved the month to %v", st.Month)
	}
	if !tt.HasText("March 2026") {
		t.Fatalf("a hidden arrow changed the visible month: %v", tt.Texts())
	}

	if err := tt.Click("March 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	if st.YearOpen {
		t.Fatal("the heading did not close the year picker")
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if st.Month.Month() != time.April {
		t.Fatalf("after closing the year picker, Next went to %v", st.Month)
	}
}

// TestCalendarYearPickerEscapeCloses checks the year grid's Escape path.
func TestCalendarYearPickerEscapeCloses(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Today: time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC)}, st)
	}, 340, 460)

	if err := tt.Click("March 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if st.YearOpen {
		t.Fatal("Escape left the year picker open")
	}
	if !tt.HasText("March 2026, Choose year") {
		t.Fatalf("the month heading is missing after Escape: %v", tt.Texts())
	}
}

// TestCalendarYearPickerFadesIn checks the native crossfade: the year grid is
// transparent on the frame it opens and painted once the transition settles,
// while reduced motion paints it at once.
func TestCalendarYearPickerFadesIn(t *testing.T) {
	newTester := func(reduced bool) (*ui.Tester, *theme.Palette) {
		st := &date.CalendarState{}
		pal := new(theme.Palette)
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			*pal = theme.Colors(c)
			ui.Column(c).FillWidth().FillHeight().Padding(20).Background(pal.Background).
				Children(func() {
					date.Calendar(c, date.CalendarOptions{
						Today: time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC),
					}, st)
				})
		}, 360, 480)
		if reduced {
			tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		}
		if err := tt.Click("March 2026, Choose year"); err != nil {
			t.Fatal(err)
		}
		return tt, pal
	}

	t.Run("animates", func(t *testing.T) {
		tt, pal := newTester(false)
		cell, ok := tt.Find("2026")
		if !ok {
			t.Fatal("the selected year cell is missing")
		}
		x, y := int(cell.X+cell.W*0.3), int(cell.Y+cell.H*0.5)
		before := tt.Image().RGBAAt(x, y)
		if diff(before.R, pal.Accent.R) < 8 && diff(before.G, pal.Accent.G) < 8 &&
			diff(before.B, pal.Accent.B) < 8 {
			t.Fatal("the year grid is already opaque on its first frame")
		}
		tt.Advance(300 * time.Millisecond)
		after := tt.Image().RGBAAt(x, y)
		if diff(after.R, pal.Accent.R) > 3 || diff(after.G, pal.Accent.G) > 3 ||
			diff(after.B, pal.Accent.B) > 3 {
			t.Fatalf("settled year fill = %d,%d,%d, want %d,%d,%d",
				after.R, after.G, after.B, pal.Accent.R, pal.Accent.G, pal.Accent.B)
		}
	})

	t.Run("reduced motion", func(t *testing.T) {
		tt, pal := newTester(true)
		cell, ok := tt.Find("2026")
		if !ok {
			t.Fatal("the selected year cell is missing")
		}
		x, y := int(cell.X+cell.W*0.3), int(cell.Y+cell.H*0.5)
		got := tt.Image().RGBAAt(x, y)
		if diff(got.R, pal.Accent.R) > 3 || diff(got.G, pal.Accent.G) > 3 ||
			diff(got.B, pal.Accent.B) > 3 {
			t.Fatalf("reduced-motion year fill = %d,%d,%d, want accent",
				got.R, got.G, got.B)
		}
	})
}

// TestRangeCalendarAcrossMonths checks that a range can start in one month and
// end in the next, and stays ordered.
func TestRangeCalendarAcrossMonths(t *testing.T) {
	st := &date.RangeCalendarState{}
	var got date.DateRange
	tt := start(t, func(c *ui.Context) {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Today: jan2026, DayLabel: label,
			OnChange: func(r date.DateRange) { got = r },
		}, st)
	}, 360, 500)

	if err := tt.Click("2026-01-30"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("February 2026") {
		t.Fatalf("the next month did not show: %v", tt.Texts())
	}
	if err := tt.Click("2026-02-03"); err != nil {
		t.Fatal(err)
	}
	if got.Start.Month() != time.January || got.Start.Day() != 30 ||
		got.End.Month() != time.February || got.End.Day() != 3 {
		t.Fatalf("cross-month range = %v..%v", got.Start, got.End)
	}
	if !st.Value.Contains(time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("the range does not cover the month boundary")
	}
	if !st.Anchor.IsZero() {
		t.Fatal("the pending endpoint was not cleared")
	}
}

// TestRangeCalendarRejectsOutOfBoundsEndpoint checks that an out-of-bounds day
// cannot start a range.
func TestRangeCalendarRejectsOutOfBoundsEndpoint(t *testing.T) {
	st := &date.RangeCalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Today:    time.Date(2026, time.February, 10, 12, 0, 0, 0, time.UTC),
			Min:      time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
			Max:      time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC),
			DayLabel: label,
		}, st)
	}, 360, 500)

	if _, ok := tt.Find("2026-03-01"); !ok {
		t.Fatal("an out-of-bounds trailing day is not rendered")
	}
	if err := tt.Click("2026-03-01"); err != nil {
		t.Fatal(err)
	}
	if !st.Anchor.IsZero() || !st.Value.Empty() {
		t.Fatalf("an out-of-bounds day started a range: anchor %v value %+v",
			st.Anchor, st.Value)
	}
	if err := tt.Click("2026-02-28"); err != nil {
		t.Fatal(err)
	}
	if st.Value.Start.Day() != 28 {
		t.Fatalf("the in-bounds endpoint = %v", st.Value.Start)
	}
}

// TestTimeFieldNoonAndMidnight checks the 12-hour boundaries that the plain
// HH:MM parser does not share.
func TestTimeFieldNoonAndMidnight(t *testing.T) {
	for _, tc := range []struct {
		text         string
		hour, minute int
	}{
		{"12:00 AM", 0, 0},
		{"12:30 PM", 12, 30},
		{"11:59 PM", 23, 59},
	} {
		t.Run(tc.text, func(t *testing.T) {
			st := &date.TimeFieldState{}
			tt := start(t, func(c *ui.Context) {
				date.TimeField(c, date.TimeFieldOptions{Label: "Alarm", Hour12: true}, st)
			}, 320, 220)
			if err := tt.Click("Alarm"); err != nil {
				t.Fatal(err)
			}
			tt.Type(tc.text)
			if st.Value.Hour() != tc.hour || st.Value.Minute() != tc.minute {
				t.Fatalf("%s committed %v, want %02d:%02d",
					tc.text, st.Value, tc.hour, tc.minute)
			}
		})
	}
}

// TestDateFieldRejectsImpossibleDate checks that a well-formed but nonexistent
// day is an error rather than a rollover into the next month.
func TestDateFieldRejectsImpossibleDate(t *testing.T) {
	var value time.Time
	tt := start(t, func(c *ui.Context) {
		date.DateField(c, date.DateFieldOptions{Label: "Start", Value: &value})
	}, 380, 260)

	if err := tt.Click("Start"); err != nil {
		t.Fatal(err)
	}
	tt.Type("2026-02-30")
	if !value.IsZero() {
		t.Fatalf("an impossible date committed %v", value)
	}
	if !tt.HasText("Invalid value") {
		t.Fatalf("no invalid message; texts: %v", tt.Texts())
	}
}
