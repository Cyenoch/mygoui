package date_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func day(m time.Month, d int) time.Time { return time.Date(2026, m, d, 0, 0, 0, 0, time.UTC) }

// TestCalendarDaysAreGridCellsSelectedWhenChosen checks what assistive
// technology hears of the day grid: a grid of cells, the chosen one selected.
func TestCalendarDaysAreGridCellsSelectedWhenChosen(t *testing.T) {
	st := &date.CalendarState{Value: day(time.March, 12)}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Label: "Days", Today: march, DayLabel: label}, st)
	}, 360, 560)
	if _, ok := tt.AccessByRole(ui.RoleGrid, "Days"); !ok {
		t.Fatal("the day grid is not a grid named Days")
	}
	chosen, ok := tt.AccessByRole(ui.RoleCell, "2026-03-12")
	if !ok || !chosen.Selected {
		t.Fatalf("the chosen day is not a selected cell: %+v", chosen)
	}
	if other, _ := tt.AccessByRole(ui.RoleCell, "2026-03-13"); other.Selected {
		t.Fatal("an unchosen day is selected")
	}
}

// TestCalendarArrowsStopAtTheBounds checks that the day keys keep the cursor
// on days between Min and Max.
func TestCalendarArrowsStopAtTheBounds(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{
			Label: "Days", Today: march, DayLabel: label,
			Min: day(time.March, 10), Max: day(time.March, 20),
		}, st)
	}, 360, 560)
	if err := tt.Click("2026-03-15"); err != nil {
		t.Fatal(err)
	}
	for _, s := range []struct {
		key  ui.Key
		want string
	}{
		{ui.KeyUp, "2026-03-10"},
		{ui.KeyLeft, "2026-03-10"},
		{ui.KeyDown, "2026-03-17"},
		{ui.KeyDown, "2026-03-20"},
		{ui.KeyRight, "2026-03-20"},
		{ui.KeyEnd, "2026-03-20"},
	} {
		tt.Key(0, s.key)
		if got := label(st.Cursor); got != s.want || st.Month.Month() != time.March {
			t.Fatalf("after %v the cursor is %s in %v, want %s", s.key, got, st.Month.Month(), s.want)
		}
	}
}

// TestCalendarHomeEndKeepTheWeekAndShiftPagesAYear checks the week and year
// keys of the day grid.
func TestCalendarHomeEndKeepTheWeekAndShiftPagesAYear(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Label: "Days", Today: march, DayLabel: label, FirstWeekday: time.Monday}, st)
	}, 360, 560)
	if err := tt.Click("2026-03-18"); err != nil { // a Wednesday
		t.Fatal(err)
	}
	for _, s := range []struct {
		mods ui.Modifiers
		key  ui.Key
		want string
	}{
		{0, ui.KeyHome, "2026-03-16"},
		{0, ui.KeyEnd, "2026-03-22"},
		{ui.Shift, ui.KeyPageDown, "2027-03-22"},
		{ui.Shift, ui.KeyPageUp, "2026-03-22"},
	} {
		tt.Key(s.mods, s.key)
		if got := label(st.Cursor); got != s.want || !tt.Focused(s.want) {
			t.Fatalf("after %v the cursor is %s, want %s", s.key, got, s.want)
		}
	}
}

// TestCalendarArrowsMirrorRightToLeft checks that Left and Right follow the
// mirrored grid.
func TestCalendarArrowsMirrorRightToLeft(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Label: "Days", Today: march, DayLabel: label}, st)
	}, 360, 560)
	tt.SetDirection(ui.RightToLeft)
	tt.Frame()
	if err := tt.Click("2026-03-15"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyLeft)
	if got := label(st.Cursor); got != "2026-03-16" {
		t.Fatalf("Left right to left moved to %s, want 2026-03-16", got)
	}
	tt.Key(0, ui.KeyRight)
	if got := label(st.Cursor); got != "2026-03-15" {
		t.Fatalf("Right right to left moved to %s, want 2026-03-15", got)
	}
}

// TestCalendarYearGridEscapeRestoresTheYear checks that Escape leaves the
// year grid without the year its arrows reached.
func TestCalendarYearGridEscapeRestoresTheYear(t *testing.T) {
	st := &date.CalendarState{}
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Label: "Days", Today: march, DayLabel: label, ShowYearPicker: true}, st)
	}, 360, 560)
	if err := tt.Click("March 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if st.YearOpen || st.Month.Year() != 2026 || !tt.Focused("March 2026, Choose year") {
		t.Fatalf("Escape: open %v, month %v", st.YearOpen, st.Month)
	}
}

// TestDateControlsTakeTheirNamesFromTheLocale checks the names the date
// controls give when the caller gives none.
func TestDateControlsTakeTheirNamesFromTheLocale(t *testing.T) {
	locale.Set(locale.Messages{ChooseYear: "Jahr wählen", OpenCalendar: "Kalender öffnen",
		StartDate: "Beginn", EndDate: "Ende"})
	t.Cleanup(func() { locale.Set(locale.Messages{}) })
	tt := start(t, func(c *ui.Context) {
		date.Calendar(c, date.CalendarOptions{Label: "Days", Today: march, DayLabel: label, ShowYearPicker: true}, nil)
		date.DateRangePicker(c, date.DateRangePickerOptions{Label: "Stay", Today: march})
	}, 420, 760)
	for _, name := range []string{"March 2026, Jahr wählen", "Kalender öffnen", "Beginn", "Ende"} {
		if _, ok := tt.Access(name); !ok {
			t.Fatalf("no control is named %q", name)
		}
	}
}

// TestRangeCalendarEscapeCancelsThePendingEndpoint checks that Escape after a
// first endpoint drops it and restores the range chosen before.
func TestRangeCalendarEscapeCancelsThePendingEndpoint(t *testing.T) {
	before := date.DateRange{Start: day(time.March, 2), End: day(time.March, 4)}
	st := &date.RangeCalendarState{Value: before}
	calls := 0
	tt := start(t, func(c *ui.Context) {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Label: "Days", Today: march, DayLabel: label,
			OnChange: func(date.DateRange) { calls++ },
		}, st)
		ui.Button(c, "After")
	}, 360, 560)
	if err := tt.Click("2026-03-20"); err != nil {
		t.Fatal(err)
	}
	if st.Anchor.IsZero() {
		t.Fatal("the first click did not start a range")
	}
	tt.Key(0, ui.KeyEscape)
	if !st.Anchor.IsZero() || st.Value != before || calls != 0 {
		t.Fatalf("Escape: anchor %v, value %+v, calls %d", st.Anchor, st.Value, calls)
	}
}

// TestRangeCalendarPreviewsTheRangeUnderThePointer checks that, with a start
// chosen, the days up to the one under the pointer are painted as the band,
// at once, and that leaving the grid removes the preview.
func TestRangeCalendarPreviewsTheRangeUnderThePointer(t *testing.T) {
	nov := time.Date(2026, time.November, 9, 12, 0, 0, 0, time.UTC)
	tt, pal := newRangeCalendar(t, time.Time{}, time.Time{}, nov, nil)
	if err := tt.Click("2026-11-09"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(200 * time.Millisecond)
	assertNoFill(t, tt, "2026-11-11", pal.Background)
	end, _ := tt.Find("2026-11-13")
	tt.Move(end.X+end.W/2, end.Y+end.H/2)
	tt.Frame()
	soft := pal.Accent.Over(pal.Background)
	r, g, b := sampleDay(t, tt, "2026-11-11", 0.2, 0.5)
	assertColor(t, "previewed 11", r, g, b, soft)
	tt.Move(1, 1)
	tt.Frame()
	assertNoFill(t, tt, "2026-11-11", pal.Background)
}

// TestDatePickerTriggersReportExpanded checks that the calendar triggers
// tell assistive technology whether their popup shows.
func TestDatePickerTriggersReportExpanded(t *testing.T) {
	var one, two bool
	tt := start(t, func(c *ui.Context) {
		date.DatePicker(c, date.DatePickerOptions{Label: "Due", Today: march, Open: &one, OpenLabel: "Due calendar"})
		date.DateRangePicker(c, date.DateRangePickerOptions{Label: "Stay", Today: march, Open: &two, OpenLabel: "Stay calendar"})
	}, 420, 760)
	for _, name := range []string{"Due calendar", "Stay calendar"} {
		if a, _ := tt.AccessByRole(ui.RoleButton, name); a.Expanded {
			t.Fatalf("%s is expanded while closed", name)
		}
	}
	one = true
	tt.Frame()
	if a, _ := tt.AccessByRole(ui.RoleButton, "Due calendar"); !a.Expanded {
		t.Fatal("the date picker trigger is not expanded while its calendar shows")
	}
	one, two = false, true
	tt.Frame()
	tt.Advance(300 * time.Millisecond)
	if a, _ := tt.AccessByRole(ui.RoleButton, "Stay calendar"); !a.Expanded {
		t.Fatal("the range picker trigger is not expanded while its calendar shows")
	}
}

// TestTimeFieldRejectsTimesOutsideItsBoundsAndStep checks Min, Max and Step:
// a time before Min, after Max or off the step is invalid and not committed.
func TestTimeFieldRejectsTimesOutsideItsBoundsAndStep(t *testing.T) {
	st := &date.TimeFieldState{}
	tt := start(t, func(c *ui.Context) {
		date.TimeField(c, date.TimeFieldOptions{
			Label: "Start", Min: 9 * time.Hour, Max: 17 * time.Hour, Step: 15 * time.Minute,
		}, st)
	}, 320, 260)
	if err := tt.Click("Start"); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"08:45", "17:15", "10:10"} {
		tt.Type(text)
		if !st.Value.IsZero() || !tt.HasText("Invalid value") {
			t.Fatalf("%s committed %v", text, st.Value)
		}
		for range len(text) {
			tt.Key(0, ui.KeyBackspace)
		}
	}
	tt.Type("10:15")
	if st.Value.Hour() != 10 || st.Value.Minute() != 15 || tt.HasText("Invalid value") {
		t.Fatalf("10:15 committed %v", st.Value)
	}
}

// TestCalendarViewsSwapOverTheStateDuration checks that the year grid takes
// the theme's state duration to fade in.
func TestCalendarViewsSwapOverTheStateDuration(t *testing.T) {
	th := theme.Default()
	th.Motion.State = 400 * time.Millisecond
	if err := theme.Use(th); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = theme.Use(theme.Default()) })
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		ui.Column(c).Fill().Padding(20).Background(pal.Background).Children(func() {
			date.Calendar(c, date.CalendarOptions{Today: march, ShowYearPicker: true}, nil)
		})
	}, 360, 480)
	if err := tt.Click("March 2026, Choose year"); err != nil {
		t.Fatal(err)
	}
	cell, _ := tt.Find("2026")
	at := func() (uint8, uint8, uint8) {
		px := tt.Image().RGBAAt(int(cell.X+cell.W*0.3), int(cell.Y+cell.H*0.5))
		return px.R, px.G, px.B
	}
	tt.Advance(250 * time.Millisecond)
	if r, g, b := at(); diff(r, pal.Primary.R) <= 3 && diff(g, pal.Primary.G) <= 3 && diff(b, pal.Primary.B) <= 3 {
		t.Fatal("the year grid settled before the state duration")
	}
	tt.Advance(250 * time.Millisecond)
	r, g, b := at()
	assertColor(t, "settled year", r, g, b, pal.Primary)
}
