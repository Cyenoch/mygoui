// Package overlay implements HeroUI v3's transient surfaces on MyGo's native
// overlay, focus and toast engines.
//
// Every constructor takes a build-scoped *ui.Context and a caller-owned open
// flag. Overlays keep the focus while they show (ui.DialogBase marks the
// backdrop modal), restore it to the element that opened them, close on
// Escape, and paint HeroUI's overlay tokens live from package theme.
//
// Toasts are managed rather than drawn by the caller: Toast adds or updates
// one toast in a window, and ToastState.View installs the native toast region
// once per view with HeroUI's look. Cards form a spring-animated deck that
// expands on hover or focus and supports horizontal swipe dismissal. MyGo
// owns timers, pause on hover and focus, keyed updates and the three-toast cap.
//
//	toasts := new(overlay.ToastState)
//	// in the view, once:
//	toasts.View(c)
//	// in an event handler:
//	overlay.Toast(c, toasts, overlay.ToastOptions{Title: "Saved"})
//
// Gestures MyGo has no arbitration for are not invented: drawers and bottom
// sheets do not implement pointer swipe dismissal. Backdrop press and Escape
// remain the native dismissal routes; see each constructor's
// documentation.
package overlay
