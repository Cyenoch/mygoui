package overlay

import (
	"math"
	"slices"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type ToastKind uint8

const (
	ToastDefault ToastKind = iota
	ToastSuccess
	ToastInfo
	ToastWarning
	ToastError
	ToastLoading
)

type ToastOptions struct {
	ID, Title, Description, Action string
	Kind                           ToastKind
	// Announcement replaces automatic speech while retaining the exact visible
	// title and description in the native toast's Data. Nil reads both normally.
	Announcement *string
	// Zero uses the native default, except loading, which remains until updated
	// or closed. Negative is sticky; a positive loading timeout is explicit.
	Timeout           time.Duration
	OnAction, OnClose func()
}

type toastDisplay struct{ title, description string }

func visibleToast(t myui.Toast) myui.Toast {
	if display, ok := t.Data.(toastDisplay); ok {
		t.Title, t.Description = display.title, display.description
	}
	return t
}

// ToastState bounds retained native toasts to 32, of which MyGo shows the newest
// three. One state belongs to one window. It retains identities only: the native
// toast collection owns payloads, callbacks, updates, pause and expiration.
type ToastState struct{ ids []string }

const ToastCapacity = 32

// Toast is called once from a UI event, never unconditionally during rendering.
// A same-ID update keeps native element and gesture identity, refreshes native
// lifetime, and announces the updated message without replaying its entrance.
func Toast(c *myui.Context, s *ToastState, o ToastOptions) string {
	if s == nil {
		return ""
	}
	if o.ID == "" || !slices.Contains(s.ids, o.ID) {
		if len(s.ids) >= ToastCapacity {
			c.CloseToast(s.ids[0])
		}
	}
	kind := ""
	switch o.Kind {
	case ToastSuccess:
		kind = "success"
	case ToastInfo:
		kind = "info"
	case ToastWarning:
		kind = "warning"
	case ToastError:
		kind = "error"
	case ToastLoading:
		kind = "loading"
		if o.Timeout == 0 {
			o.Timeout = -1
		}
	}
	id := o.ID
	title, description := o.Title, o.Description
	var data any
	if o.Announcement != nil {
		data = toastDisplay{title: title, description: description}
		title, description = *o.Announcement, ""
	}
	onClose := o.OnClose
	id = c.AddToast(myui.Toast{ID: id, Title: title, Description: description, Data: data, Type: kind, Timeout: o.Timeout, Action: o.Action, OnAction: o.OnAction,
		OnClose: func() {
			if i := slices.Index(s.ids, id); i >= 0 {
				s.ids = slices.Delete(s.ids, i, i+1)
			}
			if onClose != nil {
				onClose()
			}
		}})
	if !slices.Contains(s.ids, id) {
		s.ids = append(s.ids, id)
	}
	return id
}

func DismissToast(c *myui.Context, id string) { c.CloseToast(id) }

// The deck is geometry, not another toast store. In particular it keeps no
// message, callback, clock, Context or Element after the build. Native keyed
// transition snapshots retain an inert painted card after CloseToast has already
// released its data, including when the last card closes.
type toastDeckGeometry struct{ hoverHeight float32 }
type toastDragGeometry struct {
	offset    myui.Point
	dismissed bool
}

var toastSpring = myui.Animation{Spring: myui.SpringOptions{Duration: 320 * time.Millisecond}}

// Toaster is installed once per window. Direct Context.AddToast callers share
// this same viewport. Install before modal overlays when a modal should cover
// notifications; native modal focus containment remains authoritative.
func Toaster(c *myui.Context) myui.Element {
	w, h := c.Size()
	gutter := min(float32(theme.Space16), max(float32(0), min(w, h)/10))
	width := max(float32(1), min(float32(384), w-2*gutter))
	room := max(float32(1), h-2*gutter)
	quiet := theme.ReduceMotion(c)
	return myui.ToastViewportBase(c, func(viewport myui.Element, toasts []myui.Toast) {
		// The native viewport is transparent to input outside the actual deck.
		// A permanent parent also lets native exit snapshots finish when empty.
		deck := myui.Box(c.Key("toast-deck")).Absolute().Bottom(gutter).
			Width(width).Height(room).PassThrough()
		if c.Direction() == myui.RightToLeft {
			deck.Left(gutter)
		} else {
			deck.Right(gutter)
		}
		geometry := myui.Local(deck, "geometry", func() toastDeckGeometry { return toastDeckGeometry{} })
		expanded := viewport.Hovered() || viewport.FocusWithin()
		var sizes [3]toastCardMetrics
		var heights [3]float32
		var expandedHeight, collapsedHeight float32
		if len(toasts) > 0 {
			// A tall message scrolls inside its own card, so every visible card
			// and its controls fit even at 200% text in a short window.
			limit := max(float32(1), (room-float32(len(toasts)-1)*theme.Space10)/float32(len(toasts)))
			for i, t := range toasts {
				sizes[i] = toastMetrics(c, visibleToast(t), width, limit)
				heights[i] = sizes[i].height
				expandedHeight += heights[i]
			}
			expandedHeight += float32(len(toasts)-1) * theme.Space10
			collapsedHeight = heights[len(toasts)-1] + float32(len(toasts)-1)*theme.Space12
		}
		hull := collapsedHeight
		if expanded {
			// Hold the entire expanded hull, including gaps, while interacting.
			// Removing a tall card or inserting a shorter one cannot remove the
			// hit region under a stationary pointer and collapse/reopen the deck.
			geometry.hoverHeight = min(room, max(geometry.hoverHeight, expandedHeight))
			hull = geometry.hoverHeight
		} else {
			geometry.hoverHeight = 0
		}
		if len(toasts) == 0 {
			hull, geometry.hoverHeight = 0, 0
		}
		deck.Children(func() {
			envelope := myui.Box(c.Key("toast-envelope")).Absolute().Bottom(0).
				Width(width).Label(locales.ComponentNotifications())
			envelope.Hovered() // Hover-only: no click, focus or drag on the gaps.
			if len(toasts) == 0 {
				// Exits are native inert snapshots. A leftover hover hull must
				// not cover otherwise clickable page content while they leave.
				hull = envelope.Jump("height", 0)
			} else if !quiet {
				hull = envelope.AnimateTo("height", hull, toastSpring)
			}
			envelope.Height(hull)
			bottom := expandedHeight
			for i, t := range toasts {
				bottom -= heights[i]
				depth := len(toasts) - 1 - i
				offset, height, scale := float32(depth)*theme.Space12, heights[len(toasts)-1], 1-float32(depth)*0.04
				if expanded {
					offset, height, scale = bottom, heights[i], 1
				}
				toastCard(c, visibleToast(t), sizes[i], width, height, offset, scale, room, quiet)
				bottom -= theme.Space10
			}
		})
	})
}

// Measurements use the same native font, weights and line spacing as paint.
// MeasureText inherits the theme's base metrics even for smaller spans; using
// that exact base line height for every painted text run also covers fallback
// fonts without estimating line counts or reading last frame's bounds.
type toastCardMetrics struct {
	height, bodyWidth, icon, close, padding, paddingY, row float32
	font, actionFont, line                                 float32
	family                                                 string
}

func toastMetrics(c *myui.Context, t myui.Toast, width, limit float32) toastCardMetrics {
	m := toastCardMetrics{
		icon:       theme.Scale(c, 18),
		close:      max(float32(theme.HeightButton), theme.Scale(c, 16)),
		padding:    theme.Space14,
		paddingY:   theme.Space10,
		font:       theme.Scale(c, 16),
		actionFont: theme.Scale(c, 13),
		family:     c.Theme().Font,
	}
	if width < 240 {
		m.padding = theme.Space10
	}
	_, m.line = c.MeasureText(0, myui.Span{Text: "M", Font: m.family, Size: m.font, Weight: 400})
	m.close = min(m.close, max(float32(1), limit-2*theme.HairlineWidth))
	m.row = min(max(m.line, max(m.icon, m.close)), max(float32(1), limit-2*theme.HairlineWidth))
	m.paddingY = min(m.paddingY, max(float32(0), (limit-m.row-2*theme.HairlineWidth)/2))
	m.bodyWidth = max(float32(1), width-2*m.padding-2*theme.HairlineWidth-m.close-theme.Space10)
	if t.Type != "" {
		m.bodyWidth = max(float32(1), m.bodyWidth-m.icon-theme.Space10)
	}
	var body float32
	if t.Title != "" {
		_, body = c.MeasureText(m.bodyWidth, myui.Span{Text: t.Title, Font: m.family, Size: m.font, Weight: 700})
		body += max(float32(0), m.row-m.line)
	}
	if t.Description != "" {
		_, description := c.MeasureText(m.bodyWidth, myui.Span{Text: t.Description, Font: m.family, Size: m.font, Weight: 400})
		if t.Title != "" {
			body += theme.Space4
		} else {
			description += max(float32(0), m.row-m.line)
		}
		body += description
	}
	if t.Action != "" {
		_, action := c.MeasureText(max(float32(1), m.bodyWidth-theme.Space16-2),
			myui.Span{Text: t.Action, Font: m.family, Size: m.actionFont, Weight: 500})
		if body > 0 {
			body += theme.Space10
		}
		body += max(float32(theme.HeightButton), action+theme.Space12+2)
	}
	m.height = min(limit, max(body, m.row)+2*m.paddingY+2*theme.HairlineWidth)
	return m
}

func toastCard(c *myui.Context, t myui.Toast, m toastCardMetrics, width, height, bottom, scale, room float32, quiet bool) {
	p := theme.Colors(c)
	parts := myui.ToastBase(c, t)
	root := parts.Root.Absolute().Left(0).Size(width, height).
		Row().AlignItems(myui.Start).Padding(m.paddingY, m.padding).Gap(theme.Space10).
		Font(m.family).FontSize(m.font).FixedLineHeight(m.line).
		Radius(theme.RadiusCard).Background(p.Popover).TextColor(p.PopoverForeground).
		Border(theme.HairlineWidth, p.Border).Shadow(0, 4, 16, -4, p.Shadow).
		Label(locales.ComponentToastNamed(t.ID)).Description(t.Title).TransformOrigin(0.5, 0)
	top := room - height - bottom
	if !quiet {
		// The overlay reattaches absolute descendants after native layout
		// transitions. Animate their declared geometry, so that reattachment
		// receives the current spring values instead of overwriting them with
		// the final Bottom/Height targets on the first frame.
		top = root.AnimateTo("top", top, toastSpring)
		height = root.AnimateTo("height", height, toastSpring)
		root.Transition(myui.ElementTransition{Colors: true, Animation: &toastSpring,
			Enter: &myui.Motion{Y: 24, Scale: 0.96}, Exit: &myui.Motion{Y: 16, Scale: 0.96}}).
			AnimatePose(myui.Pose{Scale: scale}, toastSpring)
	} else {
		root.Scale(scale, scale)
	}
	root.Top(top).Height(height)
	drag := myui.Local(root, "drag", func() toastDragGeometry { return toastDragGeometry{} })
	root.Draggable(&drag.offset, myui.DraggableOptions{
		Axis: myui.Horizontal, Animation: toastSpring,
		Snap: func(projected myui.Point) myui.Point {
			if float32(math.Abs(float64(projected.X))) >= max(float32(56), width*0.28) {
				drag.dismissed = true
				return myui.Point{X: float32(math.Copysign(float64(width+32), float64(projected.X)))}
			}
			drag.dismissed = false
			return myui.Point{}
		},
	})
	// Read the maintained drag's painted offset to rebuild while settling, so
	// dismissal can finish without a view timer or a second animation engine.
	shown := root.DragOffset()
	if !root.Dragging() && drag.dismissed &&
		float32(math.Abs(float64(shown.X-drag.offset.X))) < 0.5 {
		c.CloseToast(t.ID)
		return
	}
	root.Children(func() {
		if t.Type != "" {
			myui.Box(c).Size(m.icon, m.row).Center().Shrink(0).Children(func() { toastIcon(c, t.Type, m.icon) })
		}
		// Native selectable text and buttons are the innermost input targets.
		// Their drags never become background swipes; a drag never clicks an action.
		myui.Scroll(c.Key("toast-content")).Width(m.bodyWidth).FillHeight().MinHeight(0).Shrink(0).
			Children(func() {
				myui.Column(c).Width(m.bodyWidth).MinWidth(0).MinHeight(m.row).Children(func() {
					firstLinePadding := max(float32(0), (m.row-m.line)/2)
					if t.Title != "" {
						myui.Text(c, t.Title).FillWidth().MinWidth(0).Shrink(0).FontWeight(700).
							Padding(firstLinePadding, 0).Selectable()
					}
					if t.Description != "" {
						description := myui.Text(c, t.Description).FillWidth().MinWidth(0).Shrink(0).
							FontWeight(400).TextColor(p.MutedForeground).Selectable()
						if t.Title == "" {
							description.Padding(firstLinePadding, 0)
						} else {
							description.Margin(theme.Space4, 0, 0, 0)
						}
					}
					if t.Action != "" {
						c.Key("toast-action")
						b := parts.ActionButton().Label(t.Action).AlignSelf(myui.Start).MaxWidth(m.bodyWidth).
							Padding(theme.Space6, theme.Space8).Shrink(0).
							MinHeight(theme.HeightButton).Radius(theme.RadiusControl).Border(1, p.Border).
							Background(p.Secondary).TextColor(p.SecondaryForeground)
						if t.Title != "" || t.Description != "" {
							b.Margin(theme.Space10, 0, 0, 0)
						}
						// A rigid native drag stays in the button and suppresses
						// release activation after movement, including return-inside.
						b.Draggable(nil, myui.DraggableOptions{Axis: myui.Horizontal, Bounds: &myui.Rect{}, Rigid: true})
						if b.Hovered() {
							b.Background(p.Accent)
						}
						b.Children(func() {
							myui.Text(c, t.Action).FontSize(m.actionFont).FontWeight(500).MinWidth(0).Shrink(1)
						})
					}
				})
			})
		myui.Box(c).Size(m.close, m.row).Center().Shrink(0).Children(func() {
			c.Key("toast-close")
			b := parts.CloseButton().Label(locales.ComponentCloseToast()).Size(m.close, m.close).Padding(0).Center().Shrink(0).
				Radius(theme.RadiusControl).TextColor(p.MutedForeground)
			b.Draggable(nil, myui.DraggableOptions{Axis: myui.Horizontal, Bounds: &myui.Rect{}, Rigid: true})
			if b.Hovered() {
				b.Background(p.Accent).TextColor(p.PopoverForeground)
			}
			b.Children(func() {
				myui.Icon(c, lucide.X()).Size(theme.Scale(c, 16), theme.Scale(c, 16)).Shrink(0)
			})
		})
	})
}

func toastIcon(c *myui.Context, kind string, size float32) {
	p := theme.Colors(c)
	ink, glyph, label := p.Info, lucide.Info(), locales.ComponentInfo()
	switch kind {
	case "success":
		ink, glyph, label = p.Success, lucide.CircleCheck(), locales.ComponentSuccess()
	case "warning":
		ink, glyph, label = p.Warning, lucide.TriangleAlert(), locales.ComponentWarning()
	case "error":
		ink, glyph, label = p.Destructive, lucide.CircleX(), locales.ComponentError()
	case "loading":
		mygoui.Spinner(c, mygoui.SpinnerOptions{Size: size / theme.TextScale(c), Label: locales.ComponentLoading()})
		return
	}
	myui.Icon(c, glyph).Size(size, size).Shrink(0).TextColor(ink).Label(label)
}
