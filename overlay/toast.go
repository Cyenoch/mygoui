package overlay

import (
	"fmt"
	"math"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ToastVariant is HeroUI's toast color treatment.
type ToastVariant uint8

const (
	ToastDefault ToastVariant = iota
	ToastAccent
	ToastSuccess
	ToastWarning
	ToastDanger
)

// ToastPlacement is the window edge the toast region hugs.
type ToastPlacement uint8

const (
	ToastBottom ToastPlacement = iota
	ToastTop
)

// ToastOptions configures one HeroUI toast. Zero Timeout uses MyGo's four
// seconds; a negative Timeout keeps the toast until it is closed or updated.
type ToastOptions struct {
	// ID names the toast. Showing an ID that is already visible updates it
	// in place and restarts its lifetime; an empty ID is generated.
	ID          string
	Title       string
	Description string
	Action      string
	Variant     ToastVariant
	// Placement moves the whole region; the last shown toast's placement
	// wins.
	Placement ToastPlacement
	Timeout   time.Duration
	OnAction  func()
	OnClose   func()
}

// ToastItem is the caller-visible state of one retained toast.
type ToastItem struct {
	ID          string
	Title       string
	Description string
	Action      string
	Variant     ToastVariant
}

// ToastState is the caller-owned toast manager for one window. Its zero
// value is ready to use. It retains identities and payloads and mirrors the
// native toast collection's updates and closures; timers, the pause on hover
// or focus, and the three-visible cap belong to MyGo.
type ToastState struct {
	placement ToastPlacement
	ids       []string
	items     map[string]ToastItem
	seq       int
}

// Toast shows or updates a HeroUI toast. Call it from an event handler, not
// in the view body; install the region once with ToastState.View. It returns
// the toast's ID for Dismiss.
func Toast(c *ui.Context, s *ToastState, o ToastOptions) string {
	if c == nil || s == nil {
		return ""
	}
	s.placement = o.Placement
	if s.items == nil {
		s.items = make(map[string]ToastItem, 4)
	}
	id := o.ID
	if id == "" {
		s.seq++
		id = fmt.Sprintf("overlay-toast-%d", s.seq)
	}
	item := ToastItem{
		ID: id, Title: o.Title, Description: o.Description,
		Action: o.Action, Variant: o.Variant,
	}
	if _, exists := s.items[id]; !exists {
		s.ids = append(s.ids, id)
	}
	s.items[id] = item
	onClose := o.OnClose
	return c.AddToast(ui.Toast{
		ID:          id,
		Title:       o.Title,
		Description: o.Description,
		Type:        toastType(o.Variant),
		Timeout:     o.Timeout,
		Action:      o.Action,
		OnAction:    o.OnAction,
		OnClose: func() {
			s.remove(id)
			if onClose != nil {
				onClose()
			}
		},
	})
}

// Dismiss closes one toast by ID, running its OnClose.
func (s *ToastState) Dismiss(c *ui.Context, id string) {
	if c == nil || s == nil || id == "" {
		return
	}
	c.CloseToast(id)
}

// Items returns the retained toasts, oldest first.
func (s *ToastState) Items() []ToastItem {
	if s == nil || len(s.ids) == 0 {
		return nil
	}
	out := make([]ToastItem, 0, len(s.ids))
	for _, id := range s.ids {
		if item, ok := s.items[id]; ok {
			out = append(out, item)
		}
	}
	return out
}

func (s *ToastState) remove(id string) {
	if s == nil {
		return
	}
	delete(s.items, id)
	for i, v := range s.ids {
		if v == id {
			s.ids = append(s.ids[:i], s.ids[i+1:]...)
			return
		}
	}
}

// View installs the HeroUI toast region once per view, before modal views
// when a dialog should cover notifications. The newest three form a deck
// that expands on hover or keyboard focus; MyGo owns their lifetime.
func (s *ToastState) View(c *ui.Context) ui.Element {
	if s == nil {
		return ui.Element{}
	}
	w, h := c.Size()
	width := min(float32(460), max(float32(0), w-theme.Space32))
	room := max(float32(1), h-theme.Space32)
	quiet := theme.ReduceMotion(c)
	place := s.placement
	return ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		if width <= 0 {
			return
		}
		deck := ui.Box(c.Key("toast-deck")).Absolute().Left((w - width) / 2).
			Top(theme.Space16).Width(width).Height(room).PassThrough()
		geometry := ui.Local(deck, "geometry", func() toastDeckGeometry { return toastDeckGeometry{} })
		expanded := viewport.Hovered() || viewport.FocusWithin()
		var sizes [3]toastCardMetrics
		var expandedHeight, collapsedHeight float32
		// One line height for the deck, measured once: it depends on the
		// theme, not on a card's text.
		span := ui.Span{Font: theme.FontFamily(c), Size: theme.Scale(c, theme.SM().Size), Text: "M"}
		_, line := c.MeasureText(0, span)
		for i, t := range toasts {
			// A collapsed card behind the front one takes the front's
			// height, so only the front card and an expanded deck need theirs.
			own := expanded || i == len(toasts)-1
			sizes[i] = measureToast(c, t, width, line, own)
			if own {
				expandedHeight += sizes[i].height
			}
		}
		if len(toasts) > 0 {
			expandedHeight += float32(len(toasts)-1) * theme.Space12
			collapsedHeight = sizes[len(toasts)-1].height + float32(len(toasts)-1)*theme.Space12
		}
		hull := collapsedHeight
		if expanded {
			// Keep gaps hovered even if a shorter card arrives under a stationary pointer.
			geometry.hoverHeight = max(geometry.hoverHeight, expandedHeight)
			hull = geometry.hoverHeight
		} else {
			geometry.hoverHeight = 0
		}
		if len(toasts) == 0 {
			hull, geometry.hoverHeight = 0, 0
		}
		deck.Children(func() {
			envelope := ui.Box(c.Key("toast-envelope")).Absolute().Left(0).Width(width)
			envelope.Hovered()
			if len(toasts) == 0 || quiet {
				hull = envelope.Jump("height", hull)
			} else {
				hull = envelope.AnimateTo("height", hull, toastSpring)
			}
			envelope.Height(hull)
			if place == ToastTop {
				envelope.Top(0)
			} else {
				envelope.Bottom(0)
			}
			offset := expandedHeight
			for i, t := range toasts {
				offset -= sizes[i].height
				depth := len(toasts) - 1 - i
				distance := float32(depth) * theme.Space12
				height := sizes[len(toasts)-1].height
				scale := 1 - float32(depth)*0.04
				if expanded {
					distance, height, scale = offset, sizes[i].height, 1
				}
				cardW, cardH, top := width*scale, height*scale, distance
				if place != ToastTop {
					// The bottom edge stays put, as a transform around the
					// bottom center keeps it.
					top = room - distance - cardH
				}
				toastCard(c, t, toastCardLayout{
					full: width, width: cardW, height: cardH, top: top,
					metrics: sizes[i], placement: place, quiet: quiet,
				})
				offset -= theme.Space12
			}
		})
	})
}

// Only geometry survives builds; native transitions retain inert exit paint.
type toastDeckGeometry struct{ hoverHeight float32 }

type toastDragGeometry struct {
	offset    ui.Point
	dismissed bool
}

var toastSpring = ui.Animation{Spring: ui.SpringOptions{Duration: 320 * time.Millisecond}}

// toastContentFade is HeroUI's 200 ms fade of a card's content as the deck
// expands and collapses.
var toastContentFade = ui.Animation{Duration: 200 * time.Millisecond}

type toastCardMetrics struct {
	height, bodyWidth, actionWidth, line float32
}

// toastCardLayout is one card's slot. full is the deck's width, which a card
// behind the front one is drawn narrower than; the deck shrinks that card by
// animating its own box rather than a transform, because a scaled shadow makes
// the software rasterizer integrate every pixel of its blur on each frame.
type toastCardLayout struct {
	full, width, height, top float32
	metrics                  toastCardMetrics
	placement                ToastPlacement
	quiet                    bool
}

// measureToast measures one toast card at width. line is the deck's line
// height, and withHeight tells whether the card's own body height is needed:
// a collapsed card behind the front one is drawn with the front's height.
func measureToast(c *ui.Context, t ui.Toast, width, line float32, withHeight bool) toastCardMetrics {
	span := ui.Span{Font: theme.FontFamily(c), Size: theme.Scale(c, theme.SM().Size)}
	m := toastCardMetrics{line: line}
	m.bodyWidth = width - theme.Space32 - theme.Space24 - theme.Space6
	row := float32(theme.Space24)
	if t.Action != "" {
		span.Text, span.Weight = t.Action, theme.WeightMedium
		m.actionWidth, _ = c.MeasureText(0, span)
		m.actionWidth += theme.Space32
		m.bodyWidth -= m.actionWidth + theme.Space6
		row = max(row, theme.Scale(c, theme.HeightButtonMD))
	}
	m.bodyWidth = max(float32(1), m.bodyWidth)
	var body float32
	if withHeight {
		for i, text := range [2]string{t.Title, t.Description} {
			if text == "" {
				continue
			}
			span.Text, span.Weight = text, theme.WeightNormal
			if i == 0 {
				span.Weight = theme.WeightMedium
			}
			_, height := c.MeasureText(m.bodyWidth, span)
			body += height
		}
	}
	m.height = max(body, row) + theme.Space24
	return m
}

func toastType(v ToastVariant) string {
	switch v {
	case ToastAccent:
		return "accent"
	case ToastSuccess:
		return "success"
	case ToastWarning:
		return "warning"
	case ToastDanger:
		return "danger"
	}
	return "default"
}

func variantOf(kind string) ToastVariant {
	switch kind {
	case "accent":
		return ToastAccent
	case "success":
		return ToastSuccess
	case "warning":
		return ToastWarning
	case "danger":
		return ToastDanger
	}
	return ToastDefault
}

func toastCard(c *ui.Context, t ui.Toast, layout toastCardLayout) {
	pal := theme.Colors(c)
	variant := variantOf(t.Type)
	titleInk, iconInk := toastInk(pal, variant)
	label := t.Title
	if label == "" {
		label = t.Description
	}
	parts := ui.ToastBase(c, t)
	root := parts.Root.Absolute().Label(label).
		Radius(overlayRadius()).
		Background(pal.Surface).TextColor(pal.OverlayForeground)
	root = panelElevation(c, root)
	top, height := layout.top, layout.height
	left := (layout.full - layout.width) / 2
	width := layout.width
	// HeroUI enters and leaves the deck by the card's own height, from the
	// window edge: translateY(100%) below and -100% above. The deck keeps a
	// 16 DIP gutter, so include it or the first frame leaves a sliver of the
	// card visible at the edge. The entrance is declared geometry, which MyGo
	// re-lays after transitions, and the exit is the transition's own
	// translation, which the rasterizer folds: a scale here would put the
	// card's shadow on the per-pixel transformed path every animated frame.
	dir := height + theme.Space16
	if layout.placement == ToastTop {
		dir = -height - theme.Space16
	}
	entry := ui.Local(root, "entry", func() bool { return false })
	if !layout.quiet {
		if !*entry {
			// Only the arrival translation and fade animate. A new card is
			// already at its measured size: animating height/width as well
			// made it unfold from a sliver before settling.
			root.Jump("top", top+dir)
			root.Jump("left", left)
			root.Jump("width", width)
			root.Jump("height", height)
			*entry = true
		}
		top = root.AnimateTo("top", top, toastSpring)
		left = root.AnimateTo("left", left, toastSpring)
		width = root.AnimateTo("width", width, toastSpring)
		height = root.AnimateTo("height", height, toastSpring)
		root.Transition(ui.ElementTransition{
			Colors: true, Animation: &toastSpring,
			Enter: &ui.Motion{Opacity: 0},
			Exit:  &ui.Motion{Y: dir, Opacity: 0},
		})
	} else {
		root.Jump("top", top)
		root.Jump("left", left)
		root.Jump("width", width)
		root.Jump("height", height)
		*entry = true
	}
	root.Left(left).Top(top).Width(width).Height(height)
	drag := ui.Local(root, "drag", func() toastDragGeometry { return toastDragGeometry{} })
	root.Draggable(&drag.offset, ui.DraggableOptions{
		Axis: ui.Horizontal, Animation: toastSpring,
		Snap: func(projected ui.Point) ui.Point {
			if float32(math.Abs(float64(projected.X))) >= max(float32(56), layout.width*0.28) {
				drag.dismissed = true
				return ui.Point{X: float32(math.Copysign(float64(layout.width+32), float64(projected.X)))}
			}
			drag.dismissed = false
			return ui.Point{}
		},
	})
	shown := root.DragOffset()
	if !root.Dragging() && drag.dismissed && math.Abs(float64(shown.X-drag.offset.X)) < 0.5 {
		c.CloseToast(t.ID)
		return
	}
	root.Children(func() {
		row := ui.Row(c.Key("toast-content")).Fill().AlignItems(ui.Start).
			Gap(theme.Space6).Padding(theme.Space12, theme.Space16).Clip()
		// HeroUI fades a covered card's content out, leaving only the top
		// sliver of its face; the deck's expansion reveals it again.
		reveal := float32(1)
		if layout.width < layout.full {
			reveal = 0
		}
		row.Opacity(row.AnimateTo("reveal", reveal, toastContentFade)).Children(func() {
			// Indicator: 16 DIP glyph in a 4 DIP gutter (p-1).
			box := ui.Box(c).Size(theme.Space24, theme.Space24).Shrink(0).
				Center().Padding(theme.Space4).TextColor(iconInk)
			box.Children(func() {
				ui.Icon(c, toastGlyph(variant)).
					Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
			})
			content := ui.Column(c).Width(layout.metrics.bodyWidth).MinWidth(0).AlignSelf(ui.Center).Gap(0)
			content.Children(func() {
				if t.Title != "" {
					title := style.Text(ui.Text(c, t.Title), c, theme.SM())
					title.FontWeight(theme.WeightMedium).FixedLineHeight(layout.metrics.line).TextColor(titleInk).Wrap()
				}
				if t.Description != "" {
					style.Text(ui.Text(c, t.Description), c, theme.SM()).
						FixedLineHeight(layout.metrics.line).TextColor(pal.Muted).Wrap()
				}
			})
			if t.Action != "" {
				action := parts.ActionButton().Label(t.Action).Width(layout.metrics.actionWidth).
					MinHeight(theme.Scale(c, theme.HeightButtonMD)).PaddingX(theme.Space16).
					Center().Shrink(0).Radius(theme.RadiusScale().XXXL).FocusRing(true)
				bg, hover, fg := toastActionColors(pal, variant)
				action = style.Background(action, theme.RadiusScale().XXXL, bg, hover, hover).TextColor(fg)
				action.Draggable(nil, ui.DraggableOptions{Axis: ui.Horizontal, Bounds: &ui.Rect{}, Rigid: true})
				action.Children(func() {
					style.Text(ui.Text(c, t.Action), c, theme.SM()).
						FontWeight(theme.WeightMedium).SingleLine()
				})
			}
		})
		// The close affordance floats outside the content row, without stealing
		// width from the message or action. Keep the root unclipped for its ring.
		size := theme.Scale(c, theme.Space20)
		glyph := theme.Scale(c, theme.IconXS)
		close := parts.CloseButton().Label(locale.M().Close).Absolute().
			Top(-theme.Space4).Right(-theme.Space4).
			Size(size, size).Padding(0).Center().
			Radius(theme.RadiusScale().XL).FocusRing(true).TextColor(pal.Muted)
		close.Draggable(nil, ui.DraggableOptions{Axis: ui.Horizontal, Bounds: &ui.Rect{}, Rigid: true})
		if root.Hovered() || root.FocusWithin() {
			close = close.Background(pal.Default)
			close.Children(func() {
				ui.Icon(c, lucide.X()).Size(glyph, glyph).Shrink(0)
			})
		}
	})
}

func toastGlyph(v ToastVariant) *ui.SVG {
	switch v {
	case ToastSuccess:
		return lucide.CircleCheck()
	case ToastWarning:
		return lucide.TriangleAlert()
	case ToastDanger:
		return lucide.CircleX()
	}
	return lucide.Info()
}

func toastInk(p theme.Palette, v ToastVariant) (title, icon ui.Color) {
	switch v {
	case ToastAccent:
		return p.AccentSoftForeground, p.AccentSoftForeground
	case ToastSuccess:
		return p.SuccessSoftForeground, p.SuccessSoftForeground
	case ToastWarning:
		return p.WarningSoftForeground, p.WarningSoftForeground
	case ToastDanger:
		return p.DangerSoftForeground, p.DangerSoftForeground
	}
	return p.OverlayForeground, p.OverlayForeground
}

func toastActionColors(p theme.Palette, v ToastVariant) (bg, hover, fg ui.Color) {
	switch v {
	case ToastSuccess:
		return p.Success, p.SuccessHover, p.SuccessForeground
	case ToastWarning:
		return p.Warning, p.WarningHover, p.WarningForeground
	case ToastDanger:
		return p.Danger, p.DangerHover, p.DangerForeground
	}
	return p.Accent, p.AccentHover, p.AccentForeground
}
