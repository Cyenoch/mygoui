package date_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/egoist/mygo/ui"
)

// TestCalendarDayFillFadesIn pins a chosen day taking its fill over a fade
// and a calendar built with a chosen day showing it at once.
func TestCalendarDayFillFadesIn(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Today: march, DayLabel: label}, st)
	}, 360, 420)
	before := sampleRaw(t, tt, "2026-03-20")
	if err := tt.Click("2026-03-20"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(40 * time.Millisecond)
	mid := sampleRaw(t, tt, "2026-03-20")
	tt.Advance(300 * time.Millisecond)
	end := sampleRaw(t, tt, "2026-03-20")
	if mid >= before || mid <= end {
		t.Fatalf("day fill: before %d, mid %d, settled %d", before, mid, end)
	}

	chosen := &date.CalendarState{Value: time.Date(2026, time.March, 20, 12, 0, 0, 0, time.UTC)}
	tt = start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Today: march, DayLabel: label}, chosen)
	}, 360, 420)
	if got := sampleRaw(t, tt, "2026-03-20"); got != end {
		t.Fatalf("a day built chosen shows %d, want %d at once", got, end)
	}
}

// sampleRaw is the red channel inside a day's cell, away from its number.
func sampleRaw(t *testing.T, tt *ui.Tester, day string) int {
	t.Helper()
	r, _, _ := sampleDay(t, tt, day, 0.12, 0.5)
	return int(r)
}
