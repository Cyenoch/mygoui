package navigation_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestToolbarIsOneTabStopWithWrappingArrows protects the toolbar pattern: Tab
// enters and leaves the toolbar once, the arrows of its orientation move
// between its controls and wrap, and Tab returns to the control used last.
func TestToolbarIsOneTabStopWithWrappingArrows(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Before")
		navigation.Toolbar(c, navigation.ToolbarOptions{Label: "Format", Items: func() {
			ui.Button(c, "Bold")
			navigation.ToolbarSeparator(c, mygoui.Horizontal)
			ui.Button(c, "Italic")
			ui.Button(c, "Underline")
		}})
		ui.Button(c, "After")
	}, 480, 200)
	if err := tt.Click("Before"); err != nil {
		t.Fatal(err)
	}
	steps := []struct {
		mods ui.Modifiers
		key  ui.Key
		want string
	}{
		{0, ui.KeyTab, "Bold"},
		{0, ui.KeyLeft, "Underline"},
		{0, ui.KeyRight, "Bold"},
		{0, ui.KeyEnd, "Underline"},
		{0, ui.KeyHome, "Bold"},
		{0, ui.KeyRight, "Italic"},
		{0, ui.KeyDown, "Italic"},
		{0, ui.KeyTab, "After"},
		{ui.Shift, ui.KeyTab, "Italic"},
	}
	for _, s := range steps {
		tt.Key(s.mods, s.key)
		if !tt.Focused(s.want) {
			t.Fatalf("after key %v with modifiers %v the focus is not on %s", s.key, s.mods, s.want)
		}
	}
}

// TestDisclosureTogglesWithEnterAndSpace protects keyboard operation of a
// disclosure: both keys flip the caller's state and the exposed expansion.
func TestDisclosureTogglesWithEnterAndSpace(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Before")
		navigation.Disclosure(c, navigation.DisclosureOptions{Title: "Details", Open: &open,
			Body: func() { ui.Text(c, "Hidden detail") }})
	}, 480, 300)
	if err := tt.Click("Before"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	for i, key := range []ui.Key{ui.KeyEnter, ui.KeySpace, ui.KeySpace} {
		tt.Key(0, key)
		tt.Frame()
		want := i%2 == 0
		a, ok := tt.AccessByRole(ui.RoleDisclosure, "Details")
		if !ok || open != want || a.Expanded != want {
			t.Fatalf("after key %d: open=%v expanded=%v, want %v", i+1, open, a.Expanded, want)
		}
	}
}

// TestBreadcrumbsFollowLinksFromTheKeyboard protects the trail's keyboard
// path: crumbs before the current one are links reached by Tab and followed
// with Enter, the current crumb is text marked current, and collapsed crumbs are offered in
// the overflow button's menu with their own index.
func TestBreadcrumbsFollowLinksFromTheKeyboard(t *testing.T) {
	var selected []int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Before")
		navigation.Breadcrumbs(c, navigation.BreadcrumbsOptions{
			Label: "Path", MaxVisible: 3, OnSelect: func(i int) { selected = append(selected, i) },
			Items: []navigation.BreadcrumbOptions{{Label: "Home"}, {Label: "Docs"}, {Label: "Guides"}, {Label: "Install"}},
		})
	}, 640, 300)
	if a, ok := tt.AccessByRole(ui.RoleText, "Install"); !ok || a.Focusable || a.Current != ui.CurrentPage {
		t.Fatalf("the current crumb is not plain text exposed as the current page: %+v", a)
	}
	if err := tt.Click("Before"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	if a, ok := tt.AccessByRole(ui.RoleLink, "Home"); !ok || !a.Focused {
		t.Fatal("Tab did not reach the first crumb as a link")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyTab)
	if a, ok := tt.AccessByRole(ui.RoleMenuButton, "More pages"); !ok || !a.Focused {
		t.Fatal("Tab did not reach the overflow menu button")
	}
	tt.Key(0, ui.KeyEnter)
	if err := tt.ChooseMenuItem("Guides"); err != nil {
		t.Fatal(err)
	}
	if len(selected) != 2 || selected[0] != 0 || selected[1] != 2 {
		t.Fatalf("selected %v, want [0 2]", selected)
	}
}

// TestUnnamedBreadcrumbsTakeTheLocalesName protects the trail's name: a
// trail without Label is still a named group, so assistive technology can
// tell it from the other groups of the page.
func TestUnnamedBreadcrumbsTakeTheLocalesName(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.Breadcrumbs(c, navigation.BreadcrumbsOptions{
			Items: []navigation.BreadcrumbOptions{{Label: "Home"}, {Label: "Docs"}},
		})
	}, 640, 300)
	if _, ok := tt.AccessByRole(ui.RoleGroup, "Breadcrumbs"); !ok {
		t.Fatal("an unnamed trail is not a group named Breadcrumbs")
	}
}

// TestPaginationHandsTheFocusToTheCurrentPageAtAnEnd protects keyboard paging:
// Next or Previous disabling at an end does not drop the focus, which moves
// to the page reached, exposed as the current page.
func TestPaginationHandsTheFocusToTheCurrentPageAtAnEnd(t *testing.T) {
	page := 2
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Before")
		navigation.Pagination(c, navigation.PaginationOptions{Page: &page, Pages: 3, PreviousLabel: "Previous", NextLabel: "Next"})
	}, 640, 200)
	for _, end := range []struct {
		tabs int
		want string
	}{{5, "3"}, {1, "1"}} {
		page = 2
		if err := tt.Click("Before"); err != nil {
			t.Fatal(err)
		}
		for range end.tabs {
			tt.Key(0, ui.KeyTab)
		}
		tt.Key(0, ui.KeyEnter)
		tt.Frame()
		a, ok := tt.AccessByRole(ui.RoleButton, end.want)
		if !ok || !a.Focused || a.Current != ui.CurrentPage {
			t.Fatalf("page %d reached: the focus is on %q, page %s %+v", page, focusedName(tt), end.want, a)
		}
	}
}

// TestToolbarExposesItsOrientationAndSeparators protects how assistive
// technology reads a toolbar: the way it runs, and its rules as separators
// running across it.
func TestToolbarExposesItsOrientationAndSeparators(t *testing.T) {
	for _, o := range []mygoui.Orientation{mygoui.Horizontal, mygoui.Vertical} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			navigation.Toolbar(c, navigation.ToolbarOptions{Label: "Format", Orientation: o, Items: func() {
				ui.Button(c, "Bold")
				navigation.ToolbarSeparator(c, o)
				ui.Button(c, "Italic")
			}})
		}, 480, 200)
		vertical := o == mygoui.Vertical
		bar, _ := tt.AccessByRole(ui.RoleToolbar, "Format")
		if bar.Vertical != vertical || bar.Horizontal == vertical {
			t.Fatalf("orientation %v: the toolbar reads vertical %v, horizontal %v", o, bar.Vertical, bar.Horizontal)
		}
		var rule *ui.AccessInfo
		for _, a := range tt.AccessAll() {
			if a.Role == ui.RoleSeparator {
				rule = &a
			}
		}
		if rule == nil || rule.Vertical == vertical || rule.Focusable {
			t.Fatalf("orientation %v: the rule reads as %+v", o, rule)
		}
	}
}

// TestHeaderArrowsLeaveTabToWhatLiesBetween protects the accordion pattern in
// Accordion and DisclosureGroup: Up, Down, Home and End move among the
// enabled headers, round the ends, while Tab goes from a header into the
// controls of its open panel.
func TestHeaderArrowsLeaveTabToWhatLiesBetween(t *testing.T) {
	for _, kind := range []string{"accordion", "disclosure group"} {
		expanded := []string{"one"}
		one, two := true, false
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			body := func() { ui.Button(c, "Inside") }
			if kind == "accordion" {
				navigation.Accordion(c, navigation.AccordionOptions{Expanded: &expanded, Items: []navigation.AccordionItemOptions{
					{Key: "one", Title: "One", Body: body}, {Key: "locked", Title: "Locked", Disabled: true},
					{Key: "two", Title: "Two"}, {Key: "three", Title: "Three"},
				}})
				return
			}
			navigation.DisclosureGroup(c, navigation.DisclosureGroupOptions{Multiple: true, Items: []navigation.DisclosureOptions{
				{Title: "One", Open: &one, Body: body}, {Title: "Locked", Disabled: true},
				{Title: "Two", Open: &two}, {Title: "Three"},
			}})
		}, 480, 480)
		tt.Advance(time.Second)
		if err := tt.Click("One"); err != nil {
			t.Fatal(err)
		}
		if err := tt.Click("One"); err != nil { // open again
			t.Fatal(err)
		}
		tt.Advance(time.Second)
		for _, s := range []struct {
			key  ui.Key
			want string
		}{
			{ui.KeyDown, "Two"}, {ui.KeyDown, "Three"}, {ui.KeyDown, "One"}, {ui.KeyUp, "Three"},
			{ui.KeyHome, "One"}, {ui.KeyEnd, "Three"}, {ui.KeyHome, "One"}, {ui.KeyTab, "Inside"},
		} {
			tt.Key(0, s.key)
			if got := focusedName(tt); got != s.want {
				t.Fatalf("%s: after key %v the focus is on %q, want %q", kind, s.key, got, s.want)
			}
		}
	}
}
