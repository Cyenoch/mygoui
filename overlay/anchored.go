package overlay

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
	"time"
)

type Align uint8

const (
	AlignCenter Align = iota
	AlignStart
	AlignEnd
)

type Placement struct {
	Side   Side
	Align  Align
	Offset float32
}
type PopoverOptions struct {
	Title, Description string
	Width              float32
	Placement          Placement
	Header, Footer     func()
}

func Popover(c *myui.Context, anchor myui.Element, open *bool, o PopoverOptions, body func()) myui.Element {
	if open == nil || !anchor.Valid() {
		return myui.Element{}
	}
	return myui.PopoverBase(c, anchor, open, func(panel myui.Element) {
		stylePopover(c, panel, anchor, o)
		if o.Header != nil {
			o.Header()
		} else {
			popoverHeading(c, o.Title, o.Description)
		}
		if body != nil {
			body()
		}
		if o.Footer != nil {
			o.Footer()
		}
	})
}

func stylePopover(c *myui.Context, panel, anchor myui.Element, o PopoverOptions) {
	p := theme.Colors(c)
	w, h := c.Size()
	width := o.Width
	if width <= 0 {
		width = 288
	}
	panel.Column().Width(width).MaxWidth(max(0, w-16)).MaxHeight(max(0, h-16)).Padding(16).Gap(16).
		Radius(24).Background(p.Popover).TextColor(p.PopoverForeground).Border(1, p.Border).
		Font(theme.FontFamily(c)).FontSize(theme.Scale(c, 13)).FixedLineHeight(theme.Scale(c, 18)).
		Label(o.Title).Description(o.Description)
	place(panel, anchor, o.Placement, SideBottom)
}

// Anchored previews use compact interface roles, not a dialog's title scale.
func popoverHeading(c *myui.Context, title, description string) {
	if title == "" && description == "" {
		return
	}
	myui.Column(c).Gap(4).Children(func() {
		if title != "" {
			myui.Text(c, title).FontWeight(theme.Title().Weight).
				FontSize(theme.Scale(c, theme.Title().Size)).FixedLineHeight(theme.Scale(c, 20)).
				Role(myui.RoleHeading).Level(2)
		}
		if description != "" {
			myui.Text(c, description).TextColor(theme.Colors(c).MutedForeground)
		}
	})
}

func place(panel, anchor myui.Element, p Placement, fallback Side) {
	side := p.Side
	if side == SideDefault {
		side = fallback
	}
	offset := p.Offset
	if offset == 0 {
		offset = 4
	}
	at, self := myui.AnchorBottom, myui.AnchorTop
	switch side {
	case SideTop:
		at, self = myui.AnchorTop, myui.AnchorBottom
		switch p.Align {
		case AlignStart:
			at, self = myui.AnchorTopStart, myui.AnchorBottomStart
		case AlignEnd:
			at, self = myui.AnchorTopEnd, myui.AnchorBottomEnd
		}
		panel.Margin(offset, 0, 0, 0)
	case SideRight:
		at, self = myui.AnchorRight, myui.AnchorLeft
		switch p.Align {
		case AlignStart:
			at, self = myui.AnchorTopRight, myui.AnchorTopLeft
		case AlignEnd:
			at, self = myui.AnchorBottomRight, myui.AnchorBottomLeft
		}
		panel.Margin(0, 0, 0, offset)
	case SideLeft:
		at, self = myui.AnchorLeft, myui.AnchorRight
		switch p.Align {
		case AlignStart:
			at, self = myui.AnchorTopLeft, myui.AnchorTopRight
		case AlignEnd:
			at, self = myui.AnchorBottomLeft, myui.AnchorBottomRight
		}
		panel.Margin(0, 0, 0, offset)
	default:
		switch p.Align {
		case AlignStart:
			at, self = myui.AnchorBottomStart, myui.AnchorTopStart
		case AlignEnd:
			at, self = myui.AnchorBottomEnd, myui.AnchorTopEnd
		}
		panel.Margin(offset, 0, 0, 0)
	}
	panel.AttachTo(anchor, at, self)
}

type TooltipOptions struct {
	Text      string
	Placement Placement
}

// Tooltip uses the maintained native timing: 600ms hover, immediate keyboard
// focus, 400ms adjacent-tip grace, and suppression after Escape or activation.
// Essential information must also be available in the surrounding interface.
func Tooltip(c *myui.Context, anchor myui.Element, o TooltipOptions) myui.Element {
	if !anchor.Valid() || o.Text == "" {
		return myui.Element{}
	}
	anchor.Description(o.Text)
	return myui.TooltipBase(c, anchor, func(tip myui.Element) {
		p := theme.Colors(c)
		tip.Row().MaxWidth(320).Padding(6, 12).Radius(12).Gap(6).Background(p.Foreground).TextColor(p.Background).FontSize(theme.Scale(c, 12))
		place(tip, anchor, o.Placement, SideTop)
		tip.Children(func() { myui.Text(c, o.Text) })
	})
}

type HoverCardOptions struct {
	PopoverOptions
	OpenDelay, CloseDelay time.Duration
}
type HoverCardState struct {
	Open         bool
	enter, leave time.Time
	suppressed   bool
	panel        myui.Handle
}

// HoverCard keeps focusable preview content open while either surface is hovered
// or focused. A close grace period lets the pointer cross the anchoring gap.
func HoverCard(c *myui.Context, anchor myui.Element, s *HoverCardState, o HoverCardOptions, body func()) myui.Element {
	if s == nil || !anchor.Valid() {
		return myui.Element{}
	}
	now := c.Now()
	over := anchor.Hovered() || anchor.FocusWithin() || s.panel.FocusWithin(c)
	if !over && !s.Open {
		s.suppressed = false
	}
	if over && !s.suppressed {
		s.leave = time.Time{}
		if s.enter.IsZero() {
			s.enter = now
		}
		delay := o.OpenDelay
		if delay == 0 {
			delay = 600 * time.Millisecond
		}
		if anchor.FocusWithin() || now.Sub(s.enter) >= delay {
			s.Open = true
		} else {
			c.After(delay - now.Sub(s.enter))
		}
	} else {
		s.enter = time.Time{}
	}
	if !s.Open {
		return myui.Element{}
	}
	var panel myui.Element
	myui.Overlay(c, func() {
		panel = myui.Column(c.Key("hover-card")).Bind(&s.panel).Role(myui.RolePopup)
		stylePopover(c, panel, anchor, o.PopoverOptions)
		panel.Children(func() {
			if o.Header != nil {
				o.Header()
			} else {
				popoverHeading(c, o.Title, o.Description)
			}
			if body != nil {
				body()
			}
			if o.Footer != nil {
				o.Footer()
			}
		})
		if panel.Hovered() || panel.FocusWithin() || over {
			s.leave = time.Time{}
		} else {
			if s.leave.IsZero() {
				s.leave = now
			}
			delay := o.CloseDelay
			if delay == 0 {
				delay = 300 * time.Millisecond
			}
			if now.Sub(s.leave) >= delay {
				s.Open = false
			} else {
				c.After(delay - now.Sub(s.leave))
			}
		}
		if panel.OverlayShortcut(0, myui.KeyEscape) || panel.PressedOutside() {
			s.Open = false
			s.suppressed = true
		}
	})
	return panel
}
