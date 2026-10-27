package form

import (
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SwitchSize selects the switch's desktop track and thumb size.
type SwitchSize uint8

const (
	SwitchMedium SwitchSize = iota
	SwitchSmall
	SwitchLarge
)

// SwitchLabelPosition places the label before or after the switch control.
type SwitchLabelPosition uint8

const (
	SwitchLabelAfter SwitchLabelPosition = iota
	SwitchLabelBefore
)

// SwitchOptions configures a native switch. On is caller-owned: a click or
// Space flips it and OnChange reports the new value. A nil On renders a
// disabled switch.
type SwitchOptions struct {
	On            *bool
	Label         string
	Description   string
	Error         string
	Size          SwitchSize
	LabelPosition SwitchLabelPosition
	Required      bool
	Invalid       bool
	Disabled      bool
	OnChange      func(bool)
	Handle        *ui.Handle
}

// Switch is a focusable row whose track and thumb follow the
// caller's bool. The thumb's travel and colors animate unless the desktop
// asks for reduced motion, which moves it at once.
func Switch(c *ui.Context, o SwitchOptions) ui.Element {
	on := o.On
	if on == nil {
		var off bool
		on = &off
		o.Disabled = true
	}
	t := theme.Colors(c)
	disabled := o.Disabled
	errText := choiceErrorText(o.Invalid, o.Error)

	indent := float32(0)
	if o.LabelPosition == SwitchLabelAfter {
		indent = switchMetrics(c, o.Size).w + theme.Space8
	}
	outer := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	outer.Children(func() {
		row := ui.SwitchBase(c, on).Gap(theme.Space8).Shrink(0).
			FocusRing(false).TextColor(t.Foreground).Label(o.Label)
		if o.Handle != nil {
			row.Bind(o.Handle)
		}
		if o.Description != "" {
			row.Description(o.Description)
		}
		if errText != "" {
			row.Error(errText)
		}
		if o.OnChange != nil && o.On != nil && !disabled {
			row.OnChange(func() { o.OnChange(*on) })
		}

		row.Children(func() {
			if o.LabelPosition == SwitchLabelBefore && o.Label != "" {
				Label(c, o.Label, LabelOptions{
					Required: o.Required, Invalid: o.Invalid,
				})
			}
			switchIndicator(c, on, o.Size, disabled, row)
			if o.LabelPosition == SwitchLabelAfter && o.Label != "" {
				Label(c, o.Label, LabelOptions{
					Required: o.Required, Invalid: o.Invalid,
				})
			}
		})
		if o.Description != "" {
			Description(c, o.Description).PaddingStart(indent)
		}
		if errText != "" {
			choiceErrorMessage(c, errText).PaddingStart(indent)
		}
	})
	return style.Disabled(outer, disabled)
}

// switchDimensions are the desktop metrics of one switch size.
type switchDimensions struct {
	w, h, thumbW, thumbH, radius, thumbRadius, inset float32
}

// switchMetrics resolves the Coss desktop switch and its optional sizes.
func switchMetrics(c *ui.Context, size SwitchSize) switchDimensions {
	thumb := float32(16)
	switch size {
	case SwitchSmall:
		thumb = 12
	case SwitchLarge:
		thumb = 20
	}
	thumb = theme.Scale(c, thumb)
	inset := theme.Scale(c, 1)
	return switchDimensions{
		w: 2*thumb - 2*inset, h: thumb + 2*inset,
		thumbW: thumb, thumbH: thumb, radius: (thumb + 2*inset) / 2,
		thumbRadius: thumb / 2, inset: inset,
	}
}

// switchIndicator paints the native switch's Coss track and round thumb.
func switchIndicator(c *ui.Context, on *bool, size SwitchSize, disabled bool, owner ui.Element) ui.Element {
	t := theme.Colors(c)
	m := switchMetrics(c, size)
	rtl := c.Direction() == ui.RightToLeft
	track := ui.Box(c).Size(m.w, m.h).Shrink(0)
	pos := choiceOnOff(*on)
	if !disabled && !theme.ReduceMotion(c) {
		pos = track.Animate("switch-thumb", pos, 150*time.Millisecond)
	}
	focused, pressed := owner.FocusVisible(), owner.Pressed()
	fill := t.Input
	if *on {
		fill = t.Primary
	}
	if !theme.ReduceMotion(c) {
		fill = track.AnimateColor("switch-track", fill, ui.Animation{Duration: 200 * time.Millisecond, Ease: ui.EaseOut})
	}
	track.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, fill, m.radius)
		off, end := m.inset, r.W-m.inset-m.thumbW
		if rtl {
			off, end = end, off
		}
		w := m.thumbW
		if pressed && !disabled {
			w *= 1.1
		}
		x := off + (end-off)*pos
		if *on && !rtl || !*on && rtl {
			x -= w - m.thumbW
		}
		thumb := ui.Rect{X: r.X + x, Y: r.Y + m.inset, W: w, H: m.thumbH}
		p.Shadow(thumb, m.thumbRadius, 0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.05))
		p.Fill(thumb, t.Background, m.thumbRadius)
		if focused {
			p.FocusRing(r, choiceRadii(m.radius))
		}
	})
	return track
}

// choiceOnOff returns 1 when on, else 0.
func choiceOnOff(on bool) float32 {
	if on {
		return 1
	}
	return 0
}
