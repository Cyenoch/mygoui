package mygoui

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AvatarSize selects a native avatar size.
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

// Avatar shows an already-loaded image or initials. It never fetches.
func Avatar(c *ui.Context, o AvatarOptions) ui.Element {
	t := theme.Colors(c)
	size, radius, ts := avatarMetrics(c, o.Size)
	status := statusColors(&t, o.Color)
	bg := t.Border
	if o.Soft {
		bg = status.Soft
	}
	ink := status.SoftForeground
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
					// A newly loaded portrait fades in over the initials; an
					// image already present on the first build shows at once.
					ui.Image(c.Key(avatarImage{}), o.Image).Size(size, size).Fit(ui.Cover).Role(ui.RoleNone).
						Transition(fadeTransition())
					return
				}
				fallback := style.Text(ui.Text(c.Key(avatarInitialsKey{}), st.initials), c, ts)
				fallback.TextColor(ink).SingleLine().Role(ui.RoleNone).Transition(fadeTransition())
			})
		if o.Badge != nil {
			ui.Box(c).Absolute().Right(0).Bottom(0).Children(o.Badge)
		}
	})
	return e
}

// avatarMetrics resolves the portrait size, circular radius and fallback type
// role. The overflow count shares them.
func avatarMetrics(c *ui.Context, size AvatarSize) (float32, float32, theme.TextStyle) {
	dip := float32(theme.HeightAvatarMD)
	switch size {
	case AvatarSmall:
		dip = theme.HeightAvatarSM
	case AvatarLarge:
		dip = theme.HeightAvatarLG
	}
	actual := theme.Scale(c, dip)
	return actual, actual / 2, theme.TextStyle{Size: dip * 0.375, LineHeight: dip * 0.5, Weight: theme.WeightMedium}
}

type (
	avatarStateKey    struct{}
	avatarImage       struct{}
	avatarInitialsKey struct{}
)

// avatarState keeps the initials generated for the last name and fallback.
type avatarState struct {
	name, fallback, initials string
}

func hasImage(source ui.ImageSource) bool {
	w, h := ui.ImageSize(source)
	return w > 0 && h > 0
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

// AvatarGroupOptions configures an avatar stack or grid.
type AvatarGroupOptions struct {
	Size    AvatarSize
	Overlap float32
	Ring    bool
	Grid    bool
	Total   int
	Avatars []AvatarOptions
}

// AvatarGroup overlaps avatars by 30% of the group size and optionally adds
// background-colored rings. A Total above the displayed count adds overflow,
// named by the locale's MoreCount. It is a group; name it with Label.
func AvatarGroup(c *ui.Context, o AvatarGroupOptions) ui.Element {
	var e ui.Element
	if o.Grid {
		e = ui.Row(c).Wrap().Gap(theme.Space(3)).AlignItems(ui.Center)
	} else {
		e = ui.Row(c).Gap(0).AlignItems(ui.Center)
	}
	e.Role(ui.RoleGroup)
	overlap := o.Overlap
	if overlap == 0 {
		size, _, _ := avatarMetrics(c, o.Size)
		overlap = size * 0.3
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
			if !o.Grid {
				// Fills are translucent; an opaque base keeps overlaps from mixing.
				avatar.Background(t.Background)
				if i > 0 {
					avatar.MarginStart(-overlap)
				}
			}
			avatar.Role(ui.RoleImage)
		}
		if o.Total > len(o.Avatars) {
			dip, radius, ts := avatarMetrics(c, o.Size)
			more := o.Total - len(o.Avatars)
			count := ui.Row(c).Center().Shrink(0).
				Role(ui.RoleImage).Label(fmt.Sprintf(locale.M().MoreCount, more)).
				MinWidth(dip).MinHeight(dip).PaddingX(theme.Space(1.5)).Radius(radius).
				Background(t.Default.Over(t.Background)).TextColor(t.DefaultSoftForeground)
			if len(o.Avatars) > 0 && !o.Grid {
				count.MarginStart(-overlap)
			}
			if o.Ring && !o.Grid {
				count.Shadow(0, 0, 0, 2, t.Background)
			}
			count.Children(func() {
				style.Text(ui.Text(c, "+"+strconv.Itoa(more)), c, ts).
					SingleLine().FontFeatures("tnum").Role(ui.RoleNone)
			})
		}
	})
	return e
}
