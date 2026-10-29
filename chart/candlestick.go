package chart

import (
	"fmt"
	"math"
	"slices"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

type Candle struct {
	X, Open, High, Low, Close float64
	Label                     string
}

type CandlestickOptions struct {
	Options
	Candles []Candle
	XAxis   Axis
	YAxis   Axis
}

func Candlestick(c *ui.Context, o CandlestickOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	xs, values := []float64{}, []float64{}
	for i, candle := range o.Candles {
		valid := finite(candle.X) && finite(candle.Open) && finite(candle.High) && finite(candle.Low) && finite(candle.Close)
		ordered := candle.Low <= min(candle.Open, candle.Close) && candle.High >= max(candle.Open, candle.Close)
		if !valid || !ordered || (i > 0 && candle.X <= o.Candles[i-1].X) {
			return ui.Element{}, fmt.Errorf("candlestick: invalid OHLC or unordered X at candle %d", i)
		}
		xs = append(xs, candle.X)
		values = append(values, candle.Low, candle.High)
	}
	x, err := scale(xs, nil)
	if err != nil {
		return ui.Element{}, fmt.Errorf("candlestick X: %w", err)
	}
	y, err := scale(values, nil)
	if err != nil {
		return ui.Element{}, fmt.Errorf("candlestick Y: %w", err)
	}
	step := 1.0
	if len(xs) > 1 {
		step = math.Inf(1)
		for i := 1; i < len(xs); i++ {
			step = min(step, xs[i]-xs[i-1])
		}
		x.Lo -= step / 2
		x.Hi += step / 2
	}
	ctx := PlotContext{X: x, Y: y}
	t := theme.Colors(c)
	describe := func(i int) string {
		v := o.Candles[i]
		return fmt.Sprintf("%s: open %s, high %s, low %s, close %s", v.Label, number(v.Open), number(v.High), number(v.Low), number(v.Close))
	}
	return plotFrame(c, o.Options, Axes(c, ctx, o.XAxis, o.YAxis), len(xs), describe, func(area ui.Element, inner Options) {
		area.Children(func() {
			surface(c, inner, len(xs), describe, func(p *ui.Painter, r ui.Rect, selected int) {
				Grid(p, r, ctx, t.Border)
				w := float32(step/(x.Hi-x.Lo)) * r.W * .55
				for i, candle := range o.Candles {
					px := native.X(r, x.At(candle.X))
					fill := t.Success
					if candle.Close < candle.Open {
						fill = t.Destructive
					}
					p.Line(px, native.Y(r, y.At(candle.Low)), px, native.Y(r, y.At(candle.High)), 1, fill)
					top, bottom := native.Y(r, y.At(candle.Open)), native.Y(r, y.At(candle.Close))
					p.Fill(ui.Rect{X: px - w/2, Y: min(top, bottom), W: w, H: max(1, abs(top-bottom))}, fill, 1)
					if selected == i {
						p.Line(px, r.Y, px, r.Y+r.H, 1, t.Ring)
					}
				}
			}, func(px, _ float32) int {
				v := x.Invert(px)
				i, _ := slices.BinarySearch(xs, v)
				if i > 0 && (i == len(xs) || v-xs[i-1] < xs[i]-v) {
					i--
				}
				return min(len(xs)-1, i)
			})
		})
	}), nil
}
