package layout_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// The title bar is the window's banner landmark; a scaffold's start band
// beside it stays a group, so the window has one banner.
func TestTitleBarIsTheWindowsBanner(t *testing.T) {
	var clicks int
	tt := newTester(barWith(layout.SpanWindow, &clicks), 640, 300)
	if _, ok := tt.AccessByRole(ui.RoleBanner, "Window bar"); !ok {
		t.Fatal("the title bar is not a banner named by its label")
	}
	a := &app{}
	tt = newTester(a.view, 1000, 600)
	banners := 0
	for _, n := range tt.AccessAll() {
		if n.Role == ui.RoleBanner {
			banners++
		}
	}
	if _, ok := tt.AccessByRole(ui.RoleGroup, "Mailboxes bar"); !ok || banners != 1 {
		t.Fatalf("scaffold has %d banners and start band group %v, want one banner and the band a group", banners, ok)
	}
}

// Without a Max the divider's range reaches past its value from the first
// frame: the room the window leaves the content.
func TestSplitViewDividerRangeBeforeTheFirstLayout(t *testing.T) {
	s := &shell{}
	tt := ui.NewTester(s.view, 900, 400)
	a, ok := tt.AccessByRole(ui.RoleSplitter, "Files")
	if !ok || a.Max <= a.Now {
		t.Fatalf("splitter = %+v, want a maximum above the 200 DIP width", a)
	}
}

// A scaffold's resizable sidebar draws its hairline in the sidebar's own
// border color, as a fixed one does.
func TestScaffoldSidebarHairlineUsesTheSidebarBorder(t *testing.T) {
	a := &app{}
	var p theme.Palette
	tt := newTester(func(c *ui.Context) { p = theme.Colors(c); a.view(c) }, 1000, 600)
	want := p.SidebarBorder.Over(p.Sidebar)
	img := tt.Image()
	for x := 210; x < 230; x++ {
		o := img.PixOffset(x, 300)
		if px := img.Pix[o : o+3]; px[0] == want.R && px[1] == want.G && px[2] == want.B {
			return
		}
	}
	t.Fatalf("no hairline of the sidebar border %v near the pane's edge", want)
}
