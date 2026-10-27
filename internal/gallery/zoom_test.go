package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestZoomKeepsNativeChromeAndSearchUsable(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80, Right: 138})
	for range 5 {
		tt.Key(ui.Cmd, ui.KeyEqual)
	}
	title, ok := tt.Find("Gallery title bar")
	if !ok || title.H != TitleBarHeight {
		t.Fatalf("zoom changed the native chrome height: %+v", title)
	}
	for _, label := range []string{"coss ui", "MyGo"} {
		r, found := tt.Find(label)
		if found && r.Y < title.H && (r.X < 80 || r.X+r.W > 1180-138) {
			t.Fatalf("zoomed %s overlaps native controls: %+v", label, r)
		}
	}
	if _, ok := tt.Find("Toggle navigation"); !ok {
		t.Fatal("zoom did not switch to the compact layout")
	}
	tt.Key(ui.Cmd, ui.KeyK)
	if !tt.Focused("Search components") {
		t.Fatal("zoomed search shortcut did not focus the editor")
	}
	tt.Type("composer")
	tt.Key(0, ui.KeyEnter)
	if s.Selected != "composer" {
		t.Fatalf("zoomed input did not navigate: %s", s.Selected)
	}
	tt.Key(ui.Cmd, ui.Key0)
	if s.zoomStep != 0 {
		t.Fatal("reset zoom did not restore the default")
	}
	if _, ok := tt.Find("Search components"); !ok {
		t.Fatal("reset zoom did not restore the desktop sidebar")
	}
}

func TestZoomUsesWindowOwnedStateAndStopsAtSupportedLevels(t *testing.T) {
	s := New()
	var factor float32
	tt := ui.NewTester(func(c *ui.Context) {
		zoom(c, s)
		factor = c.Zoom()
	}, 200, 100)
	for range 20 {
		tt.Key(ui.Cmd, ui.KeyEqual)
	}
	if factor != 3 {
		t.Fatalf("maximum zoom = %g", factor)
	}
	other := New()
	ui.NewTester(func(c *ui.Context) { zoom(c, other) }, 200, 100)
	if other.zoomStep != 0 {
		t.Fatal("one window inherited another window's zoom")
	}
	for range 30 {
		tt.Key(ui.Cmd, ui.KeyMinus)
	}
	if factor != 0.5 {
		t.Fatalf("minimum zoom = %g", factor)
	}
}

func TestAppearanceNotifiesHostOnlyWhenChanged(t *testing.T) {
	s := New()
	var modes []bool
	s.OnAppearanceChange = func(dark bool) { modes = append(modes, dark) }
	s.setDark(false)
	s.setDark(true)
	s.setDark(true)
	s.setDark(false)
	if len(modes) != 2 || !modes[0] || modes[1] {
		t.Fatalf("appearance changes = %v", modes)
	}
}
