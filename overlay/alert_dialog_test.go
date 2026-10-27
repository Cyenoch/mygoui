package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAlertDialogOptionalIconAndClose(t *testing.T) {
	for _, tc := range []struct {
		name      string
		options   AlertDialogOptions
		wantIcon  bool
		wantClose bool
	}{
		{name: "default"},
		{name: "show icon", options: AlertDialogOptions{ShowIcon: true}, wantIcon: true},
		{name: "explicit icon", options: AlertDialogOptions{Icon: lucide.Info()}, wantIcon: true},
		{name: "hidden icon overrides both", options: AlertDialogOptions{ShowIcon: true, Icon: lucide.Info(), HideIcon: true}},
		{name: "show close", options: AlertDialogOptions{ShowClose: true}, wantClose: true},
		{name: "hidden close overrides show", options: AlertDialogOptions{ShowClose: true, HideClose: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			var panel ui.Rect
			cancelled := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				o := tc.options
				o.Title, o.Description, o.CancelLabel = "Confirm", "Details", "Cancel"
				o.OnCancel = func() { cancelled++ }
				if p := AlertDialog(c, &open, o); p.Valid() {
					panel = p.Bounds()
				}
			}, 800, 600)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			description, ok := tt.Find("Details")
			if !ok {
				t.Fatal("missing description")
			}
			if got := description.Y-panel.Y > 80; got != tc.wantIcon {
				t.Fatalf("icon present=%v, want %v (panel=%+v description=%+v)", got, tc.wantIcon, panel, description)
			}
			if _, got := tt.Find("Close"); got != tc.wantClose {
				t.Fatalf("close present=%v, want %v", got, tc.wantClose)
			}
			if !tt.Focused("Cancel") {
				t.Fatal("Cancel did not receive initial focus")
			}
			cancel, _ := tt.Find("Cancel")
			pal := theme.ColorsForMode(false)
			want := pal.MutedBackground.Alpha(0.72).Over(pal.Popover)
			if got := tt.Image().RGBAAt(int(cancel.X+4), int(cancel.Y+cancel.H/2)); !sameColor(got, want, 3) {
				t.Fatalf("Cancel has a visible face: got %v, want footer %v", got, want)
			}
			if tc.wantClose {
				if err := tt.Click("Close"); err != nil {
					t.Fatal(err)
				}
				if open || cancelled != 1 || !tt.Focused("Open") {
					t.Fatalf("close: open=%v cancellations=%d opener focused=%v", open, cancelled, tt.Focused("Open"))
				}
			}
		})
	}
}

func TestAlertDialogWidthOptions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		options AlertDialogOptions
		width   int
		want    float32
	}{
		{"default", AlertDialogOptions{}, 800, 512},
		{"medium", AlertDialogOptions{Size: ModalMedium}, 800, 512},
		{"width override", AlertDialogOptions{Width: 448}, 800, 448},
		{"small", AlertDialogOptions{Size: ModalSmall}, 800, 384},
		{"large", AlertDialogOptions{Size: ModalLarge}, 800, 512},
		{"mobile fills window", AlertDialogOptions{Width: 200}, 320, 320},
	} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				if p := AlertDialog(c, &open, tc.options); p.Valid() {
					panel = p.Bounds()
				}
			}, tc.width, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			if panel.W != tc.want {
				t.Fatalf("panel width=%v, want %v", panel.W, tc.want)
			}
		})
	}

	open := true
	var panel ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if p := Modal(c, &open, ModalOptions{Size: ModalMedium}); p.Valid() {
			panel = p.Bounds()
		}
	}, 800, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if panel.W != 448 {
		t.Fatalf("ModalMedium width=%v, want 448", panel.W)
	}
}
