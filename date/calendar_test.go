package date_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

var march = time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC)

func label(d time.Time) string { return d.Format("2006-01-02") }

func start(t *testing.T, build func(c *ui.Context), w, h int) *ui.Tester {
	t.Helper()
	return ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Padding(20).Children(func() { build(c) })
	}, w, h)
}

// yearLabels returns the year numbers the year picker shows.
func yearLabels(tt *ui.Tester) map[int]bool {
	years := make(map[int]bool)
	for _, s := range tt.Texts() {
		if len(s) != 4 {
			continue
		}
		if y, err := strconv.Atoi(s); err == nil && y > 0 {
			years[y] = true
		}
	}
	return years
}

func TestCalendarYearPickerDefaultPage(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Today: march, DayLabel: label}, st)
	}, 360, 420)

	if err := tt.Click("March 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	years := yearLabels(tt)
	// A page around the visible year, not a scroll through two centuries.
	if len(years) < 12 || len(years) > 24 {
		t.Fatalf("year picker shows %d years, want a page of at least 12: %v", len(years), years)
	}
	for y := range years {
		if y < 2026-12 || y > 2026+12 {
			t.Fatalf("year %d is outside the page around 2026: %v", y, years)
		}
	}
	for _, y := range []int{2024, 2025, 2026, 2027} {
		if !years[y] {
			t.Fatalf("year picker is missing %d: %v", y, years)
		}
	}

	// Choosing an earlier year moves the visible month into it and collapses
	// the picker.
	if err := tt.Click("2024"); err != nil {
		t.Fatal(err)
	}
	if st.Month.Year() != 2024 || st.Month.Month() != time.March {
		t.Fatalf("month = %v, want March 2024", st.Month)
	}
	if st.YearOpen {
		t.Fatal("the year picker stayed open after choosing a year")
	}
	if !tt.HasText("March 2024") {
		t.Fatalf("heading after choosing 2024: %v", tt.Texts())
	}
	// Month navigation still works from the chosen year.
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if st.Month.Year() != 2024 || st.Month.Month() != time.April {
		t.Fatalf("after Next: month = %v, want April 2024", st.Month)
	}
}

func TestCalendarYearPickerBounds(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Today:    time.Date(2021, time.June, 15, 12, 0, 0, 0, time.UTC),
			Min:      time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			Max:      time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
			DayLabel: label,
		}, st)
	}, 360, 420)

	if err := tt.Click("June 2021, Choose year"); err != nil {
		t.Fatal(err)
	}
	years := yearLabels(tt)
	for _, y := range []int{2020, 2021, 2022, 2023} {
		if !years[y] {
			t.Fatalf("bounded year picker is missing %d: %v", y, years)
		}
	}
	for _, y := range []int{2019, 2024} {
		if years[y] {
			t.Fatalf("bounded year picker shows out-of-range %d: %v", y, years)
		}
	}

	if err := tt.Click("2023"); err != nil {
		t.Fatal(err)
	}
	if st.Month.Year() != 2023 || st.Month.Month() != time.June {
		t.Fatalf("month = %v, want June 2023", st.Month)
	}
	if st.YearOpen {
		t.Fatal("the year picker stayed open after choosing a year")
	}
}

func TestCalendarMonthNavigation(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Today: march, DayLabel: label}, st)
	}, 320, 420)

	if !tt.HasText("March 2026") {
		t.Fatalf("heading missing; texts: %v", tt.Texts())
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("April 2026") {
		t.Fatalf("next month missing; texts: %v", tt.Texts())
	}
	if st.Month.Month() != time.April {
		t.Fatalf("state month = %v, want April", st.Month)
	}
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("March 2026") {
		t.Fatalf("previous month missing; texts: %v", tt.Texts())
	}
	if st.Month.Month() != time.March {
		t.Fatalf("state month = %v, want March", st.Month)
	}
}

func TestCalendarSelectsDay(t *testing.T) {
	st := &date.CalendarState{}
	var got time.Time
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Today:    march,
			DayLabel: label,
			OnChange: func(d time.Time) { got = d },
		}, st)
	}, 320, 420)

	if err := tt.Click("2026-03-18"); err != nil {
		t.Fatal(err)
	}
	if st.Value.Day() != 18 || st.Value.Month() != time.March {
		t.Fatalf("state value = %v, want 2026-03-18", st.Value)
	}
	if got.Day() != 18 {
		t.Fatalf("callback value = %v", got)
	}
}

func TestCalendarRejectsDisabledDate(t *testing.T) {
	st := &date.CalendarState{Value: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)}
	calls := 0
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Today:         march,
			DayLabel:      label,
			DisabledDates: func(d time.Time) bool { return d.Day() == 10 },
			OnChange:      func(time.Time) { calls++ },
		}, st)
	}, 320, 420)

	if err := tt.Click("2026-03-10"); err != nil {
		t.Fatal(err)
	}
	if st.Value.Day() != 1 {
		t.Fatalf("disabled day changed value to %v", st.Value)
	}
	if calls != 0 {
		t.Fatalf("OnChange called %d times for a disabled day", calls)
	}
	if err := tt.Click("2026-03-11"); err != nil {
		t.Fatal(err)
	}
	if st.Value.Day() != 11 || calls != 1 {
		t.Fatalf("enabled day: value %v calls %d", st.Value, calls)
	}
}

func TestCalendarBounds(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Today:    march,
			Min:      time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC),
			Max:      time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC),
			DayLabel: label,
		}, st)
	}, 320, 420)

	if err := tt.Click("2026-03-01"); err != nil {
		t.Fatal(err)
	}
	if st.Value.Day() != 1 {
		t.Fatalf("in-bounds first day not selected: %v", st.Value)
	}
	if !tt.HasText("March 2026") {
		t.Fatalf("calendar left March")
	}
}

func TestRangeCalendarOrdersEndpoints(t *testing.T) {
	st := &date.RangeCalendarState{}
	var got date.DateRange
	calls := 0
	tt := start(t, func(c *ui.Context) {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Today:    march,
			DayLabel: label,
			OnChange: func(r date.DateRange) { got, calls = r, calls+1 },
		}, st)
	}, 320, 420)

	if err := tt.Click("2026-03-20"); err != nil {
		t.Fatal(err)
	}
	if calls != 0 || st.Value.Start.Day() != 20 || !st.Value.End.IsZero() {
		t.Fatalf("after first endpoint: value %+v calls %d", st.Value, calls)
	}
	if err := tt.Click("2026-03-10"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("OnChange calls = %d, want 1", calls)
	}
	if got.Start.Day() != 10 || got.End.Day() != 20 {
		t.Fatalf("range = %v..%v, want 10..20", got.Start, got.End)
	}
	if !got.Contains(time.Date(2026, time.March, 15, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("range does not contain its middle")
	}
	if got.Contains(time.Date(2026, time.March, 21, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("range contains a day after its end")
	}
	if st.Anchor.IsZero() != true {
		t.Fatal("anchor should be cleared after a complete range")
	}
}

func TestRangeCalendarStartOnlyPaintsNoBand(t *testing.T) {
	nov := time.Date(2026, time.November, 9, 12, 0, 0, 0, time.UTC)
	tt, pal := newRangeCalendar(t, time.Time{}, time.Time{}, nov, nil)

	// Choosing only the start must not paint a band through the rest of the
	// month: just the start circle, as while the end is still being chosen.
	if err := tt.Click("2026-11-09"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	r, g, b := sampleDay(t, tt, "2026-11-09", 0.2, 0.5)
	assertColor(t, "start 9 fill", r, g, b, pal.Accent)
	for _, day := range []string{"2026-11-10", "2026-11-11", "2026-11-12"} {
		assertNoFill(t, tt, day, pal.Background)
	}

	// Completing the range paints the days between the endpoints soft.
	if err := tt.Click("2026-11-13"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	soft := pal.AccentSoft.Over(pal.Background)
	for _, day := range []string{"2026-11-10", "2026-11-11", "2026-11-12"} {
		r, g, b := sampleDay(t, tt, day, 0.2, 0.5)
		assertColor(t, day+" fill", r, g, b, soft)
	}
	r, g, b = sampleDay(t, tt, "2026-11-09", 0.2, 0.5)
	assertColor(t, "start 9 fill", r, g, b, pal.Accent)
	r, g, b = sampleDay(t, tt, "2026-11-13", 0.2, 0.5)
	assertColor(t, "end 13 fill", r, g, b, pal.Accent)
}

func TestDatePickerOpensAndCommits(t *testing.T) {
	var value time.Time
	tt := start(t, func(c *ui.Context) {
		date.DatePicker(c, date.DatePickerOptions{
			Label:    "Appointment",
			Today:    march,
			DayLabel: label,
			Value:    &value,
		})
	}, 420, 560)

	if err := tt.Click("Open calendar"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("March 2026") {
		t.Fatalf("popup calendar missing; texts: %v", tt.Texts())
	}
	if err := tt.Click("2026-03-12"); err != nil {
		t.Fatal(err)
	}
	if value.Day() != 12 {
		t.Fatalf("committed value = %v, want 2026-03-12", value)
	}
	if tt.HasText("March 2026") {
		t.Fatalf("popup stayed open; texts: %v", tt.Texts())
	}
}

func TestDateFieldCommitsISOAndRejectsOutOfBounds(t *testing.T) {
	var value time.Time
	min := time.Date(2026, time.March, 10, 0, 0, 0, 0, time.UTC)
	tt := start(t, func(c *ui.Context) {
		date.DateField(c, date.DateFieldOptions{
			Label:    "Start",
			Min:      min,
			Value:    &value,
			OnChange: func(time.Time) {},
		})
	}, 420, 320)

	if err := tt.Click("Start"); err != nil {
		t.Fatal(err)
	}
	tt.Type("2026-03-05")
	if !value.IsZero() {
		t.Fatalf("out-of-bounds date committed: %v", value)
	}
	if !tt.HasText("Invalid value") {
		t.Fatalf("invalid message missing; texts: %v", tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("2026-03-15")
	if value.Day() != 15 {
		t.Fatalf("committed value = %v, want 2026-03-15", value)
	}
	if tt.HasText("Invalid value") {
		t.Fatalf("invalid message stayed; texts: %v", tt.Texts())
	}
}

func TestTimeFieldCommits(t *testing.T) {
	st := &date.TimeFieldState{}
	var got time.Time
	tt := start(t, func(c *ui.Context) {
		date.TimeField(c, date.TimeFieldOptions{
			Label:    "Alarm",
			OnChange: func(v time.Time) { got = v },
		}, st)
	}, 320, 260)

	if err := tt.Click("Alarm"); err != nil {
		t.Fatal(err)
	}
	tt.Type("14:30")
	if st.Value.Hour() != 14 || st.Value.Minute() != 30 {
		t.Fatalf("state = %v, want 14:30", st.Value)
	}
	if got.Hour() != 14 || got.Minute() != 30 {
		t.Fatalf("callback = %v, want 14:30", got)
	}
}

func TestTimeFieldHour12(t *testing.T) {
	st := &date.TimeFieldState{}
	tt := start(t, func(c *ui.Context) {
		date.TimeField(c, date.TimeFieldOptions{Label: "Alarm", Hour12: true}, st)
	}, 320, 260)

	if err := tt.Click("Alarm"); err != nil {
		t.Fatal(err)
	}
	tt.Type("09:05 PM")
	if st.Value.Hour() != 21 || st.Value.Minute() != 5 {
		t.Fatalf("state = %v, want 21:05", st.Value)
	}
}

func TestRangeCalendarPaintsEndpointsSolidAndMiddleSoft(t *testing.T) {
	st := &date.RangeCalendarState{Value: date.DateRange{
		Start: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
	}}
	var pal theme.Palette
	tt := start(t, func(c *ui.Context) {
		pal = theme.Colors(c)
		date.RangeCalendar(c, date.RangeCalendarOptions{Today: march, DayLabel: label}, st)
	}, 360, 480)
	img := tt.Image()
	at := func(day string) (uint8, uint8, uint8) {
		r, ok := tt.Find(day)
		if !ok {
			t.Fatalf("%s cell missing", day)
		}
		c := img.RGBAAt(int(r.X+r.W*0.2), int(r.Y+r.H*0.2))
		return c.R, c.G, c.B
	}
	r, g, b := at("2026-03-20")
	if r != pal.Accent.R || g != pal.Accent.G || b != pal.Accent.B {
		t.Fatalf("start fill = %d,%d,%d, want accent", r, g, b)
	}
	r, g, b = at("2026-03-21")
	want := pal.AccentSoft.Over(pal.Background)
	if diff(r, want.R) > 3 || diff(g, want.G) > 3 || diff(b, want.B) > 3 {
		t.Fatalf("middle fill = %d,%d,%d, want soft %d,%d,%d", r, g, b, want.R, want.G, want.B)
	}
	r, g, b = at("2026-03-23")
	if r != pal.Accent.R || g != pal.Accent.G || b != pal.Accent.B {
		t.Fatalf("end fill = %d,%d,%d, want accent", r, g, b)
	}
}

// sampleDay reads the pixel at a fraction of a day cell's box.
func sampleDay(t *testing.T, tt *ui.Tester, day string, fx, fy float32) (uint8, uint8, uint8) {
	t.Helper()
	cell, ok := tt.Find(day)
	if !ok {
		t.Fatalf("%s cell missing", day)
	}
	px := tt.Image().RGBAAt(int(cell.X+fx*cell.W), int(cell.Y+fy*cell.H))
	return px.R, px.G, px.B
}

func assertColor(t *testing.T, what string, r, g, b uint8, want ui.Color) {
	t.Helper()
	if diff(r, want.R) > 3 || diff(g, want.G) > 3 || diff(b, want.B) > 3 {
		t.Fatalf("%s = %d,%d,%d, want %d,%d,%d", what, r, g, b, want.R, want.G, want.B)
	}
}

// newRangeCalendar renders a range calendar over one month for sampling.
func newRangeCalendar(t *testing.T, from, to, today time.Time, disabled func(time.Time) bool) (*ui.Tester, theme.Palette) {
	t.Helper()
	st := &date.RangeCalendarState{Value: date.DateRange{Start: from, End: to}}
	var pal theme.Palette
	tt := start(t, func(c *ui.Context) {
		pal = theme.Colors(c)
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Today: today, DayLabel: label, DisabledDates: disabled,
		}, st)
	}, 360, 480)
	return tt, pal
}

// assertNoFill checks that nothing paints the border of a day cell, where a
// day, range or hover background would show. The centred number and its
// strike stay clear of the border.
func assertNoFill(t *testing.T, tt *ui.Tester, day string, bg ui.Color) {
	t.Helper()
	cell, ok := tt.Find(day)
	if !ok {
		t.Fatalf("%s cell missing", day)
	}
	const ring = 8
	img := tt.Image()
	for y := range int(cell.H) {
		for x := range int(cell.W) {
			if x >= ring && x < int(cell.W)-ring && y >= ring && y < int(cell.H)-ring {
				continue
			}
			px := img.RGBAAt(int(cell.X)+x, int(cell.Y)+y)
			if diff(px.R, bg.R)+diff(px.G, bg.G)+diff(px.B, bg.B) > 12 {
				t.Fatalf("%s painted at offset %d,%d = %d,%d,%d", day, x, y, px.R, px.G, px.B)
			}
		}
	}
}

func TestRangeCalendarLeavesDisabledDaysPlain(t *testing.T) {
	// 31 is disabled and today while the range ends the day before: the cell
	// stays plain rather than gaining a today or range fill.
	day31 := func(d time.Time) bool { return d.Day() == 31 && d.Month() == time.March }
	tt, pal := newRangeCalendar(t,
		time.Date(2026, 3, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
		day31)
	assertNoFill(t, tt, "2026-03-31", pal.Background)

	// Endpoints stay solid accent, the days between them the soft band.
	r, g, b := sampleDay(t, tt, "2026-03-27", 0.2, 0.5)
	assertColor(t, "start 27 fill", r, g, b, pal.Accent)
	r, g, b = sampleDay(t, tt, "2026-03-28", 0.2, 0.5)
	assertColor(t, "middle 28 fill", r, g, b, pal.AccentSoft.Over(pal.Background))
	r, g, b = sampleDay(t, tt, "2026-03-30", 0.2, 0.5)
	assertColor(t, "end 30 fill", r, g, b, pal.Accent)
}

func TestRangeCalendarDisabledDayStopsBand(t *testing.T) {
	// 24 is disabled in the middle of 22-27: it gets no band, and the band
	// rounds off on either side of it.
	day24 := func(d time.Time) bool { return d.Day() == 24 && d.Month() == time.March }
	tt, pal := newRangeCalendar(t,
		time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 27, 0, 0, 0, 0, time.UTC),
		march, day24)
	soft := pal.AccentSoft.Over(pal.Background)
	assertNoFill(t, tt, "2026-03-24", pal.Background)

	r, g, b := sampleDay(t, tt, "2026-03-23", 0.98, 0.05)
	assertColor(t, "23 stops before 24", r, g, b, pal.Background)
	r, g, b = sampleDay(t, tt, "2026-03-25", 0.02, 0.05)
	assertColor(t, "25 resumes after 24", r, g, b, pal.Background)
	r, g, b = sampleDay(t, tt, "2026-03-23", 0.02, 0.05)
	assertColor(t, "23 joins 22", r, g, b, soft)
	r, g, b = sampleDay(t, tt, "2026-03-25", 0.98, 0.05)
	assertColor(t, "25 joins 26", r, g, b, soft)

	r, g, b = sampleDay(t, tt, "2026-03-22", 0.2, 0.5)
	assertColor(t, "start 22 fill", r, g, b, pal.Accent)
	r, g, b = sampleDay(t, tt, "2026-03-27", 0.2, 0.5)
	assertColor(t, "end 27 fill", r, g, b, pal.Accent)
}

func TestRangeCalendarBandRoundsAtWeekEdges(t *testing.T) {
	st := &date.RangeCalendarState{Value: date.DateRange{
		Start: time.Date(2026, 3, 27, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC),
	}}
	var pal theme.Palette
	tt := start(t, func(c *ui.Context) {
		pal = theme.Colors(c)
		date.RangeCalendar(c, date.RangeCalendarOptions{Today: march, DayLabel: label}, st)
	}, 360, 480)
	soft := pal.AccentSoft.Over(pal.Background)

	// 27 is Friday and 28 Saturday ends the week, so the band wraps: square
	// where 28 continues into Saturday, rounded at the row's end and the next
	// row's start.
	r, g, b := sampleDay(t, tt, "2026-03-27", 0.98, 0.05)
	assertColor(t, "27 seam with 28", r, g, b, soft)
	r, g, b = sampleDay(t, tt, "2026-03-28", 0.02, 0.05)
	assertColor(t, "28 seam with 27", r, g, b, soft)
	r, g, b = sampleDay(t, tt, "2026-03-28", 0.98, 0.05)
	assertColor(t, "28 rounded week end", r, g, b, pal.Background)
	r, g, b = sampleDay(t, tt, "2026-03-29", 0.02, 0.05)
	assertColor(t, "29 rounded week start", r, g, b, pal.Background)
	r, g, b = sampleDay(t, tt, "2026-03-29", 0.98, 0.05)
	assertColor(t, "29 seam with 30", r, g, b, soft)
	r, g, b = sampleDay(t, tt, "2026-03-30", 0.02, 0.05)
	assertColor(t, "30 seam with 29", r, g, b, soft)
	r, g, b = sampleDay(t, tt, "2026-03-30", 0.98, 0.05)
	assertColor(t, "30 rounded cap", r, g, b, pal.Background)
}

func diff(a, b uint8) int {
	d := int(a) - int(b)
	if d < 0 {
		return -d
	}
	return d
}
