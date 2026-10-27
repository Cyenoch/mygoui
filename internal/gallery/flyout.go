package gallery

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func nativeFlyoutExample(c *ui.Context, s *State) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: "Open native flyout", Disabled: s.Window == nil,
	}).Bind(&s.flyoutTrigger).OnClick(func() {
		if s.flyout != nil {
			s.flyout.Close()
			return
		}
		var name string
		window, err := overlay.NewFlyout(overlay.FlyoutOptions{
			Parent: s.Window, Anchor: overlay.FlyoutAnchor(c, s.flyoutTrigger),
			Placement: mygo.PlacementBottomStart, Gap: 4,
			Width: 320, Height: 240, Focusable: true, Dark: &s.Dark,
			Surface: overlay.FlyoutSurfaceOptions{
				Title: "Workspace", Description: "Project settings",
				Body: func(c *ui.Context) {
					form.Input(c, &name, form.InputOptions{Label: "Name", Placeholder: "Acme", FullWidth: true})
				},
				Footer: func(c *ui.Context) {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Done", OnClick: func() { s.flyout.Close() }})
				},
			},
		})
		if err != nil {
			s.flyoutError = err.Error()
			return
		}
		s.flyout, s.flyoutError = window, ""
		window.OnClosed(func() {
			s.flyout = nil
			s.Window.Invalidate()
		})
	})
	if s.flyout != nil {
		anchor := overlay.FlyoutAnchor(c, s.flyoutTrigger)
		if anchor != s.flyoutAnchor {
			s.flyoutAnchor = anchor
			s.flyout.SetAnchor(anchor)
		}
	}
	if s.flyoutError != "" {
		ui.Text(c, s.flyoutError)
	}
}
