package gallery

import (
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestRangeCalendarSummaryDoesNotShiftGrid pins the demo layout: completing a
// range shortens the summary from "pick an end day" to one line. That must not
// resize the summary column or move the calendar beside it.
func TestRangeCalendarSummaryDoesNotShiftGrid(t *testing.T) {
	day := time.Date(2026, time.October, 14, 0, 0, 0, 0, time.UTC)
	st := &datesState{
		ready: true,
		Today: day,
		Min:   day.AddDate(0, 0, -14),
		Max:   day.AddDate(0, 0, 14),
	}
	st.Range.Value = date.DateRange{Start: day}
	st.Range.Anchor = day
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		rangeBoundedExample(c, st)
	}, 1000, 700)
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	tt.Frame()

	cell := func() ui.Rect {
		r, ok := tt.Find("Friday, October 16, 2026")
		if !ok {
			t.Fatalf("calendar day is not rendered; labels=%q", tt.Texts())
		}
		return r
	}
	before := cell()

	st.Range.Value = date.DateRange{Start: day, End: day.AddDate(0, 0, 8)}
	st.Range.Anchor = time.Time{}
	tt.Frame()
	after := cell()

	if math.Abs(float64(before.X-after.X)) > 0.5 || math.Abs(float64(before.Y-after.Y)) > 0.5 {
		t.Fatalf("completing the range moved the calendar: before=%+v after=%+v", before, after)
	}
}
