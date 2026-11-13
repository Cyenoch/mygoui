package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatisticTrend is the direction of a metric's change.
type StatisticTrend uint8

const (
	// TrendFlat marks a change with no direction.
	TrendFlat StatisticTrend = iota
	TrendUp
	TrendDown
)

// StatisticOptions configures a metric. Value, Prefix, Suffix and Delta are
// shown as given: the caller formats numbers, units and signs. Delta shows
// under the value with the icon of Trend; DeltaStatus replaces the color the
// trend implies (success up, danger down), as for a cost that rose.
type StatisticOptions struct {
	Label          string
	Value          string
	Prefix, Suffix string
	Delta          string
	Trend          StatisticTrend
	DeltaStatus    StatusColor
	// Content builds caller content under the value and its change, as a
	// sparkline or a progress bar.
	Content func()
	// Loading keeps the label and shows a skeleton block in the value's
	// place, of the value's height; the delta waits for data.
	Loading bool
}

// Statistic shows a large value under its label, with an optional change.
// A value wider than its room is cut with an ellipsis and shown whole in a
// tooltip. Assistive technology reads it as a group named by the label.
func Statistic(c *ui.Context, o StatisticOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).MinWidth(0).Gap(theme.Space4).AlignItems(ui.Start).
		Role(ui.RoleGroup).Label(o.Label).Busy(o.Loading)
	return root.Children(func() {
		if o.Label != "" {
			style.Text(ui.Text(c, o.Label), c, theme.SM()).TextColor(t.MutedForeground)
		}
		align := ui.Baseline
		if o.Loading {
			align = ui.Center
		}
		ui.Row(c).MinWidth(0).Gap(theme.Space4).AlignItems(align).Children(func() {
			statisticAffix(c, t, o.Prefix)
			if o.Loading {
				SkeletonGroup(c, SkeletonGroupOptions{}, func() { loadingBar(c, theme.H2(), 96) }).Width(theme.Scale(c, 96))
				// A hidden digit gives the row the height the platform's
				// font gives the value, which a token cannot know.
				style.Text(ui.Text(c, "0"), c, theme.H2()).SingleLine().Width(0).Opacity(0).AccessHidden(true)
			} else {
				NumberFlip(c, NumberFlipOptions{Value: o.Value, Style: theme.H2(), Color: t.Foreground}).
					Shrink(1).MinWidth(0).Tooltip(o.Value)
			}
			statisticAffix(c, t, o.Suffix)
		})
		if o.Delta != "" && !o.Loading {
			statisticDelta(c, &t, o)
		}
		if o.Content != nil {
			o.Content()
		}
	})
}

func statisticAffix(c *ui.Context, t theme.Palette, text string) {
	if text != "" {
		style.Text(ui.Text(c, text), c, theme.LG()).TextColor(t.MutedForeground).Shrink(0)
	}
}

type statisticDeltaInkKey struct{}

func statisticDelta(c *ui.Context, t *theme.Palette, o StatisticOptions) {
	icon, name, ink := lucide.Minus(), "", t.MutedForeground
	switch o.Trend {
	case TrendUp:
		icon, name, ink = lucide.TrendingUp(), locale.M().Increase, statusColors(t, ColorSuccess).Foreground
	case TrendDown:
		icon, name, ink = lucide.TrendingDown(), locale.M().Decrease, statusColors(t, ColorDanger).Foreground
	}
	if o.DeltaStatus != ColorDefault {
		ink = statusColors(t, o.DeltaStatus).Foreground
	}
	row := ui.Row(c).Gap(theme.Space4).AlignItems(ui.Center)
	row.TextColor(motion.Color(row, statisticDeltaInkKey{}, ink)).Children(func() {
		glyph := theme.Scale(c, theme.IconSM)
		// The direction is named, so it does not rest on color alone.
		mark := stateIcon(c, icon).Size(glyph, glyph).Shrink(0)
		if name == "" {
			mark.Role(ui.RoleNone)
		} else {
			mark.Label(name)
		}
		style.Text(ui.Text(c, o.Delta), c, theme.XS()).FontWeight(theme.WeightMedium).FontFeatures("tnum")
	})
}
