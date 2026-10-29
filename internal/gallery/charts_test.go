package gallery

import (
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestChartCatalogRendersEveryPublicExample(t *testing.T) {
	for _, page := range Pages() {
		if page.Group != "Charts" && page.Group != "Chart Utilities" {
			continue
		}
		t.Run(page.ID, func(t *testing.T) {
			state := New()
			state.Selected = page.ID
			tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1000, 800)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			if _, ok := tt.Find("Chart specimen"); !ok {
				t.Fatalf("no public chart rendered: %v", tt.Texts())
			}
			for _, text := range tt.Texts() {
				if strings.Contains(text, "invalid") || strings.Contains(text, "must be") {
					t.Fatalf("demo error: %s", text)
				}
			}
		})
	}
}

// Opt in to persist native headless evidence; no browser or image generator.
func TestChartVisualArtifacts(t *testing.T) {
	dir := os.Getenv("MYGOUI_CHART_ARTIFACT_DIR")
	if dir == "" {
		t.Skip("set MYGOUI_CHART_ARTIFACT_DIR to collect chart specimens")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"linechart", "barchart", "candlestickchart", "piechart", "radarchart", "radialchart", "gaugechart", "heatmapchart", "choroplethchart", "sankeychart", "sunburstchart", "chartband", "chartbrush"} {
		contact := image.NewRGBA(image.Rect(0, 0, 2160, 1500))
		for i, mode := range []struct {
			name      string
			width     int
			dark      bool
			textScale float32
		}{
			{"wide-light", 1180, false, 1}, {"wide-dark", 980, true, 1},
			{"compact-light", 640, false, 1}, {"compact-dark", 640, true, 1.3},
		} {
			state := New()
			state.Selected, state.Dark = id, mode.dark
			tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, mode.width, 750)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: mode.textScale})
			tt.Frame()
			r, ok := tt.Find("Chart specimen")
			if !ok || r.W <= 0 || r.H <= 0 || r.X < 0 || r.X+r.W > float32(mode.width) || r.Y+r.H > 750 {
				t.Fatalf("%s/%s chart clipped: %+v", id, mode.name, r)
			}
			img := tt.Image()
			writeCossPNG(t, filepath.Join(dir, id+"-"+mode.name+".png"), img)
			x, y := 0, (i/2)*750
			if i%2 == 1 {
				x = 1180
			}
			draw.Draw(contact, image.Rect(x, y, x+mode.width, y+750), img, image.Point{}, draw.Src)
		}
		writeCossPNG(t, filepath.Join(dir, id+"-contact.png"), contact)
	}
}
