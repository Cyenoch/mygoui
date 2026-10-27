package gallery

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/fonts"
	"github.com/egoist/mygo/ui"
)

// Opt in with MYGOUI_COSS_ARTIFACT_DIR. Font registration belongs to this
// explicit render entry, never package initialization or gallery.View.
func TestCossVisualArtifacts(t *testing.T) {
	dir := os.Getenv("MYGOUI_COSS_ARTIFACT_DIR")
	if dir == "" {
		t.Skip("set MYGOUI_COSS_ARTIFACT_DIR to render native gallery PNGs")
	}
	if err := fonts.Register(); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	modes := []struct {
		name  string
		width int
		dark  bool
	}{
		{"wide-light", 1180, false}, {"wide-dark", 1180, true},
		{"compact-light", 640, false}, {"compact-dark", 640, true},
	}
	for _, specimen := range []struct{ page, control, name string }{
		{"theming", "Primary 2", "theming"}, {"overview", "Browse components", "overview"}, {"button", "Default", "button"}, {"badge", "Destructive", "badge"}, {"badge", "Inbox", "badge-anchors"},
		{"frame", "Save project", "frame"}, {"group", "Repository query", "group"}, {"command", "Search actions", "command"},
		{"segmentedcontrol", "Grid view", "segmentedcontrol"}, {"composer", "Message draft", "composer"}, {"sidebar", "Projects", "sidebar"},
		{"scrollarea", "Activity", "scrollarea"}, {"scrollfade", "Masked activity", "scrollfade"},
		{"scrollfade", "Surface feed", "scrollfade-surface"},
		{"scrollfade", "Shadow feed", "scrollfade-shadow"},
		{"sheet", "Open right sheet", "sheet"}, {"previewcard", "Ada Lovelace", "previewcard"},
	} {
		sheet := image.NewRGBA(image.Rect(0, 0, 2360, 1800))
		draw.Draw(sheet, sheet.Bounds(), image.NewUniform(color.RGBA{R: 232, G: 232, B: 232, A: 255}),
			image.Point{}, draw.Src)
		for i, mode := range modes {
			s := New()
			s.Selected, s.Dark = specimen.page, mode.dark
			if specimen.page == "theming" {
				s.Primary, s.Radius = 1, 16
			}
			tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, mode.width, 900)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			r, ok := tt.Find(specimen.control)
			if !ok {
				t.Fatalf("%s/%s: missing actual component %q", specimen.page, mode.name, specimen.control)
			}
			if !contentVisible(r, float32(mode.width), 900) {
				s.Content.Y = max(0, r.Y-360)
				tt.Frame()
			}
			if specimen.page == "sheet" {
				if err := tt.Click(specimen.control); err != nil {
					t.Fatal(err)
				}
				if !tt.HasText("Workspace settings") {
					t.Fatal("sheet PNG would show a closed example")
				}
			}
			if specimen.page == "previewcard" {
				r, _ = tt.Find(specimen.control)
				tt.Move(r.X+r.W/2, r.Y+r.H/2)
				tt.Advance(300 * time.Millisecond)
				if !tt.HasText("Follow updates") {
					t.Fatal("preview PNG would show a closed example")
				}
			}
			if specimen.page == "scrollfade" {
				cossStateOf(s, "scrollfade").Scroll[3].Y = 56
			}
			tt.Advance(400 * time.Millisecond)
			tt.Frame()
			r, ok = tt.Find(specimen.control)
			if !ok || !contentVisible(r, float32(mode.width), 900) {
				t.Fatalf("%s/%s: actual component is outside the capture: %+v", specimen.page, mode.name, r)
			}
			img := tt.Image()
			writeCossPNG(t, filepath.Join(dir, specimen.name+"-"+mode.name+".png"), img)
			x, y := (i%2)*1180, (i/2)*900
			draw.Draw(sheet, image.Rect(x, y, x+mode.width, y+900), img, image.Point{}, draw.Src)
		}
		writeCossPNG(t, filepath.Join(dir, specimen.name+"-contact.png"), sheet)
	}
}

func TestWorkbenchVisualArtifacts(t *testing.T) {
	dir := os.Getenv("MYGOUI_COSS_ARTIFACT_DIR")
	if dir == "" {
		t.Skip("set MYGOUI_COSS_ARTIFACT_DIR to render native gallery PNGs")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []struct {
		name       string
		width      int
		panel      int
		nav        bool
		appearance bool
		bar        ui.TitleBar
		textScale  float32
	}{
		{"workbench-wide", 1180, 0, false, false, ui.TitleBar{Height: 52, Left: 80}, 1},
		{"workbench-category", 1180, 0, false, false, ui.TitleBar{Height: 52, Left: 80}, 1},
		{"workbench-narrow", 360, 0, false, false, ui.TitleBar{Height: 52, Right: 138}, 1},
		{"workbench-search", 360, 0, true, false, ui.TitleBar{Height: 52, Right: 138}, 1},
		{"workbench-source", 1180, 1, false, false, ui.TitleBar{Height: 52, Left: 80}, 1},
		{"workbench-appearance", 360, 0, false, true, ui.TitleBar{Height: 52, Right: 138}, 1},
		{"workbench-enlarged", 360, 0, false, false, ui.TitleBar{Height: 52, Left: 80, Right: 138}, 1.75},
		{"workbench-very-large", 360, 0, false, false, ui.TitleBar{Height: 52, Left: 80, Right: 138}, 2.5},
	} {
		s := New()
		s.Selected, s.Panel, s.NavOpen, s.AppearanceOpen = "button", mode.panel, mode.nav, mode.appearance
		if mode.name == "workbench-category" {
			s.OpenGroup = "Actions"
		}
		tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, mode.width, 900)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: mode.textScale})
		tt.SetTitleBar(mode.bar)
		writeCossPNG(t, filepath.Join(dir, mode.name+".png"), tt.Image())
	}
}

func writeCossPNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}
