package form

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Hotkey is a keyboard shortcut as MyGo reports and takes it, so it passes
// straight to OnShortcut. The zero value is no shortcut.
type Hotkey struct {
	Mods ui.Modifiers
	Key  ui.Key
}

// Keys returns the shortcut's key caps in the platform's order and naming:
// modifier glyphs on macOS, words elsewhere. It is nil for no shortcut.
func (h Hotkey) Keys() []string {
	name := hotkeyKeyName(h.Key)
	if name == "" {
		return nil
	}
	mac := runtime.GOOS == "darwin"
	var keys []string
	for _, m := range hotkeyModifiers {
		if h.Mods&m.mod == 0 {
			continue
		}
		switch {
		case mac:
			keys = append(keys, m.glyph)
		case m.mod == ui.Super && runtime.GOOS == "windows":
			keys = append(keys, "Win")
		default:
			keys = append(keys, m.word)
		}
	}
	return append(keys, name)
}

var hotkeyModifiers = [...]struct {
	mod         ui.Modifiers
	glyph, word string
}{
	{ui.Ctrl, "⌃", "Ctrl"},
	{ui.Alt, "⌥", "Alt"},
	{ui.Shift, "⇧", "Shift"},
	{ui.Super, "⌘", "Super"},
}

var hotkeyNamedKeys = map[ui.Key]string{
	ui.KeyEnter: "Enter", ui.KeyEscape: "Esc", ui.KeyBackspace: "Backspace",
	ui.KeyTab: "Tab", ui.KeySpace: "Space", ui.KeyDelete: "Delete",
	ui.KeyInsert: "Insert", ui.KeyHome: "Home", ui.KeyEnd: "End",
	ui.KeyPageUp: "PageUp", ui.KeyPageDown: "PageDown",
	ui.KeyLeft: "←", ui.KeyRight: "→", ui.KeyUp: "↑", ui.KeyDown: "↓",
	ui.KeyMinus: "-", ui.KeyEqual: "=", ui.KeyComma: ",", ui.KeyPeriod: ".",
	ui.KeySlash: "/", ui.KeySemicolon: ";", ui.KeyQuote: "'",
	ui.KeyBracketLeft: "[", ui.KeyBracketRight: "]", ui.KeyBackslash: "\\",
	ui.KeyBackquote: "`",
}

// hotkeyKeyName names a key for its key cap, empty when it has no name.
func hotkeyKeyName(key ui.Key) string {
	switch {
	case key >= ui.KeyA && key <= ui.KeyZ:
		return string(rune('A' + key - ui.KeyA))
	case key >= ui.Key0 && key <= ui.Key9:
		return string(rune('0' + key - ui.Key0))
	case key >= ui.KeyF1 && key <= ui.KeyF12:
		return "F" + strconv.Itoa(int(key-ui.KeyF1)+1)
	}
	return hotkeyNamedKeys[key]
}

// HotkeyInputOptions configures a shortcut recorder. OnChange runs after the
// shortcut is recorded or cleared.
type HotkeyInputOptions struct {
	Placeholder string
	AriaLabel   string
	Clearable   bool
	Disabled    bool
	FullWidth   bool
	OnChange    func()
}

type hotkeyRecordingKey struct{}

// HotkeyInput records a keyboard shortcut into *value. A click, Enter or
// Space starts recording; the next key pressed with a modifier is stored and
// shown as key caps. Escape or leaving the control stops recording, and
// Backspace or Delete clears a Clearable shortcut.
func HotkeyInput(c *ui.Context, value *Hotkey, o HotkeyInputOptions) ui.Element {
	t := theme.Colors(c)
	disabled := o.Disabled || value == nil
	if value == nil {
		value = &Hotkey{}
	}
	row := ui.Row(c).AlignItems(ui.Center).Shrink(0).MinHeight(theme.Scale(c, 32))
	if o.FullWidth {
		row.FillWidth()
	} else {
		row.Width(theme.Scale(c, inputDefaultWidth))
	}
	recording := ui.Local(row, hotkeyRecordingKey{}, func() bool { return false })
	set := func(h Hotkey) {
		*recording = false
		if *value == h {
			return
		}
		*value = h
		if o.OnChange != nil {
			o.OnChange()
		}
	}
	keys := value.Keys()

	ringed := false
	row.Children(func() {
		recorder := ui.ButtonBase(c).Grow(1).MinWidth(0).AlignSelf(ui.Stretch).Justify(ui.Start).
			PaddingX(theme.Space12).Cursor(ui.CursorPointer).FocusRing(false).
			Label(inputAccessibleName("", o.AriaLabel, o.Placeholder)).Value(strings.Join(keys, "+"))
		if !disabled {
			if !recorder.Focused() {
				*recording = false
			}
			recorder.OnClick(func() { *recording = !*recording })
			recorder.HandleInput(func(ev ui.InputEvent) bool {
				if ev.Kind != ui.InputKeyDown {
					return false
				}
				clears := ev.Mods == 0 && (ev.Key == ui.KeyBackspace || ev.Key == ui.KeyDelete)
				switch {
				case clears && o.Clearable:
					set(Hotkey{})
				case !*recording || ev.Key == ui.KeyTab:
					return false
				case ev.Mods == 0 && ev.Key == ui.KeyEscape:
					*recording = false
				case ev.Mods != 0 && hotkeyKeyName(ev.Key) != "":
					set(Hotkey{Mods: ev.Mods, Key: ev.Key})
				}
				return true
			})
		}
		// The ring marks recording, not mere focus: a focused recorder
		// still lets Tab, Enter and Space through.
		ringed = *recording || recorder.FocusVisible()
		recorder.Children(func() {
			if len(keys) == 0 {
				fieldText(ui.Text(c, o.Placeholder), c).TextColor(t.FieldPlaceholder).SingleLine().MinWidth(0)
				return
			}
			caps := mygoui.KbdGroup(c, mygoui.KbdOptions{Keys: keys}).AlignSelf(ui.Center)
			caps.Opacity(1 - 0.5*motion.Progress(recorder, "recording", *recording))
		})
		if o.Clearable && len(keys) > 0 && !disabled {
			hotkeyClear(c, func() { set(Hotkey{}) })
		}
	})
	return formstyle.SurfaceFocused(row, c, ringed, false, disabled)
}

// hotkeyClear is the trailing button that removes the shortcut.
func hotkeyClear(c *ui.Context, clear func()) {
	t := theme.Colors(c)
	b := ui.ButtonBase(c).Size(theme.Scale(c, theme.Space20), theme.Scale(c, theme.Space20)).
		Center().Shrink(0).Padding(0).MarginEnd(theme.Space8).
		Label(locale.M().Clear).Tooltip(locale.M().Clear).TextColor(t.FieldPlaceholder)
	b = style.Background(b, theme.RadiusScale().SM, ui.Transparent, t.Accent, t.Accent)
	b.OnClick(clear).Children(func() {
		ui.Icon(c, lucide.X()).Size(theme.Scale(c, theme.IconXS), theme.Scale(c, theme.IconXS)).Shrink(0)
	})
}
