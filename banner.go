package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// BannerOptions configures a message bar. Status tints it as an Alert's
// does. Action builds its controls after the text. OnDismiss adds a close
// button and is called when it is pressed: the caller stops building the
// banner.
type BannerOptions struct {
	Text   string
	Status StatusColor
	Icon   *ui.SVG
	// ShowIcon uses the status glyph when Icon is nil.
	ShowIcon  bool
	Action    func()
	OnDismiss func()
}

// Banner is a persistent message bar for the top of a window or a pane: one
// square row from edge to edge, ruled along its bottom. An Alert is the
// rounded callout within content. Its role is an Alert's. The text takes
// the width the icon and the close button leave, and the actions wrap below
// it when there is no room beside it. When the caller stops building it
// after OnDismiss, the bar fades out over the reveal duration; what follows it
// moves up at once unless the caller gives it a Position transition.
func Banner(c *ui.Context, o BannerOptions) ui.Element {
	t := theme.Colors(c)
	border, fill, iconInk := t.Border, t.MutedBackground, t.MutedForeground
	if o.Status != ColorDefault {
		status := statusColors(&t, o.Status)
		border, fill, iconInk = status.Solid.Alpha(0.32), status.Solid.Alpha(0.08), status.Solid
	}
	// Everything is centered on the first line, which is as tall as the
	// close button, so that a message of several lines keeps its icon, its
	// actions and its close button at the top.
	line := theme.Scale(c, bannerLine)
	m := theme.Motion()
	e := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space(2)).
		Padding(theme.Space(1.5), theme.Space(3)).MinHeight(theme.Scale(c, theme.Space(10))).
		Background(fill).BorderWidth(0, 0, theme.Hairline, 0).BorderColor(border).
		TextColor(t.Foreground).Role(messageRole(o.Status)).Label(o.Text).Transition(ui.ElementTransition{
		Duration: m.State, Ease: m.Ease, Colors: true,
		Enter: &ui.Motion{}, Exit: &ui.Motion{}, ExitDuration: m.Reveal, ExitEase: ui.EaseInOut,
	})
	return e.Children(func() {
		icon := o.Icon
		if icon == nil && o.ShowIcon {
			icon = StatusIcon(o.Status)
		}
		if icon != nil {
			glyph := theme.Scale(c, theme.IconMD)
			ui.Row(c).Height(line).Shrink(0).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon).Size(glyph, glyph).TextColor(iconInk).Role(ui.RoleNone)
			})
		}
		// The text and the actions share rows that wrap: the text asks for
		// a readable measure, and the actions go under it when they do not
		// fit beside it.
		ui.Row(c).Grow(1).Basis(0).MinWidth(0).Wrap().GapX(theme.Space(2)).AlignItems(ui.Start).Children(func() {
			pad := (line - theme.SM().LineHeight*theme.TextScale(c)) / 2
			style.Text(ui.Text(c, o.Text), c, theme.SM()).Grow(1).Basis(theme.Scale(c, 160)).MinWidth(0).
				Padding(pad, 0)
			if o.Action != nil {
				ui.Row(c).Shrink(0).MinHeight(line).AlignItems(ui.Center).Gap(theme.Space(2)).Children(o.Action)
			}
		})
		if o.OnDismiss != nil {
			ui.Row(c).Height(line).Shrink(0).AlignItems(ui.Center).Children(func() {
				CloseButton(c, CloseButtonOptions{Label: locale.M().Close, OnClick: o.OnDismiss})
			})
		}
	})
}

// bannerLine is the height of the first line of a banner, in DIPs.
const bannerLine = 28
