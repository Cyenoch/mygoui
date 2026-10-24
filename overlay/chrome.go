package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
)

// ModalBackdrop is HeroUI's modal backdrop treatment.
type ModalBackdrop uint8

const (
	// ModalBackdropOpaque dims what is behind the overlay.
	ModalBackdropOpaque ModalBackdrop = iota
	// ModalBackdropBlur dims and blurs what is behind the overlay, as
	// Tailwind's backdrop-blur-md (12 DIPs sigma).
	ModalBackdropBlur
	// ModalBackdropTransparent leaves what is behind visible and undimmed.
	ModalBackdropTransparent
)

// backdropBlur is Tailwind backdrop-blur-md, in DIPs of Gaussian sigma.
const backdropBlur = 12

// paintBackdrop applies a HeroUI backdrop variant to back. The caller keeps
// the layout (centering, padding, orientation) so Modal, Drawer, Dialog and
// BottomSheet share one paint contract.
func paintBackdrop(c *ui.Context, back ui.Element, variant ModalBackdrop) ui.Element {
	t := theme.Colors(c)
	switch variant {
	case ModalBackdropTransparent:
		return back
	case ModalBackdropBlur:
		// Material replaces the background, so the dim is painted over the
		// blur with the element's paint callback.
		back = back.Material(glass.Blur{Radius: backdropBlur})
		return back.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.Backdrop, 0)
		})
	default:
		return back.Background(t.Backdrop)
	}
}

// overlayRadius is min(32px, --radius-3xl): HeroUI's floating surface radius.
func overlayRadius() float32 {
	return min(theme.Space32, theme.RadiusScale().XXXL)
}

// panelElevation adds the overlay shadow, or the hairline inset ring dark
// panels use instead of a shadow.
func panelElevation(c *ui.Context, e ui.Element) ui.Element {
	return style.Shadow(e, c, theme.OverlayShadows(c.Theme().Dark))
}

// chrome describes the common scrollable dialog content: heading, body,
// footer and close button.
type chrome struct {
	title, description   string
	header, body, footer func()
	bodyLabel            string
	icon                 func()
	hideClose            bool
	onClose              func()

	titleStyle       theme.TextStyle
	descriptionStyle theme.TextStyle
	bodyInk          ui.Color
	padding          float32
	headerGap        float32
	headerBody       float32
	bodyFooter       float32
	closeInset       float32
	headingLevel     int
}

// fillPanel builds p's sections inside the already styled panel element.
// Header and body exist only when they have content, so the CSS sibling
// spacing (header + body 8, body + footer 20) applies without empty gaps.
func fillPanel(c *ui.Context, e ui.Element, p chrome) {
	t := theme.Colors(c)
	padding := p.padding
	if padding == 0 {
		padding = theme.Space24
	}
	headerGap := p.headerGap
	if headerGap == 0 {
		headerGap = theme.Space12
	}
	headerBody := p.headerBody
	if headerBody == 0 {
		headerBody = theme.Space8
	}
	bodyFooter := p.bodyFooter
	if bodyFooter == 0 {
		bodyFooter = theme.Space20
	}
	closeInset := p.closeInset
	if closeInset == 0 {
		closeInset = theme.Space16
	}
	level := p.headingLevel
	if level == 0 {
		level = 2
	}
	hasHeader := p.header != nil || p.title != "" || p.description != "" || p.icon != nil
	hasBody := p.body != nil
	hasFooter := p.footer != nil
	hasClose := p.onClose != nil && !p.hideClose

	e.Column().MinWidth(0).Padding(padding)
	e.Children(func() {
		if hasHeader {
			h := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Gap(headerGap)
			if hasBody {
				h.Margin(0, 0, headerBody, 0)
			} else if hasFooter {
				h.Margin(0, 0, bodyFooter, 0)
			}
			if hasClose {
				h.PaddingEnd(theme.Space32)
			}
			h.Children(func() {
				if p.header != nil {
					p.header()
					return
				}
				if p.icon != nil {
					p.icon()
				}
				inner := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4)
				inner.Children(func() {
					if p.title != "" {
						title := style.Text(ui.Text(c, p.title), c, p.titleStyle)
						title.FontWeight(theme.WeightMedium).TextColor(t.OverlayForeground).
							Role(ui.RoleHeading).Level(level)
					}
					if p.description != "" {
						style.Text(ui.Text(c, p.description), c, p.descriptionStyle).TextColor(t.Muted)
					}
				})
			})
		}
		if hasBody {
			b := ui.Scroll(c).FillWidth().Grow(1).MinHeight(0)
			if hasFooter {
				b.Margin(0, 0, bodyFooter, 0)
			}
			if p.bodyLabel != "" {
				b.Label(p.bodyLabel)
			}
			if p.bodyInk.A != 0 {
				b.TextColor(p.bodyInk)
			}
			b.Children(p.body)
		}
		if hasFooter {
			f := ui.Row(c).FillWidth().Justify(ui.End).AlignItems(ui.Center).Gap(theme.Space8).Shrink(0)
			f.Children(p.footer)
		}
		if hasClose {
			close := mygoui.CloseButton(c, mygoui.CloseButtonOptions{
				Label:   locale.M().Close,
				OnClick: p.onClose,
			})
			close.Absolute().Top(closeInset).End(closeInset)
		}
	})
}

// softPair returns the HeroUI soft surface and its foreground for a status.
func softPair(c *ui.Context, status mygoui.StatusColor) (bg, fg ui.Color) {
	t := theme.Colors(c)
	if status == mygoui.ColorDefault {
		return t.Default, t.Foreground
	}
	p := mygoui.StatusColors(c, status)
	return p.Soft, p.SoftForeground
}

// statusIcon is HeroUI's default 40 DIP status medallion (size-10,
// rounded-3xl) with a 20 DIP glyph.
func statusIcon(c *ui.Context, status mygoui.StatusColor, glyph *ui.SVG) func() {
	return func() {

		bg, fg := softPair(c, status)
		box := ui.Box(c).Size(theme.Space40, theme.Space40).Shrink(0).
			Center().Radius(overlayRadius()).Background(bg).TextColor(fg)
		box.Children(func() {
			ui.Icon(c, glyph).Size(theme.Scale(c, theme.IconLG), theme.Scale(c, theme.IconLG)).Shrink(0)
		})
	}
}

// resolveSide resolves the effective placement side, as anchored.Apply does.
func resolveSide(p mygoui.Placement, fallback mygoui.Side) mygoui.Side {
	if p.Side == mygoui.SideAuto {
		return fallback
	}
	return p.Side
}

// paintArrow draws HeroUI's 12 DIP overlay arrow as a rounded square turned
// 45°, half of it outside the panel's edge facing the anchor. The panel must
// not Clip, or the protruding half disappears.
func paintArrow(p *ui.Painter, r ui.Rect, side mygoui.Side, fill ui.Color) {
	const size = 9
	cx := r.X + r.W/2
	var cy float32
	switch side {
	case mygoui.SideTop:
		cy = r.Y + r.H
	case mygoui.SideLeft:
		cx, cy = r.X+r.W, r.Y+r.H/2
	case mygoui.SideRight:
		cx, cy = r.X, r.Y+r.H/2
	default:
		cy = r.Y
	}
	p.Transform(ui.Rotation(45), cx, cy, func() {
		p.Fill(ui.Rect{X: cx - size/2, Y: cy - size/2, W: size, H: size}, fill, 1.5)
	})
}
