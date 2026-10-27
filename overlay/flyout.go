package overlay

import (
	"fmt"
	"math"
	"runtime"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// FlyoutOptions describes a native, independently themed flyout window.
// Exactly one of Parent and Tray is required. Width and Height are window
// DIPs and must be positive; native content cannot use MyGo's FitContent.
// Body, Header and Footer must not capture a context or element from another
// build or window. The caller owns the returned window and its state.
type FlyoutOptions struct {
	Parent        *mygo.Window
	Tray          *mygo.Tray
	Anchor        mygo.Rectangle
	Placement     mygo.Placement
	Gap           int
	Width, Height int
	Focusable     bool
	Hidden        bool
	// Dark overrides the system appearance. The caller invalidates the flyout
	// after changing the pointed-to value, as with other window-owned state.
	Dark    *bool
	Surface FlyoutSurfaceOptions
}

// FlyoutSurfaceOptions describes the coss popover content of a flyout.
// Header replaces the default title and description when provided.
type FlyoutSurfaceOptions struct {
	Title, Description, Label string
	Header, Body, Footer      func(*ui.Context)
}

// NewFlyout creates a native flyout at action time, after the app is ready.
// MyGo owns window placement and dismissal; it may panic on native window
// creation failures or if called before the app is ready. WebAssembly does
// not implement native flyout placement and returns an error instead.
func NewFlyout(o FlyoutOptions) (*mygo.Window, error) {
	if runtime.GOOS == "js" {
		return nil, fmt.Errorf("overlay: native flyouts are unavailable on WebAssembly")
	}
	options, err := nativeFlyoutOptions(o)
	if err != nil {
		return nil, err
	}
	return mygo.NewFlyout(options), nil
}

func nativeFlyoutOptions(o FlyoutOptions) (mygo.FlyoutOptions, error) {
	if (o.Parent == nil) == (o.Tray == nil) {
		return mygo.FlyoutOptions{}, fmt.Errorf("overlay: flyout requires exactly one parent or tray")
	}
	if o.Width <= 0 || o.Height <= 0 {
		return mygo.FlyoutOptions{}, fmt.Errorf("overlay: flyout width and height must be positive")
	}
	switch o.Placement {
	case "", mygo.PlacementBottomStart, mygo.PlacementBottom, mygo.PlacementBottomEnd,
		mygo.PlacementTopStart, mygo.PlacementTop, mygo.PlacementTopEnd,
		mygo.PlacementRightStart, mygo.PlacementRight, mygo.PlacementRightEnd,
		mygo.PlacementLeftStart, mygo.PlacementLeft, mygo.PlacementLeftEnd:
	default:
		return mygo.FlyoutOptions{}, fmt.Errorf("overlay: unknown flyout placement %q", o.Placement)
	}
	surface := o.Surface
	return mygo.FlyoutOptions{
		Parent: o.Parent, Tray: o.Tray, Anchor: o.Anchor, Placement: o.Placement,
		Gap: o.Gap, Width: o.Width, Height: o.Height, Focusable: o.Focusable,
		Hidden: o.Hidden, Shadow: true, CornerRadius: int(math.Ceil(float64(theme.RadiusScale().LG))),
		Content: ui.View(func(c *ui.Context) {
			if o.Dark == nil {
				theme.Apply(c)
			} else {
				theme.ApplyMode(c, *o.Dark)
			}
			FlyoutSurface(c, surface)
		}),
	}, nil
}

// FlyoutAnchor converts committed handle bounds from content DIPs to the
// parent window's DIPs, including its current UI zoom. Call during a parent
// build or input action; update the window with SetAnchor if the target moves.
func FlyoutAnchor(c *ui.Context, anchor ui.Handle) mygo.Rectangle {
	r := anchor.Bounds(c)
	z := float64(c.Zoom())
	x := math.Floor(float64(r.X) * z)
	y := math.Floor(float64(r.Y) * z)
	right := math.Ceil(float64(r.X+r.W) * z)
	bottom := math.Ceil(float64(r.Y+r.H) * z)
	return mygo.Rectangle{X: int(x), Y: int(y), Width: int(right - x), Height: int(bottom - y)}
}

// FlyoutSurface builds a scrollable coss popover surface that fills its
// available bounds. Apply the theme at the window root before calling it.
// Its returned element is valid only for the current build.
func FlyoutSurface(c *ui.Context, o FlyoutSurfaceOptions) ui.Element {
	t := theme.Colors(c)
	label := o.Label
	if label == "" {
		label = o.Title
	}
	panel := ui.Box(c).Fill().MinWidth(0).MinHeight(0).
		Radius(theme.RadiusScale().LG).Background(t.Popover).TextColor(t.PopoverForeground).
		Border(theme.Hairline, t.Border)
	if label != "" {
		panel.Label(label)
	}
	if o.Description != "" {
		panel.Description(o.Description)
	}
	panel.Children(func() {
		ui.Scroll(c).Fill().MinWidth(0).MinHeight(0).Padding(theme.Space16).Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
				if o.Header != nil || o.Title != "" || o.Description != "" {
					ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4).Children(func() {
						if o.Header != nil {
							o.Header(c)
							return
						}
						if o.Title != "" {
							style.Text(ui.Text(c, o.Title), c, theme.LG()).FontWeight(theme.WeightSemibold).
								Role(ui.RoleHeading).Level(3)
						}
						if o.Description != "" {
							style.Text(ui.Text(c, o.Description), c, theme.SM()).TextColor(t.MutedForeground)
						}
					})
				}
				if o.Body != nil {
					ui.Column(c).FillWidth().MinWidth(0).Children(func() { o.Body(c) })
				}
				if o.Footer != nil {
					ui.Row(c).FillWidth().Justify(ui.End).Gap(theme.Space8).Children(func() { o.Footer(c) })
				}
			})
		})
	})
	return panel
}
