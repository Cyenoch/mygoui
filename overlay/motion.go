package overlay

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// COSS specifies 200 ms for dialogs and sheets and 450 ms for drawers.
// Popovers and tooltips use a 150 ms native transition.
const (
	dialogMotionDuration = 200 * time.Millisecond
	drawerMotionDuration = 450 * time.Millisecond
	popupMotionDuration  = 150 * time.Millisecond
)

// Popovers and tooltips fade from 98% around the anchor-facing edge.
const (
	popupEnterScale = 0.98
	popupExitScale  = 0.98
)

// Allocate modal keys before DialogBase: its backdrop reads click/shortcut
// state before invoking the styling callback. Rekeying there loses input and
// focus restoration. Anchored panels key before building their content.
type (
	dialogKey         struct{ open *bool }
	popoverKey        struct{ open *bool }
	tooltipKey        struct{ anchor uint64 }
	modalScrimKey     struct{ open *bool }
	edgeEnteredKey    struct{}
	edgeOffsetKey     struct{}
	dialogStateKey    struct{ open *bool }
	dialogEnteredKey  struct{}
	dialogProgressKey struct{}
)

type dialogState struct {
	open           bool
	resuming       bool
	exitEnabled    bool
	closedAt       time.Time
	reopenedAt     time.Time
	progress       float32
	closedProgress float32
	returnProgress float32
}

func motionFraction(elapsed time.Duration) float32 {
	return ui.EaseInOut(max(0, min(float32(elapsed)/float32(dialogMotionDuration), 1)))
}

func (s *dialogState) observe(open, reduced bool, now time.Time) {
	if s.open && !open {
		s.closedAt = now
		s.closedProgress = s.progress
		if reduced || !s.exitEnabled {
			s.closedAt = time.Time{}
		}
	} else if !s.open && open {
		s.resuming = !s.closedAt.IsZero() && now.Sub(s.closedAt) < dialogMotionDuration
		if s.resuming {
			s.reopenedAt = now
			k := motionFraction(now.Sub(s.closedAt))
			s.returnProgress = s.closedProgress + (1-s.closedProgress)*k
		}
	}
	s.open = open
}

// Keep this state outside DialogBase, whose panel is discarded on close.
func responsiveDialogState(c *ui.Context, open *bool) *dialogState {
	var state *dialogState
	ui.Overlay(c, func() {
		owner := ui.Box(c.Key(dialogStateKey{open})).Absolute().Size(0, 0).PassThrough().Role(ui.RoleNone)
		state = ui.Local(owner, dialogStateKey{}, func() dialogState { return dialogState{} })
	})
	state.observe(*open, theme.ReduceMotion(c), c.Now())
	return state
}

// modalScrim leaves the panel's backdrop transparent so its inert exit copy
// can move independently from the fading scrim.
func modalScrim(c *ui.Context, open *bool, variant ModalBackdrop, duration time.Duration, ease ui.Easing) ui.Element {
	if !*open {
		return ui.Element{}
	}
	var scrim ui.Element
	ui.Overlay(c, func() {
		scrim = ui.Box(c.Key(modalScrimKey{open})).Absolute().Left(0).Top(0).Right(0).Bottom(0).PassThrough()
		if !theme.ReduceMotion(c) {
			fade := ui.Motion{Opacity: 0}
			scrim.Transition(ui.ElementTransition{Duration: duration, Ease: ease, Enter: &fade, Exit: &fade})
		}
		paintBackdrop(c, scrim, variant)
	})
	return scrim
}

func edgeOffset(side mygoui.Side, distance float32) (float32, float32) {
	switch side {
	case mygoui.SideLeft:
		return -distance, 0
	case mygoui.SideTop:
		return 0, -distance
	case mygoui.SideBottom:
		return 0, distance
	default:
		return distance, 0
	}
}

func edgeEnter(c *ui.Context, panel ui.Element, side mygoui.Side, distance float32, duration time.Duration, ease ui.Easing) {
	entered := ui.Local(panel, edgeEnteredKey{}, func() bool { return theme.ReduceMotion(c) })
	from := float32(0)
	if !*entered {
		from = distance
		*entered = true
		c.AnimationFrame()
	}
	offset := panel.AnimateTo(edgeOffsetKey{}, from, ui.Animation{Duration: duration, Ease: ease})
	x, y := edgeOffset(side, offset)
	panel.Translate(x, y)
}

func edgeExit(c *ui.Context, back ui.Element, side mygoui.Side, distance float32, duration time.Duration, ease ui.Easing) ui.Element {
	if theme.ReduceMotion(c) {
		return back
	}
	x, y := edgeOffset(side, distance)
	enter := ui.Motion{Opacity: 0}
	exit := ui.Motion{X: x, Y: y, Opacity: 0}
	return back.Transition(ui.ElementTransition{Duration: duration, Ease: ease, Enter: &enter, Exit: &exit})
}

// responsiveDialogMotion moves the panel on entry and its retained backdrop
// subtree on exit. The separate scrim can then fade without moving.
func responsiveDialogMotion(c *ui.Context, back, panel ui.Element, state *dialogState, w, h float32) ui.Element {
	entered := ui.Local(panel, dialogEnteredKey{}, func() bool { return theme.ReduceMotion(c) })
	if theme.ReduceMotion(c) {
		*entered = true
		state.progress = 0
		state.exitEnabled = false
		return back
	}
	target := float32(0)
	if !*entered {
		target = 1
		if state.resuming {
			target = state.returnProgress
		}
		*entered = true
		if target != 0 {
			c.AnimationFrame()
		}
	}
	progress := panel.AnimateTo(dialogProgressKey{}, target, ui.Animation{Duration: dialogMotionDuration, Ease: ui.EaseInOut})
	state.progress = progress
	exit := ui.Motion{Opacity: 0}
	if w < 640 {
		offset := theme.Space16 * progress
		exit.Y = theme.Space16 - offset
		panel.Translate(0, offset)
	} else {
		scale := 1 - 0.02*progress
		panel.Scale(scale, scale)
		exit.Scale = 0.98 / scale
		if r := panel.LayoutBounds(); w > 0 && h > 0 && r.W > 0 && r.H > 0 {
			back.TransformOrigin((r.X+r.W/2)/w, (r.Y+r.H/2)/h)
		}
		if state.resuming && c.Now().Equal(state.reopenedAt) {
			retainedScale := (1 - 0.02*state.returnProgress) / (1 - 0.02*state.closedProgress)
			back.Scale(1/retainedScale, 1/retainedScale)
		}
	}
	enter := ui.Motion{Opacity: 0}
	state.exitEnabled = true
	transition := ui.ElementTransition{
		Duration: dialogMotionDuration, Ease: ui.EaseInOut, Enter: &enter, Exit: &exit,
	}
	if state.resuming && c.Now().Equal(state.reopenedAt) {
		// Settle the retained backdrop's position before its next snapshot;
		// its returning panel already carries the visible offset. MyGo uses
		// the first call for the return and the last for a later exit.
		short := transition
		short.Duration = time.Nanosecond
		back.Transition(short)
		transition.Exit = nil
		state.exitEnabled = false
	}
	return back.Transition(transition)
}

// warmOverlay builds the overlay layer now. An element's Enter applies only
// in a parent the frame before had, and the layer is rebuilt each frame: a
// component must ask for it every frame, or the first surface a window shows
// after a frame without overlays would appear without its arrival. It adds
// no element: an empty overlay only makes sure the layer has been built.
func warmOverlay(c *ui.Context) {
	ui.Overlay(c, func() {})
}

// surfaceMotion fades an overlay and keeps an inert exit snapshot.
func surfaceMotion(c *ui.Context, e ui.Element, d time.Duration) ui.Element {
	if theme.ReduceMotion(c) {
		return e
	}
	motion := ui.Motion{Opacity: 0}
	return e.Transition(ui.ElementTransition{
		Duration: d,
		Ease:     ui.EaseOut,
		Enter:    &motion,
		Exit:     &motion,
	})
}

// dialogMotion fades the native dialog backdrop, and so the panel it holds.
// DialogBase builds neither once *open is false, so the backdrop is the one
// element a Transition can keep as the dialog goes: its exit copy holds the
// last paint of the caller's body inside the panel.
//
// Modal and BottomSheet use this fade; responsive dialogs and edge sheets
// move their panels separately from the scrim.
func dialogMotion(c *ui.Context, back ui.Element, d time.Duration) ui.Element {
	return surfaceMotion(c, back, d)
}

// popupMotion fades an anchored panel and zooms it around
// the anchor point. The caller keys the panel first, before building its
// content, so that the content follows the panel's identity and the side it
// settled on can be read from its committed box. The panel is a child of the
// overlay layer, which warmOverlay keeps alive, so its Enter applies: the
// transition drives its own frames and nothing is left at a partial scale.
func popupMotion(c *ui.Context, panel ui.Element, side mygoui.Side) ui.Element {
	if theme.ReduceMotion(c) {
		return panel
	}
	panel.TransformOrigin(popupOrigin(side))
	enter := ui.Motion{Opacity: 0, Scale: popupEnterScale}
	exit := ui.Motion{Opacity: 0, Scale: popupExitScale}
	return panel.Transition(ui.ElementTransition{
		Duration: popupMotionDuration,
		Ease:     ui.EaseOut,
		Enter:    &enter,
		Exit:     &exit,
	})
}

// popupOrigin returns the transform origin for a panel on side: the point that faces the
// anchor grows and shrinks in place.
func popupOrigin(side mygoui.Side) (x, y float32) {
	switch side {
	case mygoui.SideTop:
		return 0.5, 1
	case mygoui.SideLeft:
		return 1, 0.5
	case mygoui.SideRight:
		return 0, 0.5
	}
	return 0.5, 0
}
