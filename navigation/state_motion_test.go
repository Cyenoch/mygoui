package navigation

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// State changes of navigation controls pass through intermediate frames:
// these tests read pixels a quarter of the way and at the end of a change.

const stepFrame = 70 * time.Millisecond

func sum(c ui.Color) int { return int(c.R) + int(c.G) + int(c.B) }

func rgbSum(img *image.RGBA, x, y int) int {
	c := img.RGBAAt(x, y)
	return int(c.R) + int(c.G) + int(c.B)
}

// inkSum is the darkest pixel inside r away from its edges: a label's ink.
func inkSum(img *image.RGBA, r ui.Rect) int {
	best := 3 * 255
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			best = min(best, rgbSum(img, x, y))
		}
	}
	return best
}

func textBox(t *testing.T, tt *ui.Tester, text string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(text)
	if !ok {
		t.Fatalf("%q is not on screen", text)
	}
	return r
}

// between fails unless mid lies strictly between from and to by a margin.
func between(t *testing.T, what string, from, mid, to int) {
	t.Helper()
	lo, hi := min(from, to), max(from, to)
	if mid <= lo+3 || mid >= hi-3 {
		t.Fatalf("%s: %d is not between %d and %d", what, mid, from, to)
	}
}

func TestTabLabelInkFadesToTheChosenTab(t *testing.T) {
	// Colours fade under reduced motion too.
	for _, reduced := range []bool{false, true} {
		selected := 0
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			Tabs(c, TabsOptions{Selected: &selected, Items: []TabOptions{
				{Label: "One"}, {Label: "Two"}, {Label: "Three"},
			}})
		}, 480, 120)
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduced, TextScale: 1})
		tt.Frame()
		two := textBox(t, tt, "Two")
		rest := inkSum(tt.Image(), two)
		if err := tt.Click("Two"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(stepFrame)
		mid := inkSum(tt.Image(), two)
		tt.Advance(400 * time.Millisecond)
		between(t, "tab ink", rest, mid, inkSum(tt.Image(), two))
	}
}

func TestSegmentLabelInkFadesToTheChosenSegment(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		SegmentedControl(c, SegmentedControlOptions{Selected: &selected, Items: []TabOptions{
			{Label: "Day"}, {Label: "Week"}, {Label: "Month"},
		}})
	}, 480, 120)
	week := textBox(t, tt, "Week")
	rest := inkSum(tt.Image(), week)
	if err := tt.Click("Week"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(stepFrame)
	mid := inkSum(tt.Image(), week)
	tt.Advance(400 * time.Millisecond)
	between(t, "segment ink", rest, mid, inkSum(tt.Image(), week))
}

func TestPaginationActivePageFaceFades(t *testing.T) {
	page := 1
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, true)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Background(colors.Background).Children(func() {
			Pagination(c, PaginationOptions{Page: &page, Pages: 3})
		})
	}, 480, 120)
	label := textBox(t, tt, "3")
	x, y := int(label.X)+5, int(label.Y+label.H/2)
	rest := rgbSum(tt.Image(), x, y)
	if err := tt.Click("3"); err != nil {
		t.Fatal(err)
	}
	tt.Move(2, 2)
	tt.Frame()
	tt.Advance(25 * time.Millisecond)
	mid := rgbSum(tt.Image(), x, y)
	tt.Advance(400 * time.Millisecond)
	between(t, "active page face", rest, mid, rgbSum(tt.Image(), x, y))
}

func TestPipsPagerCurrentPipDarkens(t *testing.T) {
	index := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Padding(24).Children(func() {
			PipsPager(c, PipsPagerOptions{Index: &index, Count: 3})
		})
	}, 200, 80)
	pip, ok := tt.Access("Page 2")
	if !ok {
		t.Fatal("no second pip")
	}
	x, y := int(pip.Bounds.X+pip.Bounds.W/2), int(pip.Bounds.Y+pip.Bounds.H/2)
	rest := rgbSum(tt.Image(), x, y)
	index = 1
	tt.Frame()
	tt.Advance(stepFrame)
	mid := rgbSum(tt.Image(), x, y)
	tt.Advance(400 * time.Millisecond)
	between(t, "current pip", rest, mid, rgbSum(tt.Image(), x, y))
}

func TestStatusBarMessageCrossfades(t *testing.T) {
	message := "Saved"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		StatusBar(c, StatusBarOptions{Message: message})
	}, 480, 60)
	region := ui.Rect{X: 0, Y: 4, W: 120, H: 16}
	message = "Syncing"
	tt.Frame()
	tt.Advance(stepFrame)
	mid := inkSum(tt.Image(), region)
	tt.Advance(400 * time.Millisecond)
	// Partway the new message is not yet at full strength.
	if settled := inkSum(tt.Image(), region); mid <= settled+30 {
		t.Fatalf("message ink %d partway, %d settled: the text did not fade", mid, settled)
	}
}

func TestAccordionPanelFadesWithItsHeight(t *testing.T) {
	var expanded []string
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Background(colors.Background).Children(func() {
			Accordion(c, AccordionOptions{Expanded: &expanded, Items: []AccordionItemOptions{{
				Key: "s", Title: "Section",
				Body: func() {
					ui.Box(c).FillWidth().Height(60).Background(ui.RGB(255, 0, 0))
				},
			}}})
		})
	}, 480, 240)
	title := textBox(t, tt, "Section")
	x := int(title.X) + 4
	// The first row of the body that shows, scanning down from the header.
	topmost := func() (c [4]uint8) {
		// Start under the title: subpixel text has coloured fringes.
		for y := int(title.Y + title.H); y < 230; y++ {
			if p := tt.Image().RGBAAt(x, y); int(p.R)-int(p.G) > 10 {
				return [4]uint8{p.R, p.G, p.B, p.A}
			}
		}
		return c
	}
	if err := tt.Click("Section"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	mid := topmost()
	tt.Advance(400 * time.Millisecond)
	end := topmost()
	if end[0] < 250 || end[1] > 5 {
		t.Fatalf("settled panel pixel %v is not opaque red", end)
	}
	if mid[0] == 0 || mid[1] < 20 {
		t.Fatalf("panel pixel %v partway is not blended with the page", mid)
	}
}

func TestMenuCheckDrawsInAndFadesOut(t *testing.T) {
	for _, reduced := range []bool{false, true} {
		selected := []string{}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			Menu(c, MenuOptions{
				Selection: MenuSelectMultiple, Selected: &selected,
				Sections: []MenuSectionOptions{{Items: []MenuItemOptions{
					{Key: "a", Label: "Alpha"}, {Key: "b", Label: "Bravo"},
				}}},
			})
		}, 320, 160)
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduced, TextScale: 1})
		tt.Frame()
		item, ok := tt.AccessByRole(ui.RoleMenuItemCheckBox, "Alpha")
		if !ok {
			t.Fatal("no Alpha menu item")
		}
		mark := ui.Rect{X: item.Bounds.X + 6, Y: item.Bounds.Y + 2, W: 20, H: item.Bounds.H - 4}
		inked := func() int {
			n := 0
			for y := int(mark.Y); y < int(mark.Y+mark.H); y++ {
				for x := int(mark.X); x < int(mark.X+mark.W); x++ {
					if rgbSum(tt.Image(), x, y) < 600 {
						n++
					}
				}
			}
			return n
		}
		if inked() != 0 {
			t.Fatalf("reduced=%v: an unchecked item shows a mark", reduced)
		}
		selected = []string{"a"}
		tt.Frame()
		tt.Advance(stepFrame / 2)
		mid := inked()
		tt.Advance(400 * time.Millisecond)
		full := inked()
		if mid <= 0 || mid >= full {
			t.Fatalf("reduced=%v: mark has %d inked pixels partway and %d settled", reduced, mid, full)
		}
		selected = nil
		tt.Frame()
		tt.Advance(stepFrame / 2)
		leaving := inked()
		if leaving <= 0 || leaving > full {
			t.Fatalf("reduced=%v: a leaving mark has %d inked pixels, want a fading part of %d", reduced, leaving, full)
		}
		tt.Advance(400 * time.Millisecond)
		if inked() != 0 {
			t.Fatalf("reduced=%v: the mark did not leave", reduced)
		}
	}
}

func TestStepsFaceFadesAndTickDrawsIn(t *testing.T) {
	current := 0
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Background(colors.Background).Padding(16).Children(func() {
			Steps(c, StepsOptions{Items: stepItems(), Current: current})
		})
	}, 480, 120)
	bright := func(box ui.Rect) int {
		n := 0
		for y := int(box.Y); y < int(box.Y+box.H); y++ {
			for x := int(box.X); x < int(box.X+box.W); x++ {
				if rgbSum(tt.Image(), x, y) > 3*215 {
					n++
				}
			}
		}
		return n
	}
	// The second step's face goes from muted to primary as it becomes current.
	two := textBox(t, tt, "2")
	x, y := int(two.X)-6, int(two.Y+two.H/2)
	rest := rgbSum(tt.Image(), x, y)
	one := textBox(t, tt, "1")
	dot := ui.Rect{X: one.X - 8, Y: one.Y - 8, W: 24, H: 24}
	current = 1
	tt.Frame()
	tt.Advance(stepFrame)
	between(t, "step face", rest, rgbSum(tt.Image(), x, y), sum(colors.Primary))
	if !tt.HasText("1") {
		t.Fatal("the number left before the tick drew in")
	}
	partly := bright(dot)
	tt.Advance(400 * time.Millisecond)
	settled := bright(dot)
	if tt.HasText("1") || partly >= settled || settled == 0 {
		t.Fatalf("tick: %d bright pixels partway, %d settled, number shown %v", partly, settled, tt.HasText("1"))
	}

	// A step built completed shows its tick at once.
	built := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).FillWidth().FillHeight().Background(theme.Colors(c).Background).Padding(16).Children(func() {
			Steps(c, StepsOptions{Items: stepItems(), Current: 1})
		})
	}, 480, 120)
	tt = built
	if got := bright(dot); got != settled {
		t.Fatalf("a step built completed shows %d bright pixels on its first frame, want %d", got, settled)
	}
}

func TestDocumentTabActiveFaceFadesAndClosedTabLeaves(t *testing.T) {
	f := newDocTabsFixture()
	f.tabs = f.tabs[:2]
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		f.view(c)
		colors = theme.Colors(c)
	}, 480, 80)
	tab, ok := tt.AccessByRole(ui.RoleTab, "notes.md")
	if !ok {
		t.Fatal("no notes.md tab")
	}
	x, y := int(tab.Bounds.X)+3, int(tab.Bounds.Y+tab.Bounds.H/2)
	rest := rgbSum(tt.Image(), x, y)
	f.active = "b"
	tt.Frame()
	tt.Advance(stepFrame)
	mid := rgbSum(tt.Image(), x, y)
	tt.Advance(400 * time.Millisecond)
	between(t, "active tab face", rest, mid, rgbSum(tt.Image(), x, y))
	_ = colors

	// A tab the caller closes shrinks away rather than vanishing.
	closed := tab.Bounds.W
	f.tabs = f.tabs[:1]
	f.active = "a"
	tt.Frame()
	tt.Advance(stepFrame)
	left, ok := tt.AccessByRole(ui.RoleTab, "notes.md")
	if ok && left.Bounds.W >= closed {
		t.Fatalf("the closing tab is still %v DIP wide, from %v", left.Bounds.W, closed)
	}
	tt.Advance(400 * time.Millisecond)
	if _, ok := tt.AccessByRole(ui.RoleTab, "notes.md"); ok {
		t.Fatal("the closed tab did not leave")
	}
}
