package navigation

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SplitButtonOptions configures a split button. Label names the action
// segment; MenuLabel names the arrow segment and defaults to the library's
// "More". Open is caller-owned; nil keeps it internally. Disabled stops the
// whole control, ActionDisabled the action alone, leaving the menu.
type SplitButtonOptions struct {
	Label          string
	Icon           *ui.SVG
	Variant        mygoui.ButtonVariant
	Size           theme.ButtonSize
	Disabled       bool
	ActionDisabled bool
	// Loading marks the action in flight: it shows a spinner and takes no
	// press, while the menu stays usable.
	Loading   bool
	OnClick   func()
	MenuLabel string
	Open      *bool
	Menu      MenuOptions
	// Selected is caller-owned: the key of the alternative chosen last. The
	// menu marks it, and choosing an item writes its key; the caller makes
	// that alternative the action. It replaces Menu's own selection.
	Selected  *string
	Placement mygoui.Placement
}

// splitButtonState is the arrow's handle, which anchors the menu, and the
// menu's open flag when the caller keeps none.
type splitButtonState struct {
	arrow ui.Handle
	open  bool
	// chosen mirrors Selected for the menu, which takes a slice.
	chosen []string
}

// SplitButton builds an action joined to an arrow that opens a menu of
// alternatives, as Save beside Save As. The two segments are buttons of one
// joined group in the same variant and size; each takes the focus on its
// own. Down on the arrow opens the menu at its first item, Up at its last.
// Choosing an item, Escape or a press outside closes the menu; choosing or
// Escape returns the focus to the arrow.
func SplitButton(c *ui.Context, o SplitButtonOptions) ui.Element {
	root := ui.Row(c).Shrink(0).MinWidth(0).AlignSelf(ui.Start).Role(ui.RoleGroup).Label(o.Label)
	st := ui.Local(root, "split-button", func() splitButtonState { return splitButtonState{} })
	open := o.Open
	if open == nil {
		open = &st.open
	}
	menuLabel := o.MenuLabel
	if menuLabel == "" {
		menuLabel = locale.M().More
	}
	placement := o.Placement
	if placement == (mygoui.Placement{}) {
		placement.Align = mygoui.AlignEnd
	}
	menu := o.Menu
	if o.Selected != nil {
		st.chosen = append(st.chosen[:0], *o.Selected)
		menu.Selection, menu.Selected = MenuSelectSingle, &st.chosen
		onSelect := o.Menu.OnSelect
		menu.OnSelect = func(key string) {
			*o.Selected = key
			if onSelect != nil {
				onSelect(key)
			}
		}
	}
	root.Children(func() {
		mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Separators: true},
			mygoui.ButtonOptions{
				Label: o.Label, Icon: o.Icon, Variant: o.Variant, Size: o.Size,
				Loading: o.Loading, Disabled: o.Disabled || o.ActionDisabled, OnClick: o.OnClick,
			},
			mygoui.ButtonOptions{
				Label: menuLabel, Icon: lucide.ChevronDown(), IconOnly: true,
				Variant: o.Variant, Size: o.Size, Disabled: o.Disabled, Handle: &st.arrow,
				OnClick: func() { *open = !*open },
			})
		arrow := c.Resolve(st.arrow)
		arrow.Role(ui.RoleMenuButton).Expanded(*open)
		menuPopover(c, arrow, open, menu, placement, true)
	})
	return root
}
