package gallery

import (
	"github.com/egoist/mygo/ui"
)

// foundationsStateOf keeps the shared specimen state for the foundations
// demos, so pages that show interactive rows keep their selection.
func foundationsStateOf(s *State) *foundationsState {
	return StateOf[foundationsState](s, "foundations")
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "typography", Title: "Typography", Group: "Foundations", Order: 10,
		Description: "Semantic type roles, from page headings to captions, code and prose.",
		Examples: []ComponentExample{
			{
				Title:       "Headings",
				Description: "Six semibold heading levels share one type scale and tight tracking.",
				Demo:        func(c *ui.Context, s *State) { typographyHeadingsBody(c) },
			},
			{
				Title:       "Body roles",
				Description: "Body, body SM and body XS cover prose, support and captions.",
				Demo:        func(c *ui.Context, s *State) { typographyBodyBody(c) },
			},
			{
				Title:       "Code and prose",
				Description: "The code role is an inline chip that stays selectable in Prose.",
				Demo:        func(c *ui.Context, s *State) { typographyCodeProseBody(c) },
			},
		},
		Code: `mygoui.Typography(c, "Account settings", mygoui.TypographyOptions{
    Type: theme.TypeH4,
})
mygoui.Typography(c, "Changes apply immediately.", mygoui.TypographyOptions{
    Type: theme.TypeBodySM, Muted: true,
})`,
		Params: []ComponentParam{
			{"Type", "theme.SemanticType", "TypeBody", "Semantic role: a heading level or a body, code, prose role."},
			{"Muted", "bool", "false", "Paints the secondary foreground for supporting copy."},
			{"Selectable", "bool", "false", "Lets the user select and copy the text."},
			{"Truncate", "bool", "false", "Keeps one line and ends it with an ellipsis."},
			{"Align", "ui.Align", "ui.Start", "Horizontal alignment of the text."},
			{"Weight", "int", "0", "Overrides the role's weight; zero keeps it."},
			{"Level", "int", "0", "Heading level for assistive technology; zero derives it from Type."},
		},
		Tokens: []string{"Foreground", "Muted", "Default"},
	})

	RegisterComponent(ComponentPage{
		ID: "surface", Title: "Surface", Group: "Foundations", Order: 20,
		Description: "A rounded painted panel that carries the matching foreground.",
		Examples: []ComponentExample{
			{
				Title:       "Variants",
				Description: "All four elevation tokens with matching foreground colors.",
				Demo:        func(c *ui.Context, s *State) { surfaceVariantsBody(c) },
			},
			{
				Title:       "Padding and radius",
				Description: "The same surface at three insets and corner radii.",
				Demo:        func(c *ui.Context, s *State) { surfacePaddingBody(c) },
			},
		},
		Code: `mygoui.Surface(c, mygoui.SurfaceOptions{
    Variant: mygoui.SurfaceSecondary, FillWidth: true, Padding: theme.Space16,
}, func() {
    mygoui.Typography(c, "Surface content", mygoui.TypographyOptions{
        Type: theme.TypeBodySM,
    })
})`,
		Params: []ComponentParam{
			{"Variant", "mygoui.SurfaceVariant", "SurfaceDefault", "Chooses the surface and foreground pair."},
			{"Padding", "float32", "0", "Inner padding in DIPs."},
			{"Radius", "float32", "theme radius", "Corner radius; zero keeps the theme scale."},
			{"FillWidth", "bool", "false", "Stretches the surface across its row."},
			{"Label", "string", "", "Accessible name when the surface is a landing area."},
		},
		Tokens: []string{"Surface", "SurfaceForeground", "SurfaceSecondary", "SurfaceTertiary"},
	})

	RegisterComponent(ComponentPage{
		ID: "separator", Title: "Separator", Group: "Foundations", Order: 30,
		Description: "A one DIP rule that divides sections or columns.",
		Examples: []ComponentExample{
			{
				Title:       "Horizontal",
				Description: "Horizontal rules at every emphasis.",
				Demo:        func(c *ui.Context, s *State) { separatorHorizontalBody(c) },
			},
			{
				Title:       "Vertical",
				Description: "Vertical rules divide inline labels in one row.",
				Demo:        func(c *ui.Context, s *State) { separatorVerticalBody(c) },
			},
		},
		Code: `mygoui.Separator(c, mygoui.SeparatorOptions{
    Orientation: mygoui.Horizontal,
})`,
		Params: []ComponentParam{
			{"Orientation", "mygoui.Orientation", "Horizontal", "Horizontal rules divide rows; vertical rules divide columns."},
			{"Color", "mygoui.SeparatorColor", "SeparatorDefault", "Default, secondary or tertiary separator tone."},
		},
		Tokens: []string{"Separator", "SeparatorSecondary", "SeparatorTertiary"},
	})

	RegisterComponent(ComponentPage{
		ID: "kbd", Title: "Kbd", Group: "Foundations", Order: 40,
		Description: "Keycaps and chords that document a shortcut.",
		Examples: []ComponentExample{
			{
				Title:       "Chords",
				Description: "Each chord draws its keys as separate caps.",
				Demo:        func(c *ui.Context, s *State) { kbdChordsBody(c) },
			},
			{
				Title:       "Light",
				Description: "The light keycap drops its background for inline text.",
				Demo:        func(c *ui.Context, s *State) { kbdLightBody(c) },
			},
		},
		Code: `mygoui.Kbd(c, mygoui.KbdOptions{
    Keys: []string{"Cmd", "K"},
})`,
		Params: []ComponentParam{
			{"Keys", "[]string", "nil", "Keys of the chord, drawn as separate caps."},
			{"Light", "bool", "false", "Drops the keycap background for inline use."},
		},
		Tokens: []string{"Default", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "card", Title: "Card", Group: "Foundations", Order: 50,
		Description: "A padded surface with header, content and footer slots.",
		Examples: []ComponentExample{
			{
				Title:       "Variants",
				Description: "Four card variants built from the same parts.",
				Demo:        func(c *ui.Context, s *State) { cardVariantsBody(c) },
			},
			{
				Title:       "With sections",
				Description: "Header, content and footer keep their own rhythm.",
				Demo:        func(c *ui.Context, s *State) { cardSectionsBody(c) },
			},
			{
				Title:       "Interactive",
				Description: "OnClick turns the whole card into one native action.",
				Demo: func(c *ui.Context, s *State) {
					cardInteractiveBody(c, foundationsStateOf(s))
				},
			},
		},
		Code: `mygoui.Card(c, mygoui.CardOptions{Variant: mygoui.CardSecondary}, func() {
    mygoui.CardHeader(c, func() {
        mygoui.CardTitle(c, "Deploy preview")
        mygoui.CardDescription(c, "Ships in about two minutes.")
    })
    mygoui.CardContent(c, func() { /* rows, charts, media */ })
    mygoui.CardFooter(c, func() { /* actions */ })
})`,
		Params: []ComponentParam{
			{"Variant", "mygoui.CardVariant", "CardDefault", "Default, secondary, tertiary or transparent surface."},
			{"Label", "string", "", "Accessible name for an interactive card."},
			{"OnClick", "func()", "nil", "Turns the whole card into one native action."},
			{"Handle", "*ui.Handle", "nil", "Receives focus and activation handles."},
		},
		Tokens: []string{"Surface", "SurfaceSecondary", "SurfaceSecondaryForeground", "SurfaceShadows"},
	})

	RegisterComponent(ComponentPage{
		ID: "link", Title: "Link", Group: "Foundations", Order: 60,
		Description: "An inline text action with a native focus ring.",
		Examples: []ComponentExample{
			{
				Title:       "With icon",
				Description: "A leading icon points at the destination.",
				Demo: func(c *ui.Context, s *State) {
					linkIconBody(c, foundationsStateOf(s))
				},
			},
			{
				Title:       "With tooltip",
				Description: "The tooltip explains where the link goes on hover.",
				Demo:        func(c *ui.Context, s *State) { linkTooltipBody(c) },
			},
			{
				Title:       "Disabled",
				Description: "The link dims and stops accepting activation.",
				Demo:        func(c *ui.Context, s *State) { linkDisabledBody(c) },
			},
		},
		Code: `mygoui.Link(c, mygoui.LinkOptions{
    Label: "Documentation", Icon: lucide.ArrowRight(),
    OnClick: func() { openDocs() },
})`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Visible and accessible text."},
			{"Icon", "*ui.SVG", "nil", "Leading icon, sized by the link."},
			{"Disabled", "bool", "false", "Removes activation and dims the text."},
			{"OnClick", "func()", "nil", "Runs on click, Space or Return."},
			{"Tooltip", "string", "", "Hint shown on hover."},
		},
		Tokens: []string{"Link", "Focus"},
	})

	RegisterComponent(ComponentPage{
		ID: "linkbutton", Title: "LinkButton", Group: "Foundations", Order: 70,
		Description: "A ghost text action with no hover fill, for inline navigation.",
		Examples: []ComponentExample{
			{
				Title:       "Basic",
				Description: "Inline navigation without a button surface.",
				Demo: func(c *ui.Context, s *State) {
					linkButtonBasicBody(c, foundationsStateOf(s))
				},
			},
			{
				Title:       "With icon",
				Description: "A trailing chevron marks the continuation.",
				Demo:        func(c *ui.Context, s *State) { linkButtonIconBody(c) },
			},
			{
				Title:       "Disabled",
				Description: "The action dims and stops accepting activation.",
				Demo:        func(c *ui.Context, s *State) { linkButtonDisabledBody(c) },
			},
		},
		Code: `mygoui.LinkButton(c, mygoui.LinkButtonOptions{
    Label: "View all", Icon: lucide.ChevronRight(),
    OnClick: func() { showAll() },
})`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Visible and accessible text."},
			{"Icon", "*ui.SVG", "nil", "Leading icon."},
			{"Disabled", "bool", "false", "Removes activation and dims the text."},
			{"OnClick", "func()", "nil", "Runs on click, Space or Return."},
		},
		Tokens: []string{"Foreground", "DefaultHover", "Focus"},
	})

	RegisterComponent(ComponentPage{
		ID: "listgroup", Title: "ListGroup", Group: "Foundations", Order: 80,
		Description: "A rounded surface of actionable rows with prefixes and suffixes.",
		Examples: []ComponentExample{
			{
				Title:       "Rows",
				Description: "Title, description and a disclosure chevron.",
				Demo: func(c *ui.Context, s *State) {
					listRowsBody(c, foundationsStateOf(s))
				},
			},
			{
				Title:       "With prefixes",
				Description: "A leading icon names the row's subject.",
				Demo:        func(c *ui.Context, s *State) { listPrefixBody(c) },
			},
			{
				Title:       "With suffixes",
				Description: "A trailing chip carries the row's count.",
				Demo:        func(c *ui.Context, s *State) { listSuffixBody(c) },
			},
			{
				Title:       "Disabled row",
				Description: "A disabled row keeps its layout but cannot activate.",
				Demo:        func(c *ui.Context, s *State) { listDisabledBody(c) },
			},
		},
		Code: `mygoui.ListGroup(c, mygoui.ListGroupOptions{}, func() {
    mygoui.ListGroupItem(c, mygoui.ListGroupItemOptions{
        Title: "Ada Lovelace", Description: "Owner", Chevron: true,
        OnClick: func() { open("Ada Lovelace") },
    })
})`,
		Params: []ComponentParam{
			{"Variant", "mygoui.SurfaceVariant", "SurfaceDefault", "Surface the group paints behind its rows."},
			{"Label", "string", "", "Accessible name for the group."},
			{"Title", "string", "required", "Item title (ListGroupItem)."},
			{"Description", "string", "", "Supporting line under the title."},
			{"Prefix", "func()", "nil", "Leading slot of a row (ListGroupItem)."},
			{"Suffix", "func()", "nil", "Trailing slot of a row (ListGroupItem)."},
			{"Chevron", "bool", "false", "Shows the disclosure chevron."},
			{"Disabled", "bool", "false", "Removes activation and dims the row."},
			{"OnClick", "func()", "nil", "Makes the row a native action (ListGroupItem)."},
		},
		Tokens: []string{"Surface", "SurfaceForeground", "SurfaceHover", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "emptystate", Title: "EmptyState", Group: "Foundations", Order: 90,
		Description: "A low-emphasis placeholder with a real recovery action.",
		Examples: []ComponentExample{
			{
				Title:       "Default copy",
				Description: "With no content it shows the localized empty placeholder.",
				Demo:        func(c *ui.Context, s *State) { emptyTextBody(c) },
			},
			{
				Title:       "Default",
				Description: "An icon, title and one line explain what is missing.",
				Demo:        func(c *ui.Context, s *State) { emptyDefaultBody(c) },
			},
			{
				Title:       "With action",
				Description: "The Action slot gives the reader a way forward.",
				Demo: func(c *ui.Context, s *State) {
					emptyActionBody(c, foundationsStateOf(s))
				},
			},
		},
		Code: `mygoui.EmptyState(c, mygoui.EmptyStateOptions{
    Title: "No saved searches",
    Description: "Searches you save appear here.",
    Icon: lucide.Search(),
})`,
		Params: []ComponentParam{
			{"Title", "string", "", "What is missing; empty shows the locale placeholder."},
			{"Description", "string", "", "One line that explains the next step."},
			{"Icon", "*ui.SVG", "nil", "Muted leading glyph."},
			{"Action", "func()", "nil", "Builds the recovery button."},
		},
		Tokens: []string{"Muted", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "scrollshadow", Title: "ScrollShadow", Group: "Foundations", Order: 100,
		Description: "A native scroll viewport that fades the edges holding more content.",
		Examples: []ComponentExample{
			{
				Title:       "Vertical",
				Description: "A fading viewport over a tall list of activity rows.",
				Demo: func(c *ui.Context, s *State) {
					scrollShadowVerticalBody(c, foundationsStateOf(s))
				},
			},
			{
				Title:       "Horizontal",
				Description: "The same edge fade along a wide row of cards.",
				Demo: func(c *ui.Context, s *State) {
					scrollShadowHorizontalBody(c, foundationsStateOf(s))
				},
			},
		},
		Code: `mygoui.ScrollShadow(c, &state.Feed, mygoui.ScrollShadowOptions{
    Orientation: mygoui.ScrollShadowVertical,
    Surface:     theme.Colors(c).SurfaceSecondary,
    Height:      theme.Scale(c, 180),
    Label:       "Scrollable activity",
}, func() { /* rows */ })`,
		Params: []ComponentParam{
			{"Orientation", "mygoui.ScrollShadowOrientation", "ScrollShadowVertical", "Scroll axis to fade."},
			{"Surface", "ui.Color", "Background", "The opaque surface behind the fade."},
			{"Size", "float32", "theme.Space40", "Depth of the edge fade."},
			{"Height", "float32", "0", "Fixed height; zero lets the parent size it."},
			{"Label", "string", "", "Accessible name for the scroll area."},
		},
		Tokens: []string{"Background", "SurfaceSecondary", "Scrollbar"},
	})

	RegisterComponent(ComponentPage{
		ID: "header", Title: "Header", Group: "Foundations", Order: 110,
		Description: "The small muted label that titles a section, list or menu.",
		Examples: []ComponentExample{
			{
				Title:       "Section label",
				Description: "A quiet label above a grouped set of rows.",
				Demo:        func(c *ui.Context, s *State) { headerListBody(c) },
			},
			{
				Title:       "Multiple sections",
				Description: "Several labelled groups stack with a wider section gap.",
				Demo:        func(c *ui.Context, s *State) { headerSectionsBody(c) },
			},
		},
		Code: `mygoui.Header(c, mygoui.HeaderOptions{Text: "Workspace"})`,
		Params: []ComponentParam{
			{"Text", "string", "required", "Label shown above the group."},
		},
		Tokens: []string{"Muted"},
	})
}
