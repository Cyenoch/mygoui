package date

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// DateRange is a pair of calendar dates. Either endpoint may be zero, as while
// a range is being chosen; Start after End is ordered by the widgets before it
// is reported.
type DateRange struct {
	Start, End time.Time
}

// Empty reports whether neither endpoint is set.
func (r DateRange) Empty() bool { return r.Start.IsZero() && r.End.IsZero() }

// Ordered returns the range with its endpoints in order: an end before its
// start is swapped. A missing endpoint leaves the range as it is.
func (r DateRange) Ordered() DateRange {
	if r.Start.IsZero() || r.End.IsZero() {
		return r
	}
	if compareDay(r.End, r.Start) < 0 {
		return DateRange{Start: dayStart(r.End), End: dayStart(r.Start)}
	}
	return r
}

// Contains reports whether day falls inside the range, endpoints included. A
// range with only a start contains that day and no other, as while the second
// endpoint is still being chosen.
func (r DateRange) Contains(day time.Time) bool {
	if r.Start.IsZero() {
		return false
	}
	start, end := r.Start, r.End
	if !end.IsZero() && compareDay(end, start) < 0 {
		start, end = end, start
	}
	if compareDay(day, start) < 0 {
		return false
	}
	if end.IsZero() {
		return compareDay(day, start) == 0
	}
	return compareDay(day, end) <= 0
}

// RangeCalendarState is the caller-owned state of a RangeCalendar.
type RangeCalendarState struct {
	// Value is the chosen range. Start alone marks a pending selection.
	Value DateRange
	// Anchor is the first endpoint while the second is being chosen, and zero
	// when the range is complete or idle.
	Anchor time.Time
	CalendarViewState
}

// RangeCalendarOptions configures a RangeCalendar. Its labels follow
// CalendarOptions.
type RangeCalendarOptions struct {
	Label  string
	Handle *ui.Handle

	Today         time.Time
	Min, Max      time.Time
	DisabledDates func(time.Time) bool
	FirstWeekday  time.Weekday
	Disabled      bool
	// ShowYearPicker makes the month heading open a year selector.
	ShowYearPicker bool
	// HideYearPicker overrides ShowYearPicker.
	//
	// Deprecated: leave ShowYearPicker false instead.
	HideYearPicker bool

	PrevLabel, NextLabel, YearPickerLabel string
	MonthLabel                            func(time.Time) string
	WeekdayLabel                          func(time.Weekday) string
	DayLabel                              func(time.Time) string

	// OnChange reports a completed range, ordered. Choosing the first endpoint
	// updates State.Value directly and does not call OnChange.
	OnChange func(DateRange)

	autofocus bool
}

// RangeCalendar renders a COSS-style range calendar. Choosing a day starts a
// range; choosing another completes it, swapping the endpoints when the second
// is earlier. A nil state uses a component-local copy.
func RangeCalendar(c *ui.Context, o RangeCalendarOptions, st *RangeCalendarState) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	root := ui.Box(c).MinWidth(0).Shrink(0)
	state := st
	if state == nil {
		state = ui.Local(root, rangeStateKey{}, func() RangeCalendarState { return RangeCalendarState{} })
	}
	spec := calendarSpec{
		label:          o.Label,
		today:          o.Today,
		min:            o.Min,
		max:            o.Max,
		disabledDates:  o.DisabledDates,
		firstWeekday:   o.FirstWeekday,
		disabled:       o.Disabled,
		showYearPicker: o.ShowYearPicker && !o.HideYearPicker,
		autofocus:      o.autofocus,
		prevLabel:      o.PrevLabel,
		nextLabel:      o.NextLabel,
		yearLabel:      o.YearPickerLabel,
		monthLabel:     o.MonthLabel,
		weekdayLabel:   o.WeekdayLabel,
		dayLabel:       o.DayLabel,
	}

	cal := buildCalendar(c, spec, &state.CalendarViewState, state.Value.Start, cellSelection{
		inRange: func(d time.Time) bool { return state.Value.Contains(d) },
		start:   func(d time.Time) bool { return sameDay(d, state.Value.Start) },
		end:     func(d time.Time) bool { return !state.Value.End.IsZero() && sameDay(d, state.Value.End) },
		activate: func(d time.Time) {
			if state.Anchor.IsZero() {
				state.Value = DateRange{Start: d}
				state.Anchor = d
				return
			}
			state.Anchor = time.Time{}
			state.Value = DateRange{Start: state.Value.Start, End: d}.Ordered()
			if o.OnChange != nil {
				o.OnChange(state.Value)
			}
		},
	})
	if o.Handle != nil {
		cal.Bind(o.Handle)
	}
	return cal
}

type rangeStateKey struct{}

// Clear drops the selection and any pending endpoint.
func (s *RangeCalendarState) Clear() {
	s.Value = DateRange{}
	s.Anchor = time.Time{}
}
