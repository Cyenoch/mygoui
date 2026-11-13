package overlay

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/content"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AboutLink is an action under the application's description, such as its
// website or issue tracker. Opening a URL is the caller's work.
type AboutLink struct {
	Label   string
	OnClick func()
}

// AboutSection is a titled block of text, such as credits or a licence.
type AboutSection struct {
	Title, Body string
	// Collapsible shows only the title, a disclosure closed until opened,
	// over a monospace body that scrolls when it is long: for a licence.
	Collapsible bool
}

// AboutDialogOptions describes an application. Every part is optional and
// supplied by the caller; the library reads no build or module metadata.
type AboutDialogOptions struct {
	// Icon builds the application's icon or image above its name.
	Icon                       func()
	Name, Version, Description string
	// CopyInfo is the text of the copy button beside the version, such as
	// the version, build and platform a bug report asks for; there is no
	// button without it.
	CopyInfo  string
	Links     []AboutLink
	Sections  []AboutSection
	OnDismiss func()
}

// AboutDialog presents an application's identity in a dialog while *open is
// true: icon, name, version and description stay in place while the links
// and sections scroll under them.
func AboutDialog(c *ui.Context, open *bool, o AboutDialogOptions) ui.Element {
	options := DialogOptions{
		Title:     o.Name,
		Width:     400,
		OnDismiss: o.OnDismiss,
		Header:    func() { aboutHeader(c, o) },
	}
	if len(o.Links) > 0 || len(o.Sections) > 0 {
		options.Body = func() { aboutBody(c, o) }
	}
	return Dialog(c, open, options)
}

func aboutHeader(c *ui.Context, o AboutDialogOptions) {
	t := theme.Colors(c)
	// The header keeps clear of the close button on its end side; the same
	// room on the start side keeps the identity centered in the panel.
	overlap := max(theme.Space(2)+theme.Scale(c, closeButtonSize())-theme.Space(6), 0)
	identity := ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Center).TextAlign(ui.Center).
		Gap(theme.Space(1)).PaddingStart(overlap)
	identity.Children(func() {
		if o.Icon != nil {
			ui.Box(c).Shrink(0).Margin(0, 0, theme.Space(2), 0).Children(o.Icon)
		}
		if o.Name != "" {
			style.Text(ui.Text(c, o.Name), c, theme.H4()).FontWeight(theme.WeightBold).
				Role(ui.RoleHeading).Level(2)
		}
		if o.Version != "" || o.CopyInfo != "" {
			ui.Row(c).AlignItems(ui.Center).Gap(theme.Space(1)).Children(func() {
				if o.Version != "" {
					style.Text(ui.Text(c, o.Version), c, theme.SM()).TextColor(t.MutedForeground)
				}
				if o.CopyInfo != "" {
					content.CopyButton(c, content.CopyButtonOptions{Text: o.CopyInfo, Size: theme.ButtonExtraSmall})
				}
			})
		}
		if o.Description != "" {
			style.Text(ui.Text(c, o.Description), c, theme.SM()).TextColor(t.MutedForeground).
				Margin(theme.Space(2), 0, 0, 0)
		}
	})
}

func aboutBody(c *ui.Context, o AboutDialogOptions) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space(4)).Padding(theme.Space(2), 0, 0, 0).Children(func() {
		if len(o.Links) > 0 {
			ui.Row(c).FillWidth().MinWidth(0).Wrap().Justify(ui.Center).Gap(theme.Space(4)).Children(func() {
				for i, link := range o.Links {
					mygoui.Link(c.Key(i), mygoui.LinkOptions{Label: link.Label, Icon: lucide.ArrowUpRight(), OnClick: link.OnClick})
				}
			})
		}
		for i, section := range o.Sections {
			if section.Collapsible {
				aboutCollapsible(c.Key(i), section)
				continue
			}
			ui.Column(c.Key(i)).FillWidth().MinWidth(0).Gap(theme.Space(1)).Children(func() {
				if section.Title != "" {
					style.Text(ui.Text(c, section.Title), c, theme.SM()).FontWeight(theme.WeightSemibold).
						Role(ui.RoleHeading).Level(3)
				}
				if section.Body != "" {
					style.Text(ui.Text(c, section.Body), c, theme.SM()).TextColor(t.MutedForeground)
				}
			})
		}
	})
}

// aboutBodyMaxHeight is how tall a collapsible section's text grows before it
// scrolls.
const aboutBodyMaxHeight = 160

func aboutCollapsible(c *ui.Context, section AboutSection) {
	t := theme.Colors(c)
	// The disclosure pads its trigger and panel by 8 DIPs; the wrapper gives
	// that back so both line up with the sections around it.
	ui.Column(c).MinWidth(0).Margin(0, -theme.Space(2)).Children(func() {
		navigation.Disclosure(c, navigation.DisclosureOptions{
			Title: section.Title,
			Body: func() {
				ts := theme.Code()
				ts.Size, ts.LineHeight = 12, 18
				ui.Scroll(c).FillWidth().MinWidth(0).MaxHeight(theme.Scale(c, aboutBodyMaxHeight)).
					Padding(theme.Space(3)).Radius(theme.RadiusScale().MD).Background(t.MutedBackground).
					Label(section.Title).Children(func() {
					style.Text(ui.Text(c, section.Body), c, ts).TextColor(t.MutedForeground).Selectable()
				})
			},
		})
	})
}
