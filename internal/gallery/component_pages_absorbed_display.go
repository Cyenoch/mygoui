package gallery

import (
	"slices"
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// absorbedDisplayState is the caller-owned state behind the absorbed display,
// navigation, collection and layout demos.
type absorbedDisplayState struct {
	ready bool

	step, verticalStep int

	tasks, people []list.SortableItem

	sideWidth, bottomHeight float32
	sideCollapsed           bool

	imageClicks   int
	resultNote    string
	guideHidden   bool
	guideTipShown bool
}

func absorbedDisplayStateOf(s *State) *absorbedDisplayState {
	st := StateOf[absorbedDisplayState](s, "absorbed-display")
	if !st.ready {
		st.ready = true
		st.step, st.verticalStep = 1, 1
		st.tasks = []list.SortableItem{
			{ID: "design", Label: "Design review"},
			{ID: "build", Label: "Build release"},
			{ID: "notes", Label: "Write release notes"},
			{ID: "ship", Label: "Ship"},
		}
		st.people = []list.SortableItem{
			{ID: "ada", Label: "Ada Lovelace"},
			{ID: "grace", Label: "Grace Hopper"},
			{ID: "alan", Label: "Alan Turing"},
		}
		st.sideWidth, st.bottomHeight = 180, 96
		st.guideTipShown = true
	}
	return st
}

func absorbedDisplayDemo(body func(c *ui.Context, st *absorbedDisplayState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, absorbedDisplayStateOf(s)) }
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "image", Title: "Image", Group: "Display", Order: 140,
		Description: "A caller-loaded picture with loading and failed states.",
		Examples: []ComponentExample{
			{Title: "Fit", Description: "Cover fills the box; contain shows the whole picture.",
				Demo: absorbedDisplayDemo(imageFitBody)},
			{Title: "States", Description: "Loading skeleton and failed fallback.",
				Demo: absorbedDisplayDemo(imageStatesBody)},
			{Title: "Pressable", Demo: absorbedDisplayDemo(imagePressableBody)},
		},
		Code: `mygoui.Image(c, mygoui.ImageOptions{
    Source:      photo,
    Label:       "Team offsite",
    Width:       240,
    AspectRatio: 4.0 / 3,
    Fit:         mygoui.ImageCover,
    OnClick:     openPreview,
})`,
		Params: []ComponentParam{
			{"Source", "ui.ImageSource", "nil", "Already-decoded picture; nil shows the fallback."},
			{"Label", "string", "\"\"", "Accessible name."},
			{"Width", "float32", "fill", "Box width; zero fills the parent."},
			{"Height", "float32", "0", "Box height; wins over AspectRatio."},
			{"AspectRatio", "float32", "picture ratio", "Width over height; 16:9 without a picture."},
			{"Fit", "mygoui.ImageFit", "ImageCover", "Cover or contain."},
			{"Radius", "float32", "theme LG", "Corner radius."},
			{"Status", "mygoui.ImageStatus", "ImageReady", "Ready, loading or failed."},
			{"Fallback", "func()", "nil", "Content under the failed icon."},
			{"OnClick", "func()", "nil", "Makes a shown picture a button."},
		},
		Tokens: []string{"MutedBackground", "MutedForeground", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "result", Title: "Result", Group: "Patterns", Order: 150,
		Description: "The outcome of an operation: status icon, text and actions.",
		Examples: []ComponentExample{
			{Title: "Success", Demo: absorbedDisplayDemo(resultSuccessBody)},
			{Title: "Failure", Description: "Content adds detail under the actions.",
				Demo: absorbedDisplayDemo(resultFailureBody)},
		},
		Code: `mygoui.Result(c, mygoui.ResultOptions{
    Status:      mygoui.ColorSuccess,
    Title:       "Payment received",
    Description: "A receipt is on its way to your inbox.",
    Actions: func() {
        mygoui.Button(c, mygoui.ButtonOptions{Label: "Done", OnClick: back})
    },
})`,
		Params: []ComponentParam{
			{"Status", "mygoui.StatusColor", "ColorDefault", "Tint and default icon."},
			{"Title", "string", "\"\"", "Outcome heading."},
			{"Description", "string", "\"\"", "Supporting line."},
			{"Icon", "*ui.SVG", "status icon", "Replaces the glyph and keeps the tint."},
			{"Actions", "func()", "nil", "Buttons under the text."},
			{"Content", "func()", "nil", "Supporting detail under the actions."},
		},
		Tokens: []string{"Card", "Border", "Foreground", "MutedForeground", "SuccessForeground", "DestructiveForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "guidecard", Title: "GuideCard", Group: "Patterns", Order: 160,
		Description: "A dismissible card that introduces a feature.",
		Examples: []ComponentExample{
			{Title: "With cover", Demo: absorbedDisplayDemo(guideCardCoverBody)},
			{Title: "Text only", Demo: absorbedDisplayDemo(guideCardTextBody)},
		},
		Code: `mygoui.GuideCard(c, mygoui.GuideCardOptions{
    Title:       "Meet the command palette",
    Description: "Jump to any page without leaving the keyboard.",
    Cover:       cover,
    Actions: func() {
        mygoui.Button(c, mygoui.ButtonOptions{Label: "Try it", OnClick: openDocs})
    },
    OnDismiss: dismiss,
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "Card heading and accessible name."},
			{"Description", "string", "\"\"", "Supporting line."},
			{"Cover", "ui.ImageSource", "nil", "Already-decoded picture across the top, at its own ratio."},
			{"CoverContent", "func()", "nil", "Custom cover; wins over Cover."},
			{"Actions", "func()", "nil", "Footer buttons."},
			{"OnDismiss", "func()", "nil", "Adds the close button; the caller stops rendering the card."},
		},
		Tokens: []string{"Card", "CardForeground", "Border", "MutedForeground", "Surface"},
	})

	RegisterComponent(ComponentPage{
		ID: "steps", Title: "Steps", Group: "Navigation", Order: 110,
		Description: "The numbered steps of a process joined by connectors.",
		Examples: []ComponentExample{
			{Title: "Horizontal", Demo: absorbedDisplayDemo(stepsHorizontalBody)},
			{Title: "Vertical", Description: "Pressable steps with descriptions and a failed step.",
				Demo: absorbedDisplayDemo(stepsVerticalBody)},
		},
		Code: `navigation.Steps(c, navigation.StepsOptions{
    Label: "Checkout",
    Items: []navigation.StepItem{
        {Title: "Cart"},
        {Title: "Shipping", Description: "Address and method"},
        {Title: "Payment"},
    },
    Current:  step,
    OnSelect: func(index int) { step = index },
})`,
		Params: []ComponentParam{
			{"Items", "[]navigation.StepItem", "nil", "Steps: Title, Description and Status (StepAuto or StepError)."},
			{"Current", "int", "0", "Index in progress; earlier steps are completed, len(Items) completes all."},
			{"Orientation", "navigation.StepsOrientation", "StepsHorizontal", "Row or column."},
			{"Label", "string", "\"\"", "Accessible name of the list."},
			{"OnSelect", "func(index int)", "nil", "Makes each step a button."},
		},
		Tokens: []string{"Primary", "PrimaryForeground", "MutedBackground", "MutedForeground", "Foreground", "Border", "Destructive", "Accent"},
	})

	RegisterComponent(ComponentPage{
		ID: "sortablelist", Title: "SortableList", Group: "Collections", Order: 20,
		Description: "Rows reordered by dragging a grip or with Alt+Up and Alt+Down.",
		Examples: []ComponentExample{
			{Title: "Reorder", Demo: absorbedDisplayDemo(func(c *ui.Context, st *absorbedDisplayState) {
				formsFieldStack(c, func() {
					list.SortableList(c, list.SortableListOptions{
						Label: "Release tasks", Items: st.tasks,
						OnMove: func(from, to int) { st.tasks = moveSortable(st.tasks, from, to) },
					})
				})
			})},
			{Title: "Custom rows", Demo: absorbedDisplayDemo(sortableCustomBody)},
			{Title: "Disabled", Demo: absorbedDisplayDemo(func(c *ui.Context, st *absorbedDisplayState) {
				formsFieldStack(c, func() {
					list.SortableList(c, list.SortableListOptions{
						Label: "Locked tasks", Items: st.tasks[:2], Disabled: true,
					})
				})
			})},
		},
		Code: `list.SortableList(c, list.SortableListOptions{
    Label: "Tasks",
    Items: tasks,
    OnMove: func(from, to int) {
        task := tasks[from]
        tasks = append(tasks[:from], tasks[from+1:]...)
        tasks = append(tasks[:to], append([]list.SortableItem{task}, tasks[to:]...)...)
    },
})`,
		Params: []ComponentParam{
			{"Items", "[]list.SortableItem", "nil", "Rows: stable ID, Label and optional Content."},
			{"Label", "string", "\"\"", "Accessible name of the list."},
			{"Disabled", "bool", "false", "Dims the grips and stops reordering."},
			{"OnMove", "func(from, to int)", "nil", "The item at from goes to index to, counted after its removal."},
		},
		Tokens: []string{"Card", "CardForeground", "Border", "MutedForeground", "Accent"},
	})

	RegisterComponent(ComponentPage{
		ID: "resizablepanel", Title: "ResizablePanel", Group: "Foundations", Order: 120,
		Description: "A panel resized from one edge by pointer or keyboard.",
		Examples: []ComponentExample{
			{Title: "Collapsible side", Description: "Double click or Enter on the edge toggles it.",
				Demo: absorbedDisplayDemo(resizableSideBody)},
			{Title: "Bottom edge", Demo: absorbedDisplayDemo(resizableBottomBody)},
		},
		Code: `mygoui.ResizablePanel(c, mygoui.ResizablePanelOptions{
    Label:     "Resize sidebar",
    Side:      mygoui.PanelEnd,
    Size:      &width,
    Min:       160,
    Max:       360,
    Collapsed: &collapsed,
    OnResize:  func(size float32) { save() },
}, func() {
    mygoui.Paragraph(c, "Sidebar", true)
})`,
		Params: []ComponentParam{
			{"Label", "string", "locale Resize", "Accessible name of the handle."},
			{"Side", "mygoui.PanelSide", "PanelEnd", "Edge carrying the handle: end, start, top or bottom."},
			{"Size", "*float32", "nil", "Caller-owned size in DIPs; nil keeps it in the panel."},
			{"Min", "float32", "0", "Smallest size."},
			{"Max", "float32", "0", "Largest size; nonpositive sets no upper bound."},
			{"Collapsed", "*bool", "nil", "Non-nil makes the panel collapsible."},
			{"OnResize", "func(size float32)", "nil", "Runs after the size changed."},
			{"OnCollapseChange", "func(collapsed bool)", "nil", "Runs after the panel collapsed or expanded."},
		},
		Tokens: []string{"Border", "Ring"},
	})
}

// moveSortable applies a SortableList move: the item at from goes to index
// to, counted after its removal.
func moveSortable(items []list.SortableItem, from, to int) []list.SortableItem {
	item := items[from]
	return slices.Insert(slices.Delete(items, from, from+1), to, item)
}

func imageFitBody(c *ui.Context, _ *absorbedDisplayState) {
	picture := demoPictures()[1]
	ui.Row(c).Wrap().Gap(theme.Space16).Justify(ui.Center).Children(func() {
		mygoui.Image(c, mygoui.ImageOptions{
			Source: picture.Source, Label: picture.Label + ", cover",
			Width: theme.Scale(c, 160), AspectRatio: 1,
		})
		mygoui.Image(c, mygoui.ImageOptions{
			Source: picture.Source, Label: picture.Label + ", contain",
			Width: theme.Scale(c, 160), AspectRatio: 1, Fit: mygoui.ImageContain,
		}).Background(theme.Colors(c).MutedBackground)
	})
}

func imageStatesBody(c *ui.Context, _ *absorbedDisplayState) {
	ui.Row(c).Wrap().Gap(theme.Space16).Justify(ui.Center).Children(func() {
		mygoui.Image(c, mygoui.ImageOptions{
			Label: "Loading picture", Width: theme.Scale(c, 200), Status: mygoui.ImageLoading,
		})
		mygoui.Image(c, mygoui.ImageOptions{
			Label: "Missing picture", Width: theme.Scale(c, 200), Status: mygoui.ImageFailed,
			Fallback: func() { formsMutedText(c, "Not found") },
		})
	})
}

func imagePressableBody(c *ui.Context, st *absorbedDisplayState) {
	picture := demoPictures()[0]
	ui.Column(c).Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
		mygoui.Image(c, mygoui.ImageOptions{
			Source: picture.Source, Label: "Open " + picture.Label,
			Width: theme.Scale(c, 240), OnClick: func() { st.imageClicks++ },
		})
		formsMutedText(c, strconv.Itoa(st.imageClicks)+" clicks")
	})
}

func resultSuccessBody(c *ui.Context, st *absorbedDisplayState) {
	mygoui.Result(c, mygoui.ResultOptions{
		Status: mygoui.ColorSuccess, Title: "Payment received",
		Description: "A receipt is on its way to your inbox.",
		Actions: func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "View order", OnClick: func() { st.resultNote = "Order opened" },
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Done", Variant: mygoui.ButtonOutline,
				OnClick: func() { st.resultNote = "Done" },
			})
		},
		Content: func() {
			if st.resultNote != "" {
				formsMutedText(c, st.resultNote)
			}
		},
	})
}

func resultFailureBody(c *ui.Context, st *absorbedDisplayState) {
	mygoui.Result(c, mygoui.ResultOptions{
		Status: mygoui.ColorDanger, Title: "Deploy failed",
		Description: "The build stopped at the test stage.",
		Actions: func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Retry", Icon: lucide.RotateCcw(), Variant: mygoui.ButtonOutline,
				OnClick: func() { st.resultNote = "Retrying" },
			})
		},
		Content: func() {
			mygoui.Typography(c, "exit status 1: TestCheckout", mygoui.TypographyOptions{
				Type: theme.TypeCode, Selectable: true,
			})
		},
	})
}

// guideCardFrame keeps a guide card at its reading width.
func guideCardFrame(c *ui.Context, body func()) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 340)).
		AlignItems(ui.Stretch).Children(body)
}

func guideCardCoverBody(c *ui.Context, st *absorbedDisplayState) {
	if st.guideHidden {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Show again", Variant: mygoui.ButtonOutline,
			OnClick: func() { st.guideHidden = false },
		}).AlignSelf(ui.Center)
		return
	}
	guideCardFrame(c, func() {
		mygoui.GuideCard(c, mygoui.GuideCardOptions{
			Title:       "Meet the command palette",
			Description: "Jump to any page without leaving the keyboard.",
			Cover:       demoPictures()[2].Source,
			Actions: func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Got it", Size: theme.ButtonSmall,
					OnClick: func() { st.guideHidden = true },
				})
			},
			OnDismiss: func() { st.guideHidden = true },
		})
	})
}

func guideCardTextBody(c *ui.Context, st *absorbedDisplayState) {
	if !st.guideTipShown {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Show tip", Variant: mygoui.ButtonOutline,
			OnClick: func() { st.guideTipShown = true },
		}).AlignSelf(ui.Center)
		return
	}
	guideCardFrame(c, func() {
		mygoui.GuideCard(c, mygoui.GuideCardOptions{
			Title:       "Drag to reorder",
			Description: "Hold a row's grip and move it.",
			OnDismiss:   func() { st.guideTipShown = false },
		})
	})
}

var checkoutSteps = []navigation.StepItem{{Title: "Cart"}, {Title: "Shipping"}, {Title: "Payment"}}

func stepsHorizontalBody(c *ui.Context, st *absorbedDisplayState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 560)).
		Gap(theme.Space24).AlignItems(ui.Center).Children(func() {
		navigation.Steps(c, navigation.StepsOptions{
			Label: "Checkout", Items: checkoutSteps, Current: st.step,
		})
		ui.Row(c).Gap(theme.Space8).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Back", Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
				Disabled: st.step == 0, OnClick: func() { st.step-- },
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Next", Size: theme.ButtonSmall,
				Disabled: st.step >= len(checkoutSteps), OnClick: func() { st.step++ },
			})
		})
	})
}

func stepsVerticalBody(c *ui.Context, st *absorbedDisplayState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 320)).Children(func() {
		navigation.Steps(c, navigation.StepsOptions{
			Label: "Deploy", Orientation: navigation.StepsVertical,
			Items: []navigation.StepItem{
				{Title: "Build", Description: "Compile and bundle"},
				{Title: "Test", Description: "Unit and integration"},
				{Title: "Review", Description: "Checks did not pass", Status: navigation.StepError},
				{Title: "Release", Description: "Publish to production"},
			},
			Current:  st.verticalStep,
			OnSelect: func(index int) { st.verticalStep = index },
		})
	})
}

func sortableCustomBody(c *ui.Context, st *absorbedDisplayState) {
	items := make([]list.SortableItem, len(st.people))
	for i, person := range st.people {
		items[i] = person
		items[i].Content = func() {
			ui.Row(c).Gap(theme.Space8).AlignItems(ui.Center).MinWidth(0).Children(func() {
				mygoui.Avatar(c, mygoui.AvatarOptions{Name: person.Label, Size: mygoui.AvatarSmall})
				mygoui.Typography(c, person.Label, mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Truncate: true,
				})
			})
		}
	}
	formsFieldStack(c, func() {
		list.SortableList(c, list.SortableListOptions{
			Label: "Reviewers", Items: items,
			OnMove: func(from, to int) { st.people = moveSortable(st.people, from, to) },
		})
	})
}

// resizableFrame is the bordered container a panel demo resizes within.
func resizableFrame(c *ui.Context, label string, height float32, column bool, body func()) {
	t := theme.Colors(c)
	frame := ui.Row(c)
	if column {
		frame = ui.Column(c)
	}
	frame.FillWidth().MinWidth(0).Height(theme.Scale(c, height)).Clip().
		Radius(theme.RadiusScale().LG).Border(theme.Hairline, t.Border).
		Background(t.Background).Label(label).Children(body)
}

func resizableSideBody(c *ui.Context, st *absorbedDisplayState) {
	resizableFrame(c, "Collapsible panel demo", 220, false, func() {
		mygoui.ResizablePanel(c, mygoui.ResizablePanelOptions{
			Label: "Resize sidebar", Size: &st.sideWidth, Min: 120, Max: 280,
			Collapsed: &st.sideCollapsed,
		}, func() {
			ui.Column(c).Fill().Padding(theme.Space12).Gap(theme.Space4).Children(func() {
				for _, name := range [...]string{"Inbox", "Drafts", "Archive"} {
					mygoui.Typography(c, name, mygoui.TypographyOptions{Type: theme.TypeBodySM, Truncate: true})
				}
			})
		})
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center).Justify(ui.Center).Gap(theme.Space8).Children(func() {
			formsMutedText(c, strconv.Itoa(int(st.sideWidth))+" DIP")
			label := "Collapse"
			if st.sideCollapsed {
				label = "Expand"
			}
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: label, Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
				OnClick: func() { st.sideCollapsed = !st.sideCollapsed },
			}).AlignSelf(ui.Center)
		})
	})
}

func resizableBottomBody(c *ui.Context, st *absorbedDisplayState) {
	resizableFrame(c, "Bottom edge demo", 240, true, func() {
		mygoui.ResizablePanel(c, mygoui.ResizablePanelOptions{
			Label: "Resize editor", Side: mygoui.PanelBottom,
			Size: &st.bottomHeight, Min: 56, Max: 180,
		}, func() {
			ui.Column(c).Fill().AlignItems(ui.Center).Justify(ui.Center).Children(func() { formsMutedText(c, "Editor") })
		})
		ui.Column(c).Grow(1).Basis(0).MinHeight(0).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
			formsMutedText(c, "Terminal")
		})
	})
}
