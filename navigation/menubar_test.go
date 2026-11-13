package navigation_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func menubarMenus(chosen *string) []navigation.MenubarMenuOptions {
	menu := func(items ...navigation.MenuItemOptions) navigation.MenuOptions {
		return navigation.MenuOptions{
			Sections: []navigation.MenuSectionOptions{{Items: items}},
			OnSelect: func(key string) { *chosen = key },
		}
	}
	return []navigation.MenubarMenuOptions{
		{Label: "File", Menu: menu(
			navigation.MenuItemOptions{Label: "New", Shortcut: "⌘N"},
			navigation.MenuItemOptions{Label: "Recent", Submenu: []navigation.MenuItemOptions{{Label: "notes.txt"}}},
		)},
		{Label: "Edit", Menu: menu(navigation.MenuItemOptions{Label: "Undo"})},
		{Label: "Tools", Disabled: true, Menu: menu(navigation.MenuItemOptions{Label: "Lint"})},
		{Label: "View", Menu: menu(navigation.MenuItemOptions{Label: "Zoom In"})},
	}
}

func menubarTester(open, chosen *string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.Menubar(c, navigation.MenubarOptions{Label: "Application", Open: open, Menus: menubarMenus(chosen)})
	}, 480, 320)
}

func TestMenubarArrowsSwitchMenusAndEscapeRestoresFocus(t *testing.T) {
	var open, chosen string
	tt := menubarTester(&open, &chosen)
	if err := tt.Click("File"); err != nil {
		t.Fatal(err)
	}
	if open != "File" || !tt.HasText("New") || !tt.Focused("New") {
		t.Fatalf("open %q, New shown %v focused %v", open, tt.HasText("New"), tt.Focused("New"))
	}
	tt.Key(0, ui.KeyRight)
	if open != "Edit" || tt.HasText("New") || !tt.Focused("Undo") {
		t.Fatalf("Right: open %q, New shown %v, Undo focused %v", open, tt.HasText("New"), tt.Focused("Undo"))
	}
	tt.Key(0, ui.KeyRight) // over the disabled Tools
	if open != "View" {
		t.Fatalf("Right past a disabled menu opened %q", open)
	}
	tt.Key(0, ui.KeyRight) // round to File
	tt.Key(0, ui.KeyLeft)
	if open != "View" {
		t.Fatalf("Left from the first menu opened %q", open)
	}
	tt.Key(0, ui.KeyEscape)
	if open != "" || tt.HasText("Zoom In") || !tt.Focused("View") {
		t.Fatalf("Escape: open %q, panel shown %v, title focused %v", open, tt.HasText("Zoom In"), tt.Focused("View"))
	}
}

func TestMenubarSubmenuKeepsItsArrows(t *testing.T) {
	var open, chosen string
	tt := menubarTester(&open, &chosen)
	if err := tt.Click("File"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown) // Recent
	tt.Key(0, ui.KeyRight)
	if open != "File" || !tt.HasText("notes.txt") {
		t.Fatalf("Right on a submenu item: open %q, submenu shown %v", open, tt.HasText("notes.txt"))
	}
	tt.Key(0, ui.KeyLeft)
	if open != "File" || tt.HasText("notes.txt") {
		t.Fatalf("Left in a submenu: open %q, submenu shown %v", open, tt.HasText("notes.txt"))
	}
}

func TestMenubarHoverSwitchesAnOpenMenuAndChoosingCloses(t *testing.T) {
	var open, chosen string
	tt := menubarTester(&open, &chosen)
	edit, _ := tt.Find("Edit")
	tt.Move(edit.X+edit.W/2, edit.Y+edit.H/2)
	if open != "" {
		t.Fatalf("hover opened %q with no menu showing", open)
	}
	if err := tt.Click("File"); err != nil {
		t.Fatal(err)
	}
	tt.Move(edit.X+edit.W/2, edit.Y+edit.H/2)
	tt.Frame()
	if open != "Edit" || !tt.HasText("Undo") || tt.HasText("New") {
		t.Fatalf("hover: open %q, Undo %v, New %v", open, tt.HasText("Undo"), tt.HasText("New"))
	}
	if err := tt.Click("Undo"); err != nil {
		t.Fatal(err)
	}
	if chosen != "Undo" || open != "" || tt.HasText("Undo") {
		t.Fatalf("choice %q, open %q, panel shown %v", chosen, open, tt.HasText("Undo"))
	}
}

func TestMenubarRolesForAssistiveTechnology(t *testing.T) {
	var open, chosen string
	tt := menubarTester(&open, &chosen)
	if _, ok := tt.AccessByRole(ui.RoleMenuBar, "Application"); !ok {
		t.Fatal("no menu bar named Application")
	}
	if err := tt.Click("File"); err != nil {
		t.Fatal(err)
	}
	file, ok := tt.AccessByRole(ui.RoleMenuItem, "File")
	if !ok || !file.Expanded {
		t.Fatalf("File title: %+v, found %v", file, ok)
	}
	if tools, _ := tt.AccessByRole(ui.RoleMenuItem, "Tools"); !tools.Disabled {
		t.Fatalf("Tools title: %+v", tools)
	}
}
