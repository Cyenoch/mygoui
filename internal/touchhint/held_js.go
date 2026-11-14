package touchhint

import "syscall/js"

// Held records, while held is true, that a control follows the pointer it
// is pressed with. The page reads and clears window.mygouiDragHeld.
func Held(held bool) {
	if held {
		js.Global().Set("mygouiDragHeld", true)
	}
}
