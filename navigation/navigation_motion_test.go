package navigation_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// sectionLines returns the visible heights of four stacked body lines, in
// order. A line the growing panel has not revealed has no height.
func sectionLines(tt *ui.Tester) [4]float32 {
	var h [4]float32
	for i, text := range []string{"Alpha line", "Bravo line", "Charlie line", "Delta line"} {
		if r, ok := tt.Find(text); ok {
			h[i] = r.H
		}
	}
	return h
}

// linesShow reports whether all four lines show about want DIP each.
func linesShow(got [4]float32, want float32) bool {
	for _, h := range got {
		if h < want-0.5 || h > want+0.5 {
			return false
		}
	}
	return true
}

func fourLineAccordion(c *ui.Context, expanded *[]string) {
	ui.Column(c).Fill().Padding(16).Children(func() {
		navigation.Accordion(c, navigation.AccordionOptions{
			Expanded: expanded,
			Items: []navigation.AccordionItemOptions{{
				Key: "section", Title: "Section",
				Body: func() {
					ui.Column(c).Shrink(0).Children(func() {
						for _, line := range []string{"Alpha line", "Bravo line", "Charlie line", "Delta line"} {
							ui.Text(c, line).FontSize(14).FixedLineHeight(20)
						}
					})
				},
			}},
		})
	})
}

// TestAccordionPanelGrowsContinuously opens a section and reads the clip
// between frames: the first line shows whole, the next is partly revealed and
// the rest are not yet, so the panel grows instead of appearing.
func TestAccordionPanelGrowsContinuously(t *testing.T) {
	var expanded []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		fourLineAccordion(c, &expanded)
	}, 480, 360)
	if got := sectionLines(tt); got != [4]float32{} {
		t.Fatalf("closed section shows lines of %v", got)
	}

	if err := tt.Click("Section"); err != nil {
		t.Fatal(err)
	}
	if len(expanded) != 1 || expanded[0] != "section" {
		t.Fatalf("the header did not expand the section: %v", expanded)
	}
	tt.Advance(40 * time.Millisecond)
	mid := sectionLines(tt)
	var shown float32
	for i, h := range mid {
		shown += h
		if i > 0 && h > mid[i-1]+0.5 {
			t.Fatalf("the lines are not revealed in order: %v", mid)
		}
	}
	// The panel body is four 20 DIP lines and 16 DIP of bottom padding:
	// partway through, some lines show and the clip cuts the rest.
	if shown <= 0 || shown >= 96 {
		t.Fatalf("after 40 ms the panel shows %v of 96 DIP: no growth", shown)
	}

	tt.Advance(400 * time.Millisecond)
	if got := sectionLines(tt); !linesShow(got, 20) {
		t.Fatalf("settled open lines show %v, want four of 20", got)
	}
}

// TestAccordionReversalKeepsItsProgress closes and reopens mid-animation: the
// panel carries on from the height it reached instead of snapping shut or
// jumping open, and it settles where the last choice put it.
func TestAccordionReversalKeepsItsProgress(t *testing.T) {
	var expanded []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		fourLineAccordion(c, &expanded)
	}, 480, 360)
	if err := tt.Click("Section"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(400 * time.Millisecond)

	// Closing, then reopening 48 ms in.
	if err := tt.Click("Section"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(48 * time.Millisecond)
	closing := sectionLines(tt)
	if closing[0] < 19.5 || closing[3] != 0 {
		t.Fatalf("the closing panel did not keep a partial height: %v", closing)
	}
	if err := tt.Click("Section"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	reopened := sectionLines(tt)
	if reopened[0] < 19.5 {
		t.Fatalf("reopening started at %v DIP, not where the close left it", reopened[0])
	}
	if reopened[3] != 0 {
		t.Fatalf("reopening jumped open: %v", reopened)
	}

	tt.Advance(400 * time.Millisecond)
	if got := sectionLines(tt); !linesShow(got, 20) {
		t.Fatalf("reopened section settled at %v", got)
	}

	// Closing again settles closed, and a settled panel builds nothing.
	if err := tt.Click("Section"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(400 * time.Millisecond)
	if got := sectionLines(tt); got != [4]float32{} {
		t.Fatalf("closed section settled showing %v", got)
	}
	if tt.HasText("Alpha line") {
		t.Fatal("a settled closed section still builds its panel")
	}
}

// TestAccordionReducedMotionOpensAtOnce pins the reduced-motion contract: the
// panel is fully open in the next frame, with no reveal in between.
func TestAccordionReducedMotionOpensAtOnce(t *testing.T) {
	var expanded []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		fourLineAccordion(c, &expanded)
	}, 480, 360)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Section"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := sectionLines(tt); !linesShow(got, 20) {
		t.Fatalf("reduced motion opened to %v, not all four lines", got)
	}
}

func threeTabs(c *ui.Context, selected *int) {
	ui.Column(c).FillWidth().FillHeight().Padding(24).Children(func() {
		navigation.Tabs(c, navigation.TabsOptions{
			Selected: selected,
			Items: []navigation.TabOptions{
				{Label: "One", Content: func() {}},
				{Label: "Two", Content: func() {}},
				{Label: "Three", Content: func() {}},
			},
		})
	})
}

// pillPixel is a point inside a tab that only the selection pill covers: above
// the label ink and inside the rounded corner.
func pillPixel(tab ui.Rect) (int, int) {
	return int(tab.X + tab.W/2), int(tab.Y) + 3
}

// TestTabsIndicatorSlidesAndSettles checks the pill is still under the old tab
// on the frame the new one is chosen, and that it is gone from there once the
// slide settles on the new tab.
func TestTabsIndicatorSlidesAndSettles(t *testing.T) {
	selected := 0
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		threeTabs(c, &selected)
	}, 640, 320)

	one := mustFind(t, tt, "One")
	x, y := pillPixel(one)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Segment, "initial pill")

	if err := tt.Click("Two"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if selected != 1 {
		t.Fatalf("selection=%d, want 1", selected)
	}
	one, two := mustFind(t, tt, "One"), mustFind(t, tt, "Two")
	x, y = pillPixel(one)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Segment, "pill on the frame the next tab is chosen")
	x, y = pillPixel(two)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Default, "tab the pill has not reached")

	tt.Advance(400 * time.Millisecond)
	x, y = pillPixel(one)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Default, "old tab once the slide settles")
	x, y = pillPixel(two)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Segment, "settled pill")
}

// TestTabsIndicatorSettlesAfterReversal chooses a third tab mid-slide and
// checks the indicator ends on it.
func TestTabsIndicatorSettlesAfterReversal(t *testing.T) {
	selected := 0
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		threeTabs(c, &selected)
	}, 640, 320)
	if err := tt.Click("Two"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if err := tt.Click("Three"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(400 * time.Millisecond)
	one, three := mustFind(t, tt, "One"), mustFind(t, tt, "Three")
	x, y := pillPixel(one)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Default, "first tab after the reversal")
	x, y = pillPixel(three)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Segment, "pill after the reversal")
}

// TestTabsIndicatorReducedMotion snaps the pill to a new choice at once.
func TestTabsIndicatorReducedMotion(t *testing.T) {
	selected := 0
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		threeTabs(c, &selected)
	}, 640, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Two"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	one, two := mustFind(t, tt, "One"), mustFind(t, tt, "Two")
	x, y := pillPixel(two)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Segment, "pill at once")
	x, y = pillPixel(one)
	wantNavHoverPixel(t, tt.Image(), x, y, colors.Default, "old tab at once")
}

// TestTabsSecondaryUnderline checks the accent rule sits on the selected tab's
// bottom edge, two DIP thick, and that the tab is square rather than a pill.
func TestTabsSecondaryUnderline(t *testing.T) {
	selected := 0
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).Background(colors.Background).Children(func() {
			navigation.Tabs(c, navigation.TabsOptions{
				Variant:  navigation.TabsSecondary,
				Selected: &selected,
				Items: []navigation.TabOptions{
					{Label: "One", Content: func() {}},
					{Label: "Two", Content: func() {}},
				},
			})
		})
	}, 640, 320)
	one := mustFind(t, tt, "One")
	x := int(one.X + one.W/2)
	wantNavHoverPixel(t, tt.Image(), x, int(one.Y+one.H)-1, colors.Accent, "underline")
	wantNavHoverPixel(t, tt.Image(), x, int(one.Y+one.H)-4, colors.Background, "above the underline")
	wantNavHoverPixel(t, tt.Image(), int(one.X)+6, int(one.Y)+6, colors.Background, "secondary tab face")
}

// sameRGB reports whether a pixel shows u, within a few levels.
func sameRGB(c color.RGBA, u ui.Color) bool {
	return navHoverNear(c.R, u.R) && navHoverNear(c.G, u.G) && navHoverNear(c.B, u.B)
}

// TestDropdownEscapeCloses checks Escape dismisses the panel and keeps a
// caller-owned open flag in step.
func TestDropdownEscapeCloses(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.Dropdown(c, navigation.DropdownOptions{
			Open: &open, Label: "Actions",
			Menu: navigation.MenuOptions{Sections: []navigation.MenuSectionOptions{{
				Items: []navigation.MenuItemOptions{{Key: "edit", Label: "Edit"}},
			}}},
		})
	}, 640, 480)
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.HasText("Edit") {
		t.Fatal("the trigger did not open the menu")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || tt.HasText("Edit") {
		t.Fatal("Escape did not close the panel")
	}
}

// TestDropdownEntryFade checks the panel fades in on a dropdown's first
// frame, and that reduced motion shows it at once. The panel paints Overlay;
// before it arrives the page behind shows through.
func TestDropdownEntryFade(t *testing.T) {
	for _, reduced := range []bool{false, true} {
		name := "fade"
		if reduced {
			name = "reduced motion"
		}
		var colors theme.Palette
		open := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			ui.Column(c).Fill().Padding(24).Background(colors.Background).Children(func() {
				navigation.Dropdown(c, navigation.DropdownOptions{
					Open: &open, Label: "Actions",
					Menu: navigation.MenuOptions{Sections: []navigation.MenuSectionOptions{{
						Items: []navigation.MenuItemOptions{{Key: "edit", Label: "Edit"}},
					}}},
				})
			})
		}, 640, 480)
		if reduced {
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		}
		if err := tt.Click("Actions"); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		row := mustFind(t, tt, "Edit")
		x, y := int(row.X+row.W/2), int(row.Y)+4
		panel := sameRGB(tt.Image().RGBAAt(x, y), colors.Overlay)
		if panel != reduced {
			t.Fatalf("%s: panel shown=%v on the opening frame", name, panel)
		}
		tt.Advance(300 * time.Millisecond)
		if !sameRGB(tt.Image().RGBAAt(x, y), colors.Overlay) {
			t.Fatalf("%s: the panel did not settle on its surface", name)
		}
	}
}

// TestPaginationSummaryLayout checks the summary shares the row on a wide
// pagination and stacks above the controls below the 640 DIP breakpoint.
func TestPaginationSummaryLayout(t *testing.T) {
	page := 2
	view := func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).FillWidth().Padding(24).Children(func() {
			navigation.Pagination(c, navigation.PaginationOptions{
				Page: &page, Pages: 12, Siblings: 1,
				Summary: func() { ui.Text(c, "Showing 11-20 of 120 results") },
			})
		})
	}
	wide := ui.NewTester(view, 900, 240)
	summary := mustFind(t, wide, "Showing 11-20 of 120 results")
	next := mustFind(t, wide, "Next")
	if summary.X >= next.X {
		t.Fatalf("summary at %v is not before the controls at %v", summary.X, next.X)
	}
	if d := abs32(summary.Y + summary.H/2 - (next.Y + next.H/2)); d > 2 {
		t.Fatalf("summary and controls are %v DIP apart vertically, want one row", d)
	}

	narrow := ui.NewTester(view, 600, 320)
	narrow.Frame()
	narrow.Frame()
	summary = mustFind(t, narrow, "Showing 11-20 of 120 results")
	next = mustFind(t, narrow, "Next")
	if next.Y < summary.Y+summary.H {
		t.Fatalf("narrow: controls at %v do not stack below the summary ending at %v",
			next.Y, summary.Y+summary.H)
	}
}

func mustFind(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("%q is not shown", label)
	}
	return r
}
