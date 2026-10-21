package mygoui

import (
	"slices"

	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type ToggleVariant uint8

const (
	ToggleDefault ToggleVariant = iota
	ToggleOutline
)

type ToggleSize uint8

const (
	ToggleSizeDefault ToggleSize = iota
	ToggleSM
	ToggleLG
)

type ToggleOptions struct {
	Label    string
	Icon     *myui.SVG
	IconOnly bool
	Variant  ToggleVariant
	Size     ToggleSize
	Disabled bool
	Invalid  bool
	OnChange func(bool)
	Handle   *myui.Handle
}

// Toggle binds directly to persistent caller state; callbacks run once after
// native input, without requiring a callback merely to edit the bound value.
// Source: https://ui.shadcn.com/r/styles/base-rhea/toggle.json
func Toggle(c *myui.Context, value *bool, o ToggleOptions) myui.Element {
	e := myui.ToggleBase(c, value).Label(o.Label).Disabled(o.Disabled)
	if o.Handle != nil {
		e.Bind(o.Handle)
	}
	styleToggle(c, e, *value, o)
	if o.OnChange != nil {
		e.OnChange(func() { o.OnChange(*value) })
	}
	toggleContent(c, e, o)
	return e
}

func styleToggle(c *myui.Context, e myui.Element, selected bool, o ToggleOptions) {
	t := theme.Colors(c)
	height := float32(theme.HeightInput)
	switch o.Size {
	case ToggleSM:
		height = theme.HeightButton
	case ToggleLG:
		height = theme.HeightRowCompact
	}
	height = max(height, theme.Scale(c, theme.Body().Size*1.25)+theme.Space8)
	e.MinHeight(height).MinWidth(height).PaddingX(theme.Space10).Gap(theme.Space4).
		Radius(theme.Radii().XXL).Font(theme.FontFamily(c)).FontSize(theme.Scale(c, theme.Body().Size)).
		FontWeight(500).TextColor(t.Foreground)
	if selected {
		e.Background(t.Muted)
	}
	if o.Variant == ToggleOutline {
		e.Border(theme.HairlineWidth, t.Input)
	}
	if o.Disabled {
		e.TextColor(t.DisabledForeground).Background(t.DisabledSurface)
	} else if e.Pressed() {
		e.Background(t.AccentPressed)
	} else if e.Hovered() {
		e.Background(t.AccentHover)
	}
	if o.Invalid {
		e.Border(theme.HairlineWidth, t.Destructive).Error(o.Label)
	}
}

func toggleContent(c *myui.Context, e myui.Element, o ToggleOptions) {
	e.Children(func() {
		if o.Icon != nil {
			myui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconMedium), theme.Scale(c, theme.IconMedium))
		}
		if !o.IconOnly {
			myui.Text(c, o.Label).SingleLine()
		}
	})
}

// ToggleGroupState keeps selection by stable item value across reorder/hiding.
// Zero state has no selection. Single selection permits toggling the choice off.
// Set Values outside rendering for externally controlled selection.
type ToggleGroupState struct {
	Values []string
}

func (s *ToggleGroupState) Selected(value string) bool {
	return slices.Contains(s.Values, value)
}

func (s *ToggleGroupState) toggle(value string, mode ToggleSelectionMode) {
	at := slices.Index(s.Values, value)
	if mode == ToggleSelectionSingle {
		if at >= 0 {
			s.Values = s.Values[:0]
		} else {
			s.Values = append(s.Values[:0], value)
		}
		return
	}
	if at >= 0 {
		s.Values = slices.Delete(s.Values, at, at+1)
		return
	}
	s.Values = append(s.Values, value)
}

type ToggleSelectionMode uint8

const (
	ToggleSelectionSingle ToggleSelectionMode = iota
	ToggleSelectionMultiple
)

type ToggleGroupOptions struct {
	Label            string
	Selection        ToggleSelectionMode
	Vertical         bool
	Disabled, Joined bool
	Spacing          float32
	Variant          ToggleVariant
	Size             ToggleSize
	OnChange         func()
}

type ToggleGroupItem struct {
	Value, Label string
	Icon         *myui.SVG
	Disabled     bool
	Handle       *myui.Handle
}

// ToggleGroup uses MyGo's native roving focus. Arrows follow orientation and
// direction, skip disabled items, wrap, and Home/End choose an end. Selection
// changes only on activation; Tab leaves the group as a single stop.
// Source: https://ui.shadcn.com/r/styles/base-rhea/toggle-group.json
func ToggleGroup(c *myui.Context, state *ToggleGroupState, o ToggleGroupOptions, items []ToggleGroupItem) myui.Element {
	orientation := myui.Horizontal
	group := myui.Row(c).Role(myui.RoleGroup).Label(o.Label).AlignItems(myui.Stretch)
	if o.Vertical {
		group.Column()
		orientation = myui.Vertical
	}
	spacing := o.Spacing
	if spacing == 0 && !o.Joined {
		spacing = theme.Space8
	}
	group.Gap(spacing).FocusGroup(orientation).Children(func() {
		for i, item := range items {
			c.Key(item.Value)
			disabled := o.Disabled || item.Disabled
			e := myui.ButtonBase(c).Role(myui.RoleToggleButton).Label(item.Label).
				Checked(state.Selected(item.Value)).Disabled(disabled)
			if item.Handle != nil {
				e.Bind(item.Handle)
			}
			options := ToggleOptions{
				Label: item.Label, Icon: item.Icon, Variant: o.Variant,
				Size: o.Size, Disabled: disabled,
			}
			styleToggle(c, e, state.Selected(item.Value), options)
			if o.Joined {
				joinControl(e, i, len(items), o.Vertical)
			}
			e.OnClick(func() {
				state.toggle(item.Value, o.Selection)
				if o.OnChange != nil {
					o.OnChange()
				}
			})
			toggleContent(c, e, options)
		}
	})
	return group
}
