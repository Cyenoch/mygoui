package navigation

import (
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type NavigationMenuItem struct {
	ID       overlay.MenuActionID
	Label    string
	Disabled bool
	OnSelect func()
	Items    []overlay.MenuItem
}
type NavigationMenuState struct {
	menus map[overlay.MenuActionID]*navigationDropdown
}

type navigationDropdown struct {
	menu    overlay.MenuState
	trigger ui.Handle
	wasOpen bool
}

// NavigationMenu uses the shared overlay menu engine for popup focus, arrows,
// typeahead, nested actions, Escape and return to the opener. Its top-level
// native focus group supplies Left/Right/Home/End, skipping disabled entries.
func NavigationMenu(c *ui.Context, s *NavigationMenuState, items []NavigationMenuItem) ui.Element {
	root := ui.Row(c).Role(ui.RoleGroup).Label(locales.ComponentNavigationMenu()).FocusGroup(ui.Horizontal).Gap(theme.Space4)
	root.Children(func() {
		for _, item := range items {
			if len(item.Items) == 0 {
				mygoui.Button(c.Key(item.ID), mygoui.ButtonOptions{Label: item.Label, Disabled: item.Disabled, OnClick: item.OnSelect, Variant: mygoui.ButtonGhost, Size: mygoui.ButtonSM})
				continue
			}
			if s.menus == nil {
				s.menus = make(map[overlay.MenuActionID]*navigationDropdown)
			}
			state := s.menus[item.ID]
			if state == nil {
				state = &navigationDropdown{}
				s.menus[item.ID] = state
			}
			if item.Disabled {
				state.menu.Open = false
			}
			// Dropdown content can take focus before the shared Overlay records
			// its opener. Restore the explicit trigger when this menu closes.
			if state.wasOpen && !state.menu.Open && !item.Disabled {
				state.trigger.Focus(c)
			}
			state.wasOpen = state.menu.Open
			overlay.DropdownMenu(c.Key(item.ID), &state.menu, overlay.MenuOptions{
				Label: item.Label, Disabled: item.Disabled, Handle: &state.trigger,
			}, item.Items)
		}
	})
	return root
}
