package chart

import (
	"fmt"
	"math"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// GeoPoint uses degrees in WGS84. Mercator excludes the poles.
type GeoPoint struct{ Longitude, Latitude float64 }

// Polygon contains an exterior ring followed by optional holes. Regions may
// contain multiple polygons. Coordinates must not cross the antimeridian.
type Polygon struct{ Rings [][]GeoPoint }

type Region struct {
	ID, Label string
	Value     float64
	Polygons  []Polygon
}

// MapState is a caller-owned viewport. Zero Zoom is the initial fit (1x).
type MapState struct {
	Zoom, PanX, PanY float64
}

type ChoroplethOptions struct {
	Options
	Regions  []Region
	Viewport *MapState
}

func mercator(pt GeoPoint) ui.Point {
	lat := pt.Latitude * math.Pi / 180
	return ui.Point{X: float32(pt.Longitude * math.Pi / 180), Y: float32(math.Log(math.Tan(math.Pi/4 + lat/2)))}
}

// projectedPolygon is one polygon of a region in the projected geometry:
// rings index ringEnds, and lo/hi bound its points.
type projectedPolygon struct {
	region           int
	firstRing, rings int
	lo, hi           ui.Point
}

// projectedMap is the Mercator geometry of a build's regions, projected once
// so paint and hit testing only scale and offset it.
type projectedMap struct {
	points   []ui.Point
	ringEnds []int
	polygons []projectedPolygon
	lo, hi   ui.Point
	maxValue float64
}

// ring returns the points of ring k.
func (m *projectedMap) ring(k int) []ui.Point {
	start := 0
	if k > 0 {
		start = m.ringEnds[k-1]
	}
	return m.points[start:m.ringEnds[k]]
}

func projectRegions(regions []Region) (projectedMap, error) {
	total := 0
	for _, region := range regions {
		for _, polygon := range region.Polygons {
			for _, ring := range polygon.Rings {
				total += len(ring)
			}
		}
	}
	inf := float32(math.Inf(1))
	m := projectedMap{points: make([]ui.Point, 0, total), lo: ui.Point{X: inf, Y: inf}, hi: ui.Point{X: -inf, Y: -inf}}
	ids := map[string]bool{}
	for ri, region := range regions {
		if region.ID == "" || ids[region.ID] || !finite(region.Value) || region.Value < 0 || len(region.Polygons) == 0 {
			return m, fmt.Errorf("choropleth: invalid region %d", ri)
		}
		ids[region.ID] = true
		m.maxValue = max(m.maxValue, region.Value)
		for _, polygon := range region.Polygons {
			if len(polygon.Rings) == 0 {
				return m, fmt.Errorf("choropleth: region %q has an empty polygon", region.ID)
			}
			projected := projectedPolygon{region: ri, firstRing: len(m.ringEnds), rings: len(polygon.Rings),
				lo: ui.Point{X: inf, Y: inf}, hi: ui.Point{X: -inf, Y: -inf}}
			for _, ring := range polygon.Rings {
				if len(ring) < 3 {
					return m, fmt.Errorf("choropleth: region %q needs at least three ring points", region.ID)
				}
				for i, pt := range ring {
					valid := finite(pt.Longitude) && finite(pt.Latitude) && math.Abs(pt.Longitude) <= 180 && math.Abs(pt.Latitude) <= 85.0511
					crosses := i > 0 && math.Abs(pt.Longitude-ring[i-1].Longitude) > 180
					if !valid || crosses {
						return m, fmt.Errorf("choropleth: invalid or antimeridian-crossing coordinates in %q", region.ID)
					}
					v := mercator(pt)
					m.points = append(m.points, v)
					projected.lo = ui.Point{X: min(projected.lo.X, v.X), Y: min(projected.lo.Y, v.Y)}
					projected.hi = ui.Point{X: max(projected.hi.X, v.X), Y: max(projected.hi.Y, v.Y)}
				}
				m.ringEnds = append(m.ringEnds, len(m.points))
			}
			m.polygons = append(m.polygons, projected)
			m.lo = ui.Point{X: min(m.lo.X, projected.lo.X), Y: min(m.lo.Y, projected.lo.Y)}
			m.hi = ui.Point{X: max(m.hi.X, projected.hi.X), Y: max(m.hi.Y, projected.hi.Y)}
		}
	}
	if len(regions) > 0 && (m.hi.X <= m.lo.X || m.hi.Y <= m.lo.Y) {
		return m, fmt.Errorf("choropleth: geometry must have nonzero geographic extent")
	}
	return m, nil
}

func Choropleth(c *ui.Context, o ChoroplethOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	m, err := projectRegions(o.Regions)
	if err != nil {
		return ui.Element{}, err
	}
	zoom, panX, panY := 1.0, 0.0, 0.0
	if v := o.Viewport; v != nil {
		if !finite(v.Zoom) || !finite(v.PanX) || !finite(v.PanY) || v.Zoom < 0 || v.Zoom > 16 {
			return ui.Element{}, fmt.Errorf("choropleth: viewport needs finite pan and zoom between 1 and 16")
		}
		zoom, panX, panY = max(1, v.Zoom), v.PanX, v.PanY
	}
	// view maps projected points into r: screen = origin + projected*scale,
	// with y flipped.
	view := func(r ui.Rect) (origin ui.Point, scale float32) {
		scale = min(r.W/(m.hi.X-m.lo.X), r.H/(m.hi.Y-m.lo.Y)) * .92 * float32(zoom)
		origin = ui.Point{X: r.X + r.W/2 - (m.lo.X+m.hi.X)/2*scale + float32(panX)*r.W,
			Y: r.Y + r.H/2 + (m.lo.Y+m.hi.Y)/2*scale + float32(panY)*r.H}
		return origin, scale
	}
	t := theme.Colors(c)
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8)
	root.Children(func() {
		if o.Viewport != nil {
			ui.Row(c).Gap(theme.Space4).Children(func() {
				for _, command := range []struct {
					label string
					run   func()
				}{
					{"Zoom in", func() { o.Viewport.Zoom = min(16, zoom*1.5) }},
					{"Zoom out", func() { o.Viewport.Zoom = max(1, zoom/1.5) }},
					{"Reset map", func() { *o.Viewport = MapState{} }},
				} {
					mygoui.Button(c, mygoui.ButtonOptions{Label: command.label, Size: theme.ButtonSmall, Variant: mygoui.ButtonOutline, Disabled: o.Disabled, OnClick: command.run})
				}
			})
		}
		surface(c, o.Options, len(o.Regions), func(i int) string {
			return o.Regions[i].Label + ": " + number(o.Regions[i].Value)
		}, func(p *ui.Painter, r ui.Rect, selected int) {
			origin, scale := view(r)
			for _, polygon := range m.polygons {
				fill := t.Chart1.Alpha(.2 + .8*float32(o.Regions[polygon.region].Value/max(1, m.maxValue)))
				if selected == polygon.region {
					fill = t.Chart2
				}
				path := &ui.Path{}
				for k := polygon.firstRing; k < polygon.firstRing+polygon.rings; k++ {
					for j, v := range m.ring(k) {
						x, y := origin.X+v.X*scale, origin.Y-v.Y*scale
						if j == 0 {
							path.MoveTo(x, y)
						} else {
							path.LineTo(x, y)
						}
					}
					path.Close()
				}
				p.FillPath(path, fill)
				p.StrokePath(path, 1, t.Background)
			}
		}, func(x, y float32) int {
			bounds := o.State.Handle.Bounds(c)
			origin, scale := view(ui.Rect{W: bounds.W, H: bounds.H})
			// The pointer in projected coordinates; the mapping is affine, so
			// containment is the same as on screen.
			at := ui.Point{X: (x*bounds.W - origin.X) / scale, Y: (origin.Y - y*bounds.H) / scale}
			for i := len(m.polygons) - 1; i >= 0; i-- {
				polygon := m.polygons[i]
				if at.X < polygon.lo.X || at.X > polygon.hi.X || at.Y < polygon.lo.Y || at.Y > polygon.hi.Y {
					continue
				}
				inside := false
				for k := polygon.firstRing; k < polygon.firstRing+polygon.rings; k++ {
					if polygonContains(m.ring(k), at.X, at.Y) {
						inside = !inside
					}
				}
				if inside {
					return polygon.region
				}
			}
			return -1
		})
	})
	return root, nil
}

func polygonContains(points []ui.Point, x, y float32) bool {
	inside := false
	j := len(points) - 1
	for i, a := range points {
		b := points[j]
		if (a.Y > y) != (b.Y > y) && x < (b.X-a.X)*(y-a.Y)/(b.Y-a.Y)+a.X {
			inside = !inside
		}
		j = i
	}
	return inside
}
