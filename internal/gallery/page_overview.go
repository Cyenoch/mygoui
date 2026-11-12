package gallery

import (
	"fmt"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
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

// PagesURL is where the browser build of the gallery is published.
const PagesURL = "https://cyenoch.github.io/mygoui/"

// installCommand adds the library to a module.
const installCommand = "go get github.com/Cyenoch/mygoui"

// galleryDownload is one published build of the gallery. File names are
// stable across releases, so the latest release always resolves them.
type galleryDownload struct {
	Platform, Detail, File string
	Icon                   func() *ui.SVG
}

var galleryDownloads = []galleryDownload{
	{"macOS", "Apple silicon · DMG", "coss-ui-gallery-macos-arm64.dmg", lucide.Apple},
	{"Windows", "x64 · EXE", "coss-ui-gallery-windows-amd64.exe", lucide.Monitor},
	{"Linux", "x64 · tar.gz", "coss-ui-gallery-linux-amd64.tar.gz", lucide.Terminal},
	{"Web", "Opens in the browser", "", lucide.Globe},
}

// URL is the download address of d, or the browser build for Web.
func (d galleryDownload) URL() string {
	if d.File == "" {
		return PagesURL
	}
	return RepositoryURL + "/releases/latest/download/" + d.File
}

type overviewState struct {
	InstallCopiedAt time.Time
	Email           string
	Remember        bool
	Notify          bool
}

// buildOverview is the gallery's home: what it is, how to get it, a live
// sample and the catalog.
func buildOverview(c *ui.Context, s *State) ui.Element {
	width := overviewWidth(c, s)
	st := StateOf[overviewState](s, "overview")
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space48)).
		Children(func() {
			overviewHero(c, s, st)
			overviewSection(c, "Download", func() { overviewDownloads(c, width) })
			overviewSection(c, "Live preview", func() { overviewShowcase(c, st, width) })
			overviewSection(c, "Components", func() { overviewGroupGrid(c, s, overviewGroups(Pages()), width) })
			overviewFooter(c)
		})
}

func overviewHero(c *ui.Context, s *State, st *overviewState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space16).Padding(theme.Space16, 0, 0, 0).Children(func() {
		mygoui.Badge(c, mygoui.BadgeOptions{Label: "MyGo 0.9", Variant: mygoui.BadgeOutline})
		mygoui.Typography(c, "coss ui, native in Go", mygoui.TypographyOptions{Type: theme.TypeH1})
		mygoui.Typography(c, "Accessible desktop and web components with the coss look, built on MyGo.",
			mygoui.TypographyOptions{Type: theme.TypeBody, Muted: true}).MaxWidth(theme.Scale(c, 560))
		ui.Row(c).Wrap().Gap(theme.Space8).Padding(theme.Space8, 0, 0, 0).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Browse components", EndIcon: lucide.ArrowRight(),
				OnClick: func() {
					if page, ok := pageFor("button"); ok {
						selectPage(s, page)
					}
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Customize", Icon: lucide.Palette(), Variant: mygoui.ButtonOutline,
				OnClick: func() { s.AppearanceOpen = true },
			})
		})
		commandBar(c, &st.InstallCopiedAt, "Install command", "Copy install command", installCommand).
			MaxWidth(theme.Scale(c, 560))
	})
}

// overviewSection is a titled block of the home page.
func overviewSection(c *ui.Context, title string, body func()) {
	ui.Column(c.Key(title)).FillWidth().MinWidth(0).Gap(theme.Space16).Children(func() {
		mygoui.Typography(c, title, mygoui.TypographyOptions{Type: theme.TypeH4})
		body()
	})
}

// overviewColumns fits up to max cards of at least minWidth DIPs in width.
func overviewColumns(c *ui.Context, width, minWidth float32, most int) (int, float32) {
	gap := theme.Scale(c, theme.Space12)
	cols := most
	for cols > 1 && (width-float32(cols-1)*gap)/float32(cols) < theme.Scale(c, minWidth) {
		cols--
	}
	return cols, (width - float32(cols-1)*gap) / float32(cols)
}

func overviewDownloads(c *ui.Context, width float32) {
	_, cardWidth := overviewColumns(c, width, 160, len(galleryDownloads))
	services := c.Services()
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Scale(c, theme.Space12)).Children(func() {
		for _, d := range galleryDownloads {
			trailing := lucide.Download()
			label := "Download for " + d.Platform
			if d.File == "" {
				trailing, label = lucide.ExternalLink(), "Open in the browser"
			}
			actionCard(c.Key(d.Platform), label, cardWidth, func() { services.OpenURL(d.URL()) }, func() {
				t := theme.Colors(c)
				icon := theme.Scale(c, theme.IconMD)
				ui.Row(c).FillWidth().Gap(theme.Space12).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, d.Icon()).Size(icon, icon).Shrink(0)
					ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space2).Children(func() {
						mygoui.Typography(c, d.Platform, mygoui.TypographyOptions{Type: theme.TypeH6, Truncate: true})
						mygoui.Typography(c, d.Detail, mygoui.TypographyOptions{
							Type: theme.TypeBodyXS, Muted: true, Truncate: true,
						})
					})
					ui.Icon(c, trailing).Size(icon, icon).Shrink(0).TextColor(t.MutedForeground)
				})
			})
		}
	})
}

// overviewShowcase is a small working form built from library components.
func overviewShowcase(c *ui.Context, st *overviewState, width float32) {
	cols, cardWidth := overviewColumns(c, width, 280, 2)
	if cols == 1 {
		cardWidth = width
	}
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Scale(c, theme.Space12)).AlignItems(ui.Stretch).Children(func() {
		ui.Column(c.Key("account")).Width(cardWidth).Shrink(0).Children(func() {
			mygoui.Card(c, mygoui.CardOptions{Label: "Account preview"}, func() {
				mygoui.CardHeader(c, func() {
					mygoui.CardTitle(c, "Create an account")
					mygoui.CardDescription(c, "Native fields, focus and validation.")
				})
				mygoui.CardContent(c, func() {
					ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
						form.Input(c, &st.Email, form.InputOptions{
							Label: "Email", Placeholder: "you@example.com", FullWidth: true,
						})
						form.Checkbox(c, form.CheckboxOptions{Label: "Remember me", Checked: &st.Remember})
					})
				}).Grow(1)
				mygoui.CardFooter(c, func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Create account"})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Cancel", Variant: mygoui.ButtonGhost})
				})
			}).Grow(1)
		})
		ui.Column(c.Key("status")).Width(cardWidth).Shrink(0).Children(func() {
			mygoui.Card(c, mygoui.CardOptions{Label: "Status preview"}, func() {
				mygoui.CardHeader(c, func() {
					mygoui.CardTitle(c, "Deployment")
					mygoui.CardDescription(c, "Status, progress and settings.")
				})
				mygoui.CardContent(c, func() {
					ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
						ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
							mygoui.Badge(c, mygoui.BadgeOptions{Label: "Live", Variant: mygoui.BadgeSuccess})
							mygoui.Badge(c, mygoui.BadgeOptions{Label: "v0.1.0", Variant: mygoui.BadgeOutline})
						})
						mygoui.ProgressBar(c, mygoui.ProgressOptions{
							Label: "Upload", Value: 64, Max: 100, ShowOutput: true,
						})
						form.Switch(c, form.SwitchOptions{On: &st.Notify, Label: "Email me on failure"})
					})
				}).Grow(1)
				mygoui.CardFooter(c, func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "View logs", Variant: mygoui.ButtonOutline})
				})
			}).Grow(1)
		})
	})
}

func overviewFooter(c *ui.Context) {
	t := theme.Colors(c)
	services := c.Services()
	ui.Row(c).FillWidth().Wrap().Gap(theme.Space8).AlignItems(ui.Center).
		Padding(theme.Space24, 0, 0, 0).BorderWidth(theme.Hairline, 0, 0, 0).BorderColor(t.Border).Children(func() {
		mygoui.Typography(c, "MIT licensed. Design from coss ui.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		}).Grow(1)
		mygoui.Link(c, mygoui.LinkOptions{
			Label: "Source", Icon: lucide.ExternalLink(),
			OnClick: func() { services.OpenURL(RepositoryURL) },
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
	_, cardWidth := overviewColumns(c, width, 200, 3)
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Scale(c, theme.Space12)).AlignItems(ui.Start).Children(func() {
		for _, group := range groups {
			count := fmt.Sprintf("%d components", len(group.pages))
			if len(group.pages) == 1 {
				count = "1 component"
			}
			actionCard(c.Key(group.name), "Browse "+group.name, cardWidth, func() { selectPage(s, group.pages[0]) }, func() {
				ui.Column(c).FillWidth().Gap(theme.Space4).Children(func() {
					mygoui.Typography(c, group.name, mygoui.TypographyOptions{Type: theme.TypeH6, Truncate: true})
					mygoui.Typography(c, count, mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
				})
			})
		}
	})
}

// actionCard is a bordered card that is one native button.
func actionCard(c *ui.Context, label string, width float32, onClick, body func()) {
	t := theme.Colors(c)
	radius := theme.RadiusScale().XL
	card := ui.ButtonBase(c).Width(width).Shrink(0).Padding(theme.Scale(c, theme.Space16)).
		Radius(radius).Justify(ui.Start).AlignItems(ui.Center).
		BorderWidth(theme.Hairline).BorderColor(t.Border).FocusRing(true).Label(label)
	card = style.Background(card, radius, t.Card, t.Accent, t.Accent)
	card.Children(body)
	card.OnClick(onClick)
}
