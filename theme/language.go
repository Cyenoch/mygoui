package theme

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
)

// Theme is a complete design language as one value: both palettes, the
// radius base, type, spacing, elevation, motion and effects. Default returns
// the built-in one; change what differs and publish it with Use.
//
// Preferences resolve over the published theme: a preference the host sets
// replaces the theme's value, an unset one inherits it.
type Theme struct {
	// Light and Dark are the complete palettes, derived roles included.
	Light, Dark Palette
	// Radius is the base radius; the radius scale is derived from it.
	Radius     float32
	Typography Typography
	Spacing    Spacing
	Elevation  Elevation
	Motion     MotionTokens
	Effects    EffectTokens
}

// Typography is the type of a theme: families, the base size and the roles.
// An empty family selects the native stack. Role sizes are unscaled DIPs.
type Typography struct {
	Family, HeadingFamily, MonoFamily string
	// Size is the base text size the roles are scaled against 16 DIPs.
	Size float32

	XS, SM, Base, LG, XL   TextStyle
	H1, H2, H3, H4, H5, H6 TextStyle
	Code, Label, HelpText  TextStyle
}

// Spacing is the spacing scale: Space(n) is n units.
type Spacing struct {
	Unit float32
}

// Elevation holds the shadow layers of each appearance.
type Elevation struct{ Light, Dark Shadows }

// Shadows are the shadow layers of one appearance, lowest first.
type Shadows struct {
	// Surface is a card or resting surface, Field an outlined control,
	// Overlay a popup or dialog and Tooltip a tooltip.
	Surface, Field, Overlay, Tooltip []Shadow
}

// MotionTokens are the durations and easings components share. A change
// without a token here is timed by its component.
//
// A zero duration shows its change at once; Use rejects a negative one.
// Motion reports a zero duration as one nanosecond, the shortest MyGo
// animates, so every reader lands on the frame after the change; Current
// returns the value as written.
type MotionTokens struct {
	// State is a control changing state in place: a fill, border or ink
	// fading, a mark drawing in, a thumb moving, content swapping with a
	// spinner.
	State time.Duration
	// Reveal is content opening or closing, a panel folding and an indicator
	// moving to another item.
	Reveal time.Duration
	// Roll is a changing number turning its digits.
	Roll time.Duration
	// Shift is rows and tabs moving aside for a dragged one.
	Shift time.Duration
	// Page is one page sliding to the next.
	Page time.Duration
	// Popup is an anchored popup or floating panel appearing and leaving.
	Popup time.Duration
	// Dialog is a dialog, modal or edge sheet and its scrim.
	Dialog time.Duration
	// Drawer is a drawer or bottom sheet sliding from its edge.
	Drawer time.Duration

	// Ease is the curve of a change that has no other: it decelerates into
	// the new state.
	Ease ui.Easing
	// EaseTravel is the curve of something moving between two resting
	// places: an indicator, a dialog, a highlight.
	EaseTravel ui.Easing
	// EaseDrawer is the curve of the drawer family.
	EaseDrawer ui.Easing
}

// EffectTokens are the parameters of effects several components share.
type EffectTokens struct {
	// HoverOpacity is the share of its alpha a solid fill keeps under the
	// pointer, wherever a component derives the hover from the fill as it
	// draws: buttons, status solids, the toast action, native controls, and
	// the palette roles a preference moves. The hover roles a theme's
	// palettes carry (DefaultHover, AccentHover and the status hovers) are
	// colors of that palette and are used as written.
	HoverOpacity float32
	// PopupScale is the scale a popup or wide dialog enters from.
	PopupScale float32
	// BackdropBlur is the blur of a modal backdrop in DIPs.
	BackdropBlur float32
	// FadeEdge is the depth in DIPs over which clipped content fades.
	FadeEdge float32
	// SpinnerPeriod is one turn of a spinner. Use rejects a period under a
	// millisecond.
	SpinnerPeriod time.Duration
	// SweepPeriod is one pass of an indeterminate band, SweepBand its share
	// of the track, between 0 and 1. Use rejects a period under a
	// millisecond.
	SweepPeriod time.Duration
	SweepBand   float32
	// LoadingDelay keeps a region's loading indicator back, so a load that
	// ends sooner shows none. Zero shows it at once; Use rejects a negative
	// delay.
	LoadingDelay time.Duration
	// RefreshOpacity is the opacity content keeps while it refreshes.
	RefreshOpacity float32
	// ShimmerPeriod is one loop of a loading placeholder's shimmer or
	// pulse; zero or less keeps placeholders still.
	ShimmerPeriod time.Duration
	// ShimmerPeak is how far the shimmer's band mixes a placeholder
	// towards white in light appearance, ShimmerPeakDark in dark.
	ShimmerPeak, ShimmerPeakDark float32
	// ThumbPressScale is the size a thumb grows to while the pointer holds
	// it, as a slider's or a colour area's, relative to its rest size.
	ThumbPressScale float32
}

// defaultHoverOpacity is also what the built-in palettes derive their hover
// roles with.
const defaultHoverOpacity = 0.9

// Default returns the built-in theme. The caller owns the value.
func Default() Theme {
	heading := func(size, line float32) TextStyle {
		return TextStyle{Size: size, LineHeight: line, Weight: WeightSemibold, Tracking: TrackingTight, Heading: true}
	}
	body := func(size, line float32) TextStyle {
		return TextStyle{Size: size, LineHeight: line, Weight: WeightNormal}
	}
	// A raised surface takes its light from above: the inset edge is a dark
	// lower edge in light appearance and a light upper edge in dark. An
	// inset shadow shows on the side opposite its offset.
	lightEdge := Shadow{Y: -1, Color: ui.RGBA(0, 0, 0, 0.04), Inset: true}
	darkEdge := Shadow{Y: 1, Color: ui.RGBA(255, 255, 255, 0.06), Inset: true}
	drop := Shadow{Y: 1, Blur: 2, Color: ui.RGBA(0, 0, 0, 0.05)}
	// Popups use two soft drops and tooltips two smaller ones, over the
	// surface's edge.
	large := []Shadow{{Y: 10, Blur: 15, Spread: -3, Color: ui.RGBA(0, 0, 0, 0.05)}, {Y: 4, Blur: 6, Spread: -4, Color: ui.RGBA(0, 0, 0, 0.05)}}
	medium := []Shadow{{Y: 4, Blur: 6, Spread: -1, Color: ui.RGBA(0, 0, 0, 0.05)}, {Y: 2, Blur: 4, Spread: -2, Color: ui.RGBA(0, 0, 0, 0.05)}}
	return Theme{
		Light:  buildPalette(false),
		Dark:   buildPalette(true),
		Radius: 10,
		Typography: Typography{
			Size: 16,
			XS:   body(12, 16), SM: body(14, 20), Base: body(16, 24), LG: body(18, 28), XL: body(20, 28),
			H1: heading(36, 40), H2: heading(30, 36), H3: heading(24, 32),
			H4: heading(20, 28), H5: heading(18, 28), H6: heading(16, 24),
			Code:     TextStyle{Size: 14, LineHeight: 20, Weight: WeightNormal, Mono: true},
			Label:    TextStyle{Size: 14, LineHeight: 20, Weight: WeightMedium},
			HelpText: body(12, 16),
		},
		Spacing: Spacing{Unit: SpaceUnit},
		Elevation: Elevation{
			Light: Shadows{
				Surface: []Shadow{drop, lightEdge},
				Field:   []Shadow{drop, {Y: -1, Color: ui.RGBA(0, 0, 0, 0.06), Inset: true}},
				Overlay: append(slices.Clone(large), lightEdge),
				Tooltip: append(slices.Clone(medium), lightEdge),
			},
			Dark: Shadows{
				Surface: []Shadow{darkEdge},
				Field:   []Shadow{darkEdge},
				Overlay: append(slices.Clone(large), darkEdge),
				Tooltip: append(slices.Clone(medium), darkEdge),
			},
		},
		Motion: MotionTokens{
			State:      150 * time.Millisecond,
			Reveal:     200 * time.Millisecond,
			Roll:       300 * time.Millisecond,
			Shift:      150 * time.Millisecond,
			Page:       250 * time.Millisecond,
			Popup:      150 * time.Millisecond,
			Dialog:     200 * time.Millisecond,
			Drawer:     450 * time.Millisecond,
			Ease:       ui.EaseOut,
			EaseTravel: ui.EaseInOut,
			EaseDrawer: ui.CubicBezier(0.32, 0.72, 0, 1),
		},
		Effects: EffectTokens{
			HoverOpacity:    defaultHoverOpacity,
			PopupScale:      0.98,
			BackdropBlur:    4,
			FadeEdge:        24,
			SpinnerPeriod:   time.Second,
			SweepPeriod:     1500 * time.Millisecond,
			SweepBand:       0.4,
			LoadingDelay:    150 * time.Millisecond,
			RefreshOpacity:  0.55,
			ShimmerPeriod:   2 * time.Second,
			ShimmerPeak:     0.64,
			ShimmerPeakDark: 0.04,
			ThumbPressScale: 1.2,
		},
	}
}

// Use publishes t as the design language and keeps the published
// preferences over it. Call it on the UI thread when the theme changes, then
// invalidate affected windows; every call republishes, so it does not belong
// in a view build. A nil easing takes the built-in curve.
//
// Use publishes nothing and returns an error naming the field when a
// duration is negative or a spinner or sweep period is under a millisecond:
// components divide by the periods and MyGo panics on a negative duration.
func Use(t Theme) error {
	if err := t.validate(); err != nil {
		return err
	}
	t = t.clone()
	builtin := Default().Motion
	if t.Motion.Ease == nil {
		t.Motion.Ease = builtin.Ease
	}
	if t.Motion.EaseTravel == nil {
		t.Motion.EaseTravel = builtin.EaseTravel
	}
	if t.Motion.EaseDrawer == nil {
		t.Motion.EaseDrawer = builtin.EaseDrawer
	}
	publish(t, snapshot().prefs)
	return nil
}

// instant is the shortest duration MyGo animates. A zero duration would
// select its default spring instead.
const instant = time.Nanosecond

func (t Theme) validate() error {
	m, fx := t.Motion, t.Effects
	for _, d := range []struct {
		name  string
		value time.Duration
	}{
		{"Motion.State", m.State}, {"Motion.Reveal", m.Reveal}, {"Motion.Roll", m.Roll},
		{"Motion.Shift", m.Shift}, {"Motion.Page", m.Page}, {"Motion.Popup", m.Popup},
		{"Motion.Dialog", m.Dialog}, {"Motion.Drawer", m.Drawer},
		{"Effects.LoadingDelay", fx.LoadingDelay},
	} {
		if d.value < 0 {
			return fmt.Errorf("theme: %s is %v; a duration must not be negative", d.name, d.value)
		}
	}
	for _, d := range []struct {
		name  string
		value time.Duration
	}{{"Effects.SpinnerPeriod", fx.SpinnerPeriod}, {"Effects.SweepPeriod", fx.SweepPeriod}} {
		if d.value < time.Millisecond {
			return fmt.Errorf("theme: %s is %v; a period must be at least a millisecond", d.name, d.value)
		}
	}
	return nil
}

// resolved is m as components read it: a zero duration becomes instant.
func (m MotionTokens) resolved() MotionTokens {
	for _, d := range []*time.Duration{&m.State, &m.Reveal, &m.Roll, &m.Shift, &m.Page, &m.Popup, &m.Dialog, &m.Drawer} {
		if *d == 0 {
			*d = instant
		}
	}
	return m
}

// Current returns the published theme, without the preferences resolved
// over it. The caller owns the value.
func Current() Theme { return snapshot().theme.clone() }

// Motion returns the published motion tokens as components read them.
func Motion() MotionTokens { return snapshot().motion }

// Effects returns the published effect tokens.
func Effects() EffectTokens { return snapshot().theme.Effects }

// clone detaches the caller-owned shadow layers.
func (t Theme) clone() Theme {
	for _, s := range []*Shadows{&t.Elevation.Light, &t.Elevation.Dark} {
		s.Surface = slices.Clone(s.Surface)
		s.Field = slices.Clone(s.Field)
		s.Overlay = slices.Clone(s.Overlay)
		s.Tooltip = slices.Clone(s.Tooltip)
	}
	return t
}

func (e Elevation) mode(dark bool) Shadows {
	if dark {
		return e.Dark
	}
	return e.Light
}
