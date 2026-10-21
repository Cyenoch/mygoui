package chart_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/chart"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestChartChangedDataKeyboardVisibilityAndAxes(t *testing.T) {
	data := chart.Data{Labels: []string{"Jan", "Feb", "Mar"}, Series: []chart.Series{
		{ID: "desktop", Label: "Desktop", Color: chart.Chart1, Values: []chart.Value{{Number: 10}, {Number: 20}, {Number: 0}}},
		{ID: "mobile", Label: "Mobile", Color: chart.Chart2, Values: []chart.Value{{Number: 3}, {Missing: true}, {Number: 5}}},
	}}
	s := chart.State{Selected: -1}
	o := chart.Options{Label: "Traffic", Data: data, Height: 200}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); chart.Chart(c, &s, o) }, 420, 310)
	if _, ok := tt.Find("Jan"); !ok {
		t.Fatal("visible axis did not build its category label")
	}
	before := bytes.Clone(tt.Image().Pix)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if s.Selected != 0 || !tt.HasText("Jan · Desktop: 10 · Mobile: 3") {
		t.Fatalf("keyboard authoritative reading: index=%d text=%v", s.Selected, tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	if s.Selected != 1 || !tt.HasText("Feb · Desktop: 20 · Mobile: missing") {
		t.Fatalf("missing value invented: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("Mar · Desktop: 0 · Mobile: 5") {
		t.Fatal("zero was treated as missing")
	}
	tt.Key(0, ui.KeyEscape)
	if s.Selected != -1 {
		t.Fatal("Escape did not clear reading")
	}
	o.Data.Series[0].Values[1].Number = 80
	tt.Frame()
	if bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("changed data did not update actual chart pixels")
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Feb · Desktop: 80 · Mobile: missing") {
		t.Fatal("selection read stale cached data")
	}
	if err := tt.Click("Desktop"); err != nil {
		t.Fatal(err)
	}
	if !s.Hidden["desktop"] || s.Selected != -1 {
		t.Fatal("legend did not toggle visibility and clear old reading")
	}
	if _, ok := tt.Find("Feb · Desktop: 80"); ok {
		t.Fatal("hidden series still constructed plotted marks")
	}
	o.HideAxes = true
	tt.Frame()
	if _, ok := tt.Find("Jan"); ok {
		t.Fatal("hidden axis still built labels")
	}
	if err := tt.Click("Desktop"); err != nil {
		t.Fatal(err)
	}
	if s.Hidden["desktop"] {
		t.Fatal("legend did not restore series")
	}
}

func TestAllChartKindsRenderSparseZeroEmptyAndNarrow(t *testing.T) {
	for _, kind := range []chart.ChartKind{chart.Bar, chart.Line, chart.Area, chart.Pie, chart.Radar, chart.Radial} {
		t.Run(fmt.Sprint(kind), func(t *testing.T) {
			s := chart.State{Selected: -1}
			data := chart.Data{Labels: []string{"A", "B", "C", "D"}, Series: []chart.Series{{ID: "one", Label: "One", Values: []chart.Value{{Number: 3}, {Missing: true}, {Number: 8}, {Number: 0}}}, {ID: "two", Label: "Two", Color: chart.Chart2, Values: []chart.Value{{Number: 5}, {Number: 4}}}}}
			o := chart.Options{Label: "Sample", Kind: kind, Data: data, Height: 180}
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); chart.Chart(c, &s, o) }, 240, 310)
			original := bytes.Clone(tt.Image().Pix)
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyHome)
			if s.Selected != 0 {
				t.Fatal("native keyboard did not select first reading")
			}
			if kind == chart.Pie || kind == chart.Radial {
				if !tt.HasText("One: 3") {
					t.Fatalf("polar selection=%v", tt.Texts())
				}
			} else if !tt.HasText("A · One: 3 · Two: 5") {
				t.Fatalf("selection=%v", tt.Texts())
			}
			s.Selected = -1
			for i := range o.Data.Series {
				for j := range o.Data.Series[i].Values {
					o.Data.Series[i].Values[j] = chart.Value{}
				}
			}
			tt.Frame()
			if !tt.HasText("All values are zero") {
				t.Fatalf("zero chart fabricated data: %v", tt.Texts())
			}
			if bytes.Equal(original, tt.Image().Pix) {
				t.Fatal("zero data left previous native geometry")
			}
			o.Data = chart.Data{}
			tt.Frame()
			if !tt.HasText("No data") {
				t.Fatalf("empty chart=%v", tt.Texts())
			}
			tt.SetSize(180, 310)
			tt.SetPreferences(ui.Preferences{TextScale: 2})
			tt.Frame()
			if _, ok := tt.Find("Sample"); !ok {
				t.Fatal("narrow large-text chart disappeared")
			}
		})
	}
}

func TestStackedChartVisibilityChangesScaleAndMarks(t *testing.T) {
	for _, stack := range []chart.StackMode{chart.Stacked, chart.Normalized, chart.Centered} {
		t.Run(fmt.Sprint(stack), func(t *testing.T) {
			s := chart.State{Selected: -1}
			data := chart.Data{Labels: []string{"First", "Second"}, Series: []chart.Series{{ID: "a", Label: "A", Values: []chart.Value{{Number: 10}, {Number: -5}}}, {ID: "b", Label: "B", Color: chart.Chart2, Values: []chart.Value{{Number: 20}, {Number: 15}}}}}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				chart.Chart(c, &s, chart.Options{Label: "Stack", Data: data, Stack: stack})
			}, 400, 320)
			before := bytes.Clone(tt.Image().Pix)
			tt.Click("B")
			if !s.Hidden["b"] || bytes.Equal(before, tt.Image().Pix) {
				t.Fatal("stack visibility did not update actual geometry")
			}
			tt.Key(0, ui.KeyTab)
		})
	}
}
