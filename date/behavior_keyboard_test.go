package date_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/egoist/mygo/ui"
)

// TestCalendarDayGridKeysMoveTheCursor checks the keyboard path through the
// day grid: Tab reaches it after the month buttons with the cursor on today,
// the week and month keys move the cursor, Home goes to the start of the
// week, Space chooses, and Tab leaves.
func TestCalendarDayGridKeysMoveTheCursor(t *testing.T) {
	st := &date.CalendarState{}
	var chosen []time.Time
	tt := start(t, func(c *ui.Context) {
		ui.Button(c, "Before")
		date.Calendar(c, date.CalendarOptions{
			Label: "Days", Today: march, DayLabel: label,
			OnChange: func(d time.Time) { chosen = append(chosen, d) },
		}, st)
		ui.Button(c, "After")
	}, 360, 560)

	for range 4 { // Before, Previous, Next, then the grid.
		tt.Key(0, ui.KeyTab)
	}
	if !tt.Focused("Days") || !tt.Focused("2026-03-15") {
		t.Fatal("Tab did not reach the day grid with the cursor on today")
	}
	steps := []struct {
		key  ui.Key
		want string
	}{
		{ui.KeyDown, "2026-03-22"},
		{ui.KeyUp, "2026-03-15"},
		{ui.KeyLeft, "2026-03-14"},
		{ui.KeyHome, "2026-03-08"},
		{ui.KeyPageUp, "2026-02-08"},
	}
	for _, s := range steps {
		tt.Key(0, s.key)
		if got := st.Cursor.Format("2006-01-02"); got != s.want || !tt.Focused(s.want) {
			t.Fatalf("after %v the cursor is %s, want %s", s.key, got, s.want)
		}
	}
	if len(chosen) != 0 {
		t.Fatalf("moving the cursor chose %v", chosen)
	}
	tt.Key(0, ui.KeySpace)
	if len(chosen) != 1 || label(chosen[0]) != "2026-02-08" {
		t.Fatalf("Space chose %v, want 2026-02-08", chosen)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab did not leave the day grid")
	}
}

// TestCalendarYearGridArrowsMoveTheYear checks the year grid's arrows: across
// by one year, up and down by a row of three, and Enter returns to the heading.
func TestCalendarYearGridArrowsMoveTheYear(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Label: "Days", Today: march, DayLabel: label, ShowYearPicker: true,
		}, st)
	}, 360, 560)

	if err := tt.Click("March 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("2026") {
		t.Fatal("opening the year grid did not put the cursor on the shown year")
	}
	for _, s := range []struct {
		key  ui.Key
		want int
	}{{ui.KeyRight, 2027}, {ui.KeyUp, 2024}, {ui.KeyDown, 2027}, {ui.KeyLeft, 2026}} {
		tt.Key(0, s.key)
		if st.Month.Year() != s.want {
			t.Fatalf("after %v the year is %d, want %d", s.key, st.Month.Year(), s.want)
		}
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if st.YearOpen || st.Month.Year() != 2023 || !tt.Focused("March 2023, Choose year") {
		t.Fatalf("Enter: year grid open %v, month %v", st.YearOpen, st.Month)
	}
}
