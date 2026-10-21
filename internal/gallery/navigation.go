package gallery

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type NavigationState struct {
	Accordion                                       navigation.AccordionState[string]
	Collapsible                                     navigation.CollapsibleState
	Breadcrumb                                      navigation.BreadcrumbState
	Select                                          navigation.SelectState
	Combobox                                        navigation.ComboboxState
	Command                                         navigation.CommandState
	Menu                                            navigation.NavigationMenuState
	Sidebar                                         navigation.SidebarState
	Choice                                          string
	Choices                                         []string
	Tab                                             string
	Page                                            int
	Disabled, Multiple, Borderless, GroupOpen, Line bool
	PlaygroundOpen                                  bool
	Team                                            string
	TeamMenu, AccountMenu                           overlay.MenuState
}

func NewNavigationState() NavigationState {
	return NavigationState{Choice: "go", Choices: []string{"go"}, Tab: "account", Page: 1, GroupOpen: true, PlaygroundOpen: true, Team: "Acme Studio"}
}

var sampleChoices = []navigation.Choice[string]{{Value: "go", Label: "Go", Group: "Languages"}, {Value: "typescript", Label: "TypeScript", Group: "Languages"}, {Value: "rust", Label: "Rust", Group: "Languages"}, {Value: "unavailable", Label: "Unavailable", Disabled: true}}

func navigationSample(c *ui.Context, s *NavigationState, id ComponentID, advanced bool) {
	action := func(label string) { c.AddToast(ui.Toast{Title: label}) }
	ui.Column(c).FillWidth().MinWidth(0).Gap(12).Children(func() {
		if advanced {
			ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
				mygoui.Toggle(c, &s.Disabled, mygoui.ToggleOptions{Label: "Disabled"})
				if id == "accordion" || id == "combobox" {
					mygoui.Toggle(c, &s.Multiple, mygoui.ToggleOptions{Label: "Multiple"})
				}
				if id == "accordion" {
					mygoui.Toggle(c, &s.Borderless, mygoui.ToggleOptions{Label: "Borderless"})
				}
				if id == "tabs" {
					mygoui.Toggle(c, &s.Line, mygoui.ToggleOptions{Label: "Line tabs"})
				}
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset navigation sample", Variant: mygoui.ButtonOutline, OnClick: func() { *s = NewNavigationState() }})
			})
		}
		switch id {
		case "accordion":
			mode := navigation.AccordionSingle
			if s.Multiple {
				mode = navigation.AccordionMultiple
			}
			navigation.Accordion(c, &s.Accordion, navigation.AccordionOptions{Mode: mode, Disabled: s.Disabled, Borderless: s.Borderless}, []navigation.AccordionItem[string]{
				{Value: "accessible", Label: "Is it accessible?", Content: func() { ui.Text(c, "Native disclosure triggers support keyboard focus and expansion.") }},
				{Value: "native", Label: "Is it native?", Content: func() { ui.Text(c, "These are Go views drawn by MyGo without a browser.") }},
				{Value: "disabled", Label: "Unavailable section", Disabled: true, Content: func() { ui.Text(c, "Unavailable") }},
			})
		case "collapsible":
			navigation.Collapsible(c, &s.Collapsible, navigation.CollapsibleOptions{Label: "Show project details", Disabled: s.Disabled}, func() { ui.Text(c, "Local sample project · Go · three documents") })
		case "breadcrumb":
			navigation.Breadcrumb(c, &s.Breadcrumb, navigation.BreadcrumbOptions{Label: "Sample path", MaxVisible: 3}, []navigation.BreadcrumbItem{
				{ID: "home", Label: "Home", OnSelect: func() { action("Home selected") }}, {ID: "projects", Label: "Projects", OnSelect: func() { action("Projects selected") }}, {ID: "native", Label: "Native UI", OnSelect: func() { action("Native UI selected") }}, {ID: "components", Label: "Components"},
			})
		case "select":
			navigation.Select(c, &s.Select, &s.Choice, navigation.SelectOptions[string]{Label: "Language", Placeholder: "Choose a language", Disabled: s.Disabled, Width: 260}, sampleChoices)
		case "combobox":
			navigation.Combobox(c, &s.Combobox, &s.Choices, navigation.ComboboxOptions[string]{Label: "Language", Placeholder: "Search languages…", Empty: "No matching language", Disabled: s.Disabled, Multiple: s.Multiple, Width: 260}, sampleChoices)
		case "command":
			navigation.Command(c, &s.Command, navigation.CommandOptions{Label: "Sample commands", Placeholder: "Search local commands…", Height: 180}, []navigation.CommandItem{
				{ID: "new", Label: "New project", Group: "Projects", Keywords: "create", OnSelect: func() { action("New project selected") }}, {ID: "open", Label: "Open project", Group: "Projects", OnSelect: func() { action("Open project selected") }}, {ID: "settings", Label: "Settings", Group: "Preferences", Disabled: s.Disabled, OnSelect: func() { action("Settings selected") }},
			})
		case "navigation-menu":
			navigation.NavigationMenu(c, &s.Menu, []navigation.NavigationMenuItem{{ID: "docs", Label: "Documentation", Items: []overlay.MenuItem{{ID: "guide", Label: "Getting started", OnSelect: func() { action("Getting started selected") }}, {ID: "components", Label: "Components", OnSelect: func() { action("Components selected") }}}}, {ID: "about", Label: "About", Disabled: s.Disabled, OnSelect: func() { action("About selected") }}}).Wrap().FillWidth().MinWidth(0)
		case "pagination":
			navigation.Pagination(c, navigation.PaginationOptions{Page: s.Page, Pages: 12, Disabled: s.Disabled, OnChange: func(page int) { s.Page = page }})
		case "sidebar":
			navigation.Sidebar(c, &s.Sidebar, navigation.SidebarOptions{Label: "Sample workspace", Width: 248, NarrowAt: 780,
				Header: func(collapsed bool) {
					ui.Row(c).FillWidth().MinWidth(0).Gap(8).AlignItems(ui.Center).Children(func() {
						if !collapsed {
							ui.Icon(c, lucide.GalleryVerticalEnd()).Size(24, 24).Shrink(0)
						}
						label := s.Team
						if collapsed {
							label = "WS"
						}
						overlay.DropdownMenu(c, &s.TeamMenu, overlay.MenuOptions{Label: label}, []overlay.MenuItem{
							{ID: "studio", Label: "Acme Studio", OnSelect: func() { s.Team = "Acme Studio" }},
							{ID: "personal", Label: "Personal workspace", OnSelect: func() { s.Team = "Personal workspace" }},
						}).MinWidth(0).Grow(1).Column().AlignItems(ui.Start).Justify(ui.Center).Height(theme.Scale(c, 48)).Padding(4, 6).Gap(2).Border(0, theme.Colors(c).SidebarBorder).Radius(theme.Radii().MD).Background(theme.Colors(c).Sidebar).Children(func() {
							if !collapsed {
								ui.Text(c, "Local workspace").FontSize(theme.Scale(c, 12)).TextColor(theme.Colors(c).MutedForeground).SingleLine()
								ui.Icon(c, lucide.ChevronsUpDown()).Attach(ui.AnchorEnd, ui.AnchorEnd).Size(16, 16)
							}
						})
					})
				},
				Footer: func(collapsed bool) {
					ui.Row(c).FillWidth().MinWidth(0).Gap(8).AlignItems(ui.Center).Children(func() {
						if !collapsed {
							mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Alex Morgan", Size: mygoui.AvatarSM})
						}
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							label := "Alex Morgan"
							if collapsed {
								label = "AM"
							}
							overlay.DropdownMenu(c, &s.AccountMenu, overlay.MenuOptions{Label: label}, []overlay.MenuItem{
								{ID: "profile", Label: "Account settings", OnSelect: func() { action("Account settings opened") }},
								{ID: "billing", Label: "Billing", OnSelect: func() { action("Local billing sample opened") }},
								{ID: "notifications", Label: "Notification preferences", OnSelect: func() { action("Notification preferences opened") }},
							}).FillWidth().Column().AlignItems(ui.Start).Justify(ui.Center).Height(theme.Scale(c, 48)).Padding(4, 6).Gap(2).Border(0, theme.Colors(c).SidebarBorder).Radius(theme.Radii().MD).Background(theme.Colors(c).Sidebar).Children(func() {
								if !collapsed {
									ui.Text(c, "alex@example.com").FontSize(theme.Scale(c, 12)).TextColor(theme.Colors(c).MutedForeground).SingleLine()
									ui.Icon(c, lucide.ChevronsUpDown()).Attach(ui.AnchorEnd, ui.AnchorEnd).Size(16, 16)
								}
							})
						})
					})
				},
				ContentHeader: func() { ui.Text(c, "Building your application").Bold().SingleLine() },
				Content: func() {
					ui.Column(c).FillWidth().Gap(16).Padding(20).Children(func() {
						ui.Text(c, "Overview").Bold().FontSize(theme.Scale(c, 20))
						ui.Text(c, "Keep your tools, documents and project settings in one workspace.").LineHeight(1.5)
						mygoui.Card(c, mygoui.CardOptions{Label: "Project activity"}, func() {
							mygoui.CardHeader(c, mygoui.CardHeaderOptions{Title: "Project activity", Description: "Local sample workspace"})
							mygoui.CardContent(c, func() { ui.Text(c, "Three documents are ready for review.") })
						})
					})
				},
			}, []navigation.SidebarGroup{{ID: "platform", Label: "Platform", Items: []navigation.SidebarItem{
				{ID: "playground", Label: "Playground", Icon: lucide.SquareTerminal(), Open: &s.PlaygroundOpen, Children: []navigation.SidebarItem{{ID: "history", Label: "History"}, {ID: "starred", Label: "Starred"}, {ID: "playground-settings", Label: "Settings"}}},
				{ID: "models", Label: "Models", Icon: lucide.Bot()}, {ID: "documentation", Label: "Documentation", Icon: lucide.BookOpen()}, {ID: "settings", Label: "Settings", Icon: lucide.Settings2()},
			}}}).FillWidth().Height(480)
		case "tabs":
			variant := navigation.TabsDefault
			if s.Line {
				variant = navigation.TabsLine
			}
			navigation.Tabs(c, &s.Tab, navigation.TabsOptions[string]{Label: "Sample settings", Variant: variant}, []navigation.Tab[string]{{Value: "account", Label: "Account", Content: func() { ui.Text(c, "Change your sample account preferences.") }}, {Value: "security", Label: "Security", Content: func() { ui.Text(c, "Review the local sample security settings.") }}, {Value: "disabled", Label: "Unavailable", Disabled: true}})
		}
	})
}
