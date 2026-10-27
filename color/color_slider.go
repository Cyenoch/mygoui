package color

import (
	"math"
	"strconv"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ColorSliderOptions configures one channel slider with the gradient
// track and thumb.
type ColorSliderOptions struct {
	// Label names the slider for assistive technology, such as "Hue".
	Label   string
	State   *ColorState
	Channel ColorChannel
	// ShowOutput shows the channel's current value beside the label.
	ShowOutput bool
	Vertical   bool
	Disabled   bool
	// Height is a vertical slider's length in DIPs. Zero uses 160.
	Height float32
	// Format overrides how the current value reads. It runs at build time,
	// not while painting.
	Format func(value float64) string
	// OnChange reports that the user moved the slider.
	OnChange func()
}

// ColorSlider is a single-channel color slider. It paints a gradient
// of the channel's range with a draggable thumb; arrows move by one unit (ten
// with Shift), and Home and End go to the ends.
func ColorSlider(c *ui.Context, o ColorSliderOptions) ui.Element {
	if c == nil || o.State == nil {
		return ui.Element{}
	}
	s := o.State
	ch := o.Channel
	value, lo, hi := s.ChannelValue(ch)
	format := o.Format
	if format == nil {
		format = func(v float64) string {
			n := strconv.FormatInt(int64(math.Round(v)), 10)
			switch {
			case ch.Percent():
				return n + "%"
			case ch == ChannelHue:
				return n + "°"
			}
			return n
		}
	}
	label := o.Label
	if label == "" {
		label = ch.String()
	}
	t := theme.Colors(c)
	ink := t.Foreground
	thumb := float32(16)
	rtl := c.Direction() == ui.RightToLeft
	stops := trackStops(s, ch)
	alpha := ch == ChannelAlpha || s.Value.A < 255
	valueColor := s.Value
	angle := gradientAngle(o.Vertical, rtl)

	slider := ui.Column(c).MinWidth(0)
	if o.Vertical {
		slider.Shrink(0).AlignItems(ui.Center).Gap(theme.Space8)
	} else {
		slider.FillWidth().Gap(theme.Space4)
	}
	slider.Children(func() {
		if o.Vertical {
			if o.ShowOutput {
				style.Text(ui.Text(c, format(value)), c, theme.Label()).
					TextColor(ink).SingleLine().FontFeatures("tnum")
			}
			colorSliderTrack(c, o, label, lo, hi, value, thumb, stops, alpha,
				valueColor, angle, rtl, s, format)
			style.Text(ui.Text(c, label), c, theme.Label()).TextColor(ink)
			return
		}
		if o.ShowOutput {
			ui.Row(c).FillWidth().AlignItems(ui.Center).
				Justify(ui.SpaceBetween).Children(func() {
				style.Text(ui.Text(c, label), c, theme.Label()).TextColor(ink)
				style.Text(ui.Text(c, format(value)), c, theme.Label()).
					TextColor(ink).SingleLine().FontFeatures("tnum")
			})
		} else {
			style.Text(ui.Text(c, label), c, theme.Label()).TextColor(ink)
		}
		colorSliderTrack(c, o, label, lo, hi, value, thumb, stops, alpha,
			valueColor, angle, rtl, s, format)
	})
	return style.Disabled(slider, o.Disabled)
}

// colorSliderTrack builds the focusable gradient track and its thumb.
func colorSliderTrack(c *ui.Context, o ColorSliderOptions, label string,
	lo, hi, value float64, thumb float32, stops []ui.GradientStop, alpha bool,
	valueColor ui.Color, angle float32, rtl bool, s *ColorState,
	format func(float64) string) ui.Element {
	ch := o.Channel
	var handle ui.Handle
	track := ui.Box(c).Focusable().FocusRing(false).Role(ui.RoleSlider).
		Label(label).Range(lo, hi, value).Step(1).Cursor(ui.CursorGrab).
		Disabled(o.Disabled).Bind(&handle)
	var radius float32
	if o.Vertical {
		height := o.Height
		if height <= 0 {
			height = 160
		}
		track.Size(20, height)
		radius = 10
	} else {
		track.Height(20).FillWidth()
		radius = 10
	}
	set := func(v float64) {
		before := s.Value
		s.SetChannel(ch, v)
		if o.OnChange != nil && s.Value != before {
			o.OnChange()
		}
	}
	track = track.HandleInput(func(ev ui.InputEvent) bool {
		if o.Disabled {
			return false
		}
		switch ev.Kind {
		case ui.InputPointerDown:
			if ev.Button != 0 {
				return false
			}
			r := handle.Bounds(c)
			set(fromPointer(ev.X, ev.Y, r.W, r.H, thumb, lo, hi, o.Vertical, rtl))
			return true
		case ui.InputPointerMove:
			if ev.Button != 0 {
				return false
			}
			r := handle.Bounds(c)
			set(fromPointer(ev.X, ev.Y, r.W, r.H, thumb, lo, hi, o.Vertical, rtl))
			return true
		case ui.InputKeyDown:
			step := 1.0
			if ev.Mods&ui.Shift != 0 {
				step = 10
			}
			cur, _, _ := s.ChannelValue(ch)
			if o.Vertical {
				switch ev.Key {
				case ui.KeyLeft:
					set(cur - step)
					return true
				case ui.KeyRight:
					set(cur + step)
					return true
				}
			}
			back, forth := ui.KeyLeft, ui.KeyRight
			if o.Vertical {
				back, forth = ui.KeyDown, ui.KeyUp
			} else if rtl {
				back, forth = forth, back
			}
			switch ev.Key {
			case back:
				set(cur - step)
			case forth:
				set(cur + step)
			case ui.KeyPageDown:
				set(cur - 10)
			case ui.KeyPageUp:
				set(cur + 10)
			case ui.KeyHome:
				set(lo)
			case ui.KeyEnd:
				set(hi)
			default:
				return false
			}
			return true
		}
		return false
	})
	// Snapshot the track's states in the build: a retained copy of a
	// leaving element paints in later passes, where track has expired.
	pressed, focused := track.Pressed(), track.FocusVisible()
	if pressed {
		track.Cursor(ui.CursorGrabbing)
	}
	f := float32(0)
	if hi > lo {
		f = float32((value - lo) / (hi - lo))
	}
	track.Draw(func(p *ui.Painter, r ui.Rect) {
		if alpha {
			p.Clip(r, radius, func() { checkers(p, r, 8) })
		}
		// The control insets the gradient by 10 DIP at each end and draws the
		// channel's extreme colors as caps, so the thumb travels between the
		// cap centers rather than to the bare edges.
		inset := thumb/2 + theme.Space2
		inner := r
		var lowCap, highCap ui.Rect
		var lowRadii, highRadii [4]float32
		capRadius := inset
		if o.Vertical {
			inner = ui.Rect{X: r.X, Y: r.Y + inset, W: r.W, H: r.H - 2*inset}
			lowCap = ui.Rect{X: r.X, Y: r.Y + r.H - inset, W: r.W, H: inset}
			highCap = ui.Rect{X: r.X, Y: r.Y, W: r.W, H: inset}
			lowRadii = [4]float32{0, 0, capRadius, capRadius}
			highRadii = [4]float32{capRadius, capRadius, 0, 0}
		} else {
			inner = ui.Rect{X: r.X + inset, Y: r.Y, W: r.W - 2*inset, H: r.H}
			lowCap = ui.Rect{X: r.X, Y: r.Y, W: inset, H: r.H}
			highCap = ui.Rect{X: r.X + r.W - inset, Y: r.Y, W: inset, H: r.H}
			lowRadii = [4]float32{capRadius, 0, 0, capRadius}
			highRadii = [4]float32{0, capRadius, capRadius, 0}
		}
		low, high := lowColor(stops, valueColor), highColor(stops, valueColor)
		if !o.Vertical && rtl {
			low, high = high, low
			lowCap, highCap = highCap, lowCap
			lowRadii, highRadii = highRadii, lowRadii
		}
		if inner.W > 0 && inner.H > 0 {
			p.FillGradient(inner, ui.LinearGradient{Stops: stops, Angle: angle}, 0)
		}
		fillCorners(p, lowCap, low, lowRadii[0], lowRadii[1], lowRadii[2], lowRadii[3])
		fillCorners(p, highCap, high, highRadii[0], highRadii[1], highRadii[2], highRadii[3])
		p.Stroke(r, ui.RGBA(0, 0, 0, 0.1), radius, 1)

		var cx, cy float32
		if o.Vertical {
			cy = r.Y + inset + (1-f)*(r.H-2*inset)
			cx = r.X + r.W/2
		} else {
			at := f
			if rtl {
				at = 1 - f
			}
			cx = r.X + inset + at*(r.W-2*inset)
			cy = r.Y + r.H/2
		}
		paintThumb(p, cx, cy, thumb, valueColor)
		if focused {
			half := thumb / 2
			p.FocusRing(ui.Rect{X: cx - half, Y: cy - half, W: thumb, H: thumb},
				[4]float32{half, half, half, half})
		}
	})
	return track.Value(format(value))
}

// lowColor and highColor are the track's extreme colors, for its caps: the
// first and last stop, or the current color when a track has no stops.
func lowColor(stops []ui.GradientStop, fallback ui.Color) ui.Color {
	if len(stops) == 0 {
		return fallback
	}
	return stops[0].Color
}

func highColor(stops []ui.GradientStop, fallback ui.Color) ui.Color {
	if len(stops) == 0 {
		return fallback
	}
	return stops[len(stops)-1].Color
}

// gradientAngle points from the low end of the track to the high end: right
// in a left-to-right horizontal slider, left in a right-to-left one, and up
// a vertical one.
func gradientAngle(vertical, rtl bool) float32 {
	if vertical {
		return 0
	}
	if rtl {
		return 270
	}
	return 90
}

// fromPointer maps a pointer position inside a track to a channel value. The
// thumb's center travels between the cap centers, the half thumb plus a
// 2 DIP gap, so the ends are reachable at the track's outer edge.
func fromPointer(x, y, width, height, thumb float32, lo, hi float64, vertical, rtl bool) float64 {
	inset := thumb/2 + 2
	f := float32(0)
	if vertical {
		usable := height - 2*inset
		if usable > 0 {
			f = 1 - (y-inset)/usable
		}
	} else {
		usable := width - 2*inset
		if usable > 0 {
			f = (x - inset) / usable
		}
		if rtl {
			f = 1 - f
		}
	}
	f = float32(clamp(float64(f), 0, 1))
	return lo + float64(f)*(hi-lo)
}

// hueStops is the fixed spectrum of a hue slider.
var hueStops = []ui.GradientStop{
	{Color: hsvColor(0, 1, 1), At: 0},
	{Color: hsvColor(60, 1, 1), At: 1.0 / 6},
	{Color: hsvColor(120, 1, 1), At: 2.0 / 6},
	{Color: hsvColor(180, 1, 1), At: 3.0 / 6},
	{Color: hsvColor(240, 1, 1), At: 4.0 / 6},
	{Color: hsvColor(300, 1, 1), At: 5.0 / 6},
	{Color: hsvColor(360, 1, 1), At: 1},
}

// trackStops returns the gradient of a channel's range at the state's other
// channels. It runs at build time; the paint callback only reads the slice.
func trackStops(s *ColorState, ch ColorChannel) []ui.GradientStop {
	switch ch {
	case ChannelHue:
		return hueStops
	case ChannelSaturation:
		return []ui.GradientStop{
			{Color: hsvColor(s.Hue(), 0, s.Brightness()), At: 0},
			{Color: hsvColor(s.Hue(), 1, s.Brightness()), At: 1},
		}
	case ChannelBrightness:
		return []ui.GradientStop{
			{Color: hsvColor(s.Hue(), s.Saturation(), 0), At: 0},
			{Color: hsvColor(s.Hue(), s.Saturation(), 1), At: 1},
		}
	case ChannelLightness:
		h, sat, _ := rgbToHSL(s.Value)
		return []ui.GradientStop{
			{Color: hslColor(h, sat, 0), At: 0},
			{Color: hslColor(h, sat, 0.5), At: 0.5},
			{Color: hslColor(h, sat, 1), At: 1},
		}
	case ChannelRed:
		lo, hi := s.Value, s.Value
		lo.R, hi.R = 0, 255
		return []ui.GradientStop{{Color: lo, At: 0}, {Color: hi, At: 1}}
	case ChannelGreen:
		lo, hi := s.Value, s.Value
		lo.G, hi.G = 0, 255
		return []ui.GradientStop{{Color: lo, At: 0}, {Color: hi, At: 1}}
	case ChannelBlue:
		lo, hi := s.Value, s.Value
		lo.B, hi.B = 0, 255
		return []ui.GradientStop{{Color: lo, At: 0}, {Color: hi, At: 1}}
	case ChannelAlpha:
		clear, solid := s.Value, s.Value
		clear.A = 0
		solid.A = 255
		return []ui.GradientStop{{Color: clear, At: 0}, {Color: solid, At: 1}}
	}
	return nil
}
