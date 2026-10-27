package gallery

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// datesState is the caller-owned state of the date specimens: one selection
// per control plus the open flags of the two popup pickers.
type datesState struct {
	ready bool

	Today time.Time
	Min   time.Time
	Max   time.Time

	CalDefault date.CalendarState
	Cal        date.CalendarState
	Range      date.RangeCalendarState

	CalMonday   date.CalendarState
	CalReadOnly date.CalendarState
	CalStatic   date.CalendarState
	CalYear     date.CalendarState
	RangeMonday date.RangeCalendarState
	RangeFree   date.RangeCalendarState

	Field     time.Time
	FieldFree time.Time
	Picker    time.Time
	RangePick date.DateRange

	PickerCal    date.CalendarState
	RangePickCal date.RangeCalendarState

	PickerFree        time.Time
	PickerFreeCal     date.CalendarState
	PickerFreeOpen    bool
	RangePickFree     date.DateRange
	RangePickFreeCal  date.RangeCalendarState
	RangePickFreeOpen bool

	Time   date.TimeFieldState
	Time12 date.TimeFieldState

	PickerOpen bool
	RangeOpen  bool
}

func (st *datesState) ensure(c *ui.Context) {
	if st.ready {
		return
	}
	st.ready = true
	st.Today = datesStartOfDay(c.Now())
	st.Min = st.Today.AddDate(0, 0, -10)
	st.Max = st.Today.AddDate(0, 0, 24)
	st.CalDefault = date.CalendarState{Value: st.Today}
	st.Cal = date.CalendarState{Value: datesNextWeekday(st.Today)}
	st.Range = date.RangeCalendarState{Value: date.DateRange{
		Start: datesNextWeekday(st.Today.AddDate(0, 0, 3)),
		End:   datesNextWeekday(st.Today.AddDate(0, 0, 10)),
	}}
	st.Field = datesNextWeekday(st.Today)
	st.FieldFree = datesNextWeekday(st.Today.AddDate(0, 0, 5))
	st.Picker = datesNextWeekday(st.Today.AddDate(0, 0, 2))
	st.RangePick = date.DateRange{
		Start: datesNextWeekday(st.Today.AddDate(0, 0, 4)),
		End:   datesNextWeekday(st.Today.AddDate(0, 0, 11)),
	}
	st.CalMonday = date.CalendarState{Value: datesNextWeekday(st.Today.AddDate(0, 0, 1))}
	st.CalReadOnly = date.CalendarState{Value: datesNextWeekday(st.Today)}
	st.CalStatic = date.CalendarState{Value: datesNextWeekday(st.Today.AddDate(0, 0, 2))}
	st.CalYear = date.CalendarState{Value: datesNextWeekday(st.Today.AddDate(0, 0, 2))}
	st.RangeMonday = date.RangeCalendarState{Value: date.DateRange{
		Start: datesNextWeekday(st.Today.AddDate(0, 0, 2)),
		End:   datesNextWeekday(st.Today.AddDate(0, 0, 8)),
	}}
	st.RangeFree = date.RangeCalendarState{Value: date.DateRange{
		Start: datesNextWeekday(st.Today.AddDate(0, 0, 1)),
		End:   datesNextWeekday(st.Today.AddDate(0, 0, 6)),
	}}
	st.PickerFree = datesNextWeekday(st.Today.AddDate(0, 0, 7))
	st.RangePickFree = date.DateRange{
		Start: datesNextWeekday(st.Today.AddDate(0, 0, 3)),
		End:   datesNextWeekday(st.Today.AddDate(0, 0, 9)),
	}
	st.Time = date.TimeFieldState{Value: datesClock(st.Today, 9, 30)}
	st.Time12 = date.TimeFieldState{Value: datesClock(st.Today, 16, 45)}
}

func calendarDefaultExample(c *ui.Context, st *datesState) {
	ui.Row(c).FillWidth().MinWidth(0).Justify(ui.Center).Children(func() {
		date.Calendar(c, date.CalendarOptions{
			Label: "Calendar",
			Today: st.Today,
		}, &st.CalDefault)
	})
}

// calendarBoundedExample bounds the grid and disables weekends.
func calendarBoundedExample(c *ui.Context, st *datesState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		date.Calendar(c, date.CalendarOptions{
			Label:          "Deployment day",
			Today:          st.Today,
			Min:            st.Min,
			Max:            st.Max,
			DisabledDates:  datesWeekend,
			ShowYearPicker: true,
		}, &st.Cal)
		ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, "Selected day", mygoui.TypographyOptions{Type: theme.TypeH6})
			mygoui.Typography(c, datesDayLabel(st.Cal.Value), mygoui.TypographyOptions{Type: theme.TypeBodySM})
			ui.Row(c).Wrap().Gap(theme.Space4).Children(func() {
				datesChip(c, "Weekends disabled")
				datesChip(c, "From "+st.Min.Format("Jan 2"))
				datesChip(c, "Through "+st.Max.Format("Jan 2"))
			})
		})
	})
}

// calendarMondayExample starts the week on Monday.
func calendarMondayExample(c *ui.Context, st *datesState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		date.Calendar(c, date.CalendarOptions{
			Label:          "Review day",
			Today:          st.Today,
			Min:            st.Min,
			Max:            st.Max,
			FirstWeekday:   time.Monday,
			ShowYearPicker: true,
		}, &st.CalMonday)
		ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, datesDayLabel(st.CalMonday.Value), mygoui.TypographyOptions{Type: theme.TypeBodySM})
			datesChip(c, "Monday first")
		})
	})
}

// calendarStaticHeadingExample keeps the default static month heading.
func calendarStaticHeadingExample(c *ui.Context, st *datesState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		date.Calendar(c, date.CalendarOptions{
			Label: "Static heading day",
			Today: st.Today,
			Min:   st.Min,
			Max:   st.Max,
		}, &st.CalStatic)
		ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, datesDayLabel(st.CalStatic.Value), mygoui.TypographyOptions{Type: theme.TypeBodySM})
			datesChip(c, "No year grid")
		})
	})
}

// calendarYearPickerExample enables the optional year selector.
func calendarYearPickerExample(c *ui.Context, st *datesState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		date.Calendar(c, date.CalendarOptions{
			Label:          "Year picker day",
			Today:          st.Today,
			ShowYearPicker: true,
		}, &st.CalYear)
		ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, datesDayLabel(st.CalYear.Value), mygoui.TypographyOptions{Type: theme.TypeBodySM})
			datesChip(c, "Choose year")
		})
	})
}

// calendarReadOnlyExample renders the grid at reduced opacity.
func calendarReadOnlyExample(c *ui.Context, st *datesState) {
	date.Calendar(c, date.CalendarOptions{
		Label:          "Archived day",
		Today:          st.Today,
		Min:            st.Min,
		Max:            st.Max,
		Disabled:       true,
		ShowYearPicker: true,
	}, &st.CalReadOnly)
}

// rangeBoundedExample picks two endpoints over bounded weekdays.
func rangeBoundedExample(c *ui.Context, st *datesState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Label:          "Release window",
			Today:          st.Today,
			Min:            st.Min,
			Max:            st.Max,
			DisabledDates:  datesWeekend,
			ShowYearPicker: true,
		}, &st.Range)
		ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, "Selected range", mygoui.TypographyOptions{Type: theme.TypeH6})
			// Keep both the value and the hint line at a fixed height: a
			// shorter completed range must not reflow the grid beside it.
			ui.Column(c).MinHeight(theme.Scale(c, 40)).Justify(ui.Start).Children(func() {
				mygoui.Typography(c, datesRangeLabel(st.Range.Value), mygoui.TypographyOptions{Type: theme.TypeBodySM})
			})
			ui.Column(c).MinHeight(theme.Scale(c, 32)).Justify(ui.Start).Children(func() {
				if !st.Range.Anchor.IsZero() {
					mygoui.Typography(c, "Start "+st.Range.Anchor.Format("Mon, Jan 2")+" — now pick the end day.",
						mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
				}
			})
			datesChip(c, "Weekends disabled")
		})
	})
}

// rangeMondayExample starts both grids on Monday.
func rangeMondayExample(c *ui.Context, st *datesState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Label:          "Sprint window",
			Today:          st.Today,
			Min:            st.Min,
			Max:            st.Max,
			FirstWeekday:   time.Monday,
			ShowYearPicker: true,
		}, &st.RangeMonday)
		ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, datesRangeLabel(st.RangeMonday.Value), mygoui.TypographyOptions{Type: theme.TypeBodySM})
			datesChip(c, "Monday first")
		})
	})
}

// rangeFreeExample leaves the endpoints unbounded.
func rangeFreeExample(c *ui.Context, st *datesState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		date.RangeCalendar(c, date.RangeCalendarOptions{
			Label:          "Travel window",
			Today:          st.Today,
			ShowYearPicker: true,
		}, &st.RangeFree)
		ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, datesRangeLabel(st.RangeFree.Value), mygoui.TypographyOptions{Type: theme.TypeBodySM})
			datesChip(c, "No bounds")
		})
	})
}

// dateFieldBoundedExample commits complete, in-bounds weekdays.
func dateFieldBoundedExample(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		date.DateField(c, date.DateFieldOptions{
			Label: "Launch date", Description: "YYYY-MM-DD, weekdays only",
			Value: &st.Field, Min: st.Min, Max: st.Max,
			DisabledDates: datesWeekend, Width: 200,
		})
		datesChip(c, "Launch "+st.Field.Format("Mon, Jan 2"))
	})
}

// dateFieldFreeExample accepts any complete ISO date.
func dateFieldFreeExample(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		date.DateField(c, date.DateFieldOptions{
			Label: "Renewal date", Description: "No bounds",
			Value: &st.FieldFree, Width: 200,
		})
		datesChip(c, "Renewal "+st.FieldFree.Format("Mon, Jan 2"))
	})
}

// dateFieldDisabledExample is read-only.
func dateFieldDisabledExample(c *ui.Context, st *datesState) {
	date.DateField(c, date.DateFieldOptions{
		Label: "Archived date", Description: "Disabled field",
		Value: &st.FieldFree, Placeholder: "Locked", Disabled: true, Width: 200,
	})
}

// timeField24Example edits a 24-hour clock.
func timeField24Example(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		date.TimeField(c, date.TimeFieldOptions{
			Label: "Window start", Description: "24-hour clock", Width: 140,
		}, &st.Time)
		datesChip(c, "Start "+st.Time.Value.Format("15:04"))
	})
}

// timeField12Example edits an hh:mm AM/PM clock.
func timeField12Example(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		date.TimeField(c, date.TimeFieldOptions{
			Label: "Review", Description: "12-hour clock", Hour12: true, Width: 140,
		}, &st.Time12)
		datesChip(c, "Review "+st.Time12.Value.Format("03:04 PM"))
	})
}

// timeFieldDisabledExample is read-only.
func timeFieldDisabledExample(c *ui.Context, st *datesState) {
	date.TimeField(c, date.TimeFieldOptions{
		Label: "Archived time", Description: "Disabled field", Disabled: true, Width: 140,
	}, &st.Time)
}

// datePickerBoundedExample opens a bounded popup calendar.
func datePickerBoundedExample(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).MinWidth(theme.Scale(c, 200)).Children(func() {
		date.DatePicker(c, date.DatePickerOptions{
			Label: "Deploy date", Description: "Bounded, weekdays only",
			Value: &st.Picker, State: &st.PickerCal, Open: &st.PickerOpen,
			Min: st.Min, Max: st.Max, DisabledDates: datesWeekend,
			Width: 220, OpenLabel: "Open deploy calendar", ShowYearPicker: true,
			DayLabel: func(d time.Time) string {
				return "Deploy " + d.Format("Monday, January 2, 2006")
			},
		})
		mygoui.Typography(c, datesDayLabel(st.Picker), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// datePickerFreeExample leaves the popup unbounded.
func datePickerFreeExample(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).MinWidth(theme.Scale(c, 200)).Children(func() {
		date.DatePicker(c, date.DatePickerOptions{
			Label: "Travel date", Description: "No bounds",
			Value: &st.PickerFree, State: &st.PickerFreeCal, Open: &st.PickerFreeOpen,
			Width: 220, OpenLabel: "Open travel calendar", ShowYearPicker: true,
		})
		mygoui.Typography(c, datesDayLabel(st.PickerFree), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// datePickerDisabledExample blocks the field and its popup.
func datePickerDisabledExample(c *ui.Context, st *datesState) {
	date.DatePicker(c, date.DatePickerOptions{
		Label: "Archived deploy", Description: "Disabled picker",
		Value: &st.Picker, State: &st.PickerCal, Disabled: true, Width: 220,
		ShowYearPicker: true,
	})
}

// rangePickerBoundedExample edits two bounded endpoints.
func rangePickerBoundedExample(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).MinWidth(theme.Scale(c, 240)).Children(func() {
		date.DateRangePicker(c, date.DateRangePickerOptions{
			Label: "Release window", Description: "Two editable endpoints",
			Value: &st.RangePick, State: &st.RangePickCal, Open: &st.RangeOpen,
			Min: st.Min, Max: st.Max, DisabledDates: datesWeekend,
			Width: 264, OpenLabel: "Open release calendar", ShowYearPicker: true,
			DayLabel: func(d time.Time) string {
				return "Window " + d.Format("Monday, January 2, 2006")
			},
		})
		mygoui.Typography(c, datesRangeLabel(st.RangePick), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// rangePickerFreeExample leaves both endpoints unbounded.
func rangePickerFreeExample(c *ui.Context, st *datesState) {
	ui.Column(c).Gap(theme.Space8).MinWidth(theme.Scale(c, 240)).Children(func() {
		date.DateRangePicker(c, date.DateRangePickerOptions{
			Label: "Travel window", Description: "No bounds",
			Value: &st.RangePickFree, State: &st.RangePickFreeCal, Open: &st.RangePickFreeOpen,
			Width: 264, OpenLabel: "Open travel calendar", ShowYearPicker: true,
		})
		mygoui.Typography(c, datesRangeLabel(st.RangePickFree), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// rangePickerDisabledExample is read-only.
func rangePickerDisabledExample(c *ui.Context, st *datesState) {
	date.DateRangePicker(c, date.DateRangePickerOptions{
		Label: "Archived window", Description: "Disabled picker",
		Value: &st.RangePick, State: &st.RangePickCal, Disabled: true, Width: 264,
		ShowYearPicker: true,
	})
}

func datesChip(c *ui.Context, label string) {
	mygoui.Chip(c, mygoui.ChipOptions{
		Label: label, Variant: mygoui.ChipSoft, Size: mygoui.ChipSmall,
	})
}

func datesWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func datesStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func datesNextWeekday(t time.Time) time.Time {
	for datesWeekend(t) {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func datesClock(day time.Time, hour, minute int) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.Location())
}

func datesDayLabel(t time.Time) string {
	if t.IsZero() {
		return "No day selected"
	}
	return t.Format("Mon, Jan 2, 2006")
}

func datesRangeLabel(r date.DateRange) string {
	if r.Empty() {
		return "No range selected"
	}
	if r.End.IsZero() {
		return r.Start.Format("Mon, Jan 2, 2006") + " → pick an end day"
	}
	return r.Start.Format("Jan 2, 2006") + " → " + r.End.Format("Jan 2, 2006")
}
