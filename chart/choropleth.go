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

func Choropleth(c *ui.Context, o ChoroplethOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	ids := map[string]bool{}
	loX, loY, hiX, hiY := float32(math.Inf(1)), float32(math.Inf(1)), float32(math.Inf(-1)), float32(math.Inf(-1))
	maxValue := 0.0
	for ri, region := range o.Regions {
		if region.ID == "" || ids[region.ID] || !finite(region.Value) || region.Value < 0 || len(region.Polygons) == 0 {
			return ui.Element{}, fmt.Errorf("choropleth: invalid region %d", ri)
		}
		ids[region.ID] = true
		maxValue = max(maxValue, region.Value)
		for _, polygon := range region.Polygons {
			if len(polygon.Rings) == 0 {
				return ui.Element{}, fmt.Errorf("choropleth: region %q has an empty polygon", region.ID)
			}
			for _, ring := range polygon.Rings {
				if len(ring) < 3 {
					return ui.Element{}, fmt.Errorf("choropleth: region %q needs at least three ring points", region.ID)
				}
				for i, pt := range ring {
					valid := finite(pt.Longitude) && finite(pt.Latitude) && math.Abs(pt.Longitude) <= 180 && math.Abs(pt.Latitude) <= 85.0511
					crosses := i > 0 && math.Abs(pt.Longitude-ring[i-1].Longitude) > 180
					if !valid || crosses {
						return ui.Element{}, fmt.Errorf("choropleth: invalid or antimeridian-crossing coordinates in %q", region.ID)
					}
					v := mercator(pt)
					loX, hiX, loY, hiY = min(loX, v.X), max(hiX, v.X), min(loY, v.Y), max(hiY, v.Y)
				}
			}
		}
	}
	if len(o.Regions) > 0 && (hiX <= loX || hiY <= loY) {
		return ui.Element{}, fmt.Errorf("choropleth: geometry must have nonzero geographic extent")
	}
	zoom, panX, panY := 1.0, 0.0, 0.0
	if v := o.Viewport; v != nil {
		if !finite(v.Zoom) || !finite(v.PanX) || !finite(v.PanY) || v.Zoom < 0 || v.Zoom > 16 {
			return ui.Element{}, fmt.Errorf("choropleth: viewport needs finite pan and zoom between 1 and 16")
		}
		zoom, panX, panY = max(1, v.Zoom), v.PanX, v.PanY
	}
	transform := func(r ui.Rect, pt GeoPoint) ui.Point {
		projected := mercator(pt)
		scale := min(r.W/(hiX-loX), r.H/(hiY-loY)) * .92 * float32(zoom)
		return ui.Point{X: r.X + r.W/2 + (projected.X-(loX+hiX)/2)*scale + float32(panX)*r.W,
			Y: r.Y + r.H/2 - (projected.Y-(loY+hiY)/2)*scale + float32(panY)*r.H}
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
			for i, region := range o.Regions {
				fill := t.Chart1.Alpha(.2 + .8*float32(region.Value/max(1, maxValue)))
				if selected == i {
					fill = t.Chart2
				}
				for _, polygon := range region.Polygons {
					path := &ui.Path{}
					for _, ring := range polygon.Rings {
						for j, pt := range ring {
							v := transform(r, pt)
							if j == 0 {
								path.MoveTo(v.X, v.Y)
							} else {
								path.LineTo(v.X, v.Y)
							}
						}
						path.Close()
					}
					p.FillPath(path, fill)
					p.StrokePath(path, 1, t.Background)
				}
			}
		}, func(x, y float32) int {
			bounds := o.State.Handle.Bounds(c)
			r := ui.Rect{W: bounds.W, H: bounds.H}
			for i := len(o.Regions) - 1; i >= 0; i-- {
				for _, polygon := range o.Regions[i].Polygons {
					inside := false
					for _, ring := range polygon.Rings {
						points := make([]ui.Point, len(ring))
						for j, pt := range ring {
							points[j] = transform(r, pt)
						}
						if polygonContains(points, x*r.W, y*r.H) {
							inside = !inside
						}
					}
					if inside {
						return i
					}
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
