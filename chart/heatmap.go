package chart

import (
	"fmt"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

// Cell identifies a bin by zero-based column and row. Missing cells are empty.
type Cell struct {
	Column, Row int
	Value       float64
	Label       string
}

type HeatmapOptions struct {
	Options
	Columns, Rows []string
	Cells         []Cell
	// Thresholds are the four increasing boundaries between five color levels.
	// Nil derives equally spaced thresholds from the nonnegative values.
	Thresholds []float64
}

func Heatmap(c *ui.Context, o HeatmapOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	seen := map[[2]int]bool{}
	hi := 0.0
	for i, v := range o.Cells {
		key := [2]int{v.Column, v.Row}
		inside := v.Column >= 0 && v.Column < len(o.Columns) && v.Row >= 0 && v.Row < len(o.Rows)
		if !inside || seen[key] || !finite(v.Value) || v.Value < 0 {
			return ui.Element{}, fmt.Errorf("heatmap: invalid, duplicate or out-of-range cell %d", i)
		}
		seen[key] = true
		hi = max(hi, v.Value)
	}
	thresholds := o.Thresholds
	if thresholds == nil {
		thresholds = []float64{0, hi * .25, hi * .5, hi * .75}
	} else {
		if len(thresholds) != 4 {
			return ui.Element{}, fmt.Errorf("heatmap: four level thresholds are required")
		}
		for i, v := range thresholds {
			if !finite(v) || v < 0 || (i > 0 && v <= thresholds[i-1]) {
				return ui.Element{}, fmt.Errorf("heatmap: thresholds must be finite, nonnegative and increasing")
			}
		}
	}
	t := theme.Colors(c)
	levels := []ui.Color{t.MutedBackground, t.MutedBackground.Mix(t.Chart1, .25), t.MutedBackground.Mix(t.Chart1, .5), t.MutedBackground.Mix(t.Chart1, .75), t.Chart1}
	level := func(v float64) int {
		if v == 0 {
			return 0
		}
		for i, limit := range thresholds {
			if v <= limit {
				return i
			}
		}
		return 4
	}
	x, y := native.NewBand(o.Columns...), native.NewBand(o.Rows...)
	describe := func(i int) string {
		v := o.Cells[i]
		return fmt.Sprintf("%s · %s: %s", o.Columns[v.Column], o.Rows[v.Row], number(v.Value))
	}
	return plotFrame(c, o.Options, native.Axes{
		Bottom: x.Ticks(), Left: y.Ticks(), Label: Axes(c, PlotContext{}, Axis{}, Axis{}).Label,
	}, len(o.Cells), describe, func(area ui.Element, inner Options) {
		area.Children(func() {
			surface(c, inner, len(o.Cells), describe, func(p *ui.Painter, r ui.Rect, selected int) {
				w, h := r.W/float32(max(1, len(o.Columns))), r.H/float32(max(1, len(o.Rows)))
				for i, cell := range o.Cells {
					box := ui.Rect{X: r.X + float32(cell.Column)*w + 2, Y: r.Y + r.H - float32(cell.Row+1)*h + 2, W: max(0, w-4), H: max(0, h-4)}
					p.Fill(box, levels[level(cell.Value)], 2)
					if i == selected {
						p.Stroke(box, t.Ring, 2, 2)
					}
				}
			}, func(px, py float32) int {
				column := min(len(o.Columns)-1, int(px*float32(len(o.Columns))))
				row := min(len(o.Rows)-1, int((1-py)*float32(len(o.Rows))))
				for i, cell := range o.Cells {
					if cell.Column == column && cell.Row == row {
						return i
					}
				}
				return -1
			})
		})
	}), nil
}
