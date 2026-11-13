// Package layout provides the window-level building blocks of a native
// application: a custom title bar that keeps clear of the system window
// controls, a split view of side panes beside the content, and a scaffold
// composing both with a status strip.
//
// A window that draws its own title bar is created with
// mygo.WindowOptions.TitleBarStyle; TitleBar then reads Context.TitleBar
// for the room the system controls take. The toolkit, not this package,
// moves and maximizes the window from the bar.
//
// Pane sizes and the collapsed and open flags are caller-owned pointers.
// The package keeps nothing between windows, stores nothing and does no
// work on import. Elements and contexts are build-scoped.
package layout
