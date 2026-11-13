package navigation

import (
	"math"
	"strconv"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PipsPagerOptions configures a row of pips, one per page. Index is 0-based
// and caller-owned; nil keeps it internally. MaxVisible bounds the pips built
// around the current one (default 7), however many pages there are.
type PipsPagerOptions struct {
	Index      *int
	Count      int
	MaxVisible int
	Label      string
	OnChange   func(index int)
}

// CarouselOptions configures a carousel. Index is the caller-owned 0-based
// first slide in view; nil keeps it internally. Slide builds slide index, and
// only the slides in view are built. PerView shows several slides side by
// side (default 1). Loop continues from the last slide to the first.
type CarouselOptions struct {
	Index    *int
	Count    int
	Slide    func(index int)
	PerView  int
	Loop     bool
	Label    string
	OnChange func(index int)
}

type (
	pipsIndexKey     struct{}
	carouselIndexKey struct{}
	carouselStateKey struct{}
	carouselPullKey  struct{}
	// carouselNeighbourKey is the slide a drag brings in, until a step makes
	// it a slide of its own that comes in from where it was.
	carouselNeighbourKey struct{}
)

// carouselState is the slide the last build showed and the way the carousel
// last stepped: 1 forward, -1 back. The rest is the gesture in progress: the
// pointer's pull on the slides in DIPs, whether the drag turned out vertical
// (left to the page), and the horizontal scroll of one swipe of a touchpad.
type carouselState struct {
	index int
	dir   float32

	dragging     bool
	vertical     bool
	pull, sx, sy float32
	wheelAt      time.Time
	wheelSum     float32
	wheelSpent   bool
	wheelBy      int
}

const (
	pipsDefaultVisible = 7
	carouselMotion     = 250 * time.Millisecond

	// A drag turns a page once it carries the slides a quarter of the way,
	// or flings them faster than carouselFling DIPs a second; one that moves
	// less than carouselSlop is not yet horizontal or vertical.
	carouselFling = 500
	carouselSlop  = 6
	// A swipe on a touchpad turns one page after carouselWheelStep DIPs and
	// is over once the scrolling pauses for carouselWheelQuiet, momentum
	// included. The slides resist a pull past the ends.
	carouselWheelStep  = 48
	carouselWheelQuiet = 150 * time.Millisecond
	carouselResist     = 0.35
)

// PipsPager shows one pip per page and marks the current one. A click
// chooses a page; the pips are one focus stop whose Left and Right arrows
// move among them.
func PipsPager(c *ui.Context, o PipsPagerOptions) ui.Element {
	root := ui.Row(c).Shrink(0).AlignItems(ui.Center).Role(ui.RoleGroup).FocusGroup(ui.Horizontal)
	if o.Label != "" {
		root.Label(o.Label)
	}
	index := o.Index
	if index == nil {
		index = ui.Local(root, pipsIndexKey{}, func() int { return 0 })
	}
	if o.Count <= 0 {
		return root
	}
	*index = max(0, min(*index, o.Count-1))
	first, end := pipsWindow(*index, o.Count, o.MaxVisible)
	root.Children(func() {
		for i := first; i < end; i++ {
			// A pip at the window's edge shrinks while pages continue past it.
			more := (i == first && first > 0) || (i == end-1 && end < o.Count)
			pip := buildPip(c.Key(i), i, o.Count, i == *index, more)
			if i == *index {
				continue
			}
			pip.OnClick(func() {
				*index = i
				if o.OnChange != nil {
					o.OnChange(i)
				}
			})
		}
	})
	return root
}

// pipsWindow returns the pips to build, [first, end), centered on index.
func pipsWindow(index, count, visible int) (first, end int) {
	if visible <= 0 {
		visible = pipsDefaultVisible
	}
	if count <= visible {
		return 0, count
	}
	first = max(0, min(index-visible/2, count-visible))
	return first, first + visible
}

func buildPip(c *ui.Context, i, count int, active, more bool) ui.Element {
	t := theme.Colors(c)
	side := scale(c, theme.Space16)
	pip := ui.ButtonBase(c).Size(side, side).Center().Shrink(0).Radius(side/2).FocusRing(true).
		Label(locale.M().Page+" "+strconv.Itoa(i+1)).PositionInSet(i+1, count)
	dot := scale(c, theme.Space6)
	switch {
	case active:
		dot = scale(c, theme.Space8)
		pip.Current(ui.CurrentPage)
	case more:
		dot = scale(c, theme.Space4)
	}
	// The pip that becomes current grows and darkens as the one that was
	// current shrinks and fades; hover shows at once.
	ink := selectionInk(pip, "pip-ink", active, pip.Hovered(),
		t.MutedForeground.Alpha(0.4), t.Foreground, t.MutedForeground)
	dot = pip.AnimateTo("pip-size", dot, ui.Animation{Duration: motion.State, Ease: ui.EaseOut})
	pip.Children(func() {
		ui.Box(c).Size(dot, dot).Radius(dot / 2).Background(ink)
	})
	return pip
}

// Carousel shows one slide at a time, or PerView of them, between previous
// and next buttons over a PipsPager. Left and Right step while the slides
// or the buttons have the focus; Label names the slides, which are a focus
// stop. Dragging the slides sideways or scrolling sideways over them turns
// a page, once per drag or swipe; a drag that falls short springs back, and
// a vertical one is left to the page. Slides move a slide's width as they
// come and go, and change at once under reduced motion. It never advances
// on its own.
func Carousel(c *ui.Context, o CarouselOptions) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(scale(c, theme.Space8))
	index := o.Index
	if index == nil {
		index = ui.Local(root, carouselIndexKey{}, func() int { return 0 })
	}
	perView := max(o.PerView, 1)
	last := max(o.Count-perView, 0)
	// The caller's index is only read here, so slides that reload do not
	// reset it; a step writes it back.
	current := max(0, min(*index, last))
	st := ui.Local(root, carouselStateKey{}, func() carouselState {
		return carouselState{index: current, dir: 1}
	})
	// A change made elsewhere, by the caller or the pips, goes the way the
	// index went.
	if current != st.index {
		st.dir = 1
		if current < st.index {
			st.dir = -1
		}
		st.index = current
	}
	// target is the slide a step by by goes to, if there is one.
	target := func(by int) (int, bool) {
		to := current + by
		switch {
		case to >= 0 && to <= last:
		case !o.Loop:
			return current, false
		case to < 0:
			to = last
		default:
			to = 0
		}
		return to, to != current
	}
	goTo := func(to int) {
		dir := float32(1)
		if to < current {
			dir = -1
		}
		*index, st.index, st.dir, current = to, to, dir, to
		if o.OnChange != nil {
			o.OnChange(to)
		}
	}
	step := func(by int) {
		if to, ok := target(by); ok {
			goTo(to)
		}
	}
	// forward is the inline direction of the next slide on screen.
	forward := 1
	if c.Direction() == ui.RightToLeft {
		forward = -1
	}
	if by := st.wheelBy; by != 0 {
		st.wheelBy = 0
		step(by * forward)
	}
	swiped := 0
	root.Children(func() {
		// The arrows step from the buttons and the slides; in the pips they
		// move the focus from pip to pip.
		row := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(scale(c, theme.Space8))
		row.OnShortcut(0, ui.KeyLeft, func() { step(-forward) })
		row.OnShortcut(0, ui.KeyRight, func() { step(forward) })
		row.OnShortcut(0, ui.KeyHome, func() {
			if current != 0 {
				goTo(0)
			}
		})
		row.OnShortcut(0, ui.KeyEnd, func() {
			if current != last {
				goTo(last)
			}
		})
		row.Children(func() {
			carouselArrow(c, locale.M().Previous, -forward, o.Loop || current > 0, func() { step(-1) })
			swiped = carouselSlides(c, o, st, current, perView, forward, target)
			carouselArrow(c, locale.M().Next, forward, o.Loop || current < last, func() { step(1) })
		})
		if last > 0 {
			pager := ui.Row(c).AlignSelf(ui.Center).AlignItems(ui.Center).Gap(scale(c, theme.Space8))
			pager.Children(func() {
				PipsPager(c, PipsPagerOptions{Index: index, Count: last + 1, OnChange: o.OnChange})
				// Pips are a window over the pages once there are more than
				// they hold, so the position is said as well.
				if last+1 > pipsDefaultVisible {
					style.Text(ui.Text(c, strconv.Itoa(current+1)+" / "+strconv.Itoa(last+1)), c, theme.XS()).
						TextColor(theme.Colors(c).MutedForeground).FontFeatures("tnum").SingleLine().AccessHidden(true)
				}
			})
		}
	})
	if swiped != 0 {
		step(swiped)
		c.Invalidate()
	}
	return root
}

// carouselArrow builds a step button whose chevron points toward side: -1
// left, 1 right.
func carouselArrow(c *ui.Context, label string, side int, enabled bool, onClick func()) {
	icon := lucide.ChevronRight()
	if side < 0 {
		icon = lucide.ChevronLeft()
	}
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Icon: icon, IconOnly: true, Variant: mygoui.ButtonOutline,
		Size: theme.ButtonSmall, Disabled: !enabled, OnClick: onClick,
	}).AlignSelf(ui.Center)
}

// carouselSlides builds the slides in view and reads the drag on them. It
// returns the step a drag asked for: -1 or 1 as the pages go, 0 for none.
// forward is 1 where the next slide is on the right and target says where a
// step goes, if anywhere.
func carouselSlides(c *ui.Context, o CarouselOptions, st *carouselState, index, perView, forward int, target func(by int) (int, bool)) (by int) {
	viewport := ui.Row(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Stretch).
		Clip().Focusable().FocusRing(true).
		Radius(theme.RadiusScale().LG).Role(ui.RoleGroup)
	if o.Label != "" {
		viewport.Label(o.Label)
	}
	animated := !theme.ReduceMotion(c)
	viewport.HandleInput(func(ev ui.InputEvent) bool { return carouselWheel(st, ev) })
	by = carouselDrag(viewport, st, forward, target)

	viewport.Children(func() {
		// The slides move together, so that a drag moves them as one and the
		// slides of a step come on from where the drag left them.
		track := ui.Row(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Stretch).Gap(scale(c, theme.Space8))
		switch pulled := carouselPulled(st, forward, target); {
		case st.dragging:
			track.Translate(track.Jump(carouselPullKey{}, pulled), 0)
		case animated:
			track.Translate(track.AnimateTo(carouselPullKey{}, 0, ui.Animation{Duration: motion.Reveal, Ease: ui.EaseOut}), 0)
		default:
			track.Translate(track.Jump(carouselPullKey{}, 0), 0)
		}
		// Slides only translate, so their motion needs no bitmap of the slide:
		// the GPU moves what they paint.
		shift := st.dir * float32(forward)
		enter := ui.Motion{XPercent: 100 * shift, X: scale(c, theme.Space8) * shift}
		exit := ui.Motion{XPercent: -100 * shift}
		if animated {
			// A slide leaves the way this step goes, not the way the step that
			// built it went.
			track.ChildExit(exit)
		}
		track.Children(func() {
			for i := index; i < min(index+perView, o.Count); i++ {
				slide := ui.Column(c.Key(i)).Grow(1).Basis(0).MinWidth(0).
					Role(ui.RoleGroup).PositionInSet(i+1, o.Count)
				if animated {
					slide.Transition(ui.ElementTransition{
						Position: true, Duration: carouselMotion, Ease: ui.EaseOut,
						Enter: &enter, Exit: &exit,
					})
				}
				slide.Children(func() { o.Slide(i) })
			}
			if st.dragging && perView == 1 {
				carouselNeighbour(c, o, st, viewport.LayoutBounds().W, forward, target)
			}
		})
	})
	return by
}

// carouselDrag follows the pointer pulling the slides and, once it lets go,
// returns the step the pull amounts to: past a quarter of the slides' width
// or flung fast enough, and not a vertical drag.
func carouselDrag(viewport ui.Element, st *carouselState, forward int, target func(by int) (int, bool)) int {
	dx, dy, held := viewport.Dragged()
	vx, _, ended := viewport.DragEnded()
	if held {
		st.dragging = true
		st.sx += dx
		st.sy += dy
		if !st.vertical && math.Abs(float64(st.sy)) > math.Abs(float64(st.sx)) &&
			math.Abs(float64(st.sy)) > carouselSlop {
			st.vertical = true
		}
		if !st.vertical {
			st.pull += dx
		}
	}
	if !ended && (held || !st.dragging) {
		return 0
	}
	// By the build after the pointer let go, or the system took the press.
	pull, width := st.pull, viewport.LayoutBounds().W
	st.dragging, st.vertical, st.pull, st.sx, st.sy = false, false, 0, 0, 0
	far := math.Abs(float64(pull)) > float64(width)/4
	flung := math.Abs(float64(pull)) > carouselSlop && math.Abs(float64(vx)) > carouselFling && (vx < 0) == (pull < 0)
	if !ended || !far && !flung {
		return 0
	}
	// Pulled toward the start of the screen shows what comes after.
	by := forward
	if pull > 0 {
		by = -forward
	}
	if _, ok := target(by); !ok {
		return 0
	}
	return by
}

// carouselPulled is how far the slides are shown pulled: the pointer's pull,
// held back where there is no slide to go to.
func carouselPulled(st *carouselState, forward int, target func(by int) (int, bool)) float32 {
	by := forward
	if st.pull > 0 {
		by = -forward
	}
	if _, ok := target(by); ok {
		return st.pull
	}
	return st.pull * carouselResist
}

// carouselWheel takes the sideways scrolling of a touchpad over the slides:
// one swipe, momentum included, asks for one step. Scrolling that is not
// mostly sideways is left to the page.
func carouselWheel(st *carouselState, ev ui.InputEvent) bool {
	if ev.Kind != ui.InputScroll || ev.Mods&(ui.Ctrl|ui.Cmd|ui.Alt) != 0 ||
		math.Abs(float64(ev.DX)) <= math.Abs(float64(ev.DY)) {
		return false
	}
	if ev.Time.Sub(st.wheelAt) > carouselWheelQuiet {
		st.wheelSum, st.wheelSpent = 0, false
	}
	st.wheelAt = ev.Time
	if st.wheelSpent {
		return true
	}
	st.wheelSum += ev.DX
	// A notch of a wheel is a step by itself.
	if !ev.Precise || math.Abs(float64(st.wheelSum)) >= carouselWheelStep {
		st.wheelBy, st.wheelSpent = 1, true
		if st.wheelSum < 0 {
			st.wheelBy = -1
		}
	}
	return true
}

// carouselNeighbour builds the slide the pull is bringing in, beside the
// one in view, so that the drag shows where it goes.
func carouselNeighbour(c *ui.Context, o CarouselOptions, st *carouselState, width float32, forward int, target func(by int) (int, bool)) {
	by, side := forward, float32(1)
	if st.pull > 0 {
		by, side = -forward, -1
	}
	to, ok := target(by)
	if !ok {
		return
	}
	gap := scale(c, theme.Space8)
	slide := ui.Column(c.Key(carouselNeighbourKey{})).Absolute().Top(0).Width(width).Left(side*(width+gap)).
		Role(ui.RoleGroup).PositionInSet(to+1, o.Count).AccessHidden(true)
	slide.Children(func() { o.Slide(to) })
}
