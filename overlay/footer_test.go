package overlay

import (
	"fmt"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestModalAndBottomSheetUseFullBleedFooter(t *testing.T) {
	for _, kind := range []string{"modal", "bottomsheet"} {
		for _, width := range []int{360, 800} {
			t.Run(fmt.Sprintf("%s/%d", kind, width), func(t *testing.T) {
				open := true
				var panel ui.Rect
				var pal theme.Palette
				tt := ui.NewTester(func(c *ui.Context) {
					theme.ApplyMode(c, false)
					pal = theme.Colors(c)
					footer := func() {
						mygoui.Button(c, mygoui.ButtonOptions{Label: "Cancel", Variant: mygoui.ButtonOutline, FullWidth: width < 640})
						mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", FullWidth: width < 640})
					}
					var p ui.Element
					if kind == "modal" {
						p = Modal(c, &open, ModalOptions{Title: "Panel", Footer: footer})
					} else {
						p = BottomSheet(c, &open, BottomSheetOptions{Title: "Panel", Footer: footer})
					}
					if p.Valid() {
						panel = p.LayoutBounds()
					}
				}, width, 600)
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				tt.Frame()
				tt.Frame()
				cancel, _ := tt.Find("Cancel")
				save, _ := tt.Find("Save")
				if width < 640 && save.Y+save.H > cancel.Y {
					t.Fatalf("compact footer did not stack primary above cancel: cancel=%+v save=%+v", cancel, save)
				}
				if width < 640 && (save.W != cancel.W || abs32(save.W-(panel.W-50)) > 1) {
					t.Fatalf("explicit full-width footer actions do not fill: panel=%+v save=%+v cancel=%+v", panel, save, cancel)
				}
				if width >= 640 && (save.Y != cancel.Y || save.X <= cancel.X) {
					t.Fatalf("desktop footer did not end-align actions: cancel=%+v save=%+v", cancel, save)
				}
				x, y := int(panel.X)+3, int(panel.Y+panel.H)-12
				got := tt.Image().RGBAAt(x, y)
				want := pal.MutedBackground.Alpha(0.72).Over(pal.Popover)
				if abs32(float32(got.R)-float32(want.R)) > 2 {
					t.Fatalf("footer does not reach panel edge at (%d,%d): %v, want %v", x, y, got, want)
				}
				img := tt.Image()
				for _, x := range []int{int(panel.X), int(panel.X+panel.W) - 1} {
					if got := img.RGBAAt(x, int(panel.Y+panel.H)-1); sameColor(got, want, 2) {
						t.Fatalf("footer face covers the rounded bottom corner: %v", got)
					}
				}
				if tt.FramePending() {
					t.Fatal("settled panel continues scheduling layout frames")
				}
			})
		}
	}
}

func TestShortPanelsKeepFooterActionsReachable(t *testing.T) {
	for _, kind := range []string{"modal", "bottomsheet", "dialog", "drawer", "sheet", "alert"} {
		t.Run(kind, func(t *testing.T) {
			open := false
			calls := 0
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				body := func() { ui.Text(c, "Changes are not saved yet.") }
				footer := func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Cancel", OnClick: func() { calls++ }})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: func() { calls++ }})
				}
				var p ui.Element
				switch kind {
				case "modal":
					p = Modal(c, &open, ModalOptions{Title: "Panel", Body: body, Footer: footer, PreventDismiss: true, HideClose: true})
				case "bottomsheet":
					p = BottomSheet(c, &open, BottomSheetOptions{Title: "Panel", Body: body, Footer: footer, PreventDismiss: true, HideClose: true})
				case "dialog":
					p = Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Footer: footer, PreventDismiss: true, HideClose: true})
				case "drawer":
					p = Drawer(c, &open, DrawerOptions{Title: "Panel", Body: body, Footer: footer, PreventDismiss: true})
				case "sheet":
					p = Sheet(c, &open, SheetOptions{Title: "Panel", Body: body, Footer: footer, PreventDismiss: true, HideClose: true})
				case "alert":
					p = AlertDialog(c, &open, AlertDialogOptions{Title: "Panel", Body: body, CancelLabel: "Cancel", ActionLabel: "Save", OnAction: func() { calls++ }})
				}
				if p.Valid() {
					panel = p.LayoutBounds()
				}
			}, 360, 200)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 2})
			open = true
			tt.Advance(time.Second)
			for range 4 {
				tt.Frame()
			}
			if kind == "alert" {
				cancel, _ := tt.Find("Cancel")
				if cancel.H == 0 || cancel.Y < panel.Y || cancel.Y+cancel.H > panel.Y+panel.H {
					t.Fatalf("initially focused Cancel is not revealed: cancel=%+v panel=%+v", cancel, panel)
				}
			}
			for _, action := range []string{"Cancel", "Save"} {
				for range 8 {
					tt.Key(0, ui.KeyTab)
					if tt.Focused(action) {
						break
					}
				}
				r, ok := tt.Find(action)
				if !ok || !tt.Focused(action) || r.Y < panel.Y || r.Y+r.H > panel.Y+panel.H {
					t.Fatalf("focused %s remains clipped: focused=%v action=%+v panel=%+v", action, tt.Focused(action), r, panel)
				}
				// Cancel closes an alert; Save independently verifies its action below.
				if kind != "alert" || action == "Save" {
					tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
				}
			}
			want := 2
			if kind == "alert" {
				want = 1
			}
			if calls != want {
				t.Fatalf("visible actions did not activate: %d, want %d", calls, want)
			}
		})
	}
}
