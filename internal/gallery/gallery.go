// Package gallery is a native component workbench with local sample fixtures.
package gallery

import (
	"strings"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/codeblock"
	"github.com/Cyenoch/mygoui/internal/gallery/locales"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/Cyenoch/mygoui/trace"
	"github.com/egoist/mygo/ui"
)

// State retains navigation and independent preview state across builds.
type State struct {
	Search                                   string
	Selected                                 ComponentID
	Dark, LargeText, BlueTheme, HighContrast bool
	Basic                                    BasicState
	Code                                     CodeState
	Trace                                    TraceState
	Forms                                    FormsState
	Overlays                                 OverlayState
	Conversation                             ConversationState
	previewScroll                            ui.ScrollState
	AppearanceOpen, SourceOpen               bool
	DocumentationCode                        codeblock.State
	JumpSection                              string
	DemoBasic                                BasicState
	DemoForms                                FormsState
	DemoOverlay                              OverlayState
	DemoCode                                 codeblock.State
	DemoTrace                                trace.State
	DemoConversation                         ConversationState
	Dates, DemoDates                         DateChartState
	Navigation, DemoNavigation               NavigationState
	DataLayout, DemoDataLayout               DataLayoutState
}

func New() *State {
	return &State{
		Selected: "code-block", Dark: true,
		Basic: NewBasicState(), Code: NewCodeState(), Trace: NewTraceState(), Forms: NewFormsState(),
		Overlays: NewOverlayState(), Conversation: NewConversationState(),
		DemoBasic: NewBasicState(), DemoForms: NewFormsState(), DemoOverlay: NewOverlayState(), DemoConversation: NewConversationState(),
		Dates: NewDateChartState(), DemoDates: NewDateChartState(),
		Navigation: NewNavigationState(), DemoNavigation: NewNavigationState(),
		DataLayout: NewDataLayoutState(), DemoDataLayout: NewDataLayoutState(),
	}
}

// View applies the existing theme once and constructs only the selected preview.
func View(c *ui.Context, s *State) {
	base := *c.Theme()
	base.Dark = s.Dark
	c.SetTheme(&base)
	preset := theme.Neutral
	if s.BlueTheme {
		preset = theme.Blue
	}
	if s.HighContrast {
		theme.ApplyHighContrast(c, preset)
	} else {
		theme.ApplyPreset(c, preset)
	}
	colors := theme.Colors(c)
	width, _ := c.Size()
	compact := width < theme.Scale(c, 960)
	ui.Column(c).Fill().MinWidth(0).Clip().Background(colors.Background).TextColor(colors.Foreground).Children(func() {
		ui.Column(c).Label(locales.ComponentGalleryToolbar()).FillWidth().Shrink(0).MinWidth(0).
			Padding(12, 16).Gap(8).BorderWidth(0, 0, 1, 0).BorderColor(colors.Border).Children(func() {
			ui.Row(c).FillWidth().Wrap().AlignContent(ui.Start).Gap(8).Children(func() {
				ui.Text(c, locales.ComponentGalleryBrand()).Bold().FontSize(theme.Scale(c, 14)).AlignSelf(ui.Center)
				ui.TextInput(c.Key("search"), &s.Search).Label(locales.ComponentGallerySearch()).Placeholder(locales.ComponentGallerySearchPlaceholder()).
					Grow(1).MinWidth(140).MaxWidth(320)
				mygoui.Button(c.Key("appearance"), mygoui.ButtonOptions{
					Label: appearanceLabel(s.Dark), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonSM,
					OnClick: func() { s.Dark = !s.Dark },
				})
				mygoui.Button(c.Key("display-options"), mygoui.ButtonOptions{
					Label: locales.ComponentGalleryDisplayOptions(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonSM,
					OnClick: func() { s.AppearanceOpen = !s.AppearanceOpen },
				})
			})
			if s.AppearanceOpen {
				ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
					mygoui.Toggle(c.Key("app-theme"), &s.BlueTheme, mygoui.ToggleOptions{Label: locales.ComponentGalleryBlueTheme()})
					mygoui.Toggle(c.Key("high-contrast"), &s.HighContrast, mygoui.ToggleOptions{Label: locales.ComponentGalleryHighContrast()})
					mygoui.Toggle(c.Key("large-text"), &s.LargeText, mygoui.ToggleOptions{
						Label: locales.ComponentGalleryLargeText(), OnChange: func(large bool) {
							size := float32(16)
							if large {
								size = 32
							}
							theme.SetPreferences(theme.Preferences{FontSize: size})
						},
					})
				})
			}
		})
		// The body owns only the remaining height. Scroll viewports cannot use their
		// specimen's intrinsic height to expand or center across the toolbar band.
		ui.Column(c).FillWidth().Grow(1).Basis(0).MinHeight(0).MinWidth(0).Clip().Children(func() {
			if compact {
				ui.ScrollHorizontal(c.Key("compact-navigation")).Label(locales.ComponentGalleryNavigation()).
					FillWidth().Height(max(44, theme.Scale(c, 44))).Shrink(0).Padding(6, 8).Gap(6).
					Children(func() { catalogNavigation(c, s, true) })
			}
			ui.Row(c).FillWidth().Grow(1).Basis(0).MinHeight(0).MinWidth(0).
				AlignItems(ui.Stretch).Clip().Children(func() {
				if !compact {
					ui.Scroll(c.Key("navigation")).Label(locales.ComponentGalleryNavigation()).Width(224).Shrink(0).FillHeight().
						Padding(12).Gap(4).Background(colors.Sidebar).
						BorderWidth(0, 1, 0, 0).BorderColor(colors.SidebarBorder).
						Children(func() { catalogNavigation(c, s, false) })
				}
				preview(c, s)
				if width >= theme.Scale(c, 1280) {
					outline(c, s)
				}
			})
		})
	})
	overlay.Toaster(c)
}

func appearanceLabel(dark bool) string {
	if dark {
		return locales.ComponentGalleryUseLightTheme()
	}
	return locales.ComponentGalleryUseDarkTheme()
}

func catalogNavigation(c *ui.Context, s *State, compact bool) {
	query := strings.ToLower(strings.TrimSpace(s.Search))
	matches := 0
	for _, entry := range Catalog {
		name := locales.ComponentGalleryName(strings.ReplaceAll(string(entry.ID), "-", "_"))
		family := locales.ComponentGalleryFamily(strings.ReplaceAll(entry.Family, " ", "_"))
		if query != "" && !strings.Contains(strings.ToLower(name+" "+family+" "+entry.Name+" "+entry.Family), query) {
			continue
		}
		matches++
		variant := mygoui.ButtonGhost
		if entry.ID == s.Selected {
			variant = mygoui.ButtonSecondary
		}
		id := entry.ID
		button := mygoui.Button(c.Key("nav-"+id), mygoui.ButtonOptions{
			Label: name, Variant: variant, Size: mygoui.ButtonSM,
			OnClick: func() { selectComponent(s, id) },
		})
		button.Justify(ui.Start).AlignItems(ui.Center).Shrink(0)
		if compact {
			button.MinHeight(44)
		}
		if !compact {
			button.FillWidth()
		}
	}
	if matches == 0 {
		ui.Text(c, locales.ComponentGalleryNoMatches()).TextColor(theme.Colors(c).MutedForeground)
		mygoui.Button(c, mygoui.ButtonOptions{Label: locales.ComponentGalleryClearSearch(), Variant: mygoui.ButtonOutline,
			OnClick: func() { s.Search = "" }})
	}
}
