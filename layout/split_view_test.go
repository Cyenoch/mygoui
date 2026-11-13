package layout_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/layout"
	"github.com/egoist/mygo/ui"
)

// shell is the caller-owned state of a split view under test.
type shell struct {
	size            float32
	collapsed, open bool
	note            string
	mode            layout.CompactMode
	fixed           bool
	max             float32
}

func (s *shell) view(c *ui.Context) {
	layout.SplitView(c, layout.SplitViewOptions{
		Compact: s.mode,
		Start: &layout.PaneOptions{
			Label: "Files", Size: &s.size, Width: 200, Min: 120, Max: s.max, Fixed: s.fixed,
			Collapsed: &s.collapsed, Open: &s.open,
			Body: func() { ui.Text(c, "readme.md") },
		},
		Content: func() {
			layout.PaneToggle(c, layout.PaneToggleOptions{
				Collapsed: &s.collapsed, Open: &s.open, Compact: s.mode, Label: "Toggle files",
			})
			form.Input(c, &s.note, form.InputOptions{AriaLabel: "Note"})
		},
	})
}

func splitter(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	a, ok := tt.AccessByRole(ui.RoleSplitter, name)
	if !ok {
		t.Fatalf("no splitter named %q", name)
	}
	return a.Bounds
}

func drag(t *testing.T, tt *ui.Tester, name string, dx float32) {
	t.Helper()
	r := splitter(t, tt, name)
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x+dx/2, y)
	tt.Move(x+dx, y)
	tt.Release(x+dx, y)
	// Keep the next press from reading as a double click.
	tt.Advance(time.Second)
}

func TestSplitViewDividerDragsWithinPaneAndContentBounds(t *testing.T) {
	s := &shell{max: 320}
	tt := newTester(s.view, 900, 400)
	if r := rect(t, tt, "readme.md"); s.size != 0 || !near(splitter(t, tt, "Files").X, 196) || r.X > 200 {
		t.Fatalf("at rest: size=%v divider=%+v, want the 200 DIP default without touching the caller's size", s.size, splitter(t, tt, "Files"))
	}
	drag(t, tt, "Files", 40)
	if s.size != 240 {
		t.Fatalf("after a 40 DIP drag size = %v, want 240", s.size)
	}
	drag(t, tt, "Files", 300)
	if s.size != 320 {
		t.Fatalf("past the maximum size = %v, want 320", s.size)
	}
	drag(t, tt, "Files", -230)
	if s.size != 120 || s.collapsed {
		t.Fatalf("towards the minimum: size=%v collapsed=%v, want 120 and shown", s.size, s.collapsed)
	}
	// Without a maximum the content still keeps its 240 DIP minimum.
	s.max = 0
	tt.Frame()
	drag(t, tt, "Files", 800)
	if s.size != 900-240 {
		t.Fatalf("with no maximum size = %v, want the view less the content minimum (660)", s.size)
	}
}

func TestSplitViewDividerKeysResizeCollapseAndReturn(t *testing.T) {
	s := &shell{size: 200, max: 320}
	tt := newTester(s.view, 900, 400)
	tt.Key(0, ui.KeyTab)
	if a, ok := tt.AccessByRole(ui.RoleSplitter, "Files"); !ok || !a.Focused || a.Now != 200 {
		t.Fatalf("Tab did not reach the divider named after its pane: %+v", a)
	}
	tt.Key(0, ui.KeyRight)
	if s.size != 216 {
		t.Fatalf("Right: size = %v, want 216", s.size)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyRight)
	if s.size != 320 {
		t.Fatalf("End then Right: size = %v, want the 320 maximum", s.size)
	}
	tt.Key(0, ui.KeyEnter)
	if !s.collapsed || tt.HasText("readme.md") {
		t.Fatalf("Enter: collapsed=%v, body shown=%v", s.collapsed, tt.HasText("readme.md"))
	}
	// The collapsed pane is still one Tab away, and comes back as it was.
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyTab)
	if a, _ := tt.AccessByRole(ui.RoleSplitter, "Files"); !a.Focused {
		t.Fatal("the divider of a collapsed pane is not a Tab stop")
	}
	tt.Key(0, ui.KeyEnter)
	if s.collapsed || s.size != 320 || !tt.HasText("readme.md") {
		t.Fatalf("Enter on a collapsed pane: collapsed=%v size=%v", s.collapsed, s.size)
	}
}

func TestSplitViewMirrorsInRightToLeft(t *testing.T) {
	s := &shell{size: 200, max: 320}
	tt := newTester(s.view, 900, 400)
	tt.SetDirection(ui.RightToLeft)
	if r := rect(t, tt, "readme.md"); r.X < 700 {
		t.Fatalf("start pane body = %+v, want it on the right edge", r)
	}
	drag(t, tt, "Files", -40)
	if s.size != 240 {
		t.Fatalf("dragging the divider left: size = %v, want 240", s.size)
	}
	tt.SetSize(500, 400)
	s.open = true
	tt.Frame()
	if a, ok := tt.AccessByRole(ui.RoleDialog, "Files"); !ok || !near(a.Bounds.X+a.Bounds.W, 500) {
		t.Fatalf("compact start pane = %+v (found %v), want its sheet on the right edge", a.Bounds, ok)
	}
}

func TestSplitViewCompactShowsThePaneOverTheContent(t *testing.T) {
	s := &shell{size: 200}
	tt := newTester(s.view, 900, 400)
	if err := tt.Click("Note"); err != nil {
		t.Fatal(err)
	}
	tt.Type("draft")
	tt.SetSize(500, 400)
	if tt.HasText("readme.md") {
		t.Fatal("a compact view still docks its pane")
	}
	tt.Type("ed")
	if !tt.Focused("Note") || s.note != "drafted" {
		t.Fatalf("content lost focus or text across the mode change: %q", s.note)
	}
	if note := rect(t, tt, "Note"); note.X > 100 {
		t.Fatalf("content = %+v, want it at the window's start edge", note)
	}

	s.open = true
	tt.Frame()
	body := rect(t, tt, "readme.md")
	if a, ok := tt.AccessByRole(ui.RoleDialog, "Files"); !ok || a.Bounds.X != 0 || a.Bounds.H != 400 || a.Bounds.W != 200 {
		t.Fatalf("open pane = %+v (found %v), want a 200 DIP sheet on the start edge at full height", a.Bounds, ok)
	}
	if body.X > 4 || body.Y > 4 {
		t.Fatalf("pane body = %+v, want it flush in the sheet", body)
	}
	tt.Key(0, ui.KeyEscape)
	if s.open || s.collapsed {
		t.Fatalf("Escape: open=%v collapsed=%v, want the sheet closed and the docked flag untouched", s.open, s.collapsed)
	}

	// Widening docks the pane again and forgets the open sheet.
	s.open = true
	tt.Frame()
	tt.SetSize(900, 400)
	if s.open || !tt.HasText("readme.md") {
		t.Fatalf("widened: open=%v, docked=%v", s.open, tt.HasText("readme.md"))
	}
	if _, ok := tt.AccessByRole(ui.RoleDialog, "Files"); ok {
		t.Fatal("the sheet outlived compact mode")
	}
}

func TestSplitViewCompactModeCanBeForced(t *testing.T) {
	s := &shell{mode: layout.CompactAlways}
	tt := newTester(s.view, 1400, 400)
	if tt.HasText("readme.md") {
		t.Fatal("CompactAlways docked the pane in a wide window")
	}
	s.mode = layout.CompactNever
	tt.SetSize(400, 400)
	if !tt.HasText("readme.md") {
		t.Fatal("CompactNever did not dock the pane in a narrow window")
	}
}

func TestPaneToggleFlipsTheFlagOfTheCurrentMode(t *testing.T) {
	s := &shell{size: 200}
	tt := newTester(s.view, 900, 400)
	checked := func() bool {
		a, ok := tt.AccessByRole(ui.RoleToggleButton, "Toggle files")
		if !ok {
			t.Fatal("no toggle button named Toggle files")
		}
		return a.Checked
	}
	if !checked() {
		t.Fatal("the toggle is not checked while its pane is docked")
	}
	if err := tt.Click("Toggle files"); err != nil {
		t.Fatal(err)
	}
	if !s.collapsed || s.open || checked() {
		t.Fatalf("docked toggle: collapsed=%v open=%v checked=%v", s.collapsed, s.open, checked())
	}
	// The toggle keeps the focus, so the keyboard brings the pane back.
	tt.Key(0, ui.KeySpace)
	if s.collapsed {
		t.Fatal("Space on the toggle did not expand the pane")
	}

	tt.SetSize(500, 400)
	if checked() {
		t.Fatal("the toggle is checked in a compact view whose pane is closed")
	}
	if err := tt.Click("Toggle files"); err != nil {
		t.Fatal(err)
	}
	if !s.open || s.collapsed || !tt.HasText("readme.md") {
		t.Fatalf("compact toggle: open=%v collapsed=%v", s.open, s.collapsed)
	}
	tt.Key(0, ui.KeyEscape)
	if s.open || !tt.Focused("Toggle files") {
		t.Fatalf("closing the sheet: open=%v, focus returned=%v", s.open, tt.Focused("Toggle files"))
	}
}

func TestFixedPaneHasNoDividerAndFoldsOverTime(t *testing.T) {
	s := &shell{fixed: true}
	tt := newTester(s.view, 900, 400)
	if _, ok := tt.AccessByRole(ui.RoleSplitter, ""); ok {
		t.Fatal("a fixed pane offers a divider")
	}
	note := rect(t, tt, "Note")
	if !near(note.X, 200) {
		t.Fatalf("content starts at %v, want after the 200 DIP pane", note.X)
	}

	tt.SetPreferences(ui.Preferences{TextScale: 1})
	s.collapsed = true
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if x := rect(t, tt, "Note").X; x <= 0 || x >= 200 {
		t.Fatalf("content at %v halfway through the fold, want between 0 and 200", x)
	}
	tt.Advance(300 * time.Millisecond)
	if x := rect(t, tt, "Note").X; x != 0 || tt.HasText("readme.md") {
		t.Fatalf("folded: content at %v, body shown=%v", x, tt.HasText("readme.md"))
	}

	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	s.collapsed = false
	tt.Frame()
	if x := rect(t, tt, "Note").X; !near(x, 200) {
		t.Fatalf("reduced motion: content at %v, want the pane back at once", x)
	}
}
