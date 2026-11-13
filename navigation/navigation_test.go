package navigation_test

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAccordionSingleExpansionAndKeyboardFocus(t *testing.T) {
	expanded := []string{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Accordion(c, navigation.AccordionOptions{
			Expanded: &expanded,
			Items: []navigation.AccordionItemOptions{
				{Key: "general", Title: "General", Body: func() { ui.Text(c, "General body") }},
				{Key: "privacy", Title: "Privacy", Body: func() { ui.Text(c, "Privacy body") }},
			},
		})
	}, 480, 420)

	if err := tt.Click("General"); err != nil {
		t.Fatal(err)
	}
	if len(expanded) != 1 || expanded[0] != "general" {
		t.Fatalf("after opening General: %v", expanded)
	}
	if !tt.HasText("General body") {
		t.Fatal("General panel is not shown")
	}
	if !tt.Focused("General") {
		t.Fatal("clicked header did not take the focus")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Privacy") {
		t.Fatal("Down did not move the focus to the next header")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Advance(400 * time.Millisecond)
	if len(expanded) != 1 || expanded[0] != "privacy" {
		t.Fatalf("single expansion keeps one open: %v", expanded)
	}
	if tt.HasText("General body") {
		t.Fatal("single expansion left the previous panel open")
	}
	if !tt.HasText("Privacy body") {
		t.Fatal("Privacy panel is not shown")
	}
}

func TestAccordionMultipleExpansion(t *testing.T) {
	expanded := []string{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Accordion(c, navigation.AccordionOptions{
			Expanded: &expanded,
			Multiple: true,
			Items: []navigation.AccordionItemOptions{
				{
					Key: "one", Title: "One", Body: func() { ui.Text(c, "One body") },
				},
				{
					Key: "two", Title: "Two", Body: func() { ui.Text(c, "Two body") },
				},
			},
		})
	}, 480, 420)

	if err := tt.Click("One"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Two"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(400 * time.Millisecond)
	if len(expanded) != 2 || expanded[0] != "one" || expanded[1] != "two" {
		t.Fatalf("multiple expansion: %v", expanded)
	}
	if !tt.HasText("One body") || !tt.HasText("Two body") {
		t.Fatal("multiple expansion did not keep both panels")
	}
}

func TestTabsSelectionAndPanelSwitch(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Tabs(c, navigation.TabsOptions{
			Selected: &selected,
			Label:    "Sections",
			Items: []navigation.TabOptions{
				{Label: "Overview", Content: func() { ui.Text(c, "Overview panel") }},
				{Label: "Settings", Content: func() { ui.Text(c, "Settings panel") }},
				{Label: "Advanced", Content: func() { ui.Text(c, "Advanced panel") }},
			},
		})
	}, 640, 320)

	if !tt.HasText("Overview panel") {
		t.Fatal("initial panel is not shown")
	}
	if err := tt.Click("Settings"); err != nil {
		t.Fatal(err)
	}
	if selected != 1 {
		t.Fatalf("selection=%d, want 1", selected)
	}
	if !tt.HasText("Settings panel") || tt.HasText("Overview panel") {
		t.Fatal("panel did not switch with the tab")
	}
	if !tt.Focused("Settings") {
		t.Fatal("clicked tab did not take the focus")
	}
	tt.Key(0, ui.KeyRight)
	if selected != 2 {
		t.Fatalf("Right arrow selection=%d, want 2", selected)
	}
	if !tt.HasText("Advanced panel") || tt.HasText("Settings panel") {
		t.Fatal("Right arrow did not switch the panel")
	}
	tt.Key(0, ui.KeyHome)
	if selected != 0 {
		t.Fatalf("Home selection=%d, want 0", selected)
	}
}

func TestPaginationBoundariesAndEllipsis(t *testing.T) {
	page := 1
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Pagination(c, navigation.PaginationOptions{
			Page: &page, Pages: 5, Siblings: 2, OnChange: func(int) { changes++ },
		})
	}, 640, 200)

	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if page != 2 || changes != 1 {
		t.Fatalf("Next page=%d changes=%d", page, changes)
	}
	if err := tt.Click("4"); err != nil {
		t.Fatal(err)
	}
	if page != 4 {
		t.Fatalf("page link set page=%d, want 4", page)
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if page != 5 {
		t.Fatalf("Next to last page=%d, want 5", page)
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if page != 5 {
		t.Fatalf("disabled Next moved to %d", page)
	}
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	if page != 4 {
		t.Fatalf("Previous page=%d, want 4", page)
	}
	if err := tt.Click("1"); err != nil {
		t.Fatal(err)
	}
	if page != 1 {
		t.Fatalf("page 1 link set page=%d", page)
	}
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	if page != 1 {
		t.Fatalf("disabled Previous moved to %d", page)
	}

	wide := 1
	tw := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Pagination(c, navigation.PaginationOptions{Page: &wide, Pages: 10})
	}, 640, 200)
	if !tw.HasText("More pages") {
		t.Fatal("expected an ellipsis between page groups")
	}
	if err := tw.Click("10"); err != nil {
		t.Fatal(err)
	}
	if wide != 10 {
		t.Fatalf("page link set page=%d, want 10", wide)
	}
}

func TestPaginationNarrowControlsKeepAccessibleLabels(t *testing.T) {
	page := 1
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Pagination(c, navigation.PaginationOptions{Page: &page, Pages: 4})
	}, 380, 120)
	next, ok := tt.Find("Next")
	if !ok || next.W != 32 || next.H != 32 {
		t.Fatalf("narrow Next bounds: %+v, found %v", next, ok)
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if page != 2 {
		t.Fatalf("narrow Next selected page %d", page)
	}
}

func TestPaginationSelectedPageHasVisibleDarkSurface(t *testing.T) {
	page := 2
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, true)
		ui.Column(c).Fill().Background(theme.Colors(c).Background).Center().Children(func() {
			navigation.Pagination(c, navigation.PaginationOptions{Page: &page, Pages: 5})
		})
	}, 640, 160)

	page2, ok := tt.AccessByRole(ui.RoleButton, "2")
	if !ok || page2.Current != ui.CurrentPage {
		t.Fatalf("the current page is not exposed: %+v", page2)
	}
	selected := page2.Bounds
	img := tt.Image()
	y := int(selected.Y + selected.H/2)
	inside := img.RGBAAt(int(selected.X+4), y)
	outside := img.RGBAAt(int(selected.X-2), y)
	if inside.R <= outside.R {
		t.Fatalf("selected page surface is not visible: inside=%v outside=%v", inside, outside)
	}
}

func TestDropdownActivationAndSubmenu(t *testing.T) {
	open := false
	selected := []string{}
	var calls []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Dropdown(c, navigation.DropdownOptions{
			Open:  &open,
			Label: "Actions",
			Menu: navigation.MenuOptions{
				Selection: navigation.MenuSelectMultiple,
				Selected:  &selected,
				OnSelect:  func(key string) { calls = append(calls, key) },
				Sections: []navigation.MenuSectionOptions{{
					Items: []navigation.MenuItemOptions{
						{Key: "edit", Label: "Edit", OnSelect: func() { calls = append(calls, "edit:cb") }},
						{Key: "share", Label: "Share", Submenu: []navigation.MenuItemOptions{
							{Key: "email", Label: "Email", OnSelect: func() { calls = append(calls, "email:cb") }},
							{Key: "copy", Label: "Copy link"},
						}},
					},
				}},
			},
		})
	}, 640, 480)

	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("trigger did not open the menu")
	}
	if !tt.HasText("Edit") {
		t.Fatal("menu items are not shown")
	}
	if err := tt.Click("Edit"); err != nil {
		t.Fatal(err)
	}
	if open {
		t.Fatal("choosing an item did not close the menu")
	}
	if len(calls) != 2 || calls[0] != "edit:cb" || calls[1] != "edit" {
		t.Fatalf("item activation calls=%v", calls)
	}
	if len(selected) != 1 || selected[0] != "edit" {
		t.Fatalf("selection=%v", selected)
	}

	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Share"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Email") {
		t.Fatal("submenu did not open")
	}
	if err := tt.Click("Email"); err != nil {
		t.Fatal(err)
	}
	if open {
		t.Fatal("choosing a submenu item did not close the menu")
	}
	if len(calls) != 4 || calls[3] != "email" {
		t.Fatalf("submenu activation calls=%v", calls)
	}
}

func TestDropdownKeyboardSubmenu(t *testing.T) {
	open := false
	var calls []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Dropdown(c, navigation.DropdownOptions{
			Open:  &open,
			Label: "More",
			Menu: navigation.MenuOptions{
				Sections: []navigation.MenuSectionOptions{{
					Items: []navigation.MenuItemOptions{
						{Key: "rename", Label: "Rename"},
						{Key: "export", Label: "Export", Submenu: []navigation.MenuItemOptions{
							{Key: "pdf", Label: "PDF", OnSelect: func() { calls = append(calls, "pdf") }},
						}},
					},
				}},
			},
		})
	}, 640, 480)

	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("menu did not open")
	}
	if !tt.Focused("Rename") {
		t.Fatal("opening did not focus the first item")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Export") {
		t.Fatal("Down did not move to the submenu trigger")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("PDF") {
		t.Fatal("Right did not open the submenu")
	}
	if !tt.Focused("PDF") {
		t.Fatal("submenu did not take the focus")
	}
	tt.Key(0, ui.KeyEnter)
	if len(calls) != 1 || calls[0] != "pdf" {
		t.Fatalf("submenu keyboard activation=%v", calls)
	}
	if open {
		t.Fatal("keyboard activation did not close the menu")
	}
}

func TestDisclosureGroupSingleExpansion(t *testing.T) {
	first, second := false, false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.DisclosureGroup(c, navigation.DisclosureGroupOptions{
			Items: []navigation.DisclosureOptions{
				{Title: "First", Open: &first, Body: func() { ui.Text(c, "First body") }},
				{Title: "Second", Open: &second, Body: func() { ui.Text(c, "Second body") }},
			},
		})
	}, 480, 360)

	if err := tt.Click("First"); err != nil {
		t.Fatal(err)
	}
	if !first {
		t.Fatal("First did not open")
	}
	if err := tt.Click("Second"); err != nil {
		t.Fatal(err)
	}
	if !second || first {
		t.Fatalf("single disclosure group first=%v second=%v", first, second)
	}
}

func TestTabsAlignVariantsPlaceLabels(t *testing.T) {
	for _, tc := range []struct {
		align navigation.TabsAlign
		name  string
	}{
		{navigation.TabsAlignCenter, "center"},
		{navigation.TabsAlignStart, "start"},
		{navigation.TabsAlignEnd, "end"},
	} {
		selected := 0
		align := tc.align
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Fill().Padding(16).Children(func() {
				navigation.Tabs(c, navigation.TabsOptions{
					Selected: &selected,
					Align:    align,
					Items: []navigation.TabOptions{
						{Label: "One", Content: func() { ui.Text(c, "One panel") }},
						{Label: "Two", Content: func() { ui.Text(c, "Two panel") }},
					},
				})
			})
		}, 640, 320)
		tab, ok := tt.Find("Two")
		if !ok {
			t.Fatalf("%s: second tab is not shown", tc.name)
		}
		left, right := tabInk(tt.Image(), tab)
		if left < 0 {
			t.Fatalf("%s: no label ink in the second tab", tc.name)
		}
		center := float32(left+right) / 2
		switch tc.align {
		case navigation.TabsAlignCenter:
			if d := center - (tab.X + tab.W/2); d < -8 || d > 8 {
				t.Fatalf("%s: label center %v, tab center %v", tc.name, center, tab.X+tab.W/2)
			}
		case navigation.TabsAlignStart:
			if d := float32(left) - tab.X; d < 8 || d > 28 {
				t.Fatalf("%s: label starts %v into the tab", tc.name, d)
			}
		case navigation.TabsAlignEnd:
			if d := tab.X + tab.W - float32(right); d < 8 || d > 28 {
				t.Fatalf("%s: label ends %v before the tab's end", tc.name, d)
			}
		}
	}
}

// tabInk returns the leftmost and rightmost columns of the label ink inside
// a tab, ignoring its edges.
func tabInk(img *image.RGBA, tab ui.Rect) (int, int) {
	x0, y0 := int(tab.X)+8, int(tab.Y)+6
	x1, y1 := int(tab.X+tab.W)-8, int(tab.Y+tab.H)-6
	bg := img.RGBAAt(x0, y0)
	left, right := -1, -1
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			p := img.RGBAAt(x, y)
			d := absInt(int(p.R)-int(bg.R)) + absInt(int(p.G)-int(bg.G)) + absInt(int(p.B)-int(bg.B))
			if d <= 60 {
				continue
			}
			if left < 0 || x < left {
				left = x
			}
			if x > right {
				right = x
			}
		}
	}
	return left, right
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func longMenu(n int) []navigation.MenuItemOptions {
	items := make([]navigation.MenuItemOptions, n)
	for i := range items {
		items[i] = navigation.MenuItemOptions{Key: fmt.Sprint("k", i), Label: fmt.Sprint("Item ", i)}
	}
	return items
}

func TestDropdownLongMenuKeyboardReachesEveryItem(t *testing.T) {
	open := false
	var chosen []string
	var ring ui.Color
	items := longMenu(1000)
	items[3].Disabled = true
	items[998].Separator = true
	items[10].Submenu = []navigation.MenuItemOptions{{Key: "a", Label: "Sub A"}, {Key: "b", Label: "Sub B"}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ring = c.Theme().Focus
		navigation.Dropdown(c, navigation.DropdownOptions{Open: &open, Label: "Long",
			Menu: navigation.MenuOptions{OnSelect: func(key string) { chosen = append(chosen, key) },
				Sections: []navigation.MenuSectionOptions{{Title: "All", Items: items}}}})
	}, 640, 600)
	if err := tt.Click("Long"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Item 0") || tt.HasText("Item 500") {
		t.Fatalf("opening focused Item 0 %v, built Item 500 %v", tt.Focused("Item 0"), tt.HasText("Item 500"))
	}
	for range 3 {
		tt.Key(0, ui.KeyDown)
	}
	if !tt.Focused("Item 4") {
		t.Fatal("Down did not skip the disabled item")
	}
	for range 25 {
		tt.Key(0, ui.KeyDown)
	}
	if !tt.Focused("Item 29") {
		t.Fatal("Down did not move past the first viewport")
	}
	tt.Key(0, ui.KeyEnd)
	if !tt.Focused("Item 999") {
		t.Fatal("End did not reach the last item")
	}
	if !focusRingLeftOf(tt, ring, "Item 999") {
		t.Fatal("the item End reached shows no focus ring")
	}
	if a, ok := tt.AccessByRole(ui.RoleMenuItem, "Item 999"); !ok || a.PosInSet != 999 || a.SetSize != 999 {
		t.Fatalf("last item's place for assistive technology: %d of %d", a.PosInSet, a.SetSize)
	}
	tt.Key(0, ui.KeyUp)
	if !tt.Focused("Item 997") {
		t.Fatal("Up did not skip the separator")
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyUp)
	if !tt.Focused("Item 999") {
		t.Fatal("Up from the first item did not wrap to the last")
	}
	tt.Key(0, ui.KeyEnter)
	if open || len(chosen) != 1 || chosen[0] != "k999" {
		t.Fatalf("Enter on the last item: open %v, chosen %v", open, chosen)
	}

	if err := tt.Click("Long"); err != nil {
		t.Fatal(err)
	}
	for range 9 { // Item 3 is disabled
		tt.Key(0, ui.KeyDown)
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Sub A") {
		t.Fatal("Right did not open the submenu of a long menu")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Sub B") {
		t.Fatal("Down did not move within the submenu")
	}
	tt.Key(0, ui.KeyEnter)
	if open || chosen[len(chosen)-1] != "b" {
		t.Fatalf("submenu activation: open %v, chosen %v", open, chosen)
	}
}

// focusRingLeftOf reports whether a pixel of the ring color crosses the
// row of label between the window's edge and the label.
func focusRingLeftOf(tt *ui.Tester, ring ui.Color, label string) bool {
	r, ok := tt.Find(label)
	if !ok {
		return false
	}
	img := tt.Image()
	y := int(r.Y + r.H/2)
	for x := 0; x < int(r.X); x++ {
		p := img.RGBAAt(x, y)
		if absDiff(p.R, ring.R) < 24 && absDiff(p.G, ring.G) < 24 && absDiff(p.B, ring.B) < 24 {
			return true
		}
	}
	return false
}

func absDiff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

func TestDropdownLongMenuScrollsToLastItem(t *testing.T) {
	open := true
	var chosen string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Dropdown(c, navigation.DropdownOptions{Open: &open, Label: "Long",
			Menu: navigation.MenuOptions{OnSelect: func(key string) { chosen = key },
				Sections: []navigation.MenuSectionOptions{{Items: longMenu(1000)}}}})
	}, 640, 600)
	tt.Frame()
	first, ok := tt.Find("Item 0")
	if !ok {
		t.Fatal("first item not shown")
	}
	for range 200 {
		if tt.HasText("Item 999") {
			break
		}
		tt.Scroll(first.X+first.W/2, first.Y+first.H/2, 0, 2000)
	}
	if err := tt.Click("Item 999"); err != nil {
		t.Fatal(err)
	}
	if chosen != "k999" || open {
		t.Fatalf("click on the last item: chosen %q, open %v", chosen, open)
	}
}

// A menu that grows past the long-menu size while open keeps its rows, and
// so the focus on them.
func TestDropdownKeepsFocusWhenItGrowsLongWhileOpen(t *testing.T) {
	open := false
	items := longMenu(20)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Dropdown(c, navigation.DropdownOptions{Open: &open, Label: "Grow",
			Menu: navigation.MenuOptions{Sections: []navigation.MenuSectionOptions{{Items: items}}}})
	}, 640, 600)
	if err := tt.Click("Grow"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Item 1") {
		t.Fatal("Down did not focus Item 1")
	}
	items = longMenu(60)
	tt.Frame()
	tt.Frame()
	if !tt.Focused("Item 1") {
		t.Fatal("growing the open menu dropped the focus")
	}
}

// TestTabsPanelIsATabPanelOfTheSelectedTab protects the tabs pattern for
// assistive technology: the panel shown is a tab panel, named and controlled
// by the selected tab, and follows the selection.
func TestTabsPanelIsATabPanelOfTheSelectedTab(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Tabs(c, navigation.TabsOptions{Selected: &selected, Label: "Sections", Items: []navigation.TabOptions{
			{Label: "Overview", Content: func() { ui.Text(c, "Overview panel") }},
			{Label: "Settings", Content: func() { ui.Text(c, "Settings panel") }},
		}})
	}, 640, 320)
	for _, name := range []string{"Overview", "Settings"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
		tab, _ := tt.AccessByRole(ui.RoleTab, name)
		if _, ok := tt.AccessByRole(ui.RoleTabPanel, name); !ok || len(tab.Controls) != 1 || tab.Controls[0] != name {
			t.Fatalf("%s: tab panel shown %v, the tab controls %q", name, ok, tab.Controls)
		}
	}
}
