package gallery

import (
	"github.com/egoist/mygo/ui"
)

// datesComponentState keeps every date page on the shared specimen state, so
// the calendar, range calendar, fields and pickers agree on today and bounds.
func datesComponentState(c *ui.Context, s *State) *datesState {
	st := StateOf[datesState](s, "dates")
	st.ensure(c)
	return st
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "calendar", Title: "Calendar", Group: "Dates", Order: 10,
		Description: "A single-day grid with bounded days, disabled dates and keyboard selection.",
		Examples: []ComponentExample{
			{Title: "Bounded dates", Description: "Selectable days are bounded and weekends stay disabled.",
				Demo: func(c *ui.Context, s *State) { calendarBoundedExample(c, datesComponentState(c, s)) }},
			{Title: "Monday first", Description: "The week starts on Monday instead of Sunday.",
				Demo: func(c *ui.Context, s *State) { calendarMondayExample(c, datesComponentState(c, s)) }},
			{Title: "Read-only", Description: "A disabled calendar dims its whole grid.",
				Demo: func(c *ui.Context, s *State) { calendarReadOnlyExample(c, datesComponentState(c, s)) }},
		},
		Code: `date.Calendar(c, date.CalendarOptions{
    Label: "Deployment day", Today: today,
    Min: min, Max: max, DisabledDates: isWeekend,
}, &cal)`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Names the day grid for assistive technology."},
			{"Today", "time.Time", "context clock", "The caller's current day."},
			{"Min, Max", "time.Time", "zero", "Inclusive selectable date bounds."},
			{"DisabledDates", "func(time.Time) bool", "nil", "Rejects individual days."},
			{"FirstWeekday", "time.Weekday", "time.Sunday", "First column of the week."},
			{"Disabled", "bool", "false", "Makes the whole calendar read-only."},
		},
		Tokens: []string{"Accent", "AccentHover", "AccentSoft", "AccentSoftForeground", "Default"},
	})

	RegisterComponent(ComponentPage{
		ID: "rangecalendar", Title: "RangeCalendar", Group: "Dates", Order: 20,
		Description: "Two-endpoint day selection with an anchored preview and swapped endpoints.",
		Examples: []ComponentExample{
			{Title: "Bounded weekdays", Description: "Pick a start day, then an end day inside the bounds.",
				Demo: func(c *ui.Context, s *State) { rangeBoundedExample(c, datesComponentState(c, s)) }},
			{Title: "Monday first", Description: "Both grids start their week on Monday.",
				Demo: func(c *ui.Context, s *State) { rangeMondayExample(c, datesComponentState(c, s)) }},
			{Title: "Open bounds", Description: "Any two days can form the range.",
				Demo: func(c *ui.Context, s *State) { rangeFreeExample(c, datesComponentState(c, s)) }},
		},
		Code: `date.RangeCalendar(c, date.RangeCalendarOptions{
    Label: "Release window", Today: today,
    Min: min, Max: max, DisabledDates: isWeekend,
}, &rangeState)`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Names both calendar grids for assistive technology."},
			{"Today", "time.Time", "context clock", "The caller's current day."},
			{"Min, Max", "time.Time", "zero", "Inclusive selectable date bounds."},
			{"DisabledDates", "func(time.Time) bool", "nil", "Rejects individual days."},
			{"FirstWeekday", "time.Weekday", "time.Sunday", "First column of the week."},
			{"OnChange", "func(date.DateRange)", "nil", "Reports a completed, ordered range."},
		},
		Tokens: []string{"Accent", "AccentSoft", "AccentSoftForeground", "Default", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "datefield", Title: "DateField", Group: "Dates", Order: 30,
		Description: "An ISO 8601 date editor that commits complete, in-bounds text as it is typed.",
		Examples: []ComponentExample{
			{Title: "Bounded dates", Description: "Complete in-bounds weekdays commit as they are typed.",
				Demo: func(c *ui.Context, s *State) { dateFieldBoundedExample(c, datesComponentState(c, s)) }},
			{Title: "Unbounded", Description: "Any complete ISO date is accepted.",
				Demo: func(c *ui.Context, s *State) { dateFieldFreeExample(c, datesComponentState(c, s)) }},
			{Title: "Disabled", Description: "A read-only field ignores text input.",
				Demo: func(c *ui.Context, s *State) { dateFieldDisabledExample(c, datesComponentState(c, s)) }},
		},
		Code: `date.DateField(c, date.DateFieldOptions{
    Label: "Launch date", Description: "YYYY-MM-DD, weekdays only",
    Value: &launch, Min: min, Max: max,
    DisabledDates: isWeekend, Width: 200,
})`,
		Params: []ComponentParam{
			{"Label, Description", "string", "\"\"", "Visible field labels."},
			{"Value", "*time.Time", "nil", "Caller-owned committed day."},
			{"Min, Max", "time.Time", "zero", "Inclusive committed bounds."},
			{"DisabledDates", "func(time.Time) bool", "nil", "Rejects individual days."},
			{"Width", "float32", "160", "Field width in DIPs before scaling."},
			{"Disabled", "bool", "false", "Makes the field read-only."},
		},
		Tokens: []string{"FieldBackground", "FieldHover", "FieldFocus", "FieldBorder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "timefield", Title: "TimeField", Group: "Dates", Order: 40,
		Description: "An HH:MM or hh:mm AM/PM editor that commits a complete time as it is typed.",
		Examples: []ComponentExample{
			{Title: "24-hour", Description: "An HH:MM editor that commits a complete time.",
				Demo: func(c *ui.Context, s *State) { timeField24Example(c, datesComponentState(c, s)) }},
			{Title: "12-hour", Description: "An hh:mm AM/PM editor over the same day.",
				Demo: func(c *ui.Context, s *State) { timeField12Example(c, datesComponentState(c, s)) }},
			{Title: "Disabled", Description: "A read-only time field keeps its value.",
				Demo: func(c *ui.Context, s *State) { timeFieldDisabledExample(c, datesComponentState(c, s)) }},
		},
		Code: `date.TimeField(c, date.TimeFieldOptions{
    Label: "Window start", Description: "24-hour clock",
    Hour12: false, Width: 140,
}, &start)`,
		Params: []ComponentParam{
			{"Label, Description", "string", "\"\"", "Visible field labels."},
			{"Hour12", "bool", "false", "Switches to hh:mm AM/PM editing."},
			{"Width", "float32", "104", "Field width in DIPs before scaling."},
			{"Placeholder", "string", "\"HH:MM\"", "Text shown while the field is empty."},
			{"Disabled", "bool", "false", "Makes the field read-only."},
			{"OnChange", "func(time.Time)", "nil", "Reports a committed time."},
		},
		Tokens: []string{"FieldBackground", "FieldHover", "FieldFocus", "FieldBorder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "datepicker", Title: "DatePicker", Group: "Dates", Order: 50,
		Description: "A bounded date field whose trigger opens a popup calendar.",
		Examples: []ComponentExample{
			{Title: "Bounded picker", Description: "The popup bounds days and disables weekends.",
				Demo: func(c *ui.Context, s *State) { datePickerBoundedExample(c, datesComponentState(c, s)) }},
			{Title: "Unbounded picker", Description: "The popup accepts any day.",
				Demo: func(c *ui.Context, s *State) { datePickerFreeExample(c, datesComponentState(c, s)) }},
			{Title: "Disabled", Description: "A read-only picker cannot open its popup.",
				Demo: func(c *ui.Context, s *State) { datePickerDisabledExample(c, datesComponentState(c, s)) }},
		},
		Code: `date.DatePicker(c, date.DatePickerOptions{
    Label: "Deploy date", Value: &deploy,
    State: &cal, Open: &open, Min: min, Max: max,
    Width: 220, OpenLabel: "Open deploy calendar",
})`,
		Params: []ComponentParam{
			{"Label, Description", "string", "\"\"", "Visible field labels."},
			{"Value", "*time.Time", "nil", "Caller-owned committed day."},
			{"State", "*date.CalendarState", "nil", "Visible month, cursor and year picker."},
			{"Open", "*bool", "nil", "Caller-owned popup state."},
			{"Min, Max", "time.Time", "zero", "Inclusive selectable bounds."},
			{"Width", "float32", "200", "Field width in DIPs before scaling."},
		},
		Tokens: []string{"FieldBackground", "FieldFocus", "Overlay", "AccentSoft", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "daterangepicker", Title: "DateRangePicker", Group: "Dates", Order: 60,
		Description: "Two editable endpoints plus a popup range calendar over one shared value.",
		Examples: []ComponentExample{
			{Title: "Bounded range", Description: "Two editable endpoints over one bounded popup range calendar.",
				Demo: func(c *ui.Context, s *State) { rangePickerBoundedExample(c, datesComponentState(c, s)) }},
			{Title: "Unbounded range", Description: "Both endpoints accept any day.",
				Demo: func(c *ui.Context, s *State) { rangePickerFreeExample(c, datesComponentState(c, s)) }},
			{Title: "Disabled", Description: "A read-only picker keeps its ordered range.",
				Demo: func(c *ui.Context, s *State) { rangePickerDisabledExample(c, datesComponentState(c, s)) }},
		},
		Code: `date.DateRangePicker(c, date.DateRangePickerOptions{
    Label: "Release window", Value: &window,
    State: &rangeState, Open: &open, Min: min, Max: max,
    Width: 264, OpenLabel: "Open release calendar",
})`,
		Params: []ComponentParam{
			{"Label, Description", "string", "\"\"", "Visible field labels."},
			{"Value", "*date.DateRange", "nil", "Caller-owned ordered range."},
			{"State", "*date.RangeCalendarState", "nil", "Range selection and visible month."},
			{"Open", "*bool", "nil", "Caller-owned popup state."},
			{"Min, Max", "time.Time", "zero", "Inclusive selectable bounds."},
			{"Width", "float32", "264", "Field width in DIPs before scaling."},
		},
		Tokens: []string{"FieldBackground", "FieldFocus", "Overlay", "AccentSoft", "Border"},
	})
}
