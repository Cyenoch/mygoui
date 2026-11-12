package mygoui

import "github.com/egoist/mygo"

// FlyoutHost lets anchored popups open as native flyouts of Window, so they
// can extend beyond it. Keep one per window and pass a pointer to it in a
// popup's options. A nil host, a nil Window or a platform without flyouts
// (WebAssembly) keeps the popup inside the window.
type FlyoutHost struct {
	Window *mygo.Window
}
