// Package navigation provides reusable native navigation and disclosure controls.
package navigation

import (
	"slices"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type AccordionMode uint8

const (
	AccordionSingle AccordionMode = iota
	AccordionMultiple
)

type AccordionState[T comparable] struct {
	Open     []T
	sections map[T]*accordionSection
}
type accordionSection struct {
	trigger, panel ui.Handle
	wasOpen        bool
}
type AccordionItem[T comparable] struct {
	Value    T
	Label    string
	Disabled bool
	Content  func()
}
type AccordionOptions struct {
	Mode                 AccordionMode
	Disabled, Borderless bool
	OnChange             func()
}

// Accordion retains open values through reordering. Header handles keep arrow
// navigation outside panel controls while native disclosures own activation.
func Accordion[T comparable](c *ui.Context, s *AccordionState[T], o AccordionOptions, items []AccordionItem[T]) ui.Element {
	p := theme.Colors(c)
	root := ui.Column(c).Radius(theme.Radii().XXL).Clip()
	if !o.Borderless {
		root.Border(theme.HairlineWidth, p.Border)
	}
	root.Children(func() {
		for i, item := range items {
			ui.Column(c.Key(item.Value)).Children(func() {
				at := slices.Index(s.Open, item.Value)
				open := at >= 0
				if o.Mode == AccordionSingle {
					open = at == 0
				}
				if s.sections == nil {
					s.sections = make(map[T]*accordionSection)
				}
				section := s.sections[item.Value]
				if section == nil {
					section = &accordionSection{}
					s.sections[item.Value] = section
				}
				if section.wasOpen && !open && section.panel.FocusWithin(c) {
					section.trigger.Focus(c)
				}
				parts := ui.CollapsibleBase(c.Key("disclosure"), &open)
				trigger := parts.Trigger.Bind(&section.trigger).Label(item.Label).Disabled(o.Disabled || item.Disabled).
					Justify(ui.Start).AlignItems(ui.Center).Gap(theme.Space12).Padding(theme.Space16).
					FontSize(theme.Scale(c, theme.Body().Size)).TextColor(p.Foreground)
				trigger.HandleInput(func(ev ui.InputEvent) bool {
					if ev.Kind != ui.InputKeyDown || ev.Mods != 0 {
						return false
					}
					target, step := i, 0
					switch ev.Key {
					case ui.KeyDown:
						step = 1
					case ui.KeyUp:
						step = -1
					case ui.KeyHome:
						target = -1
						step = 1
					case ui.KeyEnd:
						target = len(items)
						step = -1
					default:
						return false
					}
					for range len(items) {
						target = (target + step + len(items)) % len(items)
						if !items[target].Disabled && !o.Disabled {
							if next := s.sections[items[target].Value]; next != nil {
								next.trigger.Focus(c)
							}
							break
						}
					}
					return true
				})
				if i > 0 && !o.Borderless {
					trigger.BorderWidth(theme.HairlineWidth, 0, 0, 0).BorderColor(p.Border)
				}
				if open {
					trigger.Background(p.Muted)
				} else if trigger.Hovered() {
					trigger.Background(p.AccentHover)
				}
				trigger.Children(func() {
					ui.Text(c, item.Label).Grow(1)
					ui.Icon(c, lucide.ChevronDown()).Size(theme.IconMedium, theme.IconMedium).Rotate(180 * parts.Progress()).TextColor(p.MutedForeground)
				})
				trigger.OnChange(func() {
					if open {
						if o.Mode == AccordionSingle {
							s.Open = append(s.Open[:0], item.Value)
						} else if !slices.Contains(s.Open, item.Value) {
							s.Open = append(s.Open, item.Value)
						}
					} else if at := slices.Index(s.Open, item.Value); at >= 0 {
						s.Open = slices.Delete(s.Open, at, at+1)
					}
					if o.OnChange != nil {
						o.OnChange()
					}
				})
				parts.Panel(func() {
					// Outgoing content remains paintable during the native collapse, but can
					// no longer take focus or input after its controlling value closes.
					ui.Column(c).Bind(&section.panel).Disabled(!open || o.Disabled || item.Disabled).Padding(0, theme.Space16, theme.Space16).
						Children(func() {
							if item.Content != nil {
								item.Content()
							}
						})
				})
				section.wasOpen = open
			})
		}
	})
	return root
}

type CollapsibleState struct {
	Open           bool
	trigger, panel ui.Handle
	wasOpen        bool
}
type CollapsibleOptions struct {
	Label    string
	Disabled bool
	OnChange func(bool)
}

// Collapsible uses maintained disclosure progress, including ReduceMotion. An
// externally closed panel restores focus to its trigger before hiding children.
func Collapsible(c *ui.Context, s *CollapsibleState, o CollapsibleOptions, body func()) ui.Element {
	root := ui.Column(c)
	root.Children(func() {
		if s.wasOpen && !s.Open && s.panel.FocusWithin(c) {
			s.trigger.Focus(c)
		}
		parts := ui.CollapsibleBase(c.Key("disclosure"), &s.Open)
		p := theme.Colors(c)
		tr := parts.Trigger.Bind(&s.trigger).Label(o.Label).Disabled(o.Disabled).Justify(ui.Start).
			Padding(theme.Space8, theme.Space12).Gap(theme.Space8).Radius(theme.Radii().XXL).TextColor(p.Foreground)
		if tr.Hovered() {
			tr.Background(p.AccentHover)
		}
		tr.Children(func() {
			ui.Text(c, o.Label).Grow(1)
			ui.Icon(c, lucide.ChevronDown()).Size(theme.IconMedium, theme.IconMedium).Rotate(180 * parts.Progress())
		})
		tr.OnChange(func() {
			if o.OnChange != nil {
				o.OnChange(s.Open)
			}
		})
		parts.Panel(func() {
			ui.Column(c).Bind(&s.panel).Disabled(!s.Open || o.Disabled).Padding(theme.Space12).
				Children(func() {
					if body != nil {
						body()
					}
				})
		})
		s.wasOpen = s.Open
	})
	return root
}
