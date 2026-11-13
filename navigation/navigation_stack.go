package navigation

import (
	"strings"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// NavigationPage is one level of a NavigationStack. Key is its stable
// identity, Title names it in the header and Build builds its content.
type NavigationPage struct {
	Key   string
	Title string
	Build func()
}

// NavigationStackOptions configures a drill-down container. Pages is the
// caller-owned stack and is never modified: its last page is the one shown,
// the caller appends to go deeper and OnPop asks it to drop the last page.
// Without OnPop the stack offers no way back. Actions builds controls at the
// trailing end of the header, beside the page's title; it runs every build.
//
// KeepMounted keeps building the covered pages hidden, preserving their
// native scroll and editing state at the cost of building them every frame.
type NavigationStackOptions struct {
	Pages       []NavigationPage
	Label       string
	KeepMounted bool
	OnPop       func()
	Actions     func()
}

type (
	navStackStateKey  struct{}
	navPageEnteredKey struct{}
	navPageSlideKey   struct{}
)

// navStackState is what the last build showed, and how the shown page came:
// dir is 1 after going deeper and -1 after going back. leaving is the page a
// deeper one covers while it still slides away.
type navStackState struct {
	top     string
	depth   int
	dir     float32
	leaving string
	mounted bool
}

var navStackMotion = ui.Animation{Duration: 250 * time.Millisecond, Ease: ui.EaseOut}

// NavigationStack shows the last of Pages under a header with its title and
// a back button. A deeper page slides in over the one it covers and slides
// away again on the way back; under reduced motion pages change at once.
// The mouse's back button, Alt+Left (Alt+Right laid out right to
// left) and the back button call OnPop; Escape does not, as a page is not a
// modal and its editors keep it. When the page changes
// while the focus is inside the stack, the focus moves to the page shown.
//
// The page takes the height left under the header, so give the stack a
// height or let it grow where pages should scroll.
func NavigationStack(c *ui.Context, o NavigationStackOptions) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).MinHeight(0).Role(ui.RoleGroup)
	if o.Label != "" {
		root.Label(o.Label)
	}
	n := len(o.Pages)
	if n == 0 {
		return root
	}
	st := ui.Local(root, navStackStateKey{}, func() navStackState { return navStackState{} })
	top := o.Pages[n-1]
	changed := st.mounted && st.top != top.Key
	if changed {
		st.dir, st.leaving = -1, ""
		if n > st.depth {
			st.dir, st.leaving = 1, st.top
		}
	}
	focus := changed && root.FocusWithin()
	if changed && !focus {
		c.Announce(top.Title)
	}
	st.top, st.depth, st.mounted = top.Key, n, true

	// away is the side a covered page rests on, the inline start, and back
	// the arrow pointing there.
	away, back := float32(-100), ui.KeyLeft
	if c.Direction() == ui.RightToLeft {
		away, back = 100, ui.KeyRight
	}
	canPop := n > 1 && o.OnPop != nil
	if canPop {
		root.OnShortcut(ui.Alt, back, o.OnPop)
		root.OnShortcut(0, ui.KeyBack, o.OnPop)
	}
	root.Children(func() {
		parent := ""
		if canPop {
			parent = o.Pages[n-2].Title
		}
		navStackHeader(c, top.Title, parent, canPop, o.OnPop, o.Actions)
		stage := ui.Column(c).FillWidth().Grow(1).MinWidth(0).MinHeight(0).Clip()
		stage.Children(func() {
			for i := range o.Pages {
				page := &o.Pages[i]
				shown := i == n-1
				if !shown && !o.KeepMounted && page.Key != st.leaving {
					continue
				}
				el := navStackPage(c, st, page, shown, changed, away)
				if shown && focus {
					el.Focus()
				}
			}
		})
	})
	return root
}

// navStackHeader builds the header: a back button when there is a parent to
// return to, which names it in its tooltip, the title, and the actions.
func navStackHeader(c *ui.Context, title, parent string, canPop bool, pop, actions func()) {
	row := ui.Row(c).FillWidth().MinWidth(0).Shrink(0).AlignItems(ui.Center).
		Gap(scale(c, theme.Space4)).MinHeight(scale(c, theme.Space40))
	row.Children(func() {
		if canPop {
			icon := lucide.ChevronLeft()
			if c.Direction() == ui.RightToLeft {
				icon = lucide.ChevronRight()
			}
			// Buttons rest at the top of a row, so the button is centered
			// on the title.
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: locale.M().Back, Icon: icon, IconOnly: true,
				Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall, OnClick: pop,
				Tooltip: strings.TrimSuffix(locale.M().Back+": "+parent, ": "),
			}).AlignSelf(ui.Center)
		}
		style.Text(ui.Text(c, title), c, theme.SM()).FontWeight(theme.WeightSemibold).
			SingleLine().Grow(1).Basis(0).MinWidth(0).Role(ui.RoleHeading)
		if actions != nil {
			ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(scale(c, theme.Space4)).Children(actions)
		}
	})
}

// navStackPage builds one page at its place: the shown one in the flow, a
// covered one out of it at the inline start, hidden once it got there.
//
// Pages only translate, which the GPU does to what they paint: they need no
// bitmap of themselves (Cache), unlike a scaled panel with clips inside.
func navStackPage(c *ui.Context, st *navStackState, page *NavigationPage, shown, changed bool, away float32) ui.Element {
	el := ui.Column(c.Key(page.Key)).FillWidth().MinWidth(0).MinHeight(0).
		Role(ui.RoleGroup).Label(page.Title)
	entered := ui.Local(el, navPageEnteredKey{}, func() bool { return false })
	rest := away
	if shown {
		rest = 0
	}
	var x float32
	switch {
	case theme.ReduceMotion(c):
		x = el.Jump(navPageSlideKey{}, rest)
	case *entered:
		x = el.AnimateTo(navPageSlideKey{}, rest, navStackMotion)
	case shown && changed:
		// A page built for this change comes from the side it was on: the
		// inline end going deeper, the inline start going back.
		el.Jump(navPageSlideKey{}, -away*st.dir)
		x = el.AnimateTo(navPageSlideKey{}, rest, navStackMotion)
	default:
		x = el.Jump(navPageSlideKey{}, rest)
	}
	*entered = true
	if x != 0 {
		el.TranslatePercent(x, 0)
	}
	// Every built page keeps a transition, as MyGo copies an element that
	// had an Exit and stops declaring one. Only the shown page has the
	// Exit: once popped it is no longer the caller's to build, so an inert
	// copy of it slides back to the inline end.
	transition := ui.ElementTransition{
		Colors: true, Duration: navStackMotion.Duration, Ease: navStackMotion.Ease,
	}
	if shown {
		el.Grow(1).TabStop(false).FocusRing(false)
		transition.Exit = &ui.Motion{XPercent: -away, Opacity: 1}
	} else {
		el.Absolute().Left(0).Top(0).Right(0).Bottom(0).IgnoreInput(true).AccessHidden(true)
		if x == rest {
			el.Hidden(true)
			if st.leaving == page.Key {
				st.leaving = ""
			}
		}
	}
	if !theme.ReduceMotion(c) {
		el.Transition(transition)
	}
	el.Children(page.Build)
	return el
}
