// Package date provides COSS-inspired Gregorian date and time controls:
// Calendar, RangeCalendar, DateField, DatePicker, DateRangePicker and
// TimeField. The COSS registry's calendar and input/popover compositions are
// visual references; these controls use MyGo input, focus and popup behavior.
//
// Values and calendar view state are caller-owned when supplied. Nil state or
// value pointers use component-local state. Dates are time.Time values, and
// comparisons use only the calendar date; TimeField uses the time of day.
// A completed range is reported with its endpoints ordered.
//
// Calendars show four to six weeks, with outside days completing each week,
// at 36 DIP per day at the default text scale. The month heading is static by
// default; ShowYearPicker enables year selection. Selected days use Primary,
// range interiors and hover use Accent, and a dot marks today. A visible
// keyboard cursor receives the native focus ring. The year/day grid swap and
// popup fade respect reduced motion.
//
// # Limits
//
//   - Only Gregorian arithmetic and one visible month are supported. There is
//     no multiple-date selection, alternate calendar, week-number column,
//     multi-month view, or separate month/year dropdown as in COSS examples.
//   - RangeCalendar starts a pending range on the first activation and commits
//     on the second. It does not preview a range on hover or exclude disabled
//     dates from a committed range; disabled days cannot be chosen as endpoints
//     and interrupt the painted band.
//   - The calendar uses left-to-right day and keyboard order. FirstWeekday
//     changes the first column, but RTL mirroring is not implemented.
//   - The year picker is a native three-column grid, not a COSS dropdown. With
//     no bounds it offers twelve years around the visible year; with both
//     bounds it lists their years and scrolls when needed. Large bounded spans
//     build many cells. Grid fades are animated, but blur is not.
//   - Date fields edit ISO YYYY-MM-DD text; TimeField edits HH:MM or hh:mm AM/PM.
//     They do not provide locale-specific date parsing or native segmented
//     spinners. Pickers use editable fields with an icon trigger, rather than
//     every button-based composition in the COSS examples.
//   - MyGo dims disabled elements by 50%; COSS uses different opacity values.
//     Custom labels can localize calendar text and accessible day names.
package date
