package overlay

import (
	"math"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
)

// ModalBackdrop selects how a modal shades the content behind it. The zero
// value dims and blurs the background.
type ModalBackdrop uint8

const (
	// ModalBackdropBlur dims and blurs what is behind the overlay with a
	// 4 DIP Gaussian standard deviation.
	ModalBackdropBlur ModalBackdrop = iota
	// ModalBackdropOpaque dims what is behind the overlay without blur.
	ModalBackdropOpaque
	// ModalBackdropTransparent leaves what is behind visible and undimmed.
	ModalBackdropTransparent
)

// backdropBlur is COSS's backdrop-blur-sm, in DIPs of Gaussian sigma.
const backdropBlur = 4

// paintBackdrop applies a backdrop variant to back. The caller keeps
// the layout (centering, padding, orientation) so Modal, Drawer, Dialog and
// BottomSheet share one paint contract.
func paintBackdrop(c *ui.Context, back ui.Element, variant ModalBackdrop) ui.Element {
	t := theme.Colors(c)
	switch variant {
	case ModalBackdropTransparent:
		return back
	case ModalBackdropOpaque:
		return back.Background(t.Backdrop)
	default:
		// Material replaces the background, so the dim is painted over the
		// blur with the element's paint callback.
		back = back.Material(glass.Blur{Radius: backdropBlur})
		return back.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.Backdrop, 0)
		})
	}
}

// overlayRadius is the COSS dialog radius.
func overlayRadius() float32 {
	return theme.RadiusScale().XXL
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
	fullBleedFooter      bool
	centerHeader         bool

	titleStyle       theme.TextStyle
	descriptionStyle theme.TextStyle
	padding          float32
	headerGap        float32
	headerBody       float32
	bodyFooter       float32
	closeInset       float32
	headingLevel     int
}

// fillPanel builds p's sections inside the already styled panel element.
// Header and body exist only when they have content, so the CSS sibling
// spacing applies without empty gaps.
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
		if p.fullBleedFooter {
			bodyFooter = padding
		}
	}
	closeInset := p.closeInset
	if closeInset == 0 {
		closeInset = theme.Space8
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
			if p.centerHeader {
				h.TextAlign(ui.Center)
			}
			if hasBody {
				h.Margin(0, 0, headerBody, 0)
			} else if hasFooter {
				gap := bodyFooter
				if w, _ := c.Size(); p.fullBleedFooter && w < 640 {
					gap = theme.Space16
				}
				h.Margin(0, 0, gap, 0)
			}
			if hasClose {
				// Reserve the part of the close button that reaches the header.
				if overlap := closeInset + theme.Scale(c, closeButtonSize) - padding; overlap > 0 {
					h.PaddingEnd(overlap)
				}
			}
			h.Children(func() {
				if p.header != nil {
					p.header()
					return
				}
				if p.icon != nil {
					p.icon()
				}
				if p.title != "" {
					title := style.Text(ui.Text(c, p.title), c, p.titleStyle)
					title.FontWeight(theme.WeightBold).TextColor(t.PopoverForeground).
						Role(ui.RoleHeading).Level(level)
				}
				if p.description != "" {
					style.Text(ui.Text(c, p.description), c, p.descriptionStyle).TextColor(t.MutedForeground)
				}
			})
		}
		if hasBody {
			// Leave room for a child's focus ring inside the scroll clip.
			b := ui.Scroll(c).FillWidth().Grow(1).MinHeight(0).Padding(bodyBleed).
				Margin(0, -bodyBleed, 0, -bodyBleed)
			if hasFooter {
				b.Margin(0, -bodyBleed, bodyFooter, -bodyBleed)
			}
			if p.bodyLabel != "" {
				b.Label(p.bodyLabel)
			}
			b.Children(p.body)
		}
		if hasFooter {
			f := ui.Row(c).Justify(ui.End).AlignItems(ui.Center).Gap(theme.Space8).Shrink(0)
			if p.fullBleedFooter {
				f.AlignSelf(ui.Stretch).Margin(ui.Auto, -padding, -padding, -padding).
					Padding(theme.Space16, theme.Space24).
					Background(t.MutedBackground.Alpha(0.72)).
					BorderWidth(theme.Hairline, 0, 0, 0).BorderColor(t.Border)
				if w, _ := c.Size(); w < 640 {
					f.Column().Reverse().AlignItems(ui.Stretch)
				}
			} else {
				f.FillWidth().Background(t.MutedBackground).Radius(theme.RadiusScale().MD).Padding(theme.Space8)
			}
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

// dialogViewport positions the default dialogs in a three-row desktop grid,
// or against the bottom on narrow windows.
func dialogViewport(back, panel ui.Element, w, h float32) (float32, float32) {
	back.Grid().JustifyItems(ui.Center).AlignItems(ui.Stretch)
	panel.RowStart(2)
	if w < 640 {
		back.RowTracks(ui.Fr(1), ui.FitContent()).Padding(theme.Space48, 0, 0, 0)
		return max(w, 0), max(h-theme.Space48, 0)
	}
	back.RowTracks(ui.Fr(1), ui.FitContent(), ui.Fr(3)).Padding(theme.Space16)
	return max(w-2*theme.Space16, 0), max(h-2*theme.Space16, 0)
}

// softPair returns a soft surface and its foreground for a status.
func softPair(c *ui.Context, status mygoui.StatusColor) (bg, fg ui.Color) {
	t := theme.Colors(c)
	if status == mygoui.ColorDefault {
		return t.Default, t.Foreground
	}
	p := mygoui.StatusColors(c, status)
	return p.Soft, p.SoftForeground
}

// statusIcon builds a 40 DIP status medallion with a 20 DIP glyph.
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

// The optional overlay arrow is 12 DIPs wide and 6 DIPs tall, its base on the
// panel's edge facing the anchor, as OverlayArrow places the 12×12 art whose
// filled bezier (M0 0C5.48 8 6.5 8 12 0Z) is 12 wide and 6 tall.
const (
	arrowWidth  = 12
	arrowHeight = 6

	// closeButtonSize is mygoui.CloseButton's box, which the header reserves
	// room for.
	closeButtonSize = theme.Space24

	// bodyBleed keeps focus rings inside the scroll viewport.
	bodyBleed = 3
)

// paintArrow draws the arrow outside the panel's edge facing the
// anchor, along the edge at along of its length, where arrowAlong puts the
// anchor's center. The panel must not Clip, or the protruding part disappears.
//
// The painter fills rectangles, not paths: the bezier's outline is drawn as a
// square turned 45° with rounded corners, clipped to the arrow's half of it,
// which leaves the same 12 DIP base and 6 DIP rounded tip.
func paintArrow(p *ui.Painter, r ui.Rect, side mygoui.Side, fill ui.Color, along float32) {
	half := float32(arrowWidth) / 2
	var band ui.Rect
	var cx, cy float32
	switch side {
	case mygoui.SideTop:
		cx, cy = r.X+r.W*along, r.Y+r.H
		band = ui.Rect{X: cx - half, Y: cy, W: arrowWidth, H: arrowHeight}
	case mygoui.SideLeft:
		cx, cy = r.X+r.W, r.Y+r.H*along
		band = ui.Rect{X: cx, Y: cy - half, W: arrowHeight, H: arrowWidth}
	case mygoui.SideRight:
		cx, cy = r.X, r.Y+r.H*along
		band = ui.Rect{X: cx - arrowHeight, Y: cy - half, W: arrowHeight, H: arrowWidth}
	default:
		cx, cy = r.X+r.W*along, r.Y
		band = ui.Rect{X: cx - half, Y: cy - arrowHeight, W: arrowWidth, H: arrowHeight}
	}
	// The square's diagonal is the arrow's width, so its apex is the tip.
	side2 := half * 0.70711
	p.Clip(band, 0, func() {
		p.Transform(ui.Rotation(45), cx, cy, func() {
			p.Fill(ui.Rect{X: cx - side2, Y: cy - side2, W: 2 * side2, H: 2 * side2}, fill, 1.6)
		})
	})
}

// arrowAlong returns where along the edge the arrow points, as a fraction of
// its length: the anchor's center, as React Aria holds OverlayArrow, kept half
// an arrow plus a DIP inside either end. The layout boxes of the frame before
// are in the same coordinates as the rect a Draw callback paints.
func arrowAlong(panel, anchor ui.Element, side mygoui.Side) float32 {
	const middle = 0.5
	pr, ar := panel.LayoutBounds(), anchor.LayoutBounds()
	vertical := side == mygoui.SideTop || side == mygoui.SideBottom
	size, center, start := pr.W, ar.X+ar.W/2, pr.X
	if !vertical {
		size, center, start = pr.H, ar.Y+ar.H/2, pr.Y
	}
	if size <= 0 || ar.W <= 0 || ar.H <= 0 {
		return middle
	}
	margin := (arrowWidth/2 + 1) / size
	if margin >= middle {
		return middle
	}
	return min(1-margin, max(margin, (center-start)/size))
}

// arrowSide returns the side the panel settled on, from the boxes the frame
// before laid out: the attach engine flips a panel that would overflow the
// window, so the requested side is only the estimate for the frame a panel
// first appears in. The side is the one whose edges are closest, which holds
// for any alignment along the other axis.
func arrowSide(panel, anchor ui.Element, want mygoui.Side) mygoui.Side {
	pr, ar := panel.LayoutBounds(), anchor.LayoutBounds()
	if pr.W <= 0 || pr.H <= 0 || ar.W <= 0 || ar.H <= 0 {
		return want
	}
	best, gap := want, float32(math.Inf(1))
	consider := func(side mygoui.Side, d float32) {
		if d = float32(math.Abs(float64(d))); d < gap {
			best, gap = side, d
		}
	}
	consider(mygoui.SideBottom, pr.Y-(ar.Y+ar.H))
	consider(mygoui.SideTop, ar.Y-(pr.Y+pr.H))
	consider(mygoui.SideRight, pr.X-(ar.X+ar.W))
	consider(mygoui.SideLeft, ar.X-(pr.X+pr.W))
	return best
}
