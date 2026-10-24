package form

import (
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SwitchSize is HeroUI's switch scale: 32×16, 40×20 and 48×24 DIP tracks.
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

// Switch is HeroUI's switch: a focusable row whose track and thumb follow the
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
		indent = switchMetrics(c, o.Size).w + theme.Space12
	}
	outer := ui.Column(c).MinWidth(0).Gap(theme.Space4)
	outer.Children(func() {
		row := ui.SwitchBase(c, on).Gap(theme.Space12).Shrink(0).
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
			FieldError(c, errText).PaddingStart(indent)
		}
	})
	return style.Disabled(outer, disabled)
}

// switchDimensions are the desktop metrics of one switch size.
type switchDimensions struct {
	w, h, thumbW, thumbH, radius, thumbRadius, inset float32
}

// switchMetrics resolves HeroUI's switch sizes. The upstream values are rem,
// so they follow the text scale.
func switchMetrics(c *ui.Context, size SwitchSize) switchDimensions {
	r := theme.RadiusScale()
	d := switchDimensions{
		w: theme.Scale(c, 40), h: theme.Scale(c, 20),
		thumbW: theme.Scale(c, 22), thumbH: theme.Scale(c, 16),
		radius: r.XL, thumbRadius: r.LG, inset: theme.Scale(c, 2),
	}
	switch size {
	case SwitchSmall:
		d.w, d.h = theme.Scale(c, 32), theme.Scale(c, 16)
		d.thumbW, d.thumbH = theme.Scale(c, 16.5), theme.Scale(c, 12)
		d.radius, d.thumbRadius = r.LG, r.MD
	case SwitchLarge:
		d.w, d.h = theme.Scale(c, 48), theme.Scale(c, 24)
		d.thumbW, d.thumbH = theme.Scale(c, 27.5), theme.Scale(c, 20)
		d.radius, d.thumbRadius = r.XL, r.XL
	}
	return d
}

// switchIndicator draws a switch track and thumb, animated toward on unless
// the toolkit's reduced-motion preference is set. owner supplies the row's
// live hover, pressed and focus state.
func switchIndicator(c *ui.Context, on *bool, size SwitchSize, disabled bool, owner ui.Element) ui.Element {
	t := theme.Colors(c)
	m := switchMetrics(c, size)
	rtl := c.Direction() == ui.RightToLeft
	track := ui.Box(c).Size(m.w, m.h).Shrink(0)
	pos := choiceOnOff(*on)
	if !disabled {
		pos = track.Animate("switch-thumb", pos, 200*time.Millisecond)
	}
	// Snapshot the row's states in the build: a retained copy of a leaving
	// element paints in later passes, where owner has expired.
	hovered, pressed, focused := owner.Hovered(), owner.Pressed(), owner.FocusVisible()
	track.Draw(func(p *ui.Painter, r ui.Rect) {
		fill := t.Default
		switch {
		case *on:
			fill = t.Accent
			if hovered || pressed {
				fill = t.AccentHover
			}
		case hovered || pressed:
			fill = t.DefaultHover
		}
		p.Fill(r, fill, m.radius)

		off, end := m.inset, r.W-m.inset-m.thumbW
		if rtl {
			off, end = end, off
		}
		x := off + (end-off)*pos
		thumb := ui.Rect{
			X: r.X + x, Y: r.Y + (r.H-m.thumbH)/2,
			W: m.thumbW, H: m.thumbH,
		}
		face := ui.RGB(255, 255, 255)
		if *on {
			face = t.AccentForeground
		}
		if disabled && *on {
			face = face.Alpha(0.4)
		}
		p.Shadow(thumb, m.thumbRadius, 0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.18))
		p.Fill(thumb, face, m.thumbRadius)
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
