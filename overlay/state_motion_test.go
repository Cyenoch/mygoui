package overlay

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func solidImage(c color.RGBA) ui.ImageSource {
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	for y := range 100 {
		for x := range 200 {
			img.SetRGBA(x, y, c)
		}
	}
	return ui.NewBitmap(img)
}

func TestImageViewerZoomGlidesAndNextImageCrossfades(t *testing.T) {
	open, index := true, 0
	images := []ImageViewerItem{
		{Source: solidImage(color.RGBA{R: 255, A: 255}), Label: "Red"},
		{Source: solidImage(color.RGBA{B: 255, A: 255}), Label: "Blue"},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ImageViewer(c, &open, ImageViewerOptions{Images: images, Index: &index})
	}, 640, 480)
	tt.Frame()
	tt.Advance(400 * time.Millisecond)

	// + takes the image to 150% over a few frames.
	tt.Key(0, ui.KeyEqual)
	tt.Advance(50 * time.Millisecond)
	mid, _ := tt.AccessByRole(ui.RoleImage, "Red")
	if mid.Bounds.W <= 200 || mid.Bounds.W >= 300 {
		t.Fatalf("the image is %v DIP wide partway to 300", mid.Bounds.W)
	}
	tt.Advance(400 * time.Millisecond)
	if end, _ := tt.AccessByRole(ui.RoleImage, "Red"); end.Bounds.W != 300 {
		t.Fatalf("the zoom settled at %v DIP, want 300", end.Bounds.W)
	}

	// The next image fades in over the one leaving.
	tt.Key(0, ui.Key0)
	tt.Advance(400 * time.Millisecond)
	tt.Key(0, ui.KeyRight)
	tt.Advance(50 * time.Millisecond)
	px := tt.Image().RGBAAt(320, 240)
	if px.R < 20 || px.B < 20 {
		t.Fatalf("pixel %v partway through the step is not a mix of both images", px)
	}
	tt.Advance(400 * time.Millisecond)
	if px := tt.Image().RGBAAt(320, 240); px.R > 5 || px.B < 250 {
		t.Fatalf("pixel %v after the step is not the second image", px)
	}
}

func TestTourStepCrossfadesItsText(t *testing.T) {
	state := TourState{Open: true}
	var first, second ui.Handle
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Row(c).Padding(16).Gap(120).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Inbox"}).Bind(&first)
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Archive"}).Bind(&second)
			Tour(c, &state, []TourStep{
				{Anchor: &first, Title: "Your inbox"},
				{Anchor: &second, Title: "The archive"},
			}, TourOptions{})
		})
	}, 640, 360)
	tt.Advance(400 * time.Millisecond)
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	// The new title fades in over the second half of the step change.
	tt.Advance(200 * time.Millisecond)
	heading, ok := tt.AccessByRole(ui.RoleHeading, "The archive")
	if !ok {
		t.Fatal("the second step is not shown")
	}
	title := heading.Bounds
	pixels := func() []color.RGBA {
		img := tt.Image()
		var out []color.RGBA
		for y := int(title.Y); y < int(title.Y+title.H); y++ {
			for x := int(title.X); x < int(title.X+title.W); x++ {
				out = append(out, img.RGBAAt(x, y))
			}
		}
		return out
	}
	mid := pixels()
	tt.Advance(400 * time.Millisecond)
	settled := pixels()
	apart := 0
	for i := range mid {
		if d := int(mid[i].R) - int(settled[i].R); d > 8 || d < -8 {
			apart++
		}
	}
	if apart < 20 {
		t.Fatalf("the title partway through the step is the settled one (%d pixels differ): it did not crossfade", apart)
	}
}
