package mygoui

import (
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// ButtonVariant selects a button face, including native theme extensions.
type ButtonVariant uint8

const (
	ButtonPrimary ButtonVariant = iota
	ButtonSecondary
	ButtonTertiary
	ButtonOutline
	ButtonGhost
	ButtonDanger
	ButtonDangerSoft
	ButtonDestructiveOutline
	ButtonLink
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

	// Corners overrides the uniform coss ui radius. Nil keeps the default;
	// a zero array makes square corners.
	Corners *[4]float32
}

// buttonPressDuration is the native press-feedback transition.
const buttonPressDuration = 250 * time.Millisecond

// pressedPose applies native press feedback around the element's center.
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

// joinAdjust carries joined-group geometry into the internal builders
// without widening the public options. Its zero value is a standalone
// control.
type joinAdjust struct {
	corners    [4]float32
	hasCorners bool
	borders    [4]float32
	hasBorders bool
	grouped    bool
	// seam paints a separator at the shared edge of the next segment.
	seam     bool
	vertical bool
}

// joinCorners returns one segment's corner radii in MyGo's top-left,
// top-right, bottom-right, bottom-left order. Only the outer ends of a
// joined set keep the coss ui radius.
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

// joinBorders removes the inside edges of a joined segment, in top, right,
// bottom and left order. An explicit separator can occupy the shared edge.
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
	first := [4]float32{w, 0, w, w}
	last := [4]float32{w, w, w, 0}
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

// joinSeam paints a 1 DIP group separator without taking layout width.
func joinSeam(e ui.Element, vertical, rtl bool, ink ui.Color) ui.Element {
	if ink.A == 0 {
		return e
	}
	w := float32(theme.Hairline)
	ink = ink.Alpha(0.15)
	radius := theme.RadiusScale().SM
	return e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if vertical {
			p.Fill(ui.Rect{X: r.X + r.W*0.25, Y: r.Y, W: r.W * 0.5, H: w}, ink, radius)
			return
		}
		x := r.X
		if rtl {
			x = r.X + r.W - w
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + r.H*0.25, W: w, H: r.H * 0.5}, ink, radius)
	})
}

// Button builds a coss action on MyGo's native focusable button base.
// The pending spinner overlays invisible content, preserving the button width.
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
	radius := theme.RadiusScale().LG
	paddingX := float32(11)
	glyph := float32(theme.IconMD)
	font := theme.SM()
	font.Weight = theme.WeightMedium
	switch o.Size {
	case theme.ButtonSmall:
		paddingX = 9
	case theme.ButtonLarge:
		paddingX = 13
	case theme.ButtonExtraSmall:
		paddingX = 7
		font = theme.XS()
		font.Weight = theme.WeightMedium
		radius = theme.RadiusScale().MD
	case theme.ButtonExtraLarge:
		paddingX = 15
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
	e = style.Disabled(e, disabled)
	if o.Handle != nil {
		e.Bind(o.Handle)
	}

	bg, hover, pressed, fg := buttonColors(t, o.Variant, c.Theme().Dark)
	if pending {
		// coss ui's pending state removes the pointer faces but keeps the
		// readable surface.
		hover, pressed = bg, bg
	}
	e.TextColor(fg)
	if o.Variant == ButtonOutline || o.Variant == ButtonDestructiveOutline {
		if adj.hasBorders {
			e = e.BorderWidth(adj.borders[0], adj.borders[1], adj.borders[2], adj.borders[3]).
				BorderColor(t.Input)
		} else {
			e.Border(theme.Hairline, t.Input)
		}
	}
	if !disabled && !e.Pressed() && !adj.grouped {
		switch o.Variant {
		case ButtonPrimary, ButtonDanger:
			e = e.Shadow(0, 1, 2, 0, bg.Alpha(0.24)).InsetShadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.16))
		case ButtonOutline, ButtonDestructiveOutline:
			e = style.Shadow(e, c, theme.FieldShadows(c.Theme().Dark))
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
		row := ui.Row(c).Center().Gap(theme.Space8).MinWidth(0)
		if o.Loading {
			row.Opacity(0)
		}
		row.Children(func() {
			if o.Icon != nil {
				ui.Icon(c, o.Icon).Size(iconSize, iconSize).Shrink(0).Margin(0, -2).Opacity(0.8)
			}
			if o.Content != nil {
				o.Content()
			} else if !o.IconOnly {
				label := style.Text(ui.Text(c, o.Label), c, font).SingleLine()
				if o.Variant == ButtonLink && (e.Hovered() || e.Pressed()) {
					label.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Line(r.X, r.Y+r.H-1, r.X+r.W, r.Y+r.H-1, 1, fg) })
				}
			}
			if o.EndIcon != nil && !o.IconOnly {
				ui.Icon(c, o.EndIcon).Size(iconSize, iconSize).Shrink(0).Margin(0, -2).Opacity(0.8)
			}
		})
	}
	var overlay func()
	if o.Loading {
		spinnerInk := fg
		if o.Variant == ButtonDestructiveOutline {
			spinnerInk = t.Foreground
		}
		overlay = func() {
			ui.Box(c).Absolute().Left(0).Right(0).Top(0).Bottom(0).Center().Children(func() {
				loading.Arc(c, loading.Options{Size: iconSize, Color: spinnerInk, Duration: time.Second}).
					Label(locale.M().Loading)
			})
		}
	}
	if adj.hasCorners {
		e = style.CornerFill(c, e, theme.Space8, adj.corners, bg, hover, pressed, content, overlay)
	} else {
		e = style.Background(e, radius, bg, hover, pressed)
		e.Children(func() {
			content()
			if overlay != nil {
				overlay()
			}
		})
	}
	if adj.seam {
		e = joinSeam(e, adj.vertical, c.Direction() == ui.RightToLeft, fg)
	}
	return e
}

func buttonColors(t theme.Palette, v ButtonVariant, dark bool) (bg, hover, pressed, fg ui.Color) {
	switch v {
	case ButtonSecondary, ButtonTertiary:
		return t.Secondary, t.Secondary.Alpha(0.9), t.Secondary.Alpha(0.8), t.SecondaryForeground
	case ButtonOutline:
		if dark {
			return t.Input.Alpha(0.32), t.Input.Alpha(0.64), t.Input.Alpha(0.64), t.Foreground
		}
		return t.Popover, t.Accent.Alpha(0.5), t.Accent.Alpha(0.5), t.Foreground
	case ButtonGhost:
		return ui.Transparent, t.Accent, t.Accent, t.Foreground
	case ButtonLink:
		return ui.Transparent, ui.Transparent, ui.Transparent, t.Foreground
	case ButtonDanger:
		return t.Destructive, t.Destructive.Alpha(0.9), t.Destructive.Alpha(0.9), ui.RGB(255, 255, 255)
	case ButtonDestructiveOutline:
		bg := t.Popover
		if dark {
			bg = t.Input.Alpha(0.32)
		}
		return bg, t.Destructive.Alpha(0.04), t.Destructive.Alpha(0.04), t.DestructiveForeground
	case ButtonDangerSoft:
		return t.DangerSoft, t.DangerSoftHover, t.DangerSoftHover, t.DestructiveForeground
	}
	return t.Primary, t.Primary.Alpha(0.9), t.Primary.Alpha(0.9), t.PrimaryForeground
}

// ButtonGroupOrientation selects the join axis.
type ButtonGroupOrientation uint8

const (
	GroupHorizontal ButtonGroupOrientation = iota
	GroupVertical
)

// ButtonGroupOptions controls a joined set of coss ui buttons.
type ButtonGroupOptions struct {
	Orientation ButtonGroupOrientation
	FullWidth   bool
	Separators  bool
}

// ButtonGroup builds buttons with only the outer corners rounded, as coss ui
// does. The caller supplies complete button options, avoiding an ambient
// context and keeping each action independently focusable. Separators paint
// a half-height divider at each shared edge instead of taking layout width.
// Adjacent outline edges are removed, and buttons do not scale on press.
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
	radius := theme.RadiusScale().LG
	rtl := c.Direction() == ui.RightToLeft
	n := len(buttons)
	e.Children(func() {
		for i := range buttons {
			b := buttons[i]
			adj := joinAdjust{grouped: true, vertical: vertical}
			adj.corners, adj.hasCorners = joinCorners(vertical, i, n, radius, rtl), true
			if b.Variant == ButtonOutline || b.Variant == ButtonDestructiveOutline {
				adj.borders, adj.hasBorders = joinBorders(vertical, i, n, rtl), true
			}
			adj.seam = o.Separators && i > 0
			if o.FullWidth && vertical {
				b.FullWidth = true
			}
			item := buttonView(c, b, adj)
			if vertical {
				item.FillWidth()
			}
			if o.FullWidth && !vertical {
				item.Grow(1).Basis(0).MinWidth(0)
			}
		}
	})
	return e
}
