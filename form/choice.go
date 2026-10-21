package form

import (
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type CheckboxOptions struct {
	ControlOptions
	Mixed    *bool
	OnChange func(bool)
}

func Checkbox(c *ui.Context, value *bool, o CheckboxOptions) ui.Element {
	bound := value
	if o.ReadOnly {
		shadow := *value
		bound = &shadow
	}
	e := ui.CheckboxBase(c, bound).Gap(theme.Space8).MinHeight(theme.Metrics(c).ControlHeight).
		FontSize(theme.Scale(c, 14)).TextColor(theme.Colors(c).Foreground)
	configureControl(c, e, o.ControlOptions)
	mixed := o.Mixed != nil && *o.Mixed
	if mixed {
		e.Mixed()
	}
	e.OnChange(func() {
		if o.ReadOnly {
			return
		}
		if mixed {
			*value = true
			*o.Mixed = false
		}
		if o.OnChange != nil {
			o.OnChange(*value)
		}
	})
	colors := theme.Colors(c)
	e.Children(func() {
		box := ui.Box(c).Size(16, 16).Shrink(0).Radius(5).Background(colors.Input.Alpha(0.9))
		if *value || mixed {
			box.Background(colors.Primary)
		}
		if o.Invalid || o.Error != "" {
			box.Border(1, colors.Destructive)
		}
		box.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if mixed {
				p.Fill(ui.Rect{X: r.X + 4, Y: r.Y + 7, W: 8, H: 2}, colors.PrimaryForeground, 1)
			} else if *value {
				var path ui.Path
				path.MoveTo(r.X+3.5, r.Y+8).LineTo(r.X+6.5, r.Y+11).LineTo(r.X+12.5, r.Y+5)
				p.StrokePath(&path, 1.8, colors.PrimaryForeground)
			}
		})
		if o.Label != "" {
			ui.Text(c, o.Label)
		}
	})
	return e
}

type SwitchOptions struct {
	ControlOptions
	Small    bool
	OnChange func(bool)
}

func Switch(c *ui.Context, value *bool, o SwitchOptions) ui.Element {
	bound := value
	if o.ReadOnly {
		shadow := *value
		bound = &shadow
	}
	colors := theme.Colors(c)
	e := ui.SwitchBase(c, bound).Gap(theme.Space8).MinHeight(theme.Metrics(c).ControlHeight).
		FontSize(theme.Scale(c, 14)).TextColor(colors.Foreground)
	configureControl(c, e, o.ControlOptions)
	e.OnChange(func() {
		if o.OnChange != nil && !o.ReadOnly {
			o.OnChange(*value)
		}
	})
	w, h, thumb := float32(32), float32(20), float32(16)
	if o.Small {
		w, h, thumb = 24, 16, 12
	}
	rtl := c.Direction() == ui.RightToLeft
	on := *value
	dark := c.Theme().Dark
	e.Children(func() {
		trackColor, thumbColor := colors.Input.Alpha(0.9), colors.Background
		if on {
			trackColor = colors.Primary
		}
		if dark {
			thumbColor = colors.Foreground
			if on {
				thumbColor = colors.PrimaryForeground
			}
		}
		track := ui.Box(c.Key("track")).Size(w, h).Shrink(0).Radius(16).Background(trackColor)
		if !theme.ReduceMotion(c) {
			track.Transition(ui.ElementTransition{Colors: true, Duration: 150 * time.Millisecond})
		}
		inset := float32(2)
		if o.Invalid || o.Error != "" {
			track.Border(1, colors.Destructive)
			inset--
		}
		track.Children(func() {
			x := float32(0)
			if on != rtl {
				x = w - thumb - 4
			}
			knob := ui.Box(c.Key("thumb")).Absolute().Left(inset).Top(inset).
				Size(thumb, thumb).Radius(thumb / 2).Background(thumbColor).PassThrough()
			knob.AnimatePose(ui.Pose{X: x}, ui.Animation{Duration: 150 * time.Millisecond})
			if !theme.ReduceMotion(c) {
				knob.Transition(ui.ElementTransition{Colors: true, Duration: 150 * time.Millisecond})
			}
		})
		if o.Label != "" {
			ui.Text(c, o.Label)
		}
	})
	return e
}

type RadioGroupOptions struct {
	ControlOptions
	Horizontal bool
	OnChange   func(string)
}

type RadioItem struct {
	Value, Label, Description string
	Disabled                  bool
	Handle                    *ui.Handle
}

// RadioGroup delegates roving Tab and arrow/Home/End selection to MyGo.
func RadioGroup(c *ui.Context, value *string, o RadioGroupOptions, items []RadioItem) ui.Element {
	colors := theme.Colors(c)
	bound := value
	if o.ReadOnly {
		shadow := *value
		bound = &shadow
	}
	group := ui.RadioGroup(c, func() {
		for _, item := range items {
			c.Key(item.Value)
			e := ui.RadioBase(c, bound, item.Value).Gap(theme.Space8).MinHeight(theme.Metrics(c).ControlHeight).
				FontSize(theme.Scale(c, 14)).TextColor(colors.Foreground)
			configureControl(c, e, ControlOptions{
				Label: item.Label, Description: item.Description, Error: o.Error,
				Disabled: o.Disabled || item.Disabled, ReadOnly: o.ReadOnly,
				Invalid: o.Invalid, Loading: o.Loading, Handle: item.Handle,
			})
			e.OnChange(func() {
				if o.OnChange != nil && !o.ReadOnly {
					o.OnChange(*value)
				}
			})
			on := *value == item.Value
			e.Children(func() {
				dot := ui.Box(c).Size(16, 16).Radius(8).Shrink(0).Background(colors.Input.Alpha(0.9))
				if on {
					dot.Background(colors.Primary).Children(func() {
						ui.Box(c).Size(8, 8).Radius(4).Absolute().Left(4).Top(4).
							Background(colors.PrimaryForeground)
					})
				}
				if o.Invalid || o.Error != "" {
					dot.Border(1, colors.Destructive)
				}
				ui.Column(c).Gap(4).MinWidth(0).Children(func() {
					ui.Text(c, item.Label)
					if item.Description != "" {
						Description(c, item.Description)
					}
				})
			})
		}
	})
	group.Gap(12).FillWidth().MinWidth(0)
	if o.Horizontal {
		group.Row().Wrap()
	}
	configureControl(c, group, o.ControlOptions)
	return group
}
