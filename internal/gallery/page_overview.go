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
		Description: "MyGo UI component reference, with live examples and native additions.",
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
	{"macOS", "Apple silicon · DMG", "mygoui-gallery-macos-arm64.dmg", lucide.Apple},
	{"Windows", "x64 · EXE", "mygoui-gallery-windows-amd64.exe", lucide.Monitor},
	{"Linux", "x64 · tar.gz", "mygoui-gallery-linux-amd64.tar.gz", lucide.Terminal},
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

// buildOverview is the gallery's home: what the library is, a live sample,
// what it offers, the catalog and how to get the gallery.
func buildOverview(c *ui.Context, s *State) ui.Element {
	width := overviewWidth(c, s)
	st := StateOf[overviewState](s, "overview")
	groups := s.overviewGroups.of("", func() []pageGroup { return overviewGroups(Pages()) })
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space64)).
		Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space40)).Children(func() {
				overviewHero(c, s, st, groups, width)
				overviewShowcase(c, st, width)
			})
			overviewStats(c, groups, width)
			overviewSection(c, "Built for native apps", func() { overviewFeatures(c, width) })
			overviewSection(c, "Components", func() { overviewGroupGrid(c, s, groups, width) })
			overviewSection(c, "Get the gallery", func() { overviewDownloads(c, width) })
			overviewFooter(c)
		})
}

// overviewCompactBelow is the content width under which the home page
// stacks its rows and uses the smaller display size.
const overviewCompactBelow = 560

func componentCount(groups []pageGroup) int {
	count := 0
	for _, group := range groups {
		count += len(group.pages)
	}
	return count
}

// overviewHero is the centered opening: status pill, display headline, the
// two main actions and the install command, over a grid that fades out
// around it.
func overviewHero(c *ui.Context, s *State, st *overviewState, groups []pageGroup, width float32) {
	t := theme.Colors(c)
	compact := width < theme.Scale(c, overviewCompactBelow)
	display := theme.H1()
	display.Size, display.LineHeight = 60, 64
	top := theme.Scale(c, theme.Space64)
	switch {
	case compact:
		display.Size, display.LineHeight = 40, 44
		top = theme.Scale(c, theme.Space24)
	case width >= theme.Scale(c, 900):
		display.Size, display.LineHeight = 76, 80
		top = theme.Scale(c, 96)
	}
	cell := theme.Scale(c, 48)
	hero := ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Center).
		Gap(theme.Scale(c, theme.Space20)).Padding(top, 0, theme.Scale(c, theme.Space8), 0)
	hero.Draw(func(p *ui.Painter, r ui.Rect) {
		line := t.Border.Alpha(0.7)
		for x := r.X + r.W/2 - cell/2; x > r.X; x -= cell {
			p.Line(x, r.Y, x, r.Y+r.H, theme.Hairline, line)
			mirror := 2*(r.X+r.W/2) - x
			p.Line(mirror, r.Y, mirror, r.Y+r.H, theme.Hairline, line)
		}
		for y := r.Y + cell/2; y < r.Y+r.H; y += cell {
			p.Line(r.X, y, r.X+r.W, y, theme.Hairline, line)
		}
		clear := t.Background
		clear.A = 0
		p.FillRadialGradient(r, ui.RadialGradient{
			Radius: max(r.W, r.H) * 0.55,
			Stops: []ui.GradientStop{
				{Color: t.Background, At: 0}, {Color: clear, At: 0.45}, {Color: t.Background, At: 1},
			},
		}, 0)
	})
	hero.Children(func() {
		pill := theme.XS()
		pill.Weight = theme.WeightMedium
		ui.Row(c).AlignItems(ui.Center).Gap(theme.Space8).
			Padding(theme.Space4, theme.Space12, theme.Space4, theme.Space10).
			Radius(theme.Scale(c, 999)).Background(t.Background).
			BorderWidth(theme.Hairline).BorderColor(t.Border).Children(func() {
			dot := theme.Scale(c, 6)
			ui.Box(c).Size(dot, dot).Radius(dot).Background(t.Success)
			style.Text(ui.Text(c, fmt.Sprintf("MyGo 0.9 · %d components", componentCount(groups))), c, pill).
				TextColor(t.MutedForeground).SingleLine()
		})
		style.Text(ui.Text(c, "MyGo UI, native in Go"), c, display).
			TextColor(t.Foreground).TextAlign(ui.Center).Role(ui.RoleHeading)
		lead := theme.Base()
		lead.Size, lead.LineHeight = 18, 28
		style.Text(ui.Text(c, "Accessible desktop and web components, built on MyGo."), c, lead).
			TextColor(t.MutedForeground).TextAlign(ui.Center).MaxWidth(theme.Scale(c, 520))
		ui.Row(c).Wrap().Justify(ui.Center).Gap(theme.Space8).Padding(theme.Space8, 0, 0, 0).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Browse components", EndIcon: lucide.ArrowRight(), Size: theme.ButtonLarge,
				OnClick: func() {
					if page, ok := pageFor("button"); ok {
						selectPage(s, page)
					}
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Customize", Icon: lucide.Palette(), Variant: mygoui.ButtonOutline, Size: theme.ButtonLarge,
				OnClick: func() { s.AppearanceOpen = true },
			})
		})
		commandBar(c, &st.InstallCopiedAt, "Install command", "Copy install command", installCommand).
			MaxWidth(theme.Scale(c, 400))
	})
}

// overviewStats states the catalog in numbers taken from the registered
// pages and published builds.
func overviewStats(c *ui.Context, groups []pageGroup, width float32) {
	t := theme.Colors(c)
	stats := [...]struct{ value, label string }{
		{fmt.Sprint(componentCount(groups)), "Components"},
		{fmt.Sprint(len(groups)), "Families"},
		{fmt.Sprint(len(galleryDownloads)), "Platforms"},
		{"MIT", "License"},
	}
	cols := len(stats)
	if width < theme.Scale(c, overviewCompactBelow) {
		cols = 2
	}
	cellWidth := width / float32(cols)
	ui.Row(c).FillWidth().MinWidth(0).Wrap().
		BorderWidth(theme.Hairline, 0, theme.Hairline, 0).BorderColor(t.Border).Children(func() {
		for i, stat := range stats {
			cell := ui.Column(c.Key(stat.label)).Width(cellWidth).Shrink(0).Gap(theme.Space4).
				Padding(theme.Scale(c, theme.Space24), theme.Scale(c, theme.Space16)).
				BorderColor(t.Border)
			left, top := float32(0), float32(0)
			if i%cols > 0 {
				left = theme.Hairline
			}
			if i >= cols {
				top = theme.Hairline
			}
			cell.BorderWidth(top, 0, 0, left).Children(func() {
				mygoui.Typography(c, stat.value, mygoui.TypographyOptions{Type: theme.TypeH2})
				mygoui.Typography(c, stat.label, mygoui.TypographyOptions{Type: theme.TypeBodySM, Muted: true})
			})
		}
	})
}

var overviewFeatureList = [...]struct {
	icon        func() *ui.SVG
	title, text string
}{
	{lucide.MousePointerClick, "Native controls", "Focus, selection and text editing come from MyGo on macOS, Windows, Linux and the web."},
	{lucide.Braces, "Typed Go API", "Options are plain structs, and state stays in your code."},
	{lucide.Accessibility, "Accessible", "Built on native roles, names and keyboard navigation."},
	{lucide.Palette, "Themeable", "Light and dark palettes, accent and radius from one token set."},
}

// overviewFeatures says what the library is built on, two to a row.
func overviewFeatures(c *ui.Context, width float32) {
	t := theme.Colors(c)
	_, cardWidth := overviewColumns(c, width, 280, 2)
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Scale(c, theme.Space12)).AlignItems(ui.Stretch).Children(func() {
		for _, feature := range overviewFeatureList {
			ui.Row(c.Key(feature.title)).Width(cardWidth).Shrink(0).AlignItems(ui.Start).
				Gap(theme.Scale(c, theme.Space16)).Padding(theme.Scale(c, theme.Space20)).
				Radius(theme.RadiusScale().XL).Background(t.Card).
				BorderWidth(theme.Hairline).BorderColor(t.Border).Children(func() {
				box, icon := theme.Scale(c, 36), theme.Scale(c, theme.IconMD)
				ui.Row(c).Size(box, box).Shrink(0).AlignItems(ui.Center).Justify(ui.Center).
					Radius(theme.RadiusScale().LG).Background(t.MutedBackground).Children(func() {
					ui.Icon(c, feature.icon()).Size(icon, icon)
				})
				ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space4).Children(func() {
					mygoui.Typography(c, feature.title, mygoui.TypographyOptions{Type: theme.TypeH6})
					mygoui.Typography(c, feature.text, mygoui.TypographyOptions{Type: theme.TypeBodySM, Muted: true})
				})
			})
		}
	})
}

// overviewSection is a titled block of the home page.
func overviewSection(c *ui.Context, title string, body func()) {
	ui.Column(c.Key(title)).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space20)).Children(func() {
		mygoui.Typography(c, title, mygoui.TypographyOptions{Type: theme.TypeH3})
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

// overviewShowcase is a stage of working library components.
func overviewShowcase(c *ui.Context, st *overviewState, width float32) {
	t := theme.Colors(c)
	pad := theme.Scale(c, theme.Space12)
	width -= 2 * (pad + theme.Hairline)
	cols, cardWidth := overviewColumns(c, width, 280, 2)
	if cols == 1 {
		cardWidth = width
	}
	stage := ui.Column(c).FillWidth().MinWidth(0).Padding(pad).Gap(theme.Scale(c, theme.Space12)).
		Radius(theme.RadiusScale().XL + pad).Background(t.BackgroundSecondary).
		BorderWidth(theme.Hairline).BorderColor(t.Border).Label("Live preview")
	stage.Children(func() {
		overviewShowcaseCards(c, st, cardWidth)
		overviewShowcaseStrip(c)
	})
}

// overviewShowcaseStrip lines up the action and status variants.
func overviewShowcaseStrip(c *ui.Context) {
	t := theme.Colors(c)
	ui.Row(c).FillWidth().MinWidth(0).Wrap().AlignItems(ui.Center).Gap(theme.Space8).
		Padding(theme.Scale(c, theme.Space16), theme.Scale(c, theme.Space24)).
		Radius(theme.RadiusScale().XL).Background(t.Card).
		BorderWidth(theme.Hairline).BorderColor(t.Border).Label("Variants preview").Children(func() {
		for _, button := range [...]struct {
			label   string
			variant mygoui.ButtonVariant
		}{
			{"Default", mygoui.ButtonPrimary}, {"Secondary", mygoui.ButtonSecondary},
			{"Outline", mygoui.ButtonOutline}, {"Ghost", mygoui.ButtonGhost},
			{"Destructive", mygoui.ButtonDanger},
		} {
			mygoui.Button(c.Key(button.label), mygoui.ButtonOptions{
				Label: button.label, Variant: button.variant, Size: theme.ButtonSmall,
			})
		}
		ui.Row(c).Grow(1)
		mygoui.Badge(c, mygoui.BadgeOptions{Label: "Stable", Variant: mygoui.BadgeSuccess})
		mygoui.Badge(c, mygoui.BadgeOptions{Label: "Beta", Variant: mygoui.BadgeOutline})
	})
}

func overviewShowcaseCards(c *ui.Context, st *overviewState, cardWidth float32) {
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
				ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
					if glyph := groupIcons[group.name]; glyph != nil {
						tile, size := theme.Scale(c, 36), theme.Scale(c, theme.IconMD)
						ui.Row(c).Size(tile, tile).Shrink(0).Justify(ui.Center).AlignItems(ui.Center).
							Radius(theme.RadiusScale().LG).Background(theme.Colors(c).MutedBackground).
							Children(func() { ui.Icon(c, glyph()).Size(size, size) })
					}
					ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space4).Children(func() {
						mygoui.Typography(c, group.name, mygoui.TypographyOptions{Type: theme.TypeH6, Truncate: true})
						mygoui.Typography(c, count, mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
					})
					icon := theme.Scale(c, theme.IconSM)
					ui.Icon(c, lucide.ArrowRight()).Size(icon, icon).Shrink(0).
						TextColor(theme.Colors(c).MutedForeground)
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
