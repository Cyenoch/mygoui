package gallery

import (
	"github.com/Cyenoch/mygoui"
	"github.com/egoist/mygo/ui"
)

// displayStateOf keeps the shared specimen state for the display demos.
func displayStateOf(s *State) *displayState {
	st := StateOf[displayState](s, "display")
	st.prepare()
	return st
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "chip", Title: "Chip", Group: "Display", Order: 10,
		Description: "Small status capsules across every color and variant.",
		Examples: []ComponentExample{
			{
				Title:       "Variants",
				Description: "Primary, secondary, soft and tertiary surface treatments.",
				Demo:        func(c *ui.Context, s *State) { chipVariantsBody(c) },
			},
			{
				Title:       "Status colors",
				Description: "Error, info, success and warning.",
				Demo:        func(c *ui.Context, s *State) { chipColorsBody(c) },
			},
			{
				Title:       "Sizes",
				Description: "Small, medium and large chip densities.",
				Demo:        func(c *ui.Context, s *State) { chipSizesBody(c) },
			},
			{
				Title:       "With icons",
				Description: "A leading icon pairs with the status color.",
				Demo:        func(c *ui.Context, s *State) { chipIconsBody(c) },
			},
		},
		Code: `mygoui.Chip(c, mygoui.ChipOptions{
    Label: "Beta", Color: mygoui.ColorAccent,
    Variant: mygoui.ChipSoft, Size: mygoui.ChipSmall,
})`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Chip text."},
			{"Color", "mygoui.StatusColor", "ColorDefault", "Default, accent, info, success, warning or danger."},
			{"Variant", "mygoui.ChipVariant", "ChipSecondary", "Primary, secondary, soft or tertiary."},
			{"Size", "mygoui.ChipSize", "ChipMedium", "Small, medium or large."},
			{"Icon", "*ui.SVG", "nil", "Leading icon, sized with the chip's text."},
		},
		Tokens: []string{"DefaultSoft", "DefaultSoftForeground", "AccentSoft", "SuccessSoft"},
	})

	RegisterComponent(ComponentPage{
		ID: "badge", Title: "Badge", Group: "Display", Order: 20,
		Description: "A compact semantic label with eight Coss variants.",
		Examples: []ComponentExample{
			{
				Title:       "Variants",
				Description: "Eight semantic Coss variants.",
				Demo:        func(c *ui.Context, s *State) { badgeVariantsBody(c) },
			},
			{
				Title:       "Status colors",
				Description: "Error, info, success and warning.",
				Demo:        func(c *ui.Context, s *State) { badgeColorsBody(c) },
			},
			{
				Title:       "Sizes",
				Description: "Small, medium and large labels.",
				Demo:        func(c *ui.Context, s *State) { badgeSizesBody(c) },
			},
			{
				Title: "Native anchors",
				Demo:  func(c *ui.Context, s *State) { badgeAnchorsBody(c) },
			},
		},
		Code: `mygoui.Badge(c, mygoui.BadgeOptions{
    Label: "9", Variant: mygoui.BadgePrimary,
})`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Badge text."},
			{"Variant", "mygoui.BadgeVariant", "BadgePrimary", "Default, secondary, outline, destructive, error, info, success or warning."},
			{"Size", "mygoui.BadgeSize", "BadgeMedium", "Small, medium or large."},
		},
		Tokens: []string{"Primary", "PrimaryForeground", "Secondary", "Destructive", "Info", "Success", "Warning"},
	})

	RegisterComponent(ComponentPage{
		ID: "avatar", Title: "Avatar", Group: "Display", Order: 30,
		Description: "A person or entity portrait with initials fallback and rings.",
		Examples: []ComponentExample{
			{
				Title:       "Sizes",
				Description: "Small, medium and large portraits from generated initials.",
				Demo:        func(c *ui.Context, s *State) { avatarSizesBody(c) },
			},
			{
				Title:       "Soft colors",
				Description: "Soft status fills replace the neutral surface.",
				Demo:        func(c *ui.Context, s *State) { avatarSoftBody(c) },
			},
			{
				Title:       "Rings",
				Description: "A background-colored ring separates overlapping portraits.",
				Demo:        func(c *ui.Context, s *State) { avatarRingsBody(c) },
			},
		},
		Code: `mygoui.Avatar(c, mygoui.AvatarOptions{
    Name: "Ada Lovelace", Size: mygoui.AvatarMedium,
    Soft: true, Color: mygoui.ColorAccent, Ring: true,
})`,
		Params: []ComponentParam{
			{"Name", "string", "", "The name whose initials form the fallback."},
			{"Fallback", "string", "", "Overrides the generated initials."},
			{"Image", "ui.ImageSource", "nil", "Portrait image, when one is available."},
			{"Size", "mygoui.AvatarSize", "AvatarMedium", "Small, medium or large."},
			{"Soft", "bool", "false", "Uses the soft status fill instead of a flat one."},
			{"Ring", "bool", "false", "Draws the background-colored ring."},
			{"Color", "mygoui.StatusColor", "ColorDefault", "Status color of the fill."},
			{"Badge", "func()", "nil", "Content anchored at the bottom-right corner."},
		},
		Tokens: []string{"Default", "DefaultSoftForeground", "AccentSoft", "Background"},
	})

	RegisterComponent(ComponentPage{
		ID: "avatargroup", Title: "AvatarGroup", Group: "Display", Order: 40,
		Description: "Overlapping avatars with an optional overflow count or grid layout.",
		Examples: []ComponentExample{
			{
				Title:       "Overlap",
				Description: "Portraits overlap by 8 DIPs with separating rings.",
				Demo:        func(c *ui.Context, s *State) { avatarGroupOverlapBody(c) },
			},
			{
				Title:       "Overflow count",
				Description: "A positive Total adds a remainder count after the visible avatars.",
				Demo:        func(c *ui.Context, s *State) { avatarGroupOverflowBody(c) },
			},
			{
				Title:       "Grid",
				Description: "Grid lays the portraits out without overlap.",
				Demo:        func(c *ui.Context, s *State) { avatarGroupGridBody(c) },
			},
		},
		Code: `mygoui.AvatarGroup(c, mygoui.AvatarGroupOptions{
    Size: mygoui.AvatarSmall, Ring: true, Total: 7,
    Avatars: []mygoui.AvatarOptions{
        {Name: "Ada Lovelace"}, {Name: "Grace Hopper"},
    },
})`,
		Params: []ComponentParam{
			{"Avatars", "[]mygoui.AvatarOptions", "nil", "The avatars to overlap."},
			{"Size", "mygoui.AvatarSize", "AvatarMedium", "Size applied to every avatar and the count."},
			{"Overlap", "float32", "8", "How far each avatar overlaps the previous one."},
			{"Ring", "bool", "false", "Separates avatars with the background ring."},
			{"Grid", "bool", "false", "Lays the avatars out as a grid instead of a row."},
			{"Total", "int", "0", "Positive values add an overflow count."},
		},
		Tokens: []string{"Default", "DefaultSoftForeground", "AccentSoft", "Background"},
	})

	RegisterComponent(ComponentPage{
		ID: "tag", Title: "Tag", Group: "Display", Order: 50,
		Description: "A compact label that can be selected or removed on its own.",
		Examples: []ComponentExample{
			{
				Title:       "With icons",
				Description: "A leading icon labels the tag's subject.",
				Demo:        func(c *ui.Context, s *State) { tagIconsBody(c) },
			},
			{
				Title:       "Removable",
				Description: "A removed tag fades out where it stood while the caller drops it.",
				Demo:        func(c *ui.Context, s *State) { tagRemovableBody(c) },
			},
			{
				Title:       "Disabled",
				Description: "A disabled tag dims, keeps no hover face and stops removing.",
				Demo:        func(c *ui.Context, s *State) { tagDisabledBody(c) },
			},
		},
		Code: `mygoui.Tag(c, mygoui.TagOptions{
    Key: "engineering", Label: "Engineering", Removable: true,
    OnRemove: func() { remove("engineering") },
})`,
		Params: []ComponentParam{
			{"Key", "string", "required", "Identity used by a TagGroup and by the removal fade."},
			{"Label", "string", "required", "Visible tag text."},
			{"Icon", "*ui.SVG", "nil", "Leading icon."},
			{"Removable", "bool", "false", "Shows the remove action."},
			{"Disabled", "bool", "false", "Removes selection and removal."},
			{"OnRemove", "func()", "nil", "Runs when the remove action is used."},
		},
		Tokens: []string{"Default", "DefaultForeground", "Border", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "taggroup", Title: "TagGroup", Group: "Display", Order: 60,
		Description: "Tag collections with single, multiple or no selection.",
		Examples: []ComponentExample{
			{
				Title:       "Single selection",
				Description: "Choosing a tag clears the previous choice.",
				Demo: func(c *ui.Context, s *State) {
					tagGroupSingleBody(c, displayStateOf(s))
				},
			},
			{
				Title:       "Multiple selection",
				Description: "Each tag toggles independently and the count follows.",
				Demo: func(c *ui.Context, s *State) {
					tagGroupMultiBody(c, displayStateOf(s))
				},
			},
			{
				Title:       "Removable",
				Description: "Removal updates the caller's list and can be restored.",
				Demo: func(c *ui.Context, s *State) {
					tagGroupRemovableBody(c, displayStateOf(s))
				},
			},
		},
		Code: `mygoui.TagGroup(c, mygoui.TagGroupOptions{
    Label: "Languages", Selection: mygoui.TagSelectMultiple,
    State: &state.Tags,
}, []mygoui.TagOptions{
    {Key: "go", Label: "Go"},
    {Key: "rust", Label: "Rust"},
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Group label above the tags."},
			{"Description", "string", "", "Supporting line under the label."},
			{"Selection", "mygoui.TagSelection", "TagSelectNone", "None, single or multiple selection."},
			{"Variant", "mygoui.TagVariant", "TagDefault", "Tag surface treatment."},
			{"State", "*mygoui.TagGroupState", "nil", "Caller-owned selection set."},
			{"OnChange", "func([]string)", "nil", "Reports the selected keys."},
		},
		Tokens: []string{"Default", "AccentSoft", "Border", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "alert", Title: "Alert", Group: "Display", Order: 70,
		Description: "A status callout with an optional icon and action.",
		Examples: []ComponentExample{
			{
				Title:       "Status colors",
				Description: "Status callouts without icons.",
				Demo:        func(c *ui.Context, s *State) { alertStatusBody(c) },
			},
			{
				Title:       "With action",
				Description: "The Action slot carries a real recovery control.",
				Demo: func(c *ui.Context, s *State) {
					alertActionBody(c, displayStateOf(s))
				},
			},
			{
				Title:       "Status icon",
				Description: "Opt into the status glyph.",
				Demo:        func(c *ui.Context, s *State) { alertStatusIconBody(c) },
			},
			{
				Title:       "Custom icon",
				Description: "An explicit icon is shown even without ShowIcon.",
				Demo:        func(c *ui.Context, s *State) { alertIconBody(c) },
			},
		},
		Code: `mygoui.Alert(c, mygoui.AlertOptions{
    Title: "Deploy finished",
    Description: "The new build is live.",
    Status: mygoui.ColorSuccess,
    ShowIcon: false, // Set true to show the status glyph.
})`,
		Params: []ComponentParam{
			{"Title", "string", "required", "Headline of the callout."},
			{"Description", "string", "", "Supporting line under the title."},
			{"Status", "mygoui.StatusColor", "ColorDefault", "Default, accent, info, success, warning or danger."},
			{"ShowIcon", "bool", "false", "Shows the status glyph when Icon is nil."},
			{"Icon", "*ui.SVG", "nil", "Shows a custom icon instead of the status glyph."},
			{"Action", "func()", "nil", "Builds a trailing action."},
		},
		Tokens: []string{"Surface", "SurfaceForeground", "SuccessSoftForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "spinner", Title: "Spinner", Group: "Display", Order: 80,
		Description: "A native indeterminate indicator for pending work.",
		Examples: []ComponentExample{
			{
				Title:       "Sizes",
				Description: "Small through extra large arcs.",
				Demo:        func(c *ui.Context, s *State) { spinnerSizesBody(c) },
			},
			{
				Title:       "Colors",
				Description: "The status color vocabulary applied to the arc.",
				Demo:        func(c *ui.Context, s *State) { spinnerColorsBody(c) },
			},
		},
		Code: `mygoui.Spinner(c, mygoui.SpinnerOptions{
    Size: mygoui.SpinnerMedium, Color: mygoui.ColorAccent, Label: "Loading",
})`,
		Params: []ComponentParam{
			{"Size", "mygoui.SpinnerSize", "SpinnerMedium", "Small, medium, large or extra large."},
			{"Color", "mygoui.StatusColor", "ColorAccent", "Status color of the arc."},
			{"Label", "string", "Loading", "Accessible name announced while spinning."},
		},
		Tokens: []string{"Accent", "Success", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "skeleton", Title: "Skeleton", Group: "Display", Order: 90,
		Description: "A placeholder shape that stands in for loading content.",
		Examples: []ComponentExample{
			{
				Title:       "Shapes",
				Description: "A circle and two bars hold the eventual layout.",
				Demo:        func(c *ui.Context, s *State) { skeletonShapesBody(c) },
			},
			{
				Title:       "Animations",
				Description: "A 2 s shimmer sweep, a 2 s pulse, or no animation at all.",
				Demo:        func(c *ui.Context, s *State) { skeletonAnimationsBody(c) },
			},
		},
		Code: `mygoui.Skeleton(c, mygoui.SkeletonOptions{
    Width: 180, Height: 12, Animation: mygoui.SkeletonShimmer,
})`,
		Params: []ComponentParam{
			{"Width", "float32", "fill", "Fixed width; zero fills the row."},
			{"Height", "float32", "theme icon", "Placeholder height."},
			{"Radius", "float32", "theme radius", "Corner radius; circles derive it from the height."},
			{"Circle", "bool", "false", "Makes the placeholder a circle."},
			{"Animation", "mygoui.SkeletonAnimation", "SkeletonShimmer", "Shimmer, pulse or none."},
		},
		Tokens: []string{"SurfaceTertiary", "Surface"},
	})

	RegisterComponent(ComponentPage{
		ID: "skeletongroup", Title: "SkeletonGroup", Group: "Display", Order: 100,
		Description: "One synchronized shimmer across several placeholders.",
		Examples: []ComponentExample{
			{
				Title:       "Shimmer",
				Description: "One moving band crosses every shape in the group.",
				Demo: func(c *ui.Context, s *State) {
					skeletonGroupBody(c, mygoui.SkeletonShimmer)
				},
			},
			{
				Title:       "Pulse",
				Description: "The group fades in place instead of sweeping.",
				Demo: func(c *ui.Context, s *State) {
					skeletonGroupBody(c, mygoui.SkeletonPulse)
				},
			},
		},
		Code: `mygoui.SkeletonGroup(c, mygoui.SkeletonGroupOptions{
    Animation: mygoui.SkeletonPulse,
}, func() {
    mygoui.Skeleton(c, mygoui.SkeletonOptions{
        Height: 12, Animation: mygoui.SkeletonNone,
    })
})`,
		Params: []ComponentParam{
			{"Animation", "mygoui.SkeletonAnimation", "SkeletonShimmer", "Animation the group paints across its shapes."},
		},
		Tokens: []string{"SurfaceTertiary"},
	})

	RegisterComponent(ComponentPage{
		ID: "progressbar", Title: "ProgressBar", Group: "Display", Order: 110,
		Description: "A linear determinate or indeterminate progress indicator.",
		Examples: []ComponentExample{
			{
				Title:       "Determinate",
				Description: "A value between its bounds; each change eases the fill over 300 ms.",
				Demo: func(c *ui.Context, s *State) {
					progressBarBody(c, displayStateOf(s))
				},
			},
			{
				Title:       "Sizes",
				Description: "The three track heights with their indicator colors.",
				Demo:        func(c *ui.Context, s *State) { progressSizesBody(c) },
			},
			{
				Title:       "Indeterminate",
				Description: "A looping 2/5 band while the total is unknown.",
				Demo:        func(c *ui.Context, s *State) { progressIndeterminateBody(c) },
			},
		},
		Code: `mygoui.ProgressBar(c, mygoui.ProgressOptions{
    Label: "Upload", Value: 65, Max: 100, ShowOutput: true,
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Title above the track."},
			{"Value", "float64", "0", "Current value between Min and Max."},
			{"Min", "float64", "0", "Lower bound."},
			{"Max", "float64", "100", "Upper bound; 0 with a 0 Min means 100."},
			{"Indeterminate", "bool", "false", "Loops the indicator while the total is unknown."},
			{"Color", "mygoui.IndicatorColor", "IndicatorAccent", "Accent, default, success, warning or danger."},
			{"Size", "mygoui.ProgressSize", "ProgressMedium", "Small, medium or large track."},
			{"ShowOutput", "bool", "false", "Shows the value beside the label."},
			{"OutputText", "string", "", "Replaces the computed percentage."},
		},
		Tokens: []string{"Accent", "Default", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "progresscircle", Title: "ProgressCircle", Group: "Display", Order: 120,
		Description: "A ring that shows determinate progress or spins while pending.",
		Examples: []ComponentExample{
			{
				Title:       "Determinate",
				Description: "The ring eases to its value over 300 ms and can print it at its center.",
				Demo: func(c *ui.Context, s *State) {
					progressCircleDeterminateBody(c, displayStateOf(s))
				},
			},
			{
				Title:       "Sizes",
				Description: "Small, medium and large rings.",
				Demo:        func(c *ui.Context, s *State) { progressCircleSizesBody(c) },
			},
			{
				Title:       "Indeterminate",
				Description: "A spinning ring while the total is unknown.",
				Demo:        func(c *ui.Context, s *State) { progressCircleIndeterminateBody(c) },
			},
		},
		Code: `mygoui.ProgressCircle(c, mygoui.ProgressCircleOptions{
    Label: "Download", Value: 65, Max: 100, ShowOutput: true,
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Accessible name; defaults to the percentage."},
			{"Value", "float64", "0", "Current value between Min and Max."},
			{"Min", "float64", "0", "Lower bound."},
			{"Max", "float64", "100", "Upper bound; 0 with a 0 Min means 100."},
			{"Indeterminate", "bool", "false", "Spins the ring while the total is unknown."},
			{"Color", "mygoui.IndicatorColor", "IndicatorAccent", "Status color of the ring."},
			{"Size", "mygoui.ProgressSize", "ProgressMedium", "Small, medium or large ring."},
			{"ShowOutput", "bool", "false", "Prints the percentage at the ring's center."},
		},
		Tokens: []string{"Accent", "Default", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "meter", Title: "Meter", Group: "Display", Order: 130,
		Description: "A labelled quantity with its value and bounds made explicit.",
		Examples: []ComponentExample{
			{
				Title:       "Quantities",
				Description: "Each meter measures a different range; changes ease over 300 ms.",
				Demo: func(c *ui.Context, s *State) {
					meterQuantitiesBody(c, displayStateOf(s))
				},
			},
			{
				Title:       "Status colors",
				Description: "The indicator color separates healthy from critical ranges.",
				Demo:        func(c *ui.Context, s *State) { meterColorsBody(c) },
			},
			{
				Title:       "Custom output",
				Description: "OutputText replaces the computed percentage.",
				Demo:        func(c *ui.Context, s *State) { meterOutputBody(c) },
			},
		},
		Code: `mygoui.Meter(c, mygoui.MeterOptions{
    Label: "Storage", Value: 65, Max: 100, ShowOutput: true,
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Title above the track."},
			{"Value", "float64", "0", "Measured value between Min and Max."},
			{"Min", "float64", "0", "Lower bound."},
			{"Max", "float64", "100", "Upper bound; 0 with a 0 Min means 100."},
			{"Color", "mygoui.IndicatorColor", "IndicatorAccent", "Status color of the fill."},
			{"Size", "mygoui.ProgressSize", "ProgressMedium", "Small, medium or large track."},
			{"ShowOutput", "bool", "false", "Shows the value beside the label."},
			{"OutputText", "string", "", "Overrides the formatted value output."},
		},
		Tokens: []string{"Accent", "Default", "Foreground", "Success"},
	})
}
