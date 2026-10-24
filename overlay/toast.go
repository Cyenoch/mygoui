package overlay

import (
	"fmt"
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

// View installs the HeroUI toast region. Call it once in the view, before
// modal views when a dialog should cover notifications. MyGo passes the
// newest three toasts, oldest first; HeroUI stacks the newest at the bottom
// edge and, for ToastTop, at the top.
func (s *ToastState) View(c *ui.Context) ui.Element {
	if s == nil {
		return ui.Element{}
	}
	place := s.placement
	return ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		w, _ := c.Size()
		viewport.AlignItems(ui.Center).Gap(theme.Space12).Padding(theme.Space16)
		if place == ToastTop {
			viewport.Justify(ui.Start)
		} else {
			viewport.Justify(ui.End)
		}
		width := min(float32(460), max(float32(0), w-theme.Space32))
		if width <= 0 {
			return
		}
		for i := range toasts {
			at := i
			if place == ToastTop {
				at = len(toasts) - 1 - i
			}
			toastCard(c, toasts[at], width, place)
		}
	})
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

func toastCard(c *ui.Context, t ui.Toast, width float32, place ToastPlacement) {
	pal := theme.Colors(c)
	variant := variantOf(t.Type)
	titleInk, iconInk := toastInk(pal, variant)
	parts := ui.ToastBase(c, t)
	root := parts.Root.Width(width).Row().AlignItems(ui.Start).Shrink(0).
		Gap(theme.Space6).Padding(theme.Space12, theme.Space16).
		Radius(overlayRadius()).
		Background(pal.Surface).TextColor(pal.OverlayForeground)
	root = panelElevation(c, root)
	if !theme.ReduceMotion(c) {
		dir := float32(16)
		if place == ToastTop {
			dir = -16
		}
		// MyGo 0.6.3 re-lays absolute overlays after transitions
		// (layoutAbsolute): the card's scale and fade show, but its 16 DIP
		// lift is re-laid away and the stack position snaps when a toast
		// joins or leaves instead of sliding with it.
		root.Transition(ui.ElementTransition{
			Duration: 250 * time.Millisecond,
			Enter:    &ui.Motion{Y: dir, Scale: 0.96, Opacity: 0},
			Exit:     &ui.Motion{Y: dir, Scale: 0.96, Opacity: 0},
		})
	}
	root.Children(func() {
		// Indicator: 16 DIP glyph in a 4 DIP gutter (p-1).
		box := ui.Box(c).Size(theme.Space24, theme.Space24).Shrink(0).
			Center().Padding(theme.Space4).TextColor(iconInk)
		box.Children(func() {
			ui.Icon(c, toastGlyph(variant)).
				Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
		})
		content := ui.Column(c).Grow(1).MinWidth(0).AlignSelf(ui.Center).Gap(0)
		content.Children(func() {
			if t.Title != "" {
				title := style.Text(ui.Text(c, t.Title), c, theme.SM())
				title.FontWeight(theme.WeightMedium).TextColor(titleInk).Wrap()
			}
			if t.Description != "" {
				style.Text(ui.Text(c, t.Description), c, theme.SM()).TextColor(pal.Muted).Wrap()
			}
		})
		if t.Action != "" {
			action := parts.ActionButton().Label(t.Action).
				MinHeight(theme.Scale(c, theme.HeightButtonMD)).PaddingX(theme.Space16).
				Center().Shrink(0).Radius(theme.RadiusScale().XXXL).FocusRing(true)
			bg, hover, fg := toastActionColors(pal, variant)
			action = style.Background(action, theme.RadiusScale().XXXL, bg, hover, hover).TextColor(fg)
			action.Children(func() {
				style.Text(ui.Text(c, t.Action), c, theme.SM()).
					FontWeight(theme.WeightMedium).SingleLine()
			})
		}
		// HeroUI reveals the toast close button on hover or focus: a 20 DIP
		// icon-only button in flow at the row's top-end, inside the card. It
		// holds one lucide X at 12 DIP, centered. The build chooses the
		// revealed children and background from live state, as before: no
		// Draw callback, and so no Element that could expire while a leaving
		// toast is painted. The box keeps its size at rest, so revealing it
		// shifts neither the title nor the action.
		size := theme.Scale(c, theme.Space20)
		glyph := theme.Scale(c, theme.IconXS)
		close := parts.CloseButton().Label(locale.M().Close).
			Size(size, size).Shrink(0).Padding(0).Center().
			Radius(theme.RadiusScale().XL).FocusRing(true).
			AlignSelf(ui.Start).TextColor(pal.Muted)
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
