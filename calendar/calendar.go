// Package calendar provides Gregorian civil-date selection and editable pickers.
// The pinned MyGo i18n formatter explicitly supports Gregorian dates only;
// Persian, Hijri and Jalali arithmetic is not provided by its public API.
package calendar

import (
	"fmt"
	"slices"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Date is a civil day, independent of elapsed hours or daylight-saving changes.
// Its zero value represents no date. Valid years are 1 through 9999.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func FromTime(t time.Time) Date { y, m, d := t.Date(); return Date{y, m, d} }
func (d Date) Valid() bool {
	if d.Year < 1 || d.Year > 9999 || d.Month < time.January || d.Month > time.December || d.Day < 1 || d.Day > 31 {
		return false
	}
	return FromTime(time.Date(d.Year, d.Month, d.Day, 12, 0, 0, 0, time.UTC)) == d
}

// In returns civil noon in location, avoiding ordinary midnight DST transitions.
// Invalid dates return the zero time. The caller's location is never converted to UTC.
func (d Date) In(location *time.Location) time.Time {
	if !d.Valid() {
		return time.Time{}
	}
	if location == nil {
		location = time.Local
	}
	return time.Date(d.Year, d.Month, d.Day, 12, 0, 0, 0, location)
}
func (d Date) String() string {
	if !d.Valid() {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// Parse accepts the explicit ISO civil-date input contract, including leap-day validation.
func Parse(text string) (Date, error) {
	t, err := time.Parse(time.DateOnly, text)
	if err != nil {
		return Date{}, err
	}
	d := FromTime(t)
	if !d.Valid() {
		return Date{}, fmt.Errorf("invalid civil date %q", text)
	}
	return d, nil
}
func compare(a, b Date) int {
	if a.Year != b.Year {
		return a.Year - b.Year
	}
	if a.Month != b.Month {
		return int(a.Month - b.Month)
	}
	return a.Day - b.Day
}
func addDays(d Date, n int) Date { return FromTime(d.In(time.UTC).AddDate(0, 0, n)) }
func addMonths(d Date, n int) Date {
	first := time.Date(d.Year, d.Month+time.Month(n), 1, 12, 0, 0, 0, time.UTC)
	last := time.Date(first.Year(), first.Month()+1, 0, 12, 0, 0, 0, time.UTC).Day()
	return Date{first.Year(), first.Month(), min(d.Day, last)}
}

type Range struct{ Start, End Date }
type Mode uint8

const (
	Single Mode = iota
	Multiple
	RangeMode
)

// State belongs to one calendar in one window. It stores no build-scoped values.
type State struct {
	Date          Date
	Dates         []Date
	Range         Range
	Month, Cursor Date
	Handle        ui.Handle
}
type Options struct {
	Mode               Mode
	Label, Locale      string
	Location           *time.Location
	FirstWeekday       *time.Weekday
	Min, Max           Date
	Disabled           func(Date) bool
	ReadOnly           bool
	MonthChoice        bool
	StartYear, EndYear int
	OnChange           func()
}

func (o Options) allowed(d Date) bool {
	return d.Valid() && (!o.Min.Valid() || compare(d, o.Min) >= 0) && (!o.Max.Valid() || compare(d, o.Max) <= 0) && (o.Disabled == nil || !o.Disabled(d))
}

// Select applies the same constraints for pointer, keyboard, text and presets.
// A range crossing a disabled day is refused; reversed endpoints are ordered.
func (s *State) Select(d Date, o Options) bool {
	if s == nil || o.ReadOnly || !o.allowed(d) {
		return false
	}
	switch o.Mode {
	case Multiple:
		if i := slices.Index(s.Dates, d); i >= 0 {
			s.Dates = slices.Delete(s.Dates, i, i+1)
		} else {
			s.Dates = append(s.Dates, d)
			slices.SortFunc(s.Dates, compare)
		}
	case RangeMode:
		if !s.Range.Start.Valid() || s.Range.End.Valid() {
			s.Range = Range{Start: d}
		} else {
			r := Range{s.Range.Start, d}
			if compare(r.Start, r.End) > 0 {
				r.Start, r.End = r.End, r.Start
			}
			if !o.allowedRange(r) {
				return false
			}
			s.Range = r
		}
	default:
		if s.Date == d {
			return false
		}
		s.Date = d
	}
	s.Cursor = d
	s.Month = Date{d.Year, d.Month, 1}
	if o.OnChange != nil {
		o.OnChange()
	}
	return true
}
func (o Options) allowedRange(r Range) bool {
	if !o.allowed(r.Start) || !o.allowed(r.End) || compare(r.Start, r.End) > 0 {
		return false
	}
	if o.Disabled != nil {
		for d := r.Start; compare(d, r.End) < 0; d = addDays(d, 1) {
			if o.Disabled(d) {
				return false
			}
		}
	}
	return true
}
func (s *State) selected(d Date, o Options) bool {
	switch o.Mode {
	case Multiple:
		return slices.Contains(s.Dates, d)
	case RangeMode:
		return s.Range.Start == d || (s.Range.End.Valid() && compare(d, s.Range.Start) >= 0 && compare(d, s.Range.End) <= 0)
	default:
		return s.Date == d
	}
}
func (s *State) initialize(now time.Time, o Options) {
	if !s.Cursor.Valid() {
		s.Cursor = s.Date
		if o.Mode == Multiple && len(s.Dates) > 0 {
			s.Cursor = s.Dates[0]
		}
		if o.Mode == RangeMode {
			s.Cursor = s.Range.Start
		}
		if !s.Cursor.Valid() {
			s.Cursor = FromTime(now)
		}
	}
	if !s.Month.Valid() {
		s.Month = Date{s.Cursor.Year, s.Cursor.Month, 1}
	}
	s.Month.Day = 1
}

// Calendar extends the maintained widget's documented CLDR/grid/key behavior
// with civil values, multiple/range selection and constraints, using public native elements.
func Calendar(c *ui.Context, s *State, o Options) ui.Element {
	if s == nil {
		return ui.Element{}
	}
	now := c.Now()
	if o.Location != nil {
		now = now.In(o.Location)
	}
	s.initialize(now, o)
	locale := o.Locale
	if locale == "" {
		locale = c.Locale()
	}
	locale = c.Region().FormatLocale(locale)
	first := c.Region().FirstWeekday
	if o.FirstWeekday != nil && *o.FirstWeekday >= time.Sunday && *o.FirstWeekday <= time.Saturday {
		first = *o.FirstWeekday
	}
	if o.Label == "" {
		o.Label = locales.ComponentCalendar()
	}
	p := theme.Colors(c)
	cell := max(float32(32), theme.Scale(c, 28)+theme.Space4)
	title := i18n.FormatDate(locale, s.Month.In(o.Location), i18n.DateOptions{Skeleton: "yMMMM"})
	root := ui.Column(c).Bind(&s.Handle).Focusable().Role(ui.RoleTable).Label(o.Label).Description(title).
		Value(i18n.FormatDate(locale, s.Cursor.In(o.Location), i18n.DateOptions{Date: i18n.FullStyle})).Gap(theme.Space4).
		Padding(theme.Space12).Radius(theme.Radii().LG).Border(1, p.Border).Background(p.Background).TextColor(p.Foreground).MinWidth(0)
	move := func(d Date) {
		if !d.Valid() {
			return
		}
		step := 1
		if compare(d, s.Cursor) < 0 {
			step = -1
		}
		for n := 0; n < 36600 && !o.allowed(d); n++ {
			if (o.Min.Valid() && compare(d, o.Min) < 0) || (o.Max.Valid() && compare(d, o.Max) > 0) {
				return
			}
			d = addDays(d, step)
			if !d.Valid() {
				return
			}
		}
		if !o.allowed(d) {
			return
		}
		s.Cursor = d
		s.Month = Date{d.Year, d.Month, 1}
	}
	back, forth := ui.KeyLeft, ui.KeyRight
	if c.Direction() == ui.RightToLeft {
		back, forth = forth, back
	}
	root.OnShortcut(0, back, func() { move(addDays(s.Cursor, -1)) }).OnShortcut(0, forth, func() { move(addDays(s.Cursor, 1)) }).
		OnShortcut(0, ui.KeyUp, func() { move(addDays(s.Cursor, -7)) }).OnShortcut(0, ui.KeyDown, func() { move(addDays(s.Cursor, 7)) }).
		OnShortcut(0, ui.KeyPageUp, func() { move(addMonths(s.Cursor, -1)) }).OnShortcut(0, ui.KeyPageDown, func() { move(addMonths(s.Cursor, 1)) }).
		OnShortcut(0, ui.KeyHome, func() { move(Date{s.Cursor.Year, s.Cursor.Month, 1}) }).
		OnShortcut(0, ui.KeyEnd, func() { move(FromTime(time.Date(s.Cursor.Year, s.Cursor.Month+1, 0, 12, 0, 0, 0, time.UTC))) }).
		OnShortcut(0, ui.KeyEnter, func() { s.Select(s.Cursor, o) }).OnShortcut(0, ui.KeySpace, func() { s.Select(s.Cursor, o) })
	root.Children(func() {
		ui.Row(c).Gap(theme.Space4).Children(func() {
			previous, next := addMonths(s.Month, -1), addMonths(s.Month, 1)
			mygoui.Button(c.Key("previous"), mygoui.ButtonOptions{Label: locales.ComponentPreviousMonth(), Icon: lucide.ChevronLeft(), Size: mygoui.ButtonIconSM, Variant: mygoui.ButtonGhost, Disabled: !previous.Valid() || (o.Min.Valid() && compare(addMonths(previous, 1), Date{o.Min.Year, o.Min.Month, 1}) <= 0), OnClick: func() { s.Month = previous; s.Cursor = previous }})
			if o.MonthChoice {
				monthChoices(c, s, o, locale)
			} else {
				ui.Text(c, title).Grow(1).TextAlign(ui.Center).SingleLine().FontSize(theme.Scale(c, 14))
			}
			mygoui.Button(c.Key("next"), mygoui.ButtonOptions{Label: locales.ComponentNextMonth(), Icon: lucide.ChevronRight(), Size: mygoui.ButtonIconSM, Variant: mygoui.ButtonGhost, Disabled: !next.Valid() || (o.Max.Valid() && compare(next, Date{o.Max.Year, o.Max.Month, 1}) > 0), OnClick: func() { s.Month = next; s.Cursor = next }})
		})
		names := i18n.Weekdays(locale, i18n.Short, true)
		ui.Row(c).Children(func() {
			for i := range 7 {
				ui.Text(c, names[(int(first)+i)%7]).Width(cell).TextAlign(ui.Center).Role(ui.RoleColumnHeader).FontSize(theme.Scale(c, 12)).TextColor(p.MutedForeground)
			}
		})
		start := addDays(s.Month, -((int(s.Month.In(time.UTC).Weekday()) - int(first) + 7) % 7))
		for week := range 6 {
			ui.Row(c.Key(week)).Role(ui.RoleRow).Children(func() {
				for day := range 7 {
					d := addDays(start, week*7+day)
					chosen := s.selected(d, o)
					enabled := !o.ReadOnly && o.allowed(d)
					label := i18n.FormatDate(locale, d.In(o.Location), i18n.DateOptions{Date: i18n.FullStyle})
					b := ui.Box(c.Key(d)).Size(cell, cell).Center().Radius(theme.Radii().MD).Role(ui.RoleButton).Label(label).Selected(chosen).Disabled(!enabled)
					color := p.Foreground
					if d.Month != s.Month.Month {
						color = p.MutedForeground
					}
					if !enabled {
						color = p.DisabledForeground
					}
					if chosen {
						b.Background(p.Primary)
						color = p.PrimaryForeground
					} else if b.Hovered() {
						b.Background(p.Accent)
					}
					if d == s.Cursor {
						b.Border(1, p.Ring)
						root.ActiveDescendant(b)
					} else if d == FromTime(now) {
						b.Border(1, p.Border)
					}
					b.TextColor(color).OnClick(func() { s.Handle.Focus(c); s.Select(d, o) }).Children(func() {
						ui.Text(c, i18n.FormatDate(locale, d.In(o.Location), i18n.DateOptions{Skeleton: "d"})).FontSize(theme.Scale(c, 14)).SingleLine()
					})
				}
			})
		}
	})
	return root
}
func monthChoices(c *ui.Context, s *State, o Options, locale string) {
	label := i18n.FormatDate(locale, s.Month.In(o.Location), i18n.DateOptions{Skeleton: "MMMM"})
	ui.Button(c, label).Label(locales.ComponentChooseMonth()).Grow(1).Menu(func(menu *ui.Menu) {
		for m := time.January; m <= time.December; m++ {
			date := Date{s.Month.Year, m, 1}
			end := addDays(addMonths(date, 1), -1)
			disabled := (o.Min.Valid() && compare(end, o.Min) < 0) || (o.Max.Valid() && compare(date, o.Max) > 0)
			if menu.Item(i18n.FormatDate(locale, date.In(o.Location), i18n.DateOptions{Skeleton: "MMMM"})).Disabled(disabled).Checked(s.Month.Month == m).Chosen() {
				s.Month = date
				s.Cursor = date
			}
		}
	})
	low, high := o.StartYear, o.EndYear
	if low < 1 {
		low = 1900
	}
	if high < low {
		high = max(low, c.Now().Year()+10)
	}
	high = min(high, 9999)
	ui.Button(c, i18n.FormatDate(locale, s.Month.In(o.Location), i18n.DateOptions{Skeleton: "y"})).Label(locales.ComponentChooseYear()).Menu(func(menu *ui.Menu) {
		for y := low; y <= high; y++ {
			disabled := (o.Min.Valid() && y < o.Min.Year) || (o.Max.Valid() && y > o.Max.Year)
			if menu.Item(i18n.FormatDate(locale, Date{y, time.January, 1}.In(o.Location), i18n.DateOptions{Skeleton: "y"})).Disabled(disabled).Checked(s.Month.Year == y).Chosen() {
				s.Month = Date{y, s.Month.Month, 1}
				s.Cursor = s.Month
			}
		}
	})
}
