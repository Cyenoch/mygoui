package form

import (
	"slices"

	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ImageSelectItem is one tile. Content, when set, builds the tile's picture
// in place of Image; Label names the tile and shows under it.
type ImageSelectItem struct {
	Value    string
	Label    string
	Image    ui.ImageSource
	Content  func()
	Disabled bool
}

// ImageSelectOptions configures a row of picture tiles. Value selects one
// tile like a radio group; Values, when set instead, selects any number like
// check boxes. With neither the tiles are disabled. Width and Height size a
// tile's picture and default to 144 by 90 DIP.
type ImageSelectOptions struct {
	AriaLabel string
	Items     []ImageSelectItem
	Value     *string
	Values    *[]string
	Width     float32
	Height    float32
	Disabled  bool
	OnChange  func()
}

// ImageSelect chooses among picture tiles. A single choice is one Tab stop
// whose arrows move the choice; multiple tiles each toggle with Space. The
// chosen tiles show a primary outline and a check badge.
func ImageSelect(c *ui.Context, o ImageSelectOptions) ui.Element {
	disabled := o.Disabled || (o.Value == nil && o.Values == nil)
	if o.Width <= 0 {
		o.Width = 144
	}
	if o.Height <= 0 {
		o.Height = 90
	}
	var group ui.Element
	if o.Values != nil {
		group = ui.Row(c).Role(ui.RoleGroup)
		group.Children(func() {
			for _, it := range o.Items {
				on := slices.Contains(*o.Values, it.Value)
				tile := ui.CheckboxBase(c.Key(it.Value), &on)
				tile.OnChange(func() {
					*o.Values = imageSelectToggle(*o.Values, it.Value)
					if o.OnChange != nil {
						o.OnChange()
					}
				})
				imageSelectTile(c, tile, it, on, o)
			}
		})
	} else {
		selected := o.Value
		if selected == nil {
			selected = new(string)
		}
		group = ui.RadioGroup(c, func() {
			for _, it := range o.Items {
				tile := ui.RadioBase(c.Key(it.Value), selected, it.Value)
				if o.OnChange != nil {
					tile.OnChange(o.OnChange)
				}
				imageSelectTile(c, tile, it, *selected == it.Value, o)
			}
		})
	}
	group.Row().Wrap().MinWidth(0).Gap(theme.Space(3)).Label(o.AriaLabel)
	return style.Disabled(group, disabled)
}

// imageSelectToggle returns values with value added or removed, leaving the
// caller's previous slice untouched.
func imageSelectToggle(values []string, value string) []string {
	if !slices.Contains(values, value) {
		return append(slices.Clone(values), value)
	}
	return slices.DeleteFunc(slices.Clone(values), func(v string) bool { return v == value })
}

// imageSelectTile styles a radio or check box base as a picture with its
// label underneath.
func imageSelectTile(c *ui.Context, tile ui.Element, it ImageSelectItem, on bool, o ImageSelectOptions) {
	t := theme.Colors(c)
	radius := theme.RadiusScale().LG
	tile = style.Disabled(tile.Column().AlignItems(ui.Center).Gap(theme.Space(2)).Shrink(0).
		FocusRing(false).Cursor(ui.CursorPointer).Label(it.Label), it.Disabled)
	hovered, focused := tile.Hovered(), tile.FocusVisible()
	ink := t.MutedForeground
	if on || hovered {
		ink = t.Foreground
	}
	tile.Children(func() {
		frame := ui.Box(c).Size(theme.Scale(c, o.Width), theme.Scale(c, o.Height)).Shrink(0).
			Radius(radius).Clip().Background(t.MutedBackground)
		frame.Children(func() {
			switch {
			case it.Content != nil:
				it.Content()
			case it.Image != nil:
				ui.Image(c, it.Image).Fit(ui.Cover).Fill()
			}
		})
		badge := theme.Scale(c, 16)
		rtl := c.Direction() == ui.RightToLeft
		// The ring fades and the check draws in as the tile is chosen; both
		// leave by fading.
		chosen := motion.Progress(frame, "chosen", on)
		frame.DrawOver(func(p *ui.Painter, r ui.Rect) {
			switch {
			case hovered:
				p.Stroke(r, t.Ring, radius, theme.Hairline)
			default:
				p.Stroke(r, t.Border, radius, theme.Hairline)
			}
			if chosen > 0 {
				p.Stroke(r, t.Primary.Alpha(chosen), radius, 2)
			}
			if focused {
				formstyle.FieldRing(p, r, radius, t.Ring.Alpha(0.24))
			}
			if chosen == 0 {
				return
			}
			dot := ui.Rect{X: r.X + r.W - badge - theme.Space(1.5), Y: r.Y + theme.Space(1.5), W: badge, H: badge}
			if rtl {
				dot.X = r.X + theme.Space(1.5)
			}
			drawn := float32(1)
			if on {
				drawn = chosen
			}
			p.Fill(dot, t.Primary.Alpha(chosen), badge/2)
			mark := badge * 0.6
			p.StrokePath(choiceCheckPath(ui.Rect{
				X: dot.X + (badge-mark)/2, Y: dot.Y + (badge-mark)/2, W: mark, H: mark,
			}).Trim(0, drawn), theme.Scale(c, 1.6), t.PrimaryForeground.Alpha(chosen))
		})
		if it.Label != "" {
			style.Text(ui.Text(c, it.Label), c, theme.SM()).TextColor(ink).SingleLine()
		}
	})
}
