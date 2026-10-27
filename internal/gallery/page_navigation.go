package gallery

import (
	"fmt"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// navigationState is caller-owned page state. Every expansion, open menu and
// selection lives here so no navigation specimen keeps hidden internal state.
type navigationState struct {
	single   []string
	multiple []string

	disclosureOpen bool
	groupSingle    [3]bool
	groupMultiple  [2]bool

	tabsPrimary   int
	tabsSecondary int
	tabsVertical  int
	tabsOverflow  int

	crumbNote   string
	page        int
	pageDefault int
	pageMedium  int
	pageSmall   int
	pageSummary int

	toolbarAction string

	menuSelected   []string
	menuInit       bool
	sortSelected   []string
	sortInit       bool
	sortChoice     string
	sortOpen       bool
	sortOpenAlt    bool
	submenuOpen    bool
	submenuOpenAlt bool
	submenuChoice  string
}

func toolbarIcon(c *ui.Context, st *navigationState, label string, icon *ui.SVG) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label:    label,
		Icon:     icon,
		IconOnly: true,
		Size:     theme.ButtonSmall,
		Variant:  mygoui.ButtonGhost,
		OnClick:  func() { st.toolbarAction = label },
	})
}

func selectionSummary(prefix string, keys []string) string {
	if len(keys) == 0 {
		return prefix + ": none"
	}
	return prefix + ": " + strings.Join(keys, ", ")
}

// accordionSingleExample keeps one section open at a time.
func accordionSingleExample(c *ui.Context, st *navigationState) {
	navigation.Accordion(c, navigation.AccordionOptions{
		Label:    "Project settings",
		Expanded: &st.single,
		Items: []navigation.AccordionItemOptions{
			{Key: "general", Title: "General", Body: func() {
				mygoui.Paragraph(c, "The project name, owner and default branch.", true)
			}},
			{Key: "members", Title: "Members", Body: func() {
				mygoui.Paragraph(c, "Invite teammates and set their access level.", true)
			}},
			{Key: "billing", Title: "Billing", Body: func() {
				mygoui.Paragraph(c, "Invoices, seats and the current plan.", true)
			}},
		},
	})
}

// accordionMultipleExample keeps several surface sections open at once.
func accordionMultipleExample(c *ui.Context, st *navigationState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		navigation.Accordion(c, navigation.AccordionOptions{
			Label:    "Workspace",
			Variant:  navigation.AccordionSurface,
			Multiple: true,
			Expanded: &st.multiple,
			Items: []navigation.AccordionItemOptions{
				{Key: "notifications", Title: "Notifications", Actions: func() {
					mygoui.Chip(c, mygoui.ChipOptions{
						Label: "New", Color: mygoui.ColorAccent,
						Variant: mygoui.ChipSoft, Size: mygoui.ChipSmall,
					})
				}, Body: func() {
					mygoui.Paragraph(c, "Email, push and weekly digest preferences.", true)
				}},
				{Key: "security", Title: "Security", Body: func() {
					mygoui.Paragraph(c, "Two-factor authentication and active sessions.", true)
				}},
			},
		})
		mygoui.Typography(c, fmt.Sprintf("%d sections expanded", len(st.multiple)), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// accordionSeparatorExample removes the hairlines between sections.
func accordionSeparatorExample(c *ui.Context, st *navigationState) {
	navigation.Accordion(c, navigation.AccordionOptions{
		Label:         "Workspace",
		Variant:       navigation.AccordionSurface,
		HideSeparator: true,
		Expanded:      &st.multiple,
		Items: []navigation.AccordionItemOptions{
			{Key: "notifications", Title: "Notifications", Body: func() {
				mygoui.Paragraph(c, "A surface accordion without dividers.", true)
			}},
			{Key: "security", Title: "Security", Body: func() {
				mygoui.Paragraph(c, "Sections still share the multiple selection.", true)
			}},
		},
	})
}

// accordionDisabledExample keeps a disabled section unavailable.
func accordionDisabledExample(c *ui.Context, st *navigationState) {
	navigation.Accordion(c, navigation.AccordionOptions{
		Label:    "Project settings",
		Expanded: &st.single,
		Items: []navigation.AccordionItemOptions{
			{Key: "general", Title: "General", Body: func() {
				mygoui.Paragraph(c, "Only the enabled sections can be toggled.", true)
			}},
			{Key: "audit", Title: "Audit log", Disabled: true, Body: func() {
				mygoui.Paragraph(c, "Unavailable on the current plan.", true)
			}},
		},
	})
}

// disclosureActionsExample shows a caller-owned panel with a header action.
func disclosureActionsExample(c *ui.Context, st *navigationState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		navigation.Disclosure(c, navigation.DisclosureOptions{
			Key:   "release-notes",
			Title: "Release notes",
			Open:  &st.disclosureOpen,
			Actions: func() {
				mygoui.Chip(c, mygoui.ChipOptions{
					Label: "v3.2", Color: mygoui.ColorDefault,
					Variant: mygoui.ChipSoft, Size: mygoui.ChipSmall,
				})
			},
			Body: func() {
				mygoui.Paragraph(c, "The toggle uses Enter or Space and animates unless reduced motion is on.", true)
			},
		})
		state := "closed"
		if st.disclosureOpen {
			state = "open"
		}
		mygoui.Typography(c, "Release notes are "+state+".", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// disclosureDisabledExample dims a trigger without an action.
func disclosureDisabledExample(c *ui.Context, _ *navigationState) {
	navigation.Disclosure(c, navigation.DisclosureOptions{
		Key:      "archived",
		Title:    "Archived projects",
		Disabled: true,
		Body: func() {
			mygoui.Paragraph(c, "Archived projects cannot change until they are restored.", true)
		},
	})
}

// disclosureGroupSingleExample closes siblings when one panel opens.
func disclosureGroupSingleExample(c *ui.Context, st *navigationState) {
	navigation.DisclosureGroup(c, navigation.DisclosureGroupOptions{
		Label: "Shipping",
		Items: disclosureItems(c,
			[]string{"Standard", "Express", "Pickup"},
			[]*bool{&st.groupSingle[0], &st.groupSingle[1], &st.groupSingle[2]},
			"delivery details."),
	})
}

// disclosureGroupMultipleExample allows several panels open together.
func disclosureGroupMultipleExample(c *ui.Context, st *navigationState) {
	navigation.DisclosureGroup(c, navigation.DisclosureGroupOptions{
		Multiple: true,
		Label:    "Filters",
		Items: disclosureItems(c,
			[]string{"Status", "Assignee"},
			[]*bool{&st.groupMultiple[0], &st.groupMultiple[1]},
			"filter options."),
	})
}

// disclosureGroupDisabledExample keeps one member unavailable.
func disclosureGroupDisabledExample(c *ui.Context, st *navigationState) {
	navigation.DisclosureGroup(c, navigation.DisclosureGroupOptions{
		Multiple: true,
		Label:    "Filters",
		Items: []navigation.DisclosureOptions{
			{Key: "status", Title: "Status", Open: &st.groupMultiple[0], Body: func() {
				mygoui.Paragraph(c, "Open and closed issues.", true)
			}},
			{Key: "archived", Title: "Archived", Disabled: true, Body: func() {
				mygoui.Paragraph(c, "Archived filters cannot change.", true)
			}},
		},
	})
}

// tabsPrimaryExample is the default pill tab list with a disabled tab.
func tabsPrimaryExample(c *ui.Context, st *navigationState) {
	navigation.Tabs(c, navigation.TabsOptions{
		Label:    "Account sections",
		Selected: &st.tabsPrimary,
		Items: []navigation.TabOptions{
			{Label: "Overview", Content: func() { mygoui.Paragraph(c, "Account overview and recent activity.", true) }},
			{Label: "Password", Content: func() { mygoui.Paragraph(c, "Change the password and review devices.", true) }},
			{Label: "Sessions", Content: func() { mygoui.Paragraph(c, "Two active sessions are signed in.", true) }},
			{Label: "Danger", Disabled: true},
		},
	})
}

// tabsSecondaryExample uses the underline variant.
func tabsSecondaryExample(c *ui.Context, st *navigationState) {
	navigation.Tabs(c, navigation.TabsOptions{
		Label:    "Report tabs",
		Variant:  navigation.TabsSecondary,
		Selected: &st.tabsSecondary,
		Items: []navigation.TabOptions{
			{Label: "Summary", Content: func() { mygoui.Paragraph(c, "Revenue is up 12% this quarter.", true) }},
			{Label: "Traffic", Content: func() { mygoui.Paragraph(c, "Sessions peak on Tuesday mornings.", true) }},
			{Label: "Retention", Content: func() { mygoui.Paragraph(c, "Monthly retention holds at 84%.", true) }},
		},
	})
}

// tabsVerticalExample places the list beside its panel.
func tabsVerticalExample(c *ui.Context, st *navigationState) {
	navigation.Tabs(c, navigation.TabsOptions{
		Label:       "Settings tabs",
		Orientation: mygoui.Vertical,
		Selected:    &st.tabsVertical,
		Items: []navigation.TabOptions{
			{Label: "Profile", Content: func() { mygoui.Paragraph(c, "Name, photo and public handle.", true) }},
			{Label: "Email", Content: func() { mygoui.Paragraph(c, "Notification address and signatures.", true) }},
			{Label: "Teams", Content: func() { mygoui.Paragraph(c, "Teams you own or belong to.", true) }},
		},
	})
}

// tabsOverflowExample lets a long list scroll horizontally.
func tabsOverflowExample(c *ui.Context, st *navigationState) {
	navigation.Tabs(c, navigation.TabsOptions{
		Label:    "Workspace tabs",
		Selected: &st.tabsOverflow,
		Items: []navigation.TabOptions{
			{Label: "Overview", Content: func() { mygoui.Paragraph(c, "Workspace overview.", true) }},
			{Label: "Activity", Content: func() { mygoui.Paragraph(c, "Recent activity.", true) }},
			{Label: "Members", Content: func() { mygoui.Paragraph(c, "Members and roles.", true) }},
			{Label: "Billing", Content: func() { mygoui.Paragraph(c, "Plan and invoices.", true) }},
			{Label: "Integrations", Content: func() { mygoui.Paragraph(c, "Connected apps.", true) }},
			{Label: "Security", Content: func() { mygoui.Paragraph(c, "SSO and audit rules.", true) }},
			{Label: "Audit log", Content: func() { mygoui.Paragraph(c, "Workspace events.", true) }},
			{Label: "Settings", Content: func() { mygoui.Paragraph(c, "Workspace settings.", true) }},
		},
	})
}

// breadcrumbsSimpleExample ends with a non-interactive current page.
func breadcrumbsSimpleExample(c *ui.Context, st *navigationState) {
	navigation.Breadcrumbs(c, navigation.BreadcrumbsOptions{
		Label: "Folder path",
		OnSelect: func(i int) {
			st.crumbNote = fmt.Sprintf("selected breadcrumb %d", i+1)
		},
		Items: []navigation.BreadcrumbOptions{
			{Label: "Workspace"},
			{Label: "Projects"},
			{Label: "mygoui"},
			{Label: "Navigation"},
		},
	})
}

// breadcrumbsOverflowExample folds a long middle into a menu.
func breadcrumbsOverflowExample(c *ui.Context, st *navigationState) {
	navigation.Breadcrumbs(c, navigation.BreadcrumbsOptions{
		Label:      "Deep folder path",
		MaxVisible: 4,
		OnSelect: func(i int) {
			st.crumbNote = fmt.Sprintf("selected breadcrumb %d", i+1)
		},
		Items: []navigation.BreadcrumbOptions{
			{Label: "Home"},
			{Label: "Library"},
			{Label: "Components"},
			{Label: "Navigation"},
			{Label: "Breadcrumbs"},
			{Label: "Examples"},
		},
	})
}

// breadcrumbsSelectionExample reports the chosen link.
func breadcrumbsSelectionExample(c *ui.Context, st *navigationState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		navigation.Breadcrumbs(c, navigation.BreadcrumbsOptions{
			Label: "Workspace path",
			OnSelect: func(i int) {
				st.crumbNote = fmt.Sprintf("selected breadcrumb %d", i+1)
			},
			Items: []navigation.BreadcrumbOptions{
				{Label: "Workspace"},
				{Label: "Projects"},
				{Label: "mygoui"},
			},
		})
		note := st.crumbNote
		if note == "" {
			note = "click a link to report its index"
		}
		mygoui.Chip(c, mygoui.ChipOptions{Label: note, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft})
	})
}

// paginationDefaultExample shows the centered page controls.
func paginationDefaultExample(c *ui.Context, st *navigationState) {
	navigation.Pagination(c, navigation.PaginationOptions{
		Label: "Pages", Page: &st.pageDefault, Pages: 5, Siblings: 1,
	})
}

// paginationSummaryExample pairs a result range with the controls, as
// upstream's Pagination.Summary does: the two split a wide row and stack in a
// narrow one.
func paginationSummaryExample(c *ui.Context, st *navigationState) {
	navigation.Pagination(c, navigation.PaginationOptions{
		Label:      "Search results",
		Page:       &st.pageSummary,
		Pages:      12,
		Siblings:   1,
		Boundaries: 1,
		Summary: func() {
			mygoui.Typography(c,
				fmt.Sprintf("Showing %d-%d of 120 results", (st.pageSummary-1)*10+1, st.pageSummary*10),
				mygoui.TypographyOptions{Type: theme.TypeBodySM, Muted: true})
		},
	})
}

// paginationLargeExample shows a large twelve-page range.
func paginationLargeExample(c *ui.Context, st *navigationState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		navigation.Pagination(c, navigation.PaginationOptions{
			Label:         "Search results",
			Page:          &st.page,
			Pages:         12,
			Size:          navigation.PaginationLarge,
			Siblings:      1,
			Boundaries:    2,
			PreviousLabel: "Previous",
			NextLabel:     "Next",
		})
		mygoui.Chip(c, mygoui.ChipOptions{
			Label: fmt.Sprintf("page %d of 12", st.page),
			Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft,
		})
	})
}

// paginationMediumExample uses the default control height.
func paginationMediumExample(c *ui.Context, st *navigationState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		navigation.Pagination(c, navigation.PaginationOptions{
			Label:      "Team members",
			Page:       &st.pageMedium,
			Pages:      8,
			Siblings:   1,
			Boundaries: 2,
		})
		mygoui.Typography(c, fmt.Sprintf("Medium size, page %d", st.pageMedium), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// paginationSmallExample is a compact five-page range.
func paginationSmallExample(c *ui.Context, st *navigationState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		navigation.Pagination(c, navigation.PaginationOptions{
			Label:      "Recent results",
			Page:       &st.pageSmall,
			Pages:      5,
			Size:       navigation.PaginationSmall,
			Siblings:   0,
			Boundaries: 1,
		})
		mygoui.Typography(c, fmt.Sprintf("Small size, page %d", st.pageSmall), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// toolbarAttachedExample joins the formatting controls into one surface.
func toolbarAttachedExample(c *ui.Context, st *navigationState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		navigation.Toolbar(c, navigation.ToolbarOptions{
			Label:    "Text formatting",
			Attached: true,
			Items: func() {
				toolbarIcon(c, st, "Bold", lucide.Bold())
				toolbarIcon(c, st, "Italic", lucide.Italic())
				toolbarIcon(c, st, "Underline", lucide.Underline())
				navigation.ToolbarSeparator(c, mygoui.Horizontal)
				toolbarIcon(c, st, "Align left", lucide.AlignLeft())
				navigation.ToolbarSeparator(c, mygoui.Horizontal)
				toolbarIcon(c, st, "Copy", lucide.Copy())
			},
		})
		action := st.toolbarAction
		if action == "" {
			action = "no toolbar action yet"
		}
		mygoui.Chip(c, mygoui.ChipOptions{Label: action, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft})
	})
}

// toolbarVerticalExample stacks clipboard controls on one column.
func toolbarVerticalExample(c *ui.Context, st *navigationState) {
	navigation.Toolbar(c, navigation.ToolbarOptions{
		Label:       "Clipboard",
		Orientation: mygoui.Vertical,
		Items: func() {
			toolbarIcon(c, st, "Cut", lucide.Scissors())
			toolbarIcon(c, st, "Paste", lucide.Clipboard())
			navigation.ToolbarSeparator(c, mygoui.Vertical)
			toolbarIcon(c, st, "Delete", lucide.Trash())
		},
	})
}

// toolbarFloatingExample leaves the controls on separate surfaces.
func toolbarFloatingExample(c *ui.Context, st *navigationState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		navigation.Toolbar(c, navigation.ToolbarOptions{
			Label: "History controls",
			Items: func() {
				toolbarIcon(c, st, "Undo", lucide.Undo())
				toolbarIcon(c, st, "Redo", lucide.Redo())
				navigation.ToolbarSeparator(c, mygoui.Horizontal)
				toolbarIcon(c, st, "Search", lucide.Search())
			},
		})
		mygoui.Typography(c, "Arrow keys move between the controls in one focus stop.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// menuCheckboxExample keeps several workspace items checked.
func menuCheckboxExample(c *ui.Context, st *navigationState) {
	ui.Column(c).Width(theme.Scale(c, 320)).MaxWidth(theme.Scale(c, 520)).Gap(theme.Space8).Children(func() {
		menu := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Radius(theme.RadiusScale().XXL).
			Background(theme.Colors(c).Overlay).TextColor(theme.Colors(c).OverlayForeground)
		menu.Children(func() {
			navigation.Menu(c, navigation.MenuOptions{
				Label:     "Workspace menu",
				Selection: navigation.MenuSelectMultiple,
				Selected:  &st.menuSelected,
				Sections: []navigation.MenuSectionOptions{
					{
						Title: "Workspace",
						Items: []navigation.MenuItemOptions{
							{Label: "Overview", Icon: lucide.Settings(), Shortcut: "\u23181"},
							{Label: "Reports", Icon: lucide.Clipboard()},
							{Separator: true},
							{Label: "Notifications", Checked: true},
							{Label: "Weekly digest", Checked: true},
							{Label: "Archived items", Disabled: true},
						},
					},
				},
			})
		})
		mygoui.Typography(c, selectionSummary("checked", st.menuSelected), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// menuRadioExample keeps exactly one sort row chosen.
func menuRadioExample(c *ui.Context, st *navigationState) {
	ui.Column(c).Width(theme.Scale(c, 260)).Gap(theme.Space8).Children(func() {
		menu := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Radius(theme.RadiusScale().XXL).
			Background(theme.Colors(c).Overlay).TextColor(theme.Colors(c).OverlayForeground)
		menu.Children(func() {
			navigation.Menu(c, navigation.MenuOptions{
				Label:     "Sort order",
				Selection: navigation.MenuSelectSingle,
				Selected:  &st.sortSelected,
				OnSelect:  func(key string) { st.sortChoice = "sorted by " + key },
				Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
					{Key: "Name", Label: "Name", Shortcut: "\u23181"},
					{Key: "Modified", Label: "Last modified", Shortcut: "\u23182"},
					{Key: "Size", Label: "Size"},
				}}},
			})
		})
		mygoui.Typography(c, st.sortChoice, mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// menuSectionsExample groups items under headings.
func menuSectionsExample(c *ui.Context, _ *navigationState) {
	ui.Column(c).Width(theme.Scale(c, 280)).Children(func() {
		menu := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Radius(theme.RadiusScale().XXL).
			Background(theme.Colors(c).Overlay).TextColor(theme.Colors(c).OverlayForeground)
		menu.Children(func() {
			navigation.Menu(c, navigation.MenuOptions{
				Label: "Project menu",
				Sections: []navigation.MenuSectionOptions{
					{Title: "Project", Items: []navigation.MenuItemOptions{
						{Label: "Open project", Icon: lucide.Folder()},
						{Label: "Copy link", Icon: lucide.Copy(), Shortcut: "\u2318C"},
					}},
					{Title: "Danger zone", Items: []navigation.MenuItemOptions{
						{Label: "Archive", Icon: lucide.Archive()},
						{Label: "Delete", Icon: lucide.Trash(), Variant: navigation.MenuItemDanger},
					}},
				},
			})
		})
	})
}

// dropdownRadioExample anchors a single-choice menu to its trigger.
func dropdownRadioExample(c *ui.Context, st *navigationState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		navigation.Dropdown(c, navigation.DropdownOptions{
			Open:  &st.sortOpen,
			Label: "Sort by",
			Menu: navigation.MenuOptions{
				Label:     "Sort order",
				Selection: navigation.MenuSelectSingle,
				Selected:  &st.sortSelected,
				OnSelect: func(key string) {
					st.sortChoice = "sorted by " + key
				},
				Sections: []navigation.MenuSectionOptions{
					{Items: []navigation.MenuItemOptions{
						{Key: "Name", Label: "Name", Shortcut: "\u23181"},
						{Key: "Modified", Label: "Last modified", Shortcut: "\u23182"},
						{Key: "Size", Label: "Size"},
					}},
				},
			},
		})
		if st.sortChoice != "" {
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: st.sortChoice, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft,
			})
		}
	})
}

// dropdownIconExample adds a leading trigger icon and menu icons.
func dropdownIconExample(c *ui.Context, st *navigationState) {
	navigation.Dropdown(c, navigation.DropdownOptions{
		Open:  &st.sortOpenAlt,
		Label: "Filter",
		Icon:  lucide.Filter(),
		Menu: navigation.MenuOptions{
			Label: "Filter rows",
			Sections: []navigation.MenuSectionOptions{
				{Title: "Show", Items: []navigation.MenuItemOptions{
					{Label: "Healthy", Icon: lucide.CheckCircle()},
					{Label: "Degraded", Icon: lucide.AlertTriangle()},
					{Label: "Failed", Icon: lucide.XCircle()},
				}},
			},
		},
	})
}

// dropdownDisabledExample dims the trigger and blocks opening.
func dropdownDisabledExample(c *ui.Context, _ *navigationState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		navigation.Dropdown(c, navigation.DropdownOptions{
			Disabled: true,
			Label:    "Sort by",
			Menu: navigation.MenuOptions{
				Label:    "Sort order",
				Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{{Label: "Name"}}}},
			},
		})
		mygoui.Typography(c, "A disabled trigger cannot open its menu.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// submenuActionsExample anchors actions beside their trigger.
func submenuActionsExample(c *ui.Context, st *navigationState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			anchor := mygoui.Button(c, mygoui.ButtonOptions{
				Label:   "Share options",
				Icon:    lucide.Share(),
				Variant: mygoui.ButtonSecondary,
				OnClick: func() { st.submenuOpen = !st.submenuOpen },
			})
			navigation.SubMenu(c, anchor, navigation.SubMenuOptions{
				Open:  &st.submenuOpen,
				Label: "Share options",
				Items: []navigation.MenuItemOptions{
					{Label: "Copy link", Icon: lucide.Copy(), Shortcut: "\u2318C", OnSelect: func() {
						st.submenuChoice = "Copy link"
						st.submenuOpen = false
					}},
					{Label: "Email invite", Icon: lucide.ExternalLink(), OnSelect: func() {
						st.submenuChoice = "Email invite"
						st.submenuOpen = false
					}},
					{Label: "Embed code", OnSelect: func() {
						st.submenuChoice = "Embed code"
						st.submenuOpen = false
					}},
				},
			})
		})
		if st.submenuChoice != "" {
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: "action: " + st.submenuChoice, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft,
			})
		}
	})
}

// submenuDisabledExample leaves an unavailable action in the panel.
func submenuDisabledExample(c *ui.Context, st *navigationState) {
	anchor := mygoui.Button(c, mygoui.ButtonOptions{
		Label:   "More actions",
		Icon:    lucide.MoreHorizontal(),
		Variant: mygoui.ButtonSecondary,
		OnClick: func() { st.submenuOpenAlt = !st.submenuOpenAlt },
	})
	navigation.SubMenu(c, anchor, navigation.SubMenuOptions{
		Open:  &st.submenuOpenAlt,
		Label: "More actions",
		Items: []navigation.MenuItemOptions{
			{Label: "Duplicate", Icon: lucide.Copy(), OnSelect: func() { st.submenuOpenAlt = false }},
			{Label: "Move to archive", Icon: lucide.Archive(), Disabled: true},
			{Label: "Delete", Icon: lucide.Trash(), Variant: navigation.MenuItemDanger, OnSelect: func() {
				st.submenuOpenAlt = false
			}},
		},
	})
}

// disclosureItems builds a coordinated group's caller-owned panels.
func disclosureItems(c *ui.Context, titles []string, open []*bool, suffix string) []navigation.DisclosureOptions {
	items := make([]navigation.DisclosureOptions, 0, len(titles))
	for i, title := range titles {
		items = append(items, navigation.DisclosureOptions{
			Key: title, Title: title, Open: open[i],
			Body: func() { mygoui.Paragraph(c, fmt.Sprintf("%s %s", title, suffix), true) },
		})
	}
	return items
}
