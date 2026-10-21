package navigation

import (
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type TabsVariant uint8

const (
	TabsDefault TabsVariant = iota
	TabsLine
)

type Tab[T comparable] struct {
	Value    T
	Label    string
	Icon     *ui.SVG
	Disabled bool
	Content  func()
}
type TabsOptions[T comparable] struct {
	Label       string
	Variant     TabsVariant
	Orientation ui.Orientation
	OnChange    func(T)
}

// Tabs binds comparable identities rather than list indices. MyGo's maintained
// radio group owns arrow activation and disabled skipping; native roles describe
// the resulting tab list. Panels use RoleGroup because the pinned public toolkit
// has no TabPanel role. Hidden and disabled panels construct no controls.
func Tabs[T comparable](c *ui.Context, value *T, o TabsOptions[T], tabs []Tab[T]) ui.Element {
	p := theme.Colors(c)
	vertical := o.Orientation == ui.Vertical
	orientation := ui.Horizontal
	root := ui.Column(c).Gap(theme.Space8).MinWidth(0).FillWidth()
	if vertical {
		root.Row().AlignItems(ui.Stretch)
		orientation = ui.Vertical
	}
	// An absent/disabled selection follows the first enabled tab. Reordering keeps
	// the existing identity and never silently selects a different ordinal.
	valid := false
	for _, tab := range tabs {
		if tab.Value == *value && !tab.Disabled {
			valid = true
			break
		}
	}
	if !valid {
		for _, tab := range tabs {
			if !tab.Disabled {
				*value = tab.Value
				break
			}
		}
	}
	root.Children(func() {
		list := ui.RadioGroup(c.Key("list"), func() {
			for _, tab := range tabs {
				selected := *value == tab.Value && !tab.Disabled
				b := ui.RadioBase(c.Key(tab.Value), value, tab.Value).Role(ui.RoleTab).Label(tab.Label).Disabled(tab.Disabled).
					Justify(ui.Center).AlignItems(ui.Center).Padding(theme.Space4, theme.Space8).Gap(theme.Space6).
					MinHeight(theme.HeightButton).FontSize(theme.Scale(c, theme.Body().Size)).TextColor(p.MutedForeground)
				if selected {
					b.TextColor(p.Foreground)
					if o.Variant == TabsDefault {
						b.Background(p.Background)
					}
				} else if b.Hovered() {
					b.TextColor(p.Foreground)
				}
				if tab.Disabled {
					b.TextColor(p.DisabledForeground)
				}
				if o.Variant == TabsDefault {
					b.Radius(theme.Radii().XXL).Border(theme.HairlineWidth, ui.Color{})
					if selected {
						b.BorderColor(p.Input)
					}
				} else if vertical {
					b.BorderWidth(0, theme.HairlineWidth*2, 0, 0).BorderColor(ui.Color{})
				} else {
					b.BorderWidth(0, 0, theme.HairlineWidth*2, 0).BorderColor(ui.Color{})
				}
				if selected && o.Variant == TabsLine {
					b.BorderColor(p.Foreground)
				}
				if vertical {
					b.Justify(ui.Start)
				}
				if !theme.ReduceMotion(c) {
					b.Transition(ui.ElementTransition{Colors: true})
				}
				b.Children(func() {
					if tab.Icon != nil {
						ui.Icon(c, tab.Icon).Size(theme.IconMedium, theme.IconMedium).Shrink(0)
					}
					ui.Text(c, tab.Label).SingleLine()
				})
				b.OnChange(func() {
					if o.OnChange != nil {
						o.OnChange(*value)
					}
				})
			}
		}).Role(ui.RoleTabList).Label(o.Label).FocusGroup(orientation).Padding(theme.Space4).Gap(theme.Space4).
			AlignSelf(ui.Start).MinWidth(0).Shrink(0)
		// RadioGroup is a column by default. Its horizontal list must explicitly
		// opt out of the panel column's cross-axis stretch to keep an intrinsic strip.
		// Wrap constrains large labels without dividing the full panel width.
		if !vertical {
			list.Row().Wrap().Justify(ui.Start)
		}
		if o.Variant == TabsDefault {
			list.Radius(theme.Radii().XXL).Background(p.Muted)
		}
		panel := ui.Column(c.Key("panel")).Grow(1).MinWidth(0)
		if vertical {
			panel.Basis(0)
		} else {
			panel.FillWidth()
		}
		panel.Children(func() {
			for _, tab := range tabs {
				if *value == tab.Value && !tab.Disabled {
					ui.Column(c.Key(tab.Value)).Role(ui.RoleGroup).Label(tab.Label).Children(func() {
						if tab.Content != nil {
							tab.Content()
						}
					})
					break
				}
			}
		})
	})
	return root
}
