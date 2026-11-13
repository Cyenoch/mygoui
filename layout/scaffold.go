package layout

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ScaffoldOptions configures a window's shell. TitleBar is the bar over the
// content. Start is a full-height pane on the start edge, such as a
// sidebar; StartBar is its share of the title band, where an application
// name goes. End is a pane beside the content under the title bar, such as
// an inspector. StatusBar builds a strip under everything, such as a
// navigation.StatusBar. ContentMin and Compact are a SplitView's.
//
// The scaffold sets the Span of both bars and gives them one height;
// StartBar is borderless.
type ScaffoldOptions struct {
	TitleBar   TitleBarOptions
	Start      *PaneOptions
	StartBar   TitleBarOptions
	End        *PaneOptions
	Content    func()
	StatusBar  func()
	ContentMin float32
	Compact    CompactMode
}

// Scaffold builds the standard window shell on the window's background: a
// full-height start pane on the sidebar surface whose top band shares the
// title band, the title bar over the content, an end pane beside the
// content and a strip across the bottom. The start band reserves the start
// window controls; the title bar reserves the end ones, and the start ones
// too whenever the start pane is not docked beside it. Panes behave as a
// SplitView's, the start pane's sheet holding its band as well.
func Scaffold(c *ui.Context, o ScaffoldOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).Fill().MinWidth(0).MinHeight(0).Background(t.Background).TextColor(t.Foreground)
	m := ui.Local(root, barKey{}, func() barMeasure { return barMeasure{} })
	compact := Compact(c, o.Compact)
	bar := o.TitleBar
	// The span decides the bar's height: where the start pane will be, not
	// where its motion has it now.
	bar.Span = SpanWindow
	if o.Start != nil && !compact && (o.Start.Collapsed == nil || !*o.Start.Collapsed) {
		bar.Span = SpanEnd
	}
	g := resolveBar(c, bar, m)
	// The scaffold keeps a pane's size itself when the caller does not, to
	// know the room the pane takes.
	own := func(o *PaneOptions, end bool) *PaneOptions {
		if o == nil || o.Size != nil {
			return o
		}
		p := *o
		p.Size = ui.Local(root, paneSizeKey{end}, func() float32 { return 0 })
		return &p
	}
	start, end := own(o.Start, false), own(o.End, true)
	if start != nil {
		p, band := *start, o.StartBar
		band.Span, band.Borderless = SpanStart, true
		width := theme.Scale(c, resolvePane(c, root, start, false).width())
		p.Body = func() {
			ui.Column(c).Fill().MinWidth(0).MinHeight(0).
				Background(t.Sidebar).TextColor(t.SidebarForeground).Children(func() {
				left, right := barInsets(c, SpanStart)
				if compact {
					// A sheet as wide as the window reaches the end controls.
					l, r := barInsets(c, SpanWindow)
					window, _ := c.Size()
					beyond := max(0, window-width)
					left, right = max(left, l-beyond), max(right, r-beyond)
				}
				// A window has one banner: the title bar.
				fillBar(c, ui.Row(c), band, barGeometry{height: g.height}, nil, left, right).Role(ui.RoleGroup)
				ui.Column(c).FillWidth().Grow(1).Basis(0).MinWidth(0).MinHeight(0).Children(o.Start.Body)
			})
		}
		start = &p
	}
	contentMin := o.ContentMin
	if contentMin <= 0 {
		contentMin = contentMinimum
	}
	beside := contentMin
	if !compact {
		beside += resolvePane(c, root, end, true).docked()
	}
	sidebar := &startDock{}
	root.Children(func() {
		splitView(c, SplitViewOptions{
			Start: start, ContentMin: beside, Compact: o.Compact,
			Content: func() {
				// The bar clears what the start pane leaves of the start
				// window controls, in every frame of the pane's motion.
				left, right := barInsets(c, SpanWindow)
				if c.Direction() == ui.RightToLeft {
					right = max(0, right-sidebar.shown)
				} else {
					left = max(0, left-sidebar.shown)
				}
				fillBar(c, ui.Row(c), bar, g, m, left, right)
				splitView(c, SplitViewOptions{
					End: end, Content: o.Content, ContentMin: contentMin, Compact: o.Compact,
				}, nil).FillWidth().Grow(1).Basis(0)
			},
		}, sidebar).FillWidth().Grow(1).Basis(0)
		if o.StatusBar != nil {
			o.StatusBar()
		}
	})
	return root
}
