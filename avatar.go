package mygoui

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AvatarSize is HeroUI's avatar scale.
type AvatarSize uint8

const (
	AvatarMedium AvatarSize = iota
	AvatarSmall
	AvatarLarge
)

// AvatarOptions configures a caller-loaded avatar. Name is the accessible
// name; Fallback overrides generated initials.
type AvatarOptions struct {
	Image    ui.ImageSource
	Name     string
	Fallback string
	Size     AvatarSize
	Soft     bool
	Ring     bool
	Color    StatusColor
	Badge    func()
}

type avatarState struct {
	name, fallback, initials string
}

// Avatar shows an already-loaded image or initials. It never fetches.
func Avatar(c *ui.Context, o AvatarOptions) ui.Element {
	t := theme.Colors(c)
	size, radius, ts := float32(theme.HeightAvatarMD), theme.RadiusScale().XXXL, theme.SM()
	ts.Weight = theme.WeightMedium
	switch o.Size {
	case AvatarSmall:
		size, radius, ts = theme.HeightAvatarSM, theme.RadiusScale().XXL, theme.XS()
		ts.Weight = theme.WeightMedium
	case AvatarLarge:
		size, radius, ts = theme.HeightAvatarLG, theme.RadiusScale().XXXL, theme.Base()
		ts.Weight = theme.WeightMedium
	}
	size = theme.Scale(c, size)
	bg := t.Default
	if o.Soft {
		bg = StatusColors(c, o.Color).Soft
	}
	e := ui.Box(c).Size(size, size).Shrink(0).Radius(radius).
		Role(ui.RoleImage).Label(o.Name)
	if o.Ring {
		e.Shadow(0, 0, 0, 2, t.Background)
	}
	st := ui.Local(e, avatarStateKey{}, func() avatarState { return avatarState{} })
	if st.name != o.Name || st.fallback != o.Fallback {
		st.name, st.fallback = o.Name, o.Fallback
		st.initials = o.Fallback
		if st.initials == "" {
			st.initials = avatarInitials(o.Name)
		}
	}
	e.Children(func() {
		ui.Row(c).Size(size, size).Center().Radius(radius).Clip().Background(bg).
			Children(func() {
				if hasImage(o.Image) {
					ui.Image(c, o.Image).Size(size, size).Fit(ui.Cover).Role(ui.RoleNone)
					return
				}
				fallback := style.Text(ui.Text(c, st.initials), c, ts)
				fallback.TextColor(StatusColors(c, o.Color).Foreground).SingleLine().Role(ui.RoleNone)
			})
		if o.Badge != nil {
			ui.Box(c).Absolute().Right(0).Bottom(0).Children(o.Badge)
		}
	})
	return e
}

type avatarStateKey struct{}

func hasImage(source ui.ImageSource) bool {
	switch image := source.(type) {
	case *ui.Bitmap:
		if image != nil {
			w, h := image.Size()
			return w > 0 && h > 0
		}
	case *ui.SVG:
		if image != nil {
			w, h := image.Size()
			return w > 0 && h > 0
		}
	}
	return false
}

func avatarInitials(name string) string {
	fields := strings.Fields(name)
	if len(fields) == 0 {
		return "?"
	}
	initial := func(word string) string {
		r, _ := utf8.DecodeRuneInString(word)
		if r == utf8.RuneError {
			return ""
		}
		return string(unicode.ToUpper(r))
	}
	first := initial(fields[0])
	if first == "" {
		return "?"
	}
	if len(fields) == 1 {
		return first
	}
	return first + initial(fields[len(fields)-1])
}

// AvatarGroupOptions configures a HeroUI avatar stack or grid.
type AvatarGroupOptions struct {
	Size    AvatarSize
	Overlap float32
	Ring    bool
	Grid    bool
	Total   int
	Avatars []AvatarOptions
}

// AvatarGroup overlaps avatars by 8 DIPs and optionally separates them with a
// background-colored ring. A positive Total adds an overflow count.
func AvatarGroup(c *ui.Context, o AvatarGroupOptions) ui.Element {
	var e ui.Element
	if o.Grid {
		e = ui.Row(c).Wrap().Gap(theme.Space12).AlignItems(ui.Center)
	} else {
		e = ui.Row(c).Gap(0).AlignItems(ui.Center)
	}
	overlap := o.Overlap
	if overlap == 0 {
		overlap = theme.Space8
	}
	t := theme.Colors(c)
	e.Children(func() {
		for i := range o.Avatars {
			a := o.Avatars[i]
			if a.Size == 0 {
				a.Size = o.Size
			}
			if o.Ring && !o.Grid {
				a.Ring = true
			}
			avatar := Avatar(c, a)
			if i > 0 && !o.Grid {
				avatar.MarginStart(-overlap)
			}
			avatar.Role(ui.RoleImage)
		}
		if o.Total > len(o.Avatars) {
			count := ui.Row(c).Center().Shrink(0).
				MinWidth(theme.Scale(c, theme.HeightAvatarSM)).
				MinHeight(theme.Scale(c, theme.HeightAvatarSM)).
				PaddingX(theme.Space6).Radius(theme.RadiusScale().XXXL).
				Background(t.Default).TextColor(t.DefaultForeground)
			if len(o.Avatars) > 0 && !o.Grid {
				count.MarginStart(-overlap)
			}
			if o.Ring {
				count.Shadow(0, 0, 0, 2, t.Background)
			}
			count.Children(func() {
				style.Text(ui.Text(c, "+"+strconv.Itoa(o.Total-len(o.Avatars))), c, theme.XS()).SingleLine()
			})
		}
	})
	return e
}
