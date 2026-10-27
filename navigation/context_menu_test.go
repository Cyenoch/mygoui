package navigation

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestContextMenuOpensAtNativeTriggerAndActivates(t *testing.T) {
	var chosen []string
	tt := ui.NewTester(func(c *ui.Context) {
		anchor := ui.ButtonBase(c).Label("Document").Size(140, 60)
		anchor.Children(func() { ui.Text(c, "Document") })
		ContextMenu(anchor, MenuOptions{
			Sections: []MenuSectionOptions{{Items: []MenuItemOptions{
				{Key: "rename", Label: "Rename"},
				{Label: "Move", Submenu: []MenuItemOptions{{Key: "folder", Label: "Folder"}}},
			}}},
			OnSelect: func(key string) { chosen = append(chosen, key) },
		})
	}, 320, 200)
	if err := tt.RightClick("Document"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("Move", "Folder"); err != nil {
		t.Fatal(err)
	}
	if len(chosen) != 1 || chosen[0] != "folder" {
		t.Fatalf("right-click selection: %v", chosen)
	}
	if err := tt.Click("Document"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyContextMenu)
	if err := tt.ChooseMenuItem("Rename"); err != nil {
		t.Fatal(err)
	}
	if len(chosen) != 2 || chosen[1] != "rename" {
		t.Fatalf("keyboard selection: %v", chosen)
	}
}

func TestContextMenuTogglesCallerSelection(t *testing.T) {
	var selected []string
	tt := ui.NewTester(func(c *ui.Context) {
		anchor := ui.ButtonBase(c).Label("Document").Size(140, 60)
		ContextMenu(anchor, MenuOptions{
			Selection: MenuSelectMultiple, Selected: &selected,
			Sections: []MenuSectionOptions{{Items: []MenuItemOptions{
				{Key: "pinned", Label: "Pinned"},
			}}},
		})
	}, 320, 200)
	for _, wantSelected := range []bool{true, false} {
		if err := tt.RightClick("Document"); err != nil {
			t.Fatal(err)
		}
		if err := tt.ChooseMenuItem("Pinned"); err != nil {
			t.Fatal(err)
		}
		if got := len(selected) == 1 && selected[0] == "pinned"; got != wantSelected {
			t.Fatalf("selection %v, want pinned %v", selected, wantSelected)
		}
	}
}

func TestContextMenuOnInteractiveChild(t *testing.T) {
	var primary, secondary, chosen int
	tt := ui.NewTester(func(c *ui.Context) {
		anchor := ui.Column(c).Size(180, 80)
		anchor.Children(func() {
			ui.ButtonBase(c).Label("Child").Size(100, 40).
				HandleInput(func(event ui.InputEvent) bool {
					if event.Kind == ui.InputPointerDown && event.Button == 1 {
						secondary++
					}
					return false
				}).
				OnClick(func() { primary++ })
		})
		ContextMenu(anchor, MenuOptions{
			Sections: []MenuSectionOptions{{Items: []MenuItemOptions{
				{Label: "Open"},
				{Label: "Unavailable", Disabled: true},
				{Label: "Locked", Disabled: true, Submenu: []MenuItemOptions{{Label: "Inside"}}},
			}}},
			OnSelect: func(string) { chosen++ },
		})
	}, 320, 200)

	if err := tt.RightClick("Child"); err != nil {
		t.Fatal(err)
	}
	if secondary != 1 {
		t.Fatalf("child received %d secondary presses, want 1", secondary)
	}
	if err := tt.ChooseMenuItem("Unavailable"); err == nil {
		t.Fatal("disabled menu item was selectable")
	}
	if err := tt.ChooseMenuItem("Locked", "Inside"); err == nil {
		t.Fatal("disabled submenu was selectable")
	}
	if chosen != 0 {
		t.Fatalf("disabled item invoked %d callbacks", chosen)
	}
	if err := tt.ChooseMenuItem("Open"); err != nil {
		t.Fatal(err)
	}
	if chosen != 1 {
		t.Fatalf("right-click activation invoked %d callbacks, want 1", chosen)
	}

	if err := tt.Click("Child"); err != nil {
		t.Fatal(err)
	}
	if primary != 1 || !tt.Focused("Child") {
		t.Fatalf("primary clicks %d, child focused %v", primary, tt.Focused("Child"))
	}
	tt.Key(ui.Shift, ui.KeyF10)
	if err := tt.ChooseMenuItem("Open"); err != nil {
		t.Fatalf("Shift+F10: %v", err)
	}
	if chosen != 2 || !tt.Focused("Child") {
		t.Fatalf("keyboard activations %d, child focused %v", chosen, tt.Focused("Child"))
	}
	if err := tt.RightClick("Child"); err != nil {
		t.Fatal(err)
	}
	tt.CloseMenu()
	tt.Frame()
	if chosen != 2 {
		t.Fatalf("closing the menu invoked a callback: %d activations", chosen)
	}
}

func TestContextMenuKeepsTextSelection(t *testing.T) {
	value := "hello world"
	tt := ui.NewTester(func(c *ui.Context) {
		anchor := ui.TextInputBase(c, &value).Label("Editor").Size(180, 40)
		ContextMenu(anchor, MenuOptions{
			Sections: []MenuSectionOptions{{Items: []MenuItemOptions{{Label: "Inspect"}}}},
		})
	}, 320, 200)

	if err := tt.Click("Editor"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	bounds, ok := tt.Find("Editor")
	if !ok {
		t.Fatal("editor not found")
	}
	tt.RightClickAt(bounds.X+5, bounds.Y+bounds.H/2)
	if !tt.Focused("Editor") {
		t.Fatal("opening the menu moved focus from the editor")
	}
	if err := tt.ChooseMenuItem("Inspect"); err != nil {
		t.Fatal(err)
	}
	tt.Type("replacement")
	if value != "replacement" {
		t.Fatalf("typing after the menu produced %q; selection was not preserved", value)
	}
}

func TestContextMenuAnchorCanLeaveAndReturn(t *testing.T) {
	visible := true
	var selected []string
	var calls int
	tt := ui.NewTester(func(c *ui.Context) {
		if !visible {
			return
		}
		anchor := ui.ButtonBase(c).Label("Document").Size(140, 60)
		ContextMenu(anchor, MenuOptions{
			Selection: MenuSelectMultiple, Selected: &selected,
			Sections: []MenuSectionOptions{{Items: []MenuItemOptions{{Key: "pin", Label: "Pin"}}}},
			OnSelect: func(string) { calls++ },
		})
	}, 320, 200)

	if err := tt.RightClick("Document"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("Pin"); err != nil {
		t.Fatal(err)
	}
	visible = false
	tt.Frame()
	visible = true
	tt.Frame()
	if len(selected) != 1 || selected[0] != "pin" {
		t.Fatalf("selection lost after rebuilding anchor: %v", selected)
	}
	if err := tt.RightClick("Document"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("Pin"); err != nil {
		t.Fatal(err)
	}
	if len(selected) != 0 || calls != 2 {
		t.Fatalf("selection %v, callbacks %d after anchor returned", selected, calls)
	}
}
