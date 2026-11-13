package overlay

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TourStep is one stop of a tour. Anchor is the handle of the step's target,
// bound with Element.Bind. A step without a handle, or whose target is not
// laid out, shows its panel in the middle of the window without a spotlight,
// so the tour still goes on.
type TourStep struct {
	Anchor      *ui.Handle
	Title, Body string
	// Placement is the side of the target the panel prefers, its alignment
	// along it and its gap; SideAuto prefers below. The panel takes another
	// side where it does not fit, and stays inside the window.
	Placement mygoui.Placement
}

// TourState is the caller-owned progress of a tour: whether it shows and the
// index of its current step. Keep one per tour and window.
type TourState struct {
	Open bool
	Step int
}

// TourOptions configures a tour.
type TourOptions struct {
	// Width is the panel's width; it defaults to 320 DIPs. The panel keeps
	// it from step to step, so only its height follows the content.
	Width float32
	// OnFinish runs after Done on the last step closed the tour; OnDismiss
	// after the close button or Escape did.
	OnFinish, OnDismiss func()
}

const (
	// tourPad is the spotlight's margin around its target.
	tourPad = 6
	// tourDots is the most steps the progress shows as dots; longer tours
	// show "2 / 12".
	tourDots = 8
	// tourFadeOut is the share of a step change over which the old content
	// fades out; the new content fades in over the rest, so the two never
	// show at full strength together.
	tourFadeOut = 0.4
)

type (
	tourKey       struct{ state *TourState }
	tourViewKey   struct{}
	tourTravelKey struct{}
	tourFadeKey   struct{}
	tourHeightKey struct{}
	tourPanelKey  struct{}
	tourStackKey  struct{}
	tourLayerKey  struct{ step int }
	tourBackKey   struct{}
	tourNextKey   struct{}
	tourDotKey    struct{ step int }
)

// tourGeometry is what moves between steps: the spotlight, empty for none,
// and the panel's top left corner.
type tourGeometry struct {
	hole ui.Rect
	x, y float32
}

func (g tourGeometry) lerp(to tourGeometry, k float32) tourGeometry {
	mix := func(a, b float32) float32 { return a + (b-a)*k }
	return tourGeometry{
		hole: ui.Rect{X: mix(g.hole.X, to.hole.X), Y: mix(g.hole.Y, to.hole.Y), W: mix(g.hole.W, to.hole.W), H: mix(g.hole.H, to.hole.H)},
		x:    mix(g.x, to.x), y: mix(g.y, to.y),
	}
}

// tourBasis is what a tour's geometry is placed by: the target's box, the
// window's size and the panel's height.
type tourBasis struct {
	target ui.Rect
	w, h   float32
	panel  float32
}

// tourView is what the tour shows, kept on its scrim from frame to frame.
type tourView struct {
	started bool
	// step is the step shown; prev the one fading out, -1 for none.
	step, prev int
	// from is where the current step change started; shown is the
	// geometry of the frame before.
	from, shown tourGeometry
	// basis is what the frame before placed the geometry by.
	basis tourBasis
	// sizing is set while the content's height eases to a new step's.
	sizing bool
	// focusNext moves the focus to the forward button, as Back leaves.
	focusNext bool
}

// Tour walks state.Step through steps over a dimmed window. A spotlight cut
// out of the dim shows the step's target; one panel beside it, pointing at
// it, holds the step's title and body, the progress, Back, and Next (Done on
// the last step). On a step change the spotlight and the panel travel
// together to the next target while the content crossfades; under reduced
// motion they move at once and the content still fades.
//
// The tour is modal: presses outside its panel, the target's included, do
// nothing, and the focus stays in the panel. Right or Enter goes on, Left
// goes back (mirrored right to left), and Escape or the close button ends
// the tour; the focus then returns to where it was. Each step scrolls its
// target into view; a target still outside the window has its spotlight
// clamped to the window's edge. A Step outside steps shows nothing.
func Tour(c *ui.Context, state *TourState, steps []TourStep, o TourOptions) ui.Element {
	if state == nil {
		panic("overlay: Tour requires a caller-owned state")
	}
	anchored.Warm(c)
	if !state.Open || state.Step < 0 || state.Step >= len(steps) {
		return ui.Element{}
	}
	w, h := c.Size()
	width := o.Width
	if width <= 0 {
		width = 320
	}
	width = min(width, max(float32(0), w-2*theme.Space(2)))
	rtl := c.Direction() == ui.RightToLeft
	last := len(steps) - 1

	var panel ui.Element
	ui.Overlay(c, func() {
		scrim := ui.Box(c.Key(tourKey{state})).Absolute().Left(0).Top(0).Right(0).Bottom(0).Modal()
		// The page waits while the tour shows: a press outside the panel
		// lands here and does nothing, and no press moves the focus out of
		// the panel's buttons.
		scrim.KeepFocus().Clicked()
		fade := ui.Motion{Opacity: 0}
		m := theme.Motion()
		scrim.Transition(ui.ElementTransition{Duration: m.Reveal, Ease: m.Ease, Enter: &fade, Exit: &fade})
		view := ui.Local(scrim, tourViewKey{}, func() tourView { return tourView{prev: -1} })

		skip := func() {
			state.Open = false
			if o.OnDismiss != nil {
				o.OnDismiss()
			}
		}
		next := func() {
			if state.Step < last {
				state.Step++
				return
			}
			state.Open = false
			if o.OnFinish != nil {
				o.OnFinish()
			}
		}
		back := func() {
			if state.Step > 0 {
				state.Step--
				view.focusNext = state.Step == 0
			}
		}
		forwardKey, backKey := ui.KeyRight, ui.KeyLeft
		if rtl {
			forwardKey, backKey = backKey, forwardKey
		}
		escape := scrim.OverlayShortcut(0, ui.KeyEscape)
		forward := scrim.OverlayShortcut(0, forwardKey)
		enter := scrim.OverlayShortcut(0, ui.KeyEnter)
		backward := scrim.OverlayShortcut(0, backKey)
		switch {
		case escape:
			skip()
		case forward || enter:
			next()
		case backward:
			back()
		}
		if !state.Open {
			return
		}

		panel = ui.Column(c.Key(tourPanelKey{})).Absolute().Width(width).
			Padding(theme.Space(4)).Gap(theme.Space(4)).Role(ui.RoleDialog)
		panel = anchored.Surface(c, panel, theme.RadiusScale().LG).KeepFocus()
		step := steps[state.Step]
		var bounds ui.Rect
		if step.Anchor != nil {
			bounds = step.Anchor.Bounds(c)
		}
		// The placement reads the boxes the frame before laid out: while
		// the target, the window or the panel's height moves, build again
		// to place the panel by the boxes this frame lays out.
		basis := tourBasis{target: bounds, w: w, h: h, panel: panel.LayoutBounds().H}
		if basis != view.basis {
			view.basis = basis
			c.AnimationFrame()
		}
		target, side := tourPlace(step.Placement, bounds, width, basis.panel, w, h, rtl)
		travel := ui.Animation{Duration: m.Roll, Ease: m.EaseTravel}
		// Linear, so the two steps' opacities always add up to one.
		crossfade := ui.Animation{Duration: m.Roll, Ease: ui.Linear, NonMotion: true}
		p := scrim.AnimateTo(tourFadeKey{}, 1, crossfade)
		born := !view.started
		changed := !born && view.step != state.Step
		if born || changed {
			if changed {
				view.from = view.shown
				scrim.Jump(tourTravelKey{}, 0)
				p = tourSwapContent(view, state.Step, p)
				scrim.Jump(tourFadeKey{}, p)
				p = scrim.AnimateTo(tourFadeKey{}, 1, crossfade)
			} else {
				// Take the focus in the next frame: the focus this one
				// opens on is where it returns as the tour ends.
				view.focusNext = true
				c.AnimationFrame()
			}
			view.started, view.step = true, state.Step
			c.Announce(fmt.Sprintf(locale.M().StepOf, state.Step+1, len(steps)))
			if step.Anchor != nil {
				if e := c.Resolve(*step.Anchor); e.Valid() {
					e.ScrollIntoView()
				}
			}
		}
		if p >= 1 || view.prev >= len(steps) {
			view.prev = -1
		}
		shown := target
		if k := scrim.AnimateTo(tourTravelKey{}, 1, travel); k < 1 {
			shown = view.from.lerp(target, k)
		}
		view.shown = shown

		t := theme.Colors(c)
		hole, scrimColor := shown.hole, t.Backdrop
		radius := min(theme.RadiusScale().MD+tourPad, min(hole.W, hole.H)/2)
		// A dark window barely darkens under the scrim, so there the hole
		// also carries an outline.
		var ring ui.Color
		if c.Theme().Dark {
			ring = t.Foreground.Alpha(0.32)
		}
		scrim.Draw(func(p *ui.Painter, r ui.Rect) { paintSpotlight(p, r, hole, radius, scrimColor, ring) })

		panel.Left(shown.x).Top(shown.y).Label(step.Title).Description(step.Body)
		pr := ui.Rect{X: shown.x, Y: shown.y, W: width, H: basis.panel}
		if arrow, ok := tourArrowSide(pr, hole, side); ok {
			along := tourArrowAlong(pr, hole, arrow)
			panel.Draw(func(p *ui.Painter, r ui.Rect) {
				paintOutlinedArrow(p, r, arrow, t.Popover, t.Border, along)
			})
		}
		panel.Children(func() {
			ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space(2)).Children(func() {
				tourContent(c, view, steps, state.Step, p, changed)
				// Pull the button into the padding so its glyph lines up
				// with the text edge.
				mygoui.CloseButton(c, mygoui.CloseButtonOptions{Label: locale.M().Close, OnClick: skip}).
					Margin(-theme.Space(1), 0, 0, 0).MarginEnd(-theme.Space(1))
			})
			focus := view.focusNext && !born
			if focus {
				view.focusNext = false
			}
			tourFooter(c, state.Step, len(steps), back, next, focus)
		})
	})
	return panel
}

// tourSwapContent starts the crossfade from the step shown to step, from
// fade progress p, and returns the progress to go on from. Pressed again
// mid-change, the content the user sees keeps its opacity: a step that had
// not faded in yet is skipped, and a return to the step fading out fades it
// back in from where it is.
func tourSwapContent(view *tourView, step int, p float32) float32 {
	oldShown, newShown := tourFade(p)
	switch {
	case view.prev < 0 || p >= 1:
		view.prev, p = view.step, 0
	case oldShown > 0:
		// The old content is still leaving and the new one did not show.
	default:
		view.prev, p = view.step, tourFadeOut*(1-newShown)
	}
	if view.prev == step {
		oldShown, _ = tourFade(p)
		view.prev, p = -1, tourFadeOut+(1-tourFadeOut)*oldShown
	}
	return p
}

// tourFade returns the opacity of the content fading out and of the content
// fading in at crossfade progress p.
func tourFade(p float32) (out, in float32) {
	return 1 - min(1, p/tourFadeOut), max(0, (p-tourFadeOut)/(1-tourFadeOut))
}

// tourContent builds the step's title and body, over the previous step's
// while that fades out. The two share one cell, and the cell's height eases
// from the old content's to the new content's, so the panel neither jumps
// nor shows two texts at full strength.
func tourContent(c *ui.Context, view *tourView, steps []TourStep, at int, p float32, changed bool) {
	out, in := tourFade(p)
	if view.prev < 0 {
		in = 1
	}
	stack := ui.Box(c.Key(tourStackKey{})).Grow(1).Basis(0).MinWidth(0)
	var layer ui.Element
	stack.Children(func() {
		if view.prev >= 0 {
			tourLayer(c, steps[view.prev], view.prev, true).Opacity(out)
		}
		layer = tourLayer(c, steps[at], at, false).Opacity(in)
	})
	shown, want := stack.LayoutBounds().H, layer.LayoutBounds().H
	if changed && !view.sizing {
		view.sizing = true
		stack.Jump(tourHeightKey{}, shown)
	}
	if !view.sizing {
		return
	}
	if want <= 0 {
		// The new content lays out at the end of this frame: hold the
		// height shown until then.
		want = shown
	}
	height := stack.AnimateTo(tourHeightKey{}, want, ui.Animation{Duration: theme.Motion().Reveal, Ease: theme.Motion().EaseTravel})
	if height == want && layer.LayoutBounds().H > 0 {
		view.sizing = false
		return
	}
	stack.Height(height).Clip()
}

// tourLayer builds one step's title and body; leaving, it overlays the
// incoming layer and is hidden from assistive technology.
func tourLayer(c *ui.Context, s TourStep, step int, leaving bool) ui.Element {
	l := ui.Column(c.Key(tourLayerKey{step})).FillWidth().MinWidth(0).Gap(theme.Space(1))
	if leaving {
		l.Absolute().Left(0).Top(0).Right(0).AccessHidden(true)
	}
	l.Children(func() {
		if s.Title != "" {
			style.Text(ui.Text(c, s.Title), c, theme.SM()).FontWeight(theme.WeightSemibold).
				Role(ui.RoleHeading).Level(2)
		}
		if s.Body != "" {
			style.Text(ui.Text(c, s.Body), c, theme.SM()).TextColor(theme.Colors(c).MutedForeground)
		}
	})
	return l
}

// tourFooter builds the progress and the Back and Next buttons. Next keeps
// its place and identity from step to step, so the focus stays on it; it
// takes the focus as the tour starts and as Back leaves with the first step.
func tourFooter(c *ui.Context, at, count int, back, next func(), focus bool) {
	m := locale.M()
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Space(2)).Children(func() {
		tourProgress(c, at, count).Grow(1).Basis(0).MinWidth(0)
		if at > 0 {
			mygoui.Button(c.Key(tourBackKey{}), mygoui.ButtonOptions{
				Label: m.Back, Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall, OnClick: back,
			})
		}
		label := m.Next
		if at == count-1 {
			label = m.Done
		}
		forward := mygoui.Button(c.Key(tourNextKey{}), mygoui.ButtonOptions{
			Label: label, Size: theme.ButtonSmall, OnClick: next,
		})
		if focus {
			forward.Focus()
		}
	})
}

// tourProgress shows where the tour is: a dot per step, the current one a
// wider bar in the primary colour, or "2 / 12" for a long tour.
func tourProgress(c *ui.Context, at, count int) ui.Element {
	t := theme.Colors(c)
	label := fmt.Sprintf(locale.M().StepOf, at+1, count)
	if count > tourDots {
		text := strconv.Itoa(at+1) + " / " + strconv.Itoa(count)
		return style.Text(ui.Text(c, text), c, theme.XS()).TextColor(t.MutedForeground).
			ContentTransition(ui.ContentTransition{}).Label(label)
	}
	row := ui.Row(c).AlignItems(ui.Center).Gap(theme.Space(1)).Role(ui.RoleImage).Label(label)
	row.Children(func() {
		for i := range count {
			dot := ui.Box(c.Key(tourDotKey{i})).Shrink(0).Height(6).Radius(3)
			width, color := float32(6), t.MutedForeground.Alpha(0.3)
			if i == at {
				width, color = 16, t.Primary
			}
			dot.Width(dot.AnimateTo(tourDotKey{}, width, ui.Animation{Duration: theme.Motion().Reveal, Ease: theme.Motion().EaseTravel})).
				Background(motion.Color(dot, tourDotKey{}, color))
		}
	})
	return row
}

// tourPlace returns where a step shows: its padded target clamped to the
// window, or none, and the panel beside it on the side it fits on, inside
// the window; a step without a target centres the panel.
func tourPlace(placement mygoui.Placement, r ui.Rect, pw, ph, w, h float32, rtl bool) (tourGeometry, mygoui.Side) {
	hole, ok := tourHole(r, w, h)
	if !ok {
		return tourGeometry{hole: ui.Rect{X: w / 2, Y: h / 2}, x: (w - pw) / 2, y: (h - ph) / 2}, mygoui.SideAuto
	}
	gap := placement.Offset
	if gap == 0 {
		gap = arrowOffset
	}
	margin := float32(theme.Space(2))
	room := map[mygoui.Side]float32{
		mygoui.SideBottom: h - margin - (hole.Y + hole.H + gap) - ph,
		mygoui.SideTop:    hole.Y - gap - margin - ph,
		mygoui.SideRight:  w - margin - (hole.X + hole.W + gap) - pw,
		mygoui.SideLeft:   hole.X - gap - margin - pw,
	}
	side := resolveSide(placement, mygoui.SideBottom)
	if room[side] < 0 {
		best := side
		for _, other := range tourFallbacks(side) {
			if room[other] >= 0 {
				best = other
				break
			}
			if room[other] > room[best] {
				best = other
			}
		}
		side = best
	}
	align := placement.Align
	if rtl && align != mygoui.AlignCenter {
		align = mygoui.AlignStart + mygoui.AlignEnd - align
	}
	var x, y float32
	switch side {
	case mygoui.SideTop, mygoui.SideBottom:
		x = tourAlign(hole.X, hole.W, pw, align)
		y = hole.Y + hole.H + gap
		if side == mygoui.SideTop {
			y = hole.Y - gap - ph
		}
	default:
		y = tourAlign(hole.Y, hole.H, ph, placement.Align)
		x = hole.X + hole.W + gap
		if side == mygoui.SideLeft {
			x = hole.X - gap - pw
		}
	}
	x = max(margin, min(x, w-margin-pw))
	y = max(margin, min(y, h-margin-ph))
	return tourGeometry{hole: hole, x: x, y: y}, side
}

// tourFallbacks is the order other sides are tried in: the opposite one
// first, then the other axis.
func tourFallbacks(side mygoui.Side) []mygoui.Side {
	switch side {
	case mygoui.SideTop:
		return []mygoui.Side{mygoui.SideBottom, mygoui.SideRight, mygoui.SideLeft}
	case mygoui.SideLeft:
		return []mygoui.Side{mygoui.SideRight, mygoui.SideBottom, mygoui.SideTop}
	case mygoui.SideRight:
		return []mygoui.Side{mygoui.SideLeft, mygoui.SideBottom, mygoui.SideTop}
	}
	return []mygoui.Side{mygoui.SideTop, mygoui.SideRight, mygoui.SideLeft}
}

// tourAlign places a panel of size along an edge from start of length.
func tourAlign(start, length, size float32, align mygoui.Align) float32 {
	switch align {
	case mygoui.AlignStart:
		return start
	case mygoui.AlignEnd:
		return start + length - size
	}
	return start + (length-size)/2
}

// tourHole is the spotlight around a target's box, clamped to the window;
// it reports false for no target or one wholly outside the window.
func tourHole(r ui.Rect, w, h float32) (ui.Rect, bool) {
	if r.W <= 0 || r.H <= 0 {
		return ui.Rect{}, false
	}
	x0, y0 := max(0, r.X-tourPad), max(0, r.Y-tourPad)
	x1, y1 := min(w, r.X+r.W+tourPad), min(h, r.Y+r.H+tourPad)
	if x1 <= x0 || y1 <= y0 {
		return ui.Rect{}, false
	}
	return ui.Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}, true
}

// tourArrowSide returns the panel's edge facing the spotlight: the side the
// panel was placed on, or, mid-travel, the edge closest to the spotlight. A
// panel without a spotlight, or over it, has no arrow.
func tourArrowSide(panel, hole ui.Rect, placed mygoui.Side) (mygoui.Side, bool) {
	if hole.W < 1 || hole.H < 1 || panel.H <= 0 {
		return 0, false
	}
	gaps := map[mygoui.Side]float32{
		mygoui.SideBottom: panel.Y - (hole.Y + hole.H),
		mygoui.SideTop:    hole.Y - (panel.Y + panel.H),
		mygoui.SideRight:  panel.X - (hole.X + hole.W),
		mygoui.SideLeft:   hole.X - (panel.X + panel.W),
	}
	if gaps[placed] >= arrowHeight {
		return placed, true
	}
	best, gap := mygoui.SideAuto, float32(math.Inf(-1))
	for _, side := range []mygoui.Side{mygoui.SideBottom, mygoui.SideTop, mygoui.SideRight, mygoui.SideLeft} {
		if gaps[side] > gap {
			best, gap = side, gaps[side]
		}
	}
	return best, gap >= arrowHeight
}

// tourArrowAlong is where along the panel's edge the arrow points: at the
// spotlight's centre, kept half an arrow and its corner radius inside the
// edge's ends.
func tourArrowAlong(panel, hole ui.Rect, side mygoui.Side) float32 {
	size, center, start := panel.W, hole.X+hole.W/2, panel.X
	if side == mygoui.SideLeft || side == mygoui.SideRight {
		size, center, start = panel.H, hole.Y+hole.H/2, panel.Y
	}
	margin := (arrowWidth/2 + theme.RadiusScale().LG) / size
	if margin >= 0.5 {
		return 0.5
	}
	return min(1-margin, max(margin, (center-start)/size))
}

// paintSpotlight dims r around a rounded hole. It strokes a rounded
// rectangle whose band reaches past every corner of r, so the hole is the
// stroke's inside: one rounded-rectangle operation, not a path.
func paintSpotlight(p *ui.Painter, r, hole ui.Rect, radius float32, color, ring ui.Color) {
	if hole.W < 1 || hole.H < 1 {
		p.Fill(r, color, 0)
		return
	}
	band := 2 * (r.W + r.H)
	outer := ui.Rect{X: hole.X - band, Y: hole.Y - band, W: hole.W + 2*band, H: hole.H + 2*band}
	p.Clip(r, 0, func() { p.Stroke(outer, color, radius+band, band) })
	if ring.A > 0 {
		p.Stroke(hole, ring, radius, 1)
	}
}
