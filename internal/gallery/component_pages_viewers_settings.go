package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	uicontent "github.com/Cyenoch/mygoui/content"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// viewersState is the caller-owned state behind the viewer and settings
// demos.
type viewersState struct {
	ready bool

	json        uicontent.JSONViewerState
	jsonSource  string
	jsonRefresh int

	logLines []string
	log      uicontent.LogViewerState

	loadingGroup           bool
	wifi, bluetooth, proxy bool
	wifiSwitch, btSwitch   ui.Handle
	proxySwitch            ui.Handle
	advanced               bool
	opened                 string
}

func viewersStateOf(s *State) *viewersState {
	st := StateOf[viewersState](s, "viewers-settings")
	if !st.ready {
		st.ready = true
		st.wifi = true
		st.jsonSource = viewersJSONSource(0)
		for range 300 {
			st.logLines = append(st.logLines, viewersLogLine(len(st.logLines)))
		}
	}
	return st
}

func viewersDemo(body func(c *ui.Context, st *viewersState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, viewersStateOf(s)) }
}

const viewersDiffOld = `package greet

import "fmt"

// Greet returns a greeting for name.
func Greet(name string) string {
	if name == "" {
		return "Hello!"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Farewell returns a goodbye for name.
func Farewell(name string) string {
	return fmt.Sprintf("Goodbye, %s.", name)
}

func main() {
	fmt.Println(Greet("Ada"))
	fmt.Println(Farewell("Ada"))
}
`

const viewersDiffNew = `package greet

import (
	"fmt"
	"strings"
)

// Greet returns a greeting for name.
func Greet(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Hello!"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Farewell returns a goodbye for name.
func Farewell(name string) string {
	return fmt.Sprintf("Goodbye, %s.", name)
}

func main() {
	fmt.Println(Greet(" Ada "))
	fmt.Println(Farewell("Ada"))
}
`

var viewersLogMessages = [...]string{
	"DEBUG cache warm in 12ms",
	"INFO  GET /api/projects 200 18ms",
	"INFO  GET /api/projects/42 200 9ms",
	"WARN  slow query: 412ms SELECT * FROM builds WHERE project_id = 42",
	"INFO  POST /api/builds 201 31ms",
	"ERROR upstream timeout after 5s: artifacts.internal:443",
}

func viewersLogLine(i int) string {
	return fmt.Sprintf("12:%02d:%02d %s", i/60%60, i%60, viewersLogMessages[i%len(viewersLogMessages)])
}

func viewersLogLevel(line string) uicontent.LogLevel {
	// The level follows the 9-byte timestamp.
	switch line[9:14] {
	case "DEBUG":
		return uicontent.LogDebug
	case "INFO ":
		return uicontent.LogInfo
	case "WARN ":
		return uicontent.LogWarning
	case "ERROR":
		return uicontent.LogError
	}
	return uicontent.LogPlain
}

func viewersJSONSource(refresh int) string {
	return fmt.Sprintf(`{
  "id": 42,
  "name": "mygoui",
  "private": false,
  "license": null,
  "stars": %d,
  "description": "A generic native component library with a shadcn-inspired look: typed options, caller-owned state, keyboard and screen reader support, and no global setup. Every control is built from MyGo's own widgets.",
  "owner": {"login": "cyenoch", "site": "https://coss.com/ui/docs/components/native/json-viewer?tab=overview&ref=gallery"},
  "empty": {},
  "topics": ["native", "go", "components"],
  "releases": [
    {"tag": "v0.1.0", "assets": 3},
    {"tag": "v0.0.9", "assets": 2}
  ]
}`, 128+refresh)
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "diffviewer", Title: "DiffViewer", Group: "Viewers", Order: 220,
		Description: "A read-only diff of two texts, unified or side by side.",
		Examples: []ComponentExample{
			{Title: "Unified", Description: "Unchanged lines collapse into separators that expand; step through the changes with the arrows or F7.",
				Demo: viewersDemo(diffViewerUnifiedBody)},
			{Title: "Side by side", Description: "A fixed height scrolls both sides together.",
				Demo: viewersDemo(diffViewerSplitBody)},
		},
		Code: `content.DiffViewer(c, content.DiffViewerOptions{
    Old:    sent,
    New:    draft,
    Title:  "greet.go",
    Layout: content.DiffSplit,
    Height: 320,
})`,
		Params: []ComponentParam{
			{"Old", "string", "\"\"", "Text before the change."},
			{"New", "string", "\"\"", "Text after the change."},
			{"Title", "string", "\"\"", "Header naming the file, with line counts, a layout toggle, change stepping (F7, Shift+F7) and copy of New."},
			{"Layout", "content.DiffLayout", "DiffUnified", "Unified or side by side."},
			{"ContextLines", "int", "3", "Unchanged lines kept around a change; negative keeps none."},
			{"ExpandUnchanged", "bool", "false", "Shows every unchanged line."},
			{"NoLineNumbers", "bool", "false", "Leaves out the gutter numbers."},
			{"Wrap", "bool", "false", "Wraps long lines instead of scrolling sideways."},
			{"Height", "float32", "fit lines", "Scrolls within this many DIPs; change stepping needs it."},
			{"OnError", "func(error)", "nil", "Receives a clipboard failure of the copy."},
		},
		Tokens: []string{"Code", "CodeForeground", "Border", "Muted", "Success", "Danger", "Info"},
	})

	RegisterComponent(ComponentPage{
		ID: "jsonviewer", Title: "JSONViewer", Group: "Viewers", Order: 230,
		Description: "A JSON document as a tree that opens, closes and copies.",
		Examples: []ComponentExample{
			{Title: "Document", Description: "Search with Cmd+F; choose a row to read a long string in full and copy it with Cmd+C. Open nodes survive a refresh.",
				Demo: viewersDemo(jsonViewerDocumentBody)},
			{Title: "Invalid", Description: "The error is returned and shown in place of the tree.",
				Demo: viewersDemo(jsonViewerInvalidBody)},
		},
		Code: `if _, err := content.JSONViewer(c, content.JSONViewerOptions{
    Source: draft,
    Title:  "response.json",
}); err != nil {
    last = err.Error()
}`,
		Params: []ComponentParam{
			{"Source", "string", "required", "The document; invalid JSON returns an error."},
			{"Title", "string", "JSON", "Header text and accessible name."},
			{"Height", "float32", "rows, up to 320", "Tree height in DIPs."},
			{"ExpandDepth", "int", "2", "Levels open on the first build; negative opens only the root."},
			{"Toolbar", "bool", "false", "Search field with match count, Enter and Shift+Enter stepping, Cmd+F to focus."},
			{"OnError", "func(error)", "nil", "Receives a clipboard failure of a copy."},
			{"State", "*content.JSONViewerState", "per window", "Caller-owned open nodes, Query, scroll; ExpandAll and CollapseAll."},
		},
		Tokens: []string{"Code", "CodeForeground", "Border", "Muted", "SuccessForeground", "WarningForeground", "Primary", "DestructiveForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "logviewer", Title: "LogViewer", Group: "Viewers", Order: 240,
		Description: "Virtualized monospace log lines that follow the newest, with search and level filters.",
		Examples: []ComponentExample{
			{Title: "Follow", Description: "Filter by text or level; scroll up to hold your place and the button returns to the end.",
				Demo: viewersDemo(logViewerFollowBody)},
			{Title: "Plain", Description: "Without numbers or levels, as high as its lines.",
				Demo: viewersDemo(logViewerPlainBody)},
		},
		Code: `var lines []string // caller-owned, appended as the process writes
content.LogViewer(c, content.LogViewerOptions{
    Lines:       lines,
    LineNumbers: true,
    Follow:      true,
    Toolbar:     true,
    Height:      240,
    Level: func(i int) content.LogLevel {
        if len(lines[i]) > 0 && lines[i][0] == 'E' {
            return content.LogError
        }
        return content.LogPlain
    },
})`,
		Params: []ComponentParam{
			{"Lines", "[]string", "nil", "Caller-owned lines."},
			{"Line", "func(i int) string", "nil", "Supplies line i instead of Lines, with Count."},
			{"Count", "int", "0", "Number of lines Line supplies."},
			{"Level", "func(i int) content.LogLevel", "nil", "Tints line i: debug, info, warning or error."},
			{"LineNumbers", "bool", "false", "Shows a number gutter."},
			{"Follow", "bool", "false", "Keeps the newest line in view until the user scrolls away."},
			{"Label", "string", "\"\"", "Accessible name."},
			{"Height", "float32", "lines, up to 320", "Viewer height in DIPs."},
			{"Toolbar", "bool", "false", "Search field, level toggles, match count and wrap toggle above the lines; Cmd/Ctrl+F focuses the field, Escape clears it."},
			{"OnError", "func(error)", "nil", "Receives a clipboard failure of Copy line or Copy all."},
			{"State", "*content.LogViewerState", "per window", "Caller-owned scroll state and filter."},
			{"State.Query", "string", "\"\"", "Lists only lines holding it, whatever the case, and marks it; numbers stay the log's."},
			{"State.Levels", "content.LogLevels", "all", "Lists only lines of these levels; build with LevelSet."},
			{"State.Revision", "int", "0", "Raise it when lines change in place so the listing is scanned again."},
			{"State.NoWrap", "bool", "false", "Keeps each line on one row, cut where it ends, instead of wrapping."},
		},
		Tokens: []string{"Code", "CodeForeground", "Border", "Muted", "WarningForeground", "DestructiveForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "settingsgroup", Title: "SettingsGroup", Group: "Patterns", Order: 250,
		Description: "A titled boxed list of settings rows with trailing controls.",
		Examples: []ComponentExample{
			{Title: "Controls", Description: "The row title names its trailing control.",
				Demo: viewersDemo(settingsControlsBody)},
			{Title: "Navigation and expander", Description: "A pressable row activates whole; an expander reveals nested rows; a destructive row is inked red.",
				Demo: viewersDemo(settingsNavigationBody)},
			{Title: "Loading", Description: "Without rows they are skeletons; with rows the card dims under a line.",
				Demo: viewersDemo(settingsLoadingBody)},
		},
		Code: `var wifiSwitch ui.Handle // kept in application state
mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{
    Title: "Network",
    Rows: []mygoui.SettingsRow{
        {
            Icon: lucide.Wifi(), Title: "Wi-Fi", Subtitle: "Join known networks",
            Control: &wifiSwitch,
            Trailing: func() {
                form.Switch(c, form.SwitchOptions{On: &wifi, Handle: &wifiSwitch})
            },
        },
        {Icon: lucide.BookOpen(), Title: "Documentation", OnClick: openDocs},
    },
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "Group heading and accessible name."},
			{"Description", "string", "\"\"", "Muted text under the heading."},
			{"Action", "func()", "nil", "Control at the end of the header."},
			{"Rows", "[]mygoui.SettingsRow", "nil", "Rows of the card, separated by hairlines."},
			{"Loading", "bool", "false", "Skeleton rows while Rows is empty; otherwise dims the rows and runs a line along the card."},
			{"LoadingRows", "int", "3", "Skeleton rows shown while loading with no rows."},
			{"SettingsRow.Icon", "*ui.SVG", "nil", "Leading icon."},
			{"SettingsRow.Title", "string", "required", "Row title; names the trailing control."},
			{"SettingsRow.Subtitle", "string", "\"\"", "Muted text under the title."},
			{"SettingsRow.Trailing", "func()", "nil", "Builds the control or value at the row's end."},
			{"SettingsRow.Control", "*ui.Handle", "nil", "Handle bound to the trailing control, which the title then names."},
			{"SettingsRow.OnClick", "func()", "nil", "Makes the whole row a button with a chevron."},
			{"SettingsRow.Rows", "[]mygoui.SettingsRow", "nil", "Nested rows; makes the row an expander."},
			{"SettingsRow.Expanded", "*bool", "per window", "Caller-owned expander state."},
			{"SettingsRow.Disabled", "bool", "false", "Stops activation and dims the row."},
			{"SettingsRow.Destructive", "bool", "false", "Inks the icon and title with the danger color."},
		},
		Tokens: []string{"Card", "CardForeground", "Border", "Separator", "Muted", "SurfaceHover", "MutedBackground", "Ring"},
	})
}

func diffViewerUnifiedBody(c *ui.Context, st *viewersState) {
	uicontent.DiffViewer(c, uicontent.DiffViewerOptions{
		Old: viewersDiffOld, New: viewersDiffNew, Title: "greet.go", Height: 320,
	})
}

func diffViewerSplitBody(c *ui.Context, st *viewersState) {
	uicontent.DiffViewer(c, uicontent.DiffViewerOptions{
		Old: viewersDiffOld, New: viewersDiffNew, Title: "greet.go",
		Layout: uicontent.DiffSplit, ExpandUnchanged: true, Height: 280,
	})
}

func jsonViewerDocumentBody(c *ui.Context, st *viewersState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		if _, err := uicontent.JSONViewer(c, uicontent.JSONViewerOptions{
			Source: st.jsonSource, Title: "repository.json", State: &st.json,
			Toolbar: true, Height: 300,
		}); err != nil {
			form.FieldError(c, err.Error())
		}
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Refresh", Icon: lucide.RefreshCw(), Variant: mygoui.ButtonOutline,
			Size: theme.ButtonSmall, OnClick: func() {
				st.jsonRefresh++
				st.jsonSource = viewersJSONSource(st.jsonRefresh)
			},
		})
	})
}

func jsonViewerInvalidBody(c *ui.Context, st *viewersState) {
	// The viewer shows the returned error itself.
	_, _ = uicontent.JSONViewer(c, uicontent.JSONViewerOptions{
		Source: `{"name": "mygoui", "topics": ["native", "go",]}`, Title: "broken.json",
	})
}

func logViewerFollowBody(c *ui.Context, st *viewersState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		uicontent.LogViewer(c, uicontent.LogViewerOptions{
			Lines: st.logLines, LineNumbers: true, Follow: true, Toolbar: true, Height: 260,
			Label: "Server log", State: &st.log,
			Level: func(i int) uicontent.LogLevel { return viewersLogLevel(st.logLines[i]) },
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Add 5 lines", Icon: lucide.Plus(), Variant: mygoui.ButtonOutline,
			Size: theme.ButtonSmall,
			OnClick: func() {
				for range 5 {
					st.logLines = append(st.logLines, viewersLogLine(len(st.logLines)))
				}
			},
		})
	})
}

func logViewerPlainBody(c *ui.Context, st *viewersState) {
	uicontent.LogViewer(c, uicontent.LogViewerOptions{
		Lines: viewersLogMessages[:], Label: "Recent messages",
	})
}

func settingsLoadingBody(c *ui.Context, st *viewersState) {
	formsFieldStack(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.loadingGroup, Label: "Loading"})
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{
			Title: "Connections", Loading: st.loadingGroup,
			Rows: []mygoui.SettingsRow{
				{Icon: lucide.Wifi(), Title: "Wi-Fi", Subtitle: "Join known networks automatically"},
				{Icon: lucide.Bluetooth(), Title: "Bluetooth"},
			},
		})
	})
}

func settingsControlsBody(c *ui.Context, st *viewersState) {
	formsFieldStack(c, func() {
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{
			Title: "Connections", Description: "Radios this device may use.",
			Action: func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Reset", Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
					OnClick: func() { st.wifi, st.bluetooth = true, false },
				})
			},
			Rows: []mygoui.SettingsRow{
				{
					Icon: lucide.Wifi(), Title: "Wi-Fi", Subtitle: "Join known networks automatically",
					Control: &st.wifiSwitch,
					Trailing: func() {
						form.Switch(c, form.SwitchOptions{On: &st.wifi, Handle: &st.wifiSwitch})
					},
				},
				{
					Icon: lucide.Bluetooth(), Title: "Bluetooth", Control: &st.btSwitch,
					Trailing: func() {
						form.Switch(c, form.SwitchOptions{On: &st.bluetooth, Handle: &st.btSwitch})
					},
				},
				{
					Icon: lucide.Plane(), Title: "Airplane mode", Subtitle: "Managed by your organization",
					Disabled: true,
					Trailing: func() { form.Switch(c, form.SwitchOptions{Disabled: true}) },
				},
			},
		})
	})
}

func settingsNavigationBody(c *ui.Context, st *viewersState) {
	formsFieldStack(c, func() {
		mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{
			Title: "System",
			Rows: []mygoui.SettingsRow{
				{
					Icon: lucide.User(), Title: "Account", Subtitle: "ada@example.com",
					OnClick: func() { st.opened = "Account" },
				},
				{
					Icon: lucide.Languages(), Title: "Language",
					Trailing: func() { formsMutedText(c, "English") },
					OnClick:  func() { st.opened = "Language" },
				},
				{
					Icon: lucide.Network(), Title: "Proxy", Subtitle: "Route traffic through a server",
					Expanded: &st.advanced, Control: &st.proxySwitch,
					Trailing: func() {
						form.Switch(c, form.SwitchOptions{On: &st.proxy, Handle: &st.proxySwitch})
					},
					Rows: []mygoui.SettingsRow{
						{Title: "Server", Trailing: func() { formsMutedText(c, "proxy.internal") }},
						{Title: "Port", Trailing: func() { formsMutedText(c, "8080") }},
					},
				},
				{
					Icon: lucide.Trash2(), Title: "Erase all data", Subtitle: "This cannot be undone",
					Destructive: true, OnClick: func() { st.opened = "Erase" },
				},
			},
		})
		if st.opened != "" {
			formsMutedText(c, "Opened "+st.opened)
		}
	})
}
