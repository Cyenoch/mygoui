package overlay

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DrawerSide is the window edge a drawer slides from.
type DrawerSide uint8

const (
	// DrawerBottom is HeroUI's default drawer placement.
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
	// Width overrides the panel width of left/right drawers, which defaults
	// to HeroUI's w-80 (320 DIPs), or w-96 (384) from the 640 DIP
	// breakpoint. Height overrides the panel extent of top/bottom drawers,
	// which defaults to its content. Both stay within 85% of the window.
	Width, Height float32
	// Handle shows HeroUI's 36x4 drag-handle bar on top/bottom drawers: 8
	// DIPs under the panel's padding on a bottom drawer, and 8 DIPs off the
	// edge with the content above it on a top drawer. It is visual; pointer
	// swipe dismissal is not implemented.
	Handle         bool
	HideClose      bool
	PreventDismiss bool
	OnDismiss      func()
}

// Drawer shows a floating modal panel from one window edge, inset 16 DIPs
// from every edge with HeroUI's overlay radius. A left or right panel keeps
// the w-80/w-96 width and 85% cap; a top or bottom one spans the available
// width. The native dialog base traps and restores the focus and dismisses
// on backdrop press or Escape. The drawer takes no swipe: only the bottom
// sheet owns a grab strip, which keeps its scrolling body out of it.
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
	const margin = theme.Space16
	availW := max(float32(0), w-2*margin)
	availH := max(float32(0), h-2*margin)
	// The panel keeps the window margin and HeroUI's 85% cap, whichever is
	// the tighter.
	capW := min(w*0.85, availW)
	capH := min(h*0.85, availH)
	side := o.Side
	if *open {
		c.Key(dialogKey{open})
	}
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		back = paintBackdrop(c, back, o.Backdrop)
		back = dialogMotion(c, back, drawerMotionDuration)
		layoutEdge(back, side)
		t := theme.Colors(c)
		panel.Column().Clip().Background(t.Overlay).TextColor(t.OverlayForeground).
			MinHeight(0)
		panel = panelElevation(c, panel)
		switch side {
		case DrawerLeft, DrawerRight:
			// w-80, and w-96 from the sm breakpoint; max-w-[85vw].
			width := o.Width
			if width <= 0 {
				width = 320
				if w >= 640 {
					width = 384
				}
			}
			panel.Radius(overlayRadius()).Width(min(width, capW)).MaxWidth(capW).
				Height(availH)
		case DrawerTop:
			panel.FillWidth().MaxHeight(capH).Radius(overlayRadius())
			if o.Height > 0 {
				panel.Height(min(o.Height, capH))
			}
		default:
			panel.FillWidth().MaxHeight(capH).Radius(overlayRadius())
			if o.Height > 0 {
				panel.Height(min(o.Height, capH))
			}
		}
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
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
			hideClose:        o.HideClose,
			onClose:          close,
			titleStyle:       theme.Base(),
			descriptionStyle: theme.SM(),
			bodyInk:          t.Muted,
			padding:          theme.Space24,
			headerGap:        theme.Space12,
		})
		if o.Handle && side == DrawerTop {
			// HeroUI's top drawer keeps 8 DIPs under its handle instead of
			// the panel's padding, so the bar sits close to the edge.
			panel.Padding(theme.Space24, theme.Space24, theme.Space8, theme.Space24)
			drawerHandle(c, 0)
		}
	})
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		if o.PreventDismiss {
			*open = true
		} else if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	return p
}

// BottomSheetOptions configures native HeroUI's bottom sheet.
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

// BottomSheet shows native HeroUI's bottom panel: a floating overlay
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
	// The panel keeps the window margin and HeroUI's 85% cap, whichever is
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
			Background(t.Overlay).TextColor(t.OverlayForeground)
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
			titleStyle:       theme.LG(),
			descriptionStyle: theme.Base(),
			padding:          theme.Space20,
			headerGap:        theme.Space12,
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
			// dialog to close and fading for another 250 ms.
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

// layoutEdge overrides DialogBase's centered backdrop for a floating edge
// panel and keeps its 16 DIP margin from every window edge.
func layoutEdge(back ui.Element, side DrawerSide) {
	switch side {
	case DrawerLeft:
		back.Row().Justify(ui.Start).AlignItems(ui.Stretch).Padding(theme.Space16)
	case DrawerRight:
		back.Row().Justify(ui.End).AlignItems(ui.Stretch).Padding(theme.Space16)
	case DrawerTop:
		back.Column().Justify(ui.Start).AlignItems(ui.Stretch).Padding(theme.Space16)
	default:
		back.Column().Justify(ui.End).AlignItems(ui.Stretch).Padding(theme.Space16)
	}
}

// drawerHandle paints HeroUI's 36×4 separator bar with gapBelow DIPs under
// it, in a gutter the caller places.
func drawerHandle(c *ui.Context, gapBelow float32) {
	t := theme.Colors(c)
	gutter := ui.Box(c).FillWidth().Center().Padding(0, 0, gapBelow, 0).Shrink(0)
	gutter.Children(func() {
		ui.Box(c).Width(theme.Space32 + theme.Space4).Height(theme.Space4).
			Radius(theme.RadiusScale().XS).Background(t.Separator).Shrink(0)
	})
}
