package date

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DateRangePickerOptions configures a range field with a popup range calendar.
// Label names the field and both calendar grids.
type DateRangePickerOptions struct {
	Label, Description, Error string
	Required, Invalid         bool
	Disabled, FullWidth       bool

	Value *DateRange
	// Width overrides the 264 DIP default. FullWidth wins over Width.
	Width float32

	Min, Max      time.Time
	DisabledDates func(time.Time) bool
	FirstWeekday  time.Weekday
	Today         time.Time
	// StartPlaceholder and EndPlaceholder default to "YYYY-MM-DD".
	StartPlaceholder, EndPlaceholder string
	// StartLabel and EndLabel name the two editors. They default to "Start
	// date" and "End date".
	StartLabel, EndLabel string

	HideYearPicker                        bool
	PrevLabel, NextLabel, YearPickerLabel string
	MonthLabel                            func(time.Time) string
	WeekdayLabel                          func(time.Weekday) string
	DayLabel                              func(time.Time) string

	// State holds the range calendar's selection and view. The picker syncs
	// State.Value from Value while no endpoint is pending.
	State *RangeCalendarState
	// Open is the popup state. Nil keeps it in the component.
	Open      *bool
	OpenLabel string
	Placement mygoui.Placement

	Handle       *ui.Handle
	OnChange     func(DateRange)
	OnOpenChange func(bool)
}

// DateRangePicker renders two ISO date editors and a popup range calendar.
// Completing the range commits it to Value, reports OnChange and closes the
// popup. Editing the text commits an ordered range once both endpoints are
// valid.
func DateRangePicker(c *ui.Context, o DateRangePickerOptions) ui.Element {
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
		value = ui.Local(root, rangePickerValueKey{}, func() DateRange { return DateRange{} })
	}
	open := o.Open
	if open == nil {
		open = ui.Local(root, rangePickerOpenKey{}, func() bool { return false })
	}
	state := o.State
	if state == nil {
		state = ui.Local(root, rangePickerStateKey{}, func() RangeCalendarState { return RangeCalendarState{} })
	}
	refocus := ui.Local(root, rangePickerFocusKey{}, func() bool { return false })
	if state.Anchor.IsZero() {
		syncRangeCalendar(state, *value)
	}

	width := o.Width
	if width == 0 {
		width = theme.Scale(c, 264)
	}
	var triggerHandle ui.Handle
	spec := rangeFieldSpec{
		field: fieldSpec{
			label:       o.Label,
			description: o.Description,
			error:       o.Error,
			required:    o.Required,
			invalid:     o.Invalid,
			disabled:    o.Disabled,
			fullWidth:   o.FullWidth,
		},
		value:            value,
		min:              o.Min,
		max:              o.Max,
		disabledDates:    o.DisabledDates,
		startPlaceholder: o.StartPlaceholder,
		endPlaceholder:   o.EndPlaceholder,
		startLabel:       o.StartLabel,
		endLabel:         o.EndLabel,
		loc:              locationOf(value.Start, c),
		width:            width,
		triggerFocus:     &triggerHandle,
		onChange:         o.OnChange,
	}
	buf := newRangeBuffer(root, spec)

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
			trigger.Bind(&triggerHandle)
		}
		field, group = rangeField(c, spec, buf)
		if *refocus && trigger.Valid() {
			trigger.Focus()
			*refocus = false
		}

		if group.Valid() {
			panel := ui.PopoverBase(c, group, open, func(panel ui.Element) {
				box := popoverPanel(c, panel, group, o.Placement)
				box.Children(func() {
					RangeCalendar(c, RangeCalendarOptions{
						Label:           o.Label,
						Today:           o.Today,
						Min:             o.Min,
						Max:             o.Max,
						DisabledDates:   o.DisabledDates,
						FirstWeekday:    o.FirstWeekday,
						HideYearPicker:  o.HideYearPicker,
						PrevLabel:       o.PrevLabel,
						NextLabel:       o.NextLabel,
						YearPickerLabel: o.YearPickerLabel,
						MonthLabel:      o.MonthLabel,
						WeekdayLabel:    o.WeekdayLabel,
						DayLabel:        o.DayLabel,
						autofocus:       true,
						OnChange: func(r DateRange) {
							*value = r
							if o.OnChange != nil {
								o.OnChange(r)
							}
							closePopup(true)
							*refocus = true
						},
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

type rangePickerValueKey struct{}
type rangePickerOpenKey struct{}
type rangePickerStateKey struct{}
type rangePickerFocusKey struct{}

func syncRangeCalendar(state *RangeCalendarState, value DateRange) {
	next := value.Ordered()
	if sameRange(next, state.Value) {
		return
	}
	state.Value = next
	if !next.Start.IsZero() {
		state.Cursor = dayStart(next.Start)
		state.Month = monthStart(next.Start)
	}
}

type rangeFieldSpec struct {
	field            fieldSpec
	value            *DateRange
	min, max         time.Time
	disabledDates    func(time.Time) bool
	startPlaceholder string
	endPlaceholder   string
	startLabel       string
	endLabel         string
	loc              *time.Location
	width            float32
	suffix           func()
	// triggerFocus is the popup trigger built by suffix, when there is one:
	// the group hides its focus face while that trigger holds the focus.
	triggerFocus *ui.Handle
	onChange     func(DateRange)
}

func rangeField(c *ui.Context, spec rangeFieldSpec, buf *rangeBuffer) (ui.Element, ui.Element) {
	buf.refresh(*spec.value)
	invalid := spec.field.invalid || buf.invalid()
	errorText := spec.field.error
	if errorText == "" && invalid {
		errorText = locale.M().Invalid
	}
	var group ui.Element
	field := formField(c, spec.field, errorText, func() {
		group = rangeGroup(c, spec, buf, invalid)
	})
	return field, group
}

func rangeGroup(c *ui.Context, spec rangeFieldSpec, buf *rangeBuffer, invalid bool) ui.Element {
	t := theme.Colors(c)
	startPlaceholder, endPlaceholder := spec.startPlaceholder, spec.endPlaceholder
	if startPlaceholder == "" {
		startPlaceholder = "YYYY-MM-DD"
	}
	if endPlaceholder == "" {
		endPlaceholder = "YYYY-MM-DD"
	}
	startLabel, endLabel := spec.startLabel, spec.endLabel
	if startLabel == "" {
		startLabel = "Start date"
	}
	if endLabel == "" {
		endLabel = "End date"
	}
	group := fieldGroup(c, invalid, spec.field.disabled, spec.triggerFocus)
	if spec.width > 0 && !spec.field.fullWidth {
		group = group.Width(spec.width)
	}
	if spec.field.fullWidth {
		group = group.FillWidth()
	}
	group.Children(func() {
		start := fieldInput(c, &buf.startText, startPlaceholder)
		start.Label(startLabel)
		start.OnChange(func() { buf.commit(spec) })

		style.Text(ui.Text(c, "–"), c, theme.SM()).
			TextColor(t.FieldPlaceholder).PaddingX(theme.Space4).Shrink(0)

		end := fieldInput(c, &buf.endText, endPlaceholder)
		end.Label(endLabel)
		end.OnChange(func() { buf.commit(spec) })

		if spec.suffix != nil {
			spec.suffix()
		}
	})
	return group
}

// rangeBuffer is the editable text of both endpoints and the range it was
// last synchronised with.
type rangeBuffer struct {
	startText, endText       string
	committed                DateRange
	invalidStart, invalidEnd bool
}

func (b *rangeBuffer) invalid() bool { return b.invalidStart || b.invalidEnd }

func newRangeBuffer(root ui.Element, spec rangeFieldSpec) *rangeBuffer {
	return ui.Local(root, rangeBufferKey{}, func() rangeBuffer {
		v := spec.value.Ordered()
		return rangeBuffer{
			startText: formatISO(v.Start),
			endText:   formatISO(v.End),
			committed: v,
		}
	})
}

func (b *rangeBuffer) refresh(value DateRange) {
	next := value.Ordered()
	if sameRange(next, b.committed) {
		return
	}
	b.startText = formatISO(next.Start)
	b.endText = formatISO(next.End)
	b.committed = next
	b.invalidStart, b.invalidEnd = false, false
}

// commit parses both endpoints after an edit. An incomplete range stays quiet;
// a complete one is ordered, stored and reported.
func (b *rangeBuffer) commit(spec rangeFieldSpec) {
	b.invalidStart, b.invalidEnd = false, false
	start, startOK := parseEndpoint(b.startText, spec.loc, spec.min, spec.max, spec.disabledDates)
	end, endOK := parseEndpoint(b.endText, spec.loc, spec.min, spec.max, spec.disabledDates)
	startPartial := b.startText != "" && len(b.startText) < len("2006-01-02")
	endPartial := b.endText != "" && len(b.endText) < len("2006-01-02")
	b.invalidStart = !startOK && !startPartial
	b.invalidEnd = !endOK && !endPartial
	if startPartial || endPartial {
		return
	}

	switch {
	case b.startText == "" && b.endText == "":
		if !b.committed.Empty() {
			b.committed = DateRange{}
			*spec.value = DateRange{}
			if spec.onChange != nil {
				spec.onChange(DateRange{})
			}
		}
		return
	case b.startText == "" || b.endText == "":
		// One endpoint is still empty: wait for the other.
		return
	case !startOK || !endOK:
		return
	}
	next := DateRange{Start: start, End: end}.Ordered()
	b.startText, b.endText = formatISO(next.Start), formatISO(next.End)
	if sameRange(next, b.committed) {
		return
	}
	b.committed = next
	*spec.value = next
	if spec.onChange != nil {
		spec.onChange(next)
	}
}

type rangeBufferKey struct{}

// parseEndpoint parses one ISO endpoint. Empty text is not an error; an
// impossible, out-of-bounds or disabled date is.
func parseEndpoint(text string, loc *time.Location, min, max time.Time, disabled func(time.Time) bool) (time.Time, bool) {
	if text == "" {
		return time.Time{}, true
	}
	if len(text) < len("2006-01-02") {
		return time.Time{}, false
	}
	d, ok := parseISO(text, loc)
	if !ok || outOfBounds(d, min, max) || (disabled != nil && disabled(d)) {
		return time.Time{}, false
	}
	return d, true
}

func sameRange(a, b DateRange) bool {
	return sameDay(a.Start, b.Start) && sameDay(a.End, b.End)
}
