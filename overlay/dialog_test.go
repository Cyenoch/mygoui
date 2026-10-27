package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDialogAndAlertViewportPlacement(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		for _, tc := range []struct {
			name   string
			width  int
			height int
			mobile bool
		}{
			{"desktop", 800, 600, false},
			{"desktop breakpoint", 640, 480, false},
			{"mobile breakpoint", 639, 480, true},
			{"narrow", 320, 360, true},
		} {
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				open := true
				var panel ui.Rect
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					var p ui.Element
					if kind == "alert" {
						p = AlertDialog(c, &open, AlertDialogOptions{
							Title: "Confirmation", CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
						})
					} else {
						p = Dialog(c, &open, DialogOptions{
							Title: "Confirmation", Backdrop: ModalBackdropOpaque,
						})
					}
					if p.Valid() {
						panel = p.Bounds()
					}
				}, tc.width, tc.height)
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				tt.Frame()
				if panel.H <= 0 {
					t.Fatalf("the panel has no height: %+v", panel)
				}
				if tc.mobile {
					if panel.X != 0 || panel.W != float32(tc.width) ||
						panel.Y+panel.H != float32(tc.height) || panel.Y < 48 {
						t.Fatalf("mobile panel=%+v, want full width, bottom attached, and a 48 DIP top gap", panel)
					}
				} else {
					wantY := float32(16) + (float32(tc.height)-32-panel.H)/4
					wantX := (float32(tc.width) - 512) / 2
					if panel.W != 512 || panel.X != wantX || abs32(panel.Y-wantY) > 1 {
						t.Fatalf("desktop panel=%+v, want 512 DIP width at (%v,%v)", panel, wantX, wantY)
					}
				}
			})
		}
	}
}

func TestDialogDesktopTallBodyScrolls(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		t.Run(kind, func(t *testing.T) {
			open := true
			var panel ui.Rect
			clicked := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				body := func() {
					for range 40 {
						ui.Text(c, "Scrollable row")
					}
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Last action", OnClick: func() { clicked++ }})
				}
				var p ui.Element
				if kind == "alert" {
					p = AlertDialog(c, &open, AlertDialogOptions{Title: "Panel", Body: body, CancelLabel: "Cancel"})
				} else {
					p = Dialog(c, &open, DialogOptions{
						Title: "Panel", Body: body,
						Footer: func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Cancel"}) },
					})
				}
				if p.Valid() {
					panel = p.Bounds()
				}
			}, 800, 600)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			if panel.Y < 15 || panel.Y+panel.H > 585 {
				t.Fatalf("tall panel exceeds desktop viewport: %+v", panel)
			}
			footer, ok := tt.Find("Cancel")
			if !ok {
				t.Fatal("missing footer action")
			}
			tt.Scroll(panel.X+80, panel.Y+100, 0, 2000)
			tt.Frame()
			last, ok := tt.Find("Last action")
			if !ok || last.Y+last.H > footer.Y-16 {
				t.Fatalf("last body action did not scroll above footer: body=%+v footer=%+v", last, footer)
			}
			if err := tt.Click("Last action"); err != nil {
				t.Fatal(err)
			}
			if clicked != 1 {
				t.Fatalf("body action calls=%d, want 1", clicked)
			}
		})
	}
}

func TestAlertMobileHeaderTextIsCentered(t *testing.T) {
	for _, tc := range []struct {
		name   string
		alert  bool
		width  int
		center bool
	}{
		{"alert mobile", true, 320, true},
		{"alert desktop", true, 800, false},
		{"dialog mobile", false, 320, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				var p ui.Element
				if tc.alert {
					p = AlertDialog(c, &open, AlertDialogOptions{
						Title: "Heads up", Description: "Review the details", CancelLabel: "Cancel",
						Backdrop: ModalBackdropOpaque,
					})
				} else {
					p = Dialog(c, &open, DialogOptions{
						Title: "Heads up", Description: "Review the details", HideClose: true,
						Backdrop: ModalBackdropOpaque,
					})
				}
				if p.Valid() {
					panel = p.Bounds()
				}
			}, tc.width, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			img := tt.Image()
			for _, label := range []string{"Heads up", "Review the details"} {
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing %q", label)
				}
				if label == "Heads up" {
					r.Y = panel.Y + theme.Space24
					r.H = theme.H4().LineHeight
				}
				minX, maxX := int(panel.X+panel.W), -1
				for y := int(r.Y); y < int(r.Y+r.H); y++ {
					for x := int(panel.X) + 20; x < int(panel.X+panel.W)-20; x++ {
						pixel := img.RGBAAt(x, y)
						if pixel.R < 190 && pixel.G < 190 && pixel.B < 190 {
							minX, maxX = min(minX, x), max(maxX, x)
						}
					}
				}
				if maxX < 0 {
					t.Fatalf("no ink for %q", label)
				}
				inkCenter := float32(minX+maxX) / 2
				panelCenter := panel.X + panel.W/2
				if tc.center && abs32(inkCenter-panelCenter) > 5 {
					t.Fatalf("%q ink x=%v..%v rect=%+v centered at %v, want panel center %v", label, minX, maxX, r, inkCenter, panelCenter)
				}
				if !tc.center && inkCenter >= panelCenter-20 {
					t.Fatalf("%q ink x=%v..%v rect=%+v centered at %v, want start-aligned in panel %+v", label, minX, maxX, r, inkCenter, panel)
				}
			}
			if tc.alert {
				description, _ := tt.Find("Review the details")
				cancel, _ := tt.Find("Cancel")
				gap := cancel.Y - theme.Space16 - theme.Hairline - (description.Y + description.H)
				want := float32(theme.Space24)
				if tc.width < 640 {
					want = theme.Space16
				}
				if abs32(gap-want) > 1 {
					t.Fatalf("header/footer gap=%v, want %v", gap, want)
				}
			}
		})
	}
}
