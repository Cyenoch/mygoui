package form

import (
	"math"

	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

type SliderOptions struct {
	ControlOptions
	Min, Max, Step float64
	Vertical       bool
	Length         float32
	Unit           string
	OnChange       func(float64)
}

func Slider(c *ui.Context, value *float64, o SliderOptions) ui.Element {
	lo, hi := sliderBounds(o.Min, o.Max)
	*value = math.Max(lo, math.Min(hi, *value))
	bound := value
	if o.ReadOnly {
		shadow := *value
		bound = &shadow
	}
	colors := theme.Colors(c)
	thumbColor := colors.Background
	if c.Theme().Dark {
		thumbColor = colors.Foreground
	}
	e := ui.SliderBase(c, bound, lo, hi).Step(o.Step).FocusRing(true)
	length := o.Length
	if length <= 0 {
		length = 160
	}
	if o.Vertical {
		e.Vertical().Width(32).Height(length).PaddingY(8)
	} else {
		e.Height(32).MinWidth(length).FillWidth().PaddingX(8)
	}
	configureControl(c, e, o.ControlOptions)
	inputBorder(c, e, o.ControlOptions)
	if o.Unit != "" {
		e.Description(locales.ComponentUnit(o.Unit))
	}
	e.Value(sliderValue(c, *value, o.Unit)).OnChange(func() {
		if o.OnChange != nil && !o.ReadOnly {
			o.OnChange(*value)
		}
	})
	frac := float32((*value - lo) / (hi - lo))
	rtl := c.Direction() == ui.RightToLeft
	// inputBorder reserves one DIP inside the native content box.
	const inset = float32(9)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		center := ui.Point{X: r.X + r.W/2, Y: r.Y + r.H/2}
		if o.Vertical {
			center.Y = r.Y + inset + (r.H-2*inset)*(1-frac)
			p.Fill(ui.Rect{X: center.X - 2, Y: r.Y + inset, W: 4, H: r.H - 2*inset}, colors.Input.Alpha(0.9), 2)
			p.Fill(ui.Rect{X: center.X - 2, Y: center.Y, W: 4, H: r.Y + r.H - inset - center.Y}, colors.Primary, 2)
		} else {
			at := frac
			if rtl {
				at = 1 - frac
			}
			center.X = r.X + inset + (r.W-2*inset)*at
			track := ui.Rect{X: r.X + inset, Y: center.Y - 2, W: r.W - 2*inset, H: 4}
			p.Fill(track, colors.Input.Alpha(0.9), 2)
			filled := track
			filled.W = center.X - track.X
			if rtl {
				filled.X, filled.W = center.X, track.X+track.W-center.X
			}
			p.Fill(filled, colors.Primary, 2)
		}
		paintSliderThumb(p, center, thumbColor, colors.Foreground)
	})
	return e
}

type RangeSliderOptions struct {
	ControlOptions
	Min, Max, Step float64
	Unit           string
	OnChange       func(float64, float64)
}

// RangeSlider preserves MyGo's two independently accessible native thumbs.
// The maintained RangeSlider currently supports horizontal orientation only.
func RangeSlider(c *ui.Context, low, high *float64, o RangeSliderOptions) ui.Element {
	lo, hi := sliderBounds(o.Min, o.Max)
	*low = math.Max(lo, math.Min(hi, *low))
	*high = math.Max(*low, math.Min(hi, *high))
	left, right := low, high
	if o.ReadOnly {
		lowShadow, highShadow := *low, *high
		left, right = &lowShadow, &highShadow
	}
	colors := theme.Colors(c)
	thumbColor := colors.Background
	if c.Theme().Dark {
		thumbColor = colors.Foreground
	}
	// Match the native two-thumb hit geometry; no public RangeSlider parts or
	// thumb palette exist in the pinned renderer, so replace only its paint.
	padding := c.Theme().Space(5) / 2
	inset := padding + 1 // inputBorder's one-DIP border is outside the native content.
	e := ui.RangeSlider(c, left, right, lo, hi, o.Step).
		FillWidth().MinWidth(160).Height(32).PaddingX(padding)
	configureControl(c, e, o.ControlOptions)
	inputBorder(c, e, o.ControlOptions)
	rtl := c.Direction() == ui.RightToLeft
	lowFrac, highFrac := float32((*low-lo)/(hi-lo)), float32((*high-lo)/(hi-lo))
	if rtl {
		lowFrac, highFrac = 1-lowFrac, 1-highFrac
	}
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		cy := r.Y + r.H/2
		x0 := r.X + inset + (r.W-2*inset)*lowFrac
		x1 := r.X + inset + (r.W-2*inset)*highFrac
		p.Fill(ui.Rect{X: r.X + inset, Y: cy - 2, W: r.W - 2*inset, H: 4}, colors.Input.Alpha(0.9), 2)
		p.Fill(ui.Rect{X: min(x0, x1), Y: cy - 2, W: max(x0, x1) - min(x0, x1), H: 4}, colors.Primary, 2)
		paintSliderThumb(p, ui.Point{X: x0, Y: cy}, thumbColor, colors.Foreground)
		paintSliderThumb(p, ui.Point{X: x1, Y: cy}, thumbColor, colors.Foreground)
	})
	e.Description(locales.ComponentRangeValues(sliderValue(c, *low, o.Unit), sliderValue(c, *high, o.Unit))).
		OnChange(func() {
			if o.OnChange != nil && !o.ReadOnly {
				o.OnChange(*low, *high)
			}
		})
	if o.Unit != "" && o.Label != "" {
		e.Label(o.Label + " (" + o.Unit + ")")
	}
	return e
}

func paintSliderThumb(p *ui.Painter, center ui.Point, fill, border ui.Color) {
	thumb := ui.Rect{X: center.X - 8, Y: center.Y - 8, W: 16, H: 16}
	p.Fill(thumb, fill, 8)
	p.Stroke(thumb, border, 8, 1)
}

func sliderBounds(lo, hi float64) (float64, float64) {
	if hi <= lo {
		return 0, 100
	}
	return lo, hi
}

func sliderValue(c *ui.Context, value float64, unit string) string {
	text := i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), value, i18n.NumberOptions{MaxFractionDigits: 100})
	if unit != "" {
		return text + " " + unit
	}
	return text
}
