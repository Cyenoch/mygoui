package gallery

import (
	"fmt"
	"math/rand/v2"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// dataDisplayState is the caller-owned state behind the data display demos.
type dataDisplayState struct {
	ready bool

	deploying     bool
	query         string
	ignoreCase    bool
	notesExpanded bool
	saving        bool
	loading       bool
	invoiceName   string
	updateHidden  bool
	restarts      int
	truncate      bool
	earlierOpen   bool
	deployOpened  int

	revenue, latency, errs, churn int // churn in tenths of a percent
	count                         int
}

func dataDisplayStateOf(s *State) *dataDisplayState {
	st := StateOf[dataDisplayState](s, "data-display")
	if !st.ready {
		st.ready = true
		st.deploying = true
		st.query, st.ignoreCase = "go", true
		st.invoiceName = "Invoice 10482"
		st.revenue, st.latency, st.errs, st.churn, st.count = 12480, 248, 14, 12, 12480
	}
	return st
}

func dataDisplayDemo(body func(c *ui.Context, st *dataDisplayState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, dataDisplayStateOf(s)) }
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "descriptions", Title: "Descriptions", Group: "Display", Order: 220,
		Description: "Read-only label and value pairs in a grid.",
		Examples: []ComponentExample{
			{Title: "Plain", Description: "An item can span columns; the last of a row takes what is left.",
				Demo: dataDisplayDemo(descriptionsPlainBody)},
			{Title: "Bordered", Demo: dataDisplayDemo(descriptionsBorderedBody)},
			{Title: "Narrow", Description: "Columns drop and a table stacks its labels as the width shrinks; Truncate cuts values at one line.",
				Demo: dataDisplayDemo(descriptionsNarrowBody)},
			{Title: "Loading", Description: "Labels stay; each value waits as a bar.",
				Demo: dataDisplayDemo(descriptionsLoadingBody)},
		},
		Code: `mygoui.Descriptions(c, mygoui.DescriptionsOptions{
    Title:   "Order details",
    Columns: 3,
    Items: []mygoui.DescriptionsItem{
        {Label: "Customer", Value: "Ada Lovelace"},
        {Label: "Order", Value: "#10482"},
        {Label: "Address", Value: "12 St James's Square, London", Span: 2},
    },
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "Heading and accessible name of the group."},
			{"Items", "[]mygoui.DescriptionsItem", "nil", "Label, Value or Content, and Span in columns; an empty value shows a dash."},
			{"Columns", "int", "3", "Items in a row, fewer where the width leaves a column under MinColumnWidth."},
			{"MinColumnWidth", "float32", "200 scaled", "Narrowest a column gets before the grid drops one."},
			{"Truncate", "bool", "false", "Cuts each value at one line with a tooltip of the whole; values wrap otherwise."},
			{"Bordered", "bool", "false", "Ruled table with aligned label columns."},
			{"Loading", "bool", "false", "Keeps the labels and shows a skeleton bar for each value."},
		},
		Tokens: []string{"Foreground", "MutedForeground", "MutedBackground", "Border", "Background"},
	})

	RegisterComponent(ComponentPage{
		ID: "statistic", Title: "Statistic", Group: "Display", Order: 230,
		Description: "A large metric with its label and an optional change.",
		Examples: []ComponentExample{
			{Title: "Trends", Description: "DeltaStatus recolors a rise that is bad news. A changed value rolls.",
				Demo: dataDisplayDemo(statisticTrendsBody)},
			{Title: "Content", Description: "A value wider than its room is cut with an ellipsis and a tooltip.",
				Demo: dataDisplayDemo(statisticContentBody)},
			{Title: "Loading", Description: "The block holds the value's height.",
				Demo: dataDisplayDemo(statisticLoadingBody)},
		},
		Code: `mygoui.Statistic(c, mygoui.StatisticOptions{
    Label:  "Revenue",
    Value:  "12,480",
    Prefix: "$",
    Delta:  "12.5% vs last week",
    Trend:  mygoui.TrendUp,
})`,
		Params: []ComponentParam{
			{"Label", "string", "\"\"", "Name of the metric."},
			{"Value", "string", "\"\"", "Caller-formatted value."},
			{"Prefix", "string", "\"\"", "Muted text before the value."},
			{"Suffix", "string", "\"\"", "Muted text after the value."},
			{"Delta", "string", "\"\"", "Caller-formatted change under the value."},
			{"Trend", "mygoui.StatisticTrend", "TrendFlat", "Direction icon of the delta: flat, up or down."},
			{"DeltaStatus", "mygoui.StatusColor", "trend color", "Replaces success for up and danger for down."},
			{"Content", "func()", "nil", "Builds caller content under the value and its change."},
			{"Loading", "bool", "false", "Shows a skeleton block of the value's height and hides the delta."},
		},
		Tokens: []string{"Foreground", "MutedForeground", "SuccessSoftForeground", "DangerSoftForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "numberflip", Title: "NumberFlip", Group: "Display", Order: 235,
		Description: "A value whose digits roll to the new one.",
		Examples: []ComponentExample{
			{Title: "Counter", Description: "Digits roll in the direction of the change; reduced motion crossfades.",
				Demo: dataDisplayDemo(numberFlipCounterBody)},
		},
		Code: `mygoui.NumberFlip(c, mygoui.NumberFlipOptions{
    Value: "12,480",
    Style: theme.H2(),
})`,
		Params: []ComponentParam{
			{"Value", "string", "\"\"", "Caller-formatted value; a change rolls the digits that differ."},
			{"Style", "theme.TextStyle", "theme.H2()", "Text role of the value."},
			{"Color", "ui.Color", "Foreground", "Ink of the value."},
			{"Label", "string", "Value", "Accessible name."},
		},
		Tokens: []string{"Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "timeline", Title: "Timeline", Group: "Display", Order: 240,
		Description: "Events down a line of status markers.",
		Examples: []ComponentExample{
			{Title: "Events", Description: "Pending turns the last marker into a spinner; Recent folds the earlier items behind a control; OnClick makes an item's text a button.",
				Demo: dataDisplayDemo(timelineEventsBody)},
		},
		Code: `mygoui.Timeline(c, mygoui.TimelineOptions{
    Label:   "Release",
    Pending: true,
    Items: []mygoui.TimelineItem{
        {Title: "Build passed", Time: "09:12", Status: mygoui.ColorSuccess},
        {Title: "Deployed to staging", Time: "09:20", Icon: lucide.Rocket()},
        {Title: "Rolling out"},
    },
})`,
		Params: []ComponentParam{
			{"Items", "[]mygoui.TimelineItem", "nil", "Title, Description, Time, Status and Icon of each event."},
			{"Label", "string", "\"\"", "Accessible name of the list."},
			{"Pending", "bool", "false", "Marks the last item as in progress."},
			{"Recent", "int", "0", "Newest items shown; the earlier ones open from a control. 0 shows all."},
			{"Expanded", "*bool", "element state", "Caller-owned flag of the earlier items."},
			{"EarlierLabel", "string", "locale ShowMore", "Label of the control; format the count yourself."},
			{"TimelineItem.Content", "func()", "nil", "Caller content under the description."},
			{"TimelineItem.OnClick", "func()", "nil", "Makes the item's text a button."},
		},
		Tokens: []string{"Foreground", "MutedForeground", "Border", "Success", "Danger", "Info", "Warning"},
	})

	RegisterComponent(ComponentPage{
		ID: "statusdot", Title: "StatusDot", Group: "Display", Order: 250,
		Description: "A small colored status indicator with an optional label.",
		Examples: []ComponentExample{
			{Title: "Statuses", Description: "A live status pulses; a connecting one breathes; the rest are still.",
				Demo: dataDisplayDemo(statusDotStatusesBody)},
			{Title: "Motion", Description: "Both are still under reduced motion.",
				Demo: dataDisplayDemo(statusDotPulseBody)},
		},
		Code: `mygoui.StatusDot(c, mygoui.StatusDotOptions{
    Status: mygoui.ColorSuccess,
    Label:  "Online",
    Pulse:  true,
})`,
		Params: []ComponentParam{
			{"Status", "mygoui.StatusColor", "ColorDefault", "Dot color; the default is muted."},
			{"Label", "string", "\"\"", "Text after the dot."},
			{"Pulse", "bool", "false", "Radiates a ring from the dot, for something live."},
			{"Breathe", "bool", "false", "Dims and restores the dot, for something in progress."},
		},
		Tokens: []string{"Foreground", "MutedForeground", "Success", "Warning", "Danger", "Info"},
	})

	RegisterComponent(ComponentPage{
		ID: "highlight", Title: "Highlight", Group: "Display", Order: 260,
		Description: "Text with the matches of search queries marked.",
		Examples: []ComponentExample{
			{Title: "Search", Demo: dataDisplayDemo(highlightSearchBody)},
		},
		Code: `mygoui.Highlight(c, mygoui.HighlightOptions{
    Text:       bio,
    Queries:    []string{query},
    IgnoreCase: true,
})`,
		Params: []ComponentParam{
			{"Text", "string", "\"\"", "The paragraph."},
			{"Queries", "[]string", "nil", "Substrings to mark; trimmed, blank ones ignored."},
			{"IgnoreCase", "bool", "false", "Matches letters whatever their case."},
		},
		Tokens: []string{"Foreground", "WarningSoft", "WarningSoftForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "spoiler", Title: "Spoiler", Group: "Display", Order: 270,
		Description: "Content cut at a height, fading out, with a control that shows the rest.",
		Examples: []ComponentExample{
			{Title: "Long text", Demo: dataDisplayDemo(spoilerLongBody)},
			{Title: "Fits", Description: "No control when nothing is cut.",
				Demo: dataDisplayDemo(spoilerFitsBody)},
		},
		Code: `mygoui.Spoiler(c, mygoui.SpoilerOptions{MaxHeight: 64}, func() {
    mygoui.Typography(c, bio, mygoui.TypographyOptions{Type: theme.TypeBodySM})
})`,
		Params: []ComponentParam{
			{"Expanded", "*bool", "element state", "Caller-owned flag the control toggles."},
			{"MaxHeight", "float32", "96 scaled", "Height the collapsed content is cut at."},
			{"ShowLabel", "string", "locale ShowMore", "Label of the control while collapsed."},
			{"HideLabel", "string", "locale ShowLess", "Label of the control while expanded."},
		},
		Tokens: []string{"Link", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "loadingoverlay", Title: "LoadingOverlay", Group: "Display", Order: 280,
		Description: "A light scrim and spinner over content that is busy.",
		Examples: []ComponentExample{
			{Title: "Saving", Description: "Covered content takes no pointer or keyboard input.",
				Demo: dataDisplayDemo(loadingOverlaySavingBody)},
		},
		Code: `mygoui.LoadingOverlay(c, mygoui.LoadingOverlayOptions{
    Visible: true,
    Label:   "Saving invoice",
}, func() {
    form.Input(c, &name, form.InputOptions{Label: "Name"})
    mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: save})
})`,
		Params: []ComponentParam{
			{"Visible", "bool", "false", "Shows the cover and blocks the content."},
			{"Label", "string", "locale Loading", "Accessible name of the wait."},
			{"ShowLabel", "bool", "false", "Prints Label under the spinner."},
		},
		Tokens: []string{"Background", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "banner", Title: "Banner", Group: "Display", Order: 290,
		Description: "A persistent message bar across the top of a window or pane.",
		Examples: []ComponentExample{
			{Title: "Dismissible", Description: "Dismissing folds the bar up; actions wrap below the text when narrow.",
				Demo: dataDisplayDemo(bannerDismissibleBody)},
			{Title: "Statuses", Demo: dataDisplayDemo(bannerStatusesBody)},
		},
		Code: `mygoui.Banner(c, mygoui.BannerOptions{
    Text:     "A new version is ready to install.",
    Status:   mygoui.ColorInfo,
    ShowIcon: true,
    Action: func() {
        mygoui.Button(c, mygoui.ButtonOptions{Label: "Restart", OnClick: next})
    },
    OnDismiss: dismiss,
})`,
		Params: []ComponentParam{
			{"Text", "string", "\"\"", "The message."},
			{"Status", "mygoui.StatusColor", "ColorDefault", "Tint and default icon."},
			{"Icon", "*ui.SVG", "nil", "Leading glyph."},
			{"ShowIcon", "bool", "false", "Uses the status glyph when Icon is nil."},
			{"Action", "func()", "nil", "Controls after the text."},
			{"OnDismiss", "func()", "nil", "Adds a close button and receives its press."},
		},
		Tokens: []string{"MutedBackground", "Border", "Foreground", "Info", "Success", "Warning", "Danger"},
	})
}

func orderDescriptions(c *ui.Context) []mygoui.DescriptionsItem {
	return []mygoui.DescriptionsItem{
		{Label: "Customer", Value: "Ada Lovelace"},
		{Label: "Order", Value: "#10482"},
		{Label: "Placed", Value: "12 Nov 2026"},
		{Label: "Status", Content: func() {
			mygoui.Badge(c, mygoui.BadgeOptions{Label: "Paid"})
		}},
		{Label: "Shipping address", Value: "12 St James's Square, London SW1Y 4JH", Span: 2},
		{Label: "Notes", Value: "Leave with the concierge."},
	}
}

func descriptionsPlainBody(c *ui.Context, _ *dataDisplayState) {
	mygoui.Descriptions(c, mygoui.DescriptionsOptions{
		Title: "Order details", Columns: 3, Items: orderDescriptions(c),
	})
}

func descriptionsBorderedBody(c *ui.Context, _ *dataDisplayState) {
	mygoui.Descriptions(c, mygoui.DescriptionsOptions{
		Columns: 2, Bordered: true, Items: orderDescriptions(c),
	})
}

func statisticTrendsBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Refresh", Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
			OnClick: func() {
				st.revenue += rand.IntN(4000) - 1500
				st.latency += rand.IntN(120) - 50
				st.errs = max(0, st.errs+rand.IntN(9)-4)
				st.churn = max(0, st.churn+rand.IntN(7)-3)
			},
		})
		ui.Row(c).FillWidth().Wrap().GapX(theme.Space48).GapY(theme.Space24).Children(func() {
			mygoui.Statistic(c, mygoui.StatisticOptions{
				Label: "Revenue", Value: groupThousands(st.revenue), Prefix: "$",
				Delta: "12.5% vs last week", Trend: mygoui.TrendUp,
			})
			mygoui.Statistic(c, mygoui.StatisticOptions{
				Label: "Latency", Value: fmt.Sprint(st.latency), Suffix: "ms",
				Delta: "31 ms", Trend: mygoui.TrendUp, DeltaStatus: mygoui.ColorDanger,
			})
			mygoui.Statistic(c, mygoui.StatisticOptions{
				Label: "Errors", Value: fmt.Sprint(st.errs),
				Delta: "3 fewer", Trend: mygoui.TrendDown, DeltaStatus: mygoui.ColorSuccess,
			})
			mygoui.Statistic(c, mygoui.StatisticOptions{
				Label: "Churn", Value: fmt.Sprintf("%d.%d%%", st.churn/10, st.churn%10), Delta: "No change",
			})
		})
	})
}

func numberFlipCounterBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		mygoui.NumberFlip(c, mygoui.NumberFlipOptions{Value: groupThousands(st.count)})
		ui.Row(c).Gap(theme.Space8).Children(func() {
			step := func(label string, by func() int) {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: label, Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
					OnClick: func() { st.count = max(0, st.count+by()) },
				})
			}
			step("-1", func() int { return -1 })
			step("+1", func() int { return 1 })
			step("+100", func() int { return 100 })
			step("Randomize", func() int { return rand.IntN(20000) - st.count })
		})
	})
}

// groupThousands formats n with a comma between thousands.
func groupThousands(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0 && s[i-1] != '-'; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func descriptionsNarrowBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.truncate, Label: "Truncate"})
		mygoui.Descriptions(c, mygoui.DescriptionsOptions{
			Bordered: true, Columns: 2, Truncate: st.truncate, Items: []mygoui.DescriptionsItem{
				{Label: "Customer", Value: "Ada Lovelace"},
				{Label: "Reference", Value: "ORD-2026-000482-ALPHA-BETA-GAMMA-DELTA"},
				{Label: "Shipping address", Value: "12 St James's Square, London SW1Y 4JH, United Kingdom"},
				{Label: "Notes"},
			},
		})
	})
}

func descriptionsLoadingBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.loading, Label: "Loading"})
		mygoui.Descriptions(c, mygoui.DescriptionsOptions{
			Title: "Order details", Columns: 3, Items: orderDescriptions(c), Loading: st.loading,
		})
	})
}

func statisticContentBody(c *ui.Context, _ *dataDisplayState) {
	ui.Row(c).FillWidth().Wrap().GapX(theme.Space48).GapY(theme.Space24).Children(func() {
		ui.Column(c).Width(theme.Scale(c, 200)).MinWidth(0).Children(func() {
			mygoui.Statistic(c, mygoui.StatisticOptions{
				Label: "Storage", Value: "412", Suffix: "GB of 1 TB",
				Content: func() {
					mygoui.ProgressBar(c, mygoui.ProgressOptions{Label: "Storage used", Value: 41, Max: 100})
				},
			})
		})
		ui.Column(c).Width(theme.Scale(c, 200)).MinWidth(0).Children(func() {
			mygoui.Statistic(c, mygoui.StatisticOptions{Label: "Lifetime revenue", Value: "12,480,000,000.00", Prefix: "$"})
		})
	})
}

func statisticLoadingBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.loading, Label: "Loading"})
		ui.Row(c).FillWidth().Wrap().GapX(theme.Space48).GapY(theme.Space24).Children(func() {
			mygoui.Statistic(c, mygoui.StatisticOptions{
				Label: "Revenue", Value: "12,480", Prefix: "$", Loading: st.loading,
				Delta: "12.5% vs last week", Trend: mygoui.TrendUp,
			})
			mygoui.Statistic(c, mygoui.StatisticOptions{
				Label: "Latency", Value: "248", Suffix: "ms", Loading: st.loading,
			})
		})
	})
}

func timelineEventsBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.deploying, Label: "Rollout in progress"})
		last := mygoui.TimelineItem{Title: "Rolled out to everyone", Time: "09:41", Status: mygoui.ColorSuccess}
		if st.deploying {
			last = mygoui.TimelineItem{Title: "Rolling out", Description: "40% of traffic."}
		}
		mygoui.Timeline(c, mygoui.TimelineOptions{
			Label: "Release", Pending: st.deploying, Recent: 5, Expanded: &st.earlierOpen,
			EarlierLabel: "1 earlier",
			Items: []mygoui.TimelineItem{
				{Title: "Branch created", Time: "08:50"},
				{Title: "Build passed", Time: "09:12", Status: mygoui.ColorSuccess,
					Description: "All 412 checks are green."},
				{Title: "Deployed to staging", Time: "09:20", Status: mygoui.ColorInfo, Icon: lucide.Rocket(),
					OnClick: func() { st.deployOpened++ },
					Content: func() {
						mygoui.Tag(c, mygoui.TagOptions{Label: "v1.4.2"})
					}},
				{Title: "Smoke test failed", Time: "09:31", Status: mygoui.ColorDanger,
					Description: "Checkout returned 502 for 3 of 40 requests."},
				{Title: "Fix merged", Time: "09:38"},
				last,
			},
		})
	})
}

func statusDotStatusesBody(c *ui.Context, _ *dataDisplayState) {
	ui.Row(c).Wrap().GapX(theme.Space24).GapY(theme.Space8).Children(func() {
		mygoui.StatusDot(c, mygoui.StatusDotOptions{Label: "Offline"})
		mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: mygoui.ColorSuccess, Label: "Online", Pulse: true})
		mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: mygoui.ColorInfo, Label: "Connecting", Breathe: true})
		mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: mygoui.ColorWarning, Label: "Degraded"})
		mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: mygoui.ColorDanger, Label: "Down"})
		mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: mygoui.ColorInfo, Label: "Maintenance"})
	})
}

func statusDotPulseBody(c *ui.Context, _ *dataDisplayState) {
	ui.Row(c).Wrap().GapX(theme.Space24).GapY(theme.Space8).Children(func() {
		mygoui.StatusDot(c, mygoui.StatusDotOptions{
			Status: mygoui.ColorDanger, Label: "Recording", Pulse: true,
		})
		mygoui.StatusDot(c, mygoui.StatusDotOptions{
			Status: mygoui.ColorInfo, Label: "Syncing", Breathe: true,
		})
	})
}

func highlightSearchBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		form.Input(c, &st.query, form.InputOptions{Label: "Find", FullWidth: true})
		form.Switch(c, form.SwitchOptions{On: &st.ignoreCase, Label: "Ignore case"})
		mygoui.Highlight(c, mygoui.HighlightOptions{
			Text: "Go builds native desktop apps: go get the module, compose its components " +
				"and ship one binary. 日本語のテキストも文字単位で正しく扱います。",
			Queries: []string{st.query}, IgnoreCase: st.ignoreCase,
		})
	})
}

const spoilerNotes = "Every element is build-scoped, state stays with the caller and motion is optional. " +
	"Components share one set of tokens, so a custom surface reads like the built-in ones. " +
	"Lists virtualize their rows, overlays return focus to their trigger, and text keeps its native " +
	"selection and input methods. Nothing in the library starts a goroutine or touches the disk."

func spoilerLongBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		mygoui.Spoiler(c, mygoui.SpoilerOptions{
			Expanded: &st.notesExpanded, MaxHeight: theme.Scale(c, 72),
		}, func() {
			mygoui.Typography(c, spoilerNotes, mygoui.TypographyOptions{Type: theme.TypeBodySM})
		})
	})
}

func spoilerFitsBody(c *ui.Context, _ *dataDisplayState) {
	formsFieldStack(c, func() {
		mygoui.Spoiler(c, mygoui.SpoilerOptions{MaxHeight: theme.Scale(c, 72)}, func() {
			mygoui.Typography(c, "Nothing in the library starts a goroutine.", mygoui.TypographyOptions{
				Type: theme.TypeBodySM,
			})
		})
	})
}

func loadingOverlaySavingBody(c *ui.Context, st *dataDisplayState) {
	formsFieldStack(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.saving, Label: "Saving"})
		mygoui.LoadingOverlay(c, mygoui.LoadingOverlayOptions{
			Visible: st.saving, Label: "Saving invoice", ShowLabel: true,
		}, func() {
			ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
				form.Input(c, &st.invoiceName, form.InputOptions{Label: "Name", FullWidth: true})
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Save", OnClick: func() { st.saving = true },
				}).AlignSelf(ui.Start)
			})
		})
	})
}

func bannerDismissibleBody(c *ui.Context, st *dataDisplayState) {
	t := theme.Colors(c)
	pane := ui.Column(c).FillWidth().MinWidth(0).Radius(theme.RadiusScale().LG).Clip().
		Border(theme.Hairline, t.Border)
	pane.Children(func() {
		if !st.updateHidden {
			mygoui.Banner(c, mygoui.BannerOptions{
				Text: "A new version is ready to install.", Status: mygoui.ColorInfo, ShowIcon: true,
				Action: func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Restart", Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
						OnClick: func() { st.restarts++ },
					})
				},
				OnDismiss: func() { st.updateHidden = true },
			})
		}
		ui.Row(c).FillWidth().Padding(theme.Space16).Gap(theme.Space12).AlignItems(ui.Center).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Show banner", Variant: mygoui.ButtonOutline, Disabled: !st.updateHidden,
				OnClick: func() { st.updateHidden = false },
			})
			if st.restarts > 0 {
				mygoui.Typography(c, "Restart requested", mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Muted: true,
				})
			}
		})
	})
}

func bannerStatusesBody(c *ui.Context, _ *dataDisplayState) {
	ui.Column(c).FillWidth().MinWidth(0).Children(func() {
		mygoui.Banner(c, mygoui.BannerOptions{Text: "You are viewing a read-only copy."})
		mygoui.Banner(c, mygoui.BannerOptions{
			Text: "Changes are saved to this device.", Status: mygoui.ColorSuccess, ShowIcon: true,
		})
		mygoui.Banner(c, mygoui.BannerOptions{
			Text: "Your trial ends in 3 days.", Status: mygoui.ColorWarning, ShowIcon: true,
		})
		mygoui.Banner(c, mygoui.BannerOptions{
			Text: "Billing failed. Update your card to keep the workspace.", Status: mygoui.ColorDanger, ShowIcon: true,
		})
	})
}
