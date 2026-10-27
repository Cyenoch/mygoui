package date

import (
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DatePickerOptions configures a date field with a popup calendar. Label names
// both the field and the calendar grid.
type DatePickerOptions struct {
	Label, Description, Error string
	Required, Invalid         bool
	Disabled, FullWidth       bool

	Value *time.Time
	// Width overrides the 200 DIP default. FullWidth wins over Width.
	Width float32

	Min, Max      time.Time
	DisabledDates func(time.Time) bool
	FirstWeekday  time.Weekday
	Today         time.Time
	Placeholder   string

	// ShowYearPicker makes the calendar's month heading open a year selector.
	ShowYearPicker bool
	// HideYearPicker overrides ShowYearPicker.
	//
	// Deprecated: leave ShowYearPicker false instead.
	HideYearPicker                        bool
	PrevLabel, NextLabel, YearPickerLabel string
	MonthLabel                            func(time.Time) string
	WeekdayLabel                          func(time.Weekday) string
	DayLabel                              func(time.Time) string

	// State holds the calendar's visible month, cursor and year picker. The
	// picker keeps its Value in step with Value.
	State *CalendarState
	// Open is the popup state. Nil keeps it in the component.
	Open *bool
	// OpenLabel names the calendar trigger. Empty uses "Open calendar".
	OpenLabel string
	// Placement places the popup. The zero value uses bottom-center, as COSS
	// does.
	Placement mygoui.Placement

	Handle       *ui.Handle
	OnChange     func(time.Time)
	OnOpenChange func(bool)
}

// DatePicker renders a date field with a COSS-style popup calendar. Choosing a day
// commits it to Value, reports OnChange and closes the popup.
func DatePicker(c *ui.Context, o DatePickerOptions) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	root := ui.Box(c).MinWidth(0)
	if o.FullWidth {
		root.FillWidth()
	}
	// Keep the overlay layer present while closed so the panel's fade in and
	// out has a parent in both frames.
	ui.Overlay(c, func() {})
	value := o.Value
	if value == nil {
		value = ui.Local(root, datePickerValueKey{}, func() time.Time { return time.Time{} })
	}
	open := o.Open
	if open == nil {
		open = ui.Local(root, datePickerOpenKey{}, func() bool { return false })
	}
	state := o.State
	if state == nil {
		state = ui.Local(root, datePickerStateKey{}, func() CalendarState { return CalendarState{} })
	}
	refocus := ui.Local(root, datePickerFocusKey{}, func() bool { return false })
	syncCalendar(state, *value)

	width := o.Width
	if width == 0 {
		width = 200
	}
	width = theme.Scale(c, width)
	spec := dateFieldSpec{
		field: fieldSpec{
			label:       o.Label,
			description: o.Description,
			error:       o.Error,
			required:    o.Required,
			invalid:     o.Invalid,
			disabled:    o.Disabled,
			fullWidth:   o.FullWidth,
		},
		value:         value,
		min:           o.Min,
		max:           o.Max,
		disabledDates: o.DisabledDates,
		placeholder:   o.Placeholder,
		loc:           locationOf(*value, c),
		width:         width,
		onChange:      o.OnChange,
	}
	buf := newDateBuffer(root, spec)

	var field, group, trigger ui.Element
	closePopup := func(notify bool) {
		if *open {
			*open = false
			if notify && o.OnOpenChange != nil {
				o.OnOpenChange(false)
			}
		}
	}

	openLabel := o.OpenLabel
	if openLabel == "" {
		openLabel = "Open calendar"
	}

	root.Children(func() {
		spec.suffix = func() {
			trigger = triggerButton(c, openLabel, func() {
				if *open {
					return
				}
				*open = true
				if o.OnOpenChange != nil {
					o.OnOpenChange(true)
				}
			})
		}
		field, group = dateField(c, spec, buf)
		if *refocus && trigger.Valid() {
			trigger.Focus()
			*refocus = false
		}

		if group.Valid() {
			panel := ui.PopoverBase(c, group, open, func(panel ui.Element) {
				box := popoverPanel(c, panel, group, o.Placement)
				box.Children(func() {
					Calendar(c, CalendarOptions{
						Label:           o.Label,
						Today:           o.Today,
						Min:             o.Min,
						Max:             o.Max,
						DisabledDates:   o.DisabledDates,
						FirstWeekday:    o.FirstWeekday,
						ShowYearPicker:  o.ShowYearPicker,
						HideYearPicker:  o.HideYearPicker,
						PrevLabel:       o.PrevLabel,
						NextLabel:       o.NextLabel,
						YearPickerLabel: o.YearPickerLabel,
						MonthLabel:      o.MonthLabel,
						WeekdayLabel:    o.WeekdayLabel,
						DayLabel:        o.DayLabel,
						OnChange: func(d time.Time) {
							*value = d
							if o.OnChange != nil {
								o.OnChange(d)
							}
							closePopup(true)
							*refocus = true
						},
						autofocus: true,
					}, state)
				})
			})
			if panel.Valid() && panel.Dismissed() {
				*refocus = true
				if o.OnOpenChange != nil {
					o.OnOpenChange(false)
				}
			}
		}
	})

	if o.Handle != nil {
		field.Bind(o.Handle)
	}
	return root
}

type datePickerValueKey struct{}
type datePickerOpenKey struct{}
type datePickerStateKey struct{}
type datePickerFocusKey struct{}

// syncCalendar makes the calendar show the field's value without discarding a
// month the user navigated to for the same value.
func syncCalendar(state *CalendarState, value time.Time) {
	if sameDay(value, state.Value) {
		return
	}
	state.Value = dayStart(value)
	if !value.IsZero() {
		state.Cursor = dayStart(value)
		state.Month = monthStart(value)
	}
}

// triggerButton is the calendar icon that opens a picker's popup.
func triggerButton(c *ui.Context, label string, onClick func()) ui.Element {
	t := theme.Colors(c)
	size := theme.Scale(c, 24)
	e := ui.ButtonBase(c).Size(size, size).Center().Shrink(0).
		Radius(theme.RadiusScale().Field).FocusRing(true).Label(label).
		Cursor(ui.CursorPointer)
	e = style.Background(e, theme.RadiusScale().Field, ui.Transparent, t.Accent, t.Accent)
	if onClick != nil {
		e = e.OnClick(onClick)
	}
	e.Children(func() {
		ui.Icon(c, lucide.Calendar()).Size(theme.Scale(c, 16), theme.Scale(c, 16)).
			TextColor(t.Foreground)
	})
	return e
}

// popoverPanel gives the calendar a bordered COSS popover and an arrival fade.
func popoverPanel(c *ui.Context, panel, anchor ui.Element, place mygoui.Placement) ui.Element {
	t := theme.Colors(c)
	p := panel.Padding(theme.Space8).Radius(theme.RadiusScale().XL).Clip().Background(t.Popover).Border(theme.Hairline, t.Border)
	p = style.Shadow(p, c, theme.OverlayShadows(c.Theme().Dark))
	p = anchored.Apply(p, anchor, place, mygoui.SideBottom)
	if !theme.ReduceMotion(c) {
		fade := ui.Motion{Opacity: 0}
		p = p.Transition(ui.ElementTransition{
			Duration: 150 * time.Millisecond, Ease: ui.EaseOut, Enter: &fade, Exit: &fade,
		})
	}
	return p
}
