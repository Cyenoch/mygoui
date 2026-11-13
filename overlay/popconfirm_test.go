package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func popconfirmTester(open *bool, confirmed, cancelled *int) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(160).AlignItems(ui.Start).Children(func() {
			anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Delete", OnClick: func() { *open = true }})
			Popconfirm(c, anchor, open, PopconfirmOptions{
				Title: "Delete this file?", ConfirmLabel: "Yes, delete", Destructive: true,
				OnConfirm: func() { *confirmed++ },
				OnCancel:  func() { *cancelled++ },
			})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Elsewhere"})
		})
	}, 480, 360)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return tt
}

// A confirmation that took Enter as "yes" would run the destructive action
// for a user who only pressed the key twice.
func TestPopconfirmFocusesCancelSoEnterDoesNotConfirm(t *testing.T) {
	var open bool
	var confirmed, cancelled int
	tt := popconfirmTester(&open, &confirmed, &cancelled)
	if err := tt.Click("Delete"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !open || !tt.Focused("Cancel") {
		t.Fatalf("open=%v, Cancel focused=%v", open, tt.Focused("Cancel"))
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if open || confirmed != 0 || cancelled != 1 {
		t.Fatalf("Enter: open=%v confirmed=%d cancelled=%d", open, confirmed, cancelled)
	}
	if !tt.Focused("Delete") {
		t.Fatal("focus did not return to the anchor")
	}
}

func TestPopconfirmConfirmsOnceAndDismissalCancels(t *testing.T) {
	var open bool
	var confirmed, cancelled int
	tt := popconfirmTester(&open, &confirmed, &cancelled)

	if err := tt.Click("Delete"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.Click("Yes, delete"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if open || confirmed != 1 || cancelled != 0 {
		t.Fatalf("confirm: open=%v confirmed=%d cancelled=%d", open, confirmed, cancelled)
	}

	for name, dismiss := range map[string]func(){
		"Escape":        func() { tt.Key(0, ui.KeyEscape) },
		"outside press": func() { _ = tt.Click("Elsewhere") },
	} {
		confirmed, cancelled = 0, 0
		if err := tt.Click("Delete"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		dismiss()
		tt.Frame()
		if open || confirmed != 0 || cancelled != 1 {
			t.Fatalf("%s: open=%v confirmed=%d cancelled=%d", name, open, confirmed, cancelled)
		}
	}
}

func TestPopconfirmFocusesConfirmWhenNotDestructive(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Publish", OnClick: func() { open = true }})
		Popconfirm(c, anchor, &open, PopconfirmOptions{Title: "Publish?"})
	}, 480, 360)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Publish"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.Focused("Confirm") {
		t.Fatal("a plain confirmation does not put the focus on Confirm")
	}
}

// Closing the bubble while the confirmed work runs would hide the only sign
// that it is still running, and a second press would start it again.
func TestPopconfirmStaysOpenWhileLoading(t *testing.T) {
	var open, loading bool
	var cancelled int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(160).AlignItems(ui.Start).Children(func() {
			anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Delete", OnClick: func() { open = true }})
			Popconfirm(c, anchor, &open, PopconfirmOptions{
				Title: "Delete?", Destructive: true, Loading: loading,
				OnConfirm: func() { loading, open = true, true },
				OnCancel:  func() { cancelled++ },
			})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Elsewhere"})
		})
	}, 480, 360)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Delete"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.Click("Confirm"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	for name, dismiss := range map[string]func(){
		"Escape":        func() { tt.Key(0, ui.KeyEscape) },
		"outside press": func() { _ = tt.Click("Elsewhere") },
	} {
		dismiss()
		tt.Frame()
		tt.Frame()
		if !open || cancelled != 0 || !tt.HasText("Delete?") {
			t.Fatalf("%s while loading: open=%v cancelled=%d", name, open, cancelled)
		}
	}
	loading = false
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || cancelled != 1 {
		t.Fatalf("Escape after loading: open=%v cancelled=%d", open, cancelled)
	}
}
