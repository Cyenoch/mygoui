package content

import (
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// LogLevel tints a log line. Info and plain lines keep the text color.
type LogLevel uint8

const (
	LogPlain LogLevel = iota
	LogDebug
	LogInfo
	LogWarning
	LogError
)

// LogViewerOptions configures a log. Lines is caller-owned; Line and Count
// supply the lines instead when the log is not held as a slice. Level tints
// line i. Follow keeps the newest line in view while the user is at the
// end, and shows a button back to it once they scroll away. Height is the
// viewer's height in DIPs; zero makes it as high as its lines, up to
// 320 DIPs. Toolbar adds a search field, the number of lines listed while a
// filter is on, and a toggle for wrapping lines above the lines, with a
// toggle per level when Level is set: the levels pressed are the only ones
// listed, all of them while none is, and a plain line counts as Info.
// Without a Toolbar the application sets the filter in State. OnError receives a clipboard failure of a copy.
type LogViewerOptions struct {
	Lines       []string
	Line        func(i int) string
	Count       int
	Level       func(i int) LogLevel
	LineNumbers bool
	Follow      bool
	Label       string
	Height      float32
	Toolbar     bool
	OnError     func(error)
	// State is caller-owned scroll state; the viewer keeps its own per
	// window when nil.
	State *LogViewerState
}

// LogViewerState holds a log's scroll position and its filter. List scrolls
// it from the application: ScrollToEnd, ScrollTo, Visible and AtEnd. The
// viewer names its rows (List.Label) and, unless List.Selected is set,
// keeps the line the keys choose.
//
// Query lists the lines holding it, whatever their case, and marks it in
// them; Levels lists the lines of its levels, every line while it is the
// zero set. Line numbers stay those of the log. The listing is kept until
// the query, the levels or Revision change, or the log shrinks, and a log
// that grew is scanned from its old end only: raise Revision when lines
// change in place or the log is replaced by one of the same length.
// NoWrap shows each line on one row, cut where it ends, rather than
// wrapped.
type LogViewerState struct {
	List     ui.ListState
	Query    string
	Levels   LogLevels
	Revision int
	NoWrap   bool

	filter logFilter
	search ui.Handle
	chosen int
}

// FocusSearch focuses the toolbar's search field, as Cmd+F (Ctrl+F) does
// inside the viewer, for a Find command of the window.
func (s *LogViewerState) FocusSearch() { s.search.Focus() }

// jumpFade fades the jump button in and out as the log leaves and reaches its
// end.
func jumpFade() ui.ElementTransition {
	m := theme.Motion()
	return ui.ElementTransition{
		Duration: m.State, Ease: m.Ease, Shadows: true,
		Enter: &ui.Motion{}, Exit: &ui.Motion{},
	}
}

// logDefaultMaxHeight bounds a log without a Height, so that it stays
// virtualized however many lines it holds.
const logDefaultMaxHeight = 320

type logViewerStateKey struct{}

func (o *LogViewerOptions) count() int {
	if o.Line != nil {
		return o.Count
	}
	return len(o.Lines)
}

func (o *LogViewerOptions) line(i int) string {
	if o.Line != nil {
		return o.Line(i)
	}
	return o.Lines[i]
}

// LogViewer shows monospace log lines, building only those in view. Each
// line's text is selectable; a selection does not span lines. Up and Down
// choose a line while the lines have the focus, and the context-menu key
// opens the chosen line's menu.
func LogViewer(c *ui.Context, o LogViewerOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).AlignSelf(ui.Stretch).MinWidth(0).Clip().Radius(theme.RadiusScale().XL).
		Background(t.Code).TextColor(t.CodeForeground).
		BorderWidth(theme.Hairline).BorderColor(t.Border)
	label := o.Label
	if label == "" {
		label = locale.M().Log
	}
	root.Role(ui.RoleGroup).Label(label)
	st := o.State
	if st == nil {
		st = ui.Local(root, logViewerStateKey{}, func() LogViewerState { return LogViewerState{} })
	}
	st.List.FollowEnd = o.Follow
	if st.List.Selected == nil {
		st.chosen = -1
		st.List.Selected = &st.chosen
	}
	// A row reads as its line, then its level where the line is tinted.
	st.List.Label = func(row int) string {
		i := st.original(row)
		if o.Level != nil {
			if _, name, ok := logLevelStyle(t, o.Level(i)); ok {
				return o.line(i) + ", " + name
			}
		}
		return o.line(i)
	}
	total := o.count()
	return root.Children(func() {
		if o.Toolbar {
			logToolbar(c, t, st, &o, total)
		}
		// The toolbar's controls write the filter as they are built.
		n := st.refresh(&o, total)
		gutter := logGutterWidth(c, total)
		pad := theme.Scale(c, theme.Space(2))
		if o.Toolbar {
			root.OnShortcut(ui.Cmd, ui.KeyF, func() { st.search.Focus(c) })
		}
		if n == 0 && st.filtering() {
			style.Text(ui.Text(c, locale.M().NoResults), c, theme.SM()).TextColor(t.Muted).
				Padding(theme.Scale(c, theme.Space(4))).TextAlign(ui.Center).Unselectable()
			return
		}
		queries := []string{st.Query}
		// Log text reads left to right in a right-to-left interface too.
		list := ui.List(c, &st.List, n, func(row int) {
			logLine(c, t, &o, st, queries, st.original(row), gutter)
		}).FillWidth().MinWidth(0).Margin(pad, 0).Direction(ui.LeftToRight).Label(label)
		if o.Height > 0 {
			list.Height(o.Height)
		} else {
			list.MaxHeight(logDefaultMaxHeight)
		}
		if o.Follow && n > 0 && !st.List.AtEnd() {
			// The outline face is translucent; the backing keeps the lines
			// under it from showing through.
			ui.Row(c).Absolute().Bottom(pad).End(theme.Scale(c, theme.Space(4))).
				Background(t.Code).Radius(theme.RadiusScale().LG).Transition(jumpFade()).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: locale.M().JumpToLatest, Icon: lucide.ArrowDown(), IconOnly: true,
					Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
					// The button goes once the end shows; the lines keep the focus.
					OnClick: func() {
						st.List.ScrollToEnd()
						st.List.Focus(c)
					},
				})
			})
		}
	})
}

// logGutterWidth is the room the widest line number takes, estimated from
// the monospace advance.
func logGutterWidth(c *ui.Context, n int) float32 {
	digits := max(2, len(strconv.Itoa(n)))
	return theme.Scale(c, codeFontSize) * 0.62 * float32(digits)
}

func logLine(c *ui.Context, t theme.Palette, o *LogViewerOptions, st *LogViewerState, queries []string, i int, gutter float32) {
	ts := theme.Code()
	ts.Size, ts.LineHeight = codeFontSize, codeLineHeight
	row := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).
		Gap(theme.Scale(c, theme.Space(3))).Padding(0, theme.Scale(c, theme.Space(4)))
	row.ContextMenu(func(m *ui.Menu) {
		msg := locale.M()
		copied := func(err error) {
			if err != nil && o.OnError != nil {
				o.OnError(err)
			}
		}
		if m.Item(msg.CopyLine).Chosen() {
			c.WriteClipboardThen(o.line(i), copied)
		}
		if m.Item(msg.CopyAll).Chosen() {
			c.WriteClipboardThen(st.listedText(o), copied)
		}
	})
	row.Children(func() {
		if o.LineNumbers {
			style.Text(ui.Text(c, strconv.Itoa(i+1)), c, ts).TextColor(t.Muted).
				Width(gutter).Shrink(0).TextAlign(ui.End).Unselectable()
		}
		var text ui.Element
		if st.Query != "" {
			text = style.Text(mygoui.Highlight(c, mygoui.HighlightOptions{
				Text: o.line(i), Queries: queries, IgnoreCase: true,
			}), c, ts).TextColor(t.CodeForeground)
		} else {
			text = style.Text(ui.Text(c, o.line(i)), c, ts)
		}
		text.Grow(1).Basis(0).MinWidth(0).Selectable()
		if st.NoWrap {
			text.SingleLine()
		}
		if o.Level != nil {
			if ink, _, ok := logLevelStyle(t, o.Level(i)); ok {
				text.TextColor(ink)
			}
		}
	})
}

// listedText is the lines the viewer lists, one to a row.
func (s *LogViewerState) listedText(o *LogViewerOptions) string {
	var b strings.Builder
	for row := range s.count(o.count()) {
		if row > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(o.line(s.original(row)))
	}
	return b.String()
}

// logLevelStyle returns the ink of a level and its name, which assistive
// technology reads with the line; Info and plain lines keep the block's text
// color and are not described.
func logLevelStyle(t theme.Palette, level LogLevel) (ui.Color, string, bool) {
	m := locale.M()
	switch level {
	case LogDebug:
		return t.Muted, m.LevelDebug, true
	case LogWarning:
		return t.WarningForeground, m.LevelWarning, true
	case LogError:
		return t.DestructiveForeground, m.LevelError, true
	}
	return ui.Color{}, "", false
}
