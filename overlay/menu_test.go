package overlay

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
	"testing"
)

func TestDropdownKeyboardChoicesSkipDisabledAndNavigateSubmenus(t *testing.T) {
	var s MenuState
	checked := false
	radio := MenuChoice("small")
	disabledActions := 0
	selected := 0
	items := []MenuItem{
		{Label: "Unavailable", Disabled: true, OnSelect: func() { disabledActions++ }},
		{Label: "Pinned", Kind: MenuCheckbox, Checked: &checked},
		{ID: "size-separator", Kind: MenuSeparator},
		{Label: "Small", Kind: MenuRadio, Radio: &radio, Value: "small"},
		{Label: "Large", Kind: MenuRadio, Radio: &radio, Value: "large"},
		{Label: "Move", Items: []MenuItem{{Label: "Archive", OnSelect: func() { selected++ }}}},
	}
	tt := myui.NewTester(func(c *myui.Context) { theme.Apply(c); DropdownMenu(c, &s, MenuOptions{Label: "Actions"}, items) }, 640, 480)
	tt.Key(0, myui.KeyTab)
	tt.Key(0, myui.KeyDown)
	if !s.Open || !tt.Focused("Pinned") {
		t.Fatal("Down must open and focus first enabled choice")
	}
	tt.Key(0, myui.KeyEnter)
	if !checked || s.Open || disabledActions != 0 {
		t.Fatal("checkbox did not toggle once, or disabled choice ran")
	}
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, myui.KeyL)
	if !tt.Focused("Large") {
		t.Fatal("typeahead did not focus Large")
	}
	tt.Key(0, myui.KeyEnter)
	if radio != "large" {
		t.Fatalf("radio=%q", radio)
	}
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, myui.KeyEnd)
	if !tt.Focused("Move") {
		t.Fatal("End did not reach last enabled item")
	}
	tt.Key(0, myui.KeyRight)
	if !tt.Focused("Archive") {
		t.Fatal("Right did not enter submenu")
	}
	tt.Key(0, myui.KeyEscape)
	if !s.Open || !tt.Focused("Move") {
		t.Fatal("Escape should close only submenu and restore its item")
	}
	tt.Key(0, myui.KeyRight)
	tt.Key(0, myui.KeyEnter)
	if s.Open || selected != 1 {
		t.Fatalf("submenu selection: open=%v actions=%d", s.Open, selected)
	}
}

func TestContextMenuKeyboardCheckboxRadioAndDisabled(t *testing.T) {
	checked := false
	radio := MenuChoice("one")
	count := 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		ContextMenu(myui.Button(c, "Document"), []MenuItem{
			{Label: "Locked", Disabled: true, OnSelect: func() { count++ }},
			{Label: "Pinned", Kind: MenuCheckbox, Checked: &checked},
			{Label: "Size", Items: []MenuItem{{Label: "Two", Kind: MenuRadio, Radio: &radio, Value: "two"}}},
		})
	}, 400, 300)
	tt.Key(0, myui.KeyTab)
	tt.Key(myui.Shift, myui.KeyF10)
	if len(tt.Menu()) != 3 {
		t.Fatalf("Shift+F10 menu=%q", tt.Menu())
	}
	if err := tt.ChooseMenuItem("Locked"); err == nil {
		t.Fatal("disabled native item was selectable")
	}
	if err := tt.ChooseMenuItem("Pinned"); err != nil {
		t.Fatal(err)
	}
	if !checked || count != 0 {
		t.Fatal("native checkbox failed or disabled action ran")
	}
	tt.Key(myui.Shift, myui.KeyF10)
	if err := tt.ChooseMenuItem("Size", "Two"); err != nil {
		t.Fatal(err)
	}
	if radio != "two" {
		t.Fatal("native radio failed")
	}
}

func TestMenubarArrowsSwitchOpenMenus(t *testing.T) {
	var s MenubarState
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		Menubar(c, &s, []BarMenu{{Label: "File", Items: []MenuItem{{Label: "New", OnSelect: func() {}}}},
			{Label: "Disabled", Disabled: true}, {Label: "Edit", Items: []MenuItem{{Label: "Undo", OnSelect: func() {}}}}})
	}, 500, 300)
	tt.Key(0, myui.KeyTab)
	tt.Key(0, myui.KeyDown)
	if !tt.Focused("New") {
		t.Fatal("File menu did not open from Down")
	}
	tt.Key(0, myui.KeyRight)
	if !tt.HasText("Undo") || !tt.Focused("Undo") {
		t.Fatal("Right failed to switch to next enabled menu")
	}
	tt.Key(0, myui.KeyEscape)
	if s.Menu.Open {
		t.Fatal("Escape did not close menubar menu")
	}
}
