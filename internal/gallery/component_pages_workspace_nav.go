package gallery

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// workspaceNavState is the caller-owned state behind the document tabs,
// navigation stack, carousel and status bar demos.
type workspaceNavState struct {
	ready bool

	docs      []navigation.DocumentTab
	activeDoc string
	created   int

	manyDocs      []navigation.DocumentTab
	activeManyDoc string

	stack []string

	slide, cards, pip, many int

	message     string
	saves       int
	branch      int
	unreadCount int
	syncing     bool
	tone        navigation.StatusBarTone
}

func workspaceNavStateOf(s *State) *workspaceNavState {
	st := StateOf[workspaceNavState](s, "workspace-nav")
	if !st.ready {
		st.ready = true
		st.docs = []navigation.DocumentTab{
			{Key: "main", Title: "main.go", Icon: lucide.FileCode(), Pinned: true, Tooltip: "internal/app/main.go"},
			{Key: "readme", Title: "README.md", Icon: lucide.FileText(), Dirty: true},
			{Key: "mod", Title: "go.mod", Icon: lucide.File()},
		}
		st.activeDoc = "main"
		for i := range 40 {
			key := strconv.Itoa(i)
			st.manyDocs = append(st.manyDocs, navigation.DocumentTab{Key: key, Title: "chapter-" + key + ".md"})
		}
		st.activeManyDoc = "0"
		st.stack = []string{"Settings"}
		st.message = "Ready"
		st.unreadCount = 3
	}
	return st
}

func workspaceNavDemo(body func(c *ui.Context, st *workspaceNavState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, workspaceNavStateOf(s)) }
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "documenttabs", Title: "DocumentTabs", Group: "Patterns", Order: 170,
		Description: "A strip of open documents that close, reorder and scroll.",
		Examples: []ComponentExample{
			{Title: "Documents", Description: "Drag a tab, or press Alt+Left and Alt+Right, to move it. Right-click for more.",
				Demo: workspaceNavDemo(documentTabsBody)},
			{Title: "Overflow", Description: "The strip fades at its edges, scrolls with the wheel and lists every tab.",
				Demo: workspaceNavDemo(documentTabsOverflowBody)},
			{Title: "Shrink to fit", Demo: workspaceNavDemo(documentTabsShrinkBody)},
		},
		Code: `var tabs []navigation.DocumentTab
var active string
var closeTab func(key string)
var moveTab func(from, to int)
var newTab func()

navigation.DocumentTabs(c, navigation.DocumentTabsOptions{
    Label:    "Open documents",
    Tabs:     tabs,
    Active:   active,
    OnSelect: func(key string) { active = key },
    OnClose:  closeTab,
    OnMove:   moveTab,
    OnAdd:    newTab,
    // Close others and Close to the right arrive together.
    OnCloseMany: func(keys []string) { /* drop them */ },
    Sizing:      navigation.DocumentTabsShrink,
})`,
		Params: []ComponentParam{
			{"Tabs", "[]navigation.DocumentTab", "nil", "Open documents: Key, Title, Icon, Dirty, Pinned (icon only, kept first, no close) and Tooltip. Never modified."},
			{"Active", "string", "\"\"", "Key of the selected tab; it is kept in view."},
			{"Label", "string", "\"\"", "Accessible name of the tab list."},
			{"OnSelect", "func(key string)", "nil", "A tab was chosen by click, Enter or Space."},
			{"OnClose", "func(key string)", "nil", "Adds close buttons, middle click, Cmd/Ctrl+W and Delete."},
			{"OnMove", "func(from, to int)", "nil", "Enables dragging and Alt+Left / Alt+Right; to counts after removal."},
			{"OnAdd", "func()", "nil", "Adds the trailing add button."},
			{"OnCloseMany", "func(keys []string)", "nil", "Close others and Close to the right; pinned tabs are left out. Without it, OnClose is called for each."},
			{"Sizing", "navigation.DocumentTabSizing", "DocumentTabsFixed", "DocumentTabsShrink shares the strip among the tabs before it scrolls."},
			{"WindowShortcuts", "bool", "false", "Ctrl+Tab, Ctrl+Shift+Tab and Cmd/Ctrl+W work anywhere in the window."},
		},
		Tokens: []string{"MutedBackground", "MutedForeground", "Background", "Foreground", "Input", "Accent", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "navigationstack", Title: "NavigationStack", Group: "Patterns", Order: 180,
		Description: "A drill-down container that slides between a stack of pages.",
		Examples: []ComponentExample{
			{Title: "Drill-down", Description: "Alt+Left, the mouse back button or the back button returns.",
				Demo: workspaceNavDemo(navigationStackBody)},
		},
		Code: `var pages []navigation.NavigationPage

navigation.NavigationStack(c, navigation.NavigationStackOptions{
    Label: "Settings",
    Pages: pages,
    OnPop: func() { pages = pages[:len(pages)-1] },
    Actions: func() { /* trailing header controls */ },
})`,
		Params: []ComponentParam{
			{"Pages", "[]navigation.NavigationPage", "nil", "The stack: Key, Title and Build. The last page is shown."},
			{"Label", "string", "\"\"", "Accessible name of the stack."},
			{"KeepMounted", "bool", "false", "Builds covered pages hidden, keeping their scroll and editing state."},
			{"OnPop", "func()", "nil", "Adds the back button, whose tooltip names the page behind, the mouse's back button and Alt+Left; the caller drops the last page."},
			{"Actions", "func()", "nil", "Builds controls at the trailing end of the header."},
		},
		Tokens: []string{"Foreground", "Accent", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "carousel", Title: "Carousel", Group: "Navigation", Order: 190,
		Description: "Slides stepped with buttons, arrow keys and a pips pager.",
		Examples: []ComponentExample{
			{Title: "Looping", Demo: workspaceNavDemo(carouselLoopBody)},
			{Title: "Several per view", Demo: workspaceNavDemo(carouselPerViewBody)},
			{Title: "Many slides", Description: "Past seven pages the pips are a window, so the position is said as well. Home and End jump.",
				Demo: workspaceNavDemo(carouselManyBody)},
			{Title: "PipsPager", Description: "The pager on its own, windowed over many pages.",
				Demo: workspaceNavDemo(pipsPagerBody)},
		},
		Code: `var index int
var slide func(index int)

navigation.Carousel(c, navigation.CarouselOptions{
    Label: "Highlights",
    Index: &index,
    Count: 5,
    Loop:  true,
    Slide: slide,
})

navigation.PipsPager(c, navigation.PipsPagerOptions{Index: &index, Count: 5})`,
		Params: []ComponentParam{
			{"Index", "*int", "nil", "First slide in view, 0-based; nil keeps it internally."},
			{"Count", "int", "0", "Number of slides."},
			{"Slide", "func(index int)", "nil", "Builds one slide; only the slides in view are built."},
			{"PerView", "int", "1", "Slides shown side by side."},
			{"Loop", "bool", "false", "Continues from the last slide to the first."},
			{"Label", "string", "\"\"", "Accessible name of the slides."},
			{"OnChange", "func(index int)", "nil", "Runs after the user steps."},
			{"MaxVisible", "int", "7", "PipsPagerOptions: pips built around the current page."},
		},
		Tokens: []string{"Foreground", "MutedForeground", "Input", "Accent", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "statusbar", Title: "StatusBar", Group: "Patterns", Order: 200,
		Description: "A window's bottom strip: a transient message and permanent items.",
		Examples: []ComponentExample{
			{Title: "Editor status", Demo: workspaceNavDemo(statusBarBody)},
		},
		Code: `var message string

navigation.StatusBar(c, navigation.StatusBarOptions{
    Label:   "Status",
    Message: message,
    Tone:    navigation.StatusBarWarning,
    Items: []navigation.StatusBarItemOptions{
        {Text: "Ln 12, Col 4", Priority: 1},
        {Separator: true},
        {Icon: lucide.GitBranch(), Text: "main", Tooltip: "Switch branch", Priority: 3,
            OnClick: func() { message = "Switched branch" }},
        {Label: "Syncing", Progress: &navigation.StatusBarProgress{Indeterminate: true}},
    },
})`,
		Params: []ComponentParam{
			{"Message", "string", "\"\"", "Transient text at the leading edge; changes are announced. It keeps a minimum width."},
			{"Label", "string", "\"\"", "Accessible name of the bar."},
			{"Tone", "navigation.StatusBarTone", "StatusBarNeutral", "StatusBarWarning or StatusBarError colours the whole bar."},
			{"Items", "[]navigation.StatusBarItemOptions", "nil", "The trailing items."},
			{"Text", "string", "\"\"", "StatusBarItemOptions: item text."},
			{"Icon", "*ui.SVG", "nil", "StatusBarItemOptions: leading icon."},
			{"Label", "string", "\"\"", "StatusBarItemOptions: names an icon alone, and the button."},
			{"Tooltip", "string", "\"\"", "StatusBarItemOptions: shown on pointer rest or focus."},
			{"Progress", "*navigation.StatusBarProgress", "nil", "StatusBarItemOptions: a track, Value 0 to 1 or Indeterminate."},
			{"Priority", "int", "0", "StatusBarItemOptions: in a narrow bar the lowest leave first."},
			{"Separator", "bool", "false", "StatusBarItemOptions: a short rule; it stays only between shown items."},
			{"OnClick", "func()", "nil", "StatusBarItemOptions: makes the item a button."},
		},
		Tokens: []string{"MutedBackground", "MutedForeground", "Foreground", "Border", "Separator", "Accent", "Ring"},
	})
}

// workspaceNavFrame is the bordered window stand-in the demos sit in.
func workspaceNavFrame(c *ui.Context, body func()) ui.Element {
	t := theme.Colors(c)
	return ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 560)).
		Radius(theme.RadiusScale().XL).Border(theme.Hairline, t.Border).Clip().
		Background(t.Background).Children(body)
}

func documentTabsBody(c *ui.Context, st *workspaceNavState) {
	workspaceNavFrame(c, func() {
		ui.Column(c).FillWidth().Padding(theme.Space8).Children(func() {
			navigation.DocumentTabs(c, navigation.DocumentTabsOptions{
				Label: "Open documents", Tabs: st.docs, Active: st.activeDoc,
				OnSelect: func(key string) { st.activeDoc = key },
				OnClose: func(key string) {
					i := slices.IndexFunc(st.docs, func(d navigation.DocumentTab) bool { return d.Key == key })
					st.docs = slices.Delete(st.docs, i, i+1)
					if key == st.activeDoc && len(st.docs) > 0 {
						st.activeDoc = st.docs[min(i, len(st.docs)-1)].Key
					}
				},
				OnCloseMany: func(keys []string) {
					st.docs = slices.DeleteFunc(st.docs, func(d navigation.DocumentTab) bool { return slices.Contains(keys, d.Key) })
					if !slices.ContainsFunc(st.docs, func(d navigation.DocumentTab) bool { return d.Key == st.activeDoc }) && len(st.docs) > 0 {
						st.activeDoc = st.docs[0].Key
					}
				},
				OnMove: func(from, to int) {
					doc := st.docs[from]
					st.docs = slices.Insert(slices.Delete(st.docs, from, from+1), to, doc)
				},
				OnAdd: func() {
					st.created++
					key := "untitled-" + strconv.Itoa(st.created)
					st.docs = append(st.docs, navigation.DocumentTab{Key: key, Title: "Untitled " + strconv.Itoa(st.created)})
					st.activeDoc = key
				},
			})
		})
		title := "No document open"
		if i := slices.IndexFunc(st.docs, func(d navigation.DocumentTab) bool { return d.Key == st.activeDoc }); i >= 0 {
			title = st.docs[i].Title
		}
		ui.Box(c).FillWidth().Height(theme.Scale(c, 96)).Center().Children(func() {
			mygoui.Typography(c, title, mygoui.TypographyOptions{Type: theme.TypeBodySM, Muted: true})
		})
	})
}

func documentTabsOverflowBody(c *ui.Context, st *workspaceNavState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 420)).Gap(theme.Space12).Children(func() {
		navigation.DocumentTabs(c, navigation.DocumentTabsOptions{
			Label: "Chapters", Tabs: st.manyDocs, Active: st.activeManyDoc,
			OnSelect: func(key string) { st.activeManyDoc = key },
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Open last chapter", Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
			OnClick: func() { st.activeManyDoc = st.manyDocs[len(st.manyDocs)-1].Key },
		})
	})
}

func documentTabsShrinkBody(c *ui.Context, st *workspaceNavState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 560)).Children(func() {
		navigation.DocumentTabs(c, navigation.DocumentTabsOptions{
			Label: "Open pages", Tabs: st.manyDocs[:6], Active: st.activeManyDoc, Sizing: navigation.DocumentTabsShrink,
			OnSelect: func(key string) { st.activeManyDoc = key },
		})
	})
}

// workspaceNavSettings maps each settings page to the pages under it.
var workspaceNavSettings = map[string][]string{
	"Settings": {"Network", "Display", "Sound"},
	"Network":  {"Wi-Fi", "Ethernet", "VPN"},
	"Display":  {"Brightness", "Night light"},
}

func navigationStackBody(c *ui.Context, st *workspaceNavState) {
	pages := make([]navigation.NavigationPage, len(st.stack))
	for i, title := range st.stack {
		pages[i] = navigation.NavigationPage{Key: title, Title: title, Build: func() {
			navigationStackPageBody(c, st, title)
		}}
	}
	workspaceNavFrame(c, func() {
		navigation.NavigationStack(c, navigation.NavigationStackOptions{
			Label: "Settings", Pages: pages,
			OnPop: func() { st.stack = st.stack[:len(st.stack)-1] },
			Actions: func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Reset", Icon: lucide.RotateCcw(), IconOnly: true,
					Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
					OnClick: func() { st.stack = st.stack[:1] },
				})
			},
		}).Height(theme.Scale(c, 220)).PaddingX(theme.Space12)
	})
}

func navigationStackPageBody(c *ui.Context, st *workspaceNavState, title string) {
	children := workspaceNavSettings[title]
	if len(children) == 0 {
		ui.Box(c).Fill().Center().Children(func() {
			mygoui.Typography(c, title+" has no further settings.", mygoui.TypographyOptions{
				Type: theme.TypeBodySM, Muted: true,
			})
		})
		return
	}
	ui.Column(c).FillWidth().Gap(theme.Space2).Children(func() {
		for _, child := range children {
			mygoui.Button(c.Key(child), mygoui.ButtonOptions{
				Label: child, EndIcon: lucide.ChevronRight(), Variant: mygoui.ButtonGhost,
				FullWidth: true, Align: mygoui.ButtonAlignStart,
				OnClick: func() { st.stack = append(st.stack, child) },
			})
		}
	})
}

func carouselSlideBody(c *ui.Context, i int) {
	t := theme.Colors(c)
	ui.Box(c).FillWidth().Height(theme.Scale(c, 120)).Center().
		Radius(theme.RadiusScale().LG).Background(t.MutedBackground).Children(func() {
		mygoui.Typography(c, fmt.Sprintf("Slide %d", i+1), mygoui.TypographyOptions{Type: theme.TypeBodySM})
	})
}

func carouselLoopBody(c *ui.Context, st *workspaceNavState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 420)).Children(func() {
		navigation.Carousel(c, navigation.CarouselOptions{
			Label: "Highlights", Index: &st.slide, Count: 5, Loop: true,
			Slide: func(i int) { carouselSlideBody(c, i) },
		})
	})
}

func carouselManyBody(c *ui.Context, st *workspaceNavState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 420)).Children(func() {
		navigation.Carousel(c, navigation.CarouselOptions{
			Label: "Archive", Index: &st.many, Count: 24,
			Slide: func(i int) { carouselSlideBody(c, i) },
		})
	})
}

func carouselPerViewBody(c *ui.Context, st *workspaceNavState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 560)).Children(func() {
		navigation.Carousel(c, navigation.CarouselOptions{
			Label: "Gallery", Index: &st.cards, Count: 8, PerView: 3,
			Slide: func(i int) { carouselSlideBody(c, i) },
		})
	})
}

func pipsPagerBody(c *ui.Context, st *workspaceNavState) {
	ui.Column(c).AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Typography(c, fmt.Sprintf("Page %d of 24", st.pip+1), mygoui.TypographyOptions{
			Type: theme.TypeBodySM, Muted: true,
		})
		navigation.PipsPager(c, navigation.PipsPagerOptions{Label: "Onboarding", Index: &st.pip, Count: 24})
	})
}

func statusBarBody(c *ui.Context, st *workspaceNavState) {
	branches := []string{"main", "release", "feature/tabs"}
	workspaceNavFrame(c, func() {
		ui.Row(c).FillWidth().Height(theme.Scale(c, 96)).Center().Gap(theme.Space8).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Save", Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
				OnClick: func() {
					st.saves++
					st.message = fmt.Sprintf("Saved %d files", st.saves)
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Clear message", Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				OnClick: func() { st.message = "" },
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Sync", Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				OnClick: func() { st.syncing = !st.syncing },
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Tone", Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				OnClick: func() { st.tone = (st.tone + 1) % 3 },
			})
		})
		items := []navigation.StatusBarItemOptions{
			{Text: "Ln 12, Col 4", Priority: 1},
			{Text: "UTF-8", Priority: 1},
			{Separator: true},
			{
				Icon: lucide.GitBranch(), Text: branches[st.branch], Tooltip: "Switch branch", Priority: 3,
				OnClick: func() {
					st.branch = (st.branch + 1) % len(branches)
					st.message = "Switched to " + branches[st.branch]
				},
			},
			{
				Icon: lucide.Bell(), Text: strconv.Itoa(st.unreadCount), Label: "Notifications", Priority: 2,
				Tooltip: "Mark notifications read",
				OnClick: func() {
					st.unreadCount = 0
					st.message = "Notifications cleared"
				},
			},
		}
		if st.syncing {
			items = append(items, navigation.StatusBarItemOptions{
				Label: "Syncing", Progress: &navigation.StatusBarProgress{Indeterminate: true}, Priority: 4,
			})
		}
		navigation.StatusBar(c, navigation.StatusBarOptions{
			Label: "Status", Message: st.message, Tone: st.tone, Items: items,
		})
	})
}
