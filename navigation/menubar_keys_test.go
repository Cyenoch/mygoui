package navigation

import (
	"fmt"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func keyedMenubar(open, chosen *string) *ui.Tester {
	menu := func(items ...MenuItemOptions) MenuOptions {
		return MenuOptions{Sections: []MenuSectionOptions{{Items: items}}, OnSelect: func(key string) { *chosen = key }}
	}
	return ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		Menubar(c, MenubarOptions{Label: "Application", Open: open, Menus: []MenubarMenuOptions{
			{Label: "File", Menu: menu(
				MenuItemOptions{Label: "New", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyN}},
				MenuItemOptions{Label: "Locked", Disabled: true, Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyL}},
			)},
			{Label: "Edit", Menu: menu(MenuItemOptions{Label: "Undo"})},
			{Label: "Format", Menu: menu(MenuItemOptions{Label: "Bold"})},
		}})
	}, 480, 320)
}

func TestMenubarAcceleratorActivatesItsItemWithTheMenuClosed(t *testing.T) {
	var open, chosen string
	tt := keyedMenubar(&open, &chosen)
	tt.Key(ui.Cmd, ui.KeyL)
	if chosen != "" {
		t.Fatalf("a disabled item's accelerator chose %q", chosen)
	}
	tt.Key(ui.Cmd, ui.KeyN)
	if chosen != "New" || open != "" {
		t.Fatalf("chosen %q, open %q", chosen, open)
	}
	if err := tt.Click("File"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(menuItemShortcut(MenuItemOptions{Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyN}})) {
		t.Fatal("the row does not show its accelerator")
	}
}

func TestMenubarMnemonicsAndF10FollowWindowsConventions(t *testing.T) {
	was := mnemonicsOn
	mnemonicsOn = true
	defer func() { mnemonicsOn = was }()

	var open, chosen string
	tt := keyedMenubar(&open, &chosen)
	// Format cannot have F, which File has: it takes the next free letter.
	tt.Key(ui.Alt, ui.KeyE)
	if open != "Edit" || !tt.Focused("Undo") {
		t.Fatalf("Alt+E: open %q, Undo focused %v", open, tt.Focused("Undo"))
	}
	tt.Key(ui.Alt, ui.KeyO)
	if open != "Format" {
		t.Fatalf("Alt+O: open %q", open)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyEscape)
	open = ""
	tt.Frame()

	tt.Key(0, ui.KeyF10)
	if !tt.Focused("File") || open != "" {
		t.Fatalf("F10: File focused %v, open %q", tt.Focused("File"), open)
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Edit") {
		t.Fatal("Right after F10 did not move along the bar")
	}
}

func TestMenubarAltTapFocusesTheBarButAltWithAKeyDoesNot(t *testing.T) {
	was := mnemonicsOn
	mnemonicsOn = true
	defer func() { mnemonicsOn = was }()

	var open, chosen string
	tt := keyedMenubar(&open, &chosen)
	tt.HoldModifiers(ui.Alt)
	tt.HoldModifiers(0)
	if !tt.Focused("File") {
		t.Fatal("a tap of Alt did not focus the bar")
	}
	open = ""
	tt = keyedMenubar(&open, &chosen)
	tt.HoldModifiers(ui.Alt)
	tt.Key(ui.Alt, ui.KeyE)
	tt.Key(0, ui.KeyEscape)
	tt.HoldModifiers(0)
	if open != "" || tt.Focused("File") {
		t.Fatalf("Alt used with a key still focused the bar: open %q", open)
	}
}

func TestMenuTypeAheadFocusesTheItemTypedAndSkipsDisabledOnes(t *testing.T) {
	var open string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		Menubar(c, MenubarOptions{Open: &open, Menus: []MenubarMenuOptions{{Label: "File", Menu: MenuOptions{
			Sections: []MenuSectionOptions{{Items: []MenuItemOptions{
				{Label: "New"}, {Label: "Open"}, {Label: "Paste", Disabled: true}, {Label: "Print"}, {Label: "Quit"},
			}}},
		}}}})
	}, 480, 320)
	if err := tt.Click("File"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	tt.TypeKey(0, ui.KeyP, "p")
	if !tt.Focused("Print") {
		t.Fatal("p did not reach Print, past the disabled Paste")
	}
	tt.TypeKey(0, ui.KeyP, "p")
	if !tt.Focused("Print") {
		t.Fatal("p again left the only enabled P item")
	}
	tt.TypeKey(0, ui.KeyQ, "q")
	if !tt.Focused("Print") {
		t.Fatal("a letter right after another took a new prefix: pq matches nothing, so the focus stays")
	}
	tt.Advance(2 * time.Second)
	tt.TypeKey(0, ui.KeyQ, "q")
	if !tt.Focused("Quit") {
		t.Fatal("q after a pause did not reach Quit")
	}
}

// Alt+Shift switches the keyboard layout on Windows; it is not a tap of Alt.
func TestMenubarAltWithAnotherModifierDoesNotFocusTheBar(t *testing.T) {
	was := mnemonicsOn
	mnemonicsOn = true
	defer func() { mnemonicsOn = was }()

	var open, chosen string
	tt := keyedMenubar(&open, &chosen)
	tt.HoldModifiers(ui.Alt)
	tt.HoldModifiers(ui.Alt | ui.Shift)
	tt.HoldModifiers(ui.Alt)
	tt.HoldModifiers(0)
	if tt.Focused("File") {
		t.Fatal("Alt+Shift focused the bar")
	}
}

// TestLongMenuTypeAheadReachesRowsNotYetBuilt protects type-ahead in a menu
// long enough to build only the rows in view: typing reaches an item far
// below, which scrolls into view with the focus.
func TestLongMenuTypeAheadReachesRowsNotYetBuilt(t *testing.T) {
	var items []MenuItemOptions
	for i := range menuListRows + 8 {
		items = append(items, MenuItemOptions{Label: fmt.Sprintf("Item %02d", i)})
	}
	items = append(items, MenuItemOptions{Label: "Zoom"})
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		Dropdown(c, DropdownOptions{Label: "Actions", Open: &open, Menu: MenuOptions{
			Sections: []MenuSectionOptions{{Items: items}},
		}})
	}, 480, 480)
	tt.Frame()
	if tt.HasText("Zoom") {
		t.Fatal("the fixture builds every row: the menu is not long")
	}
	tt.TypeKey(0, ui.KeyZ, "z")
	tt.Frame()
	if !tt.Focused("Zoom") {
		t.Fatal("z did not reach Zoom below the rows built")
	}
}
