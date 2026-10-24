package navigation

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DisclosureOptions configures a single disclosure. Open is caller-owned and
// persists the expansion; the disclosure keeps it internally when nil. Body
// builds the content shown while open. Actions builds optional trailing
// header controls outside the trigger.
type DisclosureOptions struct {
	Key      string
	Title    string
	Open     *bool
	Disabled bool
	Actions  func()
	Body     func()
}

// Disclosure builds a header trigger over a panel that grows and fades into
// view over the shared collapsible animation. Enter and Space toggle it while
// focused; reduced motion snaps it open and closed.
func Disclosure(c *ui.Context, o DisclosureOptions) ui.Element {
	return disclosure(c, o, nil)
}

func disclosure(c *ui.Context, o DisclosureOptions, group *disclosureGroupScope) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).FillWidth().MinWidth(0).Shrink(0)
	if o.Key != "" {
		root.Key(o.Key)
	}
	open := o.Open
	if open == nil {
		open = ui.Local(root, "disclosure-open", func() bool { return false })
	}
	root.Children(func() {
		var parts ui.CollapsibleParts
		var tr ui.Element
		glyph := &ui.Path{}
		header := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Shrink(0)
		header.Children(func() {
			parts = ui.CollapsibleBase(c, open)
			tr = parts.Trigger.Grow(1).AlignItems(ui.Center).Gap(scale(c, theme.Space8)).
				Padding(scale(c, theme.Space4), scale(c, theme.Space8)).
				Radius(theme.RadiusScale().MD).FocusRing(true).Label(o.Title)
			progress := parts.Progress()
			tr.Children(func() {
				style.Text(ui.Text(c, o.Title), c, theme.SM()).TextColor(t.Foreground).
					FontWeight(theme.WeightMedium).Grow(1).SingleLine()
				ind := ui.Box(c).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).Shrink(0)
				ind.Draw(func(p *ui.Painter, r ui.Rect) {
					drawChevronDown(p, glyph, r, 180*progress, 1.5, t.Foreground)
				})
			})
			tr = style.Disabled(tr, o.Disabled)
			if o.Actions != nil {
				ui.Row(c).AlignItems(ui.Center).Shrink(0).Children(o.Actions)
			}
		})
		_ = tr.Clicked()
		changed := tr.Changed()
		if group != nil {
			group.note(open, changed, o.Disabled)
		}
		if panel := parts.Panel(o.Body); panel.Valid() {
			panel.FillWidth().MinWidth(0).Padding(scale(c, theme.Space8)).TextColor(t.Muted)
			fadePanel(panel, parts.Progress())
		}
	})
	return root
}

// DisclosureGroupOptions configures a group of disclosures. Multiple allows
// more than one open at a time; each disclosure keeps its own caller-owned
// Open. Items declares the disclosures in order.
type DisclosureGroupOptions struct {
	Multiple bool
	Label    string
	Items    []DisclosureOptions
}

// DisclosureGroup builds a column of disclosures and, in single mode, closes
// the others when one opens. Expansion stays in the *bool each disclosure
// was given.
func DisclosureGroup(c *ui.Context, o DisclosureGroupOptions) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Role(ui.RoleGroup)
	if o.Label != "" {
		root.Label(o.Label)
	}
	sc := &disclosureGroupScope{multiple: o.Multiple}
	root.Children(func() {
		for i := range o.Items {
			disclosure(c, o.Items[i], sc)
		}
	})
	sc.reconcile()
	return root
}

// note records the state of one disclosure and, when it just opened in
// single mode, closes the siblings already built.
func (s *disclosureGroupScope) note(open *bool, changed, disabled bool) {
	at := -1
	for i := range s.items {
		if s.items[i].open == open {
			at = i
			break
		}
	}
	if at < 0 {
		s.items = append(s.items, disclosureItem{open: open})
		at = len(s.items) - 1
	}
	s.items[at].changed = changed
	if !changed || disabled || s.multiple || !*open {
		return
	}
	for j := range s.items {
		if j != at {
			*s.items[j].open = false
		}
	}
}

// reconcile keeps at most one open disclosure in single mode, preferring the
// one toggled in this build, then the last.
func (s *disclosureGroupScope) reconcile() {
	if s.multiple {
		return
	}
	keep := -1
	for i := range s.items {
		if *s.items[i].open {
			if s.items[i].changed || keep < 0 {
				keep = i
			}
		}
	}
	for i := range s.items {
		if i != keep {
			*s.items[i].open = false
		}
	}
}
