package gallery

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// nativeFlyoutExample opens a popover as a native flyout, which may extend
// beyond the window. The browser build keeps it in the page.
func nativeFlyoutExample(c *ui.Context, s *State) {
	st := overlaysStateOf(s)
	anchor := mygoui.Button(c, mygoui.ButtonOptions{
		Label: "Open native flyout", Variant: mygoui.ButtonSecondary,
		OnClick: func() { st.flyoutOpen = !st.flyoutOpen },
	})
	overlay.Popover(c, anchor, &st.flyoutOpen, overlay.PopoverOptions{
		Title: "Workspace", Description: "Project settings",
		Width:     theme.Scale(c, 320),
		Placement: mygoui.Placement{Side: mygoui.SideBottom, Align: mygoui.AlignStart},
		Flyout:    &s.Flyouts,
		Body: func(c *ui.Context) {
			form.Input(c, &st.flyoutName, form.InputOptions{Label: "Name", Placeholder: "Acme", FullWidth: true})
		},
		Footer: func(c *ui.Context) {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Done", Size: theme.ButtonSmall, OnClick: func() { st.flyoutOpen = false },
			})
		},
	})
}
