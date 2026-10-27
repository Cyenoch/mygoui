package form

import (
	"math"
	"strconv"
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SliderMark is an optional labelled tick on a horizontal slider.
type SliderMark struct {
	Value float64
	Label string
}

// SliderState is the caller-owned value of a slider. Value drives a
// single-thumb slider; Low and High together drive a two-thumb range. Any
// nil pointer renders that slider disabled.
type SliderState struct {
	Value *float64
	Low   *float64
	High  *float64
}

// IsRange reports whether Low and High are both present.
func (s SliderState) IsRange() bool { return s.Low != nil && s.High != nil }

// SliderOptions configures a slider. A zero Min and Max give the
// upstream 0..100 range; any other empty range becomes Min..Min+1. Step snaps
// drags and keyboard moves; a non-positive Step leaves the native
// hundredth-of-the-range keyboard step.
//
// Marks draw labelled ticks below a horizontal track. Vertical sliders support
// only a single thumb: ui.RangeSlider is horizontal, so a Vertical two-thumb
// state keeps its horizontal layout. Height sets a vertical track's length.
type SliderOptions struct {
	State       SliderState
	Label       string
	Description string
	Error       string
	Min         float64
	Max         float64
	Step        float64
	Marks       []SliderMark
	Vertical    bool
	Height      float32
	Disabled    bool
	Invalid     bool
	HideOutput  bool
	Format      func(float64) string
	OnChange    func()
}

// Slider is a native slider with a primary-color track and round thumb.
// Drags and the arrows, Home and End change the caller's value; a two-thumb
// state is ordered by the native range slider, so the thumbs never cross.
func Slider(c *ui.Context, o SliderOptions) ui.Element {
	min, max := o.Min, o.Max
	if min == 0 && max == 0 {
		max = 100
	}
	if max <= min {
		max = min + 1
	}
	step := o.Step
	if step < 0 {
		step = 0
	}
	state := o.State
	disabled := o.Disabled
	var fallback float64
	if !state.IsRange() {
		if state.Value == nil {
			state.Value = state.Low
		}
		if state.Value == nil {
			state.Value = &fallback
			disabled = true
		}
	}
	if state.IsRange() {
		*state.Low = math.Max(min, math.Min(*state.Low, *state.High))
		*state.High = math.Min(max, math.Max(*state.Low, *state.High))
	}
	errText := choiceErrorText(o.Invalid, o.Error)
	if o.Format == nil {
		o.Format = sliderNumberFormat(step)
	}
	t := theme.Colors(c)
	rtl := c.Direction() == ui.RightToLeft

	root := ui.Column(c).MinWidth(0).FillWidth().Gap(theme.Space4)
	if o.Description != "" {
		root.Description(o.Description)
	}
	if errText != "" {
		root.Error(errText)
	}
	if o.Vertical && !state.IsRange() {
		root.Gap(theme.Space8)
	}

	root.Children(func() {
		if o.Vertical && !state.IsRange() {
			sliderVertical(c, o, state, min, max, step, disabled, errText, t)
			return
		}
		if !o.HideOutput || o.Label != "" {
			ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
				if o.Label != "" {
					style.Text(ui.Text(c, o.Label), c, theme.Label()).TextColor(t.Foreground)
				}
				if !o.HideOutput {
					ui.Spacer(c)
					sliderOutput(c, o, state, t)
				}
			})
		}
		if state.IsRange() {
			sliderRange(c, o, state, min, max, step, disabled, errText, t, rtl)
		} else {
			sliderSingle(c, o, state, min, max, step, disabled, errText, t, rtl)
		}
		if len(o.Marks) > 0 {
			sliderMarks(c, o, min, max, rtl, sliderEdgePad(c, state.IsRange()))
		}
		if o.Description != "" {
			Description(c, o.Description)
		}
		if errText != "" {
			FieldError(c, errText)
		}
	})
	return style.Disabled(root, disabled)
}

// sliderSingle builds a one-thumb horizontal track.
func sliderSingle(c *ui.Context, o SliderOptions, state SliderState, min, max, step float64, disabled bool, errText string, t theme.Palette, rtl bool) {
	pad := sliderEdgePad(c, false)
	trackH := theme.Scale(c, 20)
	knobW := theme.Scale(c, 16)
	innerW, innerH := knobW, knobW
	radius := knobW / 2
	innerRadius := radius

	el := ui.SliderBase(c, state.Value, min, max).Height(trackH).PaddingX(pad).
		FocusRing(false).Label(o.Label)
	if step > 0 {
		el.Step(step)
	}
	if o.Description != "" {
		el.Description(o.Description)
	}
	if errText != "" {
		el.Error(errText)
	}
	if o.OnChange != nil && !disabled {
		el.OnChange(o.OnChange)
	}
	// Snapshot the slider's states in the build: a retained copy of a
	// leaving element paints in later passes, where el has expired.
	pressed, focused := el.Pressed(), el.FocusVisible()
	scale := sliderThumbScale(el, pressed)
	el.Draw(func(p *ui.Painter, r ui.Rect) {
		f := sliderFraction(*state.Value, min, max)
		x := r.X + pad + (r.W-2*pad)*f
		if rtl {
			x = r.X + r.W - pad - (r.W-2*pad)*f
		}
		sliderPaintTrack(p, r, t, x, rtl, radius, trackH)
		thumb := ui.Rect{X: x - knobW/2, Y: r.Y + (r.H-knobW)/2, W: knobW, H: knobW}
		sliderPaintThumb(p, t, thumb, innerW, innerH, innerRadius, scale, focused, radius)
	})
}

// sliderRange builds a two-thumb track over the native range slider, whose
// knobs carry the keyboard and drag input.
func sliderRange(c *ui.Context, o SliderOptions, state SliderState, min, max, step float64, disabled bool, errText string, t theme.Palette, rtl bool) {
	pad := sliderEdgePad(c, true)
	trackH := theme.Scale(c, 20)
	knobW := theme.Scale(c, 16)
	innerW, innerH := knobW, knobW
	radius := knobW / 2
	innerRadius := radius

	el := ui.RangeSlider(c, state.Low, state.High, min, max, step).Label(o.Label)
	if o.Description != "" {
		el.Description(o.Description)
	}
	if errText != "" {
		el.Error(errText)
	}
	if o.OnChange != nil && !disabled {
		el.OnChange(o.OnChange)
	}
	focusWithin := el.FocusWithin()
	scale := sliderThumbScale(el, el.Pressed())
	el.DrawOver(func(p *ui.Painter, r ui.Rect) {
		f0, f1 := sliderFraction(*state.Low, min, max), sliderFraction(*state.High, min, max)
		x0 := r.X + pad + (r.W-2*pad)*f0
		x1 := r.X + pad + (r.W-2*pad)*f1
		if rtl {
			x0, x1 = 2*r.X+r.W-x0, 2*r.X+r.W-x1
		}
		if x1 < x0 {
			x0, x1 = x1, x0
		}
		track := ui.Rect{X: r.X, Y: r.Y + (trackH-theme.Scale(c, 4))/2, W: r.W, H: theme.Scale(c, 4)}
		p.Fill(track, t.Input, track.H/2)
		p.Fill(ui.Rect{X: x0, Y: track.Y, W: x1 - x0, H: track.H}, t.Primary, track.H/2)
		for _, x := range [2]float32{x0, x1} {
			thumb := ui.Rect{X: x - knobW/2, Y: r.Y + (trackH-knobW)/2, W: knobW, H: knobW}
			sliderPaintThumb(p, t, thumb, innerW, innerH, innerRadius, scale, focusWithin, radius)
		}
	})
}

// sliderVertical builds the vertical single-thumb slider: output above,
// a 20 DIP track, and the label below.
func sliderVertical(c *ui.Context, o SliderOptions, state SliderState, min, max, step float64, disabled bool, errText string, t theme.Palette) {
	pad := theme.Scale(c, 14)
	trackW := theme.Scale(c, 20)
	knobH := theme.Scale(c, 16)
	innerW, innerH := knobH, knobH
	height := o.Height
	if height <= 0 {
		height = theme.Scale(c, 160)
	}
	radius := knobH / 2
	innerRadius := radius

	if !o.HideOutput {
		ui.Row(c).FillWidth().Center().Children(func() {
			sliderOutput(c, o, state, t)
		})
	}
	ui.Row(c).FillWidth().Children(func() {
		ui.Spacer(c)
		el := ui.SliderBase(c, state.Value, min, max).Width(trackW).Height(height).
			PaddingY(pad).Vertical().FocusRing(false).Label(o.Label)
		if step > 0 {
			el.Step(step)
		}
		if o.Description != "" {
			el.Description(o.Description)
		}
		if errText != "" {
			el.Error(errText)
		}
		if o.OnChange != nil && !disabled {
			el.OnChange(o.OnChange)
		}
		pressed, focused := el.Pressed(), el.FocusVisible()
		scale := sliderThumbScale(el, pressed)
		el.Draw(func(p *ui.Painter, r ui.Rect) {
			f := sliderFraction(*state.Value, min, max)
			x := r.X + r.W/2
			cy := r.Y + r.H - pad - (r.H-2*pad)*f
			track := ui.Rect{X: x - theme.Scale(c, 2), Y: r.Y, W: theme.Scale(c, 4), H: r.H}
			p.Fill(track, t.Input, track.W/2)
			p.Fill(ui.Rect{X: track.X, Y: cy, W: track.W, H: track.Y + track.H - cy}, t.Primary, track.W/2)
			thumb := ui.Rect{X: x - knobH/2, Y: cy - knobH/2, W: knobH, H: knobH}
			sliderPaintThumb(p, t, thumb, innerW, innerH, innerRadius, scale, focused, radius)
		})
		ui.Spacer(c)
	})
	if o.Label != "" {
		ui.Row(c).FillWidth().Center().Children(func() {
			style.Text(ui.Text(c, o.Label), c, theme.Label()).TextColor(t.Foreground)
		})
	}
	if o.Description != "" {
		ui.Row(c).FillWidth().Center().Children(func() {
			Description(c, o.Description)
		})
	}
	if errText != "" {
		FieldError(c, errText)
	}
}

// sliderEdgePad is the horizontal room a horizontal track leaves for its
// thumb at either end; marks use it so they line up with the thumbs.
func sliderEdgePad(c *ui.Context, twoThumbs bool) float32 {
	if twoThumbs {
		return float32(theme.Space20) / 2
	}
	return theme.Scale(c, 14)
}

// sliderPaintTrack paints Coss's four-DIP rail and selected interval.
func sliderPaintTrack(p *ui.Painter, r ui.Rect, t theme.Palette, x float32, rtl bool, radius, trackH float32) {
	track := ui.Rect{X: r.X, Y: r.Y + (r.H-4)/2, W: r.W, H: 4}
	p.Fill(track, t.Input, 2)
	fill := ui.Rect{X: track.X, Y: track.Y, H: track.H}
	if rtl {
		fill.W = track.X + track.W - x
		fill.X = x
	} else {
		fill.W = x - track.X
	}
	p.Fill(fill, t.Primary, 2)
}

func sliderThumbScale(el ui.Element, pressed bool) float32 {
	scale := float32(1)
	if pressed {
		scale = 1.2
	}
	return el.Animate("slider-thumb-press", scale, 150*time.Millisecond)
}

func sliderPaintThumb(p *ui.Painter, t theme.Palette, thumb ui.Rect, innerW, innerH, innerRadius, scale float32, focused bool, radius float32) {
	w, h := innerW*scale, innerH*scale
	inner := ui.Rect{X: thumb.X + (thumb.W-w)/2, Y: thumb.Y + (thumb.H-h)/2, W: w, H: h}
	r := min(w, h) / 2
	p.Shadow(inner, r, 0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.05))
	p.Fill(inner, ui.RGB(255, 255, 255), r)
	p.Stroke(inner, t.Input, r, theme.Hairline)
	if focused {
		p.Stroke(choiceExpand(inner, 3), t.Ring.Alpha(0.24), r+3, 3)
	}
}

// sliderMarks draws labelled ticks below a horizontal slider. The ticks and
// labels sit in a row inset by the track's horizontal padding, so a mark at a
// value lines up with the thumb there instead of with the track's edge.
func sliderMarks(c *ui.Context, o SliderOptions, min, max float64, rtl bool, pad float32) {
	t := theme.Colors(c)
	height := theme.Scale(c, 16)
	tickH := theme.Scale(c, 4)
	fracs := make([]float32, len(o.Marks))
	for i, mark := range o.Marks {
		fracs[i] = sliderFraction(mark.Value, min, max)
	}
	ui.Row(c).FillWidth().Height(height).PaddingX(pad).Children(func() {
		row := ui.Box(c).FillWidth().Height(height)
		row.Draw(func(p *ui.Painter, r ui.Rect) {
			for _, f := range fracs {
				x := r.X + r.W*f
				if rtl {
					x = r.X + r.W - r.W*f
				}
				p.Fill(ui.Rect{X: x - 0.5, Y: r.Y, W: 1, H: tickH}, t.Muted, 0)
			}
		})
		row.Children(func() {
			for i, mark := range o.Marks {
				if mark.Label == "" {
					continue
				}
				style.Text(ui.Text(c, mark.Label), c, theme.XS()).TextColor(t.Muted).
					Absolute().Top(0).StartPercent(fracs[i] * 100).MarginStart(-24).Width(48).
					TextAlign(ui.Center)
			}
		})
	})
}

// sliderOutput renders the formatted value or range.
func sliderOutput(c *ui.Context, o SliderOptions, state SliderState, t theme.Palette) {
	text := ""
	format := o.Format
	if format == nil {
		format = sliderNumberFormat(0)
	}
	switch {
	case state.IsRange():
		text = format(*state.Low) + " – " + format(*state.High)
	case state.Value != nil:
		text = format(*state.Value)
	}
	style.Text(ui.Text(c, text), c, theme.SM()).FontWeight(theme.WeightMedium).
		FontFeatures("tnum").TextColor(t.Foreground).SingleLine()
}

// sliderNumberFormat returns the output formatter for a slider's step: whole
// decimals for a snapping step, at most two otherwise.
func sliderNumberFormat(step float64) func(float64) string {
	if step > 0 {
		d := sliderDecimals(step)
		return func(v float64) string { return strconv.FormatFloat(v, 'f', d, 64) }
	}
	return func(v float64) string {
		return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
	}
}

// sliderDecimals returns the decimal places a step needs, up to six.
func sliderDecimals(step float64) int {
	d := 0
	for step < 1 && step > 0 && d < 6 {
		step *= 10
		d++
	}
	return d
}

// sliderFraction maps v within lo..hi to 0..1, clamped.
func sliderFraction(v, lo, hi float64) float32 {
	if hi <= lo {
		return 0
	}
	f := (v - lo) / (hi - lo)
	if f <= 0 {
		return 0
	}
	if f >= 1 {
		return 1
	}
	return float32(f)
}
