package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type foundationsState struct {
	LinkClicks   int
	ActionClicks int
	CardClicks   int
	ListChoice   string
	Scroll       ui.ScrollState
	ScrollX      ui.ScrollState
}

// typographyHeadingsBody renders the six semantic heading levels.
func typographyHeadingsBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space4).Children(func() {
		for level := 1; level <= 6; level++ {
			mygoui.Typography(c, fmt.Sprintf("Heading %d", level), mygoui.TypographyOptions{
				Type: theme.SemanticType(int(theme.TypeH1) + level - 1),
			})
		}
	})
}

// typographyBodyBody renders the reading-length body roles.
func typographyBodyBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space4).Children(func() {
		mygoui.Typography(c, "Body copy uses the default 16/24 role for reading-length text.",
			mygoui.TypographyOptions{Type: theme.TypeBody})
		mygoui.Typography(c, "Body SM tightens supporting copy to 14/20.",
			mygoui.TypographyOptions{Type: theme.TypeBodySM})
		mygoui.Typography(c, "Body XS is reserved for captions and metadata.",
			mygoui.TypographyOptions{Type: theme.TypeBodyXS})
	})
}

// typographyCodeProseBody renders inline code and grouped prose.
func typographyCodeProseBody(c *ui.Context) {
	mygoui.Typography(c, "go test ./...", mygoui.TypographyOptions{
		Type: theme.TypeCode, Selectable: true,
	})
	mygoui.Prose(c, func() {
		mygoui.Paragraph(c, "Prose groups paragraphs at a comfortable rhythm. Each block keeps its own "+
			"semantic role so assistive technology can navigate the outline.", false)
		mygoui.Paragraph(c, "Inline code such as mygoui.Typography stays selectable inside flowing text.", true)
		mygoui.Code(c, "go build ./internal/gallery ./examples/gallery")
	})
}

// surfaceVariantsBody renders the four surface elevation tokens.
func surfaceVariantsBody(c *ui.Context) {
	variants := []struct {
		name    string
		variant mygoui.SurfaceVariant
		detail  string
	}{
		{"Surface Default", mygoui.SurfaceDefault, "Primary content plane."},
		{"Surface Secondary", mygoui.SurfaceSecondary, "Nested groups and quiet panels."},
		{"Surface Tertiary", mygoui.SurfaceTertiary, "The deepest neutral plane."},
		{"Surface Transparent", mygoui.SurfaceTransparent, "Inherits the surrounding paint."},
	}
	ui.Column(c).FillWidth().Gap(theme.Space8).Children(func() {
		for _, item := range variants {
			mygoui.Surface(c, mygoui.SurfaceOptions{
				Variant: item.variant, Padding: theme.Space12, Radius: theme.RadiusScale().XL,
				FillWidth: true, Label: item.name,
			}, func() {
				mygoui.Typography(c, item.name, mygoui.TypographyOptions{Type: theme.TypeBodySM})
				mygoui.Typography(c, item.detail, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
			})
		}
	})
}

// surfacePaddingBody renders the padding and radius options.
func surfacePaddingBody(c *ui.Context) {
	ui.Row(c).Wrap().Gap(theme.Space12).Children(func() {
		for _, item := range []struct {
			label   string
			padding float32
			radius  float32
		}{
			{"Compact", theme.Space8, theme.RadiusScale().MD},
			{"Regular", theme.Space16, theme.RadiusScale().XL},
			{"Roomy", theme.Space24, theme.RadiusScale().XXXL},
		} {
			mygoui.Surface(c, mygoui.SurfaceOptions{
				Variant: mygoui.SurfaceSecondary, Padding: item.padding,
				Radius: item.radius, Label: item.label,
			}, func() {
				mygoui.Typography(c, item.label, mygoui.TypographyOptions{Type: theme.TypeBodySM})
			})
		}
	})
}

// separatorHorizontalBody renders every horizontal emphasis.
func separatorHorizontalBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
		for _, item := range []struct {
			name  string
			color mygoui.SeparatorColor
		}{
			{"Default", mygoui.SeparatorDefault},
			{"Secondary", mygoui.SeparatorSecondary},
			{"Tertiary", mygoui.SeparatorTertiary},
		} {
			ui.Column(c).FillWidth().Gap(theme.Space4).Children(func() {
				mygoui.Typography(c, item.name, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
				mygoui.Separator(c, mygoui.SeparatorOptions{Color: item.color})
			})
		}
	})
}

// separatorVerticalBody renders vertical rules between inline labels.
func separatorVerticalBody(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Height(theme.Scale(c, theme.Space24)).
		Gap(theme.Space12).Children(func() {
		mygoui.Typography(c, "Overview", mygoui.TypographyOptions{Type: theme.TypeBodySM})
		mygoui.Separator(c, mygoui.SeparatorOptions{
			Orientation: mygoui.Vertical, Color: mygoui.SeparatorSecondary,
		})
		mygoui.Typography(c, "Activity", mygoui.TypographyOptions{Type: theme.TypeBodySM})
		mygoui.Separator(c, mygoui.SeparatorOptions{
			Orientation: mygoui.Vertical, Color: mygoui.SeparatorTertiary,
		})
		mygoui.Typography(c, "Settings", mygoui.TypographyOptions{Type: theme.TypeBodySM})
	})
}

// kbdChordsBody renders shortcut chords with their descriptions.
func kbdChordsBody(c *ui.Context) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		for _, item := range []struct {
			label string
			keys  []string
		}{
			{"Command menu", []string{"⌘", "K"}},
			{"Command palette", []string{"Ctrl", "Shift", "P"}},
			{"Dismiss", []string{"Esc"}},
		} {
			ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
				mygoui.Kbd(c, mygoui.KbdOptions{Keys: item.keys})
				mygoui.Typography(c, item.label, mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Muted: true,
				})
			})
		}
	})
}

// kbdLightBody renders the background-free keycap for inline text.
func kbdLightBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Typography(c, "Press", mygoui.TypographyOptions{Type: theme.TypeBodySM})
		mygoui.Kbd(c, mygoui.KbdOptions{Keys: []string{"Esc"}, Light: true})
		mygoui.Typography(c, "to close the panel.", mygoui.TypographyOptions{
			Type: theme.TypeBodySM,
		})
	})
}

// cardVariantsBody renders the four card surfaces with the same parts.
func cardVariantsBody(c *ui.Context) {
	cards := []struct {
		name        string
		variant     mygoui.CardVariant
		description string
	}{
		{"Default", mygoui.CardDefault, "Shadowed content card."},
		{"Secondary", mygoui.CardSecondary, "Quiet grouping on a page."},
		{"Tertiary", mygoui.CardTertiary, "Deeper neutral surface."},
		{"Transparent", mygoui.CardTransparent, "Composition without a fill."},
	}
	ui.Row(c).Wrap().Gap(theme.Space12).Children(func() {
		for _, item := range cards {
			ui.Column(c).Grow(1).Basis(0).MinWidth(theme.Scale(c, 200)).Children(func() {
				mygoui.Card(c, mygoui.CardOptions{Variant: item.variant}, func() {
					mygoui.CardHeader(c, func() {
						mygoui.CardTitle(c, item.name)
						mygoui.CardDescription(c, item.description)
					})
					mygoui.CardContent(c, func() {
						mygoui.Chip(c, mygoui.ChipOptions{
							Label: "Live", Color: mygoui.ColorSuccess, Variant: mygoui.ChipSoft,
							Size: mygoui.ChipSmall,
						})
					})
				})
			})
		}
	})
}

// cardSectionsBody renders all three slots of one card.
func cardSectionsBody(c *ui.Context) {
	ui.Column(c).FillWidth().Children(func() {
		mygoui.Card(c, mygoui.CardOptions{}, func() {
			mygoui.CardHeader(c, func() {
				mygoui.CardTitle(c, "Deploy preview")
				mygoui.CardDescription(c, "Ships in about two minutes.")
			})
			mygoui.CardContent(c, func() {
				mygoui.Typography(c, "Header, content and footer are separate slots, so media "+
					"and actions keep their own rhythm.", mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Muted: true,
				})
			})
			mygoui.CardFooter(c, func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Open preview", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				})
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Dismiss", Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				})
			})
		})
	})
}

// cardInteractiveBody renders a card that is one native action.
func cardInteractiveBody(c *ui.Context, st *foundationsState) {
	mygoui.Card(c, mygoui.CardOptions{
		Variant: mygoui.CardSecondary, Label: "Open project details",
		OnClick: func() { st.CardClicks++ },
	}, func() {
		mygoui.CardHeader(c, func() {
			mygoui.CardTitle(c, "Project Atlas")
			mygoui.CardDescription(c, "Click anywhere on the card.")
		})
	})
	mygoui.Typography(c, fmt.Sprintf("Card activations: %d", st.CardClicks), mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Muted: true,
	})
}

// linkIconBody renders the leading-icon link.
func linkIconBody(c *ui.Context, st *foundationsState) {
	mygoui.Link(c, mygoui.LinkOptions{
		Label: "Documentation", Icon: lucide.ArrowRight(),
		OnClick: func() { st.LinkClicks++ },
	})
}

// linkTooltipBody renders a link with its hover hint.
func linkTooltipBody(c *ui.Context) {
	mygoui.Link(c, mygoui.LinkOptions{
		Label: "Release notes", Tooltip: "Open the component docs",
	})
}

// linkDisabledBody renders the dimmed, non-activating link.
func linkDisabledBody(c *ui.Context) {
	mygoui.Link(c, mygoui.LinkOptions{Label: "Disabled link", Disabled: true})
}

// linkButtonBasicBody renders the inline text action.
func linkButtonBasicBody(c *ui.Context, st *foundationsState) {
	mygoui.LinkButton(c, mygoui.LinkButtonOptions{
		Label: "View all", OnClick: func() { st.LinkClicks++ },
	})
}

// linkButtonIconBody renders the inline action with a trailing icon.
func linkButtonIconBody(c *ui.Context) {
	mygoui.LinkButton(c, mygoui.LinkButtonOptions{
		Label: "Continue", Icon: lucide.ChevronRight(),
	})
}

// linkButtonDisabledBody renders the dimmed inline action.
func linkButtonDisabledBody(c *ui.Context) {
	mygoui.LinkButton(c, mygoui.LinkButtonOptions{Label: "Unavailable", Disabled: true})
}

// listRowsBody renders plain actionable rows with a chevron.
func listRowsBody(c *ui.Context, st *foundationsState) {
	mygoui.ListGroup(c, mygoui.ListGroupOptions{}, func() {
		mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
			Title: "Ada Lovelace", Description: "Owner", Chevron: true,
			OnClick: func() { st.ListChoice = "Ada Lovelace" },
		})
		mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
			Title: "Security", Description: "Managed by your organization", Chevron: true,
			OnClick: func() { st.ListChoice = "Security" },
		})
	})
	mygoui.Typography(c, "Selected row: "+listChoice(st), mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Muted: true,
	})
}

// listPrefixBody renders rows with a leading icon slot.
func listPrefixBody(c *ui.Context) {
	mygoui.ListGroup(c, mygoui.ListGroupOptions{}, func() {
		mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
			Title: "Members", Description: "Twelve people",
			Prefix: func() {
				ui.Icon(c, lucide.Users()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
					TextColor(theme.Colors(c).Muted)
			},
		})
		mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
			Title: "Alerts", Description: "Four unread updates",
			Prefix: func() {
				ui.Icon(c, lucide.Bell()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
					TextColor(theme.Colors(c).Muted)
			},
		})
	})
}

// listSuffixBody renders rows with a trailing badge.
func listSuffixBody(c *ui.Context) {
	mygoui.ListGroup(c, mygoui.ListGroupOptions{}, func() {
		mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
			Title: "Notifications", Description: "Four unread updates", Chevron: true,
			Suffix: func() {
				mygoui.Chip(c, mygoui.ChipOptions{
					Label: "4", Color: mygoui.ColorAccent, Variant: mygoui.ChipPrimary,
					Size: mygoui.ChipSmall,
				})
			},
		})
		mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
			Title: "Deployments", Description: "Two running", Chevron: true,
			Suffix: func() {
				mygoui.Chip(c, mygoui.ChipOptions{
					Label: "2", Color: mygoui.ColorSuccess, Variant: mygoui.ChipSoft,
					Size: mygoui.ChipSmall,
				})
			},
		})
	})
}

// listDisabledBody renders a row that cannot be activated.
func listDisabledBody(c *ui.Context) {
	mygoui.ListGroup(c, mygoui.ListGroupOptions{}, func() {
		mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
			Title: "Delete workspace", Description: "Unavailable for owners", Disabled: true,
			Prefix: func() {
				ui.Icon(c, lucide.Trash2()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
					TextColor(theme.Colors(c).Muted)
			},
			OnClick: func() {},
		})
	})
}

// listChoice keeps the selected-row line consistent between frames.
func listChoice(st *foundationsState) string {
	if st.ListChoice == "" {
		return "None"
	}
	return st.ListChoice
}

// emptyTextBody renders the placeholder HeroUI shows without content.
func emptyTextBody(c *ui.Context) {
	ui.Column(c).AlignItems(ui.Center).Children(func() {
		mygoui.EmptyState(c, mygoui.EmptyStateOptions{})
	})
}

// emptyDefaultBody renders the placeholder without a recovery action.
func emptyDefaultBody(c *ui.Context) {
	ui.Column(c).AlignItems(ui.Center).Children(func() {
		mygoui.EmptyState(c, mygoui.EmptyStateOptions{
			Title: "No saved searches", Description: "Searches you save appear here.",
			Icon: lucide.Search(),
		})
	})
}

// emptyActionBody renders the placeholder with a real recovery action.
func emptyActionBody(c *ui.Context, st *foundationsState) {
	ui.Column(c).AlignItems(ui.Center).Children(func() {
		mygoui.EmptyState(c, mygoui.EmptyStateOptions{
			Title: "No saved searches", Description: "Searches you save appear here.",
			Icon: lucide.Search(),
			Action: func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Create search", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
					Icon: lucide.Plus(), OnClick: func() { st.ActionClicks++ },
				})
			},
		})
	})
	mygoui.Typography(c, fmt.Sprintf("Searches created: %d", st.ActionClicks), mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Muted: true,
	})
}

// scrollShadowVerticalBody renders the vertically fading viewport.
func scrollShadowVerticalBody(c *ui.Context, st *foundationsState) {
	t := theme.Colors(c)
	mygoui.ScrollShadow(c, &st.Scroll, mygoui.ScrollShadowOptions{
		Orientation: mygoui.ScrollShadowVertical, Height: theme.Scale(c, 180),
		Surface: t.SurfaceSecondary, Label: "Scrollable activity",
	}, func() {
		ui.Column(c).Gap(theme.Space8).Padding(theme.Space4).Children(func() {
			for i := 1; i <= 12; i++ {
				ui.Row(c).AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
					ui.Box(c).Size(theme.Scale(c, 6), theme.Scale(c, 6)).
						Radius(theme.Scale(c, 3)).Background(t.Accent)
					mygoui.Typography(c, fmt.Sprintf("Build #%d finished successfully", i),
						mygoui.TypographyOptions{Type: theme.TypeBodySM})
				})
			}
		})
	})
}

// scrollShadowHorizontalBody renders the horizontally fading viewport.
func scrollShadowHorizontalBody(c *ui.Context, st *foundationsState) {
	t := theme.Colors(c)
	mygoui.ScrollShadow(c, &st.ScrollX, mygoui.ScrollShadowOptions{
		Orientation: mygoui.ScrollShadowHorizontal, Height: theme.Scale(c, 96),
		Surface: t.SurfaceSecondary, Label: "Scrollable environments",
	}, func() {
		ui.Row(c).Gap(theme.Space8).Padding(theme.Space4).Children(func() {
			for i := 1; i <= 8; i++ {
				mygoui.Surface(c, mygoui.SurfaceOptions{
					Variant: mygoui.SurfaceDefault, Padding: theme.Space12,
					Radius: theme.RadiusScale().XL, Label: fmt.Sprintf("Environment %d", i),
				}, func() {
					ui.Column(c).Width(theme.Scale(c, 140)).Gap(theme.Space4).Children(func() {
						mygoui.Typography(c, fmt.Sprintf("Environment %d", i),
							mygoui.TypographyOptions{Type: theme.TypeBodySM})
						mygoui.Typography(c, "Healthy", mygoui.TypographyOptions{
							Type: theme.TypeBodyXS, Muted: true,
						})
					})
				})
			}
		})
	})
}

// headerListBody renders a section label above its rows.
func headerListBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space4).Children(func() {
		mygoui.Header(c, mygoui.HeaderOptions{Text: "Workspace"})
		mygoui.ListGroup(c, mygoui.ListGroupOptions{}, func() {
			mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
				Title: "Members", Description: "Twelve people", Chevron: true,
			})
			mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
				Title: "Billing", Description: "Annual plan", Chevron: true,
			})
		})
	})
}

// headerSectionsBody renders several labelled sections in one column.
func headerSectionsBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
		for _, section := range []struct {
			label string
			rows  []string
		}{
			{"Account", []string{"Profile", "Sessions"}},
			{"Organization", []string{"Members", "Billing"}},
		} {
			ui.Column(c).FillWidth().Gap(theme.Space4).Children(func() {
				mygoui.Header(c, mygoui.HeaderOptions{Text: section.label})
				mygoui.ListGroup(c, mygoui.ListGroupOptions{}, func() {
					for _, row := range section.rows {
						mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{Title: row, Chevron: true})
					}
				})
			})
		}
	})
}
