package mygoui_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// useMessages publishes m for one test and restores the defaults after it.
func useMessages(t *testing.T, m locale.Messages) {
	t.Helper()
	locale.Set(m)
	t.Cleanup(func() { locale.Set(locale.Messages{}) })
}

// TestLinkIsALinkThatTakesEnterNotSpace pins the link contract: assistive
// technology hears a link, and only Enter follows it from the keyboard.
func TestLinkIsALinkThatTakesEnterNotSpace(t *testing.T) {
	follows := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Link(c, mygoui.LinkOptions{Label: "Docs", OnClick: func() { follows++ }})
	}, 240, 120)
	if _, ok := tt.AccessByRole(ui.RoleLink, "Docs"); !ok {
		t.Fatal("the link is not exposed as a link")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if follows != 0 {
		t.Fatal("Space followed the link")
	}
	tt.Key(0, ui.KeyEnter)
	if follows != 1 {
		t.Fatalf("Enter followed the link %d times, want 1", follows)
	}
}

// TestDefaultNamesFollowTheLocale pins the names the library gives when the
// caller gives none: they are the application's language, not English.
func TestDefaultNamesFollowTheLocale(t *testing.T) {
	useMessages(t, locale.Messages{Close: "Schließen", Sidebar: "Seitenleiste", ToggleSidebar: "Seitenleiste umschalten"})
	selected, collapsed := "", false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Row(c).Fill().Children(func() {
			mygoui.Sidebar(c, mygoui.SidebarOptions{Selected: &selected, Collapsed: &collapsed})
			mygoui.CloseButton(c, mygoui.CloseButtonOptions{})
		})
	}, 480, 240)
	for _, want := range []struct {
		role ui.Role
		name string
	}{
		{ui.RoleButton, "Schließen"},
		{ui.RoleNavigation, "Seitenleiste"},
		{ui.RoleToggleButton, "Seitenleiste umschalten"},
	} {
		if _, ok := tt.AccessByRole(want.role, want.name); !ok {
			t.Errorf("no %v named %q", want.role, want.name)
		}
	}
}

// TestSidebarMarksTheCurrentPageOfANavigationLandmark pins how a screen
// reader finds the sidebar and where the user is in it.
func TestSidebarMarksTheCurrentPageOfANavigationLandmark(t *testing.T) {
	selected, collapsed := "inbox", false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Groups: []mygoui.SidebarGroupOptions{{Items: []mygoui.SidebarItemOptions{
				{ID: "inbox", Label: "Inbox"}, {ID: "sent", Label: "Sent"},
			}}},
		})
	}, 300, 240)
	if _, ok := tt.AccessByRole(ui.RoleNavigation, "Sidebar"); !ok {
		t.Fatal("the sidebar is not a navigation landmark")
	}
	inbox, _ := tt.Access("Inbox")
	sent, _ := tt.Access("Sent")
	if inbox.Current != ui.CurrentPage || sent.Current != ui.CurrentNone {
		t.Fatalf("current = %v and %v, want the page on Inbox only", inbox.Current, sent.Current)
	}
	if inbox.Checked {
		t.Fatal("the current item still reads as a checked button")
	}
}

// TestTitlesAreHeadings pins heading navigation through the titles styled
// as headings.
func TestTitlesAreHeadings(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Children(func() {
			mygoui.CardTitle(c, "Card title")
			mygoui.FrameTitle(c, "Frame title")
			mygoui.CardFrameTitle(c, "Card frame title")
			mygoui.EmptyState(c, mygoui.EmptyStateOptions{Title: "No invoices"})
		})
	}, 400, 400)
	for _, name := range []string{"Card title", "Frame title", "Card frame title", "No invoices"} {
		if a, ok := tt.AccessByRole(ui.RoleHeading, name); !ok || a.Level == 0 {
			t.Errorf("%q is not a heading with a level (%+v)", name, a)
		}
	}
}

// TestListGroupIsAListOfItems pins how a screen reader counts the rows.
func TestListGroupIsAListOfItems(t *testing.T) {
	opened := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ListGroup(c, mygoui.ListGroupOptions{Label: "Account"}, func() {
			mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{Title: "Profile", OnClick: func() { opened++ }})
			mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{Title: "Plan"})
		})
	}, 320, 240)
	if _, ok := tt.AccessByRole(ui.RoleList, "Account"); !ok {
		t.Fatal("the group is not a list")
	}
	items := 0
	for _, a := range tt.AccessAll() {
		if a.Role == ui.RoleListItem {
			items++
		}
	}
	if items != 2 {
		t.Fatalf("list items = %d, want 2", items)
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, "Profile"); !ok {
		t.Fatal("the clickable row is no longer a button")
	}
	if err := tt.Click("Profile"); err != nil || opened != 1 {
		t.Fatalf("click: err=%v opened=%d", err, opened)
	}
}

// chevronApexSide reports whether the chevron drawn in the strip [x0, x1)
// of the row points to the right: its ink at the middle row lies right of
// the centre of its ink.
func chevronApexSide(img *image.RGBA, r ui.Rect, x0, x1 int) (right, found bool) {
	ink := func(x, y int) bool {
		p := img.RGBAAt(x, y)
		return int(p.R)+int(p.G)+int(p.B) < 600
	}
	lo, hi := x1, x0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := x0; x < x1; x++ {
			if ink(x, y) {
				lo, hi = min(lo, x), max(hi, x)
			}
		}
	}
	if lo > hi {
		return false, false
	}
	mid := int(r.Y + r.H/2)
	for x := x0; x < x1; x++ {
		if ink(x, mid) {
			return float32(x) > float32(lo+hi)/2, true
		}
	}
	return false, false
}

// TestListGroupChevronPointsWithTheReadingDirection pins the chevron's
// mirror in right-to-left, as the settings row's.
func TestListGroupChevronPointsWithTheReadingDirection(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			mygoui.ListGroup(c, mygoui.ListGroupOptions{Variant: mygoui.SurfaceTransparent}, func() {
				mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{Title: "Profile", Chevron: true})
			}).Width(320)
		}, 320, 120)
		tt.SetDirection(dir)
		tt.Frame()
		row, _ := tt.AccessByRole(ui.RoleListItem, "Profile")
		x0, x1 := int(row.Bounds.X+row.Bounds.W)-40, int(row.Bounds.X+row.Bounds.W)
		if dir == ui.RightToLeft {
			x0, x1 = int(row.Bounds.X), int(row.Bounds.X)+40
		}
		right, found := chevronApexSide(tt.Image(), row.Bounds, x0, x1)
		if !found {
			t.Fatalf("%v: no chevron at the end of the row", dir)
		}
		if right != (dir == ui.LeftToRight) {
			t.Errorf("%v: the chevron points the wrong way", dir)
		}
	}
}

// TestStatusIsNamedNotOnlyColoured pins the status a mark's colour shows:
// assistive technology hears it too.
func TestStatusIsNamedNotOnlyColoured(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Children(func() {
			mygoui.Timeline(c, mygoui.TimelineOptions{Items: []mygoui.TimelineItem{
				{Title: "Deploy failed", Status: mygoui.ColorDanger},
			}})
			mygoui.Result(c, mygoui.ResultOptions{Status: mygoui.ColorSuccess, Title: "Paid"})
			mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: mygoui.ColorWarning, Label: "Degraded"})
		})
	}, 480, 640)
	for _, name := range []string{"Error", "Success", "Warning"} {
		if _, ok := tt.AccessByRole(ui.RoleImage, name); !ok {
			t.Errorf("no status mark named %q", name)
		}
	}
}

// TestAvatarGroupNamesItsOverflow pins what the "+N" counts.
func TestAvatarGroupNamesItsOverflow(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.AvatarGroup(c, mygoui.AvatarGroupOptions{
			Total:   5,
			Avatars: []mygoui.AvatarOptions{{Name: "Ada Lovelace"}, {Name: "Alan Turing"}},
		}).Label("Reviewers")
	}, 320, 120)
	if _, ok := tt.AccessByRole(ui.RoleGroup, "Reviewers"); !ok {
		t.Fatal("the avatars are not a group")
	}
	if _, ok := tt.Access("3 more"); !ok {
		t.Fatal("the overflow count is not named")
	}
}

// TestDangerMessagesAreAnnouncedAssertively pins the urgency a screen
// reader gives a danger message, and the politeness of the others.
func TestDangerMessagesAreAnnouncedAssertively(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Children(func() {
			mygoui.Alert(c, mygoui.AlertOptions{Title: "Payment failed", Status: mygoui.ColorDanger})
			mygoui.Alert(c, mygoui.AlertOptions{Title: "Saved", Status: mygoui.ColorSuccess})
			mygoui.Banner(c, mygoui.BannerOptions{Text: "Offline", Status: mygoui.ColorDanger})
			mygoui.Banner(c, mygoui.BannerOptions{Text: "Update ready", Status: mygoui.ColorInfo})
		})
	}, 480, 400)
	for _, want := range []struct {
		role ui.Role
		name string
	}{
		{ui.RoleAlert, "Payment failed"}, {ui.RoleStatus, "Saved"},
		{ui.RoleAlert, "Offline"}, {ui.RoleStatus, "Update ready"},
	} {
		if _, ok := tt.AccessByRole(want.role, want.name); !ok {
			t.Errorf("%q is not a %v", want.name, want.role)
		}
	}
}

// TestButtonGroupsAreGroupsAndJoinAToolbar pins the group role, and that a
// toolbar around a group makes its buttons one Tab stop with arrow keys.
func TestButtonGroupsAreGroupsAndJoinAToolbar(t *testing.T) {
	bold := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Toolbar(c, func() {
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{},
				mygoui.ButtonOptions{Label: "Undo"}, mygoui.ButtonOptions{Label: "Redo"}).Label("History")
			mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{},
				mygoui.ToggleButtonOptions{Label: "Bold", Selected: &bold}).Label("Format")
		}).Label("Editor")
		mygoui.Button(c, mygoui.ButtonOptions{Label: "After"})
	}, 480, 160)
	for _, name := range []string{"History", "Format"} {
		if _, ok := tt.AccessByRole(ui.RoleGroup, name); !ok {
			t.Errorf("%q is not a group", name)
		}
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Undo") {
		t.Fatal("Tab does not enter the toolbar")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Redo") {
		t.Fatal("Right does not move along the group")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("the toolbar is more than one Tab stop")
	}
}

// TestHighlightKeepsTheTextWhole pins what assistive technology reads of a
// highlighted paragraph: its text, once.
func TestHighlightKeepsTheTextWhole(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Highlight(c, mygoui.HighlightOptions{Text: "Find the needle here", Queries: []string{"needle"}})
	}, 320, 120)
	if _, ok := tt.Access("Find the needle here"); !ok {
		t.Fatal("the paragraph is not read whole")
	}
}

// TestResizablePanelPagesAndEndsWithoutAMaximum pins the larger step and
// the end of a panel without a maximum: End leaves the size alone.
func TestResizablePanelPagesAndEndsWithoutAMaximum(t *testing.T) {
	size := float32(200)
	o := mygoui.ResizablePanelOptions{Size: &size, Min: 120}
	tt := resizablePanelTester(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyPageUp)
	if size != 264 {
		t.Fatalf("PageUp: size = %v, want 264", size)
	}
	tt.Key(0, ui.KeyPageDown)
	tt.Key(0, ui.KeyPageDown)
	if size != 136 {
		t.Fatalf("PageDown twice: size = %v, want 136", size)
	}
	tt.Key(0, ui.KeyEnd)
	if size != 136 {
		t.Fatalf("End without a maximum: size = %v, want 136", size)
	}
}

// TestScrollShadowRejectsANilState pins the failure of a missing state, as
// ScrollFade and ScrollArea fail.
func TestScrollShadowRejectsANilState(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("a nil scroll state did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		mygoui.ScrollShadow(c, nil, mygoui.ScrollShadowOptions{}, func() {})
	}, 100, 100)
}

// TestSkeletonGroupIsBusy pins the loading state a group of placeholders
// tells assistive technology without a busy container.
func TestSkeletonGroupIsBusy(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SkeletonGroup(c, mygoui.SkeletonGroupOptions{Animation: mygoui.SkeletonNone}, func() {
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 120})
		})
	}, 240, 120)
	if a, ok := tt.AccessByRole(ui.RoleGroup, "Loading"); !ok || !a.Busy {
		t.Fatalf("the placeholders are not a busy group (%+v)", a)
	}
}
