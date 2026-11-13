package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestTimelineListsEventsInOrder pins the list: events read top to bottom
// as items of a named list, each with its time beside the title and its
// description below.
func TestTimelineListsEventsInOrder(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Timeline(c, mygoui.TimelineOptions{Label: "Release", Items: []mygoui.TimelineItem{
			{Title: "Build passed", Time: "09:12", Description: "All checks are green.", Status: mygoui.ColorSuccess},
			{Title: "Deployed", Time: "09:20"},
		}})
	}, 360, 240)
	if tt.FramePending() {
		t.Fatal("a timeline without a pending item asks for frames")
	}
	first, time, detail := findRect(t, tt, "Build passed"), findRect(t, tt, "09:12"), findRect(t, tt, "All checks are green.")
	second := findRect(t, tt, "Deployed")
	if time.X <= first.X || time.Y >= detail.Y || detail.Y <= first.Y || second.Y <= detail.Y {
		t.Fatalf("layout: title=%+v time=%+v detail=%+v next=%+v", first, time, detail, second)
	}
	if _, ok := tt.AccessByRole(ui.RoleList, "Release"); !ok {
		t.Fatal("the timeline is not a named list")
	}
	item, ok := tt.AccessByRole(ui.RoleListItem, "Build passed")
	if !ok || item.Description != "All checks are green." {
		t.Fatalf("the event reads as %+v", item)
	}
}

// TestTimelinePendingMarksOnlyTheLastItem pins the pending state: the last
// event alone is busy and shows a loading marker.
func TestTimelinePendingMarksOnlyTheLastItem(t *testing.T) {
	pending := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Timeline(c, mygoui.TimelineOptions{Pending: pending, Items: []mygoui.TimelineItem{
			{Title: "Queued"}, {Title: "Running"},
		}})
	}, 360, 200)
	if _, ok := tt.Access(locale.M().Loading); ok {
		t.Fatal("a settled timeline shows a loading marker")
	}
	pending = true
	tt.Frame()
	if _, ok := tt.Access(locale.M().Loading); !ok {
		t.Fatal("the pending item shows no loading marker")
	}
	first, _ := tt.AccessByRole(ui.RoleListItem, "Queued")
	last, _ := tt.AccessByRole(ui.RoleListItem, "Running")
	if first.Busy || !last.Busy {
		t.Fatalf("busy: first=%v last=%v", first.Busy, last.Busy)
	}
}

// TestTimelineFoldsEarlierItemsBehindAControl pins Recent: older events are
// not shown until the control opens them, and the control opens by keyboard.
func TestTimelineFoldsEarlierItemsBehindAControl(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Timeline(c, mygoui.TimelineOptions{
			Recent: 1, EarlierLabel: "2 earlier",
			Items: []mygoui.TimelineItem{{Title: "Created"}, {Title: "Reviewed"}, {Title: "Merged"}},
		})
	}, 320, 400)
	findRect(t, tt, "Merged")
	if _, ok := tt.Find("Created"); ok {
		t.Fatal("an earlier item shows before the control opens it")
	}
	if err := tt.Click("2 earlier"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(400 * time.Millisecond)
	tt.Frame()
	findRect(t, tt, "Created")
}
