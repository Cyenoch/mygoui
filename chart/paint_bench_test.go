package chart

import (
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// benchmarkFrames measures whole frames of a chart, build and paint.
func benchmarkFrames(b *testing.B, build func(*ui.Context) (ui.Element, error)) {
	b.Helper()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if _, err := build(c); err != nil {
			b.Fatal(err)
		}
	}, 640, 480)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tt.Frame()
	}
}

func BenchmarkLineFrame(b *testing.B) {
	series := make([]Series, 3)
	for s := range series {
		points := make([]Point, 500)
		for i := range points {
			points[i] = Point{X: float64(i), Y: math.Sin(float64(i+s*40) / 30)}
		}
		series[s] = Series{ID: string(rune('a' + s)), Label: string(rune('A' + s)), Points: points}
	}
	state := State{Selected: 250}
	benchmarkFrames(b, func(c *ui.Context) (ui.Element, error) {
		return Line(c, CartesianOptions{Options: Options{Label: "Line", State: &state}, Series: series})
	})
}

func BenchmarkChoroplethFrame(b *testing.B) {
	regions := make([]Region, 40)
	for i := range regions {
		ring := make([]GeoPoint, 64)
		for k := range ring {
			a := float64(k) / 64 * 2 * math.Pi
			ring[k] = GeoPoint{Longitude: float64(i%8)*20 - 80 + 8*math.Cos(a), Latitude: float64(i/8)*15 - 30 + 6*math.Sin(a)}
		}
		regions[i] = Region{ID: string(rune('a' + i)), Label: "Region", Value: float64(i), Polygons: []Polygon{{Rings: [][]GeoPoint{ring}}}}
	}
	state := State{Selected: -1}
	benchmarkFrames(b, func(c *ui.Context) (ui.Element, error) {
		return Choropleth(c, ChoroplethOptions{Options: Options{Label: "Map", State: &state}, Regions: regions, Viewport: &MapState{}})
	})
}

func BenchmarkSankeyFrame(b *testing.B) {
	nodes := []FlowNode{{"a", "A"}, {"b", "B"}, {"c", "C"}, {"d", "D"}, {"e", "E"}, {"f", "F"}}
	links := []FlowLink{{"a", "c", 5}, {"a", "d", 3}, {"b", "c", 2}, {"b", "d", 4}, {"c", "e", 6}, {"d", "e", 3}, {"d", "f", 4}, {"c", "f", 1}}
	state := State{Selected: -1}
	benchmarkFrames(b, func(c *ui.Context) (ui.Element, error) {
		return Sankey(c, SankeyOptions{Options: Options{Label: "Flow", State: &state}, Nodes: nodes, Links: links})
	})
}
