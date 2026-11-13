package navigation

import (
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StepStatus overrides the state a step takes from its place before, at or
// after the current one.
type StepStatus uint8

const (
	// StepAuto follows StepsOptions.Current.
	StepAuto StepStatus = iota
	// StepError marks a step that failed, wherever it is.
	StepError
)

// StepsOrientation lays the steps in a row or a column.
type StepsOrientation uint8

const (
	StepsHorizontal StepsOrientation = iota
	StepsVertical
)

// StepItem is one step of a process.
type StepItem struct {
	Title       string
	Description string
	Status      StepStatus
}

// StepsOptions configures a step indicator. Current is the caller-owned
// index of the step in progress: the steps before it are completed and those
// after it upcoming; len(Items) or more completes them all. OnSelect makes
// each step a button that reports its index.
type StepsOptions struct {
	Items       []StepItem
	Current     int
	Orientation StepsOrientation
	Label       string
	OnSelect    func(index int)
}

type stepState uint8

const (
	stepUpcoming stepState = iota
	stepCurrent
	stepCompleted
	stepFailed
)

func stepStateOf(i, current int, status StepStatus) stepState {
	switch {
	case status == StepError:
		return stepFailed
	case i < current:
		return stepCompleted
	case i == current:
		return stepCurrent
	}
	return stepUpcoming
}

// Steps shows the numbered steps of a process joined by connectors. A
// completed step shows a check in place of its number, and the current step
// is exposed as such to assistive technology.
func Steps(c *ui.Context, o StepsOptions) ui.Element {
	vertical := o.Orientation == StepsVertical
	var root ui.Element
	if vertical {
		root = ui.Column(c).AlignItems(ui.Stretch)
	} else {
		root = ui.Row(c).AlignItems(ui.Start).Gap(scale(c, theme.Space8))
	}
	root.FillWidth().MinWidth(0).Role(ui.RoleList)
	if o.Label != "" {
		root.Label(o.Label)
	}
	root.Children(func() {
		for i, it := range o.Items {
			buildStep(c, o, i, it, vertical)
		}
	})
	return root
}

func buildStep(c *ui.Context, o StepsOptions, i int, it StepItem, vertical bool) {
	t := theme.Colors(c)
	state := stepStateOf(i, o.Current, it.Status)
	last := i == len(o.Items)-1
	line := t.Border
	if state == stepCompleted {
		line = t.Primary
	}
	gap := scale(c, theme.Space8)
	side := scale(c, theme.Space24)

	item := ui.Row(c.Key(i)).Gap(gap).MinWidth(0).Role(ui.RoleListItem)
	// A pressable step names its button; the item is only its container.
	describe := func(e ui.Element) {
		e.Label(it.Title)
		if it.Description != "" {
			e.Description(it.Description)
		}
		if state == stepCurrent {
			e.Current(ui.CurrentStep)
		}
	}
	if o.OnSelect == nil {
		describe(item)
	}
	if vertical {
		item.AlignItems(ui.Stretch)
	} else {
		item.AlignItems(ui.Start)
		if !last {
			item.Grow(1).Basis(0)
		}
	}
	item.Children(func() {
		var head ui.Element
		// pad is the room a pressable step keeps for its hover face.
		var pad float32
		if o.OnSelect != nil {
			pad = scale(c, theme.Space4)
			head = ui.ButtonBase(c).FocusRing(true).Padding(pad)
			describe(head)
			head = style.Background(head, theme.RadiusScale().LG, ui.Transparent, t.Accent, t.Accent)
			head.OnClick(func() { o.OnSelect(i) })
		} else {
			head = ui.Box(c)
		}
		head.Row().Justify(ui.Start).Gap(gap).MinWidth(0).Shrink(1)
		if vertical {
			head.AlignItems(ui.Stretch).Grow(1)
		} else {
			head.AlignItems(ui.Start)
		}
		head.Children(func() {
			if vertical {
				rail := ui.Column(c).AlignItems(ui.Center).Gap(scale(c, theme.Space4)).Shrink(0)
				rail.Children(func() {
					stepIndicator(c, t, i, state, side)
					if !last {
						ui.Box(c).Width(2).Grow(1).MinHeight(scale(c, theme.Space16)).
							Margin(0, 0, scale(c, theme.Space4)-pad, 0).Radius(1).Background(line)
					}
				})
			} else {
				stepIndicator(c, t, i, state, side)
			}
			stepText(c, t, it, state, side, vertical && !last)
		})
		if !vertical && !last {
			// The connector sits on the indicator's center line.
			ui.Box(c).Height(2).Grow(1).MinWidth(scale(c, theme.Space16)).
				Margin(pad+side/2-1, 0, 0, 0).Radius(1).Background(line)
		}
	})
}

func stepIndicator(c *ui.Context, t theme.Palette, i int, state stepState, side float32) {
	bg, fg := t.MutedBackground, t.MutedForeground
	switch state {
	case stepCurrent, stepCompleted:
		bg, fg = t.Primary, t.PrimaryForeground
	case stepFailed:
		bg, fg = t.Destructive, ui.RGB(255, 255, 255)
	}
	dot := ui.Box(c).Size(side, side).Radius(side / 2).Center().Shrink(0).
		Background(bg).TextColor(fg)
	dot.Children(func() {
		glyph := scale(c, theme.IconSM)
		switch state {
		case stepCompleted:
			ui.Icon(c, lucide.Check()).Size(glyph, glyph).Shrink(0)
		case stepFailed:
			ui.Icon(c, lucide.X()).Size(glyph, glyph).Shrink(0)
		default:
			style.Text(ui.Text(c, strconv.Itoa(i+1)), c, theme.XS()).
				FontWeight(theme.WeightMedium).FontFeatures("tnum")
		}
	})
}

func stepText(c *ui.Context, t theme.Palette, it StepItem, state stepState, side float32, spaced bool) {
	if it.Title == "" && it.Description == "" {
		return
	}
	title := theme.SM()
	// Center a one-line title on the indicator.
	top := max(0, (side-theme.Scale(c, title.LineHeight))/2)
	body := ui.Column(c).MinWidth(0).Shrink(1).Gap(scale(c, theme.Space2)).Padding(top, 0, 0, 0)
	if spaced {
		body.Padding(top, 0, scale(c, theme.Space16), 0)
	}
	body.Children(func() {
		if it.Title != "" {
			fg := t.Foreground
			switch state {
			case stepUpcoming:
				fg = t.MutedForeground
			case stepFailed:
				fg = t.Destructive
			}
			style.Text(ui.Text(c, it.Title), c, title).FontWeight(theme.WeightMedium).TextColor(fg)
		}
		if it.Description != "" {
			style.Text(ui.Text(c, it.Description), c, theme.XS()).TextColor(t.MutedForeground)
		}
	})
}
