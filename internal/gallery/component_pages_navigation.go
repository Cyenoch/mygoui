package gallery

import (
	"github.com/egoist/mygo/ui"
)

// navigationStateOf returns the shared navigation state with the same
// first-visit defaults the family page applies, so a component page opened
// directly starts from a populated, realistic demo.
func navigationStateOf(s *State) *navigationState {
	st := StateOf[navigationState](s, "navigation")
	if st.page == 0 {
		st.page = 4
	}
	if st.pageMedium == 0 {
		st.pageMedium = 3
	}
	if st.pageSmall == 0 {
		st.pageSmall = 2
	}
	if st.pageSummary == 0 {
		st.pageSummary = 2
	}
	if !st.menuInit {
		st.menuSelected = []string{"Notifications", "Weekly digest"}
		st.menuInit = true
	}
	if !st.sortInit {
		st.sortSelected = []string{"Modified"}
		st.sortChoice = "sorted by Modified"
		st.sortInit = true
	}
	return st
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "accordion", Title: "Accordion", Group: "Navigation", Order: 10,
		Description: "Collapsible sections with single or multiple expansion, header actions and disabled rows.",
		Examples: []ComponentExample{
			{Title: "Single expansion", Description: "One section stays open; opening another closes the rest.",
				Demo: func(c *ui.Context, s *State) { accordionSingleExample(c, navigationStateOf(s)) }},
			{Title: "Multiple open", Description: "A surface accordion keeps several sections open with a header action.",
				Demo: func(c *ui.Context, s *State) { accordionMultipleExample(c, navigationStateOf(s)) }},
			{Title: "Without separators", Description: "The hairline between sections is removed.",
				Demo: func(c *ui.Context, s *State) { accordionSeparatorExample(c, navigationStateOf(s)) }},
			{Title: "Disabled section", Description: "A section can stay visible but unavailable.",
				Demo: func(c *ui.Context, s *State) { accordionDisabledExample(c, navigationStateOf(s)) }},
		},
		Code: `expanded := []string{"general"}
navigation.Accordion(c, navigation.AccordionOptions{
    Label:    "Project settings",
    Expanded: &expanded,
    Items: []navigation.AccordionItemOptions{
        {Key: "general", Title: "General", Body: func() {
            mygoui.Paragraph(c, "The project name and owner.", true)
        }},
        {Key: "billing", Title: "Billing", Body: func() {
            mygoui.Paragraph(c, "Invoices and the current plan.", true)
        }},
    },
})`,
		Params: []ComponentParam{
			{"Variant", "navigation.AccordionVariant", "AccordionDefault", "Default on the window background, or AccordionSurface on a surface container."},
			{"Multiple", "bool", "false", "Keeps several sections open at once."},
			{"Expanded", "*[]string", "nil", "Caller-owned keys of the expanded sections; nil keeps it internal."},
			{"Items", "[]navigation.AccordionItemOptions", "required", "Sections in order, each with key, title and body."},
			{"HideSeparator", "bool", "false", "Removes the hairline between sections."},
			{"OnChange", "func(key string, expanded bool)", "nil", "Reports a section toggled after the state changed."},
		},
		Tokens: []string{"Background", "Surface", "SurfaceForeground", "Foreground", "Muted", "Separator"},
	})

	RegisterComponent(ComponentPage{
		ID: "disclosure", Title: "Disclosure", Group: "Navigation", Order: 20,
		Description: "A single header trigger over a panel that grows into view.",
		Examples: []ComponentExample{
			{Title: "With actions", Description: "One caller-owned panel plus a disabled trigger.",
				Demo: func(c *ui.Context, s *State) { disclosureActionsExample(c, navigationStateOf(s)) }},
			{Title: "Disabled trigger", Description: "A dimmed header keeps its panel closed.",
				Demo: func(c *ui.Context, s *State) { disclosureDisabledExample(c, navigationStateOf(s)) }},
		},
		Code: `open := false
navigation.Disclosure(c, navigation.DisclosureOptions{
    Key:   "release-notes",
    Title: "Release notes",
    Open:  &open,
    Actions: func() {
        mygoui.Chip(c, mygoui.ChipOptions{Label: "v3.2", Variant: mygoui.ChipSoft})
    },
    Body: func() {
        mygoui.Paragraph(c, "Enter or Space toggles the panel.", true)
    },
})`,
		Params: []ComponentParam{
			{"Title", "string", "required", "Visible trigger text and the accessible name."},
			{"Open", "*bool", "nil", "Caller-owned expansion; nil keeps it internal."},
			{"Actions", "func()", "nil", "Trailing header controls outside the disclosure toggle."},
			{"Body", "func()", "required", "Panel content revealed while open."},
			{"Disabled", "bool", "false", "Dims the trigger and blocks toggling."},
		},
		Tokens: []string{"Foreground", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "disclosuregroup", Title: "DisclosureGroup", Group: "Navigation", Order: 30,
		Description: "A coordinated group of disclosures that closes its siblings in single mode.",
		Examples: []ComponentExample{
			{Title: "Single open", Description: "Single-open and multi-open coordination from caller-owned flags.",
				Demo: func(c *ui.Context, s *State) { disclosureGroupSingleExample(c, navigationStateOf(s)) }},
			{Title: "Multiple open", Description: "Several panels can stay open at the same time.",
				Demo: func(c *ui.Context, s *State) { disclosureGroupMultipleExample(c, navigationStateOf(s)) }},
			{Title: "Disabled item", Description: "One member stays unavailable while its siblings work.",
				Demo: func(c *ui.Context, s *State) { disclosureGroupDisabledExample(c, navigationStateOf(s)) }},
		},
		Code: `open := [2]bool{}
navigation.DisclosureGroup(c, navigation.DisclosureGroupOptions{
    Label:    "Shipping",
    Multiple: false,
    Items: []navigation.DisclosureOptions{
        {Key: "standard", Title: "Standard", Open: &open[0], Body: func() {
            mygoui.Paragraph(c, "Arrives in five days.", true)
        }},
        {Key: "express", Title: "Express", Open: &open[1], Body: func() {
            mygoui.Paragraph(c, "Arrives tomorrow.", true)
        }},
    },
})`,
		Params: []ComponentParam{
			{"Multiple", "bool", "false", "Allows several items open; otherwise opening one closes the rest."},
			{"Label", "string", "\"\"", "Accessible name of the disclosure group."},
			{"Items", "[]navigation.DisclosureOptions", "required", "The disclosures and their caller-owned open flags."},
		},
		Tokens: []string{"Foreground", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "tabs", Title: "Tabs", Group: "Navigation", Order: 40,
		Description: "Primary, secondary and vertical tab lists with a scrolling overflow row.",
		Examples: []ComponentExample{
			{Title: "Primary", Description: "Pill tabs with one disabled section.",
				Demo: func(c *ui.Context, s *State) { tabsPrimaryExample(c, navigationStateOf(s)) }},
			{Title: "Secondary", Description: "An accent underline marks the selected tab.",
				Demo: func(c *ui.Context, s *State) { tabsSecondaryExample(c, navigationStateOf(s)) }},
			{Title: "Vertical", Description: "The list sits beside its panel.",
				Demo: func(c *ui.Context, s *State) { tabsVerticalExample(c, navigationStateOf(s)) }},
			{Title: "Overflow", Description: "A long list scrolls horizontally instead of wrapping.",
				Demo: func(c *ui.Context, s *State) { tabsOverflowExample(c, navigationStateOf(s)) }},
		},
		Code: `selected := 0
navigation.Tabs(c, navigation.TabsOptions{
    Label:    "Account sections",
    Selected: &selected,
    Items: []navigation.TabOptions{
        {Label: "Overview", Content: func() {
            mygoui.Paragraph(c, "Recent account activity.", true)
        }},
        {Label: "Password", Content: func() {
            mygoui.Paragraph(c, "Change the password.", true)
        }},
        {Label: "Danger", Disabled: true},
    },
})`,
		Params: []ComponentParam{
			{"Variant", "navigation.TabsVariant", "TabsPrimary", "Pill selection or an accent underline."},
			{"Orientation", "mygoui.Orientation", "mygoui.Horizontal", "Horizontal list above the panel, or a vertical list beside it."},
			{"Selected", "*int", "nil", "Caller-owned selected index; nil keeps it internal."},
			{"Items", "[]navigation.TabOptions", "required", "Tabs with label, content builder and disabled flag."},
			{"Align", "navigation.TabsAlign", "TabsAlignCenter", "Centers, starts or ends each label in its box."},
			{"MaxListHeight", "float32", "0", "Bounds a vertical list so it scrolls."},
		},
		Tokens: []string{"Accent", "Default", "Border", "Segment", "SegmentForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "breadcrumbs", Title: "Breadcrumbs", Group: "Navigation", Order: 50,
		Description: "A path of links with an overflow menu and a non-interactive current page.",
		Examples: []ComponentExample{
			{Title: "Simple path", Description: "The last item is the current page and is not interactive.",
				Demo: func(c *ui.Context, s *State) { breadcrumbsSimpleExample(c, navigationStateOf(s)) }},
			{Title: "Overflow menu", Description: "A long path keeps its ends and folds the middle into a menu.",
				Demo: func(c *ui.Context, s *State) { breadcrumbsOverflowExample(c, navigationStateOf(s)) }},
			{Title: "Selection readout", Description: "Choosing a link reports its index.",
				Demo: func(c *ui.Context, s *State) { breadcrumbsSelectionExample(c, navigationStateOf(s)) }},
		},
		Code: `navigation.Breadcrumbs(c, navigation.BreadcrumbsOptions{
    Label:      "Folder path",
    MaxVisible: 4,
    OnSelect:   func(index int) { selectCrumb(index) },
    Items: []navigation.BreadcrumbOptions{
        {Label: "Workspace"},
        {Label: "Projects"},
        {Label: "mygoui"},
        {Label: "Navigation"},
    },
})`,
		Params: []ComponentParam{
			{"Items", "[]navigation.BreadcrumbOptions", "required", "The path in order; the last item is the current page."},
			{"MaxVisible", "int", "0", "Keeps the first and last items plus an overflow menu; zero shows all."},
			{"Label", "string", "\"\"", "Accessible name of the breadcrumb group."},
			{"MoreLabel", "string", "locale MorePages", "Accessible name of the overflow button."},
			{"OnSelect", "func(index int)", "nil", "Reports the chosen breadcrumb, including from overflow."},
		},
		Tokens: []string{"Link", "Muted", "DefaultHover"},
	})

	RegisterComponent(ComponentPage{
		ID: "pagination", Title: "Pagination", Group: "Navigation", Order: 60,
		Description: "Numbered pages with ellipses, previous and next controls in three sizes and an optional summary.",
		Examples: []ComponentExample{
			{Title: "With summary", Description: "A result range beside the controls; it stacks when the page is narrow.",
				Demo: func(c *ui.Context, s *State) { paginationSummaryExample(c, navigationStateOf(s)) }},
			{Title: "Large range", Description: "Twelve pages with two boundaries and one sibling.",
				Demo: func(c *ui.Context, s *State) { paginationLargeExample(c, navigationStateOf(s)) }},
			{Title: "Medium range", Description: "The default control height with eight pages.",
				Demo: func(c *ui.Context, s *State) { paginationMediumExample(c, navigationStateOf(s)) }},
			{Title: "Small range", Description: "A compact five-page range without siblings.",
				Demo: func(c *ui.Context, s *State) { paginationSmallExample(c, navigationStateOf(s)) }},
		},
		Code: `page := 4
navigation.Pagination(c, navigation.PaginationOptions{
    Label:   "Search results",
    Page:    &page,
    Pages:   12,
    Siblings: 1,
    Summary: func() {
        mygoui.Typography(c, "Showing 31-40 of 120 results", mygoui.TypographyOptions{
            Type: theme.TypeBodySM, Muted: true,
        })
    },
})`,
		Params: []ComponentParam{
			{"Page", "*int", "nil", "Caller-owned 1-based page; nil keeps it internal."},
			{"Pages", "int", "required", "Total number of pages."},
			{"Size", "navigation.PaginationSize", "PaginationMedium", "Small, medium or large control height."},
			{"Siblings", "int", "1", "Number of pages shown beside the current one."},
			{"Boundaries", "int", "1", "Number of pages pinned at each end."},
			{"Summary", "func()", "nil", "Builds a muted line before the controls; the row stacks below 640 DIP."},
			{"OnChange", "func(page int)", "nil", "Reports the newly selected page."},
		},
		Tokens: []string{"Default", "DefaultForeground", "DefaultHover", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "toolbar", Title: "Toolbar", Group: "Navigation", Order: 70,
		Description: "A group of controls in one focus stop, with arrow navigation and separators.",
		Examples: []ComponentExample{
			{Title: "Attached", Description: "The formatting controls share one surface and one focus stop.",
				Demo: func(c *ui.Context, s *State) { toolbarAttachedExample(c, navigationStateOf(s)) }},
			{Title: "Vertical", Description: "Up and Down arrows move through a stacked toolbar.",
				Demo: func(c *ui.Context, s *State) { toolbarVerticalExample(c, navigationStateOf(s)) }},
			{Title: "Floating", Description: "The controls stay on their own surfaces with separators.",
				Demo: func(c *ui.Context, s *State) { toolbarFloatingExample(c, navigationStateOf(s)) }},
		},
		Code: `navigation.Toolbar(c, navigation.ToolbarOptions{
    Label:    "Text formatting",
    Attached: true,
    Items: func() {
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Bold", Icon: lucide.Bold(), IconOnly: true,
            Variant: mygoui.ButtonGhost, OnClick: bold,
        })
        navigation.ToolbarSeparator(c, mygoui.Horizontal)
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Copy", Icon: lucide.Copy(), IconOnly: true,
            Variant: mygoui.ButtonGhost, OnClick: copySelection,
        })
    },
})`,
		Params: []ComponentParam{
			{"Items", "func()", "required", "Builds the controls that share the focus group."},
			{"Orientation", "mygoui.Orientation", "mygoui.Horizontal", "Horizontal arrows move left and right; vertical ones move up and down."},
			{"Attached", "bool", "false", "Joins the controls into one surface with separators."},
			{"Label", "string", "\"\"", "Accessible name of the toolbar."},
		},
		Tokens: []string{"Background", "Surface", "SurfaceForeground", "Border", "Separator"},
	})

	RegisterComponent(ComponentPage{
		ID: "menu", Title: "Menu", Group: "Navigation", Order: 80,
		Description: "A vertical menu with sections, shortcuts, icons and checkbox selection.",
		Examples: []ComponentExample{
			{Title: "Checkbox selection", Description: "Sections, icons, shortcuts and multiple checked items.",
				Demo: func(c *ui.Context, s *State) { menuCheckboxExample(c, navigationStateOf(s)) }},
			{Title: "Radio selection", Description: "A single-choice menu keeps one row chosen.",
				Demo: func(c *ui.Context, s *State) { menuRadioExample(c, navigationStateOf(s)) }},
			{Title: "Section headings", Description: "Items group under headings, including a danger action.",
				Demo: func(c *ui.Context, s *State) { menuSectionsExample(c, navigationStateOf(s)) }},
		},
		Code: `selected := []string{"Notifications"}
navigation.Menu(c, navigation.MenuOptions{
    Label:     "Workspace menu",
    Selection: navigation.MenuSelectMultiple,
    Selected:  &selected,
    Sections: []navigation.MenuSectionOptions{{
        Title: "Workspace",
        Items: []navigation.MenuItemOptions{
            {Label: "Overview", Icon: lucide.Settings(), Shortcut: "⌘1"},
            {Separator: true},
            {Label: "Notifications", Checked: true},
            {Label: "Archived items", Disabled: true},
        },
    }},
})`,
		Params: []ComponentParam{
			{"Sections", "[]navigation.MenuSectionOptions", "required", "Grouped items with optional section headings."},
			{"Selection", "navigation.MenuSelection", "MenuSelectNone", "None, single (radio) or multiple (checkbox) selection."},
			{"Selected", "*[]string", "nil", "Caller-owned keys of checked items; nil keeps them internal."},
			{"Label", "string", "\"\"", "Accessible name of the menu."},
			{"OnSelect", "func(key string)", "nil", "Reports an activated item in addition to its own callback."},
		},
		Tokens: []string{"Default", "DefaultForeground", "Muted", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "dropdown", Title: "Dropdown", Group: "Navigation", Order: 90,
		Description: "A menu button and its anchored menu panel, closed by an outside press or Escape.",
		Examples: []ComponentExample{
			{Title: "Radio menu", Description: "Click the trigger, choose a row, or press Escape to close the panel.",
				Demo: func(c *ui.Context, s *State) { dropdownRadioExample(c, navigationStateOf(s)) }},
			{Title: "With icons", Description: "A leading trigger icon and menu icons.",
				Demo: func(c *ui.Context, s *State) { dropdownIconExample(c, navigationStateOf(s)) }},
			{Title: "Disabled trigger", Description: "A dimmed trigger cannot open its menu.",
				Demo: func(c *ui.Context, s *State) { dropdownDisabledExample(c, navigationStateOf(s)) }},
		},
		Code: `open := false
selected := []string{"Modified"}
navigation.Dropdown(c, navigation.DropdownOptions{
    Open:  &open,
    Label: "Sort by",
    Menu: navigation.MenuOptions{
        Label:     "Sort order",
        Selection: navigation.MenuSelectSingle,
        Selected:  &selected,
        Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
            {Key: "Name", Label: "Name"},
            {Key: "Modified", Label: "Last modified"},
            {Key: "Size", Label: "Size"},
        }}},
    },
})`,
		Params: []ComponentParam{
			{"Open", "*bool", "nil", "Caller-owned panel visibility; nil keeps it internal."},
			{"Label", "string", "required", "Trigger text and accessible name."},
			{"Menu", "navigation.MenuOptions", "required", "Contents and selection model of the anchored menu."},
			{"Icon", "*ui.SVG", "nil", "Leading icon on the trigger button."},
			{"Placement", "mygoui.Placement", "SideBottom", "Preferred side and alignment of the panel."},
			{"Disabled", "bool", "false", "Dims the trigger and blocks opening."},
		},
		Tokens: []string{"Default", "DefaultForeground", "Muted", "Danger", "Overlay"},
	})

	RegisterComponent(ComponentPage{
		ID: "submenu", Title: "SubMenu", Group: "Navigation", Order: 100,
		Description: "A nested menu panel anchored beside a trigger, closed by Escape or Left.",
		Examples: []ComponentExample{
			{Title: "With shortcuts", Description: "Escape or Left closes the panel without choosing.",
				Demo: func(c *ui.Context, s *State) { submenuActionsExample(c, navigationStateOf(s)) }},
			{Title: "Disabled action", Description: "An unavailable row stays in the nested panel.",
				Demo: func(c *ui.Context, s *State) { submenuDisabledExample(c, navigationStateOf(s)) }},
		},
		Code: `open := false
anchor := mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Share options", Icon: lucide.Share(),
    Variant: mygoui.ButtonSecondary,
    OnClick: func() { open = !open },
})
navigation.SubMenu(c, anchor, navigation.SubMenuOptions{
    Open:  &open,
    Label: "Share options",
    Items: []navigation.MenuItemOptions{
        {Label: "Copy link", Icon: lucide.Copy(), Shortcut: "⌘C", OnSelect: copyLink},
        {Label: "Email invite", Icon: lucide.ExternalLink(), OnSelect: invite},
    },
})`,
		Params: []ComponentParam{
			{"Open", "*bool", "required", "Panel visibility; nil renders nothing."},
			{"Items", "[]navigation.MenuItemOptions", "required", "Rows of the nested panel."},
			{"Label", "string", "\"\"", "Accessible name of the submenu."},
		},
		Tokens: []string{"Default", "DefaultForeground", "Muted", "Overlay"},
	})
}
