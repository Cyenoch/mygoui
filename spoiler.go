package mygoui

import (
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SpoilerOptions configures clamped content. Expanded is the caller-owned
// flag the control toggles; nil keeps the flag with the element. MaxHeight
// is the height the content is cut at while collapsed, 96 scaled DIPs when
// zero. ShowLabel and HideLabel replace the locale's control labels.
type SpoilerOptions struct {
	Expanded             *bool
	MaxHeight            float32
	ShowLabel, HideLabel string
}

type spoilerOpenKey struct{}

// spoilerSlack keeps a rounding difference from counting as overflow.
const spoilerSlack = 0.5

// Spoiler cuts body at MaxHeight, fading its last lines out, and adds a
// control that shows the rest and hides it again. Content that fits shows whole, with no control. The
// height is that of the last layout, so the control follows the content one
// frame after it changes.
func Spoiler(c *ui.Context, o SpoilerOptions, body func()) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space(2)).AlignItems(ui.Start)
	expanded := o.Expanded
	if expanded == nil {
		expanded = ui.Local(root, "expanded", func() bool { return false })
	}
	limit := o.MaxHeight
	if limit <= 0 {
		limit = theme.Scale(c, 96)
	}
	// The content's laid-out height, which its paint records.
	height := ui.Local(root, "height", func() float32 { return 0 })
	overflows := *height > limit+spoilerSlack
	open := float32(0)
	if *expanded {
		open = 1
	}
	m := theme.Motion()
	open = root.AnimateTo(spoilerOpenKey{}, open, ui.Animation{Duration: m.Reveal, Ease: m.Ease})
	services := c.Services()
	return root.Children(func() {
		view := ui.Column(c).AlignSelf(ui.Stretch).MinWidth(0).AlignItems(ui.Stretch).Clip()
		// The cut moves between the clamp and the content's height, and is
		// dropped once open, so that content that grows is not held back.
		switch {
		case open <= 0:
			view.MaxHeight(limit)
		case open < 1:
			view.MaxHeight(limit + max(*height-limit, 0)*open)
		}
		// The cut fades out the last lines instead of slicing through one,
		// and clears as the content opens. The mask works on any surface.
		if overflows && open < 1 {
			black := ui.RGB(0, 0, 0)
			from := max(1-theme.Scale(c, theme.Effects().FadeEdge)/limit, 0)
			view.MaskLinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
				{Color: black, At: from * (1 - open)}, {Color: black.Alpha(open), At: 1},
			}})
		}
		view.Children(func() {
			ui.Column(c).Shrink(0).MinWidth(0).AlignItems(ui.Stretch).
				Draw(func(_ *ui.Painter, r ui.Rect) {
					was := *height
					*height = r.H
					// Only a change of fit needs the view built again.
					if (r.H > limit+spoilerSlack) != (was > limit+spoilerSlack) {
						services.Invalidate()
					}
				}).Children(body)
		})
		if !overflows {
			return
		}
		// Focus in the content opens it. MyGo reports that focus is inside,
		// not where, so a control below the cut cannot be told from one
		// above it; a control nobody can see must not hold the keyboard.
		if !*expanded && view.FocusWithin() {
			*expanded = true
			c.AnimationFrame()
		}
		label := o.ShowLabel
		if label == "" {
			label = locale.M().ShowMore
		}
		if *expanded {
			label = o.HideLabel
			if label == "" {
				label = locale.M().ShowLess
			}
		}
		// It looks like a link but shows content in place, so it is a
		// button: Space works too.
		Link(c, LinkOptions{Label: label, OnClick: func() { *expanded = !*expanded }}).
			Role(ui.RoleButton).Expanded(*expanded).Controls(view)
	})
}
