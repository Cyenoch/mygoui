package date

import (
	"strconv"
	"sync"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// CalendarViewState is the caller-owned view state shared by Calendar and
// RangeCalendar: the visible month, the day the keyboard moves, and whether
// the year picker covers the day grid.
//
// A zero Month follows the selected value, then the context's today. A zero
// Cursor follows the selected value, then today, then the first day of the
// month. The widgets update these fields; callers may pre-set them to open a
// calendar at a particular month or day.
type CalendarViewState struct {
	Month    time.Time
	Cursor   time.Time
	YearOpen bool
}

// CalendarState is the caller-owned state of a Calendar.
type CalendarState struct {
	// Value is the chosen day. Its zero value selects nothing.
	Value time.Time
	CalendarViewState
}

// CalendarOptions configures a Calendar. Label names the day grid for
// assistive technology; HeroUI's calendar has no visible label of its own, so
// wrap it in a form.Field for a visible field label.
type CalendarOptions struct {
	Label  string
	Handle *ui.Handle

	// Today marks the caller's current day. Zero reads the context clock,
	// which headless tests control.
	Today time.Time
	// Min and Max bound selectable days, inclusive, by calendar date.
	Min, Max time.Time
	// DisabledDates rejects individual days. It runs for each visible day and
	// must not perform I/O.
	DisabledDates func(time.Time) bool
	// FirstWeekday is the first column of the week. The zero value is Sunday,
	// matching time.Weekday.
	FirstWeekday time.Weekday
	// Disabled makes the whole calendar read-only at 50% opacity.
	Disabled bool
	// HideYearPicker removes the heading's year selector and keeps the static
	// month heading.
	HideYearPicker bool

	// PrevLabel and NextLabel name the navigation buttons. Empty uses
	// locale.M().Previous and locale.M().Next.
	PrevLabel, NextLabel string
	// YearPickerLabel names the year grid and the heading's accessible
	// description. Empty uses "Choose year"; locale.M has no date-specific
	// string, so callers localize this option directly.
	YearPickerLabel string
	// MonthLabel, WeekdayLabel and DayLabel format the heading, weekday
	// headers and each cell's accessible name. Nil uses English formats
	// ("January 2006", "Mon", "Monday, January 2, 2006").
	MonthLabel   func(time.Time) string
	WeekdayLabel func(time.Weekday) string
	DayLabel     func(time.Time) string

	OnChange func(time.Time)

	// autofocus focuses the grid when it first appears, which is how the
	// pickers move the keyboard into a freshly opened popup.
	autofocus bool
}

// Calendar renders HeroUI's single-day calendar. The state carries selection
// and view position; a nil state uses a component-local copy.
func Calendar(c *ui.Context, o CalendarOptions, st *CalendarState) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	root := ui.Box(c).MinWidth(0).Shrink(0)
	state := st
	if state == nil {
		state = ui.Local(root, calendarStateKey{}, func() CalendarState { return CalendarState{} })
	}
	spec := o.spec()
	cal := buildCalendar(c, spec, &state.CalendarViewState, state.Value, cellSelection{
		selected: func(d time.Time) bool { return sameDay(d, state.Value) },
		activate: func(d time.Time) {
			state.Value = d
			spec.notify(d)
		},
	})
	if o.Handle != nil {
		cal.Bind(o.Handle)
	}
	return cal
}

func (o CalendarOptions) spec() calendarSpec {
	return calendarSpec{
		label:          o.Label,
		today:          o.Today,
		min:            o.Min,
		max:            o.Max,
		disabledDates:  o.DisabledDates,
		firstWeekday:   o.FirstWeekday,
		disabled:       o.Disabled,
		hideYearPicker: o.HideYearPicker,
		prevLabel:      o.PrevLabel,
		nextLabel:      o.NextLabel,
		yearLabel:      o.YearPickerLabel,
		monthLabel:     o.MonthLabel,
		weekdayLabel:   o.WeekdayLabel,
		dayLabel:       o.DayLabel,
		onChange:       o.OnChange,
		autofocus:      o.autofocus,
	}
}

type calendarStateKey struct{}
type calendarFocusKey struct{}
type yearScrollKey struct{}
type yearScrollInitKey struct{}
type yearFocusKey struct{}

// calendarSpec is the internal, shared configuration of both calendars.
type calendarSpec struct {
	label          string
	today          time.Time
	min, max       time.Time
	disabledDates  func(time.Time) bool
	firstWeekday   time.Weekday
	disabled       bool
	hideYearPicker bool
	autofocus      bool

	prevLabel, nextLabel, yearLabel string
	monthLabel                      func(time.Time) string
	weekdayLabel                    func(time.Weekday) string
	dayLabel                        func(time.Time) string

	onChange func(time.Time)
}

func (s calendarSpec) notify(d time.Time) {
	if s.onChange != nil {
		s.onChange(d)
	}
}

func (s calendarSpec) todayIn(c *ui.Context) time.Time {
	if !s.today.IsZero() {
		return dayStart(s.today)
	}
	return dayStart(c.Now())
}

func (s calendarSpec) prev() string {
	if s.prevLabel != "" {
		return s.prevLabel
	}
	return locale.M().Previous
}

func (s calendarSpec) next() string {
	if s.nextLabel != "" {
		return s.nextLabel
	}
	return locale.M().Next
}

func (s calendarSpec) years() string {
	if s.yearLabel != "" {
		return s.yearLabel
	}
	return "Choose year"
}

func (s calendarSpec) monthText(t time.Time) string {
	if s.monthLabel != nil {
		return s.monthLabel(t)
	}
	return t.Format("January 2006")
}

func (s calendarSpec) weekdayText(w time.Weekday) string {
	if s.weekdayLabel != nil {
		return s.weekdayLabel(w)
	}
	name := w.String()
	if len(name) > 3 {
		name = name[:3]
	}
	return name
}

func (s calendarSpec) dayText(t time.Time) string {
	if s.dayLabel != nil {
		return s.dayLabel(t)
	}
	return t.Format("Monday, January 2, 2006")
}

// disabledDay reports whether d cannot be chosen.
func (s calendarSpec) disabledDay(d time.Time) bool {
	if s.disabled {
		return true
	}
	if !s.min.IsZero() && compareDay(d, s.min) < 0 {
		return true
	}
	if !s.max.IsZero() && compareDay(d, s.max) > 0 {
		return true
	}
	return s.disabledDates != nil && s.disabledDates(d)
}

// cellSelection decides how a day looks and what choosing it does. A nil
// selected means no single-value highlight (range calendars use inRange,
// start and end instead).
type cellSelection struct {
	selected func(time.Time) bool
	inRange  func(time.Time) bool
	start    func(time.Time) bool
	end      func(time.Time) bool
	activate func(time.Time)
}

// buildCalendar builds the header and the day grid or year grid shared by
// Calendar and RangeCalendar.
func buildCalendar(c *ui.Context, spec calendarSpec, view *CalendarViewState, value time.Time, sel cellSelection) ui.Element {
	today := spec.todayIn(c)
	month := view.effectiveMonth(value, today)
	cursor := view.cursorIn(month, value, today)
	monthText := spec.monthText(month)

	width := theme.Scale(c, 252)
	cellW := width / 7
	cellH := cellW
	weekdayH := theme.Scale(c, 24)
	headerH := theme.Scale(c, 24)
	gap := theme.Scale(c, 4)
	dayGridH := weekdayH + gap + 6*cellH

	root := ui.Column(c).Width(width).MinWidth(width).Shrink(0).Gap(0)
	focusHeading := ui.Local(root, calendarFocusKey{}, func() bool { return false })

	root.Children(func() {
		buildHeader(c, spec, view, month, cursor, monthText, headerH, focusHeading)
		if view.YearOpen && !spec.hideYearPicker {
			buildYearGrid(c, spec, view, month, cursor, dayGridH, focusHeading)
		} else {
			buildDayGrid(c, spec, view, sel, month, cursor, today, cellW, cellH, weekdayH, gap)
		}
	})

	if spec.disabled {
		root = style.Disabled(root, true)
	}
	return root
}

func buildHeader(c *ui.Context, spec calendarSpec, view *CalendarViewState, month, cursor time.Time,
	monthText string, height float32, focusHeading *bool) {
	t := theme.Colors(c)
	radii := theme.RadiusScale()
	navRadius := radii.LG + 4 // HeroUI rounded-2xl for the calendar header.

	ui.Row(c).AlignItems(ui.Center).Gap(theme.Space4).
		Padding(0, 2, theme.Space16, 2).Shrink(0).Children(func() {
		open := view.YearOpen && !spec.hideYearPicker
		ink := t.Foreground
		if open {
			ink = t.AccentSoftForeground
		}

		var heading ui.Element
		if spec.hideYearPicker {
			heading = ui.Row(c).Grow(1).MinWidth(0).Height(height).AlignItems(ui.Center)
		} else {
			heading = ui.ButtonBase(c).Grow(1).MinWidth(0).Height(height).
				Radius(radii.SM).FocusRing(true).
				Label(monthText + ", " + spec.years()).Cursor(ui.CursorPointer)
			heading = heading.OnClick(func() {
				view.YearOpen = !view.YearOpen
				if !view.YearOpen {
					*focusHeading = true
				}
			})
		}
		heading.Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(theme.Space4).Children(func() {
				style.Text(ui.Text(c, monthText), c, theme.SM()).
					FontWeight(theme.WeightMedium).TextColor(ink).SingleLine()
				if !spec.hideYearPicker {
					icon := ui.Icon(c, lucide.ChevronRight()).Size(theme.Scale(c, 12), theme.Scale(c, 12)).
						TextColor(t.AccentSoftForeground)
					if open {
						icon.Rotate(90)
					}
				}
			})
		})
		if *focusHeading {
			heading.Focus()
			*focusHeading = false
		}

		previous := navButton(c, spec.prev(), navRadius)
		if canShiftMonth(month, spec.min, spec.max, -1) {
			previous = previous.OnClick(func() { shiftMonth(view, month, cursor, -1) })
		} else {
			previous = style.Disabled(previous, true)
		}
		previous.Children(func() {
			ui.Icon(c, lucide.ChevronLeft()).Size(theme.Scale(c, 16), theme.Scale(c, 16)).
				TextColor(t.AccentSoftForeground)
		})

		next := navButton(c, spec.next(), navRadius)
		if canShiftMonth(month, spec.min, spec.max, 1) {
			next = next.OnClick(func() { shiftMonth(view, month, cursor, 1) })
		} else {
			next = style.Disabled(next, true)
		}
		next.Children(func() {
			ui.Icon(c, lucide.ChevronRight()).Size(theme.Scale(c, 16), theme.Scale(c, 16)).
				TextColor(t.AccentSoftForeground)
		})
	})
}

func navButton(c *ui.Context, label string, radius float32) ui.Element {
	t := theme.Colors(c)
	size := theme.Scale(c, 24)
	e := ui.ButtonBase(c).Size(size, size).Center().Shrink(0).
		Radius(radius).FocusRing(true).Label(label).TextColor(t.AccentSoftForeground)
	return style.Background(e, radius, ui.Transparent, t.Default, t.Default)
}

func buildDayGrid(c *ui.Context, spec calendarSpec, view *CalendarViewState, sel cellSelection,
	month, cursor, today time.Time, cellW, cellH, weekdayH, gap float32) {
	t := theme.Colors(c)
	label := spec.label
	if label == "" {
		label = spec.monthText(month)
	}

	grid := ui.Column(c).Shrink(0).Focusable().FocusRing(false).Role(ui.RoleTable).Label(label)
	if spec.autofocus {
		grid = grid.AutoFocus()
	}
	var cursorCell ui.Element

	grid.Children(func() {
		ui.Row(c).Shrink(0).Height(weekdayH).Children(func() {
			for i := range 7 {
				day := time.Weekday((int(spec.firstWeekday) + i) % 7)
				cell := ui.Box(c).Width(cellW).Height(weekdayH).Center().
					Role(ui.RoleColumnHeader).Label(day.String())
				cell.Children(func() {
					style.Text(ui.Text(c, spec.weekdayText(day)), c, theme.XS()).
						FontWeight(theme.WeightMedium).TextColor(t.Muted).SingleLine()
				})
			}
		})

		lead := (int(month.Weekday()) - int(spec.firstWeekday) + 7) % 7
		start := month.AddDate(0, 0, -lead)
		ui.Column(c).Shrink(0).Gap(0).Children(func() {
			for week := range 6 {
				row := ui.Row(c).Role(ui.RoleRow)
				if week == 0 {
					row = row.Margin(gap, 0, 0, 0)
				}
				row.Children(func() {
					for col := range 7 {
						day := start.AddDate(0, 0, week*7+col)
						cell := dayCell(c, spec, view, sel, day, month, cursor, today,
							cellW, cellH, col, grid)
						if sameDay(day, cursor) {
							cursorCell = cell
						}
					}
				})
			}
		})
	})

	if cursorCell.Valid() {
		grid.ActiveDescendant(cursorCell)
	}

	grid.OnShortcut(0, ui.KeyLeft, func() { moveCursor(view, month, cursor, -1) })
	grid.OnShortcut(0, ui.KeyRight, func() { moveCursor(view, month, cursor, 1) })
	grid.OnShortcut(0, ui.KeyUp, func() { moveCursor(view, month, cursor, -7) })
	grid.OnShortcut(0, ui.KeyDown, func() { moveCursor(view, month, cursor, 7) })
	grid.OnShortcut(0, ui.KeyPageUp, func() { shiftMonth(view, month, cursor, -1) })
	grid.OnShortcut(0, ui.KeyPageDown, func() { shiftMonth(view, month, cursor, 1) })
	grid.OnShortcut(0, ui.KeyHome, func() { moveCursorTo(view, month, month) })
	grid.OnShortcut(0, ui.KeyEnd, func() {
		last := time.Date(month.Year(), month.Month(), daysIn(month), 0, 0, 0, 0, month.Location())
		moveCursorTo(view, month, last)
	})
	grid.OnShortcut(0, ui.KeyEnter, func() { activateDay(spec, view, month, sel, cursor) })
	grid.OnShortcut(0, ui.KeySpace, func() { activateDay(spec, view, month, sel, cursor) })
}

func dayCell(c *ui.Context, spec calendarSpec, view *CalendarViewState, sel cellSelection,
	day, month, cursor, today time.Time, cellW, cellH float32, col int, grid ui.Element) ui.Element {
	t := theme.Colors(c)
	outside := !sameMonth(day, month)
	disabled := spec.disabledDay(day)
	rangeMode := sel.inRange != nil
	inRange := rangeMode && sel.inRange(day)
	isStart := rangeMode && sel.start != nil && sel.start(day) && !outside
	isEnd := rangeMode && sel.end != nil && sel.end(day) && !outside
	isSelected := sel.selected != nil && sel.selected(day)
	isToday := sameDay(day, today) && !disabled
	isCursor := sameDay(day, cursor)

	ink := t.Foreground
	switch {
	case outside:
		ink = t.Muted
	case isStart || isEnd:
		ink = t.AccentForeground
	case isSelected:
		ink = t.AccentForeground
	case isToday:
		ink = t.AccentSoftForeground
	case rangeMode && inRange:
		ink = t.Foreground
	}

	// HeroUI's rounded-3xl cell: 24 DIP, clamped to half the square cell so
	// the endpoints read as circles rather than rounded squares.
	radius := min(theme.Scale(c, 24), cellH/2)
	bandEdge := theme.RadiusScale().LG
	cell := ui.Box(c).Key(dayKey{day.Year(), day.Month(), day.Day()}).
		Size(cellW, cellH).Center().Shrink(0).Radius(radius).
		Role(ui.RoleButton).Label(spec.dayText(day)).
		Selected(isSelected)
	if !disabled && sel.activate != nil {
		cell = cell.Cursor(ui.CursorPointer).OnClick(func() { activateDay(spec, view, month, sel, day) })
	}
	if disabled {
		cell = cell.Disabled(true)
	} else if outside {
		cell = cell.Opacity(0.5)
	}

	cell.Children(func() {
		text := style.Text(ui.Text(c, strconv.Itoa(day.Day())), c, theme.SM()).
			FontWeight(theme.WeightMedium).TextColor(ink).SingleLine().FontFeatures("tnum")
		if disabled && !outside {
			text.Strikethrough()
		}
	})
	// Snapshot the cell's states in the build: a retained copy of a
	// leaving element paints in later passes, where cell has expired.
	cellPressed, cellHovered := cell.Pressed(), cell.Hovered()
	cell.Draw(func(p *ui.Painter, r ui.Rect) {
		// Disabled and outside days stay plain: 50% opacity and, for disabled
		// days, a line-through, but never a hover, today, selected or range
		// fill. A disabled day inside a range breaks the band.
		if disabled || outside {
			return
		}
		if rangeMode {
			if inRange {
				// The band fills the whole cell, square in the middle of a
				// week and HeroUI's rounded-lg where the row or the range
				// ends, with the cap radius on an endpoint's outer side.
				var tl, tr, br, bl float32
				switch {
				case isStart:
					tl, bl = radius, radius
				case col == 0 || !bandContinues(spec, sel, day.AddDate(0, 0, -1), month):
					tl, bl = bandEdge, bandEdge
				}
				switch {
				case isEnd:
					tr, br = radius, radius
				case col == 6 || !bandContinues(spec, sel, day.AddDate(0, 0, 1), month):
					tr, br = bandEdge, bandEdge
				}
				fillBand(p, r, t.AccentSoft, tl, tr, br, bl)
			}
			switch {
			case isStart || isEnd:
				fill := t.Accent
				if cellPressed {
					fill = t.AccentHover
				}
				p.Fill(r, fill, radius)
			case isToday:
				fill := t.AccentSoft
				if cellHovered {
					fill = t.AccentSoftHover
				}
				p.Fill(r, fill, radius)
			case cellHovered:
				p.Fill(r, t.Default, radius)
			}
			return
		}
		switch {
		case isSelected:
			fill := t.Accent
			if cellPressed {
				fill = t.AccentHover
			}
			p.Fill(r, fill, radius)
		case isToday:
			fill := t.AccentSoft
			if cellHovered {
				fill = t.AccentSoftHover
			}
			p.Fill(r, fill, radius)
		case cellHovered:
			p.Fill(r, t.Default, radius)
		}
	})
	if isCursor && grid.FocusVisible() {
		cell.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.FocusRing(r, [4]float32{radius, radius, radius, radius})
		})
	}
	return cell
}

// bandContinues reports whether the range band covers d: in range, in the
// visible month and selectable. Disabled and outside days break the band.
func bandContinues(spec calendarSpec, sel cellSelection, d, month time.Time) bool {
	return sel.inRange(d) && sameMonth(d, month) && !spec.disabledDay(d)
}

// bandPathPool reuses the path buffer of a band with corners of unequal radii.
var bandPathPool = sync.Pool{New: func() any { return &ui.Path{} }}

// fillBand paints a rectangle with independent corner radii, as CSS's
// border-radius does; Painter.Fill takes one radius for all four corners.
func fillBand(p *ui.Painter, r ui.Rect, c ui.Color, tl, tr, br, bl float32) {
	if r.W <= 0 || r.H <= 0 {
		return
	}
	limit := min(r.W, r.H) / 2
	tl, tr, br, bl = min(tl, limit), min(tr, limit), min(br, limit), min(bl, limit)
	if tl == tr && tr == br && br == bl {
		p.Fill(r, c, tl)
		return
	}
	const k = 0.5523 // control point distance of a quarter circle
	path := bandPathPool.Get().(*ui.Path)
	path.Reset()
	path.MoveTo(r.X+tl, r.Y)
	path.LineTo(r.X+r.W-tr, r.Y)
	if tr > 0 {
		path.CubeTo(r.X+r.W-tr+tr*k, r.Y, r.X+r.W, r.Y+tr-tr*k, r.X+r.W, r.Y+tr)
	}
	path.LineTo(r.X+r.W, r.Y+r.H-br)
	if br > 0 {
		path.CubeTo(r.X+r.W, r.Y+r.H-br+br*k, r.X+r.W-br+br*k, r.Y+r.H, r.X+r.W-br, r.Y+r.H)
	}
	path.LineTo(r.X+bl, r.Y+r.H)
	if bl > 0 {
		path.CubeTo(r.X+bl-bl*k, r.Y+r.H, r.X, r.Y+r.H-bl+bl*k, r.X, r.Y+r.H-bl)
	}
	path.LineTo(r.X, r.Y+tl)
	if tl > 0 {
		path.CubeTo(r.X, r.Y+tl-tl*k, r.X+tl-tl*k, r.Y, r.X+tl, r.Y)
	}
	path.Close()
	p.FillPath(path, c)
	bandPathPool.Put(path)
}

// activateDay applies selection, following the day into its month.
func activateDay(spec calendarSpec, view *CalendarViewState, month time.Time, sel cellSelection, day time.Time) {
	if spec.disabledDay(day) || sel.activate == nil {
		return
	}
	if view != nil {
		if !sameMonth(day, month) || view.Month.IsZero() {
			view.Month = monthStart(day)
		}
		view.Cursor = day
	}
	sel.activate(day)
}

func moveCursor(view *CalendarViewState, month, cursor time.Time, days int) {
	next := cursor.AddDate(0, 0, days)
	if !sameMonth(next, month) {
		view.Month = monthStart(next)
	} else if view.Month.IsZero() {
		view.Month = month
	}
	view.Cursor = next
}

func moveCursorTo(view *CalendarViewState, month, day time.Time) {
	view.Month = month
	view.Cursor = day
}

// shiftMonth moves the visible month by delta, clamping the cursor's day to
// the new month.
func shiftMonth(view *CalendarViewState, month, cursor time.Time, delta int) {
	next := month.AddDate(0, delta, 0)
	view.Month = monthStart(next)
	view.Cursor = clampToMonth(cursor, next)
}

func canShiftMonth(month, min, max time.Time, delta int) bool {
	next := month.AddDate(0, delta, 0)
	if !min.IsZero() && compareMonth(next, min) < 0 {
		return false
	}
	if !max.IsZero() && compareMonth(next, max) > 0 {
		return false
	}
	return true
}

// buildYearGrid paints the 3-column year picker over the day grid's area.
func buildYearGrid(c *ui.Context, spec calendarSpec, view *CalendarViewState, month, cursor time.Time,
	height float32, focusHeading *bool) {
	t := theme.Colors(c)
	lo, hi := yearPage(month.Year(), spec.min, spec.max)
	selectedYear := month.Year()
	if selectedYear < lo {
		selectedYear = lo
	}
	if selectedYear > hi {
		selectedYear = hi
	}
	cursorYear := cursor.Year()
	if cursorYear < lo {
		cursorYear = lo
	}
	if cursorYear > hi {
		cursorYear = hi
	}
	cellH := theme.Scale(c, 32)
	radius := cellH / 2
	itemH := cellH + theme.Space4

	scroll := ui.Scroll(c).Height(height).Shrink(0)
	scroller := ui.Local(scroll, yearScrollKey{}, func() ui.ScrollState { return ui.ScrollState{} })
	scroll.TrackScroll(scroller)
	scrolled := ui.Local(scroll, yearScrollInitKey{}, func() bool { return false })
	if !*scrolled {
		row := float32((selectedYear-lo)/3) * itemH
		scroller.Y = max(0, row-height/2+itemH)
		*scrolled = true
	}
	wantScroll := ui.Local(scroll, yearFocusKey{}, func() bool { return false })

	var grid ui.Element
	var cursorCell ui.Element
	scroll.Children(func() {
		grid = ui.Grid(c).Columns(3).Gap(theme.Space4).Padding(theme.Space4).
			Focusable().FocusRing(false).Role(ui.RoleGroup).Label(spec.years()).AutoFocus()
		grid.Children(func() {
			for y := lo; y <= hi; y++ {
				year := y
				cell := ui.Box(c).Key(yearKey(year)).Height(cellH).Center().
					Radius(radius).Role(ui.RoleButton).Label(strconv.Itoa(year)).
					Selected(year == selectedYear).Cursor(ui.CursorPointer)
				cell = cell.OnClick(func() { chooseYear(view, month, cursor, year, focusHeading) })
				ink := t.Foreground
				if year == selectedYear {
					ink = t.AccentForeground
				}
				cell.Children(func() {
					style.Text(ui.Text(c, strconv.Itoa(year)), c, theme.SM()).
						FontWeight(theme.WeightMedium).TextColor(ink).SingleLine().FontFeatures("tnum")
				})
				// Snapshot the cell's states in the build: a retained copy of
				// a leaving element paints in later passes, where cell has
				// expired.
				cellPressed, cellHovered := cell.Pressed(), cell.Hovered()
				cell.Draw(func(p *ui.Painter, r ui.Rect) {
					switch {
					case year == selectedYear:
						fill := t.Accent
						if cellPressed || cellHovered {
							fill = t.AccentHover
						}
						p.Fill(r, fill, radius)
					case cellHovered:
						p.Fill(r, t.Default, radius)
					}
				})
				if year == cursorYear && grid.FocusVisible() {
					cell.DrawOver(func(p *ui.Painter, r ui.Rect) {
						p.FocusRing(r, [4]float32{radius, radius, radius, radius})
					})
				}
				if year == cursorYear {
					cursorCell = cell
				}
			}
		})
	})

	if cursorCell.Valid() {
		grid.ActiveDescendant(cursorCell)
	}
	moveYear := func(delta int) {
		next := cursorYear + delta
		if next < lo || next > hi {
			return
		}
		setYear(view, month, cursor, next)
		*wantScroll = true
	}
	grid.OnShortcut(0, ui.KeyLeft, func() { moveYear(-1) })
	grid.OnShortcut(0, ui.KeyRight, func() { moveYear(1) })
	grid.OnShortcut(0, ui.KeyUp, func() { moveYear(-3) })
	grid.OnShortcut(0, ui.KeyDown, func() { moveYear(3) })
	grid.OnShortcut(0, ui.KeyHome, func() { moveYear(lo - cursorYear) })
	grid.OnShortcut(0, ui.KeyEnd, func() { moveYear(hi - cursorYear) })
	grid.OnShortcut(0, ui.KeyEnter, func() { chooseYear(view, month, cursor, cursorYear, focusHeading) })
	grid.OnShortcut(0, ui.KeySpace, func() { chooseYear(view, month, cursor, cursorYear, focusHeading) })
	grid.OnShortcut(0, ui.KeyEscape, func() {
		view.YearOpen = false
		*focusHeading = true
	})
	if *wantScroll && cursorCell.Valid() {
		cursorCell.ScrollIntoView()
		*wantScroll = false
	}
}

func chooseYear(view *CalendarViewState, month, cursor time.Time, year int, focusHeading *bool) {
	setYear(view, month, cursor, year)
	view.YearOpen = false
	*focusHeading = true
}

func setYear(view *CalendarViewState, month, cursor time.Time, year int) {
	next := time.Date(year, month.Month(), 1, 0, 0, 0, 0, month.Location())
	view.Month = next
	view.Cursor = clampToMonth(cursor, next)
}

// yearPage returns the years the picker offers. Unbounded, it is one page of
// twelve years centred on the visible year, so the picker opens on a full
// grid rather than a scroll through two centuries. Min and Max bound the
// offered years: with both set the picker lists exactly the bounded years,
// and with one set the page slides inside that bound.
func yearPage(visible int, earliest, latest time.Time) (int, int) {
	const page = 12
	var lo, hi int
	switch {
	case !earliest.IsZero() && !latest.IsZero():
		return yearRange(earliest, latest)
	case !earliest.IsZero():
		lo = max(visible-page/2, earliest.Year())
		hi = lo + page - 1
	case !latest.IsZero():
		hi = min(visible+page/2-1, latest.Year())
		lo = hi - page + 1
	default:
		lo = visible - page/2
		hi = lo + page - 1
	}
	return max(lo, 1), min(hi, 9999)
}

// yearRange returns the years the two bounds cover, in order and clamped to
// the years time.Time can express.
func yearRange(from, to time.Time) (int, int) {
	lo, hi := from.Year(), to.Year()
	if lo > hi {
		lo, hi = hi, lo
	}
	return max(lo, 1), min(hi, 9999)
}

// effectiveMonth resolves the visible month without mutating the state.
func (v *CalendarViewState) effectiveMonth(value, today time.Time) time.Time {
	switch {
	case !v.Month.IsZero():
		return monthStart(v.Month)
	case !value.IsZero():
		return monthStart(value)
	}
	return monthStart(today)
}

// cursorIn resolves the keyboard cursor within month.
func (v *CalendarViewState) cursorIn(month, value, today time.Time) time.Time {
	if !v.Cursor.IsZero() && sameMonth(v.Cursor, month) {
		return dayStart(v.Cursor)
	}
	if !value.IsZero() && sameMonth(value, month) {
		return dayStart(value)
	}
	if !today.IsZero() && sameMonth(today, month) {
		return dayStart(today)
	}
	return month
}

type dayKey struct {
	year  int
	month time.Month
	day   int
}

type yearKey int

// Calendar date arithmetic. Only year, month and day participate; the clock
// and location of a stored time are preserved but never compared.
func dayStart(t time.Time) time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func monthStart(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func sameDay(a, b time.Time) bool {
	if a.IsZero() || b.IsZero() {
		return a.IsZero() && b.IsZero()
	}
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func sameMonth(a, b time.Time) bool {
	if a.IsZero() || b.IsZero() {
		return false
	}
	ay, am, _ := a.Date()
	by, bm, _ := b.Date()
	return ay == by && am == bm
}

func compareDay(a, b time.Time) int {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	switch {
	case ay != by:
		if ay < by {
			return -1
		}
		return 1
	case am != bm:
		if am < bm {
			return -1
		}
		return 1
	case ad != bd:
		if ad < bd {
			return -1
		}
		return 1
	}
	return 0
}

func compareMonth(a, b time.Time) int {
	ay, am, _ := a.Date()
	by, bm, _ := b.Date()
	av, bv := ay*12+int(am), by*12+int(bm)
	switch {
	case av < bv:
		return -1
	case av > bv:
		return 1
	}
	return 0
}

func daysIn(t time.Time) int {
	y, m, _ := t.Date()
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func clampToMonth(t, month time.Time) time.Time {
	if t.IsZero() {
		return month
	}
	day := t.Day()
	if last := daysIn(month); day > last {
		day = last
	}
	return time.Date(month.Year(), month.Month(), day, 0, 0, 0, 0, month.Location())
}
