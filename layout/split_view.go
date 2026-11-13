package layout

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// CompactWidth is the window width, in DIPs before the text scale, below
// which CompactAuto is compact.
const CompactWidth = 760

// CompactMode selects when a view shows its side panes over the content
// instead of beside it.
type CompactMode uint8

const (
	// CompactAuto is compact while the window is narrower than CompactWidth
	// scaled with the text size.
	CompactAuto CompactMode = iota
	// CompactNever always docks the panes.
	CompactNever
	// CompactAlways always shows the panes over the content.
	CompactAlways
)

// Compact reports whether mode is compact in c's window. Content that
// adapts with a SplitView or Scaffold reads the same decision here.
func Compact(c *ui.Context, mode CompactMode) bool {
	switch mode {
	case CompactNever:
		return false
	case CompactAlways:
		return true
	}
	width, _ := c.Size()
	return width < theme.Scale(c, CompactWidth)
}

// PaneOptions configures a side pane. Label names the pane and its divider
// for assistive technology. Body builds its content in a column that fills
// the pane.
//
// Width, Min and Max are DIPs before the text scale; nonpositive values use
// 256 and 160, and no maximum. Size is the caller-owned width: zero shows
// Width. A nil Size is kept by the view under its own identity: key the view
// (Context.Key) or pass Size when elements before it come and go.
// Fixed keeps the pane at Width behind a hairline, with no divider.
//
// Collapsed hides the docked pane while set; nil makes it permanent. Open
// shows the pane over the content while the view is compact; with a nil
// Open a compact view has no way to show it. The view clears Open whenever
// it is not compact.
type PaneOptions struct {
	Label     string
	Body      func()
	Size      *float32
	Width     float32
	Min, Max  float32
	Fixed     bool
	Collapsed *bool
	Open      *bool
}

// SplitViewOptions configures a split view. Start and End are the side
// panes in layout direction, either of which may be nil; Content is what
// they sit beside. ContentMin is the width, in DIPs before the text scale,
// the content keeps while a divider moves; nonpositive uses 240.
type SplitViewOptions struct {
	Start      *PaneOptions
	End        *PaneOptions
	Content    func()
	ContentMin float32
	Compact    CompactMode
}

const (
	paneWidth      = 256
	paneMin        = 160
	contentMinimum = 240
)

type (
	paneKey      struct{ end bool }
	paneSizeKey  struct{ end bool }
	paneFoldKey  struct{}
	splitBodyKey struct{}
)

// pane is one side pane resolved for a build. shown is the width it is
// docked at in this build, after the text scale; border is its hairline.
type pane struct {
	o        *PaneOptions
	end      bool
	size     *float32
	min, max float32
	shown    float32
	border   ui.Color
}

// startDock is what a scaffold needs of its start pane: the sidebar's own
// hairline, and the width the pane is docked at in this build.
type startDock struct{ shown float32 }

func (p pane) collapsed() bool { return p.o.Collapsed != nil && *p.o.Collapsed }

// width is the size the pane shows at, before the text scale.
func (p pane) width() float32 {
	w := *p.size
	if w <= 0 || p.o.Fixed {
		w = p.o.Width
		if w <= 0 {
			w = paneWidth
		}
	}
	if p.o.Fixed {
		return w
	}
	w = max(w, p.min)
	if p.max > 0 {
		w = min(w, max(p.max, p.min))
	}
	return w
}

// docked is the room the pane takes in the row, before the text scale.
func (p *pane) docked() float32 {
	if p == nil || p.collapsed() {
		return 0
	}
	return p.width()
}

// SplitView arranges a start pane, the content and an end pane in a row
// that fills its parent. A pane's divider is a splitter the user drags or
// moves with the arrow keys within Min and Max, leaving the content its
// minimum while the view has the room; Home and End go to the bounds, and
// Enter, a double click or a drag below half of Min collapses a pane that
// has a Collapsed flag. The divider stays a Tab stop while its pane is
// collapsed.
//
// A compact view shows the content alone and each pane with an Open flag as
// a modal sheet from the pane's edge while the flag is set.
func SplitView(c *ui.Context, o SplitViewOptions) ui.Element {
	return splitView(c, o, nil).Fill()
}

func splitView(c *ui.Context, o SplitViewOptions, sidebar *startDock) ui.Element {
	root := ui.Row(c).MinWidth(0).MinHeight(0).AlignItems(ui.Stretch)
	compact := Compact(c, o.Compact)
	start, end := resolvePane(c, root, o.Start, false), resolvePane(c, root, o.End, true)
	if !compact {
		fitPanes(c, root, o, start, end)
	}
	if sidebar != nil && start != nil {
		start.border = theme.Colors(c).SidebarBorder
	}
	root.Children(func() {
		if !compact {
			start.dock(c)
			if sidebar != nil && start != nil {
				sidebar.shown = start.shown
			}
		}
		ui.Column(c.Key(splitBodyKey{})).Grow(1).Basis(0).MinWidth(0).FillHeight().Children(o.Content)
		if !compact {
			end.dock(c)
		}
	})
	start.overlay(c, compact)
	end.overlay(c, compact)
	return root
}

func resolvePane(c *ui.Context, root ui.Element, o *PaneOptions, end bool) *pane {
	if o == nil {
		return nil
	}
	p := &pane{o: o, end: end, size: o.Size, min: o.Min, max: o.Max, border: theme.Colors(c).Border}
	if p.size == nil {
		p.size = ui.Local(root, paneSizeKey{end}, func() float32 { return 0 })
	}
	if p.min <= 0 {
		p.min = paneMin
	}
	return p
}

// fitPanes lowers each resizable pane's maximum so the content keeps its
// minimum in the width the view was last laid out at, the window's before
// its first layout. The start pane is served first; neither goes below its
// own minimum.
func fitPanes(c *ui.Context, root ui.Element, o SplitViewOptions, start, end *pane) {
	width := root.LayoutBounds().W
	if width <= 0 {
		width, _ = c.Size()
	}
	width /= theme.TextScale(c)
	if width <= 0 {
		return
	}
	room := width - o.ContentMin
	if o.ContentMin <= 0 {
		room = width - contentMinimum
	}
	limit := func(p *pane, room float32) {
		if p == nil || p.o.Fixed {
			return
		}
		room = max(room, p.min)
		if p.max <= 0 || p.max > room {
			p.max = room
		}
	}
	other := float32(0)
	if end != nil && !end.collapsed() {
		other = end.min
		if end.o.Fixed {
			other = end.width()
		}
	}
	limit(start, room-other)
	limit(end, room-start.docked())
}

// body builds the pane's content, named as a group.
func (p *pane) body(c *ui.Context) {
	e := ui.Column(c).Fill().MinWidth(0).MinHeight(0).Role(ui.RoleGroup)
	if p.o.Label != "" {
		e.Label(p.o.Label)
	}
	e.Children(p.o.Body)
}

// dock builds the pane in the row.
func (p *pane) dock(c *ui.Context) {
	if p == nil {
		return
	}
	if p.o.Fixed {
		p.dockFixed(c)
		return
	}
	side := mygoui.PanelEnd
	if p.end {
		side = mygoui.PanelStart
	}
	// The panel resizes a copy, so a zero caller size keeps meaning the
	// default width until the user moves the divider.
	size := p.width()
	mygoui.ResizablePanel(c.Key(paneKey{p.end}), mygoui.ResizablePanelOptions{
		Label: p.o.Label, Side: side, Size: &size, Min: p.min, Max: p.max,
		Collapsed: p.o.Collapsed, Shown: &p.shown, Border: p.border,
		OnResize: func(v float32) { *p.size = v },
	}, func() { p.body(c) })
}

// dockFixed builds a pane of constant width behind a hairline. Collapsing
// moves between the width and nothing; a text scale change resizes at once.
func (p *pane) dockFixed(c *ui.Context) {
	root := ui.Box(c.Key(paneKey{p.end})).Shrink(0).AlignSelf(ui.Stretch).BorderColor(p.border)
	target := float32(0)
	if p.collapsed() {
		target = 1
	}
	m := theme.Motion()
	fold := root.AnimateTo(paneFoldKey{}, target, ui.Animation{Duration: m.Reveal, Ease: m.Ease})
	full := theme.Scale(c, p.width())
	p.shown = full * (1 - fold)
	root.Width(p.shown)
	if fold == 1 {
		return
	}
	// The hairline never makes the pane wider than it shows.
	if hairline := min(theme.Hairline, p.shown); p.end {
		root.BorderStartWidth(hairline)
	} else {
		root.BorderEndWidth(hairline)
	}
	if fold > 0 {
		root.Clip()
	}
	root.Children(func() {
		// The body keeps its width while the pane folds, so it does not
		// reflow in every frame of the motion.
		ui.Column(c).Width(full - theme.Hairline).FillHeight().Shrink(0).Children(func() { p.body(c) })
	})
}

// overlay shows the pane as a sheet from its edge while the view is compact
// and its Open flag is set, and clears the flag otherwise.
func (p *pane) overlay(c *ui.Context, compact bool) {
	if p == nil || p.o.Open == nil {
		return
	}
	if !compact {
		*p.o.Open = false
		return
	}
	// A sheet's side follows the layout direction, as the row does.
	side := mygoui.SideLeft
	if p.end {
		side = mygoui.SideRight
	}
	overlay.Sheet(c, p.o.Open, overlay.SheetOptions{
		Side: side, Width: theme.Scale(c, p.width()), Title: p.o.Label,
		Flush: true, Body: p.o.Body,
	})
}

// PaneToggleOptions configures the button that shows and hides a pane.
// Collapsed, Open and Compact are the pane's and its view's; End selects
// the end pane's icon and default name. Label names the button and defaults
// to the locale's ToggleSidebar or ToggleInspector.
type PaneToggleOptions struct {
	Collapsed *bool
	Open      *bool
	End       bool
	Compact   CompactMode
	Label     string
}

// PaneToggle builds a small icon button for a title bar that flips the
// pane's Collapsed flag while the view docks its panes and its Open flag
// while the view is compact. It is a toggle button, checked while the pane
// shows, and disabled when the flag it would flip is nil.
func PaneToggle(c *ui.Context, o PaneToggleOptions) ui.Element {
	flag, shown := o.Collapsed, o.Collapsed != nil && !*o.Collapsed
	if Compact(c, o.Compact) {
		flag, shown = o.Open, o.Open != nil && *o.Open
	}
	label, icon := o.Label, lucide.PanelLeft()
	if label == "" {
		label = locale.M().ToggleSidebar
		if o.End {
			label = locale.M().ToggleInspector
		}
	}
	if o.End != (c.Direction() == ui.RightToLeft) {
		icon = lucide.PanelRight()
	}
	return mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Icon: icon, IconOnly: true,
		Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall, Disabled: flag == nil,
		OnClick: func() { *flag = !*flag },
	}).Role(ui.RoleToggleButton).Checked(shown).Shrink(0)
}
