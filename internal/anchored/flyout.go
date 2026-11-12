package anchored

import (
	"math"
	"runtime"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// FlyoutSpec describes a popup presented in a native flyout.
type FlyoutSpec struct {
	Placement mygoui.Placement
	// Fallback is the side used when Placement leaves it unset.
	Fallback mygoui.Side
	// MinWidth, MaxWidth and MaxHeight bound the fitted flyout in the
	// parent's content units.
	MinWidth, MaxWidth, MaxHeight float32
	Label                         string
	// Body builds the popup content in the flyout's own context. It must not
	// capture elements, handles bound or parent-keyed state of the parent
	// build; caller-owned values are shared.
	Body func(*ui.Context)
	// OnDismiss runs after the user closes the flyout and *open turns false.
	OnDismiss func()
}

// flyoutKey keys an anchor's flyout owner for one open flag.
type flyoutKey struct{ open *bool }

// flyouts holds the flyout each caller open flag presents, from the build
// that opens it until it closes. An anchor's position in the tree can change
// between builds; the caller's flag cannot. Main thread only.
var flyouts = map[*bool]*flyoutState{}

// flyoutOwner is the anchor that presented a flyout in its latest build.
// MyGo closes it after a frame that no longer builds that anchor; the flyout
// closes with it unless another anchor presented it in that frame.
type flyoutOwner struct{ open *bool }

func (o *flyoutOwner) Close() error {
	if s := flyouts[o.open]; s != nil && s.owner == o {
		s.close()
	}
	return nil
}

// flyoutState is one presented flyout.
type flyoutState struct {
	open    *bool
	owner   *flyoutOwner
	parent  *mygo.Window
	win     *mygo.Window
	anchor  mygo.Rectangle
	spec    FlyoutSpec
	zoom    float32
	dark    bool
	echoing bool // the flyout's next build follows a parent build
}

// close forgets s and closes its window after the current build, leaving
// the caller's flag as is.
func (s *flyoutState) close() {
	delete(flyouts, s.open)
	if win := s.win; win != nil {
		s.win = nil
		s.parent.Update(win.Close)
	}
}

// Flyout presents spec in a native flyout of host's window while *open is
// true and reports whether it did; false asks the caller to build the popup
// in-window, closing a flyout this flag presented before. The flyout opens
// after the build that sees *open, closes when *open turns false or the
// anchor leaves the tree, and clears *open when the user dismisses it.
func Flyout(c *ui.Context, host *mygoui.FlyoutHost, anchor ui.Element, open *bool, spec FlyoutSpec) bool {
	if host == nil || host.Window == nil || runtime.GOOS == "js" || !anchor.Valid() {
		if s := flyouts[open]; s != nil {
			s.close()
		}
		return false
	}
	owner := *ui.Local(anchor, flyoutKey{open}, func() *flyoutOwner { return &flyoutOwner{open: open} })
	s := flyouts[open]
	if s == nil {
		if !*open {
			return true
		}
		s = &flyoutState{open: open}
		flyouts[open] = s
	}
	s.owner, s.parent, s.spec, s.zoom, s.dark = owner, host.Window, spec, c.Zoom(), c.Theme().Dark
	rect := flyoutRect(c, anchor.Bounds())
	switch {
	case !*open:
		s.close()
	case s.win == nil:
		s.anchor = rect
		host.Window.Update(s.show)
	default:
		if rect != s.anchor {
			s.anchor = rect
			s.win.SetAnchor(rect)
		}
		// Parent state the body shows may have changed.
		s.echoing = true
		s.win.Invalidate()
	}
	return true
}

// show creates the flyout on the main thread after the parent build.
func (s *flyoutState) show() {
	if !*s.open || s.win != nil || flyouts[s.open] != s {
		return
	}
	side := Side(s.spec.Placement, s.spec.Fallback)
	gap := s.spec.Placement.Offset
	if gap == 0 {
		gap = defaultOffset
	}
	win := mygo.NewFlyout(mygo.FlyoutOptions{
		Parent: s.parent, Anchor: s.anchor, Placement: flyoutPlacement(side, s.spec.Placement.Align),
		Gap: s.dip(gap), Width: s.dip(s.spec.MinWidth), Height: 1,
		FitContent: true, MinWidth: s.dip(s.spec.MinWidth),
		MaxWidth: s.dip(s.spec.MaxWidth), MaxHeight: s.dip(s.spec.MaxHeight),
		Focusable: true, Hidden: true, Shadow: true,
		CornerRadius: s.dip(theme.RadiusScale().LG),
		Content:      ui.View(s.build),
	})
	s.win = win
	win.OnClosed(func() {
		if s.win != win {
			return // closed by the parent
		}
		s.win = nil
		delete(flyouts, s.open)
		*s.open = false
		if s.spec.OnDismiss != nil {
			s.spec.OnDismiss()
		}
		s.parent.Invalidate()
	})
	win.Show()
}

// build is the flyout's view. Input in the flyout changes caller state the
// parent shows, so a build of its own redraws the parent; a build that
// follows the parent's does not, which ends the exchange.
func (s *flyoutState) build(c *ui.Context) {
	theme.ApplyMode(c, s.dark)
	if s.echoing {
		s.echoing = false
	} else {
		s.parent.Invalidate()
	}
	t := theme.Colors(c)
	panel := ui.Column(c).FillWidth().MinWidth(0).MinHeight(0).Clip().
		Radius(theme.RadiusScale().LG).Background(t.Popover).TextColor(t.PopoverForeground)
	if runtime.GOOS == "linux" {
		// macOS and Windows outline a flyout with its system shadow.
		panel.Border(theme.Hairline, t.Border)
	}
	if s.spec.Label != "" {
		panel.Label(s.spec.Label)
	}
	panel.Children(func() { s.spec.Body(c) })
}

// flyoutRect converts committed content bounds to parent window DIPs.
func flyoutRect(c *ui.Context, r ui.Rect) mygo.Rectangle {
	z := float64(c.Zoom())
	x := math.Floor(float64(r.X) * z)
	y := math.Floor(float64(r.Y) * z)
	right := math.Ceil(float64(r.X+r.W) * z)
	bottom := math.Ceil(float64(r.Y+r.H) * z)
	return mygo.Rectangle{X: int(x), Y: int(y), Width: int(right - x), Height: int(bottom - y)}
}

// dip converts parent content units to DIPs.
func (s *flyoutState) dip(v float32) int { return int(math.Ceil(float64(v * s.zoom))) }

func flyoutPlacement(side mygoui.Side, align mygoui.Align) mygo.Placement {
	suffix := map[mygoui.Align]string{mygoui.AlignStart: "-start", mygoui.AlignEnd: "-end"}[align]
	switch side {
	case mygoui.SideTop:
		return mygo.Placement("top" + suffix)
	case mygoui.SideLeft:
		return mygo.Placement("left" + suffix)
	case mygoui.SideRight:
		return mygo.Placement("right" + suffix)
	}
	return mygo.Placement("bottom" + suffix)
}
