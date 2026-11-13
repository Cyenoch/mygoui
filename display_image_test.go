package mygoui_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func redImage(t *testing.T) ui.ImageSource {
	t.Helper()
	pixels := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for i := range 16 {
		pixels.SetRGBA(i%4, i/4, color.RGBA{R: 255, A: 255})
	}
	var data bytes.Buffer
	if err := png.Encode(&data, pixels); err != nil {
		t.Fatal(err)
	}
	source, err := ui.DecodeImage(data.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	return source
}

// TestImageShowsTheCallerStatusInItsBox covers the picture's three states:
// each keeps the box and the accessible name, and only a shown picture is
// pressable.
func TestImageShowsTheCallerStatusInItsBox(t *testing.T) {
	source := redImage(t)
	for _, tc := range []struct {
		name     string
		source   ui.ImageSource
		status   mygoui.ImageStatus
		picture  bool
		fallback bool
	}{
		{name: "ready", source: source, picture: true},
		{name: "loading", source: source, status: mygoui.ImageLoading},
		{name: "failed", source: source, status: mygoui.ImageFailed, fallback: true},
		{name: "no picture", source: (*ui.Bitmap)(nil), fallback: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clicks := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				mygoui.Image(c, mygoui.ImageOptions{
					Source: tc.source, Label: "Team photo", Status: tc.status,
					Width: 120, Height: 80,
					Fallback: func() { ui.Text(c, "Unavailable") },
					OnClick:  func() { clicks++ },
				})
			}, 200, 160)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			tt.Frame()
			box := findRect(t, tt, "Team photo")
			if box.W != 120 || box.H != 80 {
				t.Fatalf("box = %+v, want 120x80", box)
			}
			if tt.HasText("Unavailable") != tc.fallback {
				t.Fatalf("fallback visible = %v, want %v", !tc.fallback, tc.fallback)
			}
			red := pixel(tt.Image(), box.X+box.W/2, box.Y+box.H/2) == ui.RGB(255, 0, 0)
			if red != tc.picture {
				t.Fatalf("picture painted = %v, want %v", red, tc.picture)
			}
			tt.ClickAt(box.X+10, box.Y+10)
			if want := map[bool]int{true: 1}[tc.picture]; clicks != want {
				t.Fatalf("clicks = %d, want %d", clicks, want)
			}
		})
	}
}
