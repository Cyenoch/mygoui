package mygoui

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type AvatarSize uint8

const (
	AvatarDefault AvatarSize = iota
	AvatarSM
	AvatarLG
)

type AvatarOptions struct {
	Image          myui.ImageSource
	Name, Fallback string
	Size           AvatarSize
	Badge          func()
}

type avatarFallbackKey struct{}

type avatarFallbackState struct {
	name, explicit, text string
}

// Avatar displays an already-loaded image or initials. Loading/decoding stays
// outside the view; nil/empty images immediately expose the fallback. Name is
// the image's accessible title even when initials differ from it.
// Source: https://ui.shadcn.com/r/styles/base-rhea/avatar.json
func Avatar(c *myui.Context, o AvatarOptions) myui.Element {
	size := float32(theme.Space32)
	switch o.Size {
	case AvatarSM:
		size = theme.Space24
	case AvatarLG:
		size = theme.Space20 * 2
	}
	size = theme.Scale(c, size)
	t := theme.Colors(c)
	e := myui.Box(c).Size(size, size).Shrink(0).Role(myui.RoleImage).Label(o.Name)
	fallback := myui.Local(e, avatarFallbackKey{}, func() avatarFallbackState { return avatarFallbackState{} })
	if fallback.name != o.Name || fallback.explicit != o.Fallback || fallback.text == "" {
		fallback.name, fallback.explicit, fallback.text = o.Name, o.Fallback, o.Fallback
		if fallback.text == "" {
			fallback.text = avatarInitials(o.Name)
		}
	}
	e.Children(func() {
		myui.Row(c).Size(size, size).Center().Radius(size).Clip().
			Background(t.Muted).Border(theme.HairlineWidth, t.Border).Children(func() {
			if avatarHasImage(o.Image) {
				myui.Image(c, o.Image).Size(size, size).Fit(myui.Cover).Role(myui.RoleNone)
				return
			}
			myui.Text(c, fallback.text).SingleLine().Font(theme.FontFamily(c)).
				FontSize(theme.Scale(c, theme.Small().Size)).TextColor(t.MutedForeground).Role(myui.RoleNone)
		})
		if o.Badge != nil {
			myui.Row(c).Absolute().Right(0).Bottom(0).Children(o.Badge)
		}
	})
	return e
}

func avatarHasImage(source myui.ImageSource) bool {
	switch image := source.(type) {
	case *myui.Bitmap:
		if image != nil {
			w, h := image.Size()
			return w > 0 && h > 0
		}
	case *myui.SVG:
		if image != nil {
			w, h := image.Size()
			return w > 0 && h > 0
		}
	}
	return false
}

func avatarInitials(name string) string {
	name = strings.TrimSpace(name)
	first, _ := utf8.DecodeRuneInString(name)
	if first == utf8.RuneError || name == "" {
		return "?"
	}
	last := first
	hasLast := false
	for i := len(name); i > 0; {
		r, n := utf8.DecodeLastRuneInString(name[:i])
		i -= n
		if unicode.IsSpace(r) {
			last, _ = utf8.DecodeRuneInString(strings.TrimLeftFunc(name[i+n:], unicode.IsSpace))
			hasLast = true
			break
		}
	}
	if !hasLast {
		return string(unicode.ToUpper(first))
	}
	return string([]rune{unicode.ToUpper(first), unicode.ToUpper(last)})
}

// AvatarGroup uses a native overlapping row; each image remains individually named.
func AvatarGroup(c *myui.Context, body func()) myui.Element {
	return myui.Row(c).Gap(-theme.Space8).Children(body)
}

func AvatarGroupCount(c *myui.Context, count int) myui.Element {
	label := "+" + strconv.Itoa(count)
	return Avatar(c, AvatarOptions{Name: locales.ComponentMorePeople(count), Fallback: label})
}

type AvatarBadgeOptions struct {
	Label string
	Icon  *myui.SVG
	Size  AvatarSize
}

// AvatarBadge is normally supplied through AvatarOptions.Badge, which positions it.
func AvatarBadge(c *myui.Context, o AvatarBadgeOptions) myui.Element {
	t := theme.Colors(c)
	size := float32(theme.Space10)
	switch o.Size {
	case AvatarSM:
		size = theme.Space8
	case AvatarLG:
		size = theme.Space12
	}
	size = theme.Scale(c, size)
	return myui.Row(c).Center().Size(size, size).Radius(size).Background(t.Primary).
		Border(2, t.Background).TextColor(t.PrimaryForeground).Label(o.Label).Children(func() {
		if o.Icon != nil && o.Size != AvatarSM {
			myui.Icon(c, o.Icon).Size(size-4, size-4)
		}
	})
}
