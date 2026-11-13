package navigation

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestStatusBarAnnouncesMessageChanges(t *testing.T) {
	message := ""
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		StatusBar(c, StatusBarOptions{Label: "Status", Message: message})
	}, 480, 60)
	tt.UseScreenReader()
	message = "Saved 2 files"
	tt.Frame()
	if got := tt.Announcements(); !slices.Equal(got, []string{"Saved 2 files"}) {
		t.Fatalf("read out %q", got)
	}
	if status, ok := tt.AccessByRole(ui.RoleStatus, "Saved 2 files"); !ok || status.Bounds.H > 24 {
		t.Fatalf("status region: found %v, %+v", ok, status.Bounds)
	}
}

func TestStatusBarItemsPressOnlyWithAnAction(t *testing.T) {
	presses := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		StatusBar(c, StatusBarOptions{Message: "Ready", Items: []StatusBarItemOptions{
			{Text: "UTF-8"},
			{Separator: true},
			{Icon: lucide.GitBranch(), Text: "main", Tooltip: "Switch branch", OnClick: func() { presses++ }},
			{Icon: lucide.Bell(), Label: "Notifications", Text: "3", OnClick: func() { presses += 10 }},
		}})
	}, 480, 60)
	if _, ok := tt.AccessByRole(ui.RoleButton, "UTF-8"); ok {
		t.Fatal("an item without an action is a button")
	}
	branch, ok := tt.AccessByRole(ui.RoleButton, "main")
	if !ok || branch.Description != "Switch branch" {
		t.Fatalf("branch item: found %v, description %q", ok, branch.Description)
	}
	if err := tt.Click("main"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if presses != 11 || !tt.Focused("Notifications, 3") {
		t.Fatalf("presses %d, Notifications focused %v", presses, tt.Focused("Notifications, 3"))
	}
}

// A narrow bar keeps its message readable and gives up the items of the
// lowest priority first, leaving no rule beside a missing item.
func TestStatusBarDropsLowPriorityItemsAndKeepsTheMessage(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c)
		StatusBar(c, StatusBarOptions{Message: "Saved 2 files", Items: []StatusBarItemOptions{
			{Text: "Ln 12, Col 4"},
			{Text: "UTF-8"},
			{Separator: true},
			{Icon: lucide.GitBranch(), Text: "feature/status", Priority: 2, OnClick: func() {}},
			{Icon: lucide.Bell(), Label: "Notifications", Priority: 3, OnClick: func() {}},
		}})
	}
	wide := ui.NewTester(view, 480, 40)
	for _, name := range []string{"Ln 12, Col 4", "UTF-8", "feature/status", "Notifications"} {
		if _, ok := wide.Find(name); !ok {
			t.Errorf("%s is missing at 480 wide", name)
		}
	}
	narrow := ui.NewTester(view, 260, 40)
	narrow.Frame()
	if _, ok := narrow.Find("Notifications"); !ok {
		t.Error("the highest priority item left")
	}
	if _, ok := narrow.Find("Ln 12, Col 4"); ok {
		t.Error("a lowest priority item stayed in 260 DIP")
	}
	msg, ok := narrow.Find("Saved 2 files")
	if !ok || msg.W < 90 {
		t.Errorf("the message has %v DIP (found %v)", msg.W, ok)
	}
	for _, name := range []string{"feature/status", "Notifications"} {
		if r, ok := narrow.Find(name); ok && r.X+r.W > 260 {
			t.Errorf("%s runs past the bar: %+v", name, r)
		}
	}
}

// A item that leaves the bar must not hand its focus to the one after it.
func TestStatusBarFocusStaysWithItsItemWhenAnEarlierOneLeaves(t *testing.T) {
	items := []StatusBarItemOptions{
		{Label: "Branch", Text: "main", OnClick: func() {}},
		{Label: "Encoding", Text: "UTF-8", OnClick: func() {}},
		{Label: "Alerts", Text: "3", OnClick: func() {}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		StatusBar(c, StatusBarOptions{Items: items})
	}, 480, 60)
	if err := tt.Click("Encoding, UTF-8"); err != nil {
		t.Fatal(err)
	}
	items = items[1:]
	tt.Frame()
	if !tt.Focused("Encoding, UTF-8") {
		t.Fatal("the focus left its item when an earlier one was removed")
	}
}
