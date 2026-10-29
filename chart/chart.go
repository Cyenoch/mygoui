// Package chart provides native, coss-themed charts inspired by Bklit's MIT
// chart catalog. Data and interaction state belong to the calling window.
package chart

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

// State is window-owned inspection and legend state. Set Selected to -1 to
// start without a selection. Nil Options.State makes a chart noninteractive.
type State struct {
	Selected int
	Hidden   map[string]bool
	Handle   ui.Handle
}

// Options applies to every chart. Height is unscaled DIPs (default 240).
// Indicator draws a custom selection in plot coordinates; it replaces the
// default indicator, not the chart. Draw callbacks must only paint.
type Options struct {
	Label, Description string
	Height             float32
	State              *State
	Disabled, Loading  bool
	EmptyText          string
	OnSelect           func(int)
	OnActivate         func(int)
	Indicator          func(*ui.Painter, ui.Rect, int)
	hideReadout        bool
}

// SeriesColor resolves a series index through the host's coss chart roles.
func SeriesColor(c *ui.Context, index int) ui.Color {
	t := theme.Colors(c)
	colors := [...]ui.Color{t.Chart1, t.Chart2, t.Chart3, t.Chart4, t.Chart5}
	return colors[((index%len(colors))+len(colors))%len(colors)]
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func number(v float64) string { return strconv.FormatFloat(v, 'g', 5, 64) }

func checkOptions(o Options) error {
	if !finite(float64(o.Height)) || o.Height < 0 {
		return fmt.Errorf("chart %q: height must be finite and nonnegative", o.Label)
	}
	return nil
}

func chartHeight(c *ui.Context, o Options) float32 {
	h := o.Height
	if h == 0 {
		h = 240
	}
	return theme.Scale(c, h)
}

// surface builds one native focus stop and native text for the current value.
// Hit coordinates are fractions from the top left. Geometry is build-local.
func surface(c *ui.Context, o Options, count int, describe func(int) string,
	draw func(*ui.Painter, ui.Rect, int), hit func(float32, float32) int,
) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space8))
	root.Children(func() {
		area := ui.Box(c).FillWidth().MinWidth(0).Height(chartHeight(c, o)).
			Label(o.Label).Description(o.Description).Role(ui.RoleGroup).Disabled(o.Disabled)
		selected := -1
		if o.State != nil {
			area.Bind(&o.State.Handle)
			if !area.IsDisabled() && !o.Loading && count > 0 {
				selected = inspect(c, area, o, count, hit)
			}
		}
		if o.Loading {
			area.Children(func() { mygoui.Skeleton(c, mygoui.SkeletonOptions{}).Fill() })
			return
		}
		if count == 0 {
			text := o.EmptyText
			if text == "" {
				text = "No data"
			}
			area.Center().Children(func() {
				style.Text(ui.Text(c, text), c, theme.SM()).TextColor(t.MutedForeground)
			})
			return
		}
		if selected >= 0 {
			area.Value(describe(selected))
		}
		if o.State != nil && o.OnActivate != nil && !area.IsDisabled() {
			activate := func() {
				if i := o.State.Selected; i >= 0 && i < count {
					o.OnActivate(i)
				}
			}
			area.OnClick(activate).OnShortcut(0, ui.KeyEnter, activate)
		}
		area.FocusRing(area.Focused()).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Clip(r, 0, func() {
				draw(p, r, selected)
				if o.Indicator != nil && selected >= 0 {
					o.Indicator(p, r, selected)
				}
			})
		})
		if selected >= 0 && !o.hideReadout {
			Tooltip(c, describe(selected))
		}
	})
	return root
}

// Keep the readout outside Plot's central grid cell, so axes always match the
// drawing's height rather than growing around a newly selected tooltip.
func plotFrame(c *ui.Context, o Options, axes native.Axes, count int, describe func(int) string,
	build func(ui.Element, Options),
) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8)
	root.Children(func() {
		inner := o
		inner.hideReadout = true
		native.Plot(c, axes, func(area ui.Element) { build(area, inner) }).FillWidth().MinWidth(0).
			Gap(theme.Space12).Padding(theme.Space8, theme.Space16)
		if o.State != nil && !root.IsDisabled() && !o.Disabled && !o.Loading {
			if i := o.State.Selected; i >= 0 && i < count {
				Tooltip(c, describe(i))
			}
		}
	})
	return root
}

type pointerState struct {
	x, y    float32
	pending bool
}

func inspect(c *ui.Context, area ui.Element, o Options, count int, hit func(float32, float32) int) int {
	area.Focusable()
	st := o.State
	old := st.Selected
	if st.Selected >= count || st.Selected < -1 {
		st.Selected = -1
	}
	ptr := ui.Local(area, "chart.pointer", func() pointerState { return pointerState{} })
	// Native event routing observes committed disabled ancestors, including
	// modifiers applied after Children. Polling PointerPosition does not.
	area.HandleInput(func(event ui.InputEvent) bool {
		if event.Kind != ui.InputPointerMove {
			return false
		}
		ptr.x, ptr.y, ptr.pending = event.X, event.Y, true
		return true
	})
	choose := func(index int) {
		if st.Selected == index {
			return
		}
		st.Selected = index
		if o.OnSelect != nil {
			o.OnSelect(index)
		}
	}
	area.OnShortcut(0, ui.KeyRight, func() { choose(min(count-1, st.Selected+1)) })
	area.OnShortcut(0, ui.KeyDown, func() { choose(min(count-1, st.Selected+1)) })
	back := func() {
		index := st.Selected
		if index < 0 {
			index = count
		}
		choose(max(0, index-1))
	}
	area.OnShortcut(0, ui.KeyLeft, back).OnShortcut(0, ui.KeyUp, back).
		OnShortcut(0, ui.KeyHome, func() { choose(0) }).
		OnShortcut(0, ui.KeyEnd, func() { choose(count - 1) }).
		OnShortcut(0, ui.KeyEscape, func() { choose(-1) })
	r := area.Bounds()
	if ptr.pending {
		ptr.pending = false
		if r.W > 0 && r.H > 0 {
			st.Selected = hit(ptr.x/r.W, ptr.y/r.H)
		}
	}
	if st.Selected != old && o.OnSelect != nil {
		o.OnSelect(st.Selected)
	}
	return st.Selected
}

// Tooltip is a native, wrapping selection readout. Unlike a DOM overlay it
// participates in layout and cannot be clipped by a short native viewport.
func Tooltip(c *ui.Context, text string) ui.Element {
	t := theme.Colors(c)
	return style.Text(ui.Text(c, text), c, theme.XS()).TextColor(t.PopoverForeground).
		Background(t.Popover).Padding(theme.Space8, theme.Space12).
		Radius(theme.RadiusScale().MD).Role(ui.RoleStatus)
}

// LegendItem identifies a series. ID must match the series ID to filter it.
type LegendItem struct {
	ID, Label string
	Color     ui.Color
}

// Legend uses native toggles; nil state gives a read-only legend.
func Legend(c *ui.Context, state *State, items []LegendItem, disabled bool) ui.Element {
	return ui.Row(c).FillWidth().Wrap().Gap(theme.Space8).Children(func() {
		for i, item := range items {
			color := item.Color
			if color == (ui.Color{}) {
				color = SeriesColor(c, i)
			}
			visible := state == nil || !state.Hidden[item.ID]
			ui.Row(c.Key(item.ID)).Gap(theme.Space6).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(8, 8).Radius(2).Background(color)
				if state == nil {
					style.Text(ui.Text(c, item.Label), c, theme.XS())
					return
				}
				mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
					Label: item.Label, Selected: &visible, Disabled: disabled,
					Size: theme.ButtonSmall, Variant: mygoui.ToggleGhost,
					OnChange: func(on bool) {
						if state.Hidden == nil {
							state.Hidden = map[string]bool{}
						}
						state.Hidden[item.ID] = !on
						state.Selected = -1
					}})
			})
		}
	})
}

// BrushState is an inclusive data-coordinate viewport, not a copy of samples.
type BrushState struct{ From, To float64 }

// Brush selects an ordered domain through a native two-thumb range slider.
func Brush(c *ui.Context, state *BrushState, min, max float64, disabled bool) (ui.Element, error) {
	if state == nil || !finite(min) || !finite(max) || max <= min {
		return ui.Element{}, fmt.Errorf("chart brush: state and a finite increasing domain are required")
	}
	if !finite(state.From) || !finite(state.To) {
		return ui.Element{}, fmt.Errorf("chart brush: selection must be finite")
	}
	return form.Slider(c, form.SliderOptions{
		Label: "Range", State: form.SliderState{Low: &state.From, High: &state.To},
		Min: min, Max: max, Disabled: disabled, HideOutput: true,
	}), nil
}

func line(p *ui.Painter, a, b ui.Point, width float32, color ui.Color) {
	p.StrokePath((&ui.Path{}).MoveTo(a.X, a.Y).LineTo(b.X, b.Y), width, color)
}

func point(r ui.Rect, x, y float32) ui.Point {
	return ui.Point{X: native.X(r, x), Y: native.Y(r, y)}
}
