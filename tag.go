package mygoui

import (
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TagSelection is HeroUI's tag selection behavior.
type TagSelection uint8

const (
	TagSelectNone TagSelection = iota
	TagSelectSingle
	TagSelectMultiple
)

// TagSize is HeroUI's tag density.
type TagSize uint8

const (
	TagMedium TagSize = iota
	TagSmall
	TagLarge
)

// TagVariant is HeroUI's unselected tag surface.
type TagVariant uint8

const (
	TagDefault TagVariant = iota
	TagSurface
)

// TagOptions configures one tag.
type TagOptions struct {
	Key       string
	Label     string
	Icon      *ui.SVG
	Removable bool
	Disabled  bool
	OnRemove  func()
	OnClick   func()
}

// TagGroupState is caller-owned tag selection.
type TagGroupState struct {
	Selected map[string]bool
}

// IsSelected reports whether key is chosen.
func (s *TagGroupState) IsSelected(key string) bool {
	return s != nil && s.Selected != nil && s.Selected[key]
}

func (s *TagGroupState) toggle(key string) {
	if s == nil {
		return
	}
	if s.Selected == nil {
		s.Selected = map[string]bool{}
	}
	if s.Selected[key] {
		delete(s.Selected, key)
	} else {
		s.Selected[key] = true
	}
}

func (s *TagGroupState) clearExcept(key string) {
	if s == nil {
		return
	}
	s.Selected = map[string]bool{key: true}
}

// TagGroupOptions configures a focusable tag collection.
type TagGroupOptions struct {
	Label, Description, Error string
	Selection                 TagSelection
	Size                      TagSize
	Variant                   TagVariant
	State                     *TagGroupState
	OnChange                  func(keys []string)
}

// TagGroup lays out removable, optionally selectable tags. Tags stay smaller
// than buttons and expose native focus and removal actions.
func TagGroup(c *ui.Context, o TagGroupOptions, tags []TagOptions) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4).Label(o.Label).
		Description(o.Description)
	e.Children(func() {
		if o.Label != "" {
			style.Text(ui.Text(c, o.Label), c, theme.Label()).TextColor(t.Foreground)
		}
		ui.Row(c).FillWidth().Wrap().Gap(theme.Space6).Role(ui.RoleList).Label(o.Label).FocusGroup(ui.Horizontal).Children(func() {
			resolved := tagResolved{Size: o.Size, Variant: o.Variant}
			for i := range tags {
				tag := tags[i]
				r := resolved
				r.Selected = o.State.IsSelected(tag.Key)
				tagView(c, TagOptions{
					Key: tag.Key, Label: tag.Label, Icon: tag.Icon, Removable: tag.Removable,
					Disabled: tag.Disabled, OnRemove: tag.OnRemove,
					OnClick: func() {
						if o.Selection == TagSelectNone || o.State == nil {
							if tag.OnClick != nil {
								tag.OnClick()
							}
							return
						}
						if o.Selection == TagSelectSingle {
							o.State.clearExcept(tag.Key)
						} else {
							o.State.toggle(tag.Key)
						}
						if o.OnChange != nil {
							o.OnChange(tagKeys(tags, o.State))
						}
						if tag.OnClick != nil {
							tag.OnClick()
						}
					},
				}, r)
			}
		})
		if o.Error != "" {
			style.Text(ui.Text(c, o.Error), c, theme.HelpText()).
				TextColor(t.Danger).Role(ui.RoleStatus)
		}
	})
	return e
}

func tagKeys(tags []TagOptions, s *TagGroupState) []string {
	keys := make([]string, 0, len(tags))
	for i := range tags {
		if s.IsSelected(tags[i].Key) {
			keys = append(keys, tags[i].Key)
		}
	}
	return keys
}

// tagResolved carries the group settings a tag inherits.
type tagResolved struct {
	Selected bool
	Size     TagSize
	Variant  TagVariant
}

// tagMotionDuration is HeroUI's background-color transition for a tag. It
// doubles as the fade of a removed tag's exit copy.
const tagMotionDuration = 100 * time.Millisecond

// Tag renders one HeroUI tag. Most callers use TagGroup; Tag is available for
// a standalone label.
func Tag(c *ui.Context, o TagOptions) ui.Element {
	return tagView(c, o, tagResolved{Variant: TagDefault, Size: TagMedium})
}

// tagView is the actual builder used by TagGroup and standalone Tag. Selection
// and hover faces fade over HeroUI's 100 ms, and a tag that leaves the build
// fades and shrinks as an inert exit copy. A key keeps that exit attached to
// the tag the caller removed.
func tagView(c *ui.Context, o TagOptions, r tagResolved) ui.Element {
	t := theme.Colors(c)
	paddingX, paddingY := float32(theme.Space8), float32(theme.Space4)
	radius := theme.RadiusScale().XL
	ts := theme.XS()
	switch r.Size {
	case TagSmall:
		paddingY = theme.Space2
	case TagLarge:
		paddingX, paddingY = theme.Space10, theme.Space6
		radius = theme.RadiusScale().XXL
		ts = theme.SM()
	}
	bg, fg := t.Default, t.DefaultForeground
	if r.Variant == TagSurface {
		bg, fg = t.Surface, t.SurfaceForeground
	}
	hover := t.DefaultHover
	if r.Variant == TagSurface {
		hover = t.SurfaceHover
	}
	ts.Weight = theme.WeightMedium
	if r.Selected {
		bg, fg, hover = t.AccentSoft, t.AccentSoftForeground, t.AccentSoftHover
	}
	e := ui.Row(c).Shrink(0).Center().AlignSelf(ui.Start).Gap(theme.Space4).
		Padding(paddingY, paddingX, paddingY, paddingX).
		Radius(radius).Label(o.Label).
		Role(ui.RoleListItem).Selected(r.Selected)
	if o.Key != "" {
		e.Key(o.Key)
	}
	// Disable before the face reads the pointer, so a disabled tag never
	// paints HeroUI's hover face.
	e = style.Disabled(e, o.Disabled)
	e.Focusable().FocusRing(true).TextColor(fg)
	fill := bg
	if e.Pressed() || e.Hovered() {
		fill = hover
	}
	e.Background(fill)
	// Declare the transition in every build: omitting it reads as the tag
	// leaving and paints a fading copy over the live tag. Reduced motion
	// settles the frame instead, so the face and a removal show at once.
	if theme.ReduceMotion(c) {
		c.DisableAnimations()
	}
	e = e.Transition(ui.ElementTransition{
		Duration: tagMotionDuration, Ease: ui.EaseOut, Colors: true,
		Exit: &ui.Motion{Opacity: 0, Scale: 0.96},
	})
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick).Cursor(ui.CursorPointer)
	}
	if o.Removable && !o.Disabled {
		e.OnShortcut(0, ui.KeyDelete, func() {
			if o.OnRemove != nil {
				o.OnRemove()
			}
		})
		e.OnShortcut(0, ui.KeyBackspace, func() {
			if o.OnRemove != nil {
				o.OnRemove()
			}
		})
	}
	iconSize := theme.Scale(c, theme.IconXS)
	removeSize := theme.Scale(c, theme.Space12)
	e.Children(func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(iconSize, iconSize).Shrink(0)
		}
		style.Text(ui.Text(c, o.Label), c, ts).SingleLine()
		if o.Removable && !o.Disabled {
			remove := ui.ButtonBase(c).Size(removeSize, removeSize).
				Center().Padding(0).Radius(theme.RadiusScale().SM).FocusRing(true).
				Label(locale.M().Remove)
			remove.OnClick(func() {
				if o.OnRemove != nil {
					o.OnRemove()
				}
			})
			remove.Children(func() {
				ui.Icon(c, lucide.X()).Size(removeSize, removeSize)
			})
			remove.TextColor(fg)
		}
	})
	return e
}
