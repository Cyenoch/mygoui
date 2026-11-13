package overlay

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func viewerImages() []ImageViewerItem {
	src := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 200, 100)))
	return []ImageViewerItem{{Source: src, Label: "First"}, {Source: src, Label: "Second"}}
}

func viewerTester(open *bool, index *int, dismissed *int) *ui.Tester {
	images := viewerImages()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ImageViewer(c, open, ImageViewerOptions{
			Images: images, Index: index, OnDismiss: func() { *dismissed++ },
		})
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	return tt
}

func TestImageViewerZoomsAndStepsFromTheKeyboard(t *testing.T) {
	open, index, dismissed := true, 0, 0
	tt := viewerTester(&open, &index, &dismissed)
	first, ok := tt.AccessByRole(ui.RoleImage, "First")
	if !ok || first.Bounds != (ui.Rect{X: 220, Y: 190, W: 200, H: 100}) {
		t.Fatalf("fitted image: found %v at %+v", ok, first.Bounds)
	}
	tt.Key(0, ui.KeyEqual)
	zoomed, _ := tt.AccessByRole(ui.RoleImage, "First")
	if !tt.HasText("150%") || zoomed.Bounds.W != 300 {
		t.Fatalf("after +: texts %q, image %+v", tt.Texts(), zoomed.Bounds)
	}
	tt.Key(0, ui.KeyMinus)
	tt.Key(0, ui.KeyMinus)
	if !tt.HasText("100%") {
		t.Fatalf("zoom went below the fitted size: %q", tt.Texts())
	}
	if err := tt.Click("Zoom in"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.Key0)
	if !tt.HasText("100%") {
		t.Fatalf("0 did not reset the zoom: %q", tt.Texts())
	}

	tt.Key(0, ui.KeyEqual)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if _, ok := tt.AccessByRole(ui.RoleImage, "Second"); index != 1 || !ok || !tt.HasText("100%") {
		t.Fatalf("after Right: index %d, second shown %v, texts %q", index, ok, tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if open || dismissed != 1 {
		t.Fatalf("Escape: open %v, dismissed %d", open, dismissed)
	}
}

func TestImageViewerPansWhenZoomedAndClosesOutsideTheImage(t *testing.T) {
	open, index, dismissed := true, 0, 0
	tt := viewerTester(&open, &index, &dismissed)
	for range 4 {
		tt.Key(0, ui.KeyEqual)
	}
	// 200 DIPs at 506% overflow the 640 DIP window, so the image pans.
	before, _ := tt.AccessByRole(ui.RoleImage, "First")
	tt.Press(320, 240)
	tt.Move(280, 240)
	tt.Release(280, 240)
	tt.Frame()
	after, _ := tt.AccessByRole(ui.RoleImage, "First")
	if !open || after.Bounds.X >= before.Bounds.X {
		t.Fatalf("pan: open %v, image x %v -> %v", open, before.Bounds.X, after.Bounds.X)
	}
	tt.Key(0, ui.Key0)
	tt.ClickAt(320, 240)
	if !open {
		t.Fatal("a press on the image closed the viewer")
	}
	// A drag is not a press, wherever it ends.
	tt.Press(60, 240)
	tt.Move(100, 240)
	tt.Release(100, 240)
	tt.Frame()
	if !open {
		t.Fatal("a drag beside the image closed the viewer")
	}
	tt.ClickAt(60, 240)
	if open || dismissed != 1 {
		t.Fatalf("press outside the image: open %v, dismissed %d", open, dismissed)
	}
}
