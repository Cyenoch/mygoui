package mygoui_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestInjectedPaletteReachesComponentPaint traces one host palette through the
// components a consuming application actually shows: an action button, an
// outlined button (whose border is transparent), a selectable list, a dialog
// panel and a toast card.
func TestInjectedPaletteReachesComponentPaint(t *testing.T) {
	theme.Configure(theme.Preferences{})
	t.Cleanup(func() { theme.Configure(theme.Preferences{}) })

	canvas := ui.RGB(10, 20, 30)
	sheet := ui.RGB(40, 50, 60)
	ink := ui.RGB(230, 240, 250)
	muted := ui.RGB(120, 130, 140)
	accent := ui.RGB(70, 95, 208)
	hover := ui.RGB(30, 40, 55)
	overlayFace := ui.RGB(50, 60, 70)
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &theme.PalettePatch{
		Background: new(canvas), Surface: new(sheet), Foreground: new(ink), Muted: new(muted),
		Accent: new(accent), Hover: new(hover), Overlay: new(overlayFace),
		Border: new(ui.Transparent),
	}}})

	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Launch", OnClick: func() {}})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Outline", Variant: mygoui.ButtonOutline, OnClick: func() {}})
			list.ListBox(c, list.ListBoxOptions{
				Label: "Fruit", Selection: list.SelectionSingle, State: &state,
				Items: []list.ListBoxItem{
					{Key: "apple", Label: "Apple", Description: "Crisp"},
					{Key: "pear", Label: "Pear", Description: "Soft"},
				},
			})
		})
	}, 420, 320)
	state.List.Selected = new(0)
	tt.Frame()

	if got := theme.ColorsForMode(false).Muted; got != muted {
		t.Fatalf("list text palette muted = %v, want the injected muted ink", got)
	}

	// The primary action paints the injected accent.
	launch := findRect(t, tt, "Launch")
	if got := pixel(tt.Image(), launch.X+2, launch.Y+launch.H/2); !same(got, accent) {
		t.Errorf("primary button fill = %v, want the injected accent %v", got, accent)
	}

	// A transparent border leaves the outlined button's edge unpainted, so the
	// injected canvas shows through instead of HeroUI's border gray.
	outline := findRect(t, tt, "Outline")
	if got := pixel(tt.Image(), outline.X+2, outline.Y+outline.H/2); !same(got, canvas) {
		t.Errorf("outlined button edge = %v, want the canvas %v: a border was painted", got, canvas)
	}

	// A selected row keeps the native selection fill, which follows the
	// injected accent; an idle row shows the canvas and the injected hover fill.
	apple := findRect(t, tt, "Apple")
	if got := pixel(tt.Image(), apple.X+2, apple.Y+apple.H/2); !same(got, accent) {
		t.Errorf("selected row = %v, want the injected accent %v", got, accent)
	}
	pear := findRect(t, tt, "Pear")
	if got := pixel(tt.Image(), pear.X+2, pear.Y+pear.H/2); !same(got, canvas) {
		t.Errorf("idle row = %v, want the canvas %v", got, canvas)
	}
	tt.Move(pear.X+pear.W/2, pear.Y+pear.H/2)
	tt.Frame()
	if got := pixel(tt.Image(), pear.X+6, pear.Y+pear.H/2); !same(got, hover) {
		t.Errorf("hovered row = %v, want the injected hover fill %v", got, hover)
	}

	// The dialog panel paints its own overlay face, leaving the scrim and the
	// canvas outside it.
	open := true
	dt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		overlay.Dialog(c, &open, overlay.DialogOptions{
			Title: "Payload", Description: "Body copy",
			Body: func() { ui.Text(c, "content") },
		})
	}, 640, 480)
	dt.Frame()
	panel := findRect(t, dt, "Payload")
	if got := pixel(dt.Image(), panel.X+3, panel.Y+panel.H/2); !same(got, overlayFace) {
		t.Errorf("dialog panel = %v, want the injected overlay %v", got, overlayFace)
	}

	// The toast card paints the injected surface; its title keeps the ink.
	var toasts overlay.ToastState
	var toastCtx *ui.Context
	tt3 := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		toastCtx = c
		ui.Column(c).Fill().Background(c.Theme().Background).Children(func() {
			ui.Button(c, "Show").OnClick(func() {
				overlay.Toast(toastCtx, &toasts, overlay.ToastOptions{
					Title: "Saved", Action: "Undo", Timeout: -1,
				})
			})
			toasts.View(c)
		})
	}, 640, 480)
	tt3.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt3.Click("Show"); err != nil {
		t.Fatal(err)
	}
	tt3.Frame()
	if len(toasts.Items()) != 1 {
		t.Fatalf("toasts = %d, want 1", len(toasts.Items()))
	}
	saved := findRect(t, tt3, "Saved")
	if got := pixel(tt3.Image(), saved.X-8, saved.Y+saved.H/2); !same(got, sheet) {
		t.Errorf("toast card = %v, want the injected surface %v", got, sheet)
	}
}

func findRect(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("%q is not painted", label)
	}
	return r
}

func pixel(img *image.RGBA, x, y float32) ui.Color {
	c := img.RGBAAt(int(x), int(y))
	return ui.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

func same(got, want ui.Color) bool {
	near := func(a, b uint8) bool {
		if a > b {
			return a-b <= 2
		}
		return b-a <= 2
	}
	return near(got.R, want.R) && near(got.G, want.G) && near(got.B, want.B)
}
