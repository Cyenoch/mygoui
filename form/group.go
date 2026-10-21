package form

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type InputGroupOptions struct {
	InputOptions
	Multiline                                    bool
	MinLines, MaxLines                           int
	InlineStart, InlineEnd, BlockStart, BlockEnd func()
}

// InputGroup composes native content around one editor. Action addons keep
// their own focus and callbacks; decorative addons never impersonate buttons.
func InputGroup(c *ui.Context, value *string, o InputGroupOptions) ui.Element {
	colors := theme.Colors(c)
	outer := ui.Column(c).FillWidth().MinWidth(0).Radius(theme.Radii().XXL).Background(colors.Input.Alpha(0.5))
	outer.Children(func() {
		// Build the editor first for native Field's control association.
		ui.Column(c).FillWidth().MinWidth(0).Children(func() {
			ui.Row(c).MinWidth(0).FillWidth().Children(func() {
				ui.Row(c).Reverse().Grow(1).MinWidth(0).Children(func() {
					var input ui.Element
					if o.Multiline {
						input = Textarea(c, value, TextareaOptions{
							InputOptions: o.InputOptions, MinLines: o.MinLines, MaxLines: o.MaxLines,
						})
					} else {
						input = Input(c, value, o.InputOptions)
					}
					input.Grow(1).Width(0).Radius(0).Background(colors.Input.Alpha(0)).Border(0, colors.Input)
					if o.InlineStart != nil {
						ui.Row(c).Padding(6, 8).Gap(8).Shrink(0).Children(o.InlineStart)
					}
				})
				if o.InlineEnd != nil {
					ui.Row(c).Padding(6, 8).Gap(8).Shrink(0).Children(o.InlineEnd)
				}
			})
			if o.BlockEnd != nil {
				ui.Row(c).Padding(8, 10).Gap(8).FillWidth().Children(o.BlockEnd)
			}
		})
		if o.BlockStart != nil {
			outer.Reverse()
			ui.Row(c).Padding(8, 10).Gap(8).FillWidth().Children(o.BlockStart)
		}
	})
	outer.Disabled(o.Disabled || o.Loading)
	inputBorder(c, outer, o.ControlOptions)
	return outer
}

func InputGroupText(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).FontSize(theme.Scale(c, 14)).TextColor(theme.Colors(c).MutedForeground)
}

// InputGroupButton retains the caller's native button variant and size while
// using the source's compact addon padding and radius. Choose ButtonGhost and
// ButtonXS explicitly for the source's default addon action.
func InputGroupButton(c *ui.Context, o mygoui.ButtonOptions) ui.Element {
	e := mygoui.Button(c, o)
	if o.Size == mygoui.ButtonXS {
		e.PaddingX(6).Radius(12)
	}
	return e
}
