package overlay

import (
	"strconv"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ActionBarOptions configures a floating action bar.
type ActionBarOptions struct {
	// Label is the bar's text, as "selected" after a Count, and its
	// accessible name.
	Label string
	// Count, when positive, rolls in front of Label as the selection grows
	// or shrinks, and a change is announced.
	Count int
	// Actions are the bar's controls. Those that do not fit the width of the
	// container collapse from the end into a "more" menu; the count and the
	// close button never give way.
	Actions []navigation.OverflowItem
	// OnClose adds a close button that calls it, as does Escape while the
	// focus is in the bar; the caller clears what keeps the bar open.
	OnClose func()
	// ReturnFocus takes the focus when the bar closes with the focus inside
	// it, as after an action run from the keyboard; the bar is part of its
	// container, so nothing else gives that focus back.
	ReturnFocus *ui.Handle
}

type (
	actionBarKey   struct{}
	actionBarState struct{}
)

// actionBarMemo is what the bar remembers between builds, open or not: the
// count last shown, to announce a change, and whether the focus was inside.
type actionBarMemo struct {
	count int
	focus focusReturn
}

// actionBarRise is how far the bar travels as it appears and leaves.
const actionBarRise = 8

// ActionBar floats a toolbar at the bottom center of the container it is
// built in while open is true, which it covers without taking its input:
// build it last in a container that holds the content it floats over. It
// takes no focus and has no backdrop; the arrow keys move between its
// controls, which are one Tab stop. It rises in and out over 200 ms; under
// reduced motion it appears and leaves at once. The returned bar is absent
// while closed.
func ActionBar(c *ui.Context, open bool, o ActionBarOptions) ui.Element {
	frame := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
		Justify(ui.End).AlignItems(ui.Center).Padding(theme.Space(4)).PassThrough()
	memo := ui.Local(frame, actionBarState{}, func() actionBarMemo { return actionBarMemo{count: o.Count} })
	name := o.Label
	if o.Count > 0 {
		name = strconv.Itoa(o.Count) + " " + o.Label
	}
	changed := memo.count != o.Count
	memo.count = o.Count
	var bar ui.Element
	if !open {
		memo.focus.closed(func() {
			if o.ReturnFocus != nil {
				o.ReturnFocus.Focus(c)
			}
		})
		return bar
	}
	if changed && o.Count > 0 {
		c.Announce(name)
	}
	t := theme.Colors(c)
	frame.Children(func() {
		bar = ui.Row(c.Key(actionBarKey{})).MinWidth(0).Shrink(0).
			AlignItems(ui.Center).Gap(theme.Space(2)).
			PaddingY(theme.Space(2)).PaddingStart(theme.Space(4)).PaddingEnd(theme.Space(2)).
			Role(ui.RoleToolbar).FocusGroup(ui.Horizontal).Label(name)
		if o.OnClose != nil {
			bar.OnShortcut(0, ui.KeyEscape, o.OnClose)
		}
		// Reduced motion shows and hides it at once.
		if theme.ReduceMotion(c) {
			anchored.Still(bar)
		} else {
			bar.Transition(ui.ElementTransition{
				Duration: theme.Motion().Reveal, Ease: theme.Motion().Ease,
				Enter: &ui.Motion{Y: actionBarRise}, Exit: &ui.Motion{Y: actionBarRise},
			})
		}
		bar = anchored.Surface(c, bar, theme.RadiusScale().XL)
		memo.focus.note(bar)
		// The container is measured by the frame before.
		if w := frame.LayoutBounds().W; w > 0 {
			bar.MaxWidth(max(0, w-2*theme.Space(4)))
		}
		bar.Children(func() {
			ts := theme.SM()
			ts.Weight = theme.WeightMedium
			ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Space(1)).Children(func() {
				if o.Count > 0 {
					mygoui.NumberFlip(c, mygoui.NumberFlipOptions{Value: strconv.Itoa(o.Count), Style: ts})
				}
				if o.Label != "" {
					style.Text(ui.Text(c, o.Label), c, ts).SingleLine().Shrink(0)
				}
			})
			if len(o.Actions) > 0 {
				ui.Box(c).Width(theme.Hairline).Height(theme.Scale(c, theme.Space(5))).Shrink(0).
					Background(t.Separator).Role(ui.RoleNone)
				navigation.OverflowList(c, navigation.OverflowListOptions{Items: o.Actions}).Shrink(1)
			}
			if o.OnClose != nil {
				mygoui.CloseButton(c, mygoui.CloseButtonOptions{Label: locale.M().Close, OnClick: o.OnClose})
			}
		})
	})
	return bar
}
