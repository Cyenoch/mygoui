// Package overlay implements COSS-inspired transient surfaces on MyGo's
// native overlay, focus and toast engines.
//
// Modal, Dialog, Sheet, Drawer and BottomSheet use caller-owned open state
// and keep focus within the dialog. AlertDialog requires an explicit choice.
// Popover and PreviewCard are non-modal; Tooltip follows native hover and
// focus timing. Components read their colors live from package theme.
// Modal backdrops dim and blur the background by default; callers can select
// a dim-only or transparent backdrop.
//
// Drawer attaches to a window edge by default and supports an optional
// inset. Sheet slides 32 DIPs over 200 ms; Drawer slides from its edge over
// 450 ms. Their exits paint inert snapshots while focus returns immediately.
// Drawer dismisses with an outward swipe from the middle of its inward edge;
// ShowClose opts into a close button.
// BottomSheet keeps its floating geometry and dismisses from its handle swipe.
// A preview opens on hover or focus and allows 150 ms for the pointer to move
// between its anchor and panel.
//
// Toasts are managed rather than drawn by the caller: Toast adds or updates
// one toast in a window, and ToastState.View installs the native toast region
// once per view. Cards form a spring-animated deck that
// expands on hover or focus and supports horizontal swipe dismissal. MyGo
// owns timers, pause on hover and focus, keyed updates and the three-toast cap.
//
//	toasts := new(overlay.ToastState)
//	// in the view, once:
//	toasts.View(c)
//	// in an event handler:
//	overlay.Toast(c, toasts, overlay.ToastOptions{Title: "Saved"})
package overlay
