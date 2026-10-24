package mygoui

import (
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// ButtonVariant is HeroUI's button color variant.
type ButtonVariant uint8

const (
	ButtonPrimary ButtonVariant = iota
	ButtonSecondary
	ButtonTertiary
	ButtonOutline
	ButtonGhost
	ButtonDanger
	ButtonDangerSoft
)

// ButtonAlign places button content.
type ButtonAlign uint8

const (
	ButtonAlignCenter ButtonAlign = iota
	ButtonAlignStart
	ButtonAlignEnd
)

// ButtonOptions configures a native button. Label is always the accessible
// name, including icon-only layouts and custom content.
type ButtonOptions struct {
	Label     string
	Icon      *ui.SVG
	EndIcon   *ui.SVG
	Variant   ButtonVariant
	Size      theme.ButtonSize
	IconOnly  bool
	FullWidth bool
	Loading   bool
	Disabled  bool
	OnClick   func()
	Handle    *ui.Handle
	Tooltip   string
	Content   func()

	// Align overrides the centered label placement, for navigation rows and
	// other left-aligned actions.
	Align ButtonAlign

	// Corners overrides the uniform HeroUI radius. Nil keeps the default;
	// a zero array makes square corners.
	Corners *[4]float32
}

// buttonPressDuration is HeroUI's `transform 250ms var(--ease-smooth)`.
const buttonPressDuration = 250 * time.Millisecond

// pressedPose applies HeroUI's pressed scale around the element's center.
// The element is posed at rest in every build, so the first press animates
// from 1: MyGo shows a pose it was not posed with in the frame before at its
// target, which would otherwise snap to the pressed size. AnimatePose moves
// the committed box in paint-only frames and settles at once while the
// desktop asks for less motion.
func pressedPose(e ui.Element, pressed bool, scale float32) ui.Element {
	pose := ui.Pose{Scale: 1}
	if pressed {
		pose.Scale = scale
	}
	return e.TransformOrigin(0.5, 0.5).
		AnimatePose(pose, ui.Animation{Duration: buttonPressDuration, Ease: ui.EaseOut})
}

// buttonPressScale is HeroUI's pressed scale per size.
func buttonPressScale(size theme.ButtonSize) float32 {
	switch size {
	case theme.ButtonSmall:
		return 0.98
	case theme.ButtonLarge:
		return 0.96
	}
	return 0.97
}

// joinAdjust carries joined-group geometry into the internal builders
// without widening the public options. Its zero value is a standalone
// control.
type joinAdjust struct {
	corners    [4]float32
	hasCorners bool
	borders    [4]float32
	hasBorders bool
	// grouped turns the pressed scale off and lets seam paint the shared
	// edge HeroUI renders inside the following segment.
	grouped  bool
	seam     bool
	vertical bool
}

// joinCorners returns one segment's corner radii in MyGo's top-left,
// top-right, bottom-right, bottom-left order. Only the outer ends of a
// joined set keep the HeroUI radius.
func joinCorners(vertical bool, index, count int, radius float32, rtl bool) [4]float32 {
	if count <= 1 {
		return [4]float32{radius, radius, radius, radius}
	}
	if vertical {
		switch index {
		case 0:
			return [4]float32{radius, radius, 0, 0}
		case count - 1:
			return [4]float32{0, 0, radius, radius}
		}
		return [4]float32{}
	}
	start := [4]float32{radius, 0, 0, radius}
	end := [4]float32{0, radius, radius, 0}
	if rtl {
		start, end = end, start
	}
	switch index {
	case 0:
		return start
	case count - 1:
		return end
	}
	return [4]float32{}
}

// joinBorders returns the outline border widths, in top, right, bottom and
// left order, of one segment. HeroUI removes the border two joined outline
// buttons share so a seam stays hairline-thin.
func joinBorders(vertical bool, index, count int, rtl bool) [4]float32 {
	w := float32(theme.Hairline)
	if count <= 1 {
		return [4]float32{w, w, w, w}
	}
	if vertical {
		switch index {
		case 0:
			return [4]float32{w, w, 0, w}
		case count - 1:
			return [4]float32{0, w, w, w}
		}
		return [4]float32{0, w, 0, w}
	}
	first := [4]float32{w, w, w, 0}
	last := [4]float32{w, 0, w, w}
	if rtl {
		first, last = last, first
	}
	switch index {
	case 0:
		return first
	case count - 1:
		return last
	}
	return [4]float32{w, 0, w, 0}
}

// joinSeam paints HeroUI's group separator: a 1 DIP line across half the
// segment, 15% of its text color, at the inline start of every button after
// the first. HeroUI renders it inside the button, absolutely positioned, so
// it takes no width from the group; drawing it outside the box from
// DrawOver has the same effect and paints over the previous segment.
func joinSeam(e ui.Element, vertical, rtl bool, ink ui.Color) ui.Element {
	if ink.A == 0 {
		return e
	}
	w := float32(theme.Hairline)
	radius := theme.RadiusScale().SM
	return e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if vertical {
			p.Fill(ui.Rect{X: r.X + r.W*0.25, Y: r.Y - w, W: r.W * 0.5, H: w}, ink, radius)
			return
		}
		x := r.X - w
		if rtl {
			x = r.X + r.W
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + r.H*0.25, W: w, H: r.H * 0.5}, ink, radius)
	})
}

// Button constructs a HeroUI button on the native focusable button base.
// Hover and pressed paints follow live pointer state. A press scales the
// button to HeroUI's 0.97/0.98/0.96 and eases back on release; joined groups
// keep their size as HeroUI does. Pending buttons keep their full face but
// neither light up nor activate, matching pointer-events:none. Disabled
// buttons paint at HeroUI's 50% opacity.
func Button(c *ui.Context, o ButtonOptions) ui.Element {
	var adj joinAdjust
	if o.Corners != nil {
		adj.corners, adj.hasCorners = *o.Corners, true
	}
	return buttonView(c, o, adj)
}

// buttonView is the builder Button and ButtonGroup share, so a group can
// supply its segment geometry without an ambient scope.
func buttonView(c *ui.Context, o ButtonOptions, adj joinAdjust) ui.Element {
	t := theme.Colors(c)
	pending := o.Loading
	disabled := o.Disabled || pending
	height := theme.ButtonHeight(o.Size)
	radius := theme.RadiusScale().XXXL
	paddingX := float32(theme.Space16)
	glyph := float32(theme.IconMD)
	font := theme.SM()
	font.Weight = theme.WeightMedium
	switch o.Size {
	case theme.ButtonSmall:
		paddingX = theme.Space12
	case theme.ButtonLarge:
		font = theme.Base()
		font.Weight = theme.WeightMedium
	}
	if o.IconOnly {
		paddingX = 0
	}
	side := theme.Scale(c, height)
	iconSize := theme.Scale(c, glyph)

	e := ui.ButtonBase(c).MinHeight(side).Gap(theme.Space8).
		PaddingX(paddingX).PaddingY(0).Center().Shrink(0).
		FocusRing(true).Label(o.Label)
	if adj.hasCorners {
		e = e.Radius(adj.corners[0], adj.corners[1], adj.corners[2], adj.corners[3])
	} else {
		e = e.Radius(radius)
	}
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	switch o.Align {
	case ButtonAlignStart:
		e.Justify(ui.Start)
	case ButtonAlignEnd:
		e.Justify(ui.End)
	}
	if o.IconOnly {
		e.Width(side).MinWidth(side)
	}
	e = style.Disabled(e, o.Disabled)
	if o.Handle != nil {
		e.Bind(o.Handle)
	}

	bg, hover, pressed, fg := buttonColors(t, o.Variant)
	if pending {
		// HeroUI's pending state removes the pointer faces but keeps the
		// readable surface.
		hover, pressed = bg, bg
	}
	e.TextColor(fg)
	if o.Variant == ButtonOutline {
		if adj.hasBorders {
			e = e.BorderWidth(adj.borders[0], adj.borders[1], adj.borders[2], adj.borders[3]).
				BorderColor(t.Border)
		} else {
			e.Border(theme.Hairline, t.Border)
		}
	}
	if o.OnClick != nil && !disabled {
		e.OnClick(o.OnClick)
	}
	if o.Tooltip != "" {
		e.Tooltip(o.Tooltip)
	} else if o.IconOnly && o.Label != "" {
		e.Tooltip(o.Label)
	}

	content := func() {
		if o.Loading {
			loading.Arc(c, loading.Options{
				Size: iconSize, Color: fg,
			}).Label(locale.M().Loading)
		} else if o.Icon != nil {
			// HeroUI pulls an icon 2 DIP into the row's gap (-mx-0.5).
			ui.Icon(c, o.Icon).Size(iconSize, iconSize).Shrink(0).Margin(0, -2)
		}
		if o.Content != nil {
			o.Content()
		} else if !o.IconOnly {
			style.Text(ui.Text(c, o.Label), c, font).SingleLine()
		}
		if o.EndIcon != nil && !o.IconOnly {
			ui.Icon(c, o.EndIcon).Size(iconSize, iconSize).Shrink(0).Margin(0, -2)
		}
	}
	if adj.hasCorners {
		e = style.CornerFill(c, e, theme.Space8, adj.corners, bg, hover, pressed, content)
	} else {
		e = style.Background(e, radius, bg, hover, pressed)
		e.Children(content)
	}
	if !adj.grouped {
		e = pressedPose(e, e.Pressed() && !pending, buttonPressScale(o.Size))
	}
	if adj.seam {
		e = joinSeam(e, adj.vertical, c.Direction() == ui.RightToLeft, fg)
	}
	return e
}

func buttonColors(t theme.Palette, v ButtonVariant) (bg, hover, pressed, fg ui.Color) {
	switch v {
	case ButtonSecondary:
		return t.Default, t.DefaultHover, t.DefaultHover, t.AccentSoftForeground
	case ButtonTertiary:
		return t.Default, t.DefaultHover, t.DefaultHover, t.DefaultForeground
	case ButtonOutline:
		return ui.Transparent, t.Default.Alpha(0.6), t.Default.Alpha(0.6), t.DefaultForeground
	case ButtonGhost:
		return ui.Transparent, t.Default, t.Default, t.DefaultForeground
	case ButtonDanger:
		return t.Danger, t.DangerHover, t.DangerHover, t.DangerForeground
	case ButtonDangerSoft:
		return t.DangerSoft, t.DangerSoftHover, t.DangerSoftHover, t.DangerSoftForeground
	}
	return t.Accent, t.AccentHover, t.AccentHover, t.AccentForeground
}

// ButtonGroupOrientation selects the join axis.
type ButtonGroupOrientation uint8

const (
	GroupHorizontal ButtonGroupOrientation = iota
	GroupVertical
)

// ButtonGroupOptions controls a joined set of HeroUI buttons.
type ButtonGroupOptions struct {
	Orientation ButtonGroupOrientation
	FullWidth   bool
	Separators  bool
}

// ButtonGroup builds buttons with only the outer corners rounded, as HeroUI
// does. The caller supplies complete button options, avoiding an ambient
// context and keeping each action independently focusable. Separators paint
// HeroUI's half-height divider at each shared edge instead of taking layout
// width, outline borders collapse into one, and joined buttons do not scale
// on press.
func ButtonGroup(c *ui.Context, o ButtonGroupOptions, buttons ...ButtonOptions) ui.Element {
	vertical := o.Orientation == GroupVertical
	var e ui.Element
	if vertical {
		e = ui.Column(c).Gap(0).MinWidth(0)
	} else {
		e = ui.Row(c).Gap(0).MinWidth(0)
	}
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	radius := theme.RadiusScale().XXXL
	rtl := c.Direction() == ui.RightToLeft
	n := len(buttons)
	e.Children(func() {
		for i := range buttons {
			b := buttons[i]
			adj := joinAdjust{grouped: true, vertical: vertical}
			adj.corners, adj.hasCorners = joinCorners(vertical, i, n, radius, rtl), true
			if b.Variant == ButtonOutline {
				adj.borders, adj.hasBorders = joinBorders(vertical, i, n, rtl), true
			}
			adj.seam = o.Separators && i > 0
			if o.FullWidth && vertical {
				b.FullWidth = true
			}
			item := buttonView(c, b, adj)
			if o.FullWidth && !vertical {
				item.Grow(1).Basis(0).MinWidth(0)
			}
		}
	})
	return e
}
