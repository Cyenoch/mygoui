package navigation

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// docTabsFixture is a caller that records what the strip asks for and
// applies none of it, so each test sees the strip's reports alone.
type docTabsFixture struct {
	tabs     []DocumentTab
	active   string
	selected []string
	closed   []string
	moves    [][2]int
	added    int
}

func newDocTabsFixture() *docTabsFixture {
	return &docTabsFixture{
		tabs: []DocumentTab{
			{Key: "a", Title: "main.go"}, {Key: "b", Title: "notes.md", Dirty: true}, {Key: "c", Title: "go.mod"},
		},
		active: "a",
	}
}

func (f *docTabsFixture) view(c *ui.Context) {
	theme.Apply(c)
	DocumentTabs(c, DocumentTabsOptions{
		Label: "Documents", Tabs: f.tabs, Active: f.active,
		OnSelect: func(key string) { f.selected = append(f.selected, key) },
		OnClose:  func(key string) { f.closed = append(f.closed, key) },
		OnMove:   func(from, to int) { f.moves = append(f.moves, [2]int{from, to}) },
		OnAdd:    func() { f.added++ },
	})
}

func center(r ui.Rect) (float32, float32) { return r.X + r.W/2, r.Y + r.H/2 }

func TestDocumentTabsExposeTabsAndReportSelection(t *testing.T) {
	f := newDocTabsFixture()
	tt := ui.NewTester(f.view, 480, 80)
	if _, ok := tt.AccessByRole(ui.RoleTabList, "Documents"); !ok {
		t.Fatal("no tab list named Documents")
	}
	for title, selected := range map[string]bool{"main.go": true, "notes.md": false, "go.mod": false} {
		info, ok := tt.AccessByRole(ui.RoleTab, title)
		if !ok || info.Selected != selected {
			t.Errorf("%s: tab %v, selected %v, want %v", title, ok, info.Selected, selected)
		}
	}
	if info, _ := tt.AccessByRole(ui.RoleTab, "notes.md"); info.Description != "Modified" {
		t.Errorf("unsaved tab is described as %q", info.Description)
	}
	if err := tt.Click("go.mod"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Add"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(f.selected, []string{"c"}) || f.added != 1 {
		t.Fatalf("selected %v, added %d", f.selected, f.added)
	}
	if f.active != "a" || len(f.tabs) != 3 || len(f.closed) != 0 || len(f.moves) != 0 {
		t.Fatalf("the strip changed caller state: active %q, tabs %d, closed %v, moves %v",
			f.active, len(f.tabs), f.closed, f.moves)
	}
}

func TestDocumentTabsCloseByButtonKeysAndMiddleClick(t *testing.T) {
	f := newDocTabsFixture()
	tt := ui.NewTester(f.view, 480, 80)
	// The active tab offers its close button at rest.
	if err := tt.Click("Close main.go"); err != nil {
		t.Fatal(err)
	}
	if len(f.selected) != 0 {
		t.Fatalf("the close button also selected: %v", f.selected)
	}
	if err := tt.Click("go.mod"); err != nil {
		t.Fatal(err)
	}
	f.selected = nil
	tt.Key(ui.Cmd, ui.KeyW)
	tt.Key(0, ui.KeyDelete)
	notes, _ := tt.Find("notes.md")
	x, y := center(notes)
	tt.PressButton(2, x, y)
	tt.ReleaseButton(2, x, y+60) // let go elsewhere: no close
	tt.PressButton(2, x, y)
	tt.ReleaseButton(2, x, y)
	if want := []string{"a", "a", "c", "b"}; !slices.Equal(f.closed, want) {
		t.Fatalf("closed %v, want %v", f.closed, want)
	}
	if len(f.selected) != 0 || len(f.tabs) != 3 {
		t.Fatalf("closing selected %v or removed a tab (%d left)", f.selected, len(f.tabs))
	}
}

func TestDocumentTabsFocusMovesToTheNeighbourOfAClosedTab(t *testing.T) {
	f := newDocTabsFixture()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{
			Tabs: f.tabs, Active: f.active,
			OnClose: func(key string) {
				f.tabs = slices.DeleteFunc(f.tabs, func(tab DocumentTab) bool { return tab.Key == key })
			},
		})
	}, 480, 80)
	if err := tt.Click("go.mod"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDelete)
	if len(f.tabs) != 2 || !tt.Focused("notes.md") {
		t.Fatalf("%d tabs left, notes.md focused %v", len(f.tabs), tt.Focused("notes.md"))
	}
	tt.Key(0, ui.KeyDelete)
	if len(f.tabs) != 1 || !tt.Focused("main.go") {
		t.Fatalf("%d tabs left, main.go focused %v", len(f.tabs), tt.Focused("main.go"))
	}
}

func TestDocumentTabsReorderByKeysAndDrag(t *testing.T) {
	f := newDocTabsFixture()
	apply := true
	view := func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{
			Tabs: f.tabs, Active: f.active,
			OnMove: func(from, to int) {
				f.moves = append(f.moves, [2]int{from, to})
				if apply {
					tab := f.tabs[from]
					f.tabs = slices.Insert(slices.Delete(f.tabs, from, from+1), to, tab)
				}
			},
		})
	}
	tt := ui.NewTester(view, 480, 80)
	if err := tt.Click("main.go"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Alt, ui.KeyRight)
	if len(f.moves) != 1 || f.moves[0] != [2]int{0, 1} || f.tabs[1].Key != "a" {
		t.Fatalf("Alt+Right moved %v, order %v", f.moves, f.tabs)
	}
	if !tt.Focused("main.go") {
		t.Fatal("the focus did not follow the moved tab")
	}
	tt.Key(ui.Alt, ui.KeyRight)
	tt.Key(ui.Alt, ui.KeyRight) // already last
	if len(f.moves) != 2 || f.tabs[2].Key != "a" {
		t.Fatalf("moves %v, order %v", f.moves, f.tabs)
	}

	// Drag the last tab onto the first.
	apply, f.moves = false, nil
	from, _ := tt.Find("main.go")
	first, _ := tt.Find("notes.md")
	x, y := center(from)
	tx, _ := center(first)
	tt.Press(x, y)
	for step := 1; step <= 4; step++ {
		tt.Move(x+(tx-x)*float32(step)/4, y)
	}
	tt.Release(tx, y)
	if len(f.moves) != 1 || f.moves[0] != [2]int{2, 0} {
		t.Fatalf("drag moved %v, want [[2 0]]", f.moves)
	}
}

func TestDocumentTabsKeepTheActiveTabInView(t *testing.T) {
	var tabs []DocumentTab
	for i := range 20 {
		tabs = append(tabs, DocumentTab{Key: fmt.Sprint(i), Title: fmt.Sprintf("document-%02d.txt", i)})
	}
	active := "0"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{Tabs: tabs, Active: active, OnAdd: func() {}})
	}, 360, 80)
	add, ok := tt.Find("Add")
	if !ok || add.X+add.W > 360 {
		t.Fatalf("the add button left the window: %+v", add)
	}
	active = "19"
	tt.Frame()
	tt.Advance(time.Second) // the strip glides to it
	last, ok := tt.Find("document-19.txt")
	if !ok || last.W <= 0 || last.X < 0 || last.X+last.W > add.X {
		t.Fatalf("the active tab is not in view: %+v (add at %v)", last, add.X)
	}
}

// A closed tab leaves over a moment and the tabs after it slide into its
// place; a tab added grows the same way. Neither is left pending at rest.
func TestDocumentTabsCloseAndAddAnimate(t *testing.T) {
	f := newDocTabsFixture()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{Tabs: f.tabs, Active: "a", OnClose: func(string) {}})
	}, 480, 80)
	if tt.FramePending() {
		t.Fatal("the strip asks for frames at rest")
	}
	before, _ := tt.Find("go.mod")
	f.tabs = slices.Delete(f.tabs, 1, 2)
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	mid, _ := tt.Find("go.mod")
	if mid.X >= before.X || mid.W != before.W {
		t.Errorf("go.mod at %+v did not start sliding from %+v", mid, before)
	}
	tt.Advance(time.Second)
	after, _ := tt.Find("go.mod")
	if tt.FramePending() || after.X >= mid.X {
		t.Errorf("frames pending or go.mod did not settle: go.mod %+v then %+v", mid, after)
	}

	f.tabs = append(f.tabs, DocumentTab{Key: "d", Title: "new.txt"})
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	grown, ok := tt.Find("new.txt")
	tt.Advance(time.Second)
	full, _ := tt.Find("new.txt")
	if !ok || grown.W >= full.W || tt.FramePending() {
		t.Errorf("new tab went from %+v to %+v", grown, full)
	}
}

func manyDocTabs(n int) []DocumentTab {
	tabs := make([]DocumentTab, n)
	for i := range tabs {
		tabs[i] = DocumentTab{Key: fmt.Sprint(i), Title: fmt.Sprintf("document-%02d.txt", i)}
	}
	return tabs
}

func TestDocumentTabsContextMenuClosesTheOthersAndTheTabsAfterButNotPinned(t *testing.T) {
	tabs := manyDocTabs(5)
	tabs[0].Pinned = true
	var closed [][]string
	var single []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{
			Tabs: tabs, Active: "1",
			OnClose:     func(key string) { single = append(single, key) },
			OnCloseMany: func(keys []string) { closed = append(closed, keys) },
		})
	}, 900, 200)
	if err := tt.RightClick("document-02.txt"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Close others"); err != nil {
		t.Fatal(err)
	}
	if err := tt.RightClick("document-02.txt"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Close to the right"); err != nil {
		t.Fatal(err)
	}
	if err := tt.RightClick("document-04.txt"); err != nil {
		t.Fatal(err)
	}
	if info, ok := tt.AccessByRole(ui.RoleMenuItem, "Close to the right"); !ok || !info.Disabled {
		t.Fatalf("Close to the right on the last tab: %+v, found %v", info, ok)
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if want := [][]string{{"1", "3", "4"}, {"3", "4"}}; !slices.EqualFunc(closed, want, slices.Equal[[]string]) {
		t.Errorf("closed together %v, want %v", closed, want)
	}
	if !slices.Equal(single, []string{"4"}) {
		t.Errorf("closed singly %v", single)
	}
}

func TestDocumentTabsPinnedTabsAreIconsWithoutCloseAndStayFirst(t *testing.T) {
	tabs := manyDocTabs(3)
	tabs[0].Pinned = true
	var moves [][2]int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{
			Tabs: tabs, Active: "0", OnClose: func(string) {},
			OnMove: func(from, to int) { moves = append(moves, [2]int{from, to}) },
		})
	}, 480, 80)
	pinned, ok := tt.Find("document-00.txt")
	other, _ := tt.Find("document-01.txt")
	if !ok || pinned.W > other.W/2 {
		t.Fatalf("pinned tab %+v beside %+v", pinned, other)
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, "Close document-00.txt"); ok {
		t.Fatal("a pinned tab has a close button")
	}
	if err := tt.Click("document-01.txt"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Alt, ui.KeyLeft) // into the pinned group: refused
	if len(moves) != 0 {
		t.Fatalf("a tab crossed the pinned boundary: %v", moves)
	}
}

func TestDocumentTabsKeysChooseTabsAndCloseTheActiveOne(t *testing.T) {
	f := newDocTabsFixture()
	tt := ui.NewTester(f.view, 480, 80)
	if err := tt.Click("main.go"); err != nil {
		t.Fatal(err)
	}
	f.selected = nil
	tt.Key(ui.Ctrl, ui.KeyTab)
	tt.Key(ui.Ctrl|ui.Shift, ui.KeyTab)
	if !slices.Equal(f.selected, []string{"b", "c"}) {
		t.Fatalf("Ctrl+Tab then Ctrl+Shift+Tab chose %v", f.selected)
	}
}

func TestDocumentTabsMenuListsEveryTabAndScrollsToTheChosenOne(t *testing.T) {
	tabs := manyDocTabs(40)
	active := "0"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{Tabs: tabs, Active: active, OnSelect: func(k string) { active = k }})
	}, 480, 400)
	tt.Advance(time.Second)
	if err := tt.Click("All tabs"); err != nil {
		t.Fatalf("a strip with 40 tabs offers no list of them: %v", err)
	}
	tt.Advance(time.Second)
	if info, ok := tt.AccessByRole(ui.RoleMenuItemRadio, "document-00.txt"); !ok || !info.Checked {
		t.Fatalf("the active tab is not checked in the list: %+v, found %v", info, ok)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if active != "39" {
		t.Fatalf("chose %q", active)
	}
	last, ok := tt.Find("document-39.txt")
	if !ok || last.X < 0 || last.X+last.W > 480 {
		t.Fatalf("the chosen tab is not in view: %+v", last)
	}
}

func TestDocumentTabsWheelScrollsTheStripSideways(t *testing.T) {
	tabs := manyDocTabs(40)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{Tabs: tabs, Active: "0"})
	}, 480, 80)
	tt.Advance(time.Second)
	before, _ := tt.Find("document-03.txt")
	tt.Wheel(200, 20, 0, 160)
	tt.Advance(time.Second)
	after, ok := tt.Find("document-03.txt")
	if ok && after.X >= before.X {
		t.Fatalf("a downward wheel left the tabs at %v from %v", after.X, before.X)
	}
	tt.Wheel(200, 20, 0, -400)
	tt.Advance(time.Second)
	if back, _ := tt.Find("document-03.txt"); back.X != before.X {
		t.Fatalf("scrolling back left the tab at %v, was %v", back.X, before.X)
	}
}

// Shrinking tabs share the strip, and the strip scrolls only once they reach
// their minimum.
func TestDocumentTabsShrinkToFitBeforeScrolling(t *testing.T) {
	tabs := manyDocTabs(4)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		DocumentTabs(c, DocumentTabsOptions{Tabs: tabs, Active: "0", Sizing: DocumentTabsShrink, OnAdd: func() {}})
	}, 560, 80)
	tt.Advance(time.Second)
	first, _ := tt.Find("document-00.txt")
	last, _ := tt.Find("document-03.txt")
	add, _ := tt.Find("Add")
	if first.W >= 200 || last.X+last.W > add.X {
		t.Fatalf("four tabs in 560 DIP: first %+v, last %+v, add %+v", first, last, add)
	}
	if _, ok := tt.Find("All tabs"); ok {
		t.Fatal("tabs that fit offer a list of all tabs")
	}
	tabs = manyDocTabs(12)
	tt.Frame()
	tt.Advance(time.Second)
	if _, ok := tt.Find("All tabs"); !ok {
		t.Fatal("tabs at their minimum width offer no list of all tabs")
	}
}

// The inactive title is the theme's muted ink, lifted where it falls short of
// 4.5:1 on the strip's face.
func TestDocumentTabsInactiveTitlesReadInBothAppearances(t *testing.T) {
	for _, dark := range []bool{false, true} {
		tt := ui.NewTester(func(c *ui.Context) {
			p := theme.Colors(c)
			face := p.MutedBackground.Over(p.Background)
			if got := contrastRatio(legibleInk(p.MutedForeground, p.Foreground, face), face); got < 4.5 {
				t.Errorf("dark=%v: inactive ink at %.2f:1", dark, got)
			}
		}, 40, 40)
		tt.SetDark(dark)
		tt.Frame()
	}
}

// docTabsLiveFixture applies what the strip asks for, beside an editor the
// window shortcuts must not take the focus from.
type docTabsLiveFixture struct {
	tabs   []DocumentTab
	active string
	closed []string
}

func (f *docTabsLiveFixture) view(c *ui.Context) {
	theme.Apply(c)
	ui.Column(c).Children(func() {
		ui.Button(c, "Editor")
		DocumentTabs(c, DocumentTabsOptions{
			Tabs: f.tabs, Active: f.active, WindowShortcuts: true,
			OnSelect: func(key string) { f.active = key },
			OnClose: func(key string) {
				f.closed = append(f.closed, key)
				f.tabs = slices.DeleteFunc(f.tabs, func(tab DocumentTab) bool { return tab.Key == key })
			},
		})
	})
}

func TestDocumentTabsWindowShortcutsKeepTheEditorFocused(t *testing.T) {
	f := &docTabsLiveFixture{tabs: newDocTabsFixture().tabs, active: "a"}
	tt := ui.NewTester(f.view, 480, 120)
	if err := tt.Click("Editor"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Ctrl, ui.KeyTab)
	if f.active != "b" || !tt.Focused("Editor") {
		t.Fatalf("Ctrl+Tab: active %q, editor focused %v", f.active, tt.Focused("Editor"))
	}
	tt.Key(ui.Cmd, ui.KeyW)
	if !slices.Equal(f.closed, []string{"b"}) || !tt.Focused("Editor") {
		t.Fatalf("Cmd+W: closed %v, editor focused %v", f.closed, tt.Focused("Editor"))
	}
}

// With the focus in the strip, a tab closed by the key hands it to the tab
// that takes its place, and Cmd+W closes the active tab, never a pinned one.
func TestDocumentTabsCloseKeyActsOnTheActiveTab(t *testing.T) {
	f := &docTabsLiveFixture{tabs: newDocTabsFixture().tabs, active: "a"}
	tt := ui.NewTester(f.view, 480, 120)
	if err := tt.Click("notes.md"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("main.go"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight) // focus b, while a stays active
	tt.Key(ui.Cmd, ui.KeyW)
	if !slices.Equal(f.closed, []string{"a"}) || !tt.Focused("notes.md") {
		t.Fatalf("closed %v, notes.md focused %v", f.closed, tt.Focused("notes.md"))
	}

	f = &docTabsLiveFixture{tabs: newDocTabsFixture().tabs, active: "a"}
	tt = ui.NewTester(f.view, 480, 120)
	if err := tt.RightClick("notes.md"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if !slices.Equal(f.closed, []string{"b"}) || !tt.Focused("go.mod") {
		t.Fatalf("menu close: closed %v, go.mod focused %v", f.closed, tt.Focused("go.mod"))
	}

	f = &docTabsLiveFixture{tabs: []DocumentTab{{Key: "p", Title: "Pinned", Pinned: true}, {Key: "q", Title: "q.go"}}, active: "p"}
	tt = ui.NewTester(f.view, 480, 120)
	tt.Key(ui.Cmd, ui.KeyW)
	if len(f.closed) != 0 {
		t.Fatalf("Cmd+W closed the pinned tab: %v", f.closed)
	}
}

func TestDocumentTabsWithoutAnActiveTabStepFromTheEnds(t *testing.T) {
	f := &docTabsLiveFixture{tabs: newDocTabsFixture().tabs, active: "gone"}
	tt := ui.NewTester(f.view, 480, 120)
	tt.Key(ui.Ctrl, ui.KeyTab)
	if f.active != "a" {
		t.Fatalf("forward chose %q", f.active)
	}
	f.active = "gone"
	tt.Key(ui.Ctrl|ui.Shift, ui.KeyTab)
	if f.active != "c" {
		t.Fatalf("backward chose %q", f.active)
	}
}
