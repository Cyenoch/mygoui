package color

import (
	"math"
	"strconv"
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ColorAreaOptions configures the two-dimensional color area. X chooses
// saturation and Y chooses brightness, in the hue kept by State.
type ColorAreaOptions struct {
	// Label names the area for assistive technology, such as "Fill color".
	Label string
	// State is the color edited and shown.
	State *ColorState
	// ShowDots overlays the subtle dot pattern.
	ShowDots bool
	Disabled bool
	// Size is the square's side in DIPs. Zero uses the 224 DIP maximum.
	Size float32
	// OnChange reports that the user picked a new color.
	OnChange func()
}

// ColorArea is a saturation/brightness area: a square gradient with
// a draggable thumb. Arrows move saturation and brightness by one percent, ten
// with Shift; Home and End go to the extremes of saturation.
func ColorArea(c *ui.Context, o ColorAreaOptions) ui.Element {
	if c == nil || o.State == nil {
		return ui.Element{}
	}
	s := o.State
	size := o.Size
	if size <= 0 {
		size = 224
	}
	radius := theme.RadiusScale().XXL
	label := o.Label
	if label == "" {
		label = "Color area"
	}
	var areaHandle ui.Handle
	area := ui.Box(c).FillWidth().MaxWidth(size).AspectRatio(1).Shrink(0).Radius(radius).
		Focusable().FocusRing(false).Role(ui.RoleSlider).
		Label(label).Range(0, 100, s.Saturation()*100).Step(1).
		Cursor(ui.CursorCrosshair).Bind(&areaHandle)
	area = style.Disabled(area, o.Disabled)

	notify := func(before ui.Color) {
		if o.OnChange != nil && s.Value != before {
			o.OnChange()
		}
	}
	set := func(x, y float32, width, height float32) {
		before := s.Value
		if width > 0 {
			s.SetSaturation(float64(clamp(float64(x/width), 0, 1)))
		}
		if height > 0 {
			s.SetBrightness(1 - float64(clamp(float64(y/height), 0, 1)))
		}
		notify(before)
	}

	area = area.HandleInput(func(ev ui.InputEvent) bool {
		if o.Disabled {
			return false
		}
		switch ev.Kind {
		case ui.InputPointerDown:
			if ev.Button != 0 {
				return false
			}
			r := areaHandle.Bounds(c)
			set(ev.X, ev.Y, r.W, r.H)
			return true
		case ui.InputPointerMove:
			if ev.Button != 0 {
				return false
			}
			r := areaHandle.Bounds(c)
			set(ev.X, ev.Y, r.W, r.H)
			return true
		case ui.InputKeyDown:
			step := 0.01
			if ev.Mods&ui.Shift != 0 {
				step = 0.1
			}
			before := s.Value
			switch ev.Key {
			case ui.KeyLeft:
				s.SetSaturation(s.Saturation() - step)
			case ui.KeyRight:
				s.SetSaturation(s.Saturation() + step)
			case ui.KeyDown:
				s.SetBrightness(s.Brightness() - step)
			case ui.KeyUp:
				s.SetBrightness(s.Brightness() + step)
			case ui.KeyHome:
				s.SetSaturation(0)
			case ui.KeyEnd:
				s.SetSaturation(1)
			default:
				return false
			}
			notify(before)
			return true
		}
		return false
	})

	sat, bright := s.Saturation(), s.Brightness()
	area = area.Value(strconv.FormatInt(int64(math.Round(sat*100)), 10) + "%, " +
		strconv.FormatInt(int64(math.Round(bright*100)), 10) + "%")
	hue := hsvColor(s.Hue(), 1, 1)
	white := ui.RGB(255, 255, 255)
	clear := ui.RGBA(0, 0, 0, 0)
	black := ui.RGB(0, 0, 0)
	// Snapshot the area's states in the build: a retained copy of a
	// leaving element paints in later passes, where area has expired.
	pressed, focused := area.Pressed(), area.FocusVisible()
	thumbColor := s.Value
	thumbSize := theme.Scale(c, 16)
	area.Draw(func(p *ui.Painter, r ui.Rect) {
		p.FillGradient(r, ui.LinearGradient{From: white, To: hue, Angle: 90}, radius)
		p.FillGradient(r, ui.LinearGradient{From: clear, To: black, Angle: 180}, radius)
		p.Stroke(r, ui.RGBA(0, 0, 0, 0.1), radius, 1)
		if o.ShowDots {
			const gap = 8
			dot := ui.RGBA(255, 255, 255, 0.2)
			p.Clip(r, radius, func() {
				for y := r.Y; y < r.Y+r.H; y += gap {
					for x := r.X; x < r.X+r.W; x += gap {
						// The radial dot is 2 DIP across, one pixel of
						// body each side of the 8 DIP grid point.
						p.Fill(ui.Rect{X: x - 1, Y: y - 1, W: 2, H: 2}, dot, 1)
					}
				}
			})
		}
	})
	// The thumb is its own element so the native pose can grow it on press
	// (16 → 20 DIP, 150 ms) in paint-only frames. Tracking it as a 1 DIP
	// anchor keeps the pointer position exact and the scale centered.
	scale := float32(1)
	if pressed {
		scale = 1.25
	}
	area.Children(func() {
		thumb := ui.Box(c.Key(colorAreaThumbKey{})).Absolute().Role(ui.RoleNone).
			Size(1, 1).LeftPercent(float32(sat * 100)).TopPercent(float32((1 - bright) * 100))
		thumb.AnimatePose(ui.Pose{Scale: scale},
			ui.Animation{Duration: 150 * time.Millisecond, Ease: ui.EaseOut})
		thumb.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			paintThumb(p, cx, cy, thumbSize, thumbColor)
			if focused {
				half := thumbSize / 2
				p.FocusRing(ui.Rect{X: cx - half, Y: cy - half, W: thumbSize, H: thumbSize},
					[4]float32{half, half, half, half})
			}
		})
	})
	return area
}

type colorAreaThumbKey struct{}
