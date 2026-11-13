package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ImageStatus is the caller's load state for an Image.
type ImageStatus uint8

const (
	ImageReady ImageStatus = iota
	ImageLoading
	ImageFailed
)

// imageFace keys what the box shows, so that one face fades into the next.
type imageFace uint8

const (
	imageShown imageFace = iota
	imageLoading
	imageFailed
)

// ImageFit selects how the picture fills its box.
type ImageFit uint8

const (
	ImageCover ImageFit = iota
	ImageContain
)

// ImageOptions configures a caller-loaded image. Label is the accessible
// name. A zero Width fills the parent. Height wins over AspectRatio
// (width / height); with neither, the box takes the picture's own ratio, or
// 16:9 while there is none. A nonpositive Radius uses the theme's large radius.
type ImageOptions struct {
	Source        ui.ImageSource
	Label         string
	Width, Height float32
	AspectRatio   float32
	Fit           ImageFit
	Radius        float32
	Status        ImageStatus
	// Fallback adds content under the failed icon.
	Fallback func()
	// OnClick makes a shown picture a button, as for opening a preview.
	OnClick func()
}

// Image shows an already-decoded picture in a rounded box, a skeleton while
// the caller loads it and a fallback when loading failed or there is no
// picture. It never fetches or decodes.
func Image(c *ui.Context, o ImageOptions) ui.Element {
	t := theme.Colors(c)
	radius := o.Radius
	if radius <= 0 {
		radius = theme.RadiusScale().LG
	}
	shown := o.Status == ImageReady && hasImage(o.Source)
	var e ui.Element
	if shown && o.OnClick != nil {
		e = ui.ButtonBase(c).FocusRing(true).Cursor(ui.CursorPointer).OnClick(o.OnClick)
	} else {
		e = ui.Box(c).Role(ui.RoleImage)
	}
	e.Label(o.Label).Shrink(0).Radius(radius).Clip().Unselectable()
	if o.Width > 0 {
		e.Width(o.Width)
	} else {
		e.FillWidth().MinWidth(0)
	}
	if o.Height > 0 {
		e.Height(o.Height)
	} else {
		e.AspectRatio(imageRatio(o))
	}
	return e.Children(func() {
		switch {
		case shown:
			fit := ui.Cover
			if o.Fit == ImageContain {
				fit = ui.Contain
			}
			ui.Image(c.Key(imageShown), o.Source).Fill().Fit(fit).Role(ui.RoleNone).
				Transition(fadeTransition)
		case o.Status == ImageLoading:
			Skeleton(c.Key(imageLoading), SkeletonOptions{Radius: radius}).Fill().Transition(fadeTransition)
		default:
			ui.Column(c.Key(imageFailed)).Fill().Center().Gap(theme.Space8).Padding(theme.Space8).
				Background(t.MutedBackground).TextColor(t.MutedForeground).Transition(fadeTransition).
				Children(func() {
					size := theme.Scale(c, theme.IconLG)
					ui.Icon(c, lucide.ImageOff()).Size(size, size).Shrink(0)
					if o.Fallback != nil {
						o.Fallback()
					}
				})
		}
	})
}

// imageRatio resolves the box ratio when the caller sets no height.
func imageRatio(o ImageOptions) float32 {
	if o.AspectRatio > 0 {
		return o.AspectRatio
	}
	if w, h := ui.ImageSize(o.Source); w > 0 && h > 0 {
		return w / h
	}
	return 16.0 / 9
}
