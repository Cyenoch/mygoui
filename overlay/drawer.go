package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

var drawerEase = ui.CubicBezier(0.32, 0.72, 0, 1)

// DrawerSide is the window edge a drawer slides from.
type DrawerSide uint8

const (
	// DrawerBottom is the default drawer placement.
	DrawerBottom DrawerSide = iota
	DrawerTop
	DrawerLeft
	DrawerRight
)

// DrawerOptions configures an edge drawer.
type DrawerOptions struct {
	Side                 DrawerSide
	Backdrop             ModalBackdrop
	Title                string
	Description          string
	Header, Body, Footer func()
	BodyLabel            string
	// Width overrides the default 448 DIP side width. A side drawer keeps
	// 48 DIPs of the window visible. Height overrides the content-sized
	// extent of a top or bottom drawer.
	Width, Height float32
	// Inset floats the panel 16 DIPs from its edges at widths of 640 or more.
	Inset bool
	// Handle shows a visual bar at the center of the inward-facing edge.
	Handle bool
	// ShowClose adds a close button, which is hidden by default.
	ShowClose bool
	// HideClose overrides ShowClose.
	HideClose bool
	// PreventDismiss blocks swipe, backdrop and Escape dismissal. A close
	// button supplied with ShowClose remains an explicit dismissal action.
	PreventDismiss bool
	OnDismiss      func()
}

// Drawer shows a modal panel attached to an edge, rounding the corners facing
// the window. Inset opts into a floating panel on wider windows. Drag the
// middle of its inward-facing edge outward to dismiss; the body keeps its
// native scrolling. Focus and dismissal use the native dialog base.
func Drawer(c *ui.Context, open *bool, o DrawerOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	warmOverlay(c)
	wasOpen := *open
	close := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	w, h := c.Size()
	inset := float32(0)
	if o.Inset && w >= 640 {
		inset = theme.Space16
	}
	availW := max(float32(0), w-2*inset)
	availH := max(float32(0), h-2*inset)
	capW := max(float32(0), min(availW, w-theme.Space48))
	capH := max(float32(0), min(availH, h-theme.Space48))
	side := o.Side
	motionSide := mygoui.SideBottom
	distance := capH + inset
	switch side {
	case DrawerLeft:
		motionSide = mygoui.SideLeft
	case DrawerRight:
		motionSide = mygoui.SideRight
	case DrawerTop:
		motionSide = mygoui.SideTop
	}
	if side == DrawerLeft || side == DrawerRight {
		width := o.Width
		if width <= 0 {
			width = 448
		}
		distance = min(width, capW) + inset
	} else if o.Height > 0 {
		distance = min(o.Height, capH) + inset
	}
	scrim := modalScrim(c, open, o.Backdrop, drawerMotionDuration, drawerEase)
	if *open {
		c.Key(dialogKey{open})
	}
	swiped := false
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		if side == DrawerTop || side == DrawerBottom {
			if height := panel.LayoutBounds().H; height > 0 {
				distance = height + inset
			}
		}
		back = edgeExit(c, back, motionSide, distance, drawerMotionDuration, drawerEase)
		layoutDrawer(back, side, inset)
		t := theme.Colors(c)
		panel.Column().Clip().Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border).
			MinHeight(0)
		panel = panelElevation(c, panel)
		radius := theme.RadiusScale().XXL
		switch side {
		case DrawerLeft, DrawerRight:
			width := o.Width
			if width <= 0 {
				width = 448
			}
			panel.Width(min(width, capW)).MaxWidth(capW).Height(availH)
			if side == DrawerLeft {
				panel.Radius(0, radius, radius, 0)
			} else {
				panel.Radius(radius, 0, 0, radius)
			}
		case DrawerTop:
			panel.FillWidth().MaxHeight(capH).Radius(0, 0, radius, radius)
			if o.Height > 0 {
				panel.Height(min(o.Height, capH))
			}
		default:
			panel.FillWidth().MaxHeight(capH).Radius(radius, radius, 0, 0)
			if o.Height > 0 {
				panel.Height(min(o.Height, capH))
			}
		}
		if inset > 0 {
			panel.Radius(radius)
		}
		edgeEnter(c, panel, motionSide, distance, drawerMotionDuration, drawerEase)
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
		}
		if distance > 0 {
			offset, dismiss := drawerDrag(c, panel, distance, o)
			swiped = swiped || dismiss
			panel.Translate(offset.X, offset.Y)
			moved := offset.Y
			if side == DrawerLeft || side == DrawerRight {
				moved = offset.X
			}
			scrim.Opacity(1 - min(max(moved, -moved)/distance, 1))
		}
		if o.Handle && side == DrawerBottom {
			drawerHandle(c, theme.Space8)
		}
		fillPanel(c, panel, chrome{
			title:            o.Title,
			description:      o.Description,
			header:           o.Header,
			body:             o.Body,
			footer:           o.Footer,
			bodyLabel:        o.BodyLabel,
			hideClose:        !o.ShowClose || o.HideClose,
			onClose:          close,
			fullBleedFooter:  true,
			titleStyle:       theme.H4(),
			descriptionStyle: theme.SM(),
			padding:          theme.Space24,
			headerGap:        theme.Space8,
		})
		if o.Handle && side == DrawerTop {
			// Keep the handle near the edge, below the content.
			panel.Padding(theme.Space24, theme.Space24, theme.Space8, theme.Space24)
			handle := drawerHandle(c, 0)
			if o.Height > 0 {
				panel.Padding(theme.Space24, theme.Space24, theme.Space12, theme.Space24)
				handle.Absolute().Left(0).Right(0).Bottom(theme.Space8)
			}
		}
	})
	if swiped {
		close()
	}
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		if o.PreventDismiss {
			*open = true
		} else if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	return p
}

type (
	drawerDragKey struct{}
	drawerGrabKey struct{}
)

// drawerDrag keeps the swipe separate from the scrollable body and controls.
func drawerDrag(c *ui.Context, panel ui.Element, travel float32, o DrawerOptions) (ui.Point, bool) {
	off := ui.Local(panel, drawerDragKey{}, func() ui.Point { return ui.Point{} })
	grab := ui.Box(c.Key(drawerGrabKey{})).Absolute().Role(ui.RoleNone).KeepFocus()
	const grabLength = theme.Space48 + theme.Space24
	const grabWidth = theme.Space32
	var centerX, centerY float32
	var target ui.Point
	var extent float32
	axis := ui.Vertical
	bounds := ui.Rect{}
	switch o.Side {
	case DrawerLeft:
		grab.Right(0).TopPercent(50).Size(grabWidth, grabLength)
		centerY = -grabLength / 2
		target.X, bounds.X, bounds.W = -travel, -travel, travel
		axis, extent = ui.Horizontal, panel.LayoutBounds().W
	case DrawerRight:
		grab.Left(0).TopPercent(50).Size(grabWidth, grabLength)
		centerY = -grabLength / 2
		target.X, bounds.W = travel, travel
		axis, extent = ui.Horizontal, panel.LayoutBounds().W
	case DrawerTop:
		grab.Bottom(0).LeftPercent(50).Size(grabLength, grabWidth)
		centerX = -grabLength / 2
		target.Y, bounds.Y, bounds.H = -travel, -travel, travel
		extent = panel.LayoutBounds().H
	default:
		grab.Top(0).LeftPercent(50).Size(grabLength, grabWidth)
		centerX = -grabLength / 2
		target.Y, bounds.H = travel, travel
		extent = panel.LayoutBounds().H
	}
	dismiss := false
	grab.Draggable(off, ui.DraggableOptions{
		Axis: axis, Bounds: &bounds, Rigid: true,
		Animation: ui.Animation{Duration: drawerMotionDuration, Ease: drawerEase},
		Snap: func(p ui.Point) ui.Point {
			moved := p.Y
			if axis == ui.Horizontal {
				moved = p.X
			}
			if !o.PreventDismiss && max(moved, -moved) > sheetDismissAt(extent) {
				dismiss = true
				return target
			}
			return ui.Point{}
		},
	})
	shown := grab.DragOffset()
	// The panel follows the strip, so cancel its own transform.
	grab.Translate(centerX-shown.X, centerY-shown.Y)
	if o.Handle && axis == ui.Horizontal {
		grab.Center().Children(func() {
			ui.Box(c).Width(theme.Space4).Height(theme.Space48).
				Radius(theme.RadiusScale().XS).Background(theme.Colors(c).Input)
		})
	}
	return shown, dismiss
}

// BottomSheetOptions configures a draggable native bottom sheet.
type BottomSheetOptions struct {
	Backdrop             ModalBackdrop
	Title                string
	Description          string
	Header, Body, Footer func()
	BodyLabel            string
	// Height fixes the panel height; zero sizes it to its content. Both
	// are capped by 85% of the window or its 16 DIP gutters, whichever is
	// the tighter.
	Height float32
	// HideHandle removes the drag handle and the swipe it owns; the
	// backdrop, Escape and the close button still dismiss the sheet.
	HideHandle     bool
	HideClose      bool
	PreventDismiss bool
	OnDismiss      func()
}

// BottomSheet shows a floating bottom panel: an overlay
// surface, 16 DIPs from the left, right and bottom window edges, with all
// four corners rounded, a centered separator handle and 20 DIP padding. It
// slides up from below the window and its handle drags it down to dismiss.
// It is modal, so the native dialog base keeps and restores the focus.
func BottomSheet(c *ui.Context, open *bool, o BottomSheetOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	warmOverlay(c)
	wasOpen := *open
	close := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	w, h := c.Size()
	const margin = theme.Space16
	availW := max(float32(0), w-2*margin)
	availH := max(float32(0), h-2*margin)
	// The panel keeps the window margin and 85% cap, whichever is
	// the tighter.
	capH := min(h*0.85, availH)
	// The sheet rests inside the gutters, so a full window height of travel
	// leaves it below the window at the start and the end of a swipe.
	travel := max(h, 0)
	reduce := theme.ReduceMotion(c)
	if *open {
		c.Key(dialogKey{open})
	}
	swipedOut := false
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		back = paintBackdrop(c, back, o.Backdrop)
		back = dialogMotion(c, back, drawerMotionDuration)
		layoutEdge(back, DrawerBottom)
		t := theme.Colors(c)
		panel.Column().FillWidth().MaxWidth(availW).MaxHeight(capH).MinHeight(0).Clip().
			Radius(overlayRadius()).
			Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border)
		if o.Height > 0 {
			panel.Height(min(o.Height, capH))
		}
		panel = panelElevation(c, panel)
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
		}
		if !o.HideHandle {
			drawerHandle(c, theme.Space8)
		}
		dismiss := ui.Local(panel, sheetDismissKey{}, func() bool { return false })
		fillPanel(c, panel, chrome{
			title:       o.Title,
			description: o.Description,
			header:      o.Header,
			body:        o.Body,
			footer:      o.Footer,
			bodyLabel:   o.BodyLabel,
			hideClose:   o.HideClose,
			onClose: func() {
				if reduce || o.HideHandle || travel <= 0 {
					close()
					return
				}
				*dismiss = true
				c.AnimationFrame()
			},
			titleStyle:       theme.H4(),
			descriptionStyle: theme.SM(),
			padding:          theme.Space20,
			headerGap:        theme.Space8,
		})
		// The sheet arrives from below the window: its first frame shows it
		// there and the next sends it to its place. Both are translations,
		// so the layout never moves.
		entered := ui.Local(panel, sheetEnteredKey{}, func() bool { return reduce })
		enterY := float32(0)
		if !*entered {
			enterY = travel
			*entered = true
			c.AnimationFrame()
		}
		drop := float32(0)
		if !o.HideHandle && travel > 0 {
			var out bool
			drop, out = sheetDrag(c, panel, travel, o.PreventDismiss, *dismiss)
			swipedOut = swipedOut || out
		}
		if travel > 0 && drop > 0 {
			// The scrim leaves with the sheet instead of waiting for the
			// dialog to close and fading for another transition.
			back.Opacity(1 - min(drop/travel, 1))
		}
		panel.Translate(0, panel.AnimateTo(sheetEnterKey{}, enterY, sheetEnterAnimation)+drop)
	})
	if swipedOut {
		close()
	}
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		if o.PreventDismiss {
			*open = true
		} else if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	return p
}

// sheetEnterKey names the panel's arrival translation, sheetEnteredKey the
// frame that releases it, and sheetDragKey the grab strip's offset: all on
// the panel's own state, which goes with the dialog when it closes.
type (
	sheetEnterKey   struct{}
	sheetEnteredKey struct{}
	sheetDragKey    struct{}
	sheetDismissKey struct{}
)

// sheetEnterAnimation slides the sheet up from below the window, over the
// drawer backdrop's duration.
var sheetEnterAnimation = ui.Animation{Duration: drawerMotionDuration, Ease: ui.EaseOut}

const (
	// sheetGrabHeight is the strip that drags the sheet: the panel's 20 DIP
	// padding and the handle's 4 DIP bar and 8 DIP gutter.
	sheetGrabHeight = theme.Space20 + theme.Space4 + theme.Space8
	// sheetDismissDrag is how far a sheet is dragged down before letting go
	// dismisses it: 96 DIPs, or 30% of a sheet too short to need that.
	sheetDismissDrag = 96
	sheetDismissPart = 0.3
)

// sheetDrag builds the invisible strip that owns the sheet's downward
// swipe, over the handle and the padding above it, and returns where the
// sheet shows from its place and whether the sheet slid all the way out.
// Only that strip is draggable, so the body keeps its scrolling, and the
// strip keeps the focus on what it was.
func sheetDrag(c *ui.Context, panel ui.Element, travel float32, prevent, dismissing bool) (float32, bool) {
	off := ui.Local(panel, sheetDragKey{}, func() ui.Point { return ui.Point{} })
	if dismissing && off.Y < travel {
		// The close button asks for the same slide the handle swipe makes.
		off.Y = travel
	}
	grab := ui.Box(c).Absolute().Top(0).Start(0).Height(sheetGrabHeight).
		End(theme.Space16 + closeButtonSize + theme.Space8).
		Role(ui.RoleNone).KeepFocus()
	bounds := ui.Rect{H: travel}
	grab.Draggable(off, ui.DraggableOptions{
		Axis:      ui.Vertical,
		Bounds:    &bounds,
		Rigid:     true,
		Animation: sheetEnterAnimation,
		Snap: func(p ui.Point) ui.Point {
			if !prevent && p.Y > sheetDismissAt(panel.Bounds().H) {
				return ui.Point{Y: travel}
			}
			return ui.Point{}
		},
	})
	shown := grab.DragOffset()
	// The panel takes the same travel, which would move this strip twice:
	// cancel its own drag so it stays on the handle bar while the sheet
	// moves, and a second grab during the spring lands on it.
	grab.Translate(0, -shown.Y)
	// Draggable comes to rest exactly on the point its Snap chose. Close as
	// soon as the sheet has left the window: waiting for the spring's
	// sub-pixel precision would keep a modal offscreen for about a second.
	return shown.Y, shown.Y >= travel-1
}

// sheetDismissAt returns how far a sheet of height h is dragged down before
// letting go dismisses it: 96 DIPs or 30% of the sheet, whichever is less.
func sheetDismissAt(h float32) float32 {
	if h <= 0 {
		return sheetDismissDrag
	}
	return min(float32(sheetDismissDrag), h*sheetDismissPart)
}

// layoutDrawer places a drawer against the requested edge.
func layoutDrawer(back ui.Element, side DrawerSide, inset float32) {
	switch side {
	case DrawerLeft:
		back.Row().Justify(ui.Start).AlignItems(ui.Stretch).Padding(inset)
	case DrawerRight:
		back.Row().Justify(ui.End).AlignItems(ui.Stretch).Padding(inset)
	case DrawerTop:
		back.Column().Justify(ui.Start).AlignItems(ui.Stretch).Padding(inset)
	default:
		back.Column().Justify(ui.End).AlignItems(ui.Stretch).Padding(inset)
	}
}

// layoutEdge keeps the bottom sheet's existing floating geometry.
func layoutEdge(back ui.Element, side DrawerSide) {
	layoutDrawer(back, side, theme.Space16)
}

// drawerHandle paints a 48 by 4 DIP bar in the caller's gutter.
func drawerHandle(c *ui.Context, gapBelow float32) ui.Element {
	t := theme.Colors(c)
	gutter := ui.Box(c).FillWidth().Center().Padding(0, 0, gapBelow, 0).Shrink(0).PassThrough()
	gutter.Children(func() {
		ui.Box(c).Width(theme.Space48).Height(theme.Space4).
			Radius(theme.RadiusScale().XS).Background(t.Input).Shrink(0).PassThrough()
	})
	return gutter
}
