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
		Description: "HeroUI's components, unified into one native Go library: live demos, copyable usage, API reference and theme tokens per page.",
		Group:       "Overview",
		Order:       1,
		Build:       buildOverview,
	})
}

// buildOverview is the landing page: a hero, the three things the library
// promises, and one card per component group with its count and page links.
func buildOverview(c *ui.Context, s *State) ui.Element {
	width := overviewWidth(c)
	groups := overviewGroups()
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space32)).
		Children(func() {
			overviewHero(c, len(groups))
			mygoui.Separator(c, mygoui.SeparatorOptions{})
			ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space16)).
				Children(func() {
					ui.Row(c).FillWidth().MinWidth(0).Wrap().AlignItems(ui.End).
						Gap(theme.Scale(c, theme.Space12)).Children(func() {
						ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
							mygoui.Typography(c, "Component groups", mygoui.TypographyOptions{Type: theme.TypeH3})
						})
						mygoui.Typography(c, fmt.Sprintf("%d pages in %d groups", len(Pages()), len(groups)),
							mygoui.TypographyOptions{Type: theme.TypeBodySM, Muted: true})
					})
					overviewGroupGrid(c, s, groups, width)
				})
		})
}

// overviewHero is the title block: library name, one-line promise, counts and
// the three feature points.
func overviewHero(c *ui.Context, groups int) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space12)).Children(func() {
		mygoui.Typography(c, "Overview", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Weight: theme.WeightMedium, Muted: true,
		})
		mygoui.Typography(c, "MyGo UI", mygoui.TypographyOptions{Type: theme.TypeH1})
		mygoui.Typography(c, "One native Go library for HeroUI's components. Every page runs the real control, not a screenshot.", mygoui.TypographyOptions{
			Type: theme.TypeBody, Muted: true,
		})
		ui.Row(c).Wrap().Gap(theme.Scale(c, theme.Space16)).AlignItems(ui.Center).Children(func() {
			for _, text := range []string{
				fmt.Sprintf("%d components", len(Pages())-1),
				fmt.Sprintf("%d groups", groups),
				"Light and dark",
				"Native keyboard and clipboard",
			} {
				mygoui.Typography(c, text, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
			}
		})
		overviewFeatures(c)
	})
}

// overviewFeatures states what the library does, in three short points.
func overviewFeatures(c *ui.Context) {
	features := []struct{ title, body string }{
		{"Native controls", "Real focus, keyboard, pointer and clipboard behavior from MyGo, not a web view."},
		{"One theme", "HeroUI's light and dark palette, accent and radius tokens drive every component."},
		{"Documented", "A live example, copyable usage, API reference and theme tokens on every page."},
	}
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Scale(c, theme.Space16)).
		AlignItems(ui.Start).Children(func() {
		for _, feature := range features {
			ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 190)).
				Gap(theme.Scale(c, theme.Space4)).Children(func() {
				mygoui.Typography(c, feature.title, mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Weight: theme.WeightSemibold,
				})
				mygoui.Typography(c, feature.body, mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Muted: true,
				})
			})
		}
	})
}

// overviewGroups returns the navigation groups except Overview, in the same
// order as the sidebar.
func overviewGroups() []pageGroup {
	var groups []pageGroup
	for _, group := range pageGroups(filterPagesInOrder(Pages())) {
		if group.name == "Overview" {
			continue
		}
		groups = append(groups, group)
	}
	return groups
}

// overviewGroupGrid lays the group cards in two columns when the landing page
// is wide enough, and one otherwise.
func overviewGroupGrid(c *ui.Context, s *State, groups []pageGroup, width float32) {
	gap := theme.Scale(c, theme.Space12)
	cols := 1
	cardWidth := width
	if width >= theme.Scale(c, 700) {
		cols = 2
		cardWidth = (width - gap) / 2
	}
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(gap).AlignItems(ui.Start).Children(func() {
		for _, group := range groups {
			overviewGroupCard(c, s, group, cardWidth, cols == 2)
		}
	})
}

// overviewGroupCard is one functional group: its name, page count and a link
// for each page it contains.
func overviewGroupCard(c *ui.Context, s *State, group pageGroup, width float32, grid bool) {
	t := theme.Colors(c)
	card := ui.Column(c).MinWidth(0).Gap(theme.Scale(c, theme.Space12)).
		Padding(theme.Scale(c, theme.Space16)).Radius(theme.Scale(c, theme.Space12)).
		Clip().Background(t.Surface).TextColor(t.SurfaceForeground).
		BorderWidth(theme.Hairline).BorderColor(t.Border).
		Label("group card " + group.name)
	if grid {
		card = card.Width(width).Shrink(0)
	} else {
		card = card.FillWidth()
	}
	card.Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).
			Gap(theme.Scale(c, theme.Space8)).Children(func() {
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
				mygoui.Typography(c, group.name, mygoui.TypographyOptions{Type: theme.TypeH5})
			})
			mygoui.Typography(c, fmt.Sprintf("%d", len(group.pages)), mygoui.TypographyOptions{
				Type: theme.TypeBodyXS, Muted: true,
			})
		})
		ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Scale(c, theme.Space4)).
			Children(func() {
				for _, page := range group.pages {
					page := page
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: page.Title, Variant: mygoui.ButtonGhost,
						Size: theme.ButtonSmall,
						OnClick: func() {
							s.Selected = page.ID
							s.Content.Y = 0
						},
					})
				}
			})
	})
}
