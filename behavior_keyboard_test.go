package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// A tag group is one Tab stop: arrows walk its tags and remove buttons,
// Space selects, Backspace removes, and Tab leaves the whole group.
func TestTagGroupIsOneTabStopWalkedByArrows(t *testing.T) {
	state := &mygoui.TagGroupState{}
	removed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.TagGroup(c, mygoui.TagGroupOptions{Label: "Topics", Selection: mygoui.TagSelectMultiple, State: state},
			[]mygoui.TagOptions{
				{Key: "go", Label: "Go", Removable: true, OnRemove: func() { removed++ }},
				{Key: "ui", Label: "UI"},
			})
		mygoui.Button(c, mygoui.ButtonOptions{Label: "After"})
	}, 360, 200)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Go") {
		t.Fatal("Tab did not reach the first tag")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Remove") {
		t.Fatal("Right did not reach the tag's remove button")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("UI") {
		t.Fatal("Right did not reach the next tag")
	}
	tt.Key(0, ui.KeySpace)
	if !state.IsSelected("ui") {
		t.Fatalf("Space did not select: %v", state.Selected)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyBackspace)
	if removed != 1 {
		t.Fatalf("Backspace removed %d tags, want 1", removed)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab did not leave the group")
	}
}

// A settings expander row opens and closes from the keyboard.
func TestSettingsExpanderRowTogglesFromTheKeyboard(t *testing.T) {
	expanded := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{Rows: []mygoui.SettingsRow{{
			Title: "Advanced", Expanded: &expanded, Rows: []mygoui.SettingsRow{{Title: "Proxy"}},
		}}})
	}, 480, 320)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Advanced") {
		t.Fatal("Tab did not reach the expander")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if !expanded || !tt.HasText("Proxy") {
		t.Fatalf("Enter: expanded=%v texts=%q", expanded, tt.Texts())
	}
	tt.Key(0, ui.KeySpace)
	tt.Advance(time.Second)
	if expanded || tt.HasText("Proxy") {
		t.Fatalf("Space: expanded=%v texts=%q", expanded, tt.Texts())
	}
}

// The control folding a timeline's earlier events is a disclosure that
// opens from the keyboard.
func TestTimelineEarlierControlOpensByKeyboard(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Timeline(c, mygoui.TimelineOptions{
			Recent: 1, EarlierLabel: "2 earlier",
			Items: []mygoui.TimelineItem{{Title: "Created"}, {Title: "Reviewed"}, {Title: "Merged"}},
		})
	}, 320, 400)
	if info, ok := tt.AccessByRole(ui.RoleDisclosure, "2 earlier"); !ok || info.Expanded {
		t.Fatalf("closed control reads as %+v (found %v)", info, ok)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	findRect(t, tt, "Created")
	if info, _ := tt.AccessByRole(ui.RoleDisclosure, "2 earlier"); !info.Expanded {
		t.Fatalf("open control reads as %+v", info)
	}
}

// A card with OnClick is one named button: Tab reaches it and Enter or
// Space activates it.
func TestClickableCardIsAKeyboardButton(t *testing.T) {
	opened := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Card(c, mygoui.CardOptions{Label: "Pro plan", OnClick: func() { opened++ }}, func() {
			mygoui.CardTitle(c, "Pro")
		})
	}, 360, 240)
	if _, ok := tt.AccessByRole(ui.RoleButton, "Pro plan"); !ok {
		t.Fatalf("no button named by the card label: %+v", tt.AccessAll())
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Pro plan") {
		t.Fatal("Tab did not reach the card")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if opened != 2 {
		t.Fatalf("Enter and Space activated %d times, want 2", opened)
	}
}
