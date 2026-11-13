package content

import (
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// DiffLayout places the two texts of a diff.
type DiffLayout uint8

const (
	// DiffUnified shows the new text with deleted lines above their change.
	DiffUnified DiffLayout = iota
	// DiffSplit shows the old text beside the new one, lines aligned.
	DiffSplit
)

// DiffViewerOptions configures a read-only diff of Old against New. Its
// header names the file by Title, which also names the viewer, and shows
// the counts of added and deleted lines, a toggle between the two layouts (Layout is where it starts, and follows
// the option when that changes), a way to step through the changes while
// Height makes the viewer scroll (also F7 and Shift+F7) and a button
// copying New. ContextLines is how many unchanged lines stay around each
// change: 3 when zero, none when negative; the rest collapse into separators
// the user expands, unless ExpandUnchanged shows them all. Height makes the
// viewer scroll within that many DIPs; zero makes it as high as its lines.
// Identical texts show "No changes" in place of the lines. OnError receives
// a clipboard failure of the copy.
type DiffViewerOptions struct {
	Old, New        string
	Title           string
	Layout          DiffLayout
	ContextLines    int
	ExpandUnchanged bool
	NoLineNumbers   bool
	Wrap            bool
	Height          float32
	OnError         func(error)
}

// diffViewerState keeps the native viewer of one diff: its buffers, which
// are made again only when a text changes, the lines the user expanded and
// the scroll position, with what the header shows of the texts.
type diffViewerState struct {
	old, new          string
	view              editor.DiffState
	theme             editor.Theme
	added, deleted    int
	hunks             int
	hunk              int // the change last stepped to, -1 before any
	layout, optLayout DiffLayout
	ready             bool
}

type diffViewerStateKey struct{}

// DiffViewer shows what changed between two texts, unified or side by side,
// with line numbers, + and - before the changed lines, added and deleted
// lines tinted in the theme's success and danger roles and the words
// changed within them marked. The texts are selectable. It wraps MyGo's
// native diff viewer, which compares the texts when they change and lays
// out only the lines in view. A unified view numbers a deleted line by the
// old text and any other by the new, in one column.
func DiffViewer(c *ui.Context, o DiffViewerOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).AlignSelf(ui.Stretch).MinWidth(0).Clip().Radius(theme.RadiusScale().XL).
		Background(t.Code).TextColor(t.CodeForeground).
		BorderWidth(theme.Hairline).BorderColor(t.Border)
	if o.Title != "" {
		root.Role(ui.RoleGroup).Label(o.Title)
	}
	st := ui.Local(root, diffViewerStateKey{}, func() diffViewerState { return diffViewerState{} })
	d := &st.view
	if !st.ready || st.old != o.Old || st.new != o.New {
		st.old, st.new = o.Old, o.New
		d.Old, d.New = editor.NewBuffer(o.Old), editor.NewBuffer(o.New)
		st.added, st.deleted = d.Stats()
		st.hunks, st.hunk = len(d.Hunks()), -1
	}
	if !st.ready || st.optLayout != o.Layout {
		st.layout, st.optLayout = o.Layout, o.Layout
	}
	st.ready = true
	d.Layout = editor.DiffUnified
	if st.layout == DiffSplit {
		d.Layout = editor.DiffSplit
	}
	d.Indicators = editor.DiffClassic
	d.ContextLines, d.ExpandUnchanged = o.ContextLines, o.ExpandUnchanged
	d.NoLineNumbers, d.Wrap, d.FitContent = o.NoLineNumbers, o.Wrap, o.Height <= 0
	d.Font = ui.Font{Family: theme.MonoFontFamily(), Size: theme.Scale(c, codeFontSize)}
	d.LineHeight = float32(codeLineHeight) / codeFontSize
	d.TabSize = 4
	st.theme = editor.Theme{
		Background: t.Code, Text: t.CodeForeground, LineNumber: t.Muted,
		Added: t.Success, Deleted: t.Danger, Modified: t.Info,
	}
	d.Theme = &st.theme
	steps := o.Height > 0 && st.hunks > 0
	if steps {
		root.OnShortcut(0, ui.KeyF7, func() { st.step(1) })
		root.OnShortcut(ui.Shift, ui.KeyF7, func() { st.step(-1) })
	}
	return root.Children(func() {
		diffViewerHeader(c, t, st, &o, steps)
		if st.hunks == 0 && o.Old == o.New {
			style.Text(ui.Text(c, locale.M().NoChanges), c, theme.SM()).TextColor(t.Muted).
				Padding(theme.Scale(c, theme.Space(4))).Unselectable()
			return
		}
		view := editor.DiffView(c, d).FillWidth().MinWidth(0)
		if o.Height > 0 {
			view.Height(o.Height)
		}
	})
}

// step scrolls to the change d places after the one last shown, wrapping
// around.
func (s *diffViewerState) step(d int) {
	next := s.hunk + d
	switch {
	case s.hunk < 0 && d < 0:
		next = s.hunks - 1
	case s.hunk < 0:
		next = 0
	}
	s.hunk = (next + s.hunks) % s.hunks
	s.view.GoToHunk(s.hunk)
}

func diffViewerHeader(c *ui.Context, t theme.Palette, st *diffViewerState, o *DiffViewerOptions, steps bool) {
	m := locale.M()
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	count := theme.XS()
	count.Mono = true
	ui.Row(c).FillWidth().MinWidth(0).MinHeight(theme.Scale(c, 40)).
		AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(2))).
		Padding(0, theme.Scale(c, 10), 0, theme.Scale(c, theme.Space(4))).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).
		Children(func() {
			// The title, or the room it would take, keeps the rest at the end.
			style.Text(ui.Text(c, o.Title), c, ts).Grow(1).MinWidth(0).SingleLine().Unselectable()
			if st.hunks > 0 {
				// Counts read left to right, as the code does.
				ui.Row(c).Shrink(0).Gap(theme.Scale(c, theme.Space(2))).Direction(ui.LeftToRight).Children(func() {
					style.Text(ui.Text(c, "+"+strconv.Itoa(st.added)), c, count).
						TextColor(t.SuccessForeground).Unselectable()
					style.Text(ui.Text(c, "−"+strconv.Itoa(st.deleted)), c, count).
						TextColor(t.DestructiveForeground).Unselectable()
				})
			}
			// A button keeps its own alignment, so the row around them
			// centers them in the header.
			ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(1))).Children(func() {
				if steps {
					position := strconv.Itoa(st.hunks)
					if st.hunk >= 0 {
						position = strconv.Itoa(st.hunk+1) + " / " + position
					}
					style.Text(ui.Text(c.Key("position"), position), c, count).TextColor(t.Muted).Direction(ui.LeftToRight).
						Shrink(0).Unselectable().Role(ui.RoleStatus)
					for _, b := range []struct {
						label string
						icon  *ui.SVG
						by    int
					}{{m.Previous, lucide.ChevronUp(), -1}, {m.Next, lucide.ChevronDown(), 1}} {
						mygoui.Button(c.Key(b.label), mygoui.ButtonOptions{
							Label: b.label, Icon: b.icon, IconOnly: true, Variant: mygoui.ButtonGhost,
							Size: theme.ButtonSmall, OnClick: func() { st.step(b.by) },
						})
					}
				}
				split := st.layout == DiffSplit
				mygoui.ToggleButton(c.Key("layout"), mygoui.ToggleButtonOptions{
					Label: m.SideBySide, Icon: lucide.Columns2(), IconOnly: true, Variant: mygoui.ToggleOutline,
					Size: theme.ButtonSmall, Selected: &split,
					OnChange: func(on bool) {
						st.layout = DiffUnified
						if on {
							st.layout = DiffSplit
						}
					},
				})
				CopyButton(c.Key("copy"), CopyButtonOptions{Text: o.New, Label: copyLabel(o.Title), Size: theme.ButtonSmall, OnError: o.OnError})
			})
		})
}
