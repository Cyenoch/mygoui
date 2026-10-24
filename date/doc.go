// Package date implements HeroUI v3's Gregorian date and time controls:
// Calendar, RangeCalendar, DateField, DatePicker, DateRangePicker and
// TimeField.
//
// Every control is caller-owned. Dates are plain time.Time values; only the
// calendar date (year, month, day) matters for Calendar, RangeCalendar,
// DateField and the pickers, and only the time of day matters for TimeField.
// The view state a Calendar keeps between frames — visible month, keyboard
// cursor and year-picker visibility — is a CalendarViewState the caller passes
// in, or nil for a component-local copy. Values are written to the pointer in
// the options, and OnChange reports the same value without taking ownership of
// it.
//
// The controls follow HeroUI's desktop metrics: a 252 DIP calendar, 36 DIP
// days, a 36 DIP field with the 12 DIP field radius, 14 DIP control text and a
// 2 DIP focus ring. Hover, pressed and focus paint is read from live pointer
// and focus state while drawing. Disabled days use 50% opacity, a line-through
// and never activate.
//
// # Limits
//
//   - Only the Gregorian calendar is implemented. There is no alternate
//     calendar arithmetic, era display or non-Gregorian month grid; FirstWeekday
//     controls the week's first column only.
//   - RangeCalendar paints start, end and in-between days. HeroUI's hover
//     preview across a pending range is not implemented, because a single cell
//     cannot know the whole grid's pointer position without an extra frame of
//     state; selection still updates on click as upstream does.
//   - The calendar-year-picker crossfade and blur transition is not animated;
//     the platform's reduced-motion preference therefore has nothing to
//     suppress. The year grid keeps its 3-column, 32 DIP cell layout.
//   - Calendar keyboard navigation assumes left-to-right layout. RTL day
//     order (arrows and week columns mirrored) is not implemented.
//   - Date and time fields edit ISO 8601 text (YYYY-MM-DD, HH:MM, or hh:mm AM
//     when Hour12 is set) rather than per-segment spinners. Parsing, bounds and
//     disabled-date checks are applied as the text becomes complete.
//   - The year grid shows one page of twelve years centred on the visible
//     year, or exactly the caller's Min/Max years when both are set. Bounds
//     wider than the picker's height scroll; keep custom bounds within a few
//     centuries so the grid stays small.
package date
