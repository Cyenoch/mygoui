package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/busy"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// LoadingOverlayOptions configures a loading cover. Visible is the caller's
// flag. Label names what is loading to assistive technology; the locale's
// Loading when empty.
type LoadingOverlayOptions struct {
	Visible bool
	Label   string
	// ShowLabel prints Label under the spinner, on a plate that keeps it
	// legible over the content.
	ShowLabel bool
}

// LoadingOverlay builds body in a container and, while Visible, covers it
// with a light scrim and a spinner. The covered content keeps its place and still
// shows, but takes neither the pointer nor the keyboard and loses the focus,
// which it gets back when the cover hides unless the user moved it
// meanwhile; the container is busy to assistive technology and the cover is
// a status named by Label. Size the returned container as any other.
func LoadingOverlay(c *ui.Context, o LoadingOverlayOptions, body func()) ui.Element {
	root := ui.Column(c).MinWidth(0).AlignItems(ui.Stretch).Busy(o.Visible)
	shown := busy.Show(root, o.Visible)
	return root.Children(func() {
		if o.Visible {
			// An overlay notes the focus as it comes and, gone while no
			// element has the focus, gives it back (MyGo ui/scope.go,
			// restoreFocus): the focus the covered content loses returns
			// when the cover hides.
			ui.Overlay(c, func() { ui.Box(c).Absolute().Size(0, 0).PassThrough() })
		}
		ui.Column(c).Grow(1).MinWidth(0).MinHeight(0).AlignItems(ui.Stretch).
			IgnoreInput(o.Visible).Children(body)
		if !o.Visible && shown <= 0 {
			return
		}
		label := o.Label
		if label == "" {
			label = locale.M().Loading
		}
		cover := ui.Box(c).Absolute().Top(0).Right(0).Bottom(0).Left(0).Center().
			Background(theme.Colors(c).Background.Alpha(0.4)).
			Role(ui.RoleStatus).Label(label).Opacity(shown)
		// While visible the cover holds input from the first frame and only its
		// look waits; fading out it lets input through.
		if !o.Visible {
			cover.PassThrough()
		}
		cover.Children(func() {
			if !o.ShowLabel {
				Spinner(c, SpinnerOptions{Size: SpinnerLarge, Label: label})
				return
			}
			t := theme.Colors(c)
			ui.Column(c).MaxWidth(theme.Scale(c, 240)).AlignItems(ui.Center).Gap(theme.Space(2)).
				Padding(theme.Space(4), theme.Space(6)).Radius(theme.RadiusScale().LG).
				Background(t.Background.Alpha(0.96)).Children(func() {
				Spinner(c, SpinnerOptions{Size: SpinnerLarge, Label: label})
				style.Text(ui.Text(c, label), c, theme.SM()).TextColor(t.MutedForeground).Role(ui.RoleNone)
			})
		})
	})
}
