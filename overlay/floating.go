package overlay

import (
	"math"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// FloatingPanelOptions configures a floating panel. Position is the panel's
// top-left corner in window DIPs and caller-owned; nil keeps it in the
// panel, which then opens 16 DIPs from the window's top right corner. Width
// defaults to 320 DIPs.
type FloatingPanelOptions struct {
	Title     string
	Body      func()
	Position  *ui.Point
	Width     float32
	HideClose bool
	OnDismiss func()
	// ReturnFocus takes the focus when the panel closes with the focus
	// inside it; without it the focus returns to the element that held it
	// as the panel opened, if any.
	ReturnFocus *ui.Handle
}

type (
	floatingPanelKey struct{ open *bool }
	floatingPosKey   struct{}
)

// FloatingPanel floats a non-modal panel over the window while *open is
// true. Its title bar drags it, and Alt with an arrow moves it while the
// focus is inside, within the window. It takes no focus, has no backdrop
// and leaves the content around it usable; its close button, or Escape while
// the focus is inside it, closes it. A modal dialog keeps the focus to
// itself, so a panel shown over one cannot be focused.
func FloatingPanel(c *ui.Context, open *bool, o FloatingPanelOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	focus := overlayFocusReturn(c, open)
	if !*open {
		if o.ReturnFocus != nil {
			focus.closed(func() { o.ReturnFocus.Focus(c) })
		}
		return ui.Element{}
	}
	close := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	w, h := c.Size()
	width := o.Width
	if width <= 0 {
		width = 320
	}
	width = min(width, w)
	var panel ui.Element
	ui.Overlay(c, func() {
		panel = ui.Column(c.Key(floatingPanelKey{open})).Absolute().Width(width).MaxHeight(h).MinHeight(0).
			Role(ui.RoleDialog).Label(o.Title)
		panel = anchored.Surface(c, panel, overlayRadius()).Clip()
		panel = surfaceMotion(c, panel, theme.Motion().Popup, theme.Motion().Ease)
		pos := o.Position
		if pos == nil {
			pos = ui.Local(panel, floatingPosKey{}, func() ui.Point {
				return ui.Point{X: w - width - theme.Space(4), Y: theme.Space(4)}
			})
		}
		if panel.Shortcut(0, ui.KeyEscape) {
			close()
		}
		for _, k := range []struct {
			key    ui.Key
			dx, dy float32
		}{{ui.KeyLeft, -1, 0}, {ui.KeyRight, 1, 0}, {ui.KeyUp, 0, -1}, {ui.KeyDown, 0, 1}} {
			if panel.Shortcut(ui.Alt, k.key) {
				pos.X += k.dx * floatingStep
				pos.Y += k.dy * floatingStep
			}
		}
		panel.Children(func() {
			bar := floatingHeader(c, o.Title, !o.HideClose, close).Cursor(ui.CursorMove)
			if dx, dy, held := bar.Dragged(); held {
				pos.X += dx
				pos.Y += dy
			}
			if o.Body != nil {
				ui.Scroll(c).FillWidth().Grow(1).MinHeight(0).Padding(theme.Space(4)).Children(o.Body)
			}
		})
		// The height is the one the frame before laid out; the first frame
		// has none, so it asks for a second to settle inside the window.
		size := panel.LayoutBounds()
		if size.H <= 0 {
			c.AnimationFrame()
		}
		pos.X = max(0, min(pos.X, w-width))
		pos.Y = max(0, min(pos.Y, h-size.H))
		panel.Left(pos.X).Top(pos.Y)
		focus.note(panel)
	})
	return panel
}

// floatingHeader builds the title row floating surfaces share: the title and
// an optional close button over a hairline.
func floatingHeader(c *ui.Context, title string, closable bool, close func()) ui.Element {
	t := theme.Colors(c)
	bar := ui.Row(c).FillWidth().Shrink(0).AlignItems(ui.Center).Gap(theme.Space(2)).
		MinHeight(theme.Scale(c, theme.Space(10))).
		Padding(theme.Space(2), theme.Space(2), theme.Space(2), theme.Space(4)).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border)
	bar.Children(func() {
		style.Text(ui.Text(c, title), c, theme.SM()).FontWeight(theme.WeightSemibold).
			SingleLine().Grow(1).Basis(0).MinWidth(0)
		if closable {
			mygoui.CloseButton(c, mygoui.CloseButtonOptions{Label: locale.M().Close, OnClick: close})
		}
	})
	return bar
}

// FloatingSheetOptions configures a floating sheet. Height is caller-owned
// and follows the handle between MinHeight (160 DIPs by default) and
// MaxHeight (80% of the container by default); nil keeps it in the sheet,
// starting at MinHeight. HideHandle fixes the height.
type FloatingSheetOptions struct {
	Title      string
	Body       func()
	Height     *float32
	MinHeight  float32
	MaxHeight  float32
	HideHandle bool
	HideClose  bool
	OnDismiss  func()
	// ReturnFocus takes the focus when the sheet closes with the focus
	// inside it: the sheet is part of its container, so nothing else gives
	// that focus back.
	ReturnFocus *ui.Handle
}

type (
	floatingSheetKey  struct{}
	floatingHeightKey struct{}
	floatingFocusKey  struct{}
)

// floatingStep is how far an arrow key moves a panel or resizes a sheet.
const floatingStep = 24

// FloatingSheet presents a non-modal sheet rising from the bottom edge of
// the container it is built in, which it covers without taking its input:
// build it last in a container that holds the content it floats over. It has
// no backdrop and takes no focus. Its handle resizes it by dragging, or
// while focused with the Up and Down arrows, Home to MinHeight and End to
// MaxHeight; its close button, or Escape while the focus is inside it,
// closes it.
func FloatingSheet(c *ui.Context, open *bool, o FloatingSheetOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	frame := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Justify(ui.End).PassThrough()
	target := float32(0)
	if *open {
		target = 1
	}
	shown := frame.AnimateTo(floatingSheetKey{}, target, sheetEnterAnimation())
	focus := ui.Local(frame, floatingFocusKey{}, func() focusReturn { return focusReturn{} })
	if !*open && o.ReturnFocus != nil {
		focus.closed(func() { o.ReturnFocus.Focus(c) })
	}
	if !*open && shown <= 0 {
		return frame
	}
	close := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	lo := o.MinHeight
	if lo <= 0 {
		lo = 160
	}
	hi := float32(math.Inf(1))
	if o.MaxHeight > 0 {
		hi = o.MaxHeight
	}
	// The container is measured by the frame before: until then only the
	// caller's limits hold.
	if container := frame.LayoutBounds().H; container > 0 {
		if o.MaxHeight <= 0 {
			hi = container * 0.8
		}
		hi = min(hi, container)
	}
	height := o.Height
	if height == nil {
		height = ui.Local(frame, floatingHeightKey{}, func() float32 { return lo })
	}
	clamp := func() { *height = max(min(*height, hi), min(lo, hi)) }
	clamp()

	radius := overlayRadius()
	t := theme.Colors(c)
	frame.Children(func() {
		// The clipped box rises; the sheet inside keeps its height, so the
		// body does not reflow as it arrives.
		rise := ui.Column(c).FillWidth().Shrink(0).Clip().Radius(radius, radius, 0, 0).
			Background(t.Popover).TextColor(t.PopoverForeground).
			BorderWidth(theme.Hairline, theme.Hairline, 0, theme.Hairline).BorderColor(t.Border).
			Role(ui.RoleDialog).Label(o.Title).IgnoreInput(!*open)
		rise = panelElevation(c, rise)
		if rise.Shortcut(0, ui.KeyEscape) {
			close()
		}
		rise.Children(func() {
			sheet := ui.Column(c).FillWidth().Shrink(0)
			sheet.Children(func() {
				if !o.HideHandle {
					floatingSheetHandle(c, height, min(lo, hi), hi)
					clamp()
				}
				floatingHeader(c, o.Title, !o.HideClose, close)
				if o.Body != nil {
					ui.Scroll(c).FillWidth().Grow(1).MinHeight(0).Padding(theme.Space(4)).Children(o.Body)
				}
			})
			sheet.Height(*height)
		})
		rise.Height(*height * shown)
		if *open {
			focus.note(rise)
		}
	})
	return frame
}

// floatingSheetHandle builds the strip that resizes the sheet between lo
// and hi and applies what the pointer or the keys ask to *height; the
// caller holds it within those limits.
func floatingSheetHandle(c *ui.Context, height *float32, lo, hi float32) {
	grab := ui.Box(c).FillWidth().Shrink(0).Padding(theme.Space(2), 0, 0, 0).
		Focusable().FocusRing(true).Role(ui.RoleSplitter).Label(locale.M().Resize).
		Cursor(ui.CursorResizeNS)
	if _, dy, held := grab.Dragged(); held {
		*height -= dy
	}
	switch {
	case grab.Shortcut(0, ui.KeyUp):
		*height += floatingStep
	case grab.Shortcut(0, ui.KeyDown):
		*height -= floatingStep
	case grab.Shortcut(0, ui.KeyHome):
		*height = lo
	case grab.Shortcut(0, ui.KeyEnd) && !math.IsInf(float64(hi), 1):
		*height = hi
	}
	// Until the container is measured the sheet has no upper limit to tell.
	if !math.IsInf(float64(hi), 1) {
		grab.Range(float64(lo), float64(hi), float64(max(min(*height, hi), lo))).Step(floatingStep)
	}
	grab.Children(func() { drawerHandle(c, 0) })
}
