package layout

import (
	"math"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type CarouselState struct {
	Index         int
	offset        ui.Point
	lastIndex     int
	extent        float32
	count         int
	vertical, rtl bool
}

type CarouselOptions struct {
	Count         int
	Width, Height float32
	Vertical      bool
	Label         string
	Disabled      bool
	OnChange      func(int)
}

// Carousel places actual pages in a clipped strip. Native Draggable projects
// release velocity and springs to a page; catching a moving strip cancels its
// flight. It owns no timer or scroll state, so wheel input remains available to
// nested scroll containers. Interactive children retain native focus behavior.
func Carousel(c *ui.Context, state *CarouselState, o CarouselOptions, page func(int)) ui.Element {
	width, height := o.Width, o.Height
	if width <= 0 {
		width = 300
	}
	if height <= 0 {
		height = 180
	}
	count := max(0, o.Count)
	state.Index = min(max(0, state.Index), max(0, count-1))
	extent := width
	axis := ui.Horizontal
	rtl := c.Direction() == ui.RightToLeft
	sign := float32(-1)
	if rtl {
		sign = 1
	}
	if o.Vertical {
		extent, axis, sign = height, ui.Vertical, -1
	}
	point := func(index int) ui.Point {
		p := ui.Point{X: sign * float32(index) * extent}
		if o.Vertical {
			p.X, p.Y = 0, p.X
		}
		return p
	}
	if state.Index != state.lastIndex || state.extent != extent || state.count != count || state.vertical != o.Vertical || state.rtl != rtl {
		state.offset = point(state.Index)
		state.lastIndex, state.extent, state.count = state.Index, extent, count
		state.vertical, state.rtl = o.Vertical, rtl
	}
	if o.Disabled {
		state.offset = point(state.Index)
	}
	selectPage := func(index int) {
		index = min(max(0, index), max(0, count-1))
		if index == state.Index || o.Disabled {
			return
		}
		state.Index, state.lastIndex = index, index
		state.offset = point(index)
		if o.OnChange != nil {
			o.OnChange(index)
		}
	}
	previous, next := lucide.ArrowLeft(), lucide.ArrowRight()
	backKey, nextKey := ui.KeyLeft, ui.KeyRight
	if rtl {
		previous, next, backKey, nextKey = next, previous, nextKey, backKey
	}
	if o.Vertical {
		previous, next, backKey, nextKey = lucide.ArrowUp(), lucide.ArrowDown(), ui.KeyUp, ui.KeyDown
	}
	direction := c.Direction()
	return ui.Column(c).MinWidth(0).Gap(theme.Space8).Label(o.Label).Disabled(o.Disabled).Children(func() {
		viewport := ui.Box(c.Key("viewport")).Size(width, height).Clip().Focusable().Label(o.Label).
			OnShortcut(0, backKey, func() { selectPage(state.Index - 1) }).
			OnShortcut(0, nextKey, func() { selectPage(state.Index + 1) }).
			OnShortcut(0, ui.KeyHome, func() { selectPage(0) }).
			OnShortcut(0, ui.KeyEnd, func() { selectPage(count - 1) })
		viewport.Children(func() {
			var strip ui.Element
			if o.Vertical {
				strip = ui.Column(c.Key("pages")).Width(width).Height(height * float32(count))
			} else {
				strip = ui.Row(c.Key("pages")).Width(width * float32(count)).Height(height)
			}
			strip.Absolute().Top(0).Shrink(0)
			if !o.Vertical && rtl {
				strip.Right(0)
			} else {
				strip.Left(0)
			}
			bound := point(max(0, count-1))
			bounds := ui.Rect{X: min(0, bound.X), Y: min(0, bound.Y), W: float32(math.Abs(float64(bound.X))), H: float32(math.Abs(float64(bound.Y)))}
			if !o.Disabled && count > 1 {
				strip.Draggable(&state.offset, ui.DraggableOptions{
					Axis: axis, Bounds: &bounds, Animation: ui.Smooth,
					Snap: func(projected ui.Point) ui.Point {
						value := projected.X
						if o.Vertical {
							value = projected.Y
						}
						index := min(max(0, int(math.Round(float64(value/(sign*extent))))), count-1)
						selectPage(index)
						return point(index)
					},
				})
			} else {
				strip.Translate(state.offset.X, state.offset.Y)
			}
			strip.Children(func() {
				for i := range count {
					ui.Box(c.Key(i)).Direction(direction).Size(width, height).Shrink(0).Clip().Disabled(o.Disabled || i != state.Index).
						Label(locales.ComponentSlideOf(i+1, count)).Children(func() {
						if page != nil {
							page(i)
						}
					})
				}
			})
		})
		ui.Row(c).Gap(theme.Space8).Children(func() {
			mygoui.Button(c.Key("previous"), mygoui.ButtonOptions{Label: locales.ComponentPreviousSlide(), Icon: previous, Size: mygoui.ButtonIconSM, Variant: mygoui.ButtonOutline, Disabled: o.Disabled || state.Index <= 0 || count == 0, OnClick: func() { selectPage(state.Index - 1) }})
			label := locales.ComponentPositionOf(state.Index+1, count)
			if count == 0 {
				label = locales.ComponentNoSlides()
			}
			ui.Text(c, label).FontSize(theme.Scale(c, theme.Small().Size))
			mygoui.Button(c.Key("next"), mygoui.ButtonOptions{Label: locales.ComponentNextSlide(), Icon: next, Size: mygoui.ButtonIconSM, Variant: mygoui.ButtonOutline, Disabled: o.Disabled || state.Index >= count-1, OnClick: func() { selectPage(state.Index + 1) }})
		})
	})
}
