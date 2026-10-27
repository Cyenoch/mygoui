package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type displayState struct {
	Prepared bool

	Tags       mygoui.TagGroupState
	SingleTags mygoui.TagGroupState
	Removed    map[string]bool

	Progress    float64
	AlertClicks int
}

func (s *displayState) prepare() {
	if s.Prepared {
		return
	}
	s.Prepared = true
	s.Progress = 65
	s.Removed = map[string]bool{}
	s.Tags.Selected = map[string]bool{"go": true}
	s.SingleTags.Selected = map[string]bool{"react": true}
}

var displayColors = []struct {
	label string
	color mygoui.StatusColor
}{
	{"Default", mygoui.ColorDefault},
	{"Accent", mygoui.ColorAccent},
	{"Success", mygoui.ColorSuccess},
	{"Warning", mygoui.ColorWarning},
	{"Danger", mygoui.ColorDanger},
}

// chipVariantsBody renders the four chip surface treatments.
func chipVariantsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, item := range []struct {
			label   string
			variant mygoui.ChipVariant
		}{
			{"Secondary", mygoui.ChipSecondary},
			{"Primary", mygoui.ChipPrimary},
			{"Soft", mygoui.ChipSoft},
			{"Tertiary", mygoui.ChipTertiary},
		} {
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: item.label, Color: mygoui.ColorAccent, Variant: item.variant,
			})
		}
	})
}

// chipColorsBody renders the status color vocabulary in the soft variant.
func chipColorsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, color := range displayColors {
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: color.label, Color: color.color, Variant: mygoui.ChipSoft,
			})
		}
	})
}

// chipSizesBody renders the three chip densities.
func chipSizesBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, item := range []struct {
			label string
			size  mygoui.ChipSize
		}{
			{"Small", mygoui.ChipSmall},
			{"Medium", mygoui.ChipMedium},
			{"Large", mygoui.ChipLarge},
		} {
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: item.label, Size: item.size,
				Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft,
			})
		}
	})
}

// chipIconsBody renders chips with a leading icon.
func chipIconsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, item := range []struct {
			label string
			icon  *ui.SVG
			color mygoui.StatusColor
		}{
			{"New", lucide.Sparkles(), mygoui.ColorAccent},
			{"Shipped", lucide.CircleCheck(), mygoui.ColorSuccess},
			{"Pending", lucide.Clock(), mygoui.ColorWarning},
		} {
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: item.label, Icon: item.icon, Color: item.color,
				Variant: mygoui.ChipSoft,
			})
		}
	})
}

// badgeVariantsBody renders the eight Coss badge variants.
func badgeVariantsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, item := range []struct {
			label   string
			variant mygoui.BadgeVariant
		}{
			{"Default", mygoui.BadgePrimary},
			{"Secondary", mygoui.BadgeSecondary},
			{"Outline", mygoui.BadgeOutline},
			{"Destructive", mygoui.BadgeDestructive},
			{"Error", mygoui.BadgeError},
			{"Info", mygoui.BadgeInfo},
			{"Success", mygoui.BadgeSuccess},
			{"Warning", mygoui.BadgeWarning},
		} {
			mygoui.Badge(c, mygoui.BadgeOptions{
				Label: item.label, Variant: item.variant,
			})
		}
	})
}

// badgeColorsBody renders the status colors in the soft variant.
func badgeColorsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, item := range []struct {
			label   string
			variant mygoui.BadgeVariant
		}{
			{"Failed", mygoui.BadgeError}, {"Queued", mygoui.BadgeInfo},
			{"Deployed", mygoui.BadgeSuccess}, {"Pending", mygoui.BadgeWarning},
		} {
			mygoui.Badge(c, mygoui.BadgeOptions{
				Label: item.label, Variant: item.variant,
			})
		}
	})
}

// badgeSizesBody renders the three badge densities.
func badgeSizesBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, item := range []struct {
			label string
			size  mygoui.BadgeSize
		}{
			{"Small", mygoui.BadgeSmall},
			{"Medium", mygoui.BadgeMedium},
			{"Large", mygoui.BadgeLarge},
		} {
			mygoui.Badge(c, mygoui.BadgeOptions{
				Label: item.label, Size: item.size,
				Variant: mygoui.BadgeSecondary,
			})
		}
	})
}

func badgeAnchorsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space24).Children(func() {
		mygoui.BadgeAnchor(c, mygoui.BadgeAnchorOptions{
			Placement: mygoui.BadgeTopRight,
			Badge:     mygoui.BadgeOptions{Label: "4", Variant: mygoui.BadgeDestructive},
		}, func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Inbox", Icon: lucide.Bell(), Variant: mygoui.ButtonOutline,
			})
		})
		mygoui.BadgeAnchor(c, mygoui.BadgeAnchorOptions{
			Placement: mygoui.BadgeBottomRight,
			Badge: mygoui.BadgeOptions{
				Label: "New", Variant: mygoui.BadgeSuccess, Size: mygoui.BadgeSmall,
			},
		}, func() {
			mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Ada Lovelace"})
		})
	})
}

// avatarSizesBody renders the three avatar heights.
func avatarSizesBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space16).Children(func() {
		mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Ada Lovelace", Size: mygoui.AvatarSmall})
		mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Grace Hopper", Size: mygoui.AvatarMedium})
		mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Katherine Johnson", Size: mygoui.AvatarLarge})
	})
}

// avatarSoftBody renders the soft status fills.
func avatarSoftBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Avatar(c, mygoui.AvatarOptions{
			Name: "Alan Turing", Soft: true, Color: mygoui.ColorAccent,
		})
		mygoui.Avatar(c, mygoui.AvatarOptions{
			Name: "Margaret Hamilton", Soft: true, Color: mygoui.ColorSuccess,
		})
		mygoui.Avatar(c, mygoui.AvatarOptions{
			Name: "Barbara Liskov", Soft: true, Color: mygoui.ColorWarning,
		})
	})
}

// avatarRingsBody renders the background-colored ring around avatars.
func avatarRingsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Ada Lovelace", Ring: true})
		mygoui.Avatar(c, mygoui.AvatarOptions{
			Name: "Grace Hopper", Soft: true, Color: mygoui.ColorAccent, Ring: true,
		})
		mygoui.Avatar(c, mygoui.AvatarOptions{
			Name: "Alan Turing", Soft: true, Color: mygoui.ColorDanger, Ring: true,
		})
	})
}

// avatarGroupOverlapBody renders avatars overlapping with separating rings.
func avatarGroupOverlapBody(c *ui.Context) {
	mygoui.AvatarGroup(c, mygoui.AvatarGroupOptions{
		Size: mygoui.AvatarSmall, Ring: true,
		Avatars: []mygoui.AvatarOptions{
			{Name: "Ada Lovelace"},
			{Name: "Grace Hopper", Soft: true, Color: mygoui.ColorAccent},
			{Name: "Alan Turing", Soft: true, Color: mygoui.ColorSuccess},
			{Name: "Katherine Johnson", Soft: true, Color: mygoui.ColorWarning},
		},
	})
}

// avatarGroupOverflowBody renders the overflow-count variant.
func avatarGroupOverflowBody(c *ui.Context) {
	mygoui.AvatarGroup(c, mygoui.AvatarGroupOptions{
		Size: mygoui.AvatarSmall, Ring: true, Total: 7,
		Avatars: []mygoui.AvatarOptions{
			{Name: "Ada Lovelace"},
			{Name: "Grace Hopper", Soft: true, Color: mygoui.ColorAccent},
			{Name: "Alan Turing", Soft: true, Color: mygoui.ColorSuccess},
		},
	})
}

// avatarGroupGridBody renders the grid arrangement.
func avatarGroupGridBody(c *ui.Context) {
	mygoui.AvatarGroup(c, mygoui.AvatarGroupOptions{
		Size: mygoui.AvatarSmall, Grid: true,
		Avatars: []mygoui.AvatarOptions{
			{Name: "Ada Lovelace"},
			{Name: "Grace Hopper", Soft: true, Color: mygoui.ColorAccent},
			{Name: "Alan Turing"},
			{Name: "Katherine Johnson", Soft: true, Color: mygoui.ColorDanger},
		},
	})
}

// tagIconsBody renders standalone tags with a leading icon.
func tagIconsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Tag(c, mygoui.TagOptions{Key: "go", Label: "Go", Icon: lucide.Code()})
		mygoui.Tag(c, mygoui.TagOptions{Key: "design", Label: "Design", Icon: lucide.Sparkles()})
		mygoui.Tag(c, mygoui.TagOptions{Key: "review", Label: "Review", Icon: lucide.Check()})
	})
}

// tagRemovableBody renders tags that own a remove action.
func tagRemovableBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Tag(c, mygoui.TagOptions{
			Key: "engineering", Label: "Engineering", Removable: true, OnRemove: func() {},
		})
		mygoui.Tag(c, mygoui.TagOptions{
			Key: "research", Label: "Research", Removable: true, OnRemove: func() {},
		})
	})
}

// tagDisabledBody renders the dimmed, non-interactive tag.
func tagDisabledBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Tag(c, mygoui.TagOptions{Key: "active", Label: "Active"})
		mygoui.Tag(c, mygoui.TagOptions{Key: "archived", Label: "Archived", Disabled: true})
	})
}

// tagGroupSingleBody renders single selection with its chosen key.
func tagGroupSingleBody(c *ui.Context, st *displayState) {
	mygoui.TagGroup(c, mygoui.TagGroupOptions{
		Label: "Single selection", Description: "Pick one framework.",
		Selection: mygoui.TagSelectSingle, State: &st.SingleTags,
	}, []mygoui.TagOptions{
		{Key: "react", Label: "React"},
		{Key: "vue", Label: "Vue"},
		{Key: "svelte", Label: "Svelte"},
	})
	selected := "None"
	for _, key := range []string{"react", "vue", "svelte"} {
		if st.SingleTags.IsSelected(key) {
			selected = key
			break
		}
	}
	mygoui.Typography(c, "Selected: "+selected, mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Muted: true,
	})
}

// tagGroupMultiBody renders multiple selection with a count.
func tagGroupMultiBody(c *ui.Context, st *displayState) {
	mygoui.TagGroup(c, mygoui.TagGroupOptions{
		Label: "Multiple selection", Description: "Choose any combination.",
		Selection: mygoui.TagSelectMultiple, Size: mygoui.TagSmall,
		Variant: mygoui.TagSurface, State: &st.Tags,
	}, []mygoui.TagOptions{
		{Key: "go", Label: "Go", Icon: lucide.Code()},
		{Key: "rust", Label: "Rust", Icon: lucide.Code()},
		{Key: "swift", Label: "Swift", Icon: lucide.Code()},
		{Key: "kotlin", Label: "Kotlin", Icon: lucide.Code()},
	})
	mygoui.Typography(c, fmt.Sprintf("Selected: %d tags", len(st.Tags.Selected)),
		mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
}

// tagGroupRemovableBody renders removal backed by the caller's list.
func tagGroupRemovableBody(c *ui.Context, st *displayState) {
	removable := []mygoui.TagOptions{
		{Key: "design", Label: "Design"},
		{Key: "engineering", Label: "Engineering"},
		{Key: "research", Label: "Research"},
		{Key: "operations", Label: "Operations"},
	}
	tags := make([]mygoui.TagOptions, 0, len(removable))
	for _, tag := range removable {
		if st.Removed[tag.Key] {
			continue
		}
		tag.Removable = true
		key := tag.Key
		tag.OnRemove = func() { st.Removed[key] = true }
		tags = append(tags, tag)
	}
	mygoui.TagGroup(c, mygoui.TagGroupOptions{
		Label: "Removable", Description: "Removing a tag updates the caller's list.",
		Selection: mygoui.TagSelectNone,
	}, tags)
	if len(tags) == 0 {
		mygoui.Typography(c, "All tags removed.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	}
	if len(st.Removed) > 0 {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Restore tags", Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
			Icon: lucide.Plus(), OnClick: func() { st.Removed = map[string]bool{} },
		})
	}
}

// alertStatusBody renders one callout per status color.
func alertStatusBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
		for _, color := range displayColors {
			mygoui.Alert(c, mygoui.AlertOptions{
				Title:       color.label + " alert",
				Description: "This is the " + color.label + " status surface and foreground.",
				Status:      color.color,
			})
		}
	})
}

// alertActionBody renders a callout with a real recovery action.
func alertActionBody(c *ui.Context, st *displayState) {
	mygoui.Alert(c, mygoui.AlertOptions{
		Title:       "Deploy failed",
		Description: "The last build stopped before publishing.",
		Status:      mygoui.ColorDanger,
		Action: func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Retry deploy", Variant: mygoui.ButtonSecondary,
				Size: theme.ButtonSmall, Icon: lucide.Zap(),
				OnClick: func() { st.AlertClicks++ },
			})
		},
	})
	mygoui.Typography(c, fmt.Sprintf("Retries: %d", st.AlertClicks), mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Muted: true,
	})
}

// alertStatusIconBody renders a callout with its status glyph.
func alertStatusIconBody(c *ui.Context) {
	mygoui.Alert(c, mygoui.AlertOptions{
		Title:       "Build complete",
		Description: "Artifacts are ready to review.",
		Status:      mygoui.ColorSuccess,
		ShowIcon:    true,
	})
}

// alertIconBody renders a callout with a custom icon.
func alertIconBody(c *ui.Context) {
	mygoui.Alert(c, mygoui.AlertOptions{
		Title:       "Release is live",
		Description: "A custom icon appears without enabling the status glyph.",
		Status:      mygoui.ColorAccent,
		Icon:        lucide.Rocket(),
	})
}

// spinnerSizesBody renders the four spinner sizes.
func spinnerSizesBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space24).Children(func() {
		for _, item := range []struct {
			label string
			size  mygoui.SpinnerSize
		}{
			{"Small", mygoui.SpinnerSmall},
			{"Medium", mygoui.SpinnerMedium},
			{"Large", mygoui.SpinnerLarge},
			{"Extra large", mygoui.SpinnerExtraLarge},
		} {
			ui.Column(c).AlignItems(ui.Center).Gap(theme.Space4).Children(func() {
				mygoui.Spinner(c, mygoui.SpinnerOptions{
					Size: item.size, Label: item.label + " loading",
				})
				mygoui.Typography(c, item.label, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
			})
		}
	})
}

// spinnerColorsBody renders the status color vocabulary.
func spinnerColorsBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space24).Children(func() {
		for _, color := range displayColors {
			ui.Column(c).AlignItems(ui.Center).Gap(theme.Space4).Children(func() {
				mygoui.Spinner(c, mygoui.SpinnerOptions{
					Color: color.color, Label: color.label + " loading",
				})
				mygoui.Typography(c, color.label, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
			})
		}
	})
}

// skeletonShapesBody renders circle and bar placeholders.
func skeletonShapesBody(c *ui.Context) {
	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Skeleton(c, mygoui.SkeletonOptions{Circle: true, Width: 40, Height: 40})
		ui.Column(c).Grow(1).Basis(0).Gap(theme.Space8).Children(func() {
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Height: 12})
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 180, Height: 12})
		})
	})
}

// skeletonAnimationsBody renders one bar per animation mode.
func skeletonAnimationsBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
		for _, item := range []struct {
			label     string
			animation mygoui.SkeletonAnimation
		}{
			{"Shimmer", mygoui.SkeletonShimmer},
			{"Pulse", mygoui.SkeletonPulse},
			{"None", mygoui.SkeletonNone},
		} {
			ui.Column(c).Gap(theme.Space4).Children(func() {
				mygoui.Typography(c, item.label, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
				mygoui.Skeleton(c, mygoui.SkeletonOptions{
					Width: 220, Height: 12, Animation: item.animation,
				})
			})
		}
	})
}

// skeletonGroupBody renders a synchronized group with per-shape animation off.
func skeletonGroupBody(c *ui.Context, animation mygoui.SkeletonAnimation) {
	mygoui.SkeletonGroup(c, mygoui.SkeletonGroupOptions{Animation: animation}, func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
			mygoui.Skeleton(c, mygoui.SkeletonOptions{
				Circle: true, Width: 40, Height: 40, Animation: mygoui.SkeletonNone,
			})
			ui.Column(c).Grow(1).Basis(0).Gap(theme.Space8).Children(func() {
				mygoui.Skeleton(c, mygoui.SkeletonOptions{Height: 12, Animation: mygoui.SkeletonNone})
				mygoui.Skeleton(c, mygoui.SkeletonOptions{
					Width: 200, Height: 12, Animation: mygoui.SkeletonNone,
				})
			})
		})
	})
}

// progressBarBody renders the interactive determinate bar and its controls.
func progressBarBody(c *ui.Context, st *displayState) {
	ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
		mygoui.ProgressBar(c, mygoui.ProgressOptions{
			Label: "Upload progress", Value: st.Progress, Max: 100, ShowOutput: true,
			Color: mygoui.IndicatorAccent,
		})
		ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Advance", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				Icon: lucide.Plus(),
				OnClick: func() {
					if st.Progress >= 100 {
						st.Progress = 0
					} else {
						st.Progress += 10
					}
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Reset", Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				OnClick: func() { st.Progress = 0 },
			})
		})
	})
}

// progressSizesBody renders the three bar heights and their colors.
func progressSizesBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
		for _, item := range []struct {
			label string
			value float64
			size  mygoui.ProgressSize
			color mygoui.IndicatorColor
		}{
			{"Small", 30, mygoui.ProgressSmall, mygoui.IndicatorSuccess},
			{"Medium", 55, mygoui.ProgressMedium, mygoui.IndicatorWarning},
			{"Large", 80, mygoui.ProgressLarge, mygoui.IndicatorDanger},
		} {
			mygoui.ProgressBar(c, mygoui.ProgressOptions{
				Label: item.label, Value: item.value, Max: 100, Size: item.size,
				Color: item.color, ShowOutput: true,
			})
		}
	})
}

// progressIndeterminateBody renders the looping unknown-total bar.
func progressIndeterminateBody(c *ui.Context) {
	mygoui.ProgressBar(c, mygoui.ProgressOptions{
		Label: "Indexing documents", Indeterminate: true, Size: mygoui.ProgressLarge,
		Color: mygoui.IndicatorAccent,
	})
}

// progressCircleDeterminateBody renders a ring with centered output.
func progressCircleDeterminateBody(c *ui.Context, st *displayState) {
	mygoui.ProgressCircle(c, mygoui.ProgressCircleOptions{
		Label: "Sync progress", Value: st.Progress, Max: 100,
		Size: mygoui.ProgressLarge, ShowOutput: true, Color: mygoui.IndicatorAccent,
	})
}

// progressCircleSizesBody renders the three ring diameters.
func progressCircleSizesBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space24).Children(func() {
		for _, item := range []struct {
			label string
			size  mygoui.ProgressSize
			value float64
		}{
			{"Small", mygoui.ProgressSmall, 25},
			{"Medium", mygoui.ProgressMedium, 50},
			{"Large", mygoui.ProgressLarge, 75},
		} {
			ui.Column(c).AlignItems(ui.Center).Gap(theme.Space4).Children(func() {
				mygoui.ProgressCircle(c, mygoui.ProgressCircleOptions{
					Label: item.label + " progress", Value: item.value, Max: 100,
					Size: item.size, Color: mygoui.IndicatorSuccess,
				})
				mygoui.Typography(c, item.label, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
			})
		}
	})
}

// progressCircleIndeterminateBody renders the spinning unknown-total ring.
func progressCircleIndeterminateBody(c *ui.Context) {
	mygoui.ProgressCircle(c, mygoui.ProgressCircleOptions{
		Label: "Loading sync", Indeterminate: true,
		Size: mygoui.ProgressLarge, Color: mygoui.IndicatorWarning,
	})
}

// meterQuantitiesBody renders quantities with different ranges.
func meterQuantitiesBody(c *ui.Context, st *displayState) {
	ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
		mygoui.Meter(c, mygoui.MeterOptions{
			Label: "Storage used", Value: st.Progress, Max: 100,
			ShowOutput: true, Color: mygoui.IndicatorWarning,
		})
		mygoui.Meter(c, mygoui.MeterOptions{
			Label: "Seats assigned", Value: 8, Max: 12,
			ShowOutput: true, Size: mygoui.ProgressLarge, Color: mygoui.IndicatorSuccess,
		})
		mygoui.Meter(c, mygoui.MeterOptions{
			Label: "Request budget", Value: 92, Max: 100,
			ShowOutput: true, Size: mygoui.ProgressSmall, Color: mygoui.IndicatorDanger,
		})
	})
}

// meterColorsBody renders one value per indicator color.
func meterColorsBody(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
		for _, item := range []struct {
			label string
			color mygoui.IndicatorColor
		}{
			{"Accent", mygoui.IndicatorAccent},
			{"Success", mygoui.IndicatorSuccess},
			{"Warning", mygoui.IndicatorWarning},
			{"Danger", mygoui.IndicatorDanger},
		} {
			mygoui.Meter(c, mygoui.MeterOptions{
				Label: item.label, Value: 65, Max: 100,
				ShowOutput: true, Color: item.color,
			})
		}
	})
}

// meterOutputBody renders a meter whose output text overrides the percentage.
func meterOutputBody(c *ui.Context) {
	mygoui.Meter(c, mygoui.MeterOptions{
		Label: "Seats assigned", Value: 8, Max: 12, ShowOutput: true,
		OutputText: "8 of 12 seats", Color: mygoui.IndicatorAccent,
	})
}
