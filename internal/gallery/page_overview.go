package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func init() {
	RegisterPage(Page{
		ID:          "overview",
		Title:       "Overview",
		Description: "Coss component reference for MyGo, with live examples and native additions.",
		Group:       "Overview",
		Order:       1,
		Build:       buildOverview,
	})
}

// buildOverview puts the live catalog within reach of the first frame.
func buildOverview(c *ui.Context, s *State) ui.Element {
	width := overviewWidth(c, s)
	groups := overviewGroups(Pages())
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space24)).
		Children(func() {
			ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
				mygoui.Typography(c, "Native components", mygoui.TypographyOptions{Type: theme.TypeH1})
				ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Browse components",
						OnClick: func() {
							if page, ok := pageFor("button"); ok {
								selectPage(s, page)
							}
						},
					})
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Customize", Variant: mygoui.ButtonOutline,
						OnClick: func() { s.AppearanceOpen = true },
					})
				})
			})
			mygoui.Typography(c, "Browse", mygoui.TypographyOptions{Type: theme.TypeH5})
			overviewGroupGrid(c, s, groups, width)
		})
}

// overviewThemeControls is also the interactive theme specimen.
func overviewThemeControls(c *ui.Context, s *State) {
	ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Light", Variant: buttonFor(!s.Dark), OnClick: func() { s.setDark(false) }})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Dark", Variant: buttonFor(s.Dark), OnClick: func() { s.setDark(true) }})
			for _, radius := range []float32{4, 10, 16} {
				mygoui.Button(c.Key(radius), mygoui.ButtonOptions{
					Label: fmt.Sprintf("Radius %d", int(radius)), Variant: buttonFor(s.Radius == radius),
					OnClick: func() { s.Radius = radius },
				})
			}
		})
		swatches(c, s)
		mygoui.Frame(c, mygoui.FrameOptions{Label: "Theme preview"}, func() {
			mygoui.FrameHeader(c, func() {
				mygoui.FrameTitle(c, "Make it yours")
				mygoui.FrameDescription(c, "Primary, surface, radius")
			})
			mygoui.FramePanel(c, func() {
				ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Browse components", OnClick: func() { s.Selected = "button"; s.Content.Y = 0 }})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Customize", Variant: mygoui.ButtonOutline, OnClick: func() { s.Selected = "theming"; s.Content.Y = 0 }})
				})
			})
		})
	})
}

// overviewGroups returns the navigation groups except Overview, in the same
// order as the sidebar.
func overviewGroups(all []Page) []pageGroup {
	var groups []pageGroup
	for _, group := range pageGroups(filterPagesInOrder(all)) {
		if group.name == "Overview" {
			continue
		}
		groups = append(groups, group)
	}
	return groups
}

// overviewGroupGrid gives each component family one navigation target.
func overviewGroupGrid(c *ui.Context, s *State, groups []pageGroup, width float32) {
	gap := theme.Scale(c, theme.Space12)
	cols := 1
	if width >= theme.Scale(c, 680) {
		cols = 3
	} else if width >= theme.Scale(c, 400) {
		cols = 2
	}
	cardWidth := (width - float32(cols-1)*gap) / float32(cols)
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(gap).AlignItems(ui.Start).Children(func() {
		for _, group := range groups {
			overviewGroupCard(c, s, group, cardWidth)
		}
	})
}

func overviewGroupCard(c *ui.Context, s *State, group pageGroup, width float32) {
	t := theme.Colors(c)
	count := fmt.Sprintf("%d components", len(group.pages))
	if len(group.pages) == 1 {
		count = "1 component"
	}
	card := ui.ButtonBase(c.Key(group.name)).Width(width).MinHeight(theme.Scale(c, 88)).
		Shrink(0).Padding(theme.Scale(c, theme.Space16)).Radius(theme.RadiusScale().LG).
		Background(t.Card).BorderWidth(theme.Hairline).BorderColor(t.Border).
		FocusRing(true).Label("Browse " + group.name).Justify(ui.Start)
	card.Children(func() {
		ui.Column(c).FillWidth().Gap(theme.Space8).Children(func() {
			mygoui.Typography(c, group.name, mygoui.TypographyOptions{Type: theme.TypeH5})
			mygoui.Typography(c, count, mygoui.TypographyOptions{
				Type: theme.TypeBodyXS, Muted: true,
			})
		})
	})
	card.OnClick(func() { selectPage(s, group.pages[0]) })
}
