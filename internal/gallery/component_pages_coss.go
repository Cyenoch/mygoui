package gallery

import (
	"fmt"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type cossState struct {
	Count       int
	Choice      int
	Draft       string
	Last        string
	Name        string
	Selected    []string
	Checked     bool
	Open        [4]bool
	Scroll      [5]ui.ScrollState
	Command     [2]navigation.CommandState
	Preview     [2]overlay.PreviewCardState
	PreviewHits int
}

func cossStateOf(s *State, page string) *cossState {
	return StateOf[cossState](s, "coss:"+page)
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "frame", Title: "Frame", Group: "Foundations", Order: 51,
		Description: "A muted frame around independently composed panels.",
		Examples: []ComponentExample{
			{Title: "Panel and actions", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "frame")
				mygoui.Frame(c, mygoui.FrameOptions{Label: "Project frame"}, func() {
					mygoui.FrameHeader(c, func() {
						mygoui.FrameTitle(c, "Project settings")
						mygoui.FrameDescription(c, "Workspace preferences")
					})
					mygoui.FramePanel(c, func() {
						form.Input(c, &st.Name, form.InputOptions{
							Label: "Project name", Placeholder: "Acme", FullWidth: true,
						})
					})
					mygoui.FrameFooter(c, func() {
						mygoui.Button(c, mygoui.ButtonOptions{
							Label: "Save project", OnClick: func() { st.Count++ },
						})
						ui.Text(c, fmt.Sprintf("Saved %d", st.Count))
					})
				})
			}},
			{Title: "Custom spacing", Demo: func(c *ui.Context, s *State) {
				mygoui.Frame(c, mygoui.FrameOptions{Padding: 8, Gap: 8, Label: "Spaced frame"}, func() {
					for _, title := range []string{"Access", "Notifications"} {
						mygoui.FramePanel(c, func() {
							mygoui.FrameTitle(c, title)
							mygoui.FrameDescription(c, "Independent panel")
						}).Padding(theme.Space12)
					}
				})
			}},
		},
		Code: `mygoui.Frame(c, mygoui.FrameOptions{Padding: 8, Gap: 8}, func() {
    mygoui.FrameHeader(c, func() {
        mygoui.FrameTitle(c, "Project settings")
        mygoui.FrameDescription(c, "Workspace preferences")
    })
    mygoui.FramePanel(c, func() {
        form.Input(c, &name, form.InputOptions{Label: "Name", FullWidth: true})
    })
    mygoui.FrameFooter(c, func() {
        mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: save})
    })
})`,
		Params: []ComponentParam{
			{"Padding / Gap", "float32", "4", "Frame inset and distance between panels."},
			{"Label", "string", "", "Accessible frame name."},
			{"FramePanel", "func()", "required", "Panel content; returned element accepts native layout overrides."},
			{"FrameHeader / FrameFooter", "func()", "required", "Frame-level metadata and actions."},
			{"FrameTitle / FrameDescription", "string", "required", "Heading and supporting text."},
		},
		Tokens: []string{"MutedBackground", "Background", "Foreground", "MutedForeground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "group", Title: "Group", Group: "Actions", Order: 21,
		Description: "Joined native controls with explicit outer corners.",
		Examples: []ComponentExample{
			{Title: "Default", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "group-default")
				separator := func(key string) mygoui.GroupItem {
					return mygoui.GroupItem{Key: key, Build: func([4]float32) ui.Element {
						return mygoui.Separator(c, mygoui.SeparatorOptions{Orientation: mygoui.Vertical}).
							Height(theme.Scale(c, theme.ButtonHeight(theme.ButtonMedium))).
							Background(theme.Colors(c).Input)
					}}
				}
				ui.Column(c).FillWidth().Gap(theme.Space8).Children(func() {
					ui.Row(c).FillWidth().Justify(ui.Center).Children(func() {
						mygoui.Group(c, mygoui.GroupOptions{Label: "File actions"},
							mygoui.GroupItem{Key: "files", Build: func(corners [4]float32) ui.Element {
								return mygoui.Button(c, mygoui.ButtonOptions{
									Label: "Files", Icon: lucide.Files(), Variant: mygoui.ButtonOutline,
									Corners: &corners, OnClick: func() { st.Last = "Files" },
								})
							}},
							separator("files-media"),
							mygoui.GroupItem{Key: "media", Build: func(corners [4]float32) ui.Element {
								return mygoui.Button(c, mygoui.ButtonOptions{
									Label: "Media", Icon: lucide.Film(), Variant: mygoui.ButtonOutline,
									Corners: &corners, OnClick: func() { st.Last = "Media" },
								})
							}},
							separator("media-more"),
							mygoui.GroupItem{Key: "more", Build: func(corners [4]float32) ui.Element {
								return mygoui.Button(c, mygoui.ButtonOptions{
									Label: "More", Icon: lucide.Ellipsis(), IconOnly: true,
									Variant: mygoui.ButtonOutline, Corners: &corners,
									OnClick: func() { st.Last = "More" },
								})
							}},
						)
					})
					if st.Last != "" {
						mygoui.Typography(c, st.Last, mygoui.TypographyOptions{
							Type: theme.TypeBodyXS, Muted: true,
						}).AlignSelf(ui.Center)
					}
				})
			}},
			{Title: "Mixed controls", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "group")
				mygoui.Group(c, mygoui.GroupOptions{Label: "Repository search"},
					mygoui.GroupItem{Key: "query", Build: func(corners [4]float32) ui.Element {
						return ui.TextInput(c, &st.Name).Label("Repository query").Placeholder("Find a repository").
							Width(theme.Scale(c, 200)).MinWidth(0).
							Radius(corners[0], corners[1], corners[2], corners[3])
					}},
					mygoui.GroupItem{Key: "search", Build: func(corners [4]float32) ui.Element {
						return mygoui.Button(c, mygoui.ButtonOptions{
							Label: "Find", Corners: &corners, OnClick: func() { st.Last = st.Name },
						})
					}},
				)
				if st.Last != "" {
					ui.Text(c, st.Last)
				}
			}},
			{Title: "Vertical", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "group")
				ui.Row(c).FillWidth().Justify(ui.Center).Children(func() {
					mygoui.Group(c, mygoui.GroupOptions{Orientation: mygoui.Vertical, Label: "Zoom controls"},
						mygoui.GroupItem{Key: "in", Build: func(corners [4]float32) ui.Element {
							return mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Zoom in", Icon: lucide.ZoomIn(), IconOnly: true,
								Variant: mygoui.ButtonOutline, Corners: &corners, OnClick: func() { st.Count++ },
							})
						}},
						mygoui.GroupItem{Key: "separator", Build: func([4]float32) ui.Element {
							return mygoui.Separator(c, mygoui.SeparatorOptions{}).Background(theme.Colors(c).Input)
						}},
						mygoui.GroupItem{Key: "out", Build: func(corners [4]float32) ui.Element {
							return mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Zoom out", Icon: lucide.ZoomOut(), IconOnly: true,
								Variant: mygoui.ButtonOutline, Corners: &corners, OnClick: func() { st.Count-- },
							})
						}},
					)
				})
			}},
			{Title: "Full width", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "group")
				items := make([]mygoui.GroupItem, 0, 5)
				for i, label := range []string{"Profile", "Security", "Billing"} {
					if i > 0 {
						items = append(items, mygoui.GroupItem{Key: label + "-separator", Build: func([4]float32) ui.Element {
							return mygoui.Separator(c, mygoui.SeparatorOptions{}).Background(theme.Colors(c).Input)
						}})
					}
					items = append(items, mygoui.GroupItem{Key: label, Build: func(corners [4]float32) ui.Element {
						return mygoui.Button(c, mygoui.ButtonOptions{
							Label: label, Variant: mygoui.ButtonOutline, Corners: &corners,
							OnClick: func() { st.Last = label },
						})
					}})
				}
				mygoui.Group(c, mygoui.GroupOptions{
					Orientation: mygoui.Vertical, FullWidth: true, Label: "Account actions",
				}, items...)
			}},
		},
		Code: `mygoui.Group(c, mygoui.GroupOptions{Label: "History"},
    mygoui.GroupItem{Key: "back", Build: func(corners [4]float32) ui.Element {
        return mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Back", Variant: mygoui.ButtonOutline, Corners: &corners,
            OnClick: back,
        })
    }},
    mygoui.GroupItem{Key: "next", Build: func(corners [4]float32) ui.Element {
        return mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Next", Variant: mygoui.ButtonOutline, Corners: &corners,
            OnClick: next,
        })
    }},
)`,
		Params: []ComponentParam{
			{"Orientation", "mygoui.Orientation", "Horizontal", "Join axis."},
			{"FullWidth", "bool", "false", "Distribute controls across the available width."},
			{"Label", "string", "", "Accessible group name."},
			{"GroupItem.Key", "string", "required", "Stable identity per control."},
			{"GroupItem.Build", "func([4]float32) ui.Element", "required", "Receives top-left, top-right, bottom-right, bottom-left corners."},
		},
		Tokens: []string{"Primary", "PrimaryForeground", "Input", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "scrollarea", Title: "ScrollArea", Group: "Foundations", Order: 101,
		Description: "A native scroll viewport with optional edge treatment.",
		Examples: []ComponentExample{
			{Title: "Vertical", Demo: func(c *ui.Context, s *State) {
				scrollSpecimen(c, &cossStateOf(s, "scrollarea").Scroll[0], mygoui.ScrollVertical, nil, "Activity")
			}},
			{Title: "Horizontal with shadow", Demo: func(c *ui.Context, s *State) {
				scrollSpecimen(c, &cossStateOf(s, "scrollarea").Scroll[1], mygoui.ScrollHorizontal,
					&mygoui.ScrollFadeOptions{Mode: mygoui.ShadowFade, Edges: mygoui.FadeLeft | mygoui.FadeRight},
					"Releases")
			}},
			{Title: "Both axes", Demo: func(c *ui.Context, s *State) {
				scrollSpecimen(c, &cossStateOf(s, "scrollarea").Scroll[2], mygoui.ScrollBoth, nil, "Build matrix")
			}},
		},
		Code: `mygoui.ScrollArea(c, &scroll, mygoui.ScrollAreaOptions{
    Axis: mygoui.ScrollVertical, Height: 180, Label: "Activity",
    Fade: &mygoui.ScrollFadeOptions{
        Mode: mygoui.ShadowFade, Edges: mygoui.FadeTop | mygoui.FadeBottom,
    },
}, func() {
    for i := range 20 {
        ui.Text(c, fmt.Sprintf("Activity %02d", i+1))
    }
})`,
		Params: []ComponentParam{
			{"state", "*ui.ScrollState", "required", "Caller-owned scroll offsets and bounds."},
			{"Axis", "mygoui.ScrollAxis", "ScrollVertical", "Vertical, horizontal or both axes."},
			{"Height", "float32", "0", "Viewport height; zero follows the parent."},
			{"Label", "string", "", "Accessible viewport name."},
			{"Fade", "*mygoui.ScrollFadeOptions", "nil", "Optional shadow or surface edge fade."},
		},
		Tokens: []string{"Background", "Foreground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "scrollfade", Title: "ScrollFade", Group: "Foundations", Order: 102,
		Description: "Native alpha masking, shadow or surface blending at viewport edges.",
		Examples: []ComponentExample{
			{Title: "MaskFade over a gradient", Demo: func(c *ui.Context, s *State) {
				maskSpecimen(c, cossStateOf(s, "scrollfade"))
			}},
			{Title: "ShadowFade", Demo: func(c *ui.Context, s *State) {
				scrollSpecimen(c, &cossStateOf(s, "scrollfade").Scroll[0], mygoui.ScrollVertical,
					&mygoui.ScrollFadeOptions{Mode: mygoui.ShadowFade, Size: 32}, "Shadow feed")
			}},
			{Title: "SurfaceFade, bottom edge", Demo: func(c *ui.Context, s *State) {
				scrollSpecimen(c, &cossStateOf(s, "scrollfade").Scroll[1], mygoui.ScrollVertical,
					&mygoui.ScrollFadeOptions{
						Mode: mygoui.SurfaceFade, Edges: mygoui.FadeBottom, Size: 40,
						Surface: theme.Colors(c).Background,
					}, "Surface feed")
			}},
			{Title: "Custom horizontal shadow", Demo: func(c *ui.Context, s *State) {
				scrollSpecimen(c, &cossStateOf(s, "scrollfade").Scroll[2], mygoui.ScrollHorizontal,
					&mygoui.ScrollFadeOptions{
						Mode: mygoui.ShadowFade, Edges: mygoui.FadeRight, Size: 48,
						Color: theme.Colors(c).Primary.Alpha(0.3),
					}, "Horizontal feed")
			}},
		},
		Code: `viewport := ui.Scroll(c).TrackScroll(&scroll).Height(180).Children(func() {
    for i := range 20 {
        ui.Text(c, fmt.Sprintf("Activity %02d", i+1))
    }
})
mygoui.ScrollFade(c, viewport, &scroll, mygoui.ScrollFadeOptions{
    Mode: mygoui.MaskFade,
    Edges: mygoui.FadeTop | mygoui.FadeBottom, Size: 32,
})`,
		Params: []ComponentParam{
			{"viewport", "ui.Element", "required", "Current-build native scroll viewport."},
			{"state", "*ui.ScrollState", "required", "The same state the viewport tracks."},
			{"Mode", "mygoui.ScrollFadeMode", "ShadowFade", "MaskFade masks the native subtree; SurfaceFade blends into an opaque surface."},
			{"Edges", "mygoui.ScrollFadeEdges", "all", "Physical top, right, bottom and left bit flags."},
			{"Size", "float32", "24", "Fade depth in DIPs."},
			{"Surface / Color", "ui.Color", "theme", "SurfaceFade needs an opaque surface; ShadowFade accepts a color. MaskFade needs neither."},
		},
		Tokens: []string{"Background", "Foreground", "Primary"},
	})

	RegisterComponent(ComponentPage{
		ID: "field", Title: "Field", Group: "Forms", Order: 4,
		Description: "A label, native editor, description and validation message.",
		Examples: []ComponentExample{
			{Title: "Label and description", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "field")
				form.Field(c, form.FieldOptions{
					Label: "Repository", Description: "Lowercase letters and dashes.", Required: true,
				}, func() {
					form.Input(c, &st.Name, form.InputOptions{AriaLabel: "Repository", FullWidth: true})
				})
			}},
			{Title: "Validation", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "field")
				form.Field(c, form.FieldOptions{
					Label: "Workspace slug", Invalid: true, Error: "This slug is already taken.",
				}, func() {
					form.Input(c, &st.Draft, form.InputOptions{AriaLabel: "Workspace slug", Invalid: true, FullWidth: true})
				})
			}},
		},
		Code: `form.Field(c, form.FieldOptions{
    Label: "Repository", Description: "Lowercase letters and dashes.", Required: true,
}, func() {
    form.Input(c, &name, form.InputOptions{AriaLabel: "Repository", FullWidth: true})
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Associated native field label."},
			{"Description / Error", "string", "", "Supporting text or validation error."},
			{"Required / Invalid / Disabled", "bool", "false", "Native field state."},
		},
		Tokens: []string{"Foreground", "MutedForeground", "DestructiveForeground", "Input", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "command", Title: "Command", Group: "Navigation", Order: 11,
		Description: "Search and execute caller-owned actions with keyboard or pointer.",
		Examples: []ComponentExample{
			{Title: "Inline command", Demo: func(c *ui.Context, s *State) {
				commandSpecimen(c, cossStateOf(s, "command"), 0)
			}},
			{Title: "Custom empty state", Demo: func(c *ui.Context, s *State) {
				commandSpecimen(c, cossStateOf(s, "command"), 1)
			}},
		},
		Code: `navigation.Command(c, navigation.CommandOptions{
    State: &commands, Label: "Search actions", Placeholder: "Type a command",
    Items: []navigation.CommandItem{
        {Key: "create", Label: "Create project", Keywords: "new", Shortcut: "N"},
        {Key: "archive", Label: "Archive project", Disabled: true},
    },
    OnSelect: func(key string) { last = key },
})`,
		Params: []ComponentParam{
			{"State", "*navigation.CommandState", "required", "Query, focus handle and native list highlight."},
			{"Items", "[]navigation.CommandItem", "nil", "Key, label, keywords, shortcut and disabled state."},
			{"Label / Placeholder / Empty", "string", "", "Accessible search name, prompt and empty result."},
			{"MaxHeight", "float32", "320", "List viewport height."},
			{"OnSelect", "func(string)", "nil", "Receives the activated command key."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Accent", "MutedForeground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "contextmenu", Title: "ContextMenu", Group: "Navigation", Order: 12,
		Description: "A native secondary-click menu with keyboard access and nested actions.",
		Examples: []ComponentExample{
			{Title: "Secondary click", Demo: func(c *ui.Context, s *State) {
				contextSpecimen(c, cossStateOf(s, "contextmenu"), false)
			}},
			{Title: "Disabled and nested actions", Demo: func(c *ui.Context, s *State) {
				contextSpecimen(c, cossStateOf(s, "contextmenu"), true)
			}},
			{Title: "Themed menu", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "contextmenu-themed")
				anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Themed document actions", Variant: mygoui.ButtonOutline})
				navigation.PopoverMenu(anchor, navigation.MenuOptions{
					Selection: navigation.MenuSelectSingle, Selected: &st.Selected,
					Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
						{Key: "name", Label: "Sort by name"},
						{Key: "date", Label: "Sort by date"},
						{Label: "Delete", Disabled: true},
					}}},
				})
				ui.Text(c, "Sort: "+strings.Join(st.Selected, ", "))
			}},
		},
		Code: `anchor := mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Document actions", Variant: mygoui.ButtonOutline,
})
navigation.ContextMenu(anchor, navigation.MenuOptions{
    Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
        {Key: "rename", Label: "Rename"},
        {Label: "Delete", Disabled: true},
    }}},
    OnSelect: func(key string) { last = key },
})`,
		Params: []ComponentParam{
			{"anchor", "ui.Element", "required", "Current-build target for secondary click or the menu key."},
			{"Sections", "[]navigation.MenuSectionOptions", "nil", "Menu items, separators and submenus; PopoverMenu draws them in the current theme."},
			{"Selection / Selected", "MenuSelection / *[]string", "none", "Optional checked item state."},
			{"OnSelect", "func(string)", "nil", "Receives the activated item key."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Accent", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "segmentedcontrol", Title: "SegmentedControl", Group: "Navigation", Order: 13,
		Description: "A single-choice native control with roving keyboard focus.",
		Examples: []ComponentExample{
			{Title: "Default", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "segment-billing")
				ui.Row(c).FillWidth().Justify(ui.Center).Children(func() {
					navigation.SegmentedControl(c, navigation.SegmentedControlOptions{
						Label: "Billing period", Selected: &st.Choice,
						Items: []navigation.TabOptions{{Label: "Monthly"}, {Label: "Yearly"}},
					})
				})
			}},
			{Title: "Content view", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "segmentedcontrol")
				navigation.SegmentedControl(c, navigation.SegmentedControlOptions{
					Label: "Content view", Selected: &st.Choice,
					Items: []navigation.TabOptions{{Label: "List view"}, {Label: "Grid view"}, {Label: "Map", Disabled: true}},
				})
				label := "List"
				if st.Choice == 1 {
					label = "Grid"
				}
				ui.Text(c, "Current view: "+label)
			}},
			{Title: "Custom choices", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "segment-density")
				navigation.SegmentedControl(c, navigation.SegmentedControlOptions{
					Label: "Density", Selected: &st.Choice,
					Items: []navigation.TabOptions{{Label: "Compact"}, {Label: "Comfortable"}},
				}).FillWidth()
			}},
		},
		Code: `navigation.SegmentedControl(c, navigation.SegmentedControlOptions{
    Label: "Content view", Selected: &selected,
    Items: []navigation.TabOptions{
        {Label: "List"}, {Label: "Grid"}, {Label: "Map", Disabled: true},
    },
})`,
		Params: []ComponentParam{
			{"Selected", "*int", "required", "Caller-owned selected index."},
			{"Items", "[]navigation.TabOptions", "nil", "Labels and disabled states; content is composed after the control."},
			{"Label", "string", "", "Accessible group name."},
			{"OnChange", "func(int)", "nil", "Receives the selected index."},
		},
		Tokens: []string{"MutedBackground", "Background", "Foreground", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "previewcard", Title: "PreviewCard", Group: "Overlays", Order: 81,
		Description: "A rich preview that stays open as pointer or focus enters its content.",
		Examples: []ComponentExample{
			{Title: "Profile preview", Demo: func(c *ui.Context, s *State) {
				previewSpecimen(c, cossStateOf(s, "previewcard"), 0)
			}},
			{Title: "Custom width and placement", Demo: func(c *ui.Context, s *State) {
				previewSpecimen(c, cossStateOf(s, "previewcard"), 1)
			}},
		},
		Code: `anchor := mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Ada Lovelace", Variant: mygoui.ButtonLink,
})
overlay.PreviewCard(c, anchor, &preview, overlay.PreviewCardOptions{
    Label: "Profile preview", Width: 256,
    Placement: mygoui.Placement{Side: mygoui.SideBottom},
    Body: func() {
        mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Ada Lovelace"})
        ui.Text(c, "Ada Lovelace")
    },
})`,
		Params: []ComponentParam{
			{"state", "*overlay.PreviewCardState", "required", "Keep one state per anchor and window."},
			{"Placement", "mygoui.Placement", "auto", "Anchor side and alignment."},
			{"Width", "float32", "256", "Panel width."},
			{"Label / Body", "string / func()", "", "Accessible preview name and interactive content."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "sheet", Title: "Sheet", Group: "Overlays", Order: 82,
		Description: "An edge-aligned dialog with caller-owned open state.",
		Examples: []ComponentExample{
			{Title: "Right sheet", Demo: func(c *ui.Context, s *State) {
				sheetSpecimen(c, cossStateOf(s, "sheet"), 0, mygoui.SideRight)
			}},
			{Title: "Left sheet", Demo: func(c *ui.Context, s *State) {
				sheetSpecimen(c, cossStateOf(s, "sheet"), 1, mygoui.SideLeft)
			}},
		},
		Code: `mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Open settings", OnClick: func() { open = true },
})
overlay.Sheet(c, &open, overlay.SheetOptions{
    Side: mygoui.SideRight, Width: 320, Title: "Settings",
    Body: func() {
        form.Switch(c, form.SwitchOptions{Label: "Notifications", On: &enabled})
    },
    Footer: func() {
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Done", OnClick: func() { open = false },
        })
    },
})`,
		Params: []ComponentParam{
			{"open", "*bool", "required", "Caller-owned visibility."},
			{"Side", "mygoui.Side", "SideRight", "Window edge."},
			{"Width", "float32", "dialog default", "Panel width in DIPs."},
			{"Title / Description", "string", "", "Dialog heading and supporting text."},
			{"Header / Body / Footer", "func()", "nil", "Composed panel regions."},
			{"PreventDismiss", "bool", "false", "Requires an explicit close action."},
		},
		Tokens: []string{"Background", "Foreground", "Border", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "composer", Title: "Composer", Group: "Forms", Order: 11,
		Description: "A growing native draft editor with persistent actions and submission.",
		Examples: []ComponentExample{
			{Title: "Message", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "composer")
				form.Composer(c, &st.Draft, form.ComposerOptions{
					Label: "Message draft", Placeholder: "Write a message", SubmitLabel: "Send message",
					OnSubmit: func(text string) { st.Last = text; st.Count++; st.Draft = "" },
				})
				ui.Text(c, fmt.Sprintf("Sent %d: %s", st.Count, st.Last))
			}},
			{Title: "Custom actions", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "composer-custom")
				form.Composer(c, &st.Draft, form.ComposerOptions{
					Label: "Comment draft", Placeholder: "Add a comment", SubmitLabel: "Post comment",
					MinLines: 3, MaxLines: 5, DisableEnterSubmit: true,
					Actions: func() {
						mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
							Label: "Include context", Icon: lucide.Paperclip(), Selected: &st.Checked,
						})
					},
					OnSubmit: func(text string) { st.Last = text; st.Draft = "" },
				})
				ui.Text(c, "Posted: "+st.Last)
			}},
			{Title: "Read only", Demo: func(c *ui.Context, s *State) {
				st := cossStateOf(s, "composer-readonly")
				st.Draft = "This conversation is archived."
				form.Composer(c, &st.Draft, form.ComposerOptions{Label: "Archived draft", ReadOnly: true})
			}},
		},
		Code: `form.Composer(c, &draft, form.ComposerOptions{
    Label: "Message", Placeholder: "Write a message", SubmitLabel: "Send",
    MinLines: 2, MaxLines: 6,
    OnSubmit: func(text string) {
        sent = text
        draft = ""
    },
})`,
		Params: []ComponentParam{
			{"draft", "*string", "required", "Caller-owned text; submission only clears it when the caller does."},
			{"OnSubmit", "func(string)", "nil", "Enter or submit action; Shift+Enter inserts a newline."},
			{"MinLines / MaxLines", "int", "2 / 8", "Editor growth bounds."},
			{"Actions", "func()", "nil", "Persistent controls beside submit."},
			{"Disabled / ReadOnly / SubmitDisabled", "bool", "false", "Editing and submission states."},
			{"AllowEmpty / DisableEnterSubmit", "bool", "false", "Submission policy."},
		},
		Tokens: []string{"Input", "Background", "Primary", "PrimaryForeground", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "sidebar", Title: "Sidebar", Group: "Navigation", Order: 14,
		Description: "Native navigation with grouped actions, collapse state and header/footer slots.",
		Examples: []ComponentExample{
			{Title: "Workspace navigation", Demo: func(c *ui.Context, s *State) {
				sidebarSpecimen(c, cossStateOf(s, "sidebar"), false)
			}},
			{Title: "Compact rail", Demo: func(c *ui.Context, s *State) {
				sidebarSpecimen(c, cossStateOf(s, "sidebar-rail"), true)
			}},
		},
		Code: `mygoui.Sidebar(c, mygoui.SidebarOptions{
    Selected: &destination, Collapsed: &collapsed, Width: 220,
    Label: "Workspace navigation",
    Groups: []mygoui.SidebarGroupOptions{{
        Key: "workspace", Label: "Workspace",
        Items: []mygoui.SidebarItemOptions{
            {ID: "home", Label: "Home", Icon: lucide.House()},
            {ID: "settings", Label: "Settings", Icon: lucide.Settings()},
            {ID: "locked", Label: "Archived", Icon: lucide.Lock(), Disabled: true},
        },
    }},
    OnSelect: func(id string) { last = id },
})`,
		Params: []ComponentParam{
			{"Selected / Collapsed", "*string / *bool", "required", "Caller-owned destination and collapse state."},
			{"Groups", "[]mygoui.SidebarGroupOptions", "nil", "Stable group keys and native navigation actions."},
			{"Width / CollapsedWidth", "float32", "256 / 48", "Expanded and compact rail widths."},
			{"Header / Footer", "func()", "nil", "Slots can adapt to the caller's collapsed state."},
			{"OnSelect", "func(string)", "nil", "Receives the activated destination."},
		},
		Tokens: []string{
			"Sidebar", "SidebarForeground", "SidebarPrimary", "SidebarPrimaryForeground",
			"SidebarAccent", "SidebarAccentForeground", "SidebarBorder", "SidebarRing",
		},
	})

	RegisterComponent(ComponentPage{
		ID: "theming", Title: "Theming", Group: "Foundations", Order: 2,
		Description: "Coss semantic roles, light and dark appearances, primary and radius overrides.",
		Examples: []ComponentExample{
			{Title: "Live theme", Demo: overviewThemeControls},
			{Title: "Semantic roles", Demo: func(c *ui.Context, s *State) {
				ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Primary action"})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Secondary", Variant: mygoui.ButtonSecondary})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Accent on hover", Variant: mygoui.ButtonGhost})
					mygoui.Badge(c, mygoui.BadgeOptions{Label: "Ready", Variant: mygoui.BadgeSuccess})
				})
			}},
		},
		Code: `primary := ui.RGB(37, 99, 235)
theme.Configure(theme.Preferences{
    Radius: 10,
    Palette: theme.PaletteOverrides{
        Light: &theme.PalettePatch{Primary: &primary},
        Dark: &theme.PalettePatch{Primary: &primary},
    },
})
theme.Apply(c)`,
		Params: []ComponentParam{
			{"Radius", "float32", "10", "Base radius; controls derive their own scale."},
			{"Accent", "ui.Color", "neutral primary", "Host shortcut for overriding the primary action color."},
			{"Palette", "theme.PaletteOverrides", "Coss defaults", "Independent light and dark semantic patches."},
			{"FontFamily / HeadingFontFamily / MonoFontFamily", "string", "native", "Host typography."},
		},
		Tokens: []string{
			"Background", "Foreground", "Card", "CardForeground", "Popover", "PopoverForeground",
			"Primary", "PrimaryForeground", "Secondary", "SecondaryForeground", "MutedBackground",
			"MutedForeground", "Accent", "AccentForeground", "Destructive", "DestructiveForeground",
			"Info", "InfoForeground", "Success", "SuccessForeground", "Warning", "WarningForeground",
			"Border", "Input", "Ring", "Chart1", "Chart2", "Chart3", "Chart4", "Chart5",
			"Sidebar", "SidebarForeground", "SidebarPrimary", "SidebarPrimaryForeground",
			"SidebarAccent", "SidebarAccentForeground", "SidebarBorder", "SidebarRing",
			"Code", "CodeForeground", "CodeHighlight",
		},
	})
}

func scrollSpecimen(c *ui.Context, state *ui.ScrollState, axis mygoui.ScrollAxis, fade *mygoui.ScrollFadeOptions, label string) {
	t := theme.Colors(c)
	mygoui.ScrollArea(c, state, mygoui.ScrollAreaOptions{
		Axis: axis, Height: 172, Label: label, Fade: fade,
	}, func() {
		if axis == mygoui.ScrollHorizontal {
			ui.Row(c).Gap(theme.Space12).Padding(theme.Space12).Children(func() {
				for i := range 8 {
					ui.Column(c).Width(140).Height(120).Padding(theme.Space12).
						Radius(theme.RadiusScale().LG).Background(t.Secondary).Children(func() {
						ui.Text(c, fmt.Sprintf("Release %02d", i+1))
						ui.Text(c, "Ready to deploy")
					})
				}
			})
			return
		}
		e := ui.Column(c).Gap(theme.Space8).Padding(theme.Space12)
		if axis == mygoui.ScrollBoth {
			e.Width(840)
		}
		e.Children(func() {
			for i := range 24 {
				ui.Text(c, fmt.Sprintf("Build %02d     main     Passed     38s", i+1)).SingleLine()
				mygoui.Separator(c, mygoui.SeparatorOptions{})
			}
		})
	}).FillWidth().Background(t.Background).Border(theme.Hairline, t.Border).Radius(theme.RadiusScale().LG)
}

func maskSpecimen(c *ui.Context, st *cossState) {
	t := theme.Colors(c)
	from, to := t.InfoSoft.Over(t.Background), t.SuccessSoft.Over(t.Background)
	ui.Column(c).FillWidth().MinWidth(0).Radius(theme.RadiusScale().LG).Clip().
		Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillGradient(r, ui.LinearGradient{From: from, To: to, Angle: 90}, 0)
		}).Children(func() {
		mygoui.ScrollArea(c, &st.Scroll[3], mygoui.ScrollAreaOptions{
			Height: 192, Label: "Masked activity",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.MaskFade, Edges: mygoui.FadeTop | mygoui.FadeBottom, Size: 40,
			},
		}, func() {
			ui.Column(c).FillWidth().Padding(theme.Space12).Gap(theme.Space8).Children(func() {
				for i := range 16 {
					mygoui.Button(c.Key(i), mygoui.ButtonOptions{
						Label: fmt.Sprintf("Inspect build %02d", i+1), Variant: mygoui.ButtonOutline,
						FullWidth: true, OnClick: func() { st.Count = i + 1 },
					})
				}
			})
		}).FillWidth()
	})
	ui.Text(c, fmt.Sprintf("Inspected build %02d", st.Count))
}

func commandSpecimen(c *ui.Context, st *cossState, index int) {
	label, empty := "Search actions", "No matching actions"
	items := []navigation.CommandItem{
		{Key: "create", Label: "Create project", Keywords: "new repository", Shortcut: "N"},
		{Key: "theme", Label: "Change appearance", Keywords: "dark light theme"},
		{Key: "archive", Label: "Archive project", Disabled: true},
	}
	if index == 1 {
		label, empty = "Search files", "No matching files"
		items = []navigation.CommandItem{{Key: "readme", Label: "README.md"}, {Key: "main", Label: "main.go"}}
	}
	navigation.Command(c, navigation.CommandOptions{
		Label: label, Placeholder: "Search", Empty: empty, Items: items,
		State: &st.Command[index], MaxHeight: 144,
		OnSelect: func(key string) { st.Last = key; st.Count++ },
	})
	ui.Text(c, "Executed: "+st.Last)
}

func contextSpecimen(c *ui.Context, st *cossState, nested bool) {
	label := "Document actions"
	items := []navigation.MenuItemOptions{{Key: "rename", Label: "Rename"}, {Key: "duplicate", Label: "Duplicate"}}
	if nested {
		label = "Protected document"
		items = []navigation.MenuItemOptions{
			{Label: "Move", Submenu: []navigation.MenuItemOptions{{Key: "archive", Label: "Archive folder"}}},
			{Separator: true}, {Label: "Delete", Disabled: true},
		}
	}
	anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: label, Variant: mygoui.ButtonOutline})
	navigation.ContextMenu(anchor, navigation.MenuOptions{
		Sections: []navigation.MenuSectionOptions{{Items: items}},
		OnSelect: func(key string) { st.Last = key },
	})
	ui.Text(c, "Action: "+st.Last)
}

func previewSpecimen(c *ui.Context, st *cossState, index int) {
	label, width := "Ada Lovelace", float32(256)
	side := mygoui.SideBottom
	if index == 1 {
		label, width, side = "Release notes", 288, mygoui.SideTop
	}
	anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: label, Variant: mygoui.ButtonLink})
	overlay.PreviewCard(c, anchor, &st.Preview[index], overlay.PreviewCardOptions{
		Label: label + " preview", Width: width, Placement: mygoui.Placement{Side: side},
		Body: func() {
			mygoui.Avatar(c, mygoui.AvatarOptions{Name: label})
			mygoui.Typography(c, label, mygoui.TypographyOptions{Type: theme.TypeH5})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Follow updates", Variant: mygoui.ButtonOutline,
				OnClick: func() { st.PreviewHits++ },
			})
		},
	})
	ui.Text(c, fmt.Sprintf("Following %d", st.PreviewHits))
}

func sheetSpecimen(c *ui.Context, st *cossState, index int, side mygoui.Side) {
	label := "Open right sheet"
	if side == mygoui.SideLeft {
		label = "Open left sheet"
	}
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Variant: mygoui.ButtonOutline, OnClick: func() { st.Open[index] = true },
	})
	overlay.Sheet(c, &st.Open[index], overlay.SheetOptions{
		Side: side, Width: 320, Title: "Workspace settings",
		Body: func() {
			form.Input(c, &st.Name, form.InputOptions{Label: "Workspace name", FullWidth: true})
			form.Switch(c, form.SwitchOptions{Label: "Notifications", On: &st.Checked})
		},
		Footer: func() {
			w, _ := c.Size()
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Save settings", FullWidth: w < 640, OnClick: func() { st.Count++; st.Open[index] = false },
			})
		},
	})
	ui.Text(c, fmt.Sprintf("Changes saved %d", st.Count))
}

func sidebarSpecimen(c *ui.Context, st *cossState, rail bool) {
	if st.Name == "" {
		st.Name, st.Checked = "home", rail
	}
	ui.Row(c).FillWidth().Height(300).Gap(theme.Space12).Children(func() {
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &st.Name, Collapsed: &st.Checked, Width: 196,
			Label: "Workspace navigation", ToggleLabel: "Collapse workspace",
			Header: func() {
				if !st.Checked {
					mygoui.Typography(c, "Acme", mygoui.TypographyOptions{Type: theme.TypeH5})
				}
			},
			Footer: func() {
				if !st.Checked {
					mygoui.Typography(c, "Personal workspace", mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
				}
			},
			Groups: []mygoui.SidebarGroupOptions{{
				Key: "workspace", Label: "Workspace",
				Items: []mygoui.SidebarItemOptions{
					{ID: "home", Label: "Home", Icon: lucide.House()},
					{ID: "projects", Label: "Projects", Icon: lucide.Folder()},
					{ID: "settings", Label: "Settings", Icon: lucide.Settings()},
					{ID: "locked", Label: "Archived", Icon: lucide.Lock(), Disabled: true},
				},
			}},
			OnSelect: func(id string) { st.Last = id },
		})
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Padding(theme.Space8).Children(func() {
			ui.Text(c, "Destination: "+st.Name)
		})
	})
}
