package chart

import (
	"fmt"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

type AnnotationKind uint8

const (
	ReferenceArea AnnotationKind = iota
	ProjectionLine
	Band
	Weekends
	Deadline
	TrendLine
	Marker
)

// Annotation uses data coordinates. Band connects Points (lower) and Upper;
// TrendLine fits Points by least squares. Weekends interprets X as Unix seconds
// in Location (UTC by default). Other kinds use X1/Y1 and X2/Y2.
type Annotation struct {
	Kind           AnnotationKind
	Label          string
	X1, X2, Y1, Y2 float64
	Points, Upper  []Point
	Color          ui.Color
	Location       *time.Location
}

func validateAnnotations(annotations []Annotation) error {
	for i, a := range annotations {
		if a.Kind > Marker || !finite(a.X1) || !finite(a.X2) || !finite(a.Y1) || !finite(a.Y2) {
			return fmt.Errorf("chart: invalid annotation %d", i)
		}
		if a.Kind == Band && len(a.Points) != len(a.Upper) {
			return fmt.Errorf("chart: band %d requires aligned lower and upper bounds", i)
		}
		for j, pt := range a.Points {
			if !finite(pt.X) || !finite(pt.Y) || (j > 0 && pt.X <= a.Points[j-1].X) {
				return fmt.Errorf("chart: annotation %d point %d must be finite with increasing X", i, j)
			}
			if a.Kind == Band && (a.Upper[j].X != pt.X || !finite(a.Upper[j].Y) || a.Upper[j].Y < pt.Y) {
				return fmt.Errorf("chart: band %d has invalid upper bound at %d", i, j)
			}
		}
	}
	return nil
}

func validateWeekendDomain(x native.Linear) error {
	const century = 100 * 366 * 86400
	if x.Lo < -62135596800 || x.Hi > 253402300799 || x.Hi-x.Lo > century {
		return fmt.Errorf("chart weekends: domain must fit years 1..9999 and span at most 100 years")
	}
	return nil
}

// Grid paints data-space grid lines without replacing native axis labels.
func Grid(p *ui.Painter, r ui.Rect, ctx PlotContext, color ui.Color) {
	for _, tick := range ctx.Y.Ticks(4) {
		p.Line(r.X, native.Y(r, tick.At), r.X+r.W, native.Y(r, tick.At), 1, color)
	}
}

type BackgroundPattern uint8

const (
	BackgroundSolid BackgroundPattern = iota
	BackgroundDots
	BackgroundLines
)

// Background paints a chart backdrop with optional restrained dots or hatching.
func Background(p *ui.Painter, r ui.Rect, fill, ink ui.Color, pattern BackgroundPattern) {
	p.Fill(r, fill, 0)
	if pattern == BackgroundSolid {
		return
	}
	for y := r.Y + 8; y < r.Y+r.H; y += 12 {
		if pattern == BackgroundLines {
			p.Line(r.X, y, r.X+r.W, y, 1, ink)
			continue
		}
		for x := r.X + 8; x < r.X+r.W; x += 12 {
			p.Fill(ui.Rect{X: x, Y: y, W: 1.5, H: 1.5}, ink, 1)
		}
	}
}

func drawAnnotation(p *ui.Painter, r ui.Rect, ctx PlotContext, a Annotation, t theme.Palette, textSize float32) {
	color := a.Color
	if color == (ui.Color{}) {
		color = t.MutedForeground
	}
	pos := func(x, y float64) ui.Point { return point(r, ctx.X.At(x), ctx.Y.At(y)) }
	start, end := pos(a.X1, a.Y1), pos(a.X2, a.Y2)
	switch a.Kind {
	case ReferenceArea:
		p.Fill(ui.Rect{X: min(start.X, end.X), Y: min(start.Y, end.Y), W: abs(end.X - start.X), H: abs(end.Y - start.Y)}, color.Alpha(0.12), 0)
	case ProjectionLine:
		path := (&ui.Path{}).MoveTo(start.X, start.Y).LineTo(end.X, end.Y).Dash(0, 4, 4)
		p.StrokePath(path, 1.5, color)
	case Band:
		x, low, high := []float32{}, []float32{}, []float32{}
		for i, pt := range a.Points {
			x, low, high = append(x, ctx.X.At(pt.X)), append(low, ctx.Y.At(pt.Y)), append(high, ctx.Y.At(a.Upper[i].Y))
		}
		p.FillPath(native.Area(r, x, low, high, native.CurveLinear), color.Alpha(0.15))
	case Weekends:
		loc := a.Location
		if loc == nil {
			loc = time.UTC
		}
		from := time.Unix(int64(ctx.X.Lo), 0).In(loc)
		day := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
		// Iterate weeks, not every day in a multi-year time domain.
		day = day.AddDate(0, 0, -((int(day.Weekday()) + 1) % 7))
		for float64(day.Unix()) <= ctx.X.Hi {
			x1 := native.X(r, ctx.X.At(float64(day.Unix())))
			x2 := native.X(r, ctx.X.At(float64(day.AddDate(0, 0, 2).Unix())))
			p.Fill(ui.Rect{X: x1, Y: r.Y, W: x2 - x1, H: r.H}, color.Alpha(0.08), 0)
			day = day.AddDate(0, 0, 7)
		}
	case Deadline:
		p.Line(start.X, r.Y, start.X, r.Y+r.H, 1, color)
	case TrendLine:
		if len(a.Points) < 2 {
			return
		}
		mx, my := 0.0, 0.0
		for _, pt := range a.Points {
			mx += pt.X / float64(len(a.Points))
			my += pt.Y / float64(len(a.Points))
		}
		cov, variance := 0.0, 0.0
		for _, pt := range a.Points {
			cov += (pt.X - mx) * (pt.Y - my)
			variance += (pt.X - mx) * (pt.X - mx)
		}
		if variance == 0 {
			return
		}
		slope := cov / variance
		line(p, pos(ctx.X.Lo, my+slope*(ctx.X.Lo-mx)), pos(ctx.X.Hi, my+slope*(ctx.X.Hi-mx)), 1.5, color)
	case Marker:
		p.Fill(ui.Rect{X: start.X - 4, Y: start.Y - 4, W: 8, H: 8}, color, 4)
	}
	if a.Label != "" {
		p.Text(max(r.X, min(r.X+r.W-60, start.X+6)), max(r.Y, start.Y-textSize-6), a.Label, textSize, color)
	}
}
