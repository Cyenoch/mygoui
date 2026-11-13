package overlay

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
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

// motionFraction is how far a dialog has travelled elapsed into its motion.
func motionFraction(elapsed time.Duration) float32 {
	m := theme.Motion()
	return m.EaseTravel(max(0, min(float32(elapsed)/float32(m.Dialog), 1)))
}

func (s *dialogState) observe(open, reduced bool, now time.Time) {
	if s.open && !open {
		s.closedAt = now
		s.closedProgress = s.progress
		if reduced || !s.exitEnabled {
			s.closedAt = time.Time{}
		}
	} else if !s.open && open {
		s.resuming = !s.closedAt.IsZero() && now.Sub(s.closedAt) < theme.Motion().Dialog
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
		if theme.ReduceMotion(c) {
			anchored.Still(scrim)
		} else {
			fade := ui.Motion{Opacity: 0}
			scrim.Transition(ui.ElementTransition{Duration: duration, Ease: ease, Enter: &fade, Exit: &fade})
		}
		paintBackdrop(c, scrim, variant)
	})
	return scrim
}

// edgeOffset is the translation that moves a panel distance out through its
// edge. Left and right follow the layout direction, as the panel's place does.
func edgeOffset(c *ui.Context, side mygoui.Side, distance float32) (float32, float32) {
	if c.Direction() == ui.RightToLeft {
		switch side {
		case mygoui.SideLeft:
			side = mygoui.SideRight
		case mygoui.SideRight:
			side = mygoui.SideLeft
		}
	}
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
	x, y := edgeOffset(c, side, offset)
	panel.Translate(x, y)
}

func edgeExit(c *ui.Context, back ui.Element, side mygoui.Side, distance float32, duration time.Duration, ease ui.Easing) ui.Element {
	if theme.ReduceMotion(c) {
		return anchored.Still(back)
	}
	x, y := edgeOffset(c, side, distance)
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
		return anchored.Still(back)
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
	m := theme.Motion()
	// shrink is how much smaller than itself a wide dialog enters.
	from := theme.Effects().PopupScale
	shrink := 1 - from
	progress := panel.AnimateTo(dialogProgressKey{}, target, ui.Animation{Duration: m.Dialog, Ease: m.EaseTravel})
	state.progress = progress
	// A dialog reopened during its exit returns from where its retained
	// backdrop shows it, which holds the exit motion's share of the offset;
	// the new panel adds only the offset it had when it closed.
	reopening := state.resuming && c.Now().Equal(state.reopenedAt)
	exit := ui.Motion{Opacity: 0}
	if w < 640 {
		offset := theme.Space(4) * progress
		exit.Y = theme.Space(4) - offset
		if reopening {
			offset = theme.Space(4) * state.closedProgress
		}
		panel.Translate(0, offset)
	} else {
		scale := 1 - shrink*progress
		panel.Scale(scale, scale)
		if scale != 1 {
			// A scaled panel's content takes a CPU-rasterized clip mask per
			// frame wherever a clip inside it meets the transform, which
			// outlasts the whole motion in a browser; scale one bitmap of
			// the panel instead.
			panel.Cache()
		}
		exit.Scale = from / scale
		if r := panel.LayoutBounds(); w > 0 && h > 0 && r.W > 0 && r.H > 0 {
			back.TransformOrigin((r.X+r.W/2)/w, (r.Y+r.H/2)/h)
		}
		if reopening {
			retainedScale := (1 - shrink*state.returnProgress) / (1 - shrink*state.closedProgress)
			back.Scale(1/retainedScale, 1/retainedScale)
		}
	}
	enter := ui.Motion{Opacity: 0}
	state.exitEnabled = true
	transition := ui.ElementTransition{
		Duration: m.Dialog, Ease: m.EaseTravel, Enter: &enter, Exit: &exit,
	}
	if reopening {
		// Settle the retained backdrop's position at once. MyGo applies an
		// element's last Transition of the frame, so this frame uses only
		// the short one and takes no exit until the next build.
		transition.Duration = time.Nanosecond
		transition.Exit = nil
		state.exitEnabled = false
	}
	return back.Transition(transition)
}

// surfaceMotion fades an overlay over d with ease and keeps an inert exit
// snapshot. For a modal surface e is the native dialog backdrop, and so the
// panel it holds: DialogBase builds neither once *open is false, so the
// backdrop is the one element a Transition can keep as the dialog goes, and
// its exit copy holds the last paint of the caller's body inside the panel.
// Responsive dialogs and edge sheets move their panels separately from the
// scrim instead.
func surfaceMotion(c *ui.Context, e ui.Element, d time.Duration, ease ui.Easing) ui.Element {
	if theme.ReduceMotion(c) {
		return anchored.Still(e)
	}
	motion := ui.Motion{Opacity: 0}
	return e.Transition(ui.ElementTransition{
		Duration: d,
		Ease:     ease,
		Enter:    &motion,
		Exit:     &motion,
	})
}
