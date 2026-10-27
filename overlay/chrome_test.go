package overlay

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
)

func TestModalBackdrops(t *testing.T) {
	const width, height, edge, row = 320, 240, 48, 16
	background := func(c *ui.Context) {
		ui.Box(c).Fill().Background(ui.RGB(255, 255, 255)).Children(func() {
			ui.Box(c).Absolute().Left(0).Top(0).Width(edge).Height(height).Background(ui.RGB(0, 0, 0))
		})
	}
	builders := []struct {
		name  string
		build func(*ui.Context, *bool, ModalBackdrop)
	}{
		{"dialog", func(c *ui.Context, open *bool, backdrop ModalBackdrop) {
			Dialog(c, open, DialogOptions{Backdrop: backdrop, HideClose: true})
		}},
		{"alert", func(c *ui.Context, open *bool, backdrop ModalBackdrop) {
			AlertDialog(c, open, AlertDialogOptions{Backdrop: backdrop, HideClose: true, HideIcon: true})
		}},
		{"sheet", func(c *ui.Context, open *bool, backdrop ModalBackdrop) {
			Sheet(c, open, SheetOptions{Backdrop: backdrop, Width: 96, HideClose: true})
		}},
		{"drawer", func(c *ui.Context, open *bool, backdrop ModalBackdrop) {
			Drawer(c, open, DrawerOptions{Side: DrawerRight, Backdrop: backdrop, Width: 96, HideClose: true})
		}},
		{"modal", func(c *ui.Context, open *bool, backdrop ModalBackdrop) {
			Modal(c, open, ModalOptions{Backdrop: backdrop, HideClose: true})
		}},
		{"bottom sheet", func(c *ui.Context, open *bool, backdrop ModalBackdrop) {
			BottomSheet(c, open, BottomSheetOptions{
				Backdrop: backdrop, Height: 80, HideHandle: true, HideClose: true,
			})
		}},
	}
	for _, dark := range []bool{false, true} {
		name := "light"
		if dark {
			name = "dark"
		}
		t.Run(name, func(t *testing.T) {
			reference := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				background(c)
				ui.Overlay(c, func() {
					ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
						Material(glass.Blur{Radius: 4}).
						Draw(func(p *ui.Painter, r ui.Rect) {
							p.Fill(r, ui.RGB(0, 0, 0).Alpha(0.32), 0)
						})
				})
			}, width, height)
			reference.SetDark(dark)
			want := image.NewRGBA(image.Rect(0, 0, width, height))
			copy(want.Pix, reference.Image().Pix)
			if want.RGBAAt(edge-2, row).R == 0 {
				t.Fatal("reference did not blur across the edge")
			}
			for _, builder := range builders {
				t.Run(builder.name, func(t *testing.T) {
					open := false
					tt := ui.NewTester(func(c *ui.Context) {
						theme.Apply(c)
						background(c)
						builder.build(c, &open, ModalBackdrop(0))
					}, width, height)
					tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
					tt.SetDark(dark)
					open = true
					tt.Frame()
					got := tt.Image()
					for x := edge - 8; x < edge+8; x++ {
						if a, b := got.RGBAAt(x, row), want.RGBAAt(x, row); a != b {
							t.Fatalf("default backdrop at x=%d: got %v, want 4 DIP blur %v", x, a, b)
						}
					}
				})
			}
			for _, tc := range []struct {
				name     string
				backdrop ModalBackdrop
				white    uint8
			}{
				{"opaque", ModalBackdropOpaque, 173},
				{"transparent", ModalBackdropTransparent, 255},
			} {
				t.Run(tc.name, func(t *testing.T) {
					open := false
					tt := ui.NewTester(func(c *ui.Context) {
						theme.Apply(c)
						background(c)
						Dialog(c, &open, DialogOptions{Backdrop: tc.backdrop, HideClose: true})
					}, width, height)
					tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
					tt.SetDark(dark)
					open = true
					tt.Frame()
					img := tt.Image()
					if left, right := img.RGBAAt(edge-1, row).R, img.RGBAAt(edge, row).R; left != 0 || right != tc.white {
						t.Fatalf("backdrop edge: left=%d right=%d, want 0 and %d", left, right, tc.white)
					}
				})
			}
		})
	}
}

func TestCOSSFooterFullBleed(t *testing.T) {
	for _, kind := range []string{"dialog", "alert", "sheet", "drawer"} {
		for _, dark := range []bool{false, true} {
			name := kind + "/light"
			if dark {
				name = kind + "/dark"
			}
			t.Run(name, func(t *testing.T) {
				open := true
				var panel ui.Rect
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					footer := func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Continue"}) }
					var p ui.Element
					switch kind {
					case "dialog":
						p = Dialog(c, &open, DialogOptions{Title: "Panel", Footer: footer})
					case "alert":
						p = AlertDialog(c, &open, AlertDialogOptions{
							Title: "Panel", CancelLabel: "Cancel", ActionLabel: "Continue",
							OnAction: func() {},
						})
					case "sheet":
						p = Sheet(c, &open, SheetOptions{Title: "Panel", Footer: footer})
					case "drawer":
						p = Drawer(c, &open, DrawerOptions{Side: DrawerRight, Title: "Panel", Footer: footer})
					}
					if p.Valid() {
						panel = p.Bounds()
					}
				}, 800, 600)
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				tt.SetDark(dark)
				tt.Frame()
				action, ok := tt.Find("Continue")
				if !ok {
					t.Fatal("footer action is missing")
				}
				right := panel.X + panel.W - (action.X + action.W)
				bottom := panel.Y + panel.H - (action.Y + action.H)
				if abs32(right-25) > 2 || abs32(bottom-17) > 2 {
					t.Fatalf("footer action padding: right=%v bottom=%v panel=%+v action=%+v",
						right, bottom, panel, action)
				}
				pal := theme.ColorsForMode(dark)
				want := pal.MutedBackground.Alpha(0.72).Over(pal.Popover)
				img := tt.Image()
				y := int(action.Y - 8)
				for _, x := range []int{int(panel.X + 4), int(panel.X + panel.W - 5)} {
					if got := img.RGBAAt(x, y); !sameColor(got, want, 3) {
						t.Fatalf("footer is not full bleed at (%d,%d): got %v want %v", x, y, got, want)
					}
				}
				x := int(panel.X + 30)
				above := img.RGBAAt(x, int(action.Y-18))
				border := img.RGBAAt(x, int(action.Y-17))
				fill := img.RGBAAt(x, y)
				if border == above || border == fill {
					t.Fatalf("footer top border is missing: above=%v border=%v fill=%v", above, border, fill)
				}
			})
		}
	}
}

func TestCOSSFooterNarrowScrollAndActions(t *testing.T) {
	for _, kind := range []string{"dialog", "alert", "sheet", "drawer"} {
		t.Run(kind, func(t *testing.T) {
			var open bool
			var panel ui.Rect
			bodyActions, footerActions := 0, 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				body := func() {
					for range 25 {
						ui.Text(c, "Scrollable row")
					}
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Body action", OnClick: func() { bodyActions++ }})
				}
				footer := func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Cancel", Variant: mygoui.ButtonGhost})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: func() { footerActions++ }})
				}
				var p ui.Element
				switch kind {
				case "dialog":
					p = Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Footer: footer})
				case "alert":
					p = AlertDialog(c, &open, AlertDialogOptions{
						Title: "Panel", Body: body, CancelLabel: "Cancel", ActionLabel: "Save",
						OnAction: func() { footerActions++ },
					})
				case "sheet":
					p = Sheet(c, &open, SheetOptions{Title: "Panel", Body: body, Footer: footer})
				case "drawer":
					p = Drawer(c, &open, DrawerOptions{Side: DrawerRight, Title: "Panel", Body: body, Footer: footer})
				}
				if p.Valid() {
					panel = p.Bounds()
				}
			}, 320, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			save, ok := tt.Find("Save")
			if !ok {
				t.Fatal("missing Save")
			}
			cancel, ok := tt.Find("Cancel")
			if !ok {
				t.Fatal("missing Cancel")
			}
			if save.Y+save.H > cancel.Y || cancel.Y+cancel.H > panel.Y+panel.H-15 ||
				save.X < panel.X+23 || save.X+save.W > panel.X+panel.W-23 ||
				cancel.X < panel.X+23 || cancel.X+cancel.W > panel.X+panel.W-23 {
				t.Fatalf("stacked footer does not fit: panel=%+v save=%+v cancel=%+v", panel, save, cancel)
			}
			if kind == "alert" && (abs32(save.W-(panel.W-48)) > 2 || abs32(cancel.W-(panel.W-48)) > 2) {
				t.Fatalf("alert actions do not fill the footer: panel=%+v save=%+v cancel=%+v", panel, save, cancel)
			}

			tt.Scroll(panel.X+80, panel.Y+100, 0, 1000)
			tt.Frame()
			bodyButton, ok := tt.Find("Body action")
			if !ok || bodyButton.Y+bodyButton.H > save.Y-16 {
				t.Fatalf("body did not scroll above the fixed footer: body=%+v footer=%+v", bodyButton, save)
			}
			tt.Key(0, ui.KeyTab)
			for i := 0; i < 8 && !tt.Focused("Body action"); i++ {
				tt.Key(0, ui.KeyTab)
			}
			if !tt.Focused("Body action") {
				t.Fatal("Tab did not reach the body action")
			}
			img := tt.Image()
			ring := img.RGBAAt(int(bodyButton.X)-2, int(bodyButton.Y+bodyButton.H/2))
			if !sameColor(ring, theme.ColorsForMode(false).Focus, 12) {
				t.Fatalf("body focus ring is clipped or missing: %v", ring)
			}
			if err := tt.Click("Body action"); err != nil {
				t.Fatal(err)
			}
			if bodyActions != 1 || !open {
				t.Fatalf("body action: open=%v calls=%d", open, bodyActions)
			}
			if err := tt.Click("Save"); err != nil {
				t.Fatal(err)
			}
			if footerActions != 1 {
				t.Fatalf("footer action calls=%d, want 1", footerActions)
			}
		})
	}
}
