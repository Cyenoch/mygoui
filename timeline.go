package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TimelineItem is one event. Status colors its marker, which assistive
// technology reads by the status's name; Icon replaces the
// dot with a glyph on a tinted disc. Time is shown as given, after the
// title. Content builds caller content under the description, as a diff
// summary or an attachment. OnClick makes the text a button, as for opening
// the event's details.
type TimelineItem struct {
	Title, Description string
	Time               string
	Status             StatusColor
	Icon               *ui.SVG
	Content            func()
	OnClick            func()
}

// TimelineOptions configures a vertical list of events, in the order given.
// Pending marks the last item as still in progress: its marker is a spinner
// and the line leading to it is dashed. Recent, when positive and less than
// the number of items, shows only that many of the newest behind a control
// that reveals the earlier ones; Expanded is the caller-owned flag it
// toggles, kept with the element when nil, and EarlierLabel names the
// control, the locale's ShowMore when empty (the caller formats a count).
type TimelineOptions struct {
	Items        []TimelineItem
	Label        string
	Pending      bool
	Recent       int
	Expanded     *bool
	EarlierLabel string
}

// Timeline shows events down a line of markers. It builds every item that
// is shown: show a bounded history, and page or virtualize a log of
// unbounded length.
func Timeline(c *ui.Context, o TimelineOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Stretch).Role(ui.RoleList).Label(o.Label)
	hidden := 0
	if o.Recent > 0 && o.Recent < len(o.Items) {
		hidden = len(o.Items) - o.Recent
	}
	item := func(i int) {
		last := i == len(o.Items)-1
		timelineItem(c, &t, o.Items[i], timelinePlace{
			last:    last,
			pending: o.Pending && last,
			// The line into a pending item belongs to the item before it.
			dashed: o.Pending && i == len(o.Items)-2,
		}, i)
	}
	return root.Children(func() {
		if hidden == 0 {
			for i := range o.Items {
				item(i)
			}
			return
		}
		open := o.Expanded
		if open == nil {
			open = ui.Local(root, timelineOpenKey{}, func() bool { return false })
		}
		parts := ui.CollapsibleBase(c, open)
		timelineEarlier(c, &t, parts, o)
		panel := parts.Panel(func() {
			for i := range hidden {
				item(i)
			}
		})
		if panel.Valid() {
			panel.FillWidth().MinWidth(0)
			if progress := parts.Progress(); progress < 1 {
				panel.Opacity(max(progress, 0))
			}
		}
		for i := hidden; i < len(o.Items); i++ {
			item(i)
		}
	})
}

type timelineOpenKey struct{}

// timelineEarlier is the control that stands in the line for the items not
// shown: a chevron where a marker goes, and its label.
func timelineEarlier(c *ui.Context, t *theme.Palette, parts ui.CollapsibleParts, o TimelineOptions) {
	label := o.EarlierLabel
	if label == "" {
		label = locale.M().ShowMore
	}
	side := theme.Scale(c, theme.Space(5))
	trigger := parts.Trigger.MinWidth(0).MinHeight(side).AlignItems(ui.Center).Justify(ui.Start).
		Gap(theme.Space(3)).Cursor(ui.CursorPointer).Label(label)
	trigger = style.Background(trigger, theme.RadiusScale().MD, ui.Transparent, t.SurfaceHover, t.SurfaceHover)
	progress := parts.Progress()
	trigger.Children(func() {
		glyph := theme.Scale(c, theme.IconSM)
		ui.Box(c).Size(side, side).Shrink(0).Center().Children(func() {
			ui.Icon(c, lucide.ChevronDown()).Size(glyph, glyph).TextColor(t.MutedForeground).
				Rotate(180 * progress).Role(ui.RoleNone)
		})
		style.Text(ui.Text(c, label), c, theme.SM()).TextColor(t.MutedForeground).MinWidth(0)
	})
	_ = trigger.Clicked()
}

type timelineMarkerKind uint8

const (
	timelinePending timelineMarkerKind = iota
	timelineGlyph
	timelineDot
)

type (
	timelineFillKey struct{}
	timelineInkKey  struct{}
)

type timelinePlace struct {
	last, pending, dashed bool
}

func timelineItem(c *ui.Context, t *theme.Palette, it TimelineItem, at timelinePlace, key int) {
	side := theme.Scale(c, theme.Space(5))
	item := ui.Row(c.Key(key)).MinWidth(0).Gap(theme.Space(3)).AlignItems(ui.Stretch).
		Role(ui.RoleListItem).Label(it.Title).Busy(at.pending)
	if it.Description != "" {
		item.Description(it.Description)
	}
	item.Children(func() {
		ui.Column(c).Shrink(0).Width(side).AlignItems(ui.Center).Gap(theme.Space(1)).Children(func() {
			ui.Box(c).Size(side, side).Shrink(0).Center().Children(func() {
				timelineMarker(c, t, it, at.pending, side)
			})
			if !at.last {
				line := ui.Box(c).Width(0).Grow(1).MinHeight(theme.Space(4)).
					BorderStartWidth(theme.Hairline).BorderColor(t.Border)
				if at.dashed {
					line.BorderStyle(ui.BorderDashed)
				}
			}
		})
		outer := ui.Column(c).Grow(1).Basis(0).MinWidth(0)
		if !at.last {
			outer.Padding(0, 0, theme.Space(5), 0)
		}
		outer.Children(func() {
			if it.OnClick == nil {
				ui.Column(c).MinWidth(0).Gap(theme.Space(0.5)).Children(func() { timelineText(c, t, it) })
				return
			}
			// The face reaches past the text, so the text keeps its place.
			pad, padX := theme.Scale(c, theme.Space(1)), theme.Scale(c, theme.Space(2))
			button := ui.ButtonBase(c).Column().AlignItems(ui.Stretch).Justify(ui.Start).MinWidth(0).
				Gap(theme.Space(0.5)).Padding(pad, padX).Margin(-pad, -padX).Cursor(ui.CursorPointer).
				OnClick(it.OnClick)
			style.Background(button, theme.RadiusScale().MD, ui.Transparent, t.SurfaceHover, t.SurfaceHover).
				Children(func() { timelineText(c, t, it) })
		})
	})
}

func timelineMarker(c *ui.Context, t *theme.Palette, it TimelineItem, pending bool, side float32) {
	// Each kind of marker has its own key, so that pending giving way to
	// done crossfades.
	switch {
	case pending:
		ui.Box(c.Key(timelinePending)).Size(side, side).Center().Transition(fadeTransition()).Children(func() {
			Spinner(c, SpinnerOptions{Size: SpinnerSmall, Color: it.Status})
		})
	case it.Icon != nil:
		status := statusColors(t, it.Status)
		glyph := theme.Scale(c, theme.IconXS)
		disc := statusMark(ui.Box(c.Key(timelineGlyph)), it.Status).Size(side, side).Radius(side / 2).Center().
			Transition(fadeTransition())
		disc.Background(motion.Color(disc, timelineFillKey{}, status.Soft.Over(t.Background))).
			TextColor(motion.Color(disc, timelineInkKey{}, status.Foreground)).Children(func() {
			ui.Icon(c, it.Icon).Size(glyph, glyph).Role(ui.RoleNone)
		})
	default:
		ink := t.MutedForeground
		if it.Status != ColorDefault {
			ink = statusForeground(*t, it.Status)
		}
		dot := theme.Scale(c, theme.Space(2))
		mark := statusMark(ui.Box(c.Key(timelineDot)), it.Status).Size(dot, dot).Radius(dot / 2).
			Transition(fadeTransition())
		mark.Background(motion.Color(mark, timelineFillKey{}, ink))
	}
}

func timelineText(c *ui.Context, t *theme.Palette, it TimelineItem) {
	if it.Title != "" || it.Time != "" {
		// The time follows the title, wrapping under it where it does not fit.
		ui.Row(c).MinWidth(0).Wrap().GapX(theme.Space(2)).AlignItems(ui.Baseline).Children(func() {
			if it.Title != "" {
				style.Text(ui.Text(c, it.Title), c, theme.SM()).FontWeight(theme.WeightMedium).
					TextColor(t.Foreground).MinWidth(0)
			}
			if it.Time != "" {
				style.Text(ui.Text(c, it.Time), c, theme.XS()).TextColor(t.MutedForeground).
					FontFeatures("tnum").Shrink(0)
			}
		})
	}
	if it.Description != "" {
		style.Text(ui.Text(c, it.Description), c, theme.SM()).TextColor(t.MutedForeground).Selectable()
	}
	if it.Content != nil {
		it.Content()
	}
}
