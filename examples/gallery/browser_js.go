package main

import (
	"strings"
	"syscall/js"
)

// appKeys are the letters the gallery takes with the primary modifier: text
// editing and its search.
const appKeys = "acxvzyk"

// keepBrowserShortcuts leaves the browser its own shortcuts, as reload and
// find. MyGo's surface cancels every key pressed with the primary modifier,
// so the others are kept from it.
func keepBrowserShortcuts() {
	window := js.Global()
	mac := strings.Contains(window.Get("navigator").Get("userAgent").String(), "Mac")
	window.Call("addEventListener", "keydown", js.FuncOf(func(_ js.Value, args []js.Value) any {
		e := args[0]
		primary := e.Get("ctrlKey").Bool()
		if mac {
			primary = e.Get("metaKey").Bool()
		}
		key := strings.ToLower(e.Get("key").String())
		if primary && len(key) == 1 && !strings.Contains(appKeys, key) {
			e.Call("stopPropagation")
		}
		return nil
	}), true)
}
