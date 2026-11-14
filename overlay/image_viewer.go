package overlay

import (
	"math"
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/internal/touchhint"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ImageViewerItem is one decoded image and its accessible label.
type ImageViewerItem struct {
	Source ui.ImageSource
	Label  string
}

// ImageViewerOptions configures an image preview. Index is the caller-owned
// image shown; nil keeps it in the viewer, starting at the first. The
// library loads nothing: every Source is already decoded.
type ImageViewerOptions struct {
	Images    []ImageViewerItem
	Index     *int
	OnDismiss func()
}

const (
	viewerMaxZoom  = 8
	viewerZoomStep = 1.5
	// viewerWheelZoom is the zoom's exponent per DIP the wheel scrolls.
	viewerWheelZoom = 0.002
)

type viewerStateKey struct{}

// viewerState is the transient view of the image shown: its zoom over the
// fitted size and its pan from the center.
type viewerState struct {
	index, shown int
	zoom, x, y   float32
	held, panned bool
	// smooth has the image glide to a zoom a key or a button chose; the
	// wheel, a pinch and a drag move it as the hand does.
	smooth bool
}

func (v *viewerState) reset() { v.zoom, v.x, v.y = 1, 0, 0 }

// step zooms to a level a key or a button asked for, about the center.
func (v *viewerState) step(zoom float32) {
	v.zoomTo(zoom, 0, 0)
	v.smooth = true
}

// restore returns to the fitted view as a key or a button asked.
func (v *viewerState) restore() {
	v.reset()
	v.smooth = true
}

// zoomTo sets the zoom, keeping the image point under (ax, ay), measured
// from the viewport's center, in place.
func (v *viewerState) zoomTo(zoom, ax, ay float32) {
	zoom = min(max(zoom, 1), viewerMaxZoom)
	k := zoom / v.zoom
	v.x, v.y = ax-k*(ax-v.x), ay-k*(ay-v.y)
	v.zoom = zoom
}

// clamp keeps a zoomed image over the viewport: one that fits stays centered.
func (v *viewerState) clamp(fitW, fitH, viewW, viewH float32) {
	axis := func(at, size, view float32) float32 {
		slack := (size*v.zoom - view) / 2
		if slack <= 0 {
			return 0
		}
		return min(max(at, -slack), slack)
	}
	v.x, v.y = axis(v.x, fitW, viewW), axis(v.y, fitH, viewH)
}

// viewerFit is the size an image takes inside the window's margins, never
// larger than its own.
func viewerFit(src ui.ImageSource, w, h float32) (float32, float32) {
	iw, ih := ui.ImageSize(src)
	if iw <= 0 || ih <= 0 {
		return 0, 0
	}
	k := max(0, min((w-2*theme.Space(6))/iw, (h-2*theme.Space(16))/ih, 1))
	return iw * k, ih * k
}

// ImageViewer previews caller-provided images over a dimmed window while
// *open is true. The image fits the window; the toolbar, the + - and 0 keys
// and the wheel zoom it, and dragging pans it while zoomed. The arrows and
// the side buttons step through several images, mirrored right to left. Escape, the
// close button or a press outside the image dismisses it. It is modal: the
// native dialog base keeps and restores the focus.
func ImageViewer(c *ui.Context, open *bool, o ImageViewerOptions) ui.Element {
	if open == nil || len(o.Images) == 0 {
		return ui.Element{}
	}
	anchored.Warm(c)
	wasOpen := *open
	close := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	if *open {
		c.Key(dialogKey{open})
	}
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		back = paintBackdrop(c, back, ModalBackdropBlur)
		surfaceMotion(c, back, theme.Motion().Dialog, theme.Motion().EaseTravel)
		w, h := c.Size()
		panel.Width(w).Height(h)
		v := ui.Local(panel, viewerStateKey{}, func() viewerState { return viewerState{zoom: 1} })
		index := o.Index
		if index == nil {
			index = &v.index
		}
		n := len(o.Images)
		*index = max(0, min(*index, n-1))
		step := func(by int) {
			if to := *index + by; to >= 0 && to < n {
				*index = to
			}
		}
		// The arrows step the way the images run: right to left, Left
		// goes forward.
		earlier, later := ui.KeyLeft, ui.KeyRight
		if c.Direction() == ui.RightToLeft {
			earlier, later = later, earlier
		}
		switch {
		case panel.OverlayShortcut(0, earlier):
			step(-1)
		case panel.OverlayShortcut(0, later):
			step(1)
		case panel.OverlayShortcut(0, ui.KeyEqual), panel.OverlayShortcut(ui.Shift, ui.KeyEqual):
			v.step(v.zoom * viewerZoomStep)
		case panel.OverlayShortcut(0, ui.KeyMinus):
			v.step(v.zoom / viewerZoomStep)
		case panel.OverlayShortcut(0, ui.Key0):
			v.restore()
		}
		if v.shown != *index {
			v.shown = *index
			v.reset()
			v.smooth = false
		}
		item := o.Images[*index]
		panel.Label(item.Label)
		fitW, fitH := viewerFit(item.Source, w, h)

		// The stage takes the presses the controls above it leave.
		stage := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0)
		stage.HandleInput(func(ev ui.InputEvent) bool {
			switch ev.Kind {
			case ui.InputScroll:
				v.smooth = false
				v.zoomTo(v.zoom*float32(math.Exp(float64(-ev.DY*viewerWheelZoom))), ev.X-w/2, ev.Y-h/2)
			case ui.InputZoom:
				v.smooth = false
				v.zoomTo(v.zoom*ev.Scale, ev.X-w/2, ev.Y-h/2)
			default:
				return false
			}
			c.Invalidate()
			return true
		})
		dx, dy, held := stage.Dragged()
		touchhint.Held(held)
		if held && !v.held {
			v.panned = false
		}
		v.held = held
		if held && (dx != 0 || dy != 0) {
			v.x, v.y, v.panned, v.smooth = v.x+dx, v.y+dy, true, false
		}
		v.clamp(fitW, fitH, w, h)
		imgW, imgH := fitW*v.zoom, fitH*v.zoom
		left, top := w/2+v.x-imgW/2, h/2+v.y-imgH/2
		// The image glides to a zoom chosen by key or button, and follows
		// the hand at once otherwise.
		m := theme.Motion()
		glide := func(key string, to float32) float32 {
			if v.smooth {
				return stage.AnimateTo(key, to, ui.Animation{Duration: m.Reveal, Ease: m.EaseTravel})
			}
			return stage.Jump(key, to)
		}
		shownLeft, shownTop := glide("left", left), glide("top", top)
		shownW, shownH := glide("width", imgW), glide("height", imgH)
		if v.zoom > 1 {
			stage.Cursor(ui.CursorGrab)
		}
		if stage.Clicked() && !v.panned {
			px, py, _ := stage.PointerPosition()
			if px < left || px >= left+imgW || py < top || py >= top+imgH {
				close()
			}
		}
		stage.Children(func() {
			// The next image fades in over the one leaving.
			ui.Image(c.Key(*index), item.Source).Absolute().Left(shownLeft).Top(shownTop).Size(shownW, shownH).
				Fit(ui.Contain).Label(item.Label).IgnoreInput(true).
				Transition(ui.ElementTransition{Colors: true, Duration: m.State, Ease: m.Ease,
					Enter: &ui.Motion{}, Exit: &ui.Motion{}})
		})

		viewerControls(c, v, n, *index, step, close)
	})
	if wasOpen && !*open && p.Valid() && p.Dismissed() && o.OnDismiss != nil {
		o.OnDismiss()
	}
	return p
}

// viewerControls builds the close button, the zoom toolbar and, for several
// images, the previous and next buttons, each on a popup surface.
func viewerControls(c *ui.Context, v *viewerState, n, index int, step func(int), close func()) {
	msg := locale.M()
	button := func(icon *ui.SVG, label string, disabled bool, fn func()) {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: label, Icon: icon, IconOnly: true, Variant: mygoui.ButtonGhost,
			Size: theme.ButtonSmall, Disabled: disabled, OnClick: fn,
		})
	}
	pill := func() ui.Element {
		row := ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Space(0.5)).Padding(theme.Space(1))
		return anchored.Surface(c, row, theme.RadiusScale().XL)
	}
	edge := float32(theme.Space(4))

	pill().Absolute().Top(edge).End(edge).Children(func() {
		button(lucide.X(), msg.Close, false, close)
	})
	if n > 1 {
		// Previous sits at the start and Next at the end, each chevron
		// pointing outward, as the arrow keys step.
		side := func() ui.Element {
			return ui.Column(c).Absolute().Top(0).Bottom(0).Justify(ui.Center).PassThrough()
		}
		previous, next := lucide.ChevronLeft(), lucide.ChevronRight()
		if c.Direction() == ui.RightToLeft {
			previous, next = next, previous
		}
		side().Start(edge).Children(func() {
			pill().Children(func() {
				button(previous, msg.Previous, index == 0, func() { step(-1) })
			})
		})
		side().End(edge).Children(func() {
			pill().Children(func() {
				button(next, msg.Next, index == n-1, func() { step(1) })
			})
		})
	}
	bar := ui.Row(c).Absolute().Left(0).Right(0).Bottom(edge).Justify(ui.Center).PassThrough()
	bar.Children(func() {
		pill().Role(ui.RoleToolbar).Children(func() {
			button(lucide.ZoomOut(), msg.ZoomOut, false, func() { v.step(v.zoom / viewerZoomStep) })
			percent := strconv.Itoa(int(math.Round(float64(v.zoom)*100))) + "%"
			style.Text(ui.Text(c, percent), c, theme.XS()).FontFeatures("tnum").TextAlign(ui.Center).
				Width(theme.Scale(c, theme.Space(12))).TextColor(theme.Colors(c).MutedForeground)
			button(lucide.ZoomIn(), msg.ZoomIn, false, func() { v.step(v.zoom * viewerZoomStep) })
			button(lucide.RotateCcw(), msg.ResetZoom, false, v.restore)
		})
	})
}
