package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestSettingsRowTitleNamesItsControl pins the row's naming contract: a
// screen reader hears the row title on the trailing switch, and the title
// is a click target for it.
func TestSettingsRowTitleNamesItsControl(t *testing.T) {
	var wifi ui.Handle
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{
			Title: "Network",
			Rows: []mygoui.SettingsRow{{
				Title: "Wi-Fi", Subtitle: "Join known networks", Control: &wifi,
				Trailing: func() { form.Switch(c, form.SwitchOptions{On: &on, Handle: &wifi}) },
			}},
		})
	}, 480, 240)
	info, ok := tt.AccessByRole(ui.RoleSwitch, "Wi-Fi")
	if !ok {
		t.Fatalf("no switch named by the row title: %+v", tt.AccessAll())
	}
	if info.Description != "Join known networks" {
		t.Fatalf("switch description = %q", info.Description)
	}
	if _, ok := tt.AccessByRole(ui.RoleGroup, "Network"); !ok {
		t.Fatal("the group is not named by its title")
	}
	title := findRect(t, tt, "Wi-Fi")
	tt.ClickAt(title.X+2, title.Y+title.H/2)
	tt.Frame()
	if !on {
		t.Fatal("a click on the row title did not toggle its switch")
	}
}

// TestSettingsPressableRowActivates pins the navigation row: the whole row
// is one button for pointer and keyboard, and a disabled row is inert.
func TestSettingsPressableRowActivates(t *testing.T) {
	opened, blocked := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{Rows: []mygoui.SettingsRow{
			{Title: "Account", Subtitle: "Signed in", OnClick: func() { opened++ }},
			{Title: "Billing", Disabled: true, OnClick: func() { blocked++ }},
		}})
	}, 480, 240)
	row, ok := tt.AccessByRole(ui.RoleButton, "Account")
	if !ok {
		t.Fatalf("no button row: %+v", tt.AccessAll())
	}
	// The far end of the row, away from its texts, still activates it.
	tt.ClickAt(row.Bounds.X+row.Bounds.W-8, row.Bounds.Y+row.Bounds.H/2)
	tt.Frame()
	if opened != 1 {
		t.Fatalf("row click fired %d times", opened)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Account") {
		t.Fatal("Tab did not focus the row")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if opened != 2 {
		t.Fatalf("Enter fired %d times in total", opened)
	}
	billing := findRect(t, tt, "Billing")
	tt.ClickAt(billing.X+2, billing.Y+billing.H/2)
	tt.Frame()
	if blocked != 0 {
		t.Fatal("a disabled row activated")
	}
}

// TestSettingsExpanderRowRevealsNestedRows pins the expander: nested rows
// are absent while closed, and the caller-owned flag follows the row.
func TestSettingsExpanderRowRevealsNestedRows(t *testing.T) {
	expanded := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{Rows: []mygoui.SettingsRow{{
			Title: "Advanced", Expanded: &expanded,
			Rows: []mygoui.SettingsRow{{Title: "Proxy"}, {Title: "DNS"}},
		}}})
	}, 480, 320)
	if tt.HasText("Proxy") {
		t.Fatal("nested rows show while collapsed")
	}
	if info, ok := tt.Access("Advanced"); !ok || info.Expanded {
		t.Fatalf("collapsed expander access = %+v, found %v", info, ok)
	}
	if err := tt.Click("Advanced"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if !expanded {
		t.Fatal("the caller's flag did not follow the expander")
	}
	if !tt.HasText("Proxy") || !tt.HasText("DNS") {
		t.Fatalf("expanded rows = %q", tt.Texts())
	}
	if info, _ := tt.Access("Advanced"); !info.Expanded {
		t.Fatalf("expanded expander access = %+v", info)
	}
}

// TestSettingsGroupLoadingFirstLoadAndRefresh pins both treatments: a group
// with no rows keeps its title over skeleton rows; one with rows keeps them
// readable but takes no click until the load ends.
func TestSettingsGroupLoadingFirstLoadAndRefresh(t *testing.T) {
	clicks := 0
	o := mygoui.SettingsGroupOptions{Title: "Network", Loading: true}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SettingsGroup(c, o)
	}, 480, 360)
	group, ok := tt.AccessByRole(ui.RoleGroup, "Network")
	if !ok || !group.Busy {
		t.Fatalf("a loading group is not busy: %+v", group)
	}
	findRect(t, tt, "Network")
	three := group.Bounds.H

	o.LoadingRows = 1
	tt.Frame()
	if group, _ = tt.AccessByRole(ui.RoleGroup, "Network"); group.Bounds.H >= three {
		t.Fatalf("LoadingRows is ignored: %v rows high vs %v", group.Bounds.H, three)
	}

	o.Rows = []mygoui.SettingsRow{{Title: "Account", OnClick: func() { clicks++ }}}
	tt.Frame()
	findRect(t, tt, "Account")
	if err := tt.Click("Account"); err == nil || clicks != 0 {
		t.Fatalf("a refreshing row took a click (err=%v, clicks=%d)", err, clicks)
	}
	o.Loading = false
	tt.Frame()
	if err := tt.Click("Account"); err != nil || clicks != 1 {
		t.Fatalf("a loaded row did not take a click (err=%v, clicks=%d)", err, clicks)
	}
}

// TestSettingsRowControlWrapsUnderItsTextWhenNarrow pins the narrow card:
// the control goes below the title, which keeps its width.
func TestSettingsRowControlWrapsUnderItsTextWhenNarrow(t *testing.T) {
	var h ui.Handle
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{Rows: []mygoui.SettingsRow{{
			Title: "Join known networks", Control: &h,
			Trailing: func() { form.Switch(c, form.SwitchOptions{On: &on, Handle: &h, Label: "x"}) },
		}}})
	}, 200, 240)
	title := findRect(t, tt, "Join known networks")
	sw, _ := tt.AccessByRole(ui.RoleSwitch, "Join known networks")
	if title.W < 100 || sw.Bounds.Y < title.Y+title.H || sw.Bounds.X+sw.Bounds.W > 200 {
		t.Fatalf("title %+v, switch %+v", title, sw.Bounds)
	}
}
