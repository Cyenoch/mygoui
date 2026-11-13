package navigation

import (
	"fmt"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// navStackFixture is a caller owning a two-level stack: a list page whose
// button goes deeper, and a detail page.
type navStackFixture struct {
	pages       []NavigationPage
	keepMounted bool
	pops        int
}

func (f *navStackFixture) view(c *ui.Context) {
	theme.Apply(c)
	if f.pages == nil {
		f.pages = []NavigationPage{{Key: "list", Title: "Settings", Build: func() {
			owner := ui.Column(c)
			count := ui.Local(owner, "count", func() int { return 0 })
			owner.Children(func() {
				ui.Button(c, fmt.Sprintf("Count %d", *count)).OnClick(func() { *count++ })
				ui.Button(c, "Open network").OnClick(func() {
					f.pages = append(f.pages, NavigationPage{Key: "network", Title: "Network", Build: func() {
						ui.Button(c, "Reset adapter")
					}})
				})
			})
		}}}
	}
	ui.Column(c).Fill().Children(func() {
		NavigationStack(c, NavigationStackOptions{
			Label: "Preferences", Pages: f.pages, KeepMounted: f.keepMounted,
			OnPop: func() {
				f.pops++
				f.pages = f.pages[:len(f.pages)-1]
			},
		}).Grow(1)
	})
}

func TestNavigationStackPushesAndPopsWithFocus(t *testing.T) {
	f := &navStackFixture{}
	tt := ui.NewTester(f.view, 360, 240)
	if _, ok := tt.Find("Back"); ok {
		t.Fatal("the root page offers a way back")
	}
	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if page, ok := tt.AccessByRole(ui.RoleGroup, "Network"); !ok || !page.Focused {
		t.Fatalf("the focus did not move into the pushed page: found %v, %+v", ok, page)
	}
	if _, ok := tt.AccessByRole(ui.RoleHeading, "Network"); !ok || tt.HasText("Open network") {
		t.Fatalf("after the push the stack shows %q", tt.Texts())
	}
	tt.Key(ui.Alt, ui.KeyLeft)
	tt.Advance(time.Second)
	if page, ok := tt.AccessByRole(ui.RoleGroup, "Settings"); f.pops != 1 || !ok || !page.Focused {
		t.Fatalf("Alt+Left: pops %d, focus returned %v (%+v)", f.pops, ok, page)
	}
	tt.Key(ui.Alt, ui.KeyLeft)
	if f.pops != 1 {
		t.Fatal("Alt+Left popped the root page")
	}

	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	tt.Key(ui.Alt, ui.KeyLeft)
	tt.Advance(time.Second)
	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if err := tt.Click("Back"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if f.pops != 3 || len(f.pages) != 1 || tt.HasText("Reset adapter") {
		t.Fatalf("pops %d, depth %d, texts %q", f.pops, len(f.pages), tt.Texts())
	}

	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	tt.Key(0, ui.KeyBack)
	tt.Advance(time.Second)
	if f.pops != 4 {
		t.Fatal("the mouse's back button did not go back")
	}

	// Laid out right to left, back is to the right.
	tt.SetDirection(ui.RightToLeft)
	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	tt.Key(ui.Alt, ui.KeyLeft)
	if f.pops != 4 {
		t.Fatal("Alt+Left went back in a right-to-left layout")
	}
	tt.Key(ui.Alt, ui.KeyRight)
	if f.pops != 5 {
		t.Fatal("Alt+Right did not go back in a right-to-left layout")
	}
}

func TestNavigationStackSlidesUnlessMotionIsReduced(t *testing.T) {
	f := &navStackFixture{}
	tt := ui.NewTester(f.view, 360, 240)
	rest, _ := tt.Find("Open network")
	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(48 * time.Millisecond)
	leaving, _ := tt.Find("Open network")
	if leaving.W <= 0 || leaving.X >= rest.X {
		t.Fatalf("mid push the covered page is at %+v, from %+v: it should be sliding away", leaving, rest)
	}
	// The covered page is still built, so no copy of it may leave toward
	// the side the new page comes from.
	img := tt.Image()
	for y := int(rest.Y); y < int(rest.Y+rest.H); y++ {
		for x := 300; x < 360; x++ {
			if r, _, _, _ := img.At(x, y).RGBA(); r < 0xf000 {
				t.Fatalf("a copy of the covered page paints at %d,%d", x, y)
			}
		}
	}
	tt.Advance(time.Second)
	if coming, _ := tt.Find("Reset adapter"); coming.W <= 0 || tt.HasText("Open network") || tt.FramePending() {
		t.Fatalf("the push settled at %+v, shows %q, pending %v", coming, tt.Texts(), tt.FramePending())
	}

	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Key(ui.Alt, ui.KeyLeft)
	if tt.HasText("Reset adapter") || !tt.HasText("Open network") {
		t.Fatalf("reduced motion did not change pages at once: %q", tt.Texts())
	}
}

func TestNavigationStackKeepMountedPreservesCoveredPageState(t *testing.T) {
	f := &navStackFixture{keepMounted: true}
	tt := ui.NewTester(f.view, 360, 240)
	if err := tt.Click("Count 0"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if _, found := tt.Find("Count 1"); found {
		t.Fatal("the covered page still exposes an input target")
	}
	tt.Key(ui.Alt, ui.KeyLeft)
	tt.Advance(time.Second)
	if !tt.HasText("Count 1") {
		t.Fatalf("the covered page lost its state: %q", tt.Texts())
	}
}

func TestNavigationStackBackButtonIsCenteredOnTheTitle(t *testing.T) {
	f := &navStackFixture{}
	tt := ui.NewTester(f.view, 320, 200)
	if err := tt.Click("Open network"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	back, _ := tt.Find("Back")
	title, _ := tt.Find("Network")
	if _, by := center(back); by < title.Y+title.H/2-1 || by > title.Y+title.H/2+1 {
		t.Fatalf("Back %+v is not centered on the title %+v", back, title)
	}
}

func TestNavigationStackActionsSitInTheHeaderOfEveryPage(t *testing.T) {
	var edits int
	pages := []NavigationPage{
		{Key: "a", Title: "Settings", Build: func() {}},
		{Key: "b", Title: "Network", Build: func() {}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		NavigationStack(c, NavigationStackOptions{
			Pages: pages, OnPop: func() {},
			Actions: func() { ui.Button(c, "Edit").OnClick(func() { edits++ }) },
		}).Height(200)
	}, 320, 240)
	if err := tt.Click("Edit"); err != nil {
		t.Fatalf("no action in the header: %v", err)
	}
	back, ok := tt.AccessByRole(ui.RoleButton, "Back")
	edit, _ := tt.AccessByRole(ui.RoleButton, "Edit")
	if edits != 1 || !ok || edit.Bounds.X < back.Bounds.X {
		t.Fatalf("edits %d; back %v, edit %v", edits, back.Bounds, edit.Bounds)
	}
}

// A pushed page is not a modal: Escape belongs to the editor inside it.
func TestNavigationStackEscapeLeavesTheEditorItsEditing(t *testing.T) {
	pops := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		NavigationStack(c, NavigationStackOptions{
			Pages: []NavigationPage{
				{Key: "list", Title: "Settings", Build: func() {}},
				{Key: "edit", Title: "Name", Build: func() { ui.Button(c, "Field") }},
			},
			OnPop: func() { pops++ },
		})
	}, 360, 240)
	if err := tt.Click("Field"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	if pops != 0 {
		t.Fatal("Escape popped the page")
	}
	tt.Key(ui.Alt, ui.KeyLeft)
	if pops != 1 {
		t.Fatalf("Alt+Left popped %d times", pops)
	}
}
