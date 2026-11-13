package gallery

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// absorbedOverlaysState is the caller-owned state behind the floating and
// image viewer demos. Every open flag starts false.
type absorbedOverlaysState struct {
	ready bool

	panelOpen bool
	panelPos  ui.Point
	notes     int

	sheetOpen   bool
	sheetHeight float32

	viewerOpen  bool
	viewerIndex int
}

func absorbedOverlaysStateOf(s *State) *absorbedOverlaysState {
	st := StateOf[absorbedOverlaysState](s, "absorbed-overlays")
	if !st.ready {
		st.ready = true
		st.panelPos = ui.Point{X: 320, Y: 220}
		st.sheetHeight = 180
	}
	return st
}

func absorbedOverlaysDemo(body func(c *ui.Context, st *absorbedOverlaysState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, absorbedOverlaysStateOf(s)) }
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "floatingpanel", Title: "FloatingPanel", Group: "Overlays", Order: 90,
		Description: "A non-modal panel dragged by its title bar within the window.",
		Examples: []ComponentExample{
			{Title: "Draggable", Demo: absorbedOverlaysDemo(floatingPanelBody)},
		},
		Code: `overlay.FloatingPanel(c, &open, overlay.FloatingPanelOptions{
    Title:    "Inspector",
    Position: &position,
    Width:    280,
    Body: func() {
        mygoui.Paragraph(c, "The page behind stays usable.", true)
    },
    OnDismiss: save,
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "Title bar text and accessible name."},
			{"Body", "func()", "nil", "Scrolling content."},
			{"Position", "*ui.Point", "nil", "Caller-owned top-left corner in window DIPs; nil opens top right."},
			{"Width", "float32", "320", "Panel width in DIPs."},
			{"HideClose", "bool", "false", "Leaves out the close button."},
			{"OnDismiss", "func()", "nil", "Runs when the close button or Escape closes it."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "floatingsheet", Title: "FloatingSheet", Group: "Overlays", Order: 100,
		Description: "A non-modal sheet rising from the bottom of its container.",
		Examples: []ComponentExample{
			{Title: "Resizable", Description: "The handle drags or takes the Up and Down arrows.",
				Demo: absorbedOverlaysDemo(floatingSheetBody)},
		},
		Code: `ui.Column(c).FillWidth().Height(320).Children(func() {
    mygoui.Button(c, mygoui.ButtonOptions{
        Label: "Details", OnClick: func() { open = true },
    })
    // Built last, the sheet covers the content it floats over.
    overlay.FloatingSheet(c, &open, overlay.FloatingSheetOptions{
        Title:     "Details",
        Height:    &height,
        MinHeight: 120,
        Body: func() {
            mygoui.Paragraph(c, "Drag the handle to resize.", true)
        },
    })
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "Header text and accessible name."},
			{"Body", "func()", "nil", "Scrolling content."},
			{"Height", "*float32", "nil", "Caller-owned height; nil keeps it in the sheet."},
			{"MinHeight", "float32", "160", "Smallest height in DIPs."},
			{"MaxHeight", "float32", "80% of container", "Largest height in DIPs."},
			{"HideHandle", "bool", "false", "Fixes the height."},
			{"HideClose", "bool", "false", "Leaves out the close button."},
			{"OnDismiss", "func()", "nil", "Runs when the close button or Escape closes it."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "imageviewer", Title: "ImageViewer", Group: "Overlays", Order: 110,
		Description: "A modal preview that zooms, pans and steps through images.",
		Examples: []ComponentExample{
			{Title: "Gallery", Demo: absorbedOverlaysDemo(imageViewerBody)},
		},
		Code: `mygoui.Image(c, mygoui.ImageOptions{
    Source: photos[0].Source, Label: "Open preview", Width: 160,
    OnClick: func() { index, open = 0, true },
})
overlay.ImageViewer(c, &open, overlay.ImageViewerOptions{
    Images: photos,
    Index:  &index,
})`,
		Params: []ComponentParam{
			{"Images", "[]overlay.ImageViewerItem", "nil", "Already-decoded Source and accessible Label of each image."},
			{"Index", "*int", "nil", "Caller-owned image shown; nil keeps it in the viewer."},
			{"OnDismiss", "func()", "nil", "Runs when Escape, the close button or an outside press closes it."},
		},
		Tokens: []string{"Backdrop", "Popover", "PopoverForeground", "Border", "MutedForeground"},
	})
}

func floatingPanelBody(c *ui.Context, st *absorbedOverlaysState) {
	label := "Open panel"
	if st.panelOpen {
		label = "Close panel"
	}
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Variant: mygoui.ButtonOutline,
		OnClick: func() { st.panelOpen = !st.panelOpen },
	}).AlignSelf(ui.Center)
	overlay.FloatingPanel(c, &st.panelOpen, overlay.FloatingPanelOptions{
		Title: "Inspector", Position: &st.panelPos, Width: 280,
		Body: func() {
			ui.Column(c).FillWidth().Gap(theme.Space12).AlignItems(ui.Start).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Add note", Size: theme.ButtonSmall, OnClick: func() { st.notes++ },
				})
				for i := range st.notes {
					formsMutedText(c.Key(i), "Note")
				}
			})
		},
	})
}

func floatingSheetBody(c *ui.Context, st *absorbedOverlaysState) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 520)).Height(theme.Scale(c, 320)).
		Clip().Radius(theme.RadiusScale().LG).Border(theme.Hairline, t.Border).
		Background(t.Background).Label("Floating sheet container").Children(func() {
		ui.Column(c).Fill().Padding(theme.Space24).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Open sheet", Variant: mygoui.ButtonOutline,
				Disabled: st.sheetOpen, OnClick: func() { st.sheetOpen = true },
			}).AlignSelf(ui.Center)
		})
		overlay.FloatingSheet(c, &st.sheetOpen, overlay.FloatingSheetOptions{
			Title: "Details", Height: &st.sheetHeight, MinHeight: 120,
			Body: func() {
				ui.Column(c).FillWidth().Gap(theme.Space8).Children(func() {
					for _, line := range [...]string{"Created today", "Edited by Ada", "Shared with 3 people"} {
						formsMutedText(c, line)
					}
				})
			},
		})
	})
}

func imageViewerBody(c *ui.Context, st *absorbedOverlaysState) {
	pictures := demoPictures()
	images := make([]overlay.ImageViewerItem, len(pictures))
	ui.Row(c).Wrap().Gap(theme.Space12).Justify(ui.Center).Children(func() {
		for i, p := range pictures {
			images[i] = overlay.ImageViewerItem{Source: p.Source, Label: p.Label}
			mygoui.Image(c.Key(i), mygoui.ImageOptions{
				Source: p.Source, Label: "Preview " + p.Label, Width: theme.Scale(c, 120), AspectRatio: 4.0 / 3,
				OnClick: func() { st.viewerIndex, st.viewerOpen = i, true },
			})
		}
	})
	overlay.ImageViewer(c, &st.viewerOpen, overlay.ImageViewerOptions{Images: images, Index: &st.viewerIndex})
}
