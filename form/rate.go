package form

import (
	"math"
	"strconv"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// RateOptions configures a star rating. Count defaults to five stars.
// AllowHalf steps by half a star; AllowClear lets a click on the current
// rating clear it. ReadOnly shows the rating without taking input.
type RateOptions struct {
	AriaLabel  string
	Count      int
	AllowHalf  bool
	AllowClear bool
	ReadOnly   bool
	Disabled   bool
	OnChange   func()
}

const (
	rateStarSize = 20
	rateStarGap  = 4
)

// Rate is a star rating of *value out of Count. The pointer previews and a
// click sets the rating; the arrows step it, Page Up and Page Down move it a
// whole star, and Home and End take it to none and all.
func Rate(c *ui.Context, value *float64, o RateOptions) ui.Element {
	t := theme.Colors(c)
	count := o.Count
	if count <= 0 {
		count = 5
	}
	step := 1.0
	if o.AllowHalf {
		step = 0.5
	}
	disabled := o.Disabled || value == nil
	if value == nil {
		var none float64
		value = &none
	}
	rtl := c.Direction() == ui.RightToLeft
	size, gap := theme.Scale(c, rateStarSize), theme.Scale(c, rateStarGap)
	width := float32(count)*(size+gap) - gap

	var root ui.Element
	if o.ReadOnly {
		root = ui.Box(c)
	} else {
		root = ui.ButtonBase(c).FocusRing(false).Cursor(ui.CursorPointer)
	}
	root = style.Disabled(root.Size(width, size).Shrink(0).Role(ui.RoleSlider).Label(o.AriaLabel), disabled)

	shown, previewing := *value, false
	if !o.ReadOnly && !disabled {
		set := func(v float64) {
			v = math.Min(math.Max(v, 0), float64(count))
			if v == *value {
				return
			}
			*value = v
			if o.OnChange != nil {
				o.OnChange()
			}
		}
		x, _, over := root.PointerPosition()
		if rtl {
			x = width - x
		}
		pointed := rateAt(x, size, gap, count, o.AllowHalf)
		if over {
			shown, previewing = pointed, true
		}
		root.OnClick(func() {
			switch {
			case !over: // Enter and Space click too, without a star to choose.
			case o.AllowClear && pointed == *value:
				set(0)
			default:
				set(pointed)
			}
		})
		forward, back := ui.KeyRight, ui.KeyLeft
		if rtl {
			forward, back = back, forward
		}
		// Snap a value off the grid before stepping from it.
		at := math.Round(*value/step) * step
		for _, k := range [...]struct {
			key ui.Key
			to  float64
		}{
			{forward, at + step}, {ui.KeyUp, at + step},
			{back, at - step}, {ui.KeyDown, at - step},
			{ui.KeyPageUp, at + 1}, {ui.KeyPageDown, at - 1},
			{ui.KeyHome, 0}, {ui.KeyEnd, float64(count)},
		} {
			root.OnShortcut(0, k.key, func() { set(k.to) })
		}
		// Enter and Space would click the star under a resting pointer.
		root.OnShortcut(0, ui.KeyEnter, func() {})
		root.OnShortcut(0, ui.KeySpace, func() {})
	}
	shown = math.Min(math.Max(shown, 0), float64(count))
	root.Range(0, float64(count), *value).Step(step).
		Value(strconv.FormatFloat(*value, 'f', -1, 64) + " / " + strconv.Itoa(count))

	focused := root.FocusVisible()
	// Stars fill as the rating is set; the pointer's preview follows it at once.
	stars := ui.Local(root, rateStarsKey{}, func() rateStars { return rateStars{} })
	fills := stars.fills[:0]
	for i := range count {
		fill := float32(math.Min(math.Max(shown-float64(i), 0), 1))
		if previewing {
			fill = root.Jump(rateStarKey(i), fill)
		} else {
			fill = root.AnimateTo(rateStarKey(i), fill, ui.Animation{Duration: theme.Motion().State, Ease: theme.Motion().Ease})
		}
		fills = append(fills, fill)
	}
	stars.fills = fills
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		paths := stars.at(r, count, size, gap, rtl)
		for i := range count {
			box := paths[i].box
			star := paths[i].path
			fill := fills[i]
			if fill < 1 {
				rateFillStar(p, star, t.Input)
			}
			if fill > 0 {
				part := box
				part.W = size * fill
				if rtl {
					part.X = box.X + size - part.W
				}
				p.Clip(part, 0, func() { rateFillStar(p, star, t.Warning) })
			}
		}
		if focused {
			p.FocusRing(r, choiceRadii(theme.RadiusScale().SM))
		}
	})
	return root
}

// rateStarKey names the animated fill of one star.
type rateStarKey int

// rateStarsKey keeps a rating's star fills and paths across frames, which
// paint the same stars until the rating moves or resizes.
type rateStarsKey struct{}

type rateStars struct {
	fills []float32
	r     ui.Rect
	rtl   bool
	paths []rateStar
}

type rateStar struct {
	box  ui.Rect
	path *ui.Path
}

// at returns the stars laid out in r, building them again only when the
// layout changed.
func (s *rateStars) at(r ui.Rect, count int, size, gap float32, rtl bool) []rateStar {
	if r == s.r && rtl == s.rtl && len(s.paths) == count {
		return s.paths
	}
	s.r, s.rtl, s.paths = r, rtl, s.paths[:0]
	for i := range count {
		box := ui.Rect{X: r.X + float32(i)*(size+gap), Y: r.Y, W: size, H: size}
		if rtl {
			box.X = r.X + r.W - size - float32(i)*(size+gap)
		}
		s.paths = append(s.paths, rateStar{box, rateStarPath(box)})
	}
	return s.paths
}

// rateAt maps a pointer offset from the first star's leading edge to the
// rating it points at.
func rateAt(x, size, gap float32, count int, half bool) float64 {
	i := min(max(int(x/(size+gap)), 0), count-1)
	if half && x-float32(i)*(size+gap) < size/2 {
		return float64(i) + 0.5
	}
	return float64(i + 1)
}

// rateFillStar fills a star and strokes it for rounded points.
func rateFillStar(p *ui.Painter, star *ui.Path, ink ui.Color) {
	p.FillPath(star, ink)
	p.StrokePath(star, 1.5, ink)
}

// rateStarPath is a five-pointed star inset in box for its stroke.
func rateStarPath(box ui.Rect) *ui.Path {
	const inner = 0.47
	cx, cy := box.X+box.W/2, box.Y+box.H/2+box.H*0.04
	outer := box.W/2 - 1.5
	var p ui.Path
	for i := range 10 {
		radius := outer
		if i%2 == 1 {
			radius = outer * inner
		}
		angle := -math.Pi/2 + float64(i)*math.Pi/5
		x, y := cx+radius*float32(math.Cos(angle)), cy+radius*float32(math.Sin(angle))
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	return p.Close()
}
