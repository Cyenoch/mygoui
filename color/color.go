// Package color implements HeroUI v3's color controls with native MyGo
// painting, keyboard handling and accessibility. Callers own every value; the
// package performs no I/O and builds no images.
package color

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"
)

// ColorChannel is one axis of a color a slider edits. Hue runs from 0 to 360
// degrees; saturation, brightness, lightness and alpha run from 0 to 100; the
// red, green and blue channels run from 0 to 255.
type ColorChannel uint8

// Color channels.
const (
	ChannelHue ColorChannel = iota
	ChannelSaturation
	ChannelBrightness
	ChannelLightness
	ChannelRed
	ChannelGreen
	ChannelBlue
	ChannelAlpha
)

// String returns the channel's technical name.
func (ch ColorChannel) String() string {
	switch ch {
	case ChannelHue:
		return "Hue"
	case ChannelSaturation:
		return "Saturation"
	case ChannelBrightness:
		return "Brightness"
	case ChannelLightness:
		return "Lightness"
	case ChannelRed:
		return "Red"
	case ChannelGreen:
		return "Green"
	case ChannelBlue:
		return "Blue"
	case ChannelAlpha:
		return "Opacity"
	}
	return "Color"
}

// Range returns the channel's inclusive bounds.
func (ch ColorChannel) Range() (lo, hi float64) {
	switch ch {
	case ChannelHue:
		return 0, 360
	case ChannelRed, ChannelGreen, ChannelBlue:
		return 0, 255
	}
	return 0, 100
}

// Percent reports whether the channel's values read as percentages.
func (ch ColorChannel) Percent() bool {
	switch ch {
	case ChannelHue, ChannelRed, ChannelGreen, ChannelBlue:
		return false
	}
	return true
}

// ColorState is the caller-owned color a family of color controls shares. It
// keeps the hue and saturation of achromatic colors, so dragging brightness to
// zero and back restores the chosen hue rather than producing white.
//
// The zero value is black. Construct one with NewColorState.
type ColorState struct {
	// Value is the current sRGB color, with straight alpha. Set it through
	// Set, the channel setters or SetHex so the HSV members follow; a direct
	// assignment is reconciled on the next read.
	Value ui.Color

	hue, saturation, brightness float64
	seen                        ui.Color
	ready                       bool
	field                       colorFieldState
}

// sync reconciles a direct assignment to Value with the HSV members.
func (s *ColorState) sync() {
	if !s.ready {
		s.remember()
		s.seen = s.Value
		s.ready = true
		return
	}
	if s.Value != s.seen {
		s.Set(s.Value)
	}
}

type colorFieldState struct {
	text    string
	seen    ui.Color
	ready   bool
	invalid bool
}

// NewColorState returns a state holding v.
func NewColorState(v ui.Color) *ColorState {
	s := &ColorState{Value: v}
	s.remember()
	s.seen = v
	s.ready = true
	return s
}

// remember derives the HSV members from Value. A black value keeps the
// saturation, and any achromatic value keeps the hue, so a channel can return
// to where it came from.
func (s *ColorState) remember() {
	h, sat, val := rgbToHSV(s.Value)
	if s.Value.R == s.Value.G && s.Value.G == s.Value.B {
		// Achromatic: keep the hue that produced this gray, and the
		// saturation under black, so the channel can return.
		h = s.hue
		if val > 0 {
			sat = 0
		} else {
			sat = s.saturation
		}
	}
	s.hue, s.saturation, s.brightness = h, sat, val
}

// Set replaces the color, deriving its hue, saturation and brightness. An
// external gray or black cannot express a hue; the previous one is kept.
func (s *ColorState) Set(v ui.Color) {
	s.Value = v
	s.remember()
	s.seen = v
	s.ready = true
}

// Hue returns the hue in degrees, from 0 to 360.
func (s *ColorState) Hue() float64 { s.sync(); return s.hue }

// Saturation returns the HSV saturation, from 0 to 1.
func (s *ColorState) Saturation() float64 { s.sync(); return s.saturation }

// Brightness returns the HSV value, from 0 to 1.
func (s *ColorState) Brightness() float64 { s.sync(); return s.brightness }

// Alpha returns the opacity, from 0 to 1.
func (s *ColorState) Alpha() float64 { s.sync(); return float64(s.Value.A) / 255 }

// SetHue sets the hue in degrees and recomputes the color.
func (s *ColorState) SetHue(h float64) {
	s.sync()
	s.hue = clamp(h, 0, 360)
	s.Value = s.compose()
	s.seen = s.Value
}

// SetSaturation sets the HSV saturation from 0 to 1.
func (s *ColorState) SetSaturation(v float64) {
	s.sync()
	s.saturation = clamp(v, 0, 1)
	s.Value = s.compose()
	s.seen = s.Value
}

// SetBrightness sets the HSV value from 0 to 1.
func (s *ColorState) SetBrightness(v float64) {
	s.sync()
	s.brightness = clamp(v, 0, 1)
	s.Value = s.compose()
	s.seen = s.Value
}

// SetAlpha sets the opacity from 0 to 1.
func (s *ColorState) SetAlpha(a float64) {
	s.sync()
	s.Value.A = uint8(math.Round(clamp(a, 0, 1) * 255))
	s.seen = s.Value
}

// compose returns the color of the HSV members with the current alpha.
func (s *ColorState) compose() ui.Color {
	c := hsvColor(s.hue, s.saturation, s.brightness)
	c.A = s.Value.A
	return c
}

// ChannelValue returns the channel's current value and bounds.
func (s *ColorState) ChannelValue(ch ColorChannel) (value, lo, hi float64) {
	s.sync()
	lo, hi = ch.Range()
	switch ch {
	case ChannelHue:
		return s.hue, lo, hi
	case ChannelSaturation:
		return s.saturation * 100, lo, hi
	case ChannelBrightness:
		return s.brightness * 100, lo, hi
	case ChannelLightness:
		_, _, l := rgbToHSL(s.Value)
		return l * 100, lo, hi
	case ChannelRed:
		return float64(s.Value.R), lo, hi
	case ChannelGreen:
		return float64(s.Value.G), lo, hi
	case ChannelBlue:
		return float64(s.Value.B), lo, hi
	case ChannelAlpha:
		return s.Alpha() * 100, lo, hi
	}
	return 0, lo, hi
}

// SetChannel sets one channel, leaving the others as they are.
func (s *ColorState) SetChannel(ch ColorChannel, v float64) {
	s.sync()
	lo, hi := ch.Range()
	v = clamp(v, lo, hi)
	switch ch {
	case ChannelHue:
		s.SetHue(v)
	case ChannelSaturation:
		s.SetSaturation(v / 100)
	case ChannelBrightness:
		s.SetBrightness(v / 100)
	case ChannelLightness:
		h, sat, _ := rgbToHSL(s.Value)
		c := hslColor(h, sat, v/100)
		c.A = s.Value.A
		s.Set(c)
	case ChannelRed:
		s.Value.R = uint8(math.Round(v))
		s.Set(s.Value)
	case ChannelGreen:
		s.Value.G = uint8(math.Round(v))
		s.Set(s.Value)
	case ChannelBlue:
		s.Value.B = uint8(math.Round(v))
		s.Set(s.Value)
	case ChannelAlpha:
		s.SetAlpha(v / 100)
	}
}

// Hex returns the color as "#rrggbb", or "#rrggbbaa" when translucent.
func (s *ColorState) Hex() string { return FormatHex(s.Value) }

// SetHex parses a hex color into the state.
func (s *ColorState) SetHex(text string) error {
	c, err := ParseHex(text)
	if err != nil {
		return err
	}
	s.Set(c)
	return nil
}

// ParseHex parses "#rgb", "#rgba", "#rrggbb" and "#rrggbbaa", with or without
// the leading hash. It reports an error rather than panicking on other input.
func ParseHex(text string) (ui.Color, error) {
	h := strings.TrimSpace(text)
	h = strings.TrimPrefix(h, "#")
	var v [8]uint8
	for i := 0; i < len(h) && i < len(v); i++ {
		d, ok := hexDigit(h[i])
		if !ok {
			return ui.Color{}, fmt.Errorf("color: invalid hex color %q", text)
		}
		v[i] = d
	}
	switch len(h) {
	case 3, 4:
		a := uint8(15)
		if len(h) == 4 {
			a = v[3]
		}
		return ui.Color{R: v[0] * 17, G: v[1] * 17, B: v[2] * 17, A: a * 17}, nil
	case 6, 8:
		a := uint8(255)
		if len(h) == 8 {
			a = v[6]<<4 | v[7]
		}
		return ui.Color{R: v[0]<<4 | v[1], G: v[2]<<4 | v[3], B: v[4]<<4 | v[5], A: a}, nil
	}
	return ui.Color{}, fmt.Errorf("color: invalid hex color %q", text)
}

// FormatHex returns "#rrggbb", or "#rrggbbaa" when the color is translucent.
func FormatHex(c ui.Color) string {
	if c.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

func hexDigit(c byte) (uint8, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

func clamp(v, lo, hi float64) float64 { return max(lo, min(hi, v)) }

// hsvColor converts HSV, hue in degrees and the rest from 0 to 1, to an opaque
// sRGB color.
func hsvColor(h, s, v float64) ui.Color {
	h = math.Mod(math.Mod(h, 360)+360, 360)
	ch := v * s
	x := ch * (1 - math.Abs(math.Mod(h/60, 2)-1))
	var r, g, b float64
	switch {
	case h < 60:
		r, g = ch, x
	case h < 120:
		r, g = x, ch
	case h < 180:
		g, b = ch, x
	case h < 240:
		g, b = x, ch
	case h < 300:
		r, b = x, ch
	default:
		r, b = ch, x
	}
	m := v - ch
	return ui.Color{
		R: uint8(math.Round(clamp(r+m, 0, 1) * 255)),
		G: uint8(math.Round(clamp(g+m, 0, 1) * 255)),
		B: uint8(math.Round(clamp(b+m, 0, 1) * 255)),
		A: 255,
	}
}

// rgbToHSV converts a color to hue in degrees, saturation and brightness.
func rgbToHSV(c ui.Color) (h, s, v float64) {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	hi, lo := max(r, max(g, b)), min(r, min(g, b))
	d := hi - lo
	switch {
	case d == 0:
	case hi == r:
		h = 60 * math.Mod((g-b)/d, 6)
	case hi == g:
		h = 60 * ((b-r)/d + 2)
	default:
		h = 60 * ((r-g)/d + 4)
	}
	if h < 0 {
		h += 360
	}
	if hi > 0 {
		s = d / hi
	}
	return h, s, hi
}

// rgbToHSL converts a color to hue, saturation and lightness, each fraction
// except the hue.
func rgbToHSL(c ui.Color) (h, s, l float64) {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	hi, lo := max(r, max(g, b)), min(r, min(g, b))
	d := hi - lo
	l = (hi + lo) / 2
	switch {
	case d == 0:
	case hi == r:
		h = 60 * math.Mod((g-b)/d, 6)
	case hi == g:
		h = 60 * ((b-r)/d + 2)
	default:
		h = 60 * ((r-g)/d + 4)
	}
	if h < 0 {
		h += 360
	}
	if d != 0 {
		s = d / (1 - math.Abs(2*l-1))
	}
	return h, s, l
}

// hslColor converts HSL, hue in degrees and the rest from 0 to 1, to an opaque
// sRGB color.
func hslColor(h, s, l float64) ui.Color {
	h = math.Mod(math.Mod(h, 360)+360, 360)
	ch := (1 - math.Abs(2*l-1)) * s
	x := ch * (1 - math.Abs(math.Mod(h/60, 2)-1))
	var r, g, b float64
	switch {
	case h < 60:
		r, g = ch, x
	case h < 120:
		r, g = x, ch
	case h < 180:
		g, b = ch, x
	case h < 240:
		g, b = x, ch
	case h < 300:
		r, b = x, ch
	default:
		r, b = ch, x
	}
	m := l - ch/2
	return ui.Color{
		R: uint8(math.Round(clamp(r+m, 0, 1) * 255)),
		G: uint8(math.Round(clamp(g+m, 0, 1) * 255)),
		B: uint8(math.Round(clamp(b+m, 0, 1) * 255)),
		A: 255,
	}
}
