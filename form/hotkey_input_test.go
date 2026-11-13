package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func hotkeyTester(value *Hotkey, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			HotkeyInput(c, value, HotkeyInputOptions{
				AriaLabel: "Shortcut",
				Clearable: true,
				OnChange:  func() { *changes++ },
			})
		})
	}, 360, 160)
}

func TestHotkeyInputRecordsChord(t *testing.T) {
	var value Hotkey
	changes := 0
	tt := hotkeyTester(&value, &changes)
	if err := tt.Click("Shortcut"); err != nil {
		t.Fatal(err)
	}
	// A key without a modifier is no shortcut.
	tt.Key(0, ui.KeyK)
	if value != (Hotkey{}) {
		t.Fatalf("a bare key recorded %+v", value)
	}
	tt.Key(ui.Ctrl|ui.Shift, ui.KeyK)
	want := Hotkey{Mods: ui.Ctrl | ui.Shift, Key: ui.KeyK}
	if value != want || changes != 1 {
		t.Fatalf("recorded %+v with %d changes, want %+v and 1", value, changes, want)
	}
	if !tt.HasText("K") {
		t.Fatalf("the shortcut is not shown as keys: %v", tt.Texts())
	}
	// Recording ended with the chord: the next one is not taken.
	tt.Key(ui.Ctrl, ui.KeyJ)
	if value != want {
		t.Fatalf("a chord after recording replaced the shortcut with %+v", value)
	}
}

func TestHotkeyInputEscapeCancelsRecording(t *testing.T) {
	value := Hotkey{Mods: ui.Alt, Key: ui.KeyP}
	changes := 0
	tt := hotkeyTester(&value, &changes)
	if err := tt.Click("Shortcut"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Key(ui.Ctrl, ui.KeyJ)
	if want := (Hotkey{Mods: ui.Alt, Key: ui.KeyP}); value != want || changes != 0 {
		t.Fatalf("after Escape the shortcut is %+v with %d changes, want it unchanged", value, changes)
	}
}

func TestHotkeyInputBackspaceClears(t *testing.T) {
	value := Hotkey{Mods: ui.Alt, Key: ui.KeyP}
	changes := 0
	tt := hotkeyTester(&value, &changes)
	if err := tt.Click("Shortcut"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyBackspace)
	if value != (Hotkey{}) || changes != 1 {
		t.Fatalf("Backspace left %+v with %d changes, want no shortcut and 1", value, changes)
	}
}

func TestHotkeyKeysNameEveryKeyKind(t *testing.T) {
	got := Hotkey{Mods: ui.Shift, Key: ui.KeyF12}.Keys()
	if len(got) != 2 || got[1] != "F12" {
		t.Fatalf("Shift+F12 keys are %v", got)
	}
	if got := (Hotkey{Mods: ui.Ctrl, Key: ui.KeyZ}).Keys(); got[len(got)-1] != "Z" {
		t.Fatalf("Ctrl+Z keys are %v", got)
	}
	if got := (Hotkey{Mods: ui.Ctrl, Key: ui.Key9}).Keys(); got[len(got)-1] != "9" {
		t.Fatalf("Ctrl+9 keys are %v", got)
	}
}

func TestHotkeyInputExposesRecording(t *testing.T) {
	value := Hotkey{Mods: ui.Ctrl, Key: ui.KeyK}
	changes := 0
	tt := hotkeyTester(&value, &changes)
	formTextClick(t, tt, "Shortcut")
	if info, ok := tt.AccessByRole(ui.RoleToggleButton, "Shortcut"); !ok || !info.Checked {
		t.Fatalf("recording: found %v, pressed %v", ok, info.Checked)
	}
	tt.Key(0, ui.KeyEscape)
	if info, _ := tt.AccessByRole(ui.RoleToggleButton, "Shortcut"); info.Checked {
		t.Fatal("the recorder still reads as pressed after Escape")
	}
}
