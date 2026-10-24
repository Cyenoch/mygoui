package navigation

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AccordionVariant selects HeroUI's accordion surface.
type AccordionVariant uint8

const (
	// AccordionDefault separates sections on the window background.
	AccordionDefault AccordionVariant = iota
	// AccordionSurface places the sections on a surface container.
	AccordionSurface
)

// AccordionOptions configures an accordion. Expanded is caller-owned; when it
// is nil the accordion keeps its expansion internally. Multiple allows more
// than one open section; the default is one.
type AccordionOptions struct {
	Variant       AccordionVariant
	Multiple      bool
	Expanded      *[]string
	HideSeparator bool
	Label         string
	// Items declares the sections in order.
	Items    []AccordionItemOptions
	OnChange func(key string, expanded bool)
}

// AccordionItemOptions configures one section. Key identifies it in Expanded
// and defaults to Title. Actions builds trailing header controls outside the
// disclosure toggle. Body builds the panel.
type AccordionItemOptions struct {
	Key           string
	Title         string
	Disabled      bool
	HideSeparator bool
	Actions       func()
	Body          func()
}

// Accordion builds a column of sections. Up and Down move the keyboard focus
// between headers, Home and End to the first and the last; Enter and Space
// toggle the focused section. The panel grows and its indicator turns with
// the shared collapsible animation, which reduced motion disables.
func Accordion(c *ui.Context, o AccordionOptions) ui.Element {
	t := theme.Colors(c)
	radius := min(scale(c, theme.Space32), theme.RadiusScale().XXXL)
	root := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Radius(radius).Clip().
		Role(ui.RoleGroup).FocusGroup(ui.Vertical)
	if o.Label != "" {
		root.Label(o.Label)
	}
	if o.Variant == AccordionSurface {
		root.Background(t.Surface).TextColor(t.SurfaceForeground)
	}
	expanded := o.Expanded
	if expanded == nil {
		expanded = ui.Local(root, "accordion-expanded", func() []string { return nil })
	}
	sc := &accordionScope{
		variant:       o.Variant,
		multiple:      o.Multiple,
		expanded:      expanded,
		hideSeparator: o.HideSeparator,
		onChange:      o.OnChange,
	}
	root.Children(func() {
		for i := range o.Items {
			accordionItem(c, sc, i, o.Items[i])
		}
	})
	return root
}

// accordionItem builds one section. A separator is drawn above every section
// but the first, matching upstream's bottom-aligned rule.
func accordionItem(c *ui.Context, sc *accordionScope, index int, o AccordionItemOptions) ui.Element {
	t := theme.Colors(c)
	key := o.Key
	if key == "" {
		key = o.Title
	}
	item := ui.Column(c).Shrink(0).MinWidth(0)
	if key != "" {
		item.Key(key)
	}
	expanded := sc.expanded
	open := false
	for _, k := range *expanded {
		if k == key {
			open = true
			break
		}
	}
	variant := sc.variant
	hide := o.HideSeparator || sc.hideSeparator
	item.Children(func() {
		if index > 0 && !hide {
			sep := ui.Box(c).Height(hairline).FillWidth().Shrink(0).Background(t.Separator)
			if variant == AccordionSurface {
				sep.Background(t.SurfaceForeground.Alpha(0.06)).MarginX(scale(c, theme.Space12))
			}
			sep.Role(ui.RoleNone)
		}
		var parts ui.CollapsibleParts
		var tr ui.Element
		headerRadius := theme.RadiusScale().XL
		header := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).
			PaddingX(scale(c, theme.Space16)).Radius(headerRadius).Clip()
		hover := t.Foreground.Alpha(0.03)
		if variant == AccordionSurface {
			hover = t.Default
		}
		header.Children(func() {
			parts = ui.CollapsibleBase(c, &open)
			tr = parts.Trigger.Grow(1).AlignItems(ui.Center).Gap(scale(c, theme.Space8)).
				PaddingY(scale(c, theme.Space16)).FocusRing(false).Label(o.Title)
			// Snapshot the trigger's states in the build: a retained copy
			// of a leaving element paints in later passes, where tr has
			// expired.
			trHovered, trPressed, trFocused := tr.Hovered(), tr.Pressed(), tr.FocusVisible()
			header.Draw(func(p *ui.Painter, r ui.Rect) {
				if open || o.Disabled {
					return
				}
				if trPressed || trHovered {
					p.Fill(r, hover, headerRadius)
				}
			})
			tr.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if trFocused {
					inset := [4]float32{theme.RadiusScale().MD, theme.RadiusScale().MD, theme.RadiusScale().MD, theme.RadiusScale().MD}
					p.FocusRing(ui.Rect{X: r.X + 4, Y: r.Y + 4, W: r.W - 8, H: r.H - 8}, inset)
				}
			})
			progress := parts.Progress()
			tr.Children(func() {
				style.Text(ui.Text(c, o.Title), c, theme.SM()).TextColor(t.Foreground).
					FontWeight(theme.WeightMedium).Grow(1).SingleLine()
				ind := ui.Box(c).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).Shrink(0)
				ind.Draw(func(p *ui.Painter, r ui.Rect) {
					drawChevronDown(p, r, 180*progress, 1.5, t.Muted)
				})
			})
			tr = style.Disabled(tr, o.Disabled)
			if o.Actions != nil {
				ui.Row(c).AlignItems(ui.Center).Shrink(0).Children(o.Actions)
			}
		})
		// Clicked applies the collapsible's value input, which toggles open
		// and marks the trigger changed, so the scope sees the new value in
		// this frame.
		_ = tr.Clicked()
		if tr.Changed() {
			sc.set(key, open)
			if sc.onChange != nil {
				sc.onChange(key, open)
			}
		}
		if panel := parts.Panel(o.Body); panel.Valid() {
			panel.FillWidth().MinWidth(0).Gap(scale(c, theme.Space8)).
				Padding(0, scale(c, theme.Space16), scale(c, theme.Space16), scale(c, theme.Space16)).
				TextColor(t.Muted)
		}
	})
	return item
}
