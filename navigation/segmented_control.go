package navigation

import (
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SegmentedControlOptions configures a single-choice segmented control.
// Content on an item is ignored; it is intended for Tabs instead.
type SegmentedControlOptions struct {
	Selected *int
	Items    []TabOptions
	Label    string
	OnChange func(index int)
}

// SegmentedControl builds a radio group whose choice is caller-owned when
// Selected is supplied. Arrow keys move the native focus and selection. The
// selected surface moves and resizes over 200 ms, or snaps for reduced motion.
func SegmentedControl(c *ui.Context, o SegmentedControlOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Box(c).Shrink(0).MinWidth(0)
	selected := o.Selected
	if selected == nil {
		selected = ui.Local(root, "segment-selected", func() int { return 0 })
	}
	boxes := ui.Local(root, "segment-boxes", func() []ui.Rect { return make([]ui.Rect, len(o.Items)) })
	indicator := ui.Local(root, "segment-indicator", func() tabsIndicator {
		return tabsIndicator{seen: -1, sel: -1}
	})
	if len(*boxes) != len(o.Items) {
		*boxes = make([]ui.Rect, len(o.Items))
		*indicator = tabsIndicator{seen: -1, sel: -1}
	}
	root.Children(func() {
		parts := ui.SegmentedBase(c, selected, len(o.Items))
		track := parts.Track.AlignItems(ui.Center).Padding(scale(c, 2)).Gap(scale(c, 2)).
			Radius(theme.RadiusScale().LG).Background(t.MutedBackground).Label(o.Label)
		trackBox := track.LayoutBounds()
		placed := *selected
		known := placed >= 0 && placed < len(*boxes) && (*boxes)[placed].W > 0
		face, radius := segmentBackground(c, t), theme.RadiusScale().MD
		shadows := theme.SurfaceShadows(c.Theme().Dark)
		var paintedTrack ui.Rect
		track.Children(func() {
			for i, item := range o.Items {
				segment := parts.Segment(i).Height(scale(c, 30)).PaddingX(scale(c, 10)).
					Radius(theme.RadiusScale().MD).FocusRing(true).Label(item.Label)
				fg := t.MutedForeground.Alpha(0.72)
				if !item.Disabled && segment.Hovered() {
					fg = t.MutedForeground
				}
				if i == *selected {
					fg = t.Foreground
					if !known {
						segment.Background(face)
						segment = style.Shadow(segment, c, shadows)
					}
				}
				segment.Children(func() {
					style.Text(ui.Text(c, item.Label), c, theme.SM()).
						FontWeight(theme.WeightMedium).TextColor(fg).SingleLine()
				})
				style.Disabled(segment, item.Disabled)
				box := segment.LayoutBounds()
				if box.W > 0 && box.H > 0 {
					box.X -= trackBox.X
					box.Y -= trackBox.Y
				}
				(*boxes)[i] = box
				segment.DrawOver(func(_ *ui.Painter, measured ui.Rect) {
					if paintedTrack.W <= 0 {
						return
					}
					measured.X -= paintedTrack.X
					measured.Y -= paintedTrack.Y
					if layoutChanged(box, measured) {
						c.Invalidate()
					}
				})
			}
		})
		var rect ui.Rect
		if known {
			target := (*boxes)[placed]
			if target.W <= 0 || target.H <= 0 {
				known = false
			} else {
				if indicator.sel < 0 {
					from := target
					if indicator.seen >= 0 && indicator.seen < len(*boxes) && (*boxes)[indicator.seen].W > 0 {
						from = (*boxes)[indicator.seen]
					}
					track.Jump(tabsX, from.X)
					track.Jump(tabsY, from.Y)
					track.Jump(tabsW, from.W)
					track.Jump(tabsH, from.H)
				}
				animate := !theme.ReduceMotion(c) && (indicator.sel != placed || indicator.box == target)
				follow := func(key any, value float32) float32 {
					if animate {
						return track.AnimateTo(key, value, ui.Animation{Duration: 200 * time.Millisecond, Ease: ui.EaseInOut})
					}
					return track.Jump(key, value)
				}
				rect = ui.Rect{
					X: follow(tabsX, target.X), Y: follow(tabsY, target.Y),
					W: follow(tabsW, target.W), H: follow(tabsH, target.H),
				}
				*indicator = tabsIndicator{seen: *selected, sel: placed, box: target}
			}
		}
		if !known {
			indicator.seen = *selected
		}
		track.Draw(func(p *ui.Painter, trackRect ui.Rect) {
			paintedTrack = trackRect
			if !known {
				return
			}
			shown := rect
			shown.X += trackRect.X
			shown.Y += trackRect.Y
			p.Fill(shown, face, radius)
			for _, sh := range shadows {
				p.Shadow(shown, radius, sh.X, sh.Y, sh.Blur, sh.Spread, sh.Color)
			}
		})
		if track.Changed() && o.OnChange != nil {
			o.OnChange(*selected)
		}
	})
	return root
}
