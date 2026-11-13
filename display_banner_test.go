package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestBannerSpansItsPaneAndRunsItsControls pins the bar: it runs from edge
// to edge as a status, its action works, and the dismiss button exists only
// with OnDismiss and reports the press to the caller.
func TestBannerSpansItsPaneAndRunsItsControls(t *testing.T) {
	restarts, dismissed := 0, 0
	dismissible := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		o := mygoui.BannerOptions{
			Text: "An update is ready.", Status: mygoui.ColorInfo, ShowIcon: true,
			Action: func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Restart", OnClick: func() { restarts++ }})
			},
		}
		if dismissible {
			o.OnDismiss = func() { dismissed++ }
		}
		mygoui.Banner(c, o)
	}, 480, 120)
	bar, ok := tt.AccessByRole(ui.RoleStatus, "An update is ready.")
	if !ok || bar.Bounds.X != 0 || bar.Bounds.Y != 0 || bar.Bounds.W != 480 {
		t.Fatalf("the banner reads as %+v (found=%v), want a status across the pane", bar, ok)
	}
	if text, action := findRect(t, tt, "An update is ready."), findRect(t, tt, "Restart"); action.X <= text.X {
		t.Fatalf("the action is not after the text: %+v %+v", text, action)
	}
	if err := tt.Click("Restart"); err != nil || restarts != 1 {
		t.Fatalf("action: err=%v restarts=%d", err, restarts)
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, locale.M().Close); ok {
		t.Fatal("a banner without OnDismiss shows a dismiss button")
	}
	dismissible = true
	tt.Frame()
	if err := tt.Click(locale.M().Close); err != nil || dismissed != 1 {
		t.Fatalf("dismiss: err=%v dismissed=%d", err, dismissed)
	}
}

// TestBannerKeepsItsTextReadableWhenNarrow pins the narrow bar: the action
// goes below the text instead of crushing it to a letter a line.
func TestBannerKeepsItsTextReadableWhenNarrow(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Banner(c, mygoui.BannerOptions{
			Text: "Update ready", ShowIcon: true, OnDismiss: func() {},
			Action: func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Restart now"}) },
		})
	}, 240, 200)
	bar, _ := tt.AccessByRole(ui.RoleStatus, "Update ready")
	action := findRect(t, tt, "Restart now")
	if action.X > bar.Bounds.W/4 || action.Y < 24 {
		t.Fatalf("bar %+v, action %+v: want the action under the text", bar.Bounds, action)
	}
}
